package network

import (
	"github.com/virak-cloud/cli/cmd/network/create"
	"github.com/virak-cloud/cli/cmd/network/firewall"
	"github.com/virak-cloud/cli/cmd/network/instance"
	"github.com/virak-cloud/cli/cmd/network/lb"
	publicip "github.com/virak-cloud/cli/cmd/network/public-ip"

	"github.com/virak-cloud/cli/cmd/network/vpn"

	"github.com/spf13/cobra"
)

// NetworkCmd is the parent command for all network related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks in a zone",
	Long:  `Perform operations on networks such as create, list, show, and delete.`,
}

// init registers the subcommands for the `network` command.
func init() {
	NetworkCmd.AddCommand(create.NetworkCreateCmd)
	NetworkCmd.AddCommand(firewall.NetworkFirewallCmd)
	NetworkCmd.AddCommand(instance.NetworkInstanceCmd)
	NetworkCmd.AddCommand(lb.NetworkLbCmd)
	NetworkCmd.AddCommand(publicip.NetworkPublicIPCmd)
	NetworkCmd.AddCommand(vpn.NetworkVpnCmd)
}
