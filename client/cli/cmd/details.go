package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bryk-io/x/did"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
)

var didDetailsCmd = &cobra.Command{
	Use:     "details",
	Short:   "Display the current information available on an existing DID",
	Example: "didctl did details [DID reference name]",
	Aliases: []string{"info"},
	RunE:    runDidDetailsCmd,
}

func init() {
	didCmd.AddCommand(didDetailsCmd)
}

func runDidDetailsCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
	}

	// Get store handler
	st, err := getClientStore()
	if err != nil {
		return err
	}
	defer func() {
		_ = st.Close()
	}()

	// Retrieve identifier
	name := sanitize.Name(args[0])
	e := st.Get(name)
	if e == nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
	}
	id := &did.Identifier{}
	if err = id.Decode(e.Contents); err != nil {
		return errors.New("failed to decode entry contents")
	}

	// Present its LD document as output
	info, _ := json.MarshalIndent(id.Document(), "", "  ")
	fmt.Printf("%s\n", info)
	return nil
}
