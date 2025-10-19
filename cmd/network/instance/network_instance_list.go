package instance

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// networkInstanceListOptions holds the options for the `network instance list` command.
// These options are populated from command-line flags.
type networkInstanceListOptions struct {
	// ZoneID is the ID of the zone where the network and instance are located.
	// This is optional if a default zone is set in the config.
	ZoneID     string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to list instances for.
	NetworkID  string `flag:"networkId" usage:"Network ID"`
	// InstanceID is the ID of the instance to list network connections for.
	InstanceID string `flag:"instanceId" usage:"Instance ID"`
}

var networkInstanceListOpt networkInstanceListOptions

// NetworkInstanceListCmd represents the `network instance list` command.
// It lists all instances connected to a network.
var NetworkInstanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instances connected to a network",
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
		if err := cli.LoadFromCobraFlags(cmd, &networkInstanceListOpt); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.ListNetworkInstances(zoneID, networkInstanceListOpt.NetworkID, networkInstanceListOpt.InstanceID)
		if err != nil {
			slog.Error("failed to list instances", "error", err)
			return fmt.Errorf("failed to list instances: %w", err)
		}
		if len(resp.Data) == 0 {
			fmt.Println("No instances connected to this network.")
			return nil
		}
		presenter.RenderInstanceNetworkList(resp.Data)
		return nil
	},
}

// init registers the `network instance list` command with the parent `network instance` command
// and binds the flags for the `networkInstanceListOptions` struct.
func init() {
	NetworkInstanceCmd.AddCommand(NetworkInstanceListCmd)
	_ = cli.BindFlagsFromStruct(NetworkInstanceListCmd, &networkInstanceListOpt)
}
