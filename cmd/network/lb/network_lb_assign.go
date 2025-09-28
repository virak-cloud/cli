package lb

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

type lbAssignOptions struct {
	ZoneID             string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	NetworkID          string `flag:"networkId" usage:"Network ID for the load balancer"`
	RuleID             string `flag:"ruleId" usage:"Load balancer rule ID"`
	InstanceNetworkIds string `flag:"instanceNetworkIds" usage:"Comma-separated instance network IDs to assign"`
}

var lbAssignOpts lbAssignOptions

// NetworkLbAssignCmd is the command for assigning instances to a load balancing rule.
var NetworkLbAssignCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign instances to a load balancing rule",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("ruleId"),
			cli.Required("instanceNetworkIds"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &lbAssignOpts); err != nil {
			return err
		}
		instanceIds := strings.Split(lbAssignOpts.InstanceNetworkIds, ",")
		for i := range instanceIds {
			instanceIds[i] = strings.TrimSpace(instanceIds[i])
			if instanceIds[i] == "" {
				return fmt.Errorf("instanceNetworkIds contains empty value")
			}
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.AssignLoadBalancerRule(zoneID, lbAssignOpts.NetworkID, lbAssignOpts.RuleID, instanceIds)
		if err != nil {
			slog.Error("failed to assign instances to load balancer rule", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if resp.Data.Success {
			slog.Info("instances assigned to load balancer rule successfully")
			fmt.Println("Instances assigned to load balancer rule successfully.")
		} else {
			slog.Error("assign instances to load balancer rule unsuccessful", "response", resp)
			return fmt.Errorf("failed to assign instances to load balancer rule")
		}
		return nil
	},
}

func init() {
	NetworkLbCmd.AddCommand(NetworkLbAssignCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbAssignCmd, &lbAssignOpts)
}
