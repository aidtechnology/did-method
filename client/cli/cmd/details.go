package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var didDetailsCmd = &cobra.Command{
	Use:     "details",
	Aliases: []string{"info"},
	Example: "bryk-id did details record-name",
	Short:   "Display the current information available on an existing DID",
	RunE:    runDidDetailsCmd,
}

func init() {
	didCmd.AddCommand(didDetailsCmd)
}

func runDidDetailsCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("show DID details")
	return nil
}
