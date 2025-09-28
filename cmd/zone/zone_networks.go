package zone

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

type zoneNetworksOptions struct {
	ZoneID string `flag:"zoneId" desc:"Zone ID to use (optional if default.zoneId is set in config, overrides positional argument if set)"`
}

var networksOpt zoneNetworksOptions

// networksCmd represents the networks command
var networksCmd = &cobra.Command{
	Use:   "networks",
	Short: "List networks for a specific zone",
	Long:  `Get a list of networks for a specific zone. You must provide the zone id as an argument or use --default-zone or --zoneId flag.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.OneOf("default-zone", "true", "false"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())
		if err := cli.LoadFromCobraFlags(cmd, &networksOpt); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		networks, err := httpClient.ListNetworks(zoneID)
		if err != nil {
			slog.Error("failed to get zone networks", "error", err, "zoneId", zoneID)
			return fmt.Errorf("error: %w", err)
		}
		slog.Info("successfully retrieved zone networks", "zoneId", zoneID, "count", len(networks.Data))
		fmt.Println("Networks for Zone:")
		for i, net := range networks.Data {
			fmt.Printf("[%d] Name: %s, ID: %s, Status: %s\n", i+1, net.Name, net.ID, net.Status)
			fmt.Printf("    Offering: %s (%s), Type: %s, Rate: %d Mbps\n", net.NetworkOffering.DisplayName, net.NetworkOffering.Name, net.NetworkOffering.Type, net.NetworkOffering.NetworkRate)
		}
		return nil
	},
}

func init() {
	ZoneCmd.AddCommand(networksCmd)
	_ = cli.BindFlagsFromStruct(networksCmd, &networksOpt)
}
