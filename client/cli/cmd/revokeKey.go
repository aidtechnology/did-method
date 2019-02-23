package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var revokeKeyCmd = &cobra.Command{
	Use:     "key-revoke",
	Example: "bryk-id did key-revoke --did sample-record --key-name iadb-account",
	Short:   "Revoke an existing cryptographic key for the DID",
	RunE:    runRevokeKeyCmd,
}

func init() {
	didCmd.AddCommand(revokeKeyCmd)
}

func runRevokeKeyCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("revoke an existing on a DID")
	return nil
}
