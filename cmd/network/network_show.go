package network

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/internal/presenter"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type showOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default zone from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	NetworkID   string `flag:"networkId" usage:"Network ID to show"`
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
