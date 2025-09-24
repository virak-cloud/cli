package instance

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/internal/presenter"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

// networkInstanceListOptions holds the flags for the 'network instance list' command.
type networkInstanceListOptions struct {
	ZoneID      string `flag:"zoneId" usage:"Zone ID"`
	DefaultZone bool   `flag:"default-zone" usage:"Use default zone"`
	NetworkID   string `flag:"networkId" usage:"Network ID"`
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
}

var networkInstanceListOpt networkInstanceListOptions

// NetworkInstanceListCmd is the command for listing instances connected to a network.
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

func init() {
	NetworkInstanceCmd.AddCommand(NetworkInstanceListCmd)
	_ = cli.BindFlagsFromStruct(NetworkInstanceListCmd, &networkInstanceListOpt)
}
