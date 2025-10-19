package create

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// createL2NetworkOptions holds the options for the `network create l2` command.
// These options are populated from command-line flags.
type createL2NetworkOptions struct {
	// ZoneID is the ID of the zone where the network will be created.
	// This is optional if a default zone is set in the config.
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// NetworkOfferingID is the ID of the L2 network offering to use.
	NetworkOfferingID string `flag:"network-offering-id" usage:"Network offering ID"`
	// Name is the name of the new L2 network.
	Name              string `flag:"name" usage:"Network name"`
}

var l2NetworkOptions createL2NetworkOptions

// networkCreateL2Cmd represents the `network create l2` command.
// It creates a new Layer 2 network in a specified zone.
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
		
		// Validate that the network offering is of type L2
		serviceOfferings, err := httpClient.GetL2NetworkServiceOfferings(zoneID)
		if err != nil {
			slog.Error("failed to get L2 network service offerings", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		// Check if the provided network offering ID exists and is of type L2
		isValidL2Offering := false
		for _, offering := range serviceOfferings.Data {
			if offering.ID == l2NetworkOptions.NetworkOfferingID {
				if offering.Type != "L2" {
					return fmt.Errorf("network offering '%s' is not of type L2 (found type: %s)", l2NetworkOptions.NetworkOfferingID, offering.Type)
				}
				isValidL2Offering = true
				break
			}
		}

		if !isValidL2Offering {
			return fmt.Errorf("network offering ID '%s' not found or is not a valid L2 network offering", l2NetworkOptions.NetworkOfferingID)
		}

		// Call the HTTP method and handle response
		_, err = httpClient.CreateL2Network(zoneID, l2NetworkOptions.NetworkOfferingID, l2NetworkOptions.Name)
		if err != nil {
			slog.Error("failed to create L2 network", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("L2 network created successfully")
		fmt.Println("L2 network created successfully")
		return nil
	},
}

// init registers the `network create l2` command with the parent `network create` command
// and binds the flags for the `createL2NetworkOptions` struct.
func init() {
	_ = cli.BindFlagsFromStruct(networkCreateL2Cmd, &l2NetworkOptions)
	NetworkCreateCmd.AddCommand(networkCreateL2Cmd)
}
