package staticnat

import (
	"fmt"
	"log/slog"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type disableStaticIpOptions struct {
	DefaultZone       bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use"`
	NetworkID         string `flag:"networkId" usage:"Network ID"`
	NetworkPublicIPID string `flag:"networkPublicIpId" usage:"Network Public IP ID"`
}

var disableOpts disableStaticIpOptions

// NetworkPublicIPStaticNatDisableCmd represents the disable subcommand
var NetworkPublicIPStaticNatDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable static NAT for a public IP",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling (required for this command)
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("networkPublicIpId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		// Load options from flags
		if err := cli.LoadFromCobraFlags(cmd, &disableOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		// Call the HTTP method and handle response
		resp, err := httpClient.DisableNetworkPublicIpStaticNat(zoneID, disableOpts.NetworkID, disableOpts.NetworkPublicIPID)
		if err != nil {
			slog.Error("failed to disable static NAT", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		if resp.Data.Success {
			slog.Info("static NAT disable operation completed successfully")
			fmt.Println("Static NAT disable started successfully.")
		} else {
			slog.Error("static NAT disable unsuccessful", "response", resp)
			return fmt.Errorf("failed to start static NAT disable")
		}
		return nil
	},
}

func init() {
	NetworkPublicIPStaticNatCmd.AddCommand(NetworkPublicIPStaticNatDisableCmd)
	_ = cli.BindFlagsFromStruct(NetworkPublicIPStaticNatDisableCmd, &disableOpts)
}
