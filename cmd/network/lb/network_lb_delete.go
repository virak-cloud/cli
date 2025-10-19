package lb

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// lbDeleteOptions holds the options for the `network lb delete` command.
// These options are populated from command-line flags.
type lbDeleteOptions struct {
	// ZoneID is the ID of the zone where the load balancer is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network where the load balancer is located.
	NetworkID string `flag:"networkId" usage:"Network ID for the load balancer"`
	// RuleID is the ID of the load balancer rule to be deleted.
	RuleID    string `flag:"ruleId" usage:"Load balancer rule ID"`
}

var lbDeleteOpts lbDeleteOptions

// NetworkLbDeleteCmd represents the `network lb delete` command.
// It deletes a load balancing rule by its ID.
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

// init registers the `network lb delete` command with the parent `network lb` command
// and binds the flags for the `lbDeleteOptions` struct.
func init() {
	NetworkLbCmd.AddCommand(NetworkLbDeleteCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbDeleteCmd, &lbDeleteOpts)
}
