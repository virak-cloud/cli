package instance

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

// vmImageListOptions holds the options for the `instance vm-image-list` command.
// These options are populated from command-line flags.
type vmImageListOptions struct {
	// ZoneID is the ID of the zone to list VM images from.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
}

var vmImageListOpt vmImageListOptions

// instanceVMImageListCmd represents the `instance vm-image-list` command.
// It lists all available VM images for instances in a specified zone.
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

// renderInstanceVMImages renders a table of instance VM images.
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

// init registers the `instance vm-image-list` command with the parent `instance` command
// and binds the flags for the `vmImageListOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceVMImageListCmd)
	_ = cli.BindFlagsFromStruct(instanceVMImageListCmd, &vmImageListOpt)
}
