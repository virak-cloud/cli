package dns

import "github.com/spf13/cobra"

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage your domains",
	Long:  "Manage your domains",
}

func init() {
	DnsCmd.AddCommand(domainCmd)
}
