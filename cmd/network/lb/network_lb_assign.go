package lb

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

// lbAssignOptions holds the options for the `network lb assign` command.
// These options are populated from command-line flags.
type lbAssignOptions struct {
	// ZoneID is the ID of the zone where the load balancer is located.
	// This is optional if a default zone is set in the config.
	ZoneID             string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network where the load balancer is located.
	NetworkID          string `flag:"networkId" usage:"Network ID for the load balancer"`
	// RuleID is the ID of the load balancer rule to assign instances to.
	RuleID             string `flag:"ruleId" usage:"Load balancer rule ID"`
	// InstanceNetworkIds is a comma-separated list of instance network IDs to assign to the rule.
	InstanceNetworkIds string `flag:"instanceNetworkIds" usage:"Comma-separated instance network IDs to assign"`
}

var lbAssignOpts lbAssignOptions

// NetworkLbAssignCmd represents the `network lb assign` command.
// It assigns one or more instances to a load balancing rule.
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

// init registers the `network lb assign` command with the parent `network lb` command
// and binds the flags for the `lbAssignOptions` struct.
func init() {
	NetworkLbCmd.AddCommand(NetworkLbAssignCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbAssignCmd, &lbAssignOpts)
}
