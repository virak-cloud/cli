package haproxy

import "github.com/spf13/cobra"

// NetworkLbHaproxyCmd is the parent command for all HAProxy related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var NetworkLbHaproxyCmd = &cobra.Command{
	Use:   "haproxy",
	Short: "HAProxy reports for network load balancer",
}

// init registers the subcommands for the `network lb haproxy` command.
func init() {
	NetworkLbHaproxyCmd.AddCommand(NetworkLbHaproxyLiveCmd)
	NetworkLbHaproxyCmd.AddCommand(NetworkLbHaproxyLogCmd)
}
