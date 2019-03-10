package cmd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bryk-io/did-method/client/store"
	"github.com/bryk-io/did-method/proto"
	"github.com/bryk-io/x/did"
	"github.com/bryk-io/x/net/rpc"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			name:      "key",
			usage:     "cryptographic key to use for the sync operation",
			flagKey:   "sync.key",
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

	// Get selected key for the sync operation
	key, err := getSyncKey(id)
	if err != nil {
		return err
	}
	fmt.Printf("Key selected for the operation: %s\n", key.ID)

	// Update proof
	fmt.Println("Updating record proof...")
	if err = id.AddProof(key.ID, didDomainValue); err != nil {
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
	ticket := proto.NewTicket(safe, key.ID)
	start := time.Now()
	challenge, err := ticket.Solve(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to generate request ticket: %s", err)
	}
	fmt.Printf("Ticket obtained: %s\n", challenge)
	fmt.Printf("Time: %s (rounds completed %d)\n", time.Since(start), ticket.Nonce())

	// Sign ticket
	ch, _ := hex.DecodeString(challenge)
	ticket.Signature, err = key.Sign(ch)
	if err != nil {
		return fmt.Errorf("failed to generate request ticket: %s", err)
	}

	// Verify on client's side
	if err = ticket.Verify(nil); err != nil {
		return fmt.Errorf("failed to verify ticket: %s", err)
	}

	// Submit request
	node := viper.GetString("node")
	fmt.Printf("Establishing connection to the network with node: %s\n", node)
	var opts []rpc.ClientOption
	opts = append(opts, rpc.WaitForReady())
	opts = append(opts, rpc.WithUserAgent("bryk-id-client"))
	opts = append(opts, rpc.WithTimeout(5 * time.Second))
	conn, err := rpc.NewClientConnection(node, opts...)
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

	// Update local record if sync was successful
	if res.Ok {
		contents, err := id.Encode()
		if err != nil {
			return fmt.Errorf("failed to encode identifier: %s", err)
		}
		return st.Update(name, contents)
	}
	return nil
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
