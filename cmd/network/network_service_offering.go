package network

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

type listServiceOfferringOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	Type   string `flag:"type" usage:"Filter by service offering type: l2, l3, or all (default: all)"`
}

var listOfferingOpts listServiceOfferringOptions

// NetworkServiceOfferingCmd is the command to list network service offerings.
var NetworkServiceOfferingCmd = &cobra.Command{
	Use:     "service-offering",
	Aliases: []string{"serviceofferings", "service-offerings", "serviceoffering", "offering", "offerings", "service-offering list"},
	Short:   "List available network service offerings for a zone",
	Long:    `List available network service offerings in a zone. Filter by type using --type flag (l2, l3, or all).`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("type"),
			cli.OneOf("type", "l2", "l3", "all"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.LoadFromCobraFlags(cmd, &listOfferingOpts); err != nil {
			return err
		}
		// Validate type argument
		if listOfferingOpts.Type != "" {
			lowerType := strings.ToLower(listOfferingOpts.Type)
			if lowerType != "l2" && lowerType != "l3" && lowerType != "all" {
				return fmt.Errorf("invalid type '%s'. Must be one of: l2, l3, all", listOfferingOpts.Type)
			}
			listOfferingOpts.Type = lowerType
		} else {
			listOfferingOpts.Type = "all"
		}
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		var resp *responses.NetworkServiceOfferingListResponse
		var err error

		// Call appropriate HTTP method based on type
		switch listOfferingOpts.Type {
		case "l2":
			resp, err = httpClient.GetL2NetworkServiceOfferings(zoneID)
			if err != nil {
				slog.Error("failed to list L2 network service offerings", "error", err)
				return fmt.Errorf("error: %w", err)
			}
		case "l3":
			resp, err = httpClient.GetL3NetworkServiceOfferings(zoneID)
			if err != nil {
				slog.Error("failed to list L3 network service offerings", "error", err)
				return fmt.Errorf("error: %w", err)
			}
		default: // "all"
			resp, err = httpClient.ListNetworkServiceOfferings(zoneID)
			if err != nil {
				slog.Error("failed to list network service offerings", "error", err)
				return fmt.Errorf("error: %w", err)
			}
		}

		if len(resp.Data) == 0 {
			switch listOfferingOpts.Type {
			case "l2":
				fmt.Println("No L2 network service offerings found.")
			case "l3":
				fmt.Println("No L3 network service offerings found.")
			default:
				fmt.Println("No network service offerings found.")
			}
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
