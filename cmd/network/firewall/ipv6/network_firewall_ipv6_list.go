package ipv6

import (
	"fmt"
	"log/slog"
	"os"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type firewallIPv6ListOptions struct {
	ZoneId    string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	NetworkId string `flag:"networkId" desc:"Network ID (required)"`
}

var firewallIPv6ListOpts firewallIPv6ListOptions

var NetworkFirewallIPv6ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List IPv6 firewall rules for a network",
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

		if err := cli.LoadFromCobraFlags(cmd, &firewallIPv6ListOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListIPv6FirewallRules(zoneId, firewallIPv6ListOpts.NetworkId)
		if err != nil {
			slog.Error("failed to list IPv6 firewall rules", "error", err)
			return fmt.Errorf("failed to list IPv6 firewall rules: %w", err)
		}

		if len(resp.Data) == 0 {
			fmt.Println("No IPv6 firewall rules found.")
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
	NetworkFirewallIPv6Cmd.AddCommand(NetworkFirewallIPv6ListCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv6ListCmd, &firewallIPv6ListOpts)
}
