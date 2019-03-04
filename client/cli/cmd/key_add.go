package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bryk-io/id/client/store"
	"github.com/bryk-io/x/did"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addKeyCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a new cryptographic key for the DID",
	Example: "bryk-id did key add [DID reference name] --name my-new-key --type ed --authentication",
	RunE:    runAddKeyCmd,
}

func init() {
	params := []cParam{
		{
			name:      "name",
			usage:     "name to be assigned to the newly added key",
			flagKey:   "key-add.name",
			byDefault: "key-#",
		},
		{
			name:      "type",
			usage:     "type of cryptographic key, either RSA (rsa) or Ed25519 (ed)",
			flagKey:   "key-add.type",
			byDefault: "ed",
		},
		{
			name:      "authentication",
			usage:     "enable this key for authentication purposes",
			flagKey:   "key-add.authentication",
			byDefault: false,
		},
	}
	if err := setupCommandParams(addKeyCmd, params); err != nil {
		log.Fatal(err)
	}
	keyCmd.AddCommand(addKeyCmd)
}

func runAddKeyCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
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

	// Set parameters
	keyName := viper.GetString("key-add.name")
	if strings.Count(keyName, "#") > 1 {
		return errors.New("invalid key name")
	}
	if strings.Count(keyName, "#") == 1 {
		keyName = strings.Replace(keyName, "#", fmt.Sprintf("%d", len(id.Keys()) + 1), 1)
	}
	keyName = sanitize.Name(keyName)
	keyType := did.KeyTypeEd
	keyEnc := did.EncodingHex
	if viper.GetString("key-add.type") == "rsa" {
		keyType = did.KeyTypeRSA
		keyEnc = did.EncodingBase64
	}

	// Add key
	if err = id.AddNewKey(keyName, keyType, keyEnc); err != nil {
		return fmt.Errorf("failed to add new key: %s", err)
	}
	if viper.GetBool("key-add.authentication") {
		if err = id.AddAuthenticationKey(keyName); err != nil {
			return fmt.Errorf("failed to establish key for authentication purposes: %s", err)
		}
	}
	if err = id.AddProof("master", didDomainValue); err != nil {
		return fmt.Errorf("failed to generate proof: %s", err)
	}

	// Update record
	contents, err := id.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode identifier: %s", err)
	}
	return st.Update(name, contents)
}
