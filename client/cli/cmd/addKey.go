package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addKeyCmd = &cobra.Command{
	Use:     "key-add",
	Example: "bryk-id did key-add --did sample-record --key-name iadb-account --key-type ed",
	Short:   "Add a new cryptographic key for the DID",
	RunE:    runAddKeyCmd,
}

func init() {
	didCmd.AddCommand(addKeyCmd)
}

func runAddKeyCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("add a new key to an existing DID")
	return nil
}
