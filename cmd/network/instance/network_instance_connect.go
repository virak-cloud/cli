package instance

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

// Step 1: Define options struct
type networkInstanceConnectOptions struct {
	ZoneID     string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	NetworkID  string `flag:"networkId" usage:"Network ID"`
	InstanceID string `flag:"instanceId" usage:"Instance ID"`
}

var networkInstanceConnectOpt networkInstanceConnectOptions

// NetworkInstanceConnectCmd is the command for connecting an instance to a network.
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

func init() {
	NetworkInstanceCmd.AddCommand(NetworkInstanceConnectCmd)
	_ = cli.BindFlagsFromStruct(NetworkInstanceConnectCmd, &networkInstanceConnectOpt)
}
