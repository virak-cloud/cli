package instance

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type volumeListOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
}

var volumeListOpt volumeListOptions

var instanceVolumeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List volumes in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cli.Preflight(true)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &volumeListOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListInstanceVolumes(zoneID)
		if err != nil {
			slog.Error("failed to list volumes", "error", err, "zoneID", zoneID)
			return fmt.Errorf("failed to list volumes: %w", err)
		}
		renderInstanceVolumes(resp)
		return nil
	},
}

func renderInstanceVolumes(resp *responses.InstanceVolumeListResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Size", "Status"})
	for _, vol := range resp.Data {
		table.Append([]string{
			vol.ID,
			vol.Name,
			fmt.Sprintf("%d", vol.Size),
			vol.Status,
		})
	}
	table.Render()
}

func init() {
	instanceVolumeCmd.AddCommand(instanceVolumeListCmd)
	_ = cli.BindFlagsFromStruct(instanceVolumeListCmd, &volumeListOpt)
}
