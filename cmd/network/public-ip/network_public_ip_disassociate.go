package publicip

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type disassociateOptions struct {
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	NetworkID         string `flag:"networkId" usage:"Network ID to disassociate the public IP from"`
	NetworkPublicIPID string `flag:"publicIpId" usage:"Public IP ID to disassociate from the network"`
}

var disassociateOpts disassociateOptions

// NetworkPublicIPDisassociateCmd represents the disassociate subcommand
var NetworkPublicIPDisassociateCmd = &cobra.Command{
	Use:   "disassociate",
	Short: "Disassociate a public IP from a network",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("networkId"),
			cli.Required("publicIpId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &disassociateOpts); err != nil {
			return err
		}

		client := http.NewClient(token)
		resp, err := client.DisassociateNetworkPublicIp(zoneID, disassociateOpts.NetworkID, disassociateOpts.NetworkPublicIPID)
		if err != nil {
			slog.Error("failed to disassociate public IP", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if resp.Data.Success {
			slog.Info("Public IP disassociation started successfully.")
			fmt.Println("Public IP disassociation started successfully.")
		} else {
			slog.Error("public IP disassociation unsuccessful", "response", resp)
			return fmt.Errorf("failed to start public IP disassociation")
		}
		return nil
	},
}

func init() {
	NetworkPublicIPCmd.AddCommand(NetworkPublicIPDisassociateCmd)
	_ = cli.BindFlagsFromStruct(NetworkPublicIPDisassociateCmd, &disassociateOpts)
}
