package ipv6

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// firewallIPv6ListOptions holds the options for the `network firewall ipv6 list` command.
// These options are populated from command-line flags.
type firewallIPv6ListOptions struct {
	// ZoneId is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneId    string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	// NetworkId is the ID of the network to list the firewall rules for.
	NetworkId string `flag:"networkId" desc:"Network ID (required)"`
}

var firewallIPv6ListOpts firewallIPv6ListOptions

// NetworkFirewallIPv6ListCmd represents the `network firewall ipv6 list` command.
// It lists all IPv6 firewall rules for a network.
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

// init registers the `network firewall ipv6 list` command with the parent `network firewall ipv6` command
// and binds the flags for the `firewallIPv6ListOptions` struct.
func init() {
	NetworkFirewallIPv6Cmd.AddCommand(NetworkFirewallIPv6ListCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv6ListCmd, &firewallIPv6ListOpts)
}
