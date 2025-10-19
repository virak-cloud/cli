package vpn

import (
	"github.com/spf13/cobra"
)

// NetworkVpnCmd is the parent command for all VPN related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkVpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Manage VPN for a network",
}

func init() {

}
