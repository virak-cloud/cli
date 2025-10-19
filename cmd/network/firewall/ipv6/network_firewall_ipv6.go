package ipv6

import (
	"github.com/spf13/cobra"
)

// NetworkFirewallIPv6Cmd is the parent command for all IPv6 firewall rule commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkFirewallIPv6Cmd = &cobra.Command{
	Use:   "ipv6",
	Short: "Manage IPv6 firewall rules",
}
