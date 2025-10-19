package instance

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

// networkInstanceDisconnectOptions holds the options for the `network instance disconnect` command.
// These options are populated from command-line flags.
type networkInstanceDisconnectOptions struct {
	// ZoneID is the ID of the zone where the network and instance are located.
	// This is optional if a default zone is set in the config.
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to disconnect the instance from.
	NetworkID         string `flag:"networkId" usage:"Network ID"`
	// InstanceID is the ID of the instance to be disconnected.
	InstanceID        string `flag:"instanceId" usage:"Instance ID"`
	// InstanceNetworkID is the ID of the instance's network interface to be disconnected.
	InstanceNetworkID string `flag:"instanceNetworkId" usage:"Instance Network ID"`
}

var networkInstanceDisConnectOpt networkInstanceDisconnectOptions

// NetworkInstanceDisconnectCmd represents the `network instance disconnect` command.
// It disconnects an instance from a network.
var NetworkInstanceDisconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect an instance from a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("instanceId"),
			cli.Required("instanceNetworkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &networkInstanceDisConnectOpt); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.DisconnectInstanceFromNetwork(zoneID, networkInstanceDisConnectOpt.NetworkID, networkInstanceDisConnectOpt.InstanceID, networkInstanceDisConnectOpt.InstanceNetworkID)
		if err != nil {
			if strings.Contains(err.Error(), "instance_network_id you dont have access") {
				slog.Error("invalid instance_network_id or no access", "error", err)
				return fmt.Errorf("the provided instance_network_id is invalid or you do not have access. Please check and try again")
			}
			slog.Error("failed to disconnect instance from network", "error", err)
			return fmt.Errorf("failed to disconnect instance: %w", err)
		}
		if resp.Data.Success {
			slog.Info("instance disconnected from network successfully")
			fmt.Println("Instance successfully disconnected from network.")
		} else {
			slog.Error("disconnect instance from network unsuccessful", "response", resp)
			return fmt.Errorf("failed to disconnect instance from network")
		}
		return nil
	},
}

// init registers the `network instance disconnect` command with the parent `network instance` command
// and binds the flags for the `networkInstanceDisconnectOptions` struct.
func init() {
	NetworkInstanceCmd.AddCommand(NetworkInstanceDisconnectCmd)
	_ = cli.BindFlagsFromStruct(NetworkInstanceDisconnectCmd, &networkInstanceDisConnectOpt)
}
