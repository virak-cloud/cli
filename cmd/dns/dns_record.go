package dns

import "github.com/spf13/cobra"

// recordCmd is the parent command for all DNS record related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Manage domain records",
	Long:  "Manage domain records",
}

func init() {
	DnsCmd.AddCommand(recordCmd)
}
