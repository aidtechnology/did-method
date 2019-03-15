package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bryk-io/x/did"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/sha3"
)

var signCmd = &cobra.Command{
	Use:     "sign",
	Short:   "Produce a linked digital signature",
	Example: "bryk-did did sign [DID reference name] --input \"contents to sign\"",
	RunE:    runSignCmd,
}

func init() {
	params := []cParam{
		{
			name:      "input",
			usage:     "contents to sign, if longer than 32 bytes a SHA3-256 will be generated",
			flagKey:   "sign.input",
			byDefault: "",
		},
		{
			name:      "key",
			usage:     "key to use to produce the signature",
			flagKey:   "sign.key",
			byDefault: "master",
		},
		{
			name:      "domain",
			usage:     "domain value to use when producing LD signatures",
			flagKey:   "sign.domain",
			byDefault: didDomainValue,
		},
	}
	if err := setupCommandParams(signCmd, params); err != nil {
		panic(err)
	}
	keyCmd.AddCommand(signCmd)
}

func runSignCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
	}

	// Get input, CLI takes precedence, from standard input otherwise
	var input []byte
	input = []byte(viper.GetString("sign.input"))
	if len(input) == 0 {
		input, _ = getPipedInput()
	}
	if len(input) == 0 {
		return errors.New("no input passed in to sign")
	}
	if len(input) > 32 {
		digest := sha3.New256()
		digest.Write(input)
		input = digest.Sum(nil)
	}

	// Get store handler
	st, err := getClientStore()
	if err != nil {
		return err
	}
	defer st.Close()

	// Retrieve identifier
	name := sanitize.Name(args[0])
	record := st.Get(name)
	if record == nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
	}
	id := &did.Identifier{}
	if err = id.Decode(record.Contents); err != nil {
		return errors.New("failed to decode entry contents")
	}

	// Get key
	key := id.Key(viper.GetString("sign.key"))
	if key == nil {
		return fmt.Errorf("selected key is not available on the DID: %s", viper.GetString("sign.key"))
	}

	// Sign
	sld, err := key.ProduceSignatureLD(input, viper.GetString("sign.domain"))
	if err != nil {
		return fmt.Errorf("failed to produce signature: %s", err)
	}
	js, err := json.MarshalIndent(sld, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to produce signature: %s", err)
	}
	fmt.Printf("%s\n", js)
	return nil
}