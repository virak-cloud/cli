package zone

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

// zoneResourcesOptions holds the options for the `zone resources` command.
// These options are populated from command-line flags.
type zoneResourcesOptions struct {
	// ZoneID is the ID of the zone to list resources for.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" desc:"Zone ID to use (optional if default.zoneId is set in config, overrides positional argument if set)"`
}

var resourcesOpt zoneResourcesOptions

// resourcesCmd represents the `zone resources` command.
// It lists the resources for a specific zone.
var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "List of resources for a specific zone",
	Long: `Get a list of resources for a specific zone.
You must provide the zone id as an argument, or use --default-zone or --zoneId flag.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.OneOf("default-zone", "true", "false"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &resourcesOpt); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resources, err := httpClient.GetZoneCustomerResource(zoneID)
		if err != nil {
			slog.Error("failed to get zone resources", "error", err, "zoneId", zoneID)
			return fmt.Errorf("error: %w", err)
		}
		slog.Info("successfully retrieved zone resources", "zoneId", zoneID)
		fmt.Println("Resources for Zone:")
		fmt.Printf("  Memory: %d/%d (Megabyte)\n", resources.InstanceResourceCollected.Memory.Collected, resources.InstanceResourceCollected.Memory.Total)
		fmt.Printf("  CPU Number: %d/%d (Core)\n", resources.InstanceResourceCollected.CPUNumber.Collected, resources.InstanceResourceCollected.CPUNumber.Total)
		fmt.Printf("  Data Volume: %d/%d (Gigabyte)\n", resources.InstanceResourceCollected.DataVolume.Collected, resources.InstanceResourceCollected.DataVolume.Total)
		fmt.Printf("  VM Limit: %d/%d\n", resources.InstanceResourceCollected.VMLimit.Collected, resources.InstanceResourceCollected.VMLimit.Total)
		return nil
	},
}

// init registers the `zone resources` command with the parent `zone` command
// and binds the flags for the `zoneResourcesOptions` struct.
func init() {
	ZoneCmd.AddCommand(resourcesCmd)
	_ = cli.BindFlagsFromStruct(resourcesCmd, &resourcesOpt)
}
