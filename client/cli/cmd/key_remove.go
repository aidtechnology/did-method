package cmd

import (
	"errors"
	"fmt"

	"github.com/bryk-io/did-method/client/store"
	"github.com/bryk-io/x/did"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeKeyCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an existing cryptographic key for the DID",
	Example: "bryk-id did key remove [DID reference name] [key name]",
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
	st, err := store.NewLocalStore(viper.GetString("home"))
	if err != nil {
		return err
	}
	defer st.Close()

	// Get identifier
	name := sanitize.Name(args[0])
	e := st.Get(name)
	if e == nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
	}
	id := &did.Identifier{}
	if err = id.Decode(e.Contents); err != nil {
		return errors.New("failed to decode entry contents")
	}

	// Remove key
	if len(id.Keys()) >= 2 {
		_ = id.RemoveAuthenticationKey(sanitize.Name(args[1]))
	}
	if err = id.RemoveKey(sanitize.Name(args[1])); err != nil {
		return fmt.Errorf("failed to remove key: %s", name)
	}

	// Update record
	contents, err := id.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode identifier: %s", err)
	}
	return st.Update(name, contents)
}
