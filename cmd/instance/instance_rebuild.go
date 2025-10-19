package instance

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// rebuildOptions holds the options for the `instance rebuild` command.
// These options are populated from command-line flags.
type rebuildOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID     string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to be rebuilt.
	InstanceID string `flag:"instance-id" usage:"ID of the instance to rebuild"`
	// VMImageID is the ID of the new VM image to use for rebuilding the instance.
	VMImageID  string `flag:"vm-image-id" usage:"ID of the new VM image"`
}

var rebuildOpt rebuildOptions

// instanceRebuildCmd represents the `instance rebuild` command.
// It rebuilds a virtual machine instance with a new VM image.
var instanceRebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild an instance with a new VM image",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("instance-id"),
			cli.Required("vm-image-id"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &rebuildOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.RebuildInstance(zoneID, rebuildOpt.InstanceID, rebuildOpt.VMImageID)
		if err != nil {
			slog.Error("failed to rebuild instance", "error", err, "zoneID", zoneID, "instanceID", rebuildOpt.InstanceID)
			return fmt.Errorf("failed to rebuild instance: %w", err)
		}

		if resp.Data.Success {
			fmt.Println("Instance rebuild request accepted.")
		} else {
			fmt.Println("Instance rebuild failed.")
		}
		return nil
	},
}

// init registers the `instance rebuild` command with the parent `instance` command
// and binds the flags for the `rebuildOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceRebuildCmd)
	_ = cli.BindFlagsFromStruct(instanceRebuildCmd, &rebuildOpt)
}
