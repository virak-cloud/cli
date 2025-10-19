package create

import (
	"github.com/spf13/cobra"
)

// NetworkCreateCmd is the parent command for all network creation commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new network in a zone",
	Long:  `Create a new Layer 2 (L2) or Layer 3 (L3) network in a specified zone.`,
}

func init() {
}
