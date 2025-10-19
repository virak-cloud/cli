package ipv4

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// firewallIPv4DeleteOptions holds the options for the `network firewall ipv4 delete` command.
// These options are populated from command-line flags.
type firewallIPv4DeleteOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to delete the firewall rule from.
	NetworkID string `flag:"networkId" desc:"Network ID (required)"`
	// RuleID is the ID of the firewall rule to be deleted.
	RuleID    string `flag:"ruleId" desc:"Firewall Rule ID (required)"`
}

var firewallIPv4DeleteOpts firewallIPv4DeleteOptions

// NetworkFirewallIPv4DeleteCmd represents the `network firewall ipv4 delete` command.
// It deletes an IPv4 firewall rule from a network.
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

// init registers the `network firewall ipv4 delete` command with the parent `network firewall ipv4` command
// and binds the flags for the `firewallIPv4DeleteOptions` struct.
func init() {
	NetworkFirewallIPv4Cmd.AddCommand(NetworkFirewallIPv4DeleteCmd)
	_ = cli.BindFlagsFromStruct(NetworkFirewallIPv4DeleteCmd, &firewallIPv4DeleteOpts)
}
