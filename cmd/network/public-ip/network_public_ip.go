package publicip

import (
	"github.com/spf13/cobra"

	staticnat "github.com/virak-cloud/cli/cmd/network/public-ip/staticnat"
)

// NetworkPublicIPCmd is the parent command for all public IP related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkPublicIPCmd = &cobra.Command{
	Use:     "public-ip",
	Aliases: []string{"publicip", "ip"},
	Short:   "Manage public IPs in a network",
	Long:    `Perform operations on public IPs such as list, associate, disassociate, enable/disable static NAT.`,
}

// init registers the subcommands for the `network public-ip` command.
func init() {

	NetworkPublicIPCmd.AddCommand(staticnat.NetworkPublicIPStaticNatCmd)
}
