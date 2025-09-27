package instance

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"
	"virak-cli/pkg/http/responses"
)

type vmImageListOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
}

var vmImageListOpt vmImageListOptions

var instanceVMImageListCmd = &cobra.Command{
	Use:     "vm-image-list",
	Aliases: []string{"vm", "images"},
	Short:   "List available VM images for instances in a zone",
	PreRunE: cli.Preflight(true),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &vmImageListOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListInstanceVMImages(zoneID)
		if err != nil {
			slog.Error("failed to list instance VM images", "error", err, "zoneID", zoneID)
			return fmt.Errorf("failed to list instance VM images: %w", err)
		}
		renderInstanceVMImages(resp)
		return nil
	},
}

func renderInstanceVMImages(resp *responses.InstanceVMImageListResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Type", "OS Name", "OS Version", "Available", "Category"})
	for _, image := range resp.Data {
		table.Append([]string{
			image.ID,
			image.Name,
			image.Type,
			image.OSName,
			image.OSVersion,
			fmt.Sprintf("%t", image.IsAvailable),
			image.Category,
		})
	}
	table.Render()
}

func init() {
	InstanceCmd.AddCommand(instanceVMImageListCmd)
	_ = cli.BindFlagsFromStruct(instanceVMImageListCmd, &vmImageListOpt)
}
