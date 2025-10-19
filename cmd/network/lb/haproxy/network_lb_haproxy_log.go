package haproxy

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// lbHaproxyLogOptions holds the options for the `network lb haproxy log` command.
// These options are populated from command-line flags.
type lbHaproxyLogOptions struct {
	// ZoneID is the ID of the zone where the load balancer is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to get the HAProxy logs for.
	NetworkID string `flag:"networkId" usage:"Network ID for the load balancer"`
}

var lbHaproxyLogOpts lbHaproxyLogOptions

// NetworkLbHaproxyLogCmd represents the `network lb haproxy log` command.
// It gets HAProxy logs for a network.
var NetworkLbHaproxyLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Get HAProxy logs for the network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &lbHaproxyLogOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.GetHaproxyLog(zoneId, lbHaproxyLogOpts.NetworkID)
		if err != nil {
			slog.Error("failed to get HAProxy log report", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if len(resp.Data) == 0 {
			fmt.Println("No HAProxy logs found.")
			return nil
		}
		fmt.Println("HAProxy Logs:")
		for i, log := range resp.Data {
			fmt.Printf("%d: %v\n", i+1, log)
		}
		return nil
	},
}

// init registers the `network lb haproxy log` command with the parent `network lb haproxy` command
// and binds the flags for the `lbHaproxyLogOptions` struct.
func init() {
	NetworkLbHaproxyCmd.AddCommand(NetworkLbHaproxyLogCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbHaproxyLogCmd, &lbHaproxyLogOpts)
}
