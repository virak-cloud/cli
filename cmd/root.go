package cmd

import (
	"errors"
	"fmt"
	"os"
	bucket "virak-cli/cmd/bucket"
	"virak-cli/cmd/cluster"
	"virak-cli/cmd/dns"
	"virak-cli/cmd/instance"
	"virak-cli/cmd/network"
	"virak-cli/cmd/zone"
	"virak-cli/internal/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var disableLog bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "virak-cli",
	Short: "A command-line interface for interacting with the Virak Cloud API, built with the Go programming language.",
	Long:  `The vk-cloud CLI is a command-line interface that allows you to manage your Virak Cloud resources directly from your terminal.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !disableLog {
			logger.InitLogger()
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&disableLog, "disable-log", false, "Disable logging")
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(bucket.ObjectStorageCmd)
	RootCmd.AddCommand(instance.InstanceCmd)
	RootCmd.AddCommand(network.NetworkCmd)
	RootCmd.AddCommand(zone.ZoneCmd)
	RootCmd.AddCommand(cluster.KubernetesClusterCmd)
	RootCmd.AddCommand(dns.DnsCmd)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to find home directory:", err)
		os.Exit(1)
	}

	viper.SetConfigName(".virak-cli")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Config file not found; create it
			// Set default values for other configuration options as needed
			viper.SetDefault("auth.token", "")
			viper.SetDefault("default.zoneId", "")
			viper.SetDefault("default.zoneName", "")
			err := viper.SafeWriteConfig()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unable to write default config:", err)
				return
			}
		}
	}
}
