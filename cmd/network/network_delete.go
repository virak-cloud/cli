package network

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// deleteOptions holds the options for the `network delete` command.
// These options are populated from command-line flags.
type deleteOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to be deleted.
	NetworkID string `flag:"networkId" usage:"Network ID to delete"`
}

var deleteOpts deleteOptions

// networkDeleteCmd represents the `network delete` command.
// It deletes a specific network in a zone.
var networkDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific network in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &deleteOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.DeleteNetwork(zoneID, deleteOpts.NetworkID)
		if err != nil {
			slog.Error("failed to delete network", "error", err)
			return fmt.Errorf("error deleting network: %w", err)
		}

		if resp == nil || !resp.Data.Success {
			slog.Error("network deletion unsuccessful", "response", resp)
			return fmt.Errorf("network deletion failed")
		}

		fmt.Println("network deleted successfully")
		slog.Info("network deleted successfully")
		return nil
	},
}

// init registers the `network delete` command with the parent `network` command
// and binds the flags for the `deleteOptions` struct.
func init() {
	_ = cli.BindFlagsFromStruct(networkDeleteCmd, &deleteOpts)
	NetworkCmd.AddCommand(networkDeleteCmd)
}
