package instance

import (
	"fmt"
	"log/slog"
	"os"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"
	"virak-cli/pkg/http/responses"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type volumeServiceOfferingListOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
}

var volumeSoListOpt volumeServiceOfferingListOptions

var instanceVolumeServiceOfferingListCmd = &cobra.Command{
	Use:     "service-offering-list",
	Aliases: []string{"offering", "offerings", "service-offerings"},
	Short:   "List volume service offerings",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cli.Preflight(true)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &volumeSoListOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListInstanceVolumeServiceOfferings(zoneID)
		if err != nil {
			slog.Error("failed to list volume service offerings", "error", err, "zoneID", zoneID)
			return fmt.Errorf("failed to list volume service offerings: %w", err)
		}
		renderInstanceVolumeServiceOfferings(resp)
		return nil
	},
}

func renderInstanceVolumeServiceOfferings(resp *responses.InstanceVolumeServiceOfferingListResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "Size", "Price", "Public", "Featured"})
	for _, offering := range resp.Data {
		table.Append([]string{
			offering.ID,
			offering.Name,
			offering.Description,
			offering.Size,
			offering.Price,
			fmt.Sprintf("%v", offering.IsPublic),
			fmt.Sprintf("%v", offering.IsFeatured),
		})
	}
	table.Render()
}

func init() {
	instanceVolumeCmd.AddCommand(instanceVolumeServiceOfferingListCmd)
	_ = cli.BindFlagsFromStruct(instanceVolumeServiceOfferingListCmd, &volumeSoListOpt)
}
