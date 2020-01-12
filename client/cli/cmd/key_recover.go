package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.bryk.io/x/cli"
	"go.bryk.io/x/crypto/shamir"
)

var recoverKeyCmd = &cobra.Command{
	Use:     "recover",
	Short:   "Recover a previously generated Ed25519 cryptographic key",
	Example: "didctl did key recover --passphrase",
	RunE:    runRecoverKeyCmd,
}

func init() {
	params := []cli.Param{
		{
			Name:      "passphrase",
			Usage:     "use a passphrase to recover the key",
			FlagKey:   "key-recover.passphrase",
			ByDefault: false,
		},
		{
			Name:      "shared-secret",
			Usage:     "provide a comma separated list of share files",
			FlagKey:   "key-recover.shares",
			ByDefault: "",
		},
	}
	if err := cli.SetupCommandParams(recoverKeyCmd, params); err != nil {
		panic(err)
	}
	// keyCmd.AddCommand(recoverKeyCmd)
}

func runRecoverKeyCmd(_ *cobra.Command, _ []string) error {
	// Validate parameters
	pp := viper.GetBool("key-recover.passphrase")
	shares := strings.TrimSpace(viper.GetString("key-recover.shares"))
	if !pp && shares == "" {
		return errors.New("you must specify a recovery mechanism")
	}
	if pp && shares != "" {
		return errors.New("only one recovery mechanism might be used")
	}

	// Recover secret
	secret, err := recoverSecret(pp, strings.Split(shares, ","))
	if err != nil {
		return err
	}

	// Recover key
	kp, err := keyFromMaterial(secret)
	if err != nil {
		return fmt.Errorf("failed to recreate key: %s", err)
	}
	defer kp.Destroy()
	fmt.Printf("\nkey recovered: %x\n", *kp.Private)
	return nil
}

func recoverSecret(pp bool, shareFile []string) ([]byte, error) {
	// Use passphrase
	if pp {
		secret, err := cli.ReadSecure("\nEnter the passphrase used when creating the key: ")
		if err != nil {
			return nil, err
		}
		confirmation, err := cli.ReadSecure("\nConfirm the provided value: ")
		if err != nil {
			return nil, err
		}
		fmt.Println("")
		if !bytes.Equal(secret, confirmation) {
			return nil, errors.New("the values provided are not equal")
		}
		return secret, nil
	}

	// Use secret sharing
	var parts [][]byte
	for _, s := range shareFile {
		c, err := ioutil.ReadFile(filepath.Clean(s))
		if err != nil {
			return nil, fmt.Errorf("failed to load the share: %s", s)
		}
		parts = append(parts, c)
	}
	return shamir.Combine(parts)
}
