package user

import (
	"github.com/spf13/cobra"
)

// UserSSHKeyCmd represents the ssh-key command
var UserSSHKeyCmd = &cobra.Command{
	Use:   "ssh-key",
	Short: "Manage SSH keys",
}

func init() {
	UserCmd.AddCommand(UserSSHKeyCmd)
}