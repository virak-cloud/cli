package instance

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// networkInstanceConnectOptions holds the options for the `network instance connect` command.
// These options are populated from command-line flags.
type networkInstanceConnectOptions struct {
	// ZoneID is the ID of the zone where the network and instance are located.
	// This is optional if a default zone is set in the config.
	ZoneID     string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to connect the instance to.
	NetworkID  string `flag:"networkId" usage:"Network ID"`
	// InstanceID is the ID of the instance to be connected.
	InstanceID string `flag:"instanceId" usage:"Instance ID"`
}

var networkInstanceConnectOpt networkInstanceConnectOptions

// NetworkInstanceConnectCmd represents the `network instance connect` command.
// It connects an instance to a network.
var NetworkInstanceConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect an instance to a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("instanceId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &networkInstanceConnectOpt); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.ConnectInstanceToNetwork(zoneID, networkInstanceConnectOpt.NetworkID, networkInstanceConnectOpt.InstanceID)
		if err != nil {
			slog.Error("failed to connect instance to network", "error", err)
			return fmt.Errorf("failed to connect instance: %w", err)
		}
		if resp.Data.Success {
			slog.Info("instance connected to network successfully")
			fmt.Println("Instance successfully connected to network.")
		} else {
			slog.Error("connect instance to network unsuccessful", "response", resp)
			return fmt.Errorf("failed to connect instance to network")
		}
		return nil
	},
}

// init registers the `network instance connect` command with the parent `network instance` command
// and binds the flags for the `networkInstanceConnectOptions` struct.
func init() {
	NetworkInstanceCmd.AddCommand(NetworkInstanceConnectCmd)
	_ = cli.BindFlagsFromStruct(NetworkInstanceConnectCmd, &networkInstanceConnectOpt)
}
