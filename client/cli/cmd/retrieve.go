package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bryk-io/did-method/proto"
	"github.com/bryk-io/x/did"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var retrieveCmd = &cobra.Command{
	Use:     "retrieve",
	Short:   "Retrieve the DID document of an existing identifier",
	Example: "bryk-did retrieve --verify [existing DID]",
	Aliases: []string{"get", "resolve"},
	RunE:    runRetrieveCmd,
}

func init() {
	params := []cParam{
		{
			name:      "verify",
			usage:     "Verify the proofs included in the received DID Document",
			flagKey:   "retrieve.verify",
			byDefault: false,
		},
	}
	if err := setupCommandParams(retrieveCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(retrieveCmd)
}

func runRetrieveCmd(_ *cobra.Command, args []string) error {
	// Check params
	if len(args) != 1 {
		return errors.New("you must specify a DID to retrieve")
	}

	// Parse input
	id, err := did.Parse(args[0])
	if err != nil {
		return fmt.Errorf("the provided value is not a valid DID: %s", args[0])
	}

	// Validate method
	if id.Method() != "bryk" {
		return errors.New("only 'bryk' identifiers are supported")
	}

	// Get network connection
	ll := getLogger()
	conn, err := getClientConnection(ll)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Retrieve subject
	ll.Debug("retrieving record")
	client := proto.NewAgentClient(conn)
	res, err := client.Retrieve(context.TODO(), &proto.Request{Subject: id.Subject()})
	if err != nil {
		return fmt.Errorf("failed to retrieve DID records: %s", err)
	}
	if !res.Ok {
		return errors.New("no information available for the provided DID")
	}

	// Decode contents
	ll.Debug("decoding contents")
	peer := &did.Identifier{}
	if err = peer.Decode(res.Contents); err != nil {
		return fmt.Errorf("failed to decode DID records: %s", err)
	}
	js, err := json.MarshalIndent(peer.Document(), "", "  ")
	if err != nil {
		return fmt.Errorf("failed to decode DID records: %s", err)
	}

	// Verify DID Document
	if viper.GetBool("retrieve.verify") {
		ll.Info("verifying the received DID document")
		if err := peer.VerifyProof(nil); err != nil {
			return err
		}
		ll.Info("integrity proof is valid")
	}

	fmt.Printf("%s\n", js)
	return nil
}
