package network

import (
	"virak-cli/cmd/network/create"
	"virak-cli/cmd/network/firewall"
	"virak-cli/cmd/network/instance"
	"virak-cli/cmd/network/lb"
	publicip "virak-cli/cmd/network/public-ip"

	"virak-cli/cmd/network/vpn"

	"github.com/spf13/cobra"
)

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks in a zone",
	Long:  `Perform operations on networks such as create, list, show, and delete.`,
}

func init() {
	NetworkCmd.AddCommand(create.NetworkCreateCmd)
	NetworkCmd.AddCommand(firewall.NetworkFirewallCmd)
	NetworkCmd.AddCommand(instance.NetworkInstanceCmd)
	NetworkCmd.AddCommand(lb.NetworkLbCmd)
	NetworkCmd.AddCommand(publicip.NetworkPublicIPCmd)
	NetworkCmd.AddCommand(vpn.NetworkVpnCmd)
}
