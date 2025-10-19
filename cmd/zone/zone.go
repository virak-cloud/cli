package zone

import (
	"github.com/spf13/cobra"
)

// ZoneCmd is the parent command for all zone related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var ZoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Manage zones",
	Long:  `The zone command allows you to manage your zones, including listing zones, viewing networks, resources, and services.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("zone called")
	//},
}

func init() {

}
