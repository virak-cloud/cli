package staticnat

import (
	"github.com/spf13/cobra"
)

// NetworkPublicIPStaticNatCmd is the parent command for all static NAT related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkPublicIPStaticNatCmd = &cobra.Command{
	Use:   "staticnat",
	Short: "Manage static NAT for public IPs in a network",
	Long:  `Enable or disable static NAT for public IPs in a network.`,
}

func init() {

}
