package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addServiceCmd = &cobra.Command{
	Use:     "service-add",
	Example: "bryk-id did service-add --did sample-record --name \"service name\" --endpoint https://www.agency.com/user_id",
	Short:   "Register a new service entry for the DID",
	RunE:    runAddServiceCmd,
}

func init() {
	didCmd.AddCommand(addServiceCmd)
}

func runAddServiceCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("register a new service")
	return nil
}
