package haproxy

import "github.com/spf13/cobra"

var NetworkLbHaproxyCmd = &cobra.Command{
	Use:   "haproxy",
	Short: "HAProxy reports for network load balancer",
}

func init() {
	NetworkLbHaproxyCmd.AddCommand(NetworkLbHaproxyLiveCmd)
	NetworkLbHaproxyCmd.AddCommand(NetworkLbHaproxyLogCmd)
}
