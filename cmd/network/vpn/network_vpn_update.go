package vpn

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// vpnUpdateOptions holds the options for the `network vpn update` command.
// These options are populated from command-line flags.
type vpnUpdateOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to update VPN credentials for.
	NetworkID string `flag:"networkId" usage:"Network ID for the VPN"`
}

var vpnUpdateOpts vpnUpdateOptions

// NetworkVpnUpdateCmd represents the `network vpn update` command.
// It updates VPN credentials for a network.
var NetworkVpnUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VPN credentials for a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd, cli.Required("networkId"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &vpnUpdateOpts); err != nil {
			return err
		}
		client := http.NewClient(token)
		resp, err := client.UpdateNetworkVpnCredentials(zoneID, vpnUpdateOpts.NetworkID)
		if err != nil {
			if strings.Contains(err.Error(), "The network should be Implemented") {
				slog.Error("network not implemented, need to connect instance", "error", err)
				return fmt.Errorf("first need to connect an instance to this network")
			}
			slog.Error("failed to update VPN credentials", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		if resp.Data.Success {
			slog.Info("VPN credentials update started successfully")
			fmt.Println("VPN credentials update started successfully.")
		} else {
			slog.Error("VPN credentials update unsuccessful", "response", resp)
			return fmt.Errorf("failed to start VPN credentials update")
		}
		return nil
	},
}

// init registers the `network vpn update` command with the parent `network vpn` command
// and binds the flags for the `vpnUpdateOptions` struct.
func init() {
	NetworkVpnCmd.AddCommand(NetworkVpnUpdateCmd)
	_ = cli.BindFlagsFromStruct(NetworkVpnUpdateCmd, &vpnUpdateOpts)
}
