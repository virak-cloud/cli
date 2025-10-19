package ipv6

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// firewallIPv6DeleteOptions holds the options for the `network firewall ipv6 delete` command.
// These options are populated from command-line flags.
type firewallIPv6DeleteOptions struct {
	// ZoneId is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneId    string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	// NetworkId is the ID of the network to delete the firewall rule from.
	NetworkId string `flag:"networkId" desc:"Network ID (required)"`
	// RuleId is the ID of the firewall rule to be deleted.
	RuleId    string `flag:"ruleId" desc:"Firewall Rule ID (required)"`
}

var firewallIPv6DeleteOpts firewallIPv6DeleteOptions

// NetworkFirewallIPv6DeleteCmd represents the `network firewall ipv6 delete` command.
// It deletes an IPv6 firewall rule from a network.
var NetworkFirewallIPv6DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an IPv6 firewall rule from a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		// Validate required flags
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("ruleId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())

		// Load options from flags
		if err := cli.LoadFromCobraFlags(cmd, &firewallIPv6DeleteOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.DeleteIPv6FirewallRule(zoneId, firewallIPv6DeleteOpts.NetworkId, firewallIPv6DeleteOpts.RuleId)
		if err != nil {
			slog.Error("failed to delete IPv6 firewall rule", "error", err)
			return fmt.Errorf("failed to delete IPv6 firewall rule: %w", err)
		}

		if resp.Data.Success {
			fmt.Println("IPv6 firewall rule deleted successfully.")
		} else {
			return fmt.Errorf("failed to delete IPv6 firewall rule")
		}
		return nil
	},
}

// init registers the `network firewall ipv6 delete` command with the parent `network firewall ipv6` command
// and binds the flags for the `firewallIPv6DeleteOptions` struct.
func init() {
	NetworkFirewallIPv6Cmd.AddCommand(NetworkFirewallIPv6DeleteCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv6DeleteCmd, &firewallIPv6DeleteOpts)
}
