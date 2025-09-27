package bucket

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/internal/presenter"

	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

// ListOptions contains options for listing object storage buckets.
type ListOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
}

var listOpts ListOptions

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

func init() {
	ObjectStorageCmd.AddCommand(objectStorageListCmd)
	_ = cli.BindFlagsFromStruct(objectStorageListCmd, &listOpts)

}
