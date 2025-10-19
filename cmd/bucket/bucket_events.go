package bucket

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// eventsOptions holds the options for the `bucket events` command.
// These options are populated from command-line flags.
type eventsOptions struct {
	// ZoneID is the ID of the zone where the bucket is located.
	// This is optional if a default zone is set in the config.
	ZoneID   string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// BucketID is the ID of the bucket to retrieve events for.
	// If not provided, events for all buckets in the zone will be retrieved.
	BucketID string `flag:"bucketId" usage:"Id of the bucket"`
}

var eventOpt eventsOptions

// objectStorageEventsCmd represents the `bucket events` command.
// It retrieves events for a specific object storage bucket or for all buckets in a zone.
var objectStorageEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List object storage events in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling (required for this command)
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &eventOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)

		if eventOpt.BucketID != "" {
			eventsResponse, err := httpClient.GetObjectStorageBucketEvents(zoneID, eventOpt.BucketID)
			if err != nil {
				slog.Error("failed to get object storage bucket events", "error", err, "zoneID", zoneID, "bucketId", eventOpt.BucketID)
				return fmt.Errorf("error: %w", err)

			}
			slog.Info("successfully retrieved object storage bucket events", "zoneID", zoneID, "count", len(eventsResponse.Data))
			presenter.RenderBucketEvents(eventsResponse.Data)
		} else {
			eventsResponse, err := httpClient.GetObjectStorageEvents(zoneID)
			if err != nil {
				slog.Error("failed to get object storage events", "error", err, "zoneID", zoneID)
				fmt.Println("Error:", err)
				return fmt.Errorf("error: %w", err)
			}
			slog.Info("successfully retrieved object storage events", "zoneID", zoneID, "count", len(eventsResponse.Data))
			presenter.RenderBucketEvents(eventsResponse.Data)
		}
		return nil
	},
}

// init registers the `bucket events` command with the parent `bucket` command
// and binds the flags for the `eventsOptions` struct.
func init() {
	ObjectStorageCmd.AddCommand(objectStorageEventsCmd)

	_ = cli.BindFlagsFromStruct(objectStorageEventsCmd, &eventOpt)

}
