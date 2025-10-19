package bucket

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// ListOptions holds the options for the `bucket list` command.
// These options are populated from command-line flags.
type ListOptions struct {
	// ZoneID is the ID of the zone to list buckets from.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
}

var listOpts ListOptions

// objectStorageListCmd represents the `bucket list` command.
// It lists all object storage buckets in a specified zone.
var objectStorageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List object storage buckets in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		bucketsResponse, err := httpClient.GetObjectStorageBuckets(zoneID)
		if err != nil {
			slog.Error("failed to get object storage buckets", "error", err, "zoneID", zoneID)
			return fmt.Errorf("Error: %w", err)

		}

		slog.Info("successfully retrieved object storage buckets", "zoneID", zoneID, "count", len(bucketsResponse.Data))
		presenter.RenderBucketList(bucketsResponse.Data)

		return nil
	},
}

// init registers the `bucket list` command with the parent `bucket` command
// and binds the flags for the `ListOptions` struct.
func init() {
	ObjectStorageCmd.AddCommand(objectStorageListCmd)
	_ = cli.BindFlagsFromStruct(objectStorageListCmd, &listOpts)

}
