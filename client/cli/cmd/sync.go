package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	protov1 "github.com/aidtechnology/did-method/proto/did/v1"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/pkg/cli"
	"go.bryk.io/pkg/did"
	"go.bryk.io/pkg/errors"
	xlog "go.bryk.io/pkg/log"
	"go.bryk.io/pkg/net/rpc"
	"go.bryk.io/pkg/otel"
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Publish a DID instance to the processing network",
	Example: "didctl sync [DID reference name]",
	Aliases: []string{"publish", "update", "upload", "push"},
	RunE:    runSyncCmd,
}

func init() {
	if err := cli.SetupCommandParams(syncCmd, conf.Overrides("client"), viper.GetViper()); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(syncCmd)
}

func runSyncCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
	}

	// Get store handler
	st, err := getClientStore()
	if err != nil {
		return err
	}

	// Retrieve identifier
	name := sanitize.Name(args[0])
	id, err := st.Get(name)
	if err != nil {
		return errors.Errorf("no available record for reference name: %s", name)
	}

	// Get selected key for the sync operation
	key, err := getSyncKey(id)
	if err != nil {
		return err
	}
	log.Debugf("key selected for the operation: %s", key.ID)

	// Generate request ticket
	log.Infof("publishing: %s", name)
	ticket, err := getRequestTicket(id, key)
	if err != nil {
		return err
	}

	// automatically instrument client if agent OTEL settings are available
	// in the same environment. Often the case for dev/testing.
	var clOpts []rpc.ClientOption
	if conf.Agent.OTEL != nil {
		oop, err := otel.NewOperator(conf.OTEL(log)...)
		if err == nil {
			clOpts = append(clOpts, rpc.WithClientObservability(oop))
			defer oop.Shutdown(context.Background())
		}
	}

	// Get client connection
	conn, err := getClientConnection(clOpts...)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	// Build request
	req := &protov1.ProcessRequest{
		Task:   protov1.ProcessRequest_TASK_PUBLISH,
		Ticket: ticket,
	}
	if viper.GetBool("sync.deactivate") {
		req.Task = protov1.ProcessRequest_TASK_DEACTIVATE
	}

	// Submit request
	log.Info("submitting request to the network")
	client := protov1.NewAgentAPIClient(conn)
	res, err := client.Process(context.Background(), req)
	if err != nil {
		return fmt.Errorf("network return an error: %w", err)
	}
	log.Debugf("request status: %v", res.Ok)
	if !res.Ok {
		return nil
	}

	// Update local record if sync was successful
	return st.Update(name, id)
}

func getRequestTicket(id *did.Identifier, key *did.VerificationKey) (*protov1.Ticket, error) {
	diff := uint(viper.GetInt("client.pow"))
	log.WithFields(xlog.Fields{"pow": diff}).Info("generating request ticket")

	// Create new ticket
	ticket, err := protov1.NewTicket(id, key.ID)
	if err != nil {
		return nil, err
	}

	// Solve PoW challenge
	start := time.Now()
	challenge := ticket.Solve(context.Background(), diff)
	log.Debugf("ticket obtained: %s", challenge)
	log.Debugf("time: %s (rounds completed %d)", time.Since(start), ticket.Nonce())
	ch, _ := hex.DecodeString(challenge)

	// Sign ticket
	if ticket.Signature, err = key.Sign(ch); err != nil {
		return nil, fmt.Errorf("failed to generate request ticket: %w", err)
	}

	// Verify on client's side
	if err = ticket.Verify(diff); err != nil {
		return nil, fmt.Errorf("failed to verify ticket: %w", err)
	}

	return ticket, nil
}

func getSyncKey(id *did.Identifier) (*did.VerificationKey, error) {
	// Get selected key for the sync operation
	key := id.VerificationMethod(viper.GetString("sync.key"))
	if key == nil {
		return nil, errors.New("invalid key selected")
	}

	// Verify the key is enabled for authentication
	isAuth := false
	for _, k := range id.GetVerificationRelationship(did.AuthenticationVM) {
		if k == key.ID {
			isAuth = true
			break
		}
	}
	if !isAuth {
		return nil, errors.New("key selected is not enabled for authentication purposes")
	}
	return key, nil
}
