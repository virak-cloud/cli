package instance

import (
	"strings"

	"github.com/spf13/cobra"
)

// InstanceCmd is the root command for managing instances.
var InstanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"vm", "virtual-machine"},
	Short:   "Manage instances in a zone",
}

// SplitAndTrim splits a string by comma and trims spaces, removing empty parts.
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
