package ipv6

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type firewallIPv6DeleteOptions struct {
	ZoneId    string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	NetworkId string `flag:"networkId" desc:"Network ID (required)"`
	RuleId    string `flag:"ruleId" desc:"Firewall Rule ID (required)"`
}

var firewallIPv6DeleteOpts firewallIPv6DeleteOptions

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

func init() {
	NetworkFirewallIPv6Cmd.AddCommand(NetworkFirewallIPv6DeleteCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv6DeleteCmd, &firewallIPv6DeleteOpts)
}
