package cluster

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// serviceOfferingsListOptions holds the options for the `cluster offering` command.
// These options are populated from command-line flags.
type serviceOfferingsListOptions struct {
	// ZoneID is the ID of the zone to list service offerings from.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
}

var serviceOfferingsListOpts serviceOfferingsListOptions

// kubernetesServiceOfferingsListCmd represents the `cluster offering` command.
// It lists all available Kubernetes service offerings in a specified zone.
var kubernetesServiceOfferingsListCmd = &cobra.Command{
	Use:     "offering",
	Aliases: []string{"service-offerings", "offerings-list", "offerings", "service-offerings-list"},
	Short:   "List available kubernetes service offerings",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		offerings, err := httpClient.GetKubernetesServiceOfferings(zoneID)
		if err != nil {
			slog.Error("failed to get kubernetes service offerings", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID", "Name", "Available", "CPU Cores", "CPU MHz", "RAM MB", "Disk GB",
			"Network Rate", "Price Up (per hour)", "Price Down (per hour)",
		})

		for _, o := range offerings.Data {
			table.Append([]string{
				o.ID,
				o.Name,
				fmt.Sprintf("%t", o.IsAvailable),
				fmt.Sprintf("%d", o.Hardware.CPUCore),
				fmt.Sprintf("%d", o.Hardware.CPUSpeedMHz),
				fmt.Sprintf("%d", o.Hardware.MemoryMB),
				fmt.Sprintf("%d", o.Hardware.RootDiskSizeGB),
				fmt.Sprintf("%d", o.Hardware.NetworkRate),
				fmt.Sprintf("%d", o.HourlyPrice.Up),
				fmt.Sprintf("%d", o.HourlyPrice.Down),
			})
		}
		table.Render()

		return nil
	},
}

// init registers the `cluster offering` command with the parent `cluster` command
// and binds the flags for the `serviceOfferingsListOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesServiceOfferingsListCmd)
	_ = cli.BindFlagsFromStruct(kubernetesServiceOfferingsListCmd, &serviceOfferingsListOpts)
}
