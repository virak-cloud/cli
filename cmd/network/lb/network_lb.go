package lb

import (
	"github.com/spf13/cobra"

	haproxy "github.com/virak-cloud/cli/cmd/network/lb/haproxy"
)

// NetworkLbCmd is the root command for network load balancer operations.
var NetworkLbCmd = &cobra.Command{
	Use:     "loadbalance",
	Aliases: []string{"lb", "load-balancer"},
	Short:   "Manage network load balancer rules",
}

func init() {

	NetworkLbCmd.AddCommand(haproxy.NetworkLbHaproxyCmd)
}
