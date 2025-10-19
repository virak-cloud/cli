package instance

import (
	"strings"

	"github.com/spf13/cobra"
)

// InstanceCmd is the parent command for all instance related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var InstanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"vm", "virtual-machine"},
	Short:   "Manage instances in a zone",
}

// SplitAndTrim splits a comma-separated string into a slice of strings,
// trimming whitespace from each part and removing any empty strings.
func SplitAndTrim(s string) []string {
	parts := []string{}
	for _, p := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func init() {

}
