package bucket

import (
	"github.com/spf13/cobra"
)

// ObjectStorageCmd is the parent command for all object storage related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var ObjectStorageCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Manage object storage resources",
}

func init() {

}
