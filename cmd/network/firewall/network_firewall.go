package firewall

import (
	"github.com/spf13/cobra"

	ipv4 "virak-cli/cmd/network/firewall/ipv4"
	ipv6 "virak-cli/cmd/network/firewall/ipv6"
)

var NetworkFirewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Manage network firewalls",
}

func init() {
	NetworkFirewallCmd.AddCommand(ipv4.NetworkFirewallIPv4Cmd)
	NetworkFirewallCmd.AddCommand(ipv6.NetworkFirewallIPv6Cmd)
}
