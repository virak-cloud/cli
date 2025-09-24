package instance

import (
	"github.com/spf13/cobra"
)

// NetworkInstanceCmd is the root command for managing network instances.
var NetworkInstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage instances connected to a network",
	Long:  `Perform operations on instances connected to a network such as connect, disconnect, and list.`,
}
