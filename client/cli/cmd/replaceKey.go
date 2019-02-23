package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var replaceKeyCmd = &cobra.Command{
	Use:     "key-replace",
	Example: "bryk-id did key-replace --did sample-record --key-name iadb-account",
	Short:   "Replace an existing cryptographic key for the DID",
	RunE:    runReplaceKeyCmd,
}

func init() {
	didCmd.AddCommand(replaceKeyCmd)
}

func runReplaceKeyCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("replace an existing on a DID")
	return nil
}
