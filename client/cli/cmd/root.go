package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	homeDir string
)

var rootCmd = &cobra.Command{
	Use:           "bryk-id",
	Short:         "Bryk Identity: client application",
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() { 
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file ($HOME/.bryk-id/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&homeDir, "home", "", "home directory ($HOME/.bryk-id)")
	if err := viper.BindPFlag("home", rootCmd.PersistentFlags().Lookup("home")); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		homeDir = path.Join(home, ".bryk-id")
		viper.AddConfigPath(homeDir)
		viper.SetConfigName("config")
	}

	// ENV
	viper.SetEnvPrefix("bryk-id")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil && viper.ConfigFileUsed() != "" {
		fmt.Println("failed to load configuration file:", viper.ConfigFileUsed())
	}
}
