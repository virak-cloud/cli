package haproxy

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// lbHaproxyLiveOptions holds the options for the `network lb haproxy live` command.
// These options are populated from command-line flags.
type lbHaproxyLiveOptions struct {
	// ZoneID is the ID of the zone where the load balancer is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to get the HAProxy live report for.
	NetworkID string `flag:"networkId" usage:"Network ID for the load balancer"`
}

var lbHaproxyLiveOpts lbHaproxyLiveOptions

// NetworkLbHaproxyLiveCmd represents the `network lb haproxy live` command.
// It gets live HAProxy statistics for a network.
var NetworkLbHaproxyLiveCmd = &cobra.Command{
	Use:   "live",
	Short: "Get live HAProxy statistics for the network",
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
		if err := cli.LoadFromCobraFlags(cmd, &lbHaproxyLiveOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.GetHaproxyLive(zoneId, lbHaproxyLiveOpts.NetworkID)
		if err != nil {
			slog.Error("failed to get HAProxy live report", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		fmt.Printf("Updated At: %d\n", resp.Data.UpdatedAt)
		if len(resp.Data.Rules) == 0 {
			fmt.Println("No HAProxy rules found.")
			return nil
		}
		fmt.Println("ID\tName\tAlgorithm\tPublicPort\tPrivatePort\tStatus")
		for _, rule := range resp.Data.Rules {
			fmt.Printf("%s\t%s\t%s\t%d\t%d\t%s\n", rule.ID, rule.Name, rule.Algorithm, rule.PublicPort, rule.PrivatePort, rule.Status)
		}
		return nil
	},
}

// init registers the `network lb haproxy live` command with the parent `network lb haproxy` command
// and binds the flags for the `lbHaproxyLiveOptions` struct.
func init() {
	NetworkLbHaproxyCmd.AddCommand(NetworkLbHaproxyLiveCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbHaproxyLiveCmd, &lbHaproxyLiveOpts)
}
