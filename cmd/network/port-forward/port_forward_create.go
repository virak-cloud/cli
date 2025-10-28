package portforward

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/virak-cloud/cli/internal/cli"

	"github.com/spf13/cobra"
)

type portForwardCreateOptions struct {
	ZoneID      string `flag:"zoneId" desc:"Zone ID (optional if default.zoneId is set in config)"`
	NetworkID   string `flag:"networkId" desc:"Network ID (required)"`
	Protocol    string `flag:"protocol" desc:"Protocol (TCP/UDP) [required]"`
	PublicPort  int    `flag:"publicPort" desc:"Public port [required]"`
	PrivatePort int    `flag:"privatePort" desc:"Private port [required]"`
	PrivateIP   string `flag:"privateIp" desc:"Private IP address [required]"`
}

var portForwardCreateOpts portForwardCreateOptions

var NetworkPortForwardCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a port forwarding rule for a network",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.IsUlid("networkId"),
			cli.Required("protocol"),
			cli.OneOf("protocol", "TCP", "UDP"),
			cli.Required("publicPort"),
			cli.Required("privatePort"),
			cli.Required("privateIp"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneId := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &portForwardCreateOpts); err != nil {
			return err
		}

		request := map[string]interface{}{
			"network_id":  portForwardCreateOpts.NetworkID,
			"protocol":    portForwardCreateOpts.Protocol,
			"public_port": portForwardCreateOpts.PublicPort,
			"private_port": portForwardCreateOpts.PrivatePort,
			"private_ip":   portForwardCreateOpts.PrivateIP,
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.CreatePortForward(zoneId, request)
		if err != nil {
			slog.Error("failed to create port forwarding rule", "error", err)
			return fmt.Errorf("failed to create port forwarding rule: %w", err)
		}

		if resp.Data.Success {
			fmt.Println("Port forwarding rule created successfully.")
		} else {
			fmt.Println("Failed to create port forwarding rule.")
		}
		return nil
	},
}

func init() {
	NetworkPortForwardCmd.AddCommand(NetworkPortForwardCreateCmd)
	_ = cli.BindFlagsFromStruct(NetworkPortForwardCreateCmd, &portForwardCreateOpts)
}