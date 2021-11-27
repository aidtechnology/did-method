package cmd

import (
	"errors"
	"fmt"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Permanently delete a local identifier",
	Example: "didctl delete [DID reference name]",
	Aliases: []string{"del", "rm", "remove"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("you must specify a DID reference name")
		}

		// Get store handler
		st, err := getClientStore()
		if err != nil {
			return err
		}

		// Delete identifier
		name := sanitize.Name(args[0])
		if err = st.Delete(name); err != nil {
			return fmt.Errorf("failed to remove entry: %w", err)
		}
		log.Infof("identifier successfully deleted: %s", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
