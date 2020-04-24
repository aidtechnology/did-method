package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/x/cli"
	"golang.org/x/crypto/sha3"
)

var signCmd = &cobra.Command{
	Use:     "sign",
	Short:   "Produce a linked digital signature",
	Example: "didctl sign [DID reference name] --input \"contents to sign\"",
	RunE:    runSignCmd,
}

func init() {
	params := []cli.Param{
		{
			Name:      "input",
			Usage:     "contents to sign, if longer than 32 bytes a SHA3-256 will be generated",
			FlagKey:   "sign.input",
			ByDefault: "",
			Short:     "i",
		},
		{
			Name:      "key",
			Usage:     "key to use to produce the signature",
			FlagKey:   "sign.key",
			ByDefault: "master",
			Short:     "k",
		},
		{
			Name:      "domain",
			Usage:     "domain value to use when producing LD signatures",
			FlagKey:   "sign.domain",
			ByDefault: didDomainValue,
			Short:     "d",
		},
	}
	if err := cli.SetupCommandParams(signCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(signCmd)
}

func runSignCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must specify a DID reference name")
	}

	// Get input, CLI takes precedence, from standard input otherwise
	input := []byte(viper.GetString("sign.input"))
	if len(input) == 0 {
		input, _ = cli.ReadPipedInput(maxPipeInputSize)
	}
	if len(input) == 0 {
		return errors.New("no input passed in to sign")
	}

	// Hash input
	hi := sha3.Sum256(input)
	input = hi[:]

	// Get store handler
	st, err := getClientStore()
	if err != nil {
		return err
	}

	// Retrieve identifier
	name := sanitize.Name(args[0])
	id, err := st.Get(name)
	if err != nil {
		return fmt.Errorf("no available record under the provided reference name: %s", name)
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
