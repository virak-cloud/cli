package staticnat

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"

	"github.com/spf13/cobra"
)

// disableStaticIpOptions holds the options for the `network public-ip staticnat disable` command.
// These options are populated from command-line flags.
type disableStaticIpOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network.
	NetworkID         string `flag:"networkId" usage:"Network ID"`
	// NetworkPublicIPID is the ID of the public IP to disable static NAT for.
	NetworkPublicIPID string `flag:"networkPublicIpId" usage:"Network Public IP ID"`
}

var disableOpts disableStaticIpOptions

// NetworkPublicIPStaticNatDisableCmd represents the `network public-ip staticnat disable` command.
// It disables static NAT for a public IP.
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

// init registers the `network public-ip staticnat disable` command with the parent `network public-ip staticnat` command
// and binds the flags for the `disableStaticIpOptions` struct.
func init() {
	NetworkPublicIPStaticNatCmd.AddCommand(NetworkPublicIPStaticNatDisableCmd)
	_ = cli.BindFlagsFromStruct(NetworkPublicIPStaticNatDisableCmd, &disableOpts)
}
