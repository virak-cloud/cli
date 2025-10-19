package publicip

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// associateOptions holds the options for the `network public-ip associate` command.
// These options are populated from command-line flags.
type associateOptions struct {
	// ZoneID is the ID of the zone where the network is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkID is the ID of the network to associate the public IP with.
	NetworkID string `flag:"networkId" usage:"Network ID to associate the public IP with"`
}

var associateOpts associateOptions

// NetworkPublicIPAssociateCmd represents the `network public-ip associate` command.
// It associates a new public IP with a network.
var NetworkPublicIPAssociateCmd = &cobra.Command{
	Use:   "associate",
	Short: "Associate a new public IP with a network",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &associateOpts); err != nil {
			return err
		}

		client := http.NewClient(token)
		resp, err := client.AssociateNetworkPublicIp(zoneID, associateOpts.NetworkID)
		if err != nil {
			slog.Error("failed to associate public IP", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if resp.Data.Success {
			slog.Info("Public IP association started successfully.")
			fmt.Println("Public IP association started successfully.")
		} else {
			slog.Error("public IP association unsuccessful", "response", resp)
			return fmt.Errorf("failed to start public IP association")
		}
		return nil
	},
}

// init registers the `network public-ip associate` command with the parent `network public-ip` command
// and binds the flags for the `associateOptions` struct.
func init() {
	NetworkPublicIPCmd.AddCommand(NetworkPublicIPAssociateCmd)
	_ = cli.BindFlagsFromStruct(NetworkPublicIPAssociateCmd, &associateOpts)
}
