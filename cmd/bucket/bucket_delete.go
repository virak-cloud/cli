package bucket

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// deleteOptions holds the options for the `bucket delete` command.
// These options are populated from command-line flags.
type deleteOptions struct {
	// ZoneID is the ID of the zone where the bucket is located.
	// This is optional if a default zone is set in the config.
	ZoneID   string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// BucketID is the ID of the bucket to be deleted.
	BucketID string `flag:"bucketId" usage:"Id of the bucket"`
}

var deleteOpt deleteOptions

// objectStorageDeleteCmd represents the `bucket delete` command.
// It deletes an object storage bucket in a specified zone.
var objectStorageDeleteCmd = &cobra.Command{
	Use:   "delete [zoneId] [bucketId]",
	Short: "Delete an object storage bucket",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("bucketId"),
			cli.IsUlid("bucketId"),
		)
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &deleteOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		err := httpClient.DeleteObjectStorageBucket(zoneID, deleteOpt.BucketID)
		if err != nil {
			slog.Error("failed to delete object storage bucket", "error", err, "zoneID", zoneID, "bucketId", deleteOpt.BucketID)
			return fmt.Errorf("error: %w", err)

		}

		slog.Info("object storage bucket deleted successfully", "zoneId", zoneID, "bucketId", deleteOpt.BucketID)
		fmt.Println("Object storage bucket deleted successfully.")
		return nil
	},
}

// init registers the `bucket delete` command with the parent `bucket` command
// and binds the flags for the `deleteOptions` struct.
func init() {
	ObjectStorageCmd.AddCommand(objectStorageDeleteCmd)
	_ = cli.BindFlagsFromStruct(objectStorageDeleteCmd, &deleteOpt)
}
