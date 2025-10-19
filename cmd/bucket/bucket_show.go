package bucket

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// showOptions holds the options for the `bucket show` command.
// These options are populated from command-line flags.
type showOptions struct {
	// ZoneID is the ID of the zone where the bucket is located.
	// This is optional if a default zone is set in the config.
	ZoneID   string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// BucketID is the ID of the bucket to be shown.
	BucketID string `flag:"bucketId" usage:"Id of the bucket"`
}

var showOpt showOptions

// objectStorageShowCmd represents the `bucket show` command.
// It shows the details of a specific object storage bucket.
var objectStorageShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of an object storage bucket",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.IsUlid("bucketId"),
			cli.Required("bucketId"),
		)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &showOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		bucketResponse, err := httpClient.GetObjectStorageBucket(zoneID, showOpt.BucketID)
		if err != nil {
			slog.Error("failed to get object storage bucket", "error", err, "zoneID", zoneID, "bucketId", showOpt.BucketID)
			return fmt.Errorf("error: %w", err)

		}

		slog.Info("successfully retrieved object storage bucket", "zoneID", zoneID, "bucketId", showOpt.BucketID)
		presenter.RenderBucketDetail(bucketResponse.Data)
		return nil
	},
}

// init registers the `bucket show` command with the parent `bucket` command
// and binds the flags for the `showOptions` struct.
func init() {
	ObjectStorageCmd.AddCommand(objectStorageShowCmd)
	_ = cli.BindFlagsFromStruct(objectStorageShowCmd, &showOpt)

}
