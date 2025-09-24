package create

import (
	"github.com/spf13/cobra"
)

// NetworkCreateCmd represents the create command for networks
var NetworkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new network in a zone",
	Long:  `Create a new Layer 2 (L2) or Layer 3 (L3) network in a specified zone.`,
}

func init() {
}
