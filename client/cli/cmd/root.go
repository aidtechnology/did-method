package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aidtechnology/did-method/client/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	xlog "go.bryk.io/pkg/log"
	"go.bryk.io/x/errors"
)

var (
	log     xlog.Logger        // app main logger
	conf    *internal.Settings // app settings management
	cfgFile = ""               // configuration file used
	homeDir = ""               // home directory used
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
	// catch any panics
	defer func() {
		if err := errors.FromRecover(recover()); err != nil {
			log.Warning("recovered panic")
			fmt.Printf("%+v", err)
			os.Exit(1)
		}
	}()
	// execute command
	if err := rootCmd.Execute(); err != nil {
		log.WithField("error", err).Error("command failed")
		os.Exit(1)
	}
}

func init() {
	log = xlog.WithZero(xlog.ZeroOptions{PrettyPrint: true})
	conf = new(internal.Settings)
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file ($HOME/.didctl/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&homeDir, "home", "", "home directory ($HOME/.didctl)")
	if err := viper.BindPFlag("client.home", rootCmd.PersistentFlags().Lookup("home")); err != nil {
		panic(err)
	}
}

func initConfig() {
	// Used for ENV variables prefix and home directories
	var appIdentifier = "didctl"

	// Find home directory
	home := homeDir
	if home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		home = h
	}

	// Set default values
	conf.SetDefaults(viper.GetViper(), home, appIdentifier)

	// Set configuration file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(filepath.Join(home, fmt.Sprintf(".%s", appIdentifier)))
		viper.AddConfigPath(filepath.Join(home, appIdentifier))
		viper.AddConfigPath(filepath.Join("/etc", appIdentifier))
		viper.AddConfigPath(".")
	}

	// ENV
	viper.SetEnvPrefix(appIdentifier)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil && viper.ConfigFileUsed() != "" {
		log.WithField("file", viper.ConfigFileUsed()).Error("failed to load configuration file")
	}
	if cf := viper.ConfigFileUsed(); cf != "" {
		log.WithField("file", cf).Info("configuration loaded")
	}

	// Load configuration into "settings" helper
	if err := conf.Load(viper.GetViper()); err != nil {
		panic(err)
	}
}
