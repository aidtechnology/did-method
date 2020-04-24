package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/x/ccg/did"
	"go.bryk.io/x/cli"
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
	_, err := did.Parse(args[0])
	if err != nil {
		return err
	}

	// Retrieve record
	log.Info("retrieving record")
	response, err := resolve(args[0])
	if err != nil {
		return fmt.Errorf("failed to resolve DID: %s", err)
	}

	// Parse response as a DID document. In case of error show warning
	// and print contents as-is.
	doc := &did.Document{}
	if err := json.Unmarshal(response, doc); err != nil {
		log.Warning("invalid DID document received")
		fmt.Printf("%s\n", response)
		return nil
	}

	// Get DID instance from the received document. In case of error show warning
	// and print contents as-is.
	sid, err := did.FromDocument(doc)
	if err != nil {
		log.Warning("invalid DID document received")
		fmt.Printf("%s\n", response)
		return nil
	}

	// Integrity checks
	if viper.GetBool("retrieve.verify") {
		if err := sid.VerifyProof(nil); err != nil {
			return err
		}
		log.Info("integrity proof is valid")
	}

	// Print out received response
	output, err := json.MarshalIndent(sid.SafeDocument(), "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", output)
	return nil
}
