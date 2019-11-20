package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bryk-io/did-method/resolver"
	"github.com/bryk-io/x/cli"
	"github.com/bryk-io/x/did"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var retrieveCmd = &cobra.Command{
	Use:     "retrieve",
	Short:   "Retrieve the DID document of an existing identifier",
	Example: "didctl retrieve --verify [existing DID]",
	Aliases: []string{"get", "resolve"},
	RunE:    runRetrieveCmd,
}

func init() {
	params := []cli.Param{
		{
			Name:      "verify",
			Usage:     "Verify the proofs included in the received DID Document",
			FlagKey:   "retrieve.verify",
			ByDefault: false,
		},
	}
	if err := cli.SetupCommandParams(retrieveCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(retrieveCmd)
}

func runRetrieveCmd(_ *cobra.Command, args []string) error {
	// Check params
	if len(args) != 1 {
		return errors.New("you must specify a DID to retrieve")
	}

	// Verify the provided value is a valid DID string
	id, err := did.Parse(args[0])
	if err != nil {
		return err
	}

	// Retrieve subject
	var sid *did.Identifier
	ll := getLogger()
	ll.Info("retrieving record")

	if id.Method() == "bryk" {
		// Use RPC client
		sid, err = retrieveSubject(id.Subject(), ll)
		if err != nil {
			return err
		}
	} else {
		// Use global resolver
		response, err := resolver.Get(args[0])
		if err != nil {
			return err
		}
		doc := &did.Document{}
		if err := json.Unmarshal(response, doc); err != nil {
			return fmt.Errorf("failed to decode received document: %s", err)
		}
		sid, err = did.FromDocument(doc)
		if err != nil {
			return err
		}
	}

	// Verify DID Document
	if viper.GetBool("retrieve.verify") {
		if err := sid.VerifyProof(nil); err != nil {
			return err
		}
		ll.Info("integrity proof is valid")
	}

	// Print out received response
	output, err := json.MarshalIndent(sid.Document(), "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", output)
	return nil
}
