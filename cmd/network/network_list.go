package network

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/internal/presenter"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type listOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default zone from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
}

var listOpts listOptions

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

func init() {
	_ = cli.BindFlagsFromStruct(networkListCmd, &listOpts)
	NetworkCmd.AddCommand(networkListCmd)
}
