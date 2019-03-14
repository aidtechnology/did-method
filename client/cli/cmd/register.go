package cmd

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/bryk-io/did-method/client/store"
	"github.com/bryk-io/x/crypto/shamir"
	"github.com/bryk-io/x/did"
	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var registerCmd = &cobra.Command{
	Use:     "register",
	Short:   "Creates a new DID locally",
	Example: "bryk-did register [DID reference name]",
	Aliases: []string{"create", "new"},
	RunE:    runRegisterCmd,
}

func init() {
	params := []cParam{
		{
			name:      "recovery-mode",
			usage:     "choose a recovery mechanism for your primary key, 'passphrase' or 'secret-sharing'",
			flagKey:   "register.recovery-mode",
			byDefault: "secret-sharing",
		},
		{
			name:      "secret-sharing",
			usage:     "specify the number of shares and threshold value in the following format: shares,threshold",
			flagKey:   "register.secret-sharing",
			byDefault: "3,2",
		},
		{
			name:      "tag",
			usage:     "specify a tag value for the identifier instance",
			flagKey:   "register.tag",
			byDefault: "",
		},
	}
	if err := setupCommandParams(registerCmd, params); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(registerCmd)
}

func runRegisterCmd(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("a reference name for the DID is required")
	}
	name := sanitize.Name(args[0])

	// Get store handler
	ll := getLogger()
	st, err := getClientStore()
	if err != nil {
		return err
	}
	defer st.Close()

	// Check for duplicates
	if st.Get(name) != nil {
		return fmt.Errorf("there's already an entry with reference name: %s", name)
	}

	// Get key secret from the user
	ll.Info("obtaining secret material for the master private key")
	secret, err := getSecret(name)
	if err != nil {
		return err
	}

	// Generate master key from available secret
	masterKey, err := keyFromMaterial(secret)
	if err != nil {
		return err
	}
	defer masterKey.Destroy()
	pk := make([]byte, 64)
	copy(pk, masterKey.Private[:])

	// Generate base identifier instance
	ll.Info("generating new identifier")
	id, err := did.NewIdentifierWithMode("bryk", viper.GetString("register.tag"), did.ModeUUID)
	if err != nil {
		return err
	}
	ll.Debug("adding master key")
	if err = id.AddExistingKey("master", pk, did.KeyTypeEd, did.EncodingHex); err != nil {
		return err
	}
	ll.Debug("setting master key as authentication mechanism")
	if err = id.AddAuthenticationKey("master"); err != nil {
		return err
	}
	ll.Debug("generating initial integrity proof")
	if err = id.AddProof("master", didDomainValue); err != nil {
		return err
	}

	// Save instance in the store
	ll.Info("adding entry to local store")
	contents, err := id.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode DID: %s", err)
	}
	record := &store.Entry{
		Name:     name,
		Recovery: viper.GetString("register.recovery-mode"),
		Contents: contents,
	}
	return st.Save(name, record)
}

func getSecret(name string) ([]byte, error) {
	switch viper.GetString("register.recovery-mode") {
	case "secret-sharing":
		// Use random bytes as original secret
		secret := make([]byte, 128)
		if _, err := rand.Read(secret); err != nil {
			return nil, err
		}

		// Spilt secret and save shares to local files
		shares, err := splitSecret(secret, viper.GetString("register.secret-sharing"))
		if err != nil {
			return nil, err
		}
		for i, k := range shares {
			share := fmt.Sprintf("%s.share_%d.bin", name, i+1)
			if err := ioutil.WriteFile(share, k, 0400); err != nil {
				return nil, fmt.Errorf("failed to save share '%s': %s", share, err)
			}
		}
		return secret, nil
	case "passphrase":
		secret, err := secureAsk("\nEnter a secure passphrase: ")
		if err != nil {
			return nil, err
		}
		confirmation, err := secureAsk("\nConfirm the provided value: ")
		if err != nil {
			return nil, err
		}
		fmt.Println("")
		if !bytes.Equal(secret, confirmation) {
			return nil, errors.New("the values provided are not equal")
		}
		return secret, nil
	}
	return nil, errors.New("invalid recovery mode")
}

func splitSecret(secret []byte, conf string) ([][]byte, error) {
	// Load configuration
	sssConf := strings.Split(conf, ",")
	if len(sssConf) != 2 {
		return nil, errors.New("invalid secret sharing configuration value")
	}

	// Validate configuration
	shares, err := strconv.Atoi(sssConf[0])
	if err != nil {
		return nil, fmt.Errorf("invalid number shares: %s", sssConf[0])
	}
	threshold, err := strconv.Atoi(sssConf[1])
	if err != nil {
		return nil, fmt.Errorf("invalid threshold value: %s", sssConf[1])
	}
	if threshold >= shares {
		return nil, fmt.Errorf("threshold value (%d) should be smaller than the total number of shares (%d)", threshold, shares)
	}

	// Split secret
	return shamir.Split(secret, shares, threshold)
}
