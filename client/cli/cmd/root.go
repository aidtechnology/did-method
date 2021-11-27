package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bryk-io/did-method/resolver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	xlog "go.bryk.io/pkg/log"
)

var (
	log              xlog.Logger
	cfgFile          = ""
	homeDir          = ""
	didDomainValue   = "did.bryk.io"
	defaultProviders = []*resolver.Provider{
		{
			Method:   "bryk",
			Endpoint: "https://did.bryk.io/v1/retrieve/{{.Method}}/{{.Subject}}",
			Protocol: "http",
		},
	}
)

var rootCmd = &cobra.Command{
	Use:           "didctl",
	Short:         "Bryk DID Method: Client",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long: `DID Controller

Reference client implementation for the "bryk" DID method. The platform allows
entities to fully manage Decentralized Identifiers as described on the version
v1.0 of the specification.

For more information:
https://github.com/bryk-io/did-method`,
}

// Execute will process the CLI invocation.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	log = xlog.WithZero(xlog.ZeroOptions{
		PrettyPrint: true,
		ErrorField:  "error",
	})
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file ($HOME/.didctl/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&homeDir, "home", "", "home directory ($HOME/.didctl)")
	if err := viper.BindPFlag("client.home", rootCmd.PersistentFlags().Lookup("home")); err != nil {
		panic(err)
	}
}

func initConfig() {
	// Find home directory
	home := homeDir
	if home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		home = h
	}

	// Set default values
	viper.SetDefault("resolver", defaultProviders)
	viper.SetDefault("client.timeout", 5)
	viper.SetDefault("client.home", filepath.Join(home, ".didctl"))

	// Set configuration file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(filepath.Join(home, ".didctl"))
		viper.AddConfigPath("/etc/didctl")
		viper.SetConfigName("config")
	}

	// ENV
	viper.SetEnvPrefix("didctl")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil && viper.ConfigFileUsed() != "" {
		fmt.Println("failed to load configuration file:", viper.ConfigFileUsed())
	}
}
