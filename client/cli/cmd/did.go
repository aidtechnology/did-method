package cmd

import "github.com/spf13/cobra"

var didCmd = &cobra.Command{
	Use:   "did",
	Short: "Manage existing DIDs",
}

func init() {
	rootCmd.AddCommand(didCmd)
}
