package dns

import (
	"github.com/spf13/cobra"
)

// DnsCmd is the parent command for all DNS related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var DnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage your DNS service",
	Long:  "Manage your DNS service",
}

func init() {

}
