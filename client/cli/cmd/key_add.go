package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bryk-io/x/cli"
	"github.com/bryk-io/x/did"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addKeyCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a new cryptographic key for the DID",
	Example: "didctl did key add [DID reference name] --name my-new-key --type ed --authentication",
	RunE:    runAddKeyCmd,
}

func init() {
	params := []cli.Param{
		{
			Name:      "name",
			Usage:     "name to be assigned to the newly added key",
			FlagKey:   "key-add.name",
			ByDefault: "key-#",
		},
		{
			Name:      "type",
			Usage:     "type of cryptographic key, either RSA (rsa) or Ed25519 (ed)",
			FlagKey:   "key-add.type",
			ByDefault: "ed",
		},
		{
			Name:      "authentication",
			Usage:     "enable this key for authentication purposes",
			FlagKey:   "key-add.authentication",
			ByDefault: false,
		},
	}
	if err := cli.SetupCommandParams(addKeyCmd, params); err != nil {
		panic(err)
	}
	keyCmd.AddCommand(addKeyCmd)
}

func runAddKeyCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
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
	ll.Info("adding new key")
	ll.Debugf("retrieving entry with reference name: %s", name)
	e := st.Get(name)
	if e == nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
	}
	id := &did.Identifier{}
	if err = id.Decode(e.Contents); err != nil {
		return errors.New("failed to decode entry contents")
	}

	// Set parameters
	ll.Debug("validating parameters")
	keyName := viper.GetString("key-add.name")
	if strings.Count(keyName, "#") > 1 {
		return errors.New("invalid key name")
	}
	if strings.Count(keyName, "#") == 1 {
		keyName = strings.Replace(keyName, "#", fmt.Sprintf("%d", len(id.Keys())+1), 1)
	}
	keyName = sanitize.Name(keyName)
	keyType := did.KeyTypeEd
	keyEnc := did.EncodingHex
	if viper.GetString("key-add.type") == "rsa" {
		keyType = did.KeyTypeRSA
		keyEnc = did.EncodingBase64
	}

	// Add key
	ll.Debugf("adding new key with name: %s", keyName)
	if err = id.AddNewKey(keyName, keyType, keyEnc); err != nil {
		return fmt.Errorf("failed to add new key: %s", err)
	}
	if viper.GetBool("key-add.authentication") {
		ll.Info("setting new key as authentication mechanism")
		if err = id.AddAuthenticationKey(keyName); err != nil {
			return fmt.Errorf("failed to establish key for authentication purposes: %s", err)
		}
	}

	// Update record
	ll.Info("updating local record")
	contents, err := id.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode identifier: %s", err)
	}
	return st.Update(name, contents)
}
