package lb

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// lbDeassignOptions holds the options for the `network lb deassign` command.
// These options are populated from command-line flags.
type lbDeassignOptions struct {
	// ZoneID is the ID of the zone where the load balancer is located.
	// This is optional if a default zone is set in the config.
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network where the load balancer is located.
	NetworkID         string `flag:"networkId" usage:"Network ID for the load balancer"`
	// RuleID is the ID of the load balancer rule to de-assign the instance from.
	RuleID            string `flag:"ruleId" usage:"Load balancer rule ID"`
	// InstanceNetworkID is the ID of the instance network to de-assign from the rule.
	InstanceNetworkID string `flag:"instanceNetworkId" usage:"Instance network ID to de-assign"`
}

var lbDeassignOpts lbDeassignOptions

// NetworkLbDeassignCmd represents the `network lb deassign` command.
// It de-assigns an instance from a load balancing rule.
var NetworkLbDeassignCmd = &cobra.Command{
	Use:   "deassign",
	Short: "De-assign an instance from a load balancing rule",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("ruleId"),
			cli.Required("instanceNetworkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &lbDeassignOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.DeassignLoadBalancerRule(zoneID, lbDeassignOpts.NetworkID, lbDeassignOpts.RuleID, lbDeassignOpts.InstanceNetworkID)
		if err != nil {
			slog.Error("failed to de-assign instance from load balancer rule", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if resp.Data.Success {
			slog.Info("instance de-assigned from load balancer rule successfully")
			fmt.Println("Instance de-assigned from load balancer rule successfully.")
		} else {
			slog.Error("failed to de-assign instance from load balancer rule", "response", resp)
			return fmt.Errorf("de-assign failed")
		}
		return nil
	},
}

// init registers the `network lb deassign` command with the parent `network lb` command
// and binds the flags for the `lbDeassignOptions` struct.
func init() {
	NetworkLbCmd.AddCommand(NetworkLbDeassignCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbDeassignCmd, &lbDeassignOpts)
}
