package create

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type createL2NetworkOptions struct {
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	NetworkOfferingID string `flag:"network-offering-id" usage:"Network offering ID"`
	Name              string `flag:"name" usage:"Network name"`
}

var l2NetworkOptions createL2NetworkOptions

var networkCreateL2Cmd = &cobra.Command{
	Use:   "l2",
	Short: "Create a new L2 network in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling (required for this command)
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd,
			cli.Required("network-offering-id"),
			cli.Required("name"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &l2NetworkOptions); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		// Call the HTTP method and handle response
		_, err := httpClient.CreateL2Network(zoneID, l2NetworkOptions.NetworkOfferingID, l2NetworkOptions.Name)
		if err != nil {
			slog.Error("failed to create L2 network", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("L2 network created successfully")
		fmt.Println("L2 network created successfully")
		return nil
	},
}

func init() {
	_ = cli.BindFlagsFromStruct(networkCreateL2Cmd, &l2NetworkOptions)
	NetworkCreateCmd.AddCommand(networkCreateL2Cmd)
}
