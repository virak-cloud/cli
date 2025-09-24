package dns

import "github.com/spf13/cobra"

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Manage domain records",
	Long:  "Manage domain records",
}

func init() {
	DnsCmd.AddCommand(recordCmd)
}
