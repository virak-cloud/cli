package network

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

type showOptions struct {
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	NetworkID string `flag:"networkId" usage:"Network ID to show"`
}

var showOpts showOptions

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

func init() {
	_ = cli.BindFlagsFromStruct(networkShowCmd, &showOpts)
	NetworkCmd.AddCommand(networkShowCmd)
}
