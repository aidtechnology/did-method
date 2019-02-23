package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeKeyCmd = &cobra.Command{
	Use:     "key-remove",
	Example: "bryk-id did key-remove --did sample-record --key-name iadb-account",
	Short:   "Remove an existing cryptographic key for the DID",
	RunE:    runRemoveKeyCmd,
}

func init() {
	didCmd.AddCommand(removeKeyCmd)
}

func runRemoveKeyCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("remove an existing on a DID")
	return nil
}
