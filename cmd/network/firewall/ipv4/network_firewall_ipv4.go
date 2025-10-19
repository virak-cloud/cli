package ipv4

import (
	"github.com/spf13/cobra"
)

// NetworkFirewallIPv4Cmd is the parent command for all IPv4 firewall rule commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkFirewallIPv4Cmd = &cobra.Command{
	Use:   "ipv4",
	Short: "Manage IPv4 firewall rules",
}
