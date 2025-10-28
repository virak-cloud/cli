package user

import (
	"github.com/spf13/cobra"
)

// TokenCmd is the subcommand for managing user tokens.
var TokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage user tokens",
}

func init() {
	UserCmd.AddCommand(TokenCmd)
}