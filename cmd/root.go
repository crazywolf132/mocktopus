package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configLocation string = ""
)

var rootCmd = &cobra.Command{
	Use:   "mocktopus",
	Short: "A mocking / stubbing tool for the command line",
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVarP(&configLocation, "config", "c", "", "Custom location for the config file.")
}

func LoadConfig() *viper.Viper {
	if configLocation == "" {
		configDir, _ := os.UserConfigDir()
		configLocation, _ = filepath.Abs(filepath.Join(configDir, "mocktopus", "mocktopus.toml"))
	} else {
		configLocation, _ = filepath.Abs(configLocation)
	}
	fmt.Println("Loading config from: ", configLocation)

	config := viper.New()
	config.SetConfigType("toml")
	config.SetConfigFile(configLocation)
	config.ReadInConfig()

	return config
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
