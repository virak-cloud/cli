package vpn

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// vpnEnableOptions holds the options for the `network vpn enable` command.
// These options are populated from command-line flags.
type vpnEnableOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to enable VPN for.
	NetworkID string `flag:"networkId" usage:"Network ID for the VPN"`
}

var vpnEnableOpts vpnEnableOptions

// NetworkVpnEnableCmd represents the `network vpn enable` command.
// It enables VPN for a network.
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

// init registers the `network vpn enable` command with the parent `network vpn` command
// and binds the flags for the `vpnEnableOptions` struct.
func init() {
	NetworkVpnCmd.AddCommand(NetworkVpnEnableCmd)
	_ = cli.BindFlagsFromStruct(NetworkVpnEnableCmd, &vpnEnableOpts)
}
