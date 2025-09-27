package vpn

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type vpnShowOptions struct {
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	NetworkID string `flag:"networkId" usage:"Network ID for the VPN"`
}

var vpnShowOpts vpnShowOptions

// NetworkVpnShowCmd is the command to show VPN details for a network.
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

func init() {
	NetworkVpnCmd.AddCommand(NetworkVpnShowCmd)
	_ = cli.BindFlagsFromStruct(NetworkVpnShowCmd, &vpnShowOpts)
}
