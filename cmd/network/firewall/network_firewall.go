package firewall

import (
	"github.com/spf13/cobra"

	ipv4 "github.com/virak-cloud/cli/cmd/network/firewall/ipv4"
	ipv6 "github.com/virak-cloud/cli/cmd/network/firewall/ipv6"
)

// NetworkFirewallCmd is the parent command for all network firewall related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkFirewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Manage network firewalls",
}

// init registers the subcommands for the `network firewall` command.
func init() {
	NetworkFirewallCmd.AddCommand(ipv4.NetworkFirewallIPv4Cmd)
	NetworkFirewallCmd.AddCommand(ipv6.NetworkFirewallIPv6Cmd)
}
