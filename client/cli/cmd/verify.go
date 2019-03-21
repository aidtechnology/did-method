package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/bryk-io/x/cli"
	"github.com/bryk-io/x/did"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/sha3"
)

var verifyCmd = &cobra.Command{
	Use:     "verify",
	Short:   "Check the validity of a SignatureLD document",
	Example: "bryk-id verify [signature file] --input \"contents to verify\"",
	RunE:    runVerifyCmd,
}

func init() {
	params := []cli.Param{
		{
			Name:      "input",
			Usage:     "original contents to run the verification against",
			FlagKey:   "verify.input",
			ByDefault: "",
		},
	}
	if err := cli.SetupCommandParams(verifyCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(verifyCmd)
}

func runVerifyCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("you must provide the signature file to verify")
	}

	// Get input, CLI takes precedence, from standard input otherwise
	var input []byte
	input = []byte(viper.GetString("verify.input"))
	if len(input) == 0 {
		input, _ = cli.ReadPipedInput(maxPipeInputSize)
	}
	if len(input) == 0 {
		return errors.New("no input passed in to verify")
	}
	if len(input) > 32 {
		digest := sha3.New256()
		digest.Write(input)
		input = digest.Sum(nil)
	}

	// Load signature file
	ll := getLogger()
	ll.Info("verifying LD signature")
	ll.Debug("load signature file")
	entry, err := ioutil.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read the signature file: %s", err)
	}
	ll.Debug("decoding contents")
	sig := &did.SignatureLD{}
	if err = json.Unmarshal(entry, sig); err != nil {
		return fmt.Errorf("invalid signature file: %s", err)
	}

	// Validate signature creator
	ll.Debug("validating signature creator")
	id, err := did.Parse(sig.Creator)
	if err != nil {
		return fmt.Errorf("invalid signature creator: %s", err)
	}
	if id.Method() != "bryk" {
		return fmt.Errorf("only 'bryk' DID are supported: %s", id)
	}

	// Retrieve subject
	peer, err := retrieveSubject(id.Subject(), ll)
	if err != nil {
		return err
	}

	// Get creator's key
	ck := peer.Key(sig.Creator)
	if ck == nil {
		return fmt.Errorf("creator key is not available on the DID Document: %s", sig.Creator)
	}

	// Verify signature
	if !ck.VerifySignatureLD(input, sig) {
		return errors.New("signature is invalid")
	}
	ll.Info("signature is valid")
	return nil
}
