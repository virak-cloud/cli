package dns

import (
	"github.com/spf13/cobra"
)

var DnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage your DNS service",
	Long:  "Manage your DNS service",
}

func init() {

}
