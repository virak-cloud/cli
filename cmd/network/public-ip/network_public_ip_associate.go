package publicip

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type associateOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	NetworkID   string `flag:"networkId" usage:"Network ID to associate the public IP with"`
}

var associateOpts associateOptions

// NetworkPublicIPAssociateCmd represents the associate subcommand
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

func init() {
	NetworkPublicIPCmd.AddCommand(NetworkPublicIPAssociateCmd)
	_ = cli.BindFlagsFromStruct(NetworkPublicIPAssociateCmd, &associateOpts)
}
