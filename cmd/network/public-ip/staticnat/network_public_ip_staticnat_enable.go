package staticnat

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

type enableStaticIpOptions struct {
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	NetworkID         string `flag:"networkId" usage:"Network ID"`
	NetworkPublicIPID string `flag:"networkPublicIpId" usage:"Network Public IP ID"`
	InstanceID        string `flag:"instanceId" usage:"Instance ID"`
}

var enableOpts enableStaticIpOptions

// NetworkPublicIPStaticNatEnableCmd represents the enable subcommand
var NetworkPublicIPStaticNatEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable static NAT for a public IP",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling (required for this command)
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("networkPublicIpId"),
			cli.Required("instanceId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		// Load options from flags
		if err := cli.LoadFromCobraFlags(cmd, &enableOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		// Call the HTTP method and handle response
		resp, err := httpClient.EnableNetworkPublicIpStaticNat(zoneID, enableOpts.NetworkID, enableOpts.NetworkPublicIPID, enableOpts.InstanceID)
		if err != nil {
			slog.Error("failed to enable static NAT", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		if resp.Data.Success {
			slog.Info("static NAT enable operation completed successfully")
			fmt.Println("Static NAT enable started successfully.")
		} else {
			slog.Error("static NAT enable unsuccessful", "response", resp)
			return fmt.Errorf("failed to start static NAT enable")
		}
		return nil
	},
}

func init() {
	NetworkPublicIPStaticNatCmd.AddCommand(NetworkPublicIPStaticNatEnableCmd)
	_ = cli.BindFlagsFromStruct(NetworkPublicIPStaticNatEnableCmd, &enableOpts)
}
