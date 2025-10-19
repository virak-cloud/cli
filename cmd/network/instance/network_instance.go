package instance

import (
	"github.com/spf13/cobra"
)

// NetworkInstanceCmd is the parent command for all network instance related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkInstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage instances connected to a network",
	Long:  `Perform operations on instances connected to a network such as connect, disconnect, and list.`,
}
