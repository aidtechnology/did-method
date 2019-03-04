package cmd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/bryk-io/id/client/store"
	"github.com/bryk-io/id/proto"
	"github.com/bryk-io/x/did"
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

	// Generate request ticket
	fmt.Printf("Publishing: %s\n", name)
	fmt.Println("Generating request ticket...")
	ticket := &proto.Ticket{
		Timestamp:  time.Now().Unix(),
		Content:    record.Contents,
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

	// -> Submit request
	return nil
}
