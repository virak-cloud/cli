package ipv6

import (
	"github.com/spf13/cobra"
)

var NetworkFirewallIPv6Cmd = &cobra.Command{
	Use:   "ipv6",
	Short: "Manage IPv6 firewall rules",
}
