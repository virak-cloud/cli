package ipv4

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

type firewallIPv4ListOptions struct {
	ZoneID    string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	NetworkID string `flag:"networkId" desc:"Network ID (required)"`
}

var firewallIPv4ListOpts firewallIPv4ListOptions

var NetworkFirewallIPv4ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List IPv4 firewall rules for a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.IsUlid("networkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &firewallIPv4ListOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListIPv4FirewallRules(zoneId, firewallIPv4ListOpts.NetworkID)
		if err != nil {
			slog.Error("failed to list IPv4 firewall rules", "error", err)
			return fmt.Errorf("failed to list IPv4 firewall rules: %w", err)
		}

		if len(resp.Data) == 0 {
			fmt.Println("No IPv4 firewall rules found.")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Protocol", "TrafficType", "Source", "Destination", "Status", "CreatedAt"})

		for _, rule := range resp.Data {
			table.Append([]string{rule.ID, rule.Protocol, rule.TrafficType, rule.IPSource, rule.IPDestination, rule.Status, fmt.Sprintf("%d", rule.CreatedAt)})
		}
		table.Render()
		return nil
	},
}

func init() {
	NetworkFirewallIPv4Cmd.AddCommand(NetworkFirewallIPv4ListCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv4ListCmd, &firewallIPv4ListOpts)
}
