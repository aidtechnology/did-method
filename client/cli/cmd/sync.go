package cmd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	didpb "github.com/bryk-io/did-method/proto"
	"github.com/kennygrant/sanitize"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/x/cli"
	"go.bryk.io/x/did"
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Publish a DID instance to the processing network",
	Example: "didctl sync [DID reference name]",
	Aliases: []string{"publish", "update", "upload"},
	RunE:    runSyncCmd,
}

func init() {
	params := []cli.Param{
		{
			Name:      "key",
			Usage:     "cryptographic key to use for the sync operation",
			FlagKey:   "sync.key",
			ByDefault: "master",
		},
		{
			Name:      "deactivate",
			Usage:     "instruct the network agent to deactivate the identifier",
			FlagKey:   "sync.deactivate",
			ByDefault: false,
		},
		{
			Name:      "pow",
			Usage:     "set the required request ticket difficulty level",
			FlagKey:   "sync.pow",
			ByDefault: 24,
		},
	}
	if err := cli.SetupCommandParams(syncCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(syncCmd)
}

func runSyncCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
	}

	// Get store handler
	ll := getLogger()
	st, err := getClientStore()
	if err != nil {
		return err
	}
	defer func() {
		_ = st.Close()
	}()

	// Retrieve identifier
	name := sanitize.Name(args[0])
	record := st.Get(name)
	if record == nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
	}
	id := &did.Identifier{}
	if err = id.Decode(record.Contents); err != nil {
		return errors.New("failed to decode entry contents")
	}

	// Get selected key for the sync operation
	key, err := getSyncKey(id)
	if err != nil {
		return err
	}
	ll.Debugf("key selected for the operation: %s", key.ID)

	// Update proof
	ll.Info("updating record proof")
	if err = id.AddProof(key.ID, didDomainValue); err != nil {
		return fmt.Errorf("failed to generate proof: %s", err)
	}

	// Generate request ticket
	ll.Infof("publishing: %s", name)
	ticket, err := getRequestTicket(id, key, ll)
	if err != nil {
		return err
	}

	// Get client connection
	conn, err := getClientConnection(ll)
	if err != nil {
		return fmt.Errorf("failed to establish connection: %s", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// Build request
	req := &didpb.Request{
		Task:   didpb.Request_PUBLISH,
		Ticket: ticket,
	}
	if viper.GetBool("sync.deactivate") {
		req.Task = didpb.Request_DEACTIVATE
	}

	// Submit request
	ll.Info("submitting request to the network")
	client := didpb.NewAgentAPIClient(conn)
	res, err := client.Process(context.TODO(), req)
	if err != nil {
		return fmt.Errorf("network return an error: %s", err)
	}
	ll.Debugf("request status: %v", res.Ok)
	if !res.Ok {
		return nil
	}

	// Update local record if sync was successful
	contents, err := id.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode identifier: %s", err)
	}
	return st.Update(name, contents)
}

func getRequestTicket(id *did.Identifier, key *did.PublicKey, ll *log.Logger) (*didpb.Ticket, error) {
	ll.Info("generating request ticket")
	diff := uint(viper.GetInt("sync.pow"))
	ticket := didpb.NewTicket(id, key.ID)
	start := time.Now()
	challenge := ticket.Solve(context.TODO(), diff)
	ll.Debugf("ticket obtained: %s", challenge)
	ll.Debugf("time: %s (rounds completed %d)", time.Since(start), ticket.Nonce())
	ch, _ := hex.DecodeString(challenge)

	// Sign ticket
	var err error
	if ticket.Signature, err = key.Sign(ch); err != nil {
		return nil, fmt.Errorf("failed to generate request ticket: %s", err)
	}

	// Verify on client's side
	if err = ticket.Verify(nil, diff); err != nil {
		return nil, fmt.Errorf("failed to verify ticket: %s", err)
	}

	return ticket, nil
}

func getSyncKey(id *did.Identifier) (*did.PublicKey, error) {
	// Get selected key for the sync operation
	key := id.Key(viper.GetString("sync.key"))
	if key == nil {
		return nil, errors.New("invalid key selected")
	}

	// Verify the key is enabled for authentication
	isAuth := false
	for _, k := range id.AuthenticationKeys() {
		if k == key.ID {
			isAuth = true
			break
		}
	}
	if !isAuth {
		return nil, errors.New("the key selected is not enabled for authentication purposes")
	}
	return key, nil
}
