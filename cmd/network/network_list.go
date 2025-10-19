package network

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// listOptions holds the options for the `network list` command.
// These options are populated from command-line flags.
type listOptions struct {
	// ZoneID is the ID of the zone to list networks from.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
}

var listOpts listOptions

// networkListCmd represents the `network list` command.
// It lists all networks in a specified zone.
var networkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all networks in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &listOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.ListNetworks(zoneId)
		if err != nil {
			slog.Error("failed to list networks", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		presenter.RenderNetworkList(resp.Data)
		return nil
	},
}

// init registers the `network list` command with the parent `network` command
// and binds the flags for the `listOptions` struct.
func init() {
	_ = cli.BindFlagsFromStruct(networkListCmd, &listOpts)
	NetworkCmd.AddCommand(networkListCmd)
}
