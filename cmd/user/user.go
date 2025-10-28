package user

import (
	"github.com/spf13/cobra"
)

// UserCmd is the root command for managing user-related operations.
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user profile and authentication",
}

func init() {

}