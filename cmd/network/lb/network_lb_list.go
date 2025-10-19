package lb

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/virak-cloud/cli/internal/cli"
	http "github.com/virak-cloud/cli/pkg/http"

	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

// lbListOptions holds the options for the `network lb list` command.
// These options are populated from command-line flags.
type lbListOptions struct {
	// ZoneID is the ID of the zone where the load balancer is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to list the load balancer rules for.
	NetworkID string `flag:"networkId" usage:"Network ID for the load balancer"`
}

var lbListOpts lbListOptions

// NetworkLbListCmd represents the `network lb list` command.
// It lists all load balancing rules for a network.
var NetworkLbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all load balancing rules for a network",
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
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &lbListOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.ListLoadBalancerRules(zoneID, lbListOpts.NetworkID)
		if err != nil {
			slog.Error("failed to list load balancer rules", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if len(resp.Data) == 0 {
			fmt.Println("No load balancer rules found.")
			return nil
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Algorithm", "PublicPort", "PrivatePort", "Status"})
		for _, rule := range resp.Data {
			table.Append([]string{
				rule.ID,
				rule.Name,
				rule.Algorithm,
				fmt.Sprintf("%d", rule.PublicPort),
				fmt.Sprintf("%d", rule.PrivatePort),
				rule.Status,
			})
		}
		table.Render()
		return nil
	},
}

// init registers the `network lb list` command with the parent `network lb` command
// and binds the flags for the `lbListOptions` struct.
func init() {
	NetworkLbCmd.AddCommand(NetworkLbListCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbListCmd, &lbListOpts)
}
