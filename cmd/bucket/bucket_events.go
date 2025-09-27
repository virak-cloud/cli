package bucket

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/internal/presenter"

	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type eventsOptions struct {
	ZoneID   string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	BucketID string `flag:"bucketId" usage:"Id of the bucket"`
}

var eventOpt eventsOptions
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

func init() {
	ObjectStorageCmd.AddCommand(objectStorageEventsCmd)

	_ = cli.BindFlagsFromStruct(objectStorageEventsCmd, &eventOpt)

}
