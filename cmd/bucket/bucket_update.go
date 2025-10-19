package bucket

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// updateOptions holds the options for the `bucket update` command.
// These options are populated from command-line flags.
type updateOptions struct {
	// ZoneID is the ID of the zone where the bucket is located.
	// This is optional if a default zone is set in the config.
	ZoneID   string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// BucketID is the ID of the bucket to be updated.
	BucketID string `flag:"bucketId" usage:"Id of the bucket"`
	// Policy is the new access policy for the bucket.
	// It can be either "Private" or "Public".
	Policy   string `flag:"policy" default:"Private" usage:"Policy (Private|Public)"`
}

var updateOpt updateOptions

// objectStorageUpdateCmd represents the `bucket update` command.
// It updates the policy of an existing object storage bucket.
var objectStorageUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an object storage bucket",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}

		return cli.Validate(cmd,
			cli.OneOf("policy", "Private", "Public"),
			cli.IsUlid("bucketId"),
			cli.Required("bucketId"),
		)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &updateOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.UpdateObjectStorageBucket(zoneID, updateOpt.BucketID, updateOpt.Policy)
		if err != nil {
			slog.Error("failed to update object storage bucket", "error", err, "zoneID", zoneID, "bucketId", updateOpt.BucketID, "policy", updateOpt.Policy)
			return fmt.Errorf("error: %w", err)

		}

		slog.Info("bucket update request accepted", "zoneID", zoneID, "bucketId", updateOpt.BucketID, "policy", updateOpt.Policy)
		fmt.Println("Bucket update request accepted. Operation is asynchronous; check the bucket list or show endpoint for status.")
		return nil
	},
}

// init registers the `bucket update` command with the parent `bucket` command
// and binds the flags for the `updateOptions` struct.
func init() {
	ObjectStorageCmd.AddCommand(objectStorageUpdateCmd)
	_ = cli.BindFlagsFromStruct(objectStorageUpdateCmd, &updateOpt)
}
