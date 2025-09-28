package vpn

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

type vpnEnableOptions struct {
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	NetworkID string `flag:"networkId" usage:"Network ID for the VPN"`
}

var vpnEnableOpts vpnEnableOptions

// NetworkVpnEnableCmd is the command to enable VPN for a network.
var NetworkVpnEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable VPN for a network",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd, cli.Required("networkId"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &vpnEnableOpts); err != nil {
			return err
		}
		client := http.NewClient(token)
		resp, err := client.EnableNetworkVpn(zoneID, vpnEnableOpts.NetworkID)
		if err != nil {
			if strings.Contains(err.Error(), "The network should be Implemented") {
				slog.Error("network not implemented, need to connect instance", "error", err)
				return fmt.Errorf("first need to connect an instance to this network")
			}
			slog.Error("failed to enable VPN", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if resp.Data.Success {
			slog.Info("VPN enable started successfully")
			fmt.Println("VPN enable started successfully.")
		} else {
			slog.Error("VPN enable unsuccessful", "response", resp)
			return fmt.Errorf("failed to start VPN enable")
		}
		return nil
	},
}

func init() {
	NetworkVpnCmd.AddCommand(NetworkVpnEnableCmd)
	_ = cli.BindFlagsFromStruct(NetworkVpnEnableCmd, &vpnEnableOpts)
}
