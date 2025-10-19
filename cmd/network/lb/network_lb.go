package lb

import (
	"github.com/spf13/cobra"

	haproxy "github.com/virak-cloud/cli/cmd/network/lb/haproxy"
)

// NetworkLbCmd is the parent command for all network load balancer related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkLbCmd = &cobra.Command{
	Use:     "loadbalance",
	Aliases: []string{"lb", "load-balancer"},
	Short:   "Manage network load balancer rules",
}

// init registers the subcommands for the `network lb` command.
func init() {

	NetworkLbCmd.AddCommand(haproxy.NetworkLbHaproxyCmd)
}
