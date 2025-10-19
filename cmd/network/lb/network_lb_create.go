package lb

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// lbCreateOptions holds the options for the `network lb create` command.
// These options are populated from command-line flags.
type lbCreateOptions struct {
	// ZoneID is the ID of the zone where the load balancer will be created.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to create the load balancer in.
	NetworkID   string `flag:"networkId" usage:"Network ID for the load balancer"`
	// PublicIPID is the ID of the public IP to associate with the load balancer.
	PublicIPID  string `flag:"publicIpId" usage:"Public IP ID to use"`
	// Name is the name of the load balancer rule.
	Name        string `flag:"name" usage:"Name of the load balancer rule"`
	// Algorithm is the load balancing algorithm to use.
	Algorithm   string `flag:"algorithm" usage:"Algorithm (e.g., roundrobin)"`
	// PublicPort is the public port for the load balancer.
	PublicPort  int    `flag:"publicPort" usage:"Public port for the load balancer"`
	// PrivatePort is the private port to forward traffic to.
	PrivatePort int    `flag:"privatePort" usage:"Private port for the load balancer"`
}

var lbCreateOpts lbCreateOptions

// NetworkLbCreateCmd represents the `network lb create` command.
// It creates a new load balancing rule for a network.
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

// init registers the `network lb create` command with the parent `network lb` command
// and binds the flags for the `lbCreateOptions` struct.
func init() {
	NetworkLbCmd.AddCommand(NetworkLbCreateCmd)
	_ = cli.BindFlagsFromStruct(NetworkLbCreateCmd, &lbCreateOpts)
}
