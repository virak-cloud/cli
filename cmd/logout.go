package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from the Virak Cloud API",
	Long:  `Logout command removes your authentication token from the Virak Cloud CLI config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("auth.token", "")
		err := viper.SafeWriteConfig()
		if err != nil {
			slog.Error("failed to clear token", "error", err)
			return fmt.Errorf("failed to clear token: %w", err)
		}
		slog.Info("logout successful")
		fmt.Println("You have been logged out.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
