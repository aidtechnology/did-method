package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Example: "bryk-id list",
	Short:   "List registered DIDs",
	RunE:    runListCmd,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runListCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("list registered DIDs")
	return nil
}
