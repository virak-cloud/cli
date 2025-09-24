package lb

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type lbDeassignOptions struct {
	DefaultZone       bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use"`
	NetworkID         string `flag:"networkId" usage:"Network ID for the load balancer"`
	RuleID            string `flag:"ruleId" usage:"Load balancer rule ID"`
	InstanceNetworkID string `flag:"instanceNetworkId" usage:"Instance network ID to de-assign"`
}

var lbDeassignOpts lbDeassignOptions

// NetworkLbDeassignCmd is the command for de-assigning an instance from a load balancing rule.
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

func init() {
	NetworkLbCmd.AddCommand(NetworkLbDeassignCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbDeassignCmd, &lbDeassignOpts)
}
