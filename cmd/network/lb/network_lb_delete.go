package lb

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type lbDeleteOptions struct {
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	NetworkID string `flag:"networkId" usage:"Network ID for the load balancer"`
	RuleID    string `flag:"ruleId" usage:"Load balancer rule ID"`
}

var lbDeleteOpts lbDeleteOptions

// NetworkLbDeleteCmd is the command for deleting a load balancing rule.
var NetworkLbDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a load balancing rule by its ID",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("ruleId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &lbDeleteOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		resp, err := httpClient.DeleteLoadBalancerRule(zoneID, lbDeleteOpts.NetworkID, lbDeleteOpts.RuleID)
		if err != nil {
			slog.Error("failed to delete load balancer rule", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if resp.Data.Success {
			slog.Info("load balancer rule deleted successfully")
			fmt.Println("Load balancer rule deleted successfully.")
		} else {
			slog.Error("load balancer rule deletion unsuccessful", "response", resp)
			return fmt.Errorf("failed to delete load balancer rule")
		}
		return nil
	},
}

func init() {
	NetworkLbCmd.AddCommand(NetworkLbDeleteCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbDeleteCmd, &lbDeleteOpts)
}
