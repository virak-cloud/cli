package create

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type createL3NetworkOptions struct {
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	NetworkOfferingID string `flag:"network-offering-id" usage:"Network offering ID"`
	Name              string `flag:"name" usage:"Network name"`
	Gateway           string `flag:"gateway" usage:"Gateway IP address"`
	Netmask           string `flag:"netmask" usage:"Netmask"`
}

var l3NetworkOptions createL3NetworkOptions

var networkCreateL3Cmd = &cobra.Command{
	Use:   "l3",
	Short: "Create a new L3 network in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling (required for this command)
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd,
			cli.Required("network-offering-id"),
			cli.IsUlid("network-offering-id"),
			cli.Required("name"),
			cli.Required("gateway"),
			cli.Required("netmask"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		// Load l2NetworkOptions from flags

		if err := cli.LoadFromCobraFlags(cmd, &l3NetworkOptions); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		// Call the HTTP method and handle response
		_, err := httpClient.CreateL3Network(zoneID, l3NetworkOptions.NetworkOfferingID, l3NetworkOptions.Name, l3NetworkOptions.Gateway, l3NetworkOptions.Netmask)
		if err != nil {
			slog.Error("failed to create L3 network", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("L3 network created successfully")
		fmt.Println("L3 network created successfully")
		return nil
	},
}

func init() {
	_ = cli.BindFlagsFromStruct(networkCreateL3Cmd, &l3NetworkOptions)
	NetworkCreateCmd.AddCommand(networkCreateL3Cmd)
}
