package ipv4

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type firewallIPv4DeleteOptions struct {
	DefaultZone bool   `flag:"default-zone" desc:"Use default zone from config"`
	ZoneID      string `flag:"zoneId" desc:"Zone ID"`
	NetworkID   string `flag:"networkId" desc:"Network ID (required)"`
	RuleID      string `flag:"ruleId" desc:"Firewall Rule ID (required)"`
}

var firewallIPv4DeleteOpts firewallIPv4DeleteOptions

var NetworkFirewallIPv4DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an IPv4 firewall rule from a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.IsUlid("networkId"),
			cli.Required("ruleId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &firewallIPv4DeleteOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.DeleteIPv4FirewallRule(zoneId, firewallIPv4DeleteOpts.NetworkID, firewallIPv4DeleteOpts.RuleID)
		if err != nil {
			slog.Error("failed to delete IPv4 firewall rule", "error", err)
			return fmt.Errorf("failed to delete IPv4 firewall rule: %w", err)
		}

		if resp.Data.Success {
			fmt.Println("IPv4 firewall rule deleted successfully.")
		} else {
			fmt.Println("Failed to delete IPv4 firewall rule.")
		}
		return nil
	},
}

func init() {
	NetworkFirewallIPv4Cmd.AddCommand(NetworkFirewallIPv4DeleteCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv4DeleteCmd, &firewallIPv4DeleteOpts)
}
