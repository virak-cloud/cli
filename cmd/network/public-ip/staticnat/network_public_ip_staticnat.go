package staticnat

import (
	"github.com/spf13/cobra"
)

// NetworkPublicIPStaticNatCmd represents the staticnat subgroup under network public ip
var NetworkPublicIPStaticNatCmd = &cobra.Command{
	Use:   "staticnat",
	Short: "Manage static NAT for public IPs in a network",
	Long:  `Enable or disable static NAT for public IPs in a network.`,
}

func init() {

}
