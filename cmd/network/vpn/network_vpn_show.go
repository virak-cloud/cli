package vpn

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// vpnShowOptions holds the options for the `network vpn show` command.
// These options are populated from command-line flags.
type vpnShowOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to show VPN details for.
	NetworkID string `flag:"networkId" usage:"Network ID for the VPN"`
}

var vpnShowOpts vpnShowOptions

// NetworkVpnShowCmd represents the `network vpn show` command.
// It shows VPN details for a network.
var NetworkVpnShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show VPN details for a network",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd, cli.Required("networkId"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &vpnShowOpts); err != nil {
			return err
		}
		client := http.NewClient(token)
		resp, err := client.GetNetworkVpnDetails(zoneID, vpnShowOpts.NetworkID)
		if err != nil {
			slog.Error("failed to get VPN details", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		fmt.Printf("VPN IP Address: %s\nUsername: %s\nPassword: %s\nPreshared Key: %s\nStatus: %s\n",
			resp.Data.IPAddress, resp.Data.Username, resp.Data.Password, resp.Data.PresharedKey, resp.Data.Status)
		return nil
	},
}

// init registers the `network vpn show` command with the parent `network vpn` command
// and binds the flags for the `vpnShowOptions` struct.
func init() {
	NetworkVpnCmd.AddCommand(NetworkVpnShowCmd)
	_ = cli.BindFlagsFromStruct(NetworkVpnShowCmd, &vpnShowOpts)
}
