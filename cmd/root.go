package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var dataDir string

var rootCmd = &cobra.Command{
	Use:   "pitwall",
	Short: "Pitwall is a live-time TUI for Formula 1 races",
	Long:  "Pitwall is a real-time terminal user interface for Formula 1 races",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeConfig(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultDataDir := filepath.Join(home, ".pitwall", "data")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pitwall/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", defaultDataDir, "data directory (default is $HOME/.pitwall/data")
}

func initializeConfig(cmd *cobra.Command) error {
	viper.SetEnvPrefix("PITWALL")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.pitwall")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
