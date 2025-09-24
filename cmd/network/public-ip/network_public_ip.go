package publicip

import (
	"github.com/spf13/cobra"

	staticnat "virak-cli/cmd/network/public-ip/staticnat"
)

// NetworkPublicIPCmd represents the public IP command group
var NetworkPublicIPCmd = &cobra.Command{
	Use:     "public ip",
	Aliases: []string{"publicip", "ip"},
	Short:   "Manage public IPs in a network",
	Long:    `Perform operations on public IPs such as list, associate, disassociate, enable/disable static NAT.`,
}

func init() {

	NetworkPublicIPCmd.AddCommand(staticnat.NetworkPublicIPStaticNatCmd)
}
