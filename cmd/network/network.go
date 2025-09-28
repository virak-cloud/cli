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
