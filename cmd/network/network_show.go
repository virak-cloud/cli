package network

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// showOptions holds the options for the `network show` command.
// These options are populated from command-line flags.
type showOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to be shown.
	NetworkID string `flag:"networkId" usage:"Network ID to show"`
}

var showOpts showOptions

// networkShowCmd represents the `network show` command.
// It shows the details of a specific network in a zone.
var networkShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a specific network in a zone",
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
		if err := cli.LoadFromCobraFlags(cmd, &showOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.ShowNetwork(zoneID, showOpts.NetworkID)
		if err != nil {
			slog.Error("failed to show network", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		presenter.RenderNetworkDetail(resp.Data)
		return nil
	},
}

// init registers the `network show` command with the parent `network` command
// and binds the flags for the `showOptions` struct.
func init() {
	_ = cli.BindFlagsFromStruct(networkShowCmd, &showOpts)
	NetworkCmd.AddCommand(networkShowCmd)
}
