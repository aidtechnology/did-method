package cmd

import (
	"errors"
	"fmt"

	"github.com/bryk-io/x/did"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
)

var removeServiceCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an existing service entry for the DID",
	Example: "bryk-id did service remove [DID reference name] [service name]",
	RunE:    runRemoveServiceCmd,
}

func init() {
	serviceCmd.AddCommand(removeServiceCmd)
}

func runRemoveServiceCmd(_ *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("you must specify [DID reference name] [service name]")
	}

	// Get store handler
	st, err := getClientStore()
	if err != nil {
		return err
	}
	defer st.Close()

	// Get identifier
	ll := getLogger()
	name := sanitize.Name(args[0])
	ll.Info("removing existing service")
	ll.Debugf("retrieving entry with reference name: %s", name)
	e := st.Get(name)
	if e == nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
	}
	id := &did.Identifier{}
	if err = id.Decode(e.Contents); err != nil {
		return errors.New("failed to decode entry contents")
	}

	// Remove service
	sName := sanitize.Name(args[1])
	ll.Debugf("deleting service with name: %s", sName)
	if err = id.RemoveService(sName); err != nil {
		return fmt.Errorf("failed to remove service: %s", sName)
	}

	// Update record
	ll.Info("updating local record")
	contents, err := id.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode identifier: %s", err)
	}
	return st.Update(name, contents)
}
