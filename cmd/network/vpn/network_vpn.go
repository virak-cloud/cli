package vpn

import (
	"github.com/spf13/cobra"
)

// NetworkVpnCmd is the root command for VPN management.
var NetworkVpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Manage VPN for a network",
}

func init() {

}
