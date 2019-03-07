package cmd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bryk-io/id/client/store"
	"github.com/bryk-io/id/proto"
	"github.com/bryk-io/x/did"
	"github.com/bryk-io/x/net/rpc"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	e "golang.org/x/crypto/ed25519"
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Publish a DID instance to the processing network",
	Example: "bryk-id sync [DID reference name]",
	RunE:    runSyncCmd,
}

func init() {
	params := []cParam{
		{
			name:      "proof-key",
			usage:     "name of the key to use to generate the identifier proof",
			flagKey:   "sync.proof",
			byDefault: "master",
		},
	}
	if err := setupCommandParams(syncCmd, params); err != nil {
		log.Fatal(err)
	}
	rootCmd.AddCommand(syncCmd)
}

func runSyncCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
	}

	// Get store handler
	st, err := store.NewLocalStore(viper.GetString("home"))
	if err != nil {
		return err
	}
	defer st.Close()

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

	// Update proof
	if err = id.AddProof(viper.GetString("sync.proof"), didDomainValue); err != nil {
		return fmt.Errorf("failed to generate proof: %s", err)
	}

	// Get safe contents to synchronize with the network
	safe, err := id.SafeEncode()
	if err != nil {
		return fmt.Errorf("failed to safely export identifier instance: %s", err)
	}

	// Generate request ticket
	fmt.Printf("Publishing: %s\n", name)
	fmt.Println("Generating request ticket...")
	ticket := &proto.Ticket{
		Timestamp:  time.Now().Unix(),
		Content:    safe,
		NonceValue: 0,
	}
	start := time.Now()
	challenge, err := ticket.Solve(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to generate request ticket: %s", err)
	}
	fmt.Printf("Ticket obtained: %s\n", challenge)
	fmt.Printf("Time: %s (rounds completed %d)\n", time.Since(start), ticket.Nonce())

	// Sign ticket
	key := id.Key("master")
	if key == nil {
		return errors.New("no master key set for the DID")
	}
	pvt := e.PrivateKey(key.Private)
	ch, _ := hex.DecodeString(challenge)
	ticket.Signature = e.Sign(pvt, ch)

	// Verify on client's side
	if err = ticket.Verify(); err != nil {
		return fmt.Errorf("failed to verify ticket: %s", err)
	}

	// Submit request
	fmt.Println("Establishing connection to the network")
	var opts []rpc.ClientOption
	opts = append(opts, rpc.WaitForReady())
	opts = append(opts, rpc.WithUserAgent("bryk-id-client"))
	conn, err := rpc.NewClientConnection(viper.GetString("node"), opts...)
	if err != nil {
		return fmt.Errorf("failed to establish connection: %s", err)
	}
	defer conn.Close()

	// Submit request
	fmt.Println("Submitting request to the network")
	client := proto.NewMethodClient(conn)
	res, err := client.Process(context.TODO(), ticket)
	if err != nil {
		return fmt.Errorf("network return an error: %s", err)
	}
	fmt.Printf("Final request status: %v\n", res.Ok)
	return nil
}
