package network

import (
	"fmt"
	"log/slog"
	"os"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type listServiceOfferringOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
}

var listOfferingOpts listServiceOfferringOptions

// NetworkServiceOfferingCmd is the command to list network service offerings.
var NetworkServiceOfferingCmd = &cobra.Command{
	Use:     "service-offering",
	Aliases: []string{"serviceofferings", "service-offerings", "serviceoffering", "offering", "offerings", "service-offering list"},
	Short:   "List available network service offerings for a zone",
	Long:    `List available network service offerings in a zone.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.LoadFromCobraFlags(cmd, &listOpts); err != nil {
			return err
		}
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListNetworkServiceOfferings(zoneID)
		if err != nil {
			slog.Error("failed to list network service offerings", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		if len(resp.Data) == 0 {
			fmt.Println("No network service offerings found.")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Display Name", "Price", "Overprice", "Plan(GiB)", "Rate(Mbps)", "Type", "Protocol", "Desc", "DisplayNameFA"})
		for _, offering := range resp.Data {
			row := []string{
				offering.ID,
				offering.Name,
				offering.DisplayName,
				fmt.Sprintf("%.2f", offering.HourlyStartedPrice),
				fmt.Sprintf("%.2f", offering.TrafficTransferOverprice),
				fmt.Sprintf("%d", offering.TrafficTransferPlan),
				fmt.Sprintf("%d", offering.NetworkRate),
				offering.Type,
				offering.InternetProtocol,
				offering.Description,
				offering.DisplayNameFA,
			}
			table.Append(row)
		}
		table.Render()
		return nil
	},
}

func init() {
	_ = cli.BindFlagsFromStruct(NetworkServiceOfferingCmd, &listOfferingOpts)
	NetworkCmd.AddCommand(NetworkServiceOfferingCmd)
}
