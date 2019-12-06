package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"go.bryk.io/x/did"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List registered DIDs",
	Example: "didctl list",
	RunE:    runListCmd,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runListCmd(_ *cobra.Command, _ []string) error {
	// Get store handler
	st, err := getClientStore()
	if err != nil {
		return err
	}
	defer func() {
		_ = st.Close()
	}()

	// Get list of entries
	list := st.List()
	if len(list) == 0 {
		fmt.Println("no DIDs registered for the moment")
		return nil
	}

	// Show list of registered entries
	table := tabwriter.NewWriter(os.Stdout, 8, 0, 4, ' ', tabwriter.TabIndent)
	_, _ = fmt.Fprintf(table, "%s\t%s\t%s\n", "Reference Name", "Recovery Mode", "DID")
	for _, e := range list {
		id := &did.Identifier{}
		if err := id.Decode(e.Contents); err != nil {
			continue
		}
		_, _ = fmt.Fprintf(table, "%s\t%s\t%s\n", e.Name, e.Recovery, id.DID())
	}
	return table.Flush()
}
