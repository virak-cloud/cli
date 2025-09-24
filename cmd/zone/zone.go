package zone

import (
	"github.com/spf13/cobra"
)

// ZoneCmd represents the zone command
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
