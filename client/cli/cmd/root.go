package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	hd "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	homeDir        string
	didDomainValue = "did.bryk.io"
	defaultNode    = "rpc-did.bryk.io:80"
)

var rootCmd = &cobra.Command{
	Use:           "bryk-did",
	Short:         "Bryk DID Method: Client",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long: `Bryk DID Method: Client

Reference client implementation for the "bryk" DID method. The platform allows
entities to fully manage Decentralized Identifiers as described on the version
v0.11 of the specification.

For more information:
https://w3c-ccg.github.io/did-spec`,
}

// Execute will process the CLI invocation
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ll := getLogger()
		ll.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file ($HOME/.bryk-did/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&homeDir, "home", "", "home directory ($HOME/.bryk-did)")
	if err := viper.BindPFlag("home", rootCmd.PersistentFlags().Lookup("home")); err != nil {
		panic(err)
	}
}

func initConfig() {
	home := ""
	if homeDir == "" {
		h, err := hd.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		home = h
	} else {
		home = homeDir
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		homeDir = path.Join(home, ".bryk-did")
		viper.AddConfigPath(homeDir)
		viper.SetConfigName("config")
	}

	// ENV
	viper.SetEnvPrefix("bryk-did")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	viper.SetDefault("node", defaultNode)
	if err := viper.ReadInConfig(); err != nil && viper.ConfigFileUsed() != "" {
		fmt.Println("failed to load configuration file:", viper.ConfigFileUsed())
	}
}
