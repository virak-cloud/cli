package zone

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"
)

type zoneServicesOptions struct {
	ZoneID string `flag:"zoneId" desc:"Zone ID to use (optional if default.zoneId is set in config, overrides positional argument if set)"`
}

var servicesOpt zoneServicesOptions

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Show active services for a specific zone",
	Long:  `Get a list of active services for a specific zone from Virak API. You must provide the zone ID as an argument or use --default-zone flag.`,
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
		if err := cli.LoadFromCobraFlags(cmd, &servicesOpt); err != nil {
			return err
		}
		httpClient := http.NewClient(token)
		services, err := httpClient.GetZoneActiveServices(zoneID)
		if err != nil {
			slog.Error("failed to get zone active services", "error", err, "zoneId", zoneID)
			return fmt.Errorf("error: %w", err)
		}
		slog.Info("successfully retrieved zone active services", "zoneId", zoneID)
		fmt.Println("Active Services for Zone:")
		fmt.Printf("  Instance: %v\n", services.Instance)
		fmt.Printf("  DataVolume: %v\n", services.DataVolume)
		fmt.Printf("  Network: %v\n", services.Network)
		fmt.Printf("  ObjectStorage: %v\n", services.ObjectStorage)
		fmt.Printf("  K8s: %v\n", services.K8s)
		return nil
	},
}

func init() {
	ZoneCmd.AddCommand(servicesCmd)
	_ = cli.BindFlagsFromStruct(servicesCmd, &servicesOpt)
}
