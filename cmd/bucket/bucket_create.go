package bucket

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"virak-cli/internal/cli"
	httpc "virak-cli/pkg/http"
)

type createOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	Name        string `flag:"name" usage:"Name of the bucket"`
	Policy      string `flag:"policy" default:"Private" usage:"Policy (Private|Public)"`
}

var createOpt createOptions

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

func init() {
	ObjectStorageCmd.AddCommand(bucketCreateCmd)

	_ = cli.BindFlagsFromStruct(bucketCreateCmd, &createOpt)
}
