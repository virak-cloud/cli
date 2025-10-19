package dns

import "github.com/spf13/cobra"

// domainCmd is the parent command for all DNS domain related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage your domains",
	Long:  "Manage your domains",
}

func init() {
	DnsCmd.AddCommand(domainCmd)
}
