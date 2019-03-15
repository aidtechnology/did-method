package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/bryk-io/did-method/proto"
	"github.com/bryk-io/x/did"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verifyCmd = &cobra.Command{
	Use:     "verify",
	Short:   "Check the validity of a SignatureLD document",
	Example: "bryk-id verify [signature file] --input \"contents to verify\"",
	RunE:    runVerifyCmd,
}

func init() {
	params := []cParam{
		{
			name:      "input",
			usage:     "original contents to run the verification against",
			flagKey:   "verify.input",
			byDefault: "",
		},
	}
	if err := setupCommandParams(verifyCmd, params); err != nil {
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
		input, _ = getPipedInput()
	}
	if len(input) == 0 {
		return errors.New("no input passed in to verify")
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

	// Get network connection
	conn, err := getClientConnection(ll)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Retrieve subject
	ll.Debug("retrieving record")
	client := proto.NewAgentClient(conn)
	res, err := client.Retrieve(context.TODO(), &proto.Query{Subject: id.Subject()})
	if err != nil {
		return fmt.Errorf("failed to retrieve DID records: %s", err)
	}
	if !res.Ok {
		return errors.New("no information available for the provided DID")
	}

	// Decode contents
	ll.Debug("decoding contents")
	doc := &did.Document{}
	if err = doc.Decode(res.Contents); err != nil {
		return fmt.Errorf("failed to decode received DID Document: %s", err)
	}
	peer, err := did.FromDocument(doc)
	if err != nil {
		return fmt.Errorf("failed to decode received DID Document: %s", err)
	}
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
