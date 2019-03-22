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
	Example: "bryk-did retrieve --verify [existing DID]",
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

	// Retrieve subject
	ll := getLogger()
	ll.Info("retrieving record")
	response, err := resolver.Get(args[0])
	if err != nil {
		return err
	}

	// Verify DID Document
	if viper.GetBool("retrieve.verify") {
		ll.Info("verifying the received DID document")
		doc := &did.Document{}
		if err := json.Unmarshal(response, doc); err != nil {
			return fmt.Errorf("failed to decode received document: %s", err)
		}
		peer, err := did.FromDocument(doc)
		if err != nil {
			return fmt.Errorf("failed to restore DID document: %s", err)
		}
		if err := peer.VerifyProof(nil); err != nil {
			return err
		}
		ll.Info("integrity proof is valid")
	}

	// Print out received response
	fmt.Printf("%s\n", response)
	return nil
}
