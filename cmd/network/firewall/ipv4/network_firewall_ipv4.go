package ipv4

import (
	"github.com/spf13/cobra"
)

var NetworkFirewallIPv4Cmd = &cobra.Command{
	Use:   "ipv4",
	Short: "Manage IPv4 firewall rules",
}
