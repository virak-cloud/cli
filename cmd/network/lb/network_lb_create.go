package lb

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type lbCreateOptions struct {
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	NetworkID   string `flag:"networkId" usage:"Network ID for the load balancer"`
	PublicIPID  string `flag:"publicIpId" usage:"Public IP ID to use"`
	Name        string `flag:"name" usage:"Name of the load balancer rule"`
	Algorithm   string `flag:"algorithm" usage:"Algorithm (e.g., roundrobin)"`
	PublicPort  int    `flag:"publicPort" usage:"Public port for the load balancer"`
	PrivatePort int    `flag:"privatePort" usage:"Private port for the load balancer"`
}

var lbCreateOpts lbCreateOptions

// NetworkLbCreateCmd is the command for creating a new load balancing rule.
var NetworkLbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new load balancing rule for a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("publicIpId"),
			cli.Required("name"),
			cli.Required("algorithm"),
			cli.Required("publicPort"),
			cli.Required("privatePort"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &lbCreateOpts); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		_, err := httpClient.CreateLoadBalancerRule(
			zoneID,
			lbCreateOpts.NetworkID,
			lbCreateOpts.PublicIPID,
			lbCreateOpts.Name,
			lbCreateOpts.Algorithm,
			lbCreateOpts.PublicPort,
			lbCreateOpts.PrivatePort,
		)
		if err != nil {
			slog.Error("failed to create load balancer rule", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		slog.Info("load balancer rule created successfully")
		fmt.Println("Load balancer rule created successfully")
		return nil
	},
}

func init() {
	NetworkLbCmd.AddCommand(NetworkLbCreateCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbCreateCmd, &lbCreateOpts)
}
