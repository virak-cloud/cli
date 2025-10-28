package portforward

import (
	"github.com/spf13/cobra"
)

var NetworkPortForwardCmd = &cobra.Command{
	Use:   "port-forward",
	Short: "Manage network port forwarding rules",
}