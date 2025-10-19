package bucket

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

// createOptions holds the options for the `bucket create` command.
// These options are populated from command-line flags.

type createOptions struct {
	// ZoneID is the ID of the zone where the bucket will be created.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// Name is the name of the bucket to be created.
	Name string `flag:"name" usage:"Name of the bucket"`
	// Policy is the access policy for the bucket.
	// It can be either "Private" or "Public".
	Policy string `flag:"policy" default:"Private" usage:"Policy (Private|Public)"`
}

var createOpt createOptions

// bucketCreateCmd represents the `bucket create` command.
// It creates a new object storage bucket in a specified zone.
var bucketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an object storage bucket in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling (required for this command)
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd,
			cli.Required("name"),
			cli.OneOf("policy", "Private", "Public"),
			// When not using --default-zone, zoneId is required (Preflight already ensures it)
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &createOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		if _, err := httpClient.CreateObjectStorageBucket(zoneID, createOpt.Name, createOpt.Policy); err != nil {
			slog.Error("failed to create object storage bucket", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("object storage bucket creation request accepted", "zoneID", zoneID, "name", createOpt.Name, "policy", createOpt.Policy)
		fmt.Println("Object storage bucket creation request accepted. Operation is asynchronous; check the bucket list for status.")
		return nil
	},
}

// init registers the `bucket create` command with the parent `bucket` command
// and binds the flags for the `createOptions` struct.
func init() {
	ObjectStorageCmd.AddCommand(bucketCreateCmd)

	_ = cli.BindFlagsFromStruct(bucketCreateCmd, &createOpt)
}
