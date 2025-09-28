package bucket

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

type showOptions struct {
	ZoneID   string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	BucketID string `flag:"bucketId" usage:"Id of the bucket"`
}

var showOpt showOptions
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

func init() {
	ObjectStorageCmd.AddCommand(objectStorageShowCmd)
	_ = cli.BindFlagsFromStruct(objectStorageShowCmd, &showOpt)

}
