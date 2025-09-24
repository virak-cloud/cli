package instance

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type rebuildOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	InstanceID  string `flag:"instance-id" usage:"ID of the instance to rebuild"`
	VMImageID   string `flag:"vm-image-id" usage:"ID of the new VM image"`
}

var rebuildOpt rebuildOptions

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

func init() {
	InstanceCmd.AddCommand(instanceRebuildCmd)
	_ = cli.BindFlagsFromStruct(instanceRebuildCmd, &rebuildOpt)
}
