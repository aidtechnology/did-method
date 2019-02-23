package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:     "register",
	Example: "bryk-id register",
	Short:   "Register a new DID with the network",
	RunE:    runRegisterCmd,
}

func init() {
	rootCmd.AddCommand(registerCmd)
}

func runRegisterCmd(cmd *cobra.Command, args []string) error {
	fmt.Println("register new DID")
	return nil
}
