package cmd

import (
	"errors"
	"fmt"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"go.bryk.io/x/did"
)

var removeKeyCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an existing cryptographic key for the DID",
	Example: "didctl did key remove [DID reference name] [key name]",
	RunE:    runRemoveKeyCmd,
}

func init() {
	keyCmd.AddCommand(removeKeyCmd)
}

func runRemoveKeyCmd(_ *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("you must specify [DID reference name] [key name]")
	}

	// Get store handler
	ll := getLogger()
	st, err := getClientStore()
	if err != nil {
		return err
	}
	defer func() {
		_ = st.Close()
	}()

	// Get identifier
	name := sanitize.Name(args[0])
	keyName := sanitize.Name(args[1])
	ll.Info("removing existing key")
	ll.Debugf("retrieving entry with reference name: %s", name)
	e := st.Get(name)
	if e == nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
	}
	id := &did.Identifier{}
	if err = id.Decode(e.Contents); err != nil {
		return errors.New("failed to decode entry contents")
	}

	// Remove key
	ll.Debug("validating parameters")
	if len(id.Keys()) >= 2 {
		_ = id.RemoveAuthenticationKey(keyName)
	}
	if err = id.RemoveKey(keyName); err != nil {
		return fmt.Errorf("failed to remove key: %s", keyName)
	}

	// Update record
	ll.Debug("key removed")
	ll.Info("updating local record")
	contents, err := id.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode identifier: %s", err)
	}
	return st.Update(name, contents)
}
