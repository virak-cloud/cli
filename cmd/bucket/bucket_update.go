package bucket

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

type updateOptions struct {
	ZoneID   string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	BucketID string `flag:"bucketId" usage:"Id of the bucket"`
	Policy   string `flag:"policy" default:"Private" usage:"Policy (Private|Public)"`
}

var updateOpt updateOptions

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

func init() {
	ObjectStorageCmd.AddCommand(objectStorageUpdateCmd)
	_ = cli.BindFlagsFromStruct(objectStorageUpdateCmd, &updateOpt)
}
