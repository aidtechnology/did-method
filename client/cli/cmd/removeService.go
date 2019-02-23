package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeServiceCmd = &cobra.Command{
	Use:     "service-remove",
	Example: "bryk-id did service-remove --did sample-record --name \"service name\"",
	Short:   "Remove an existing service entry for the DID",
	RunE:    runRemoveServiceCmd,
}

func init() {
	didCmd.AddCommand(removeServiceCmd)
}

func runRemoveServiceCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("remove an existing service")
	return nil
}
