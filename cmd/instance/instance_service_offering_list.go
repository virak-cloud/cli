// Package instance provides commands for managing instances.
package instance

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"
	"virak-cli/pkg/http/responses"
)

type serviceOfferingListOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	Available   bool   `flag:"available" usage:"Show only available service offerings"`
	Columns     string `flag:"columns" usage:"Comma-separated list of columns to display"`
}

var soListOpt serviceOfferingListOptions

var validColumns = map[string]struct {
	Header string
	Value  func(offering responses.InstanceServiceOffering) string
}{
	"id":        {"ID", func(o responses.InstanceServiceOffering) string { return o.ID }},
	"name":      {"Name", func(o responses.InstanceServiceOffering) string { return o.Name }},
	"category":  {"Category", func(o responses.InstanceServiceOffering) string { return o.Category }},
	"cpu":       {"CPU", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%d", o.Hardware.CPUCore) }},
	"memory":    {"Memory (MB)", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%d", o.Hardware.MemoryMB) }},
	"storage":   {"Storage (GB)", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%d", o.Hardware.RootDiskSizeGB) }},
	"cpu_speed": {"CPU Speed (MHz)", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%d", o.Hardware.CPUSpeedMHz) }},
	"network":   {"Network Rate", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%d", o.Hardware.NetworkRate) }},
	"disk_iops": {"Disk IOPS", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%d", o.Hardware.DiskIOPS) }},
	"suggested": {"Suggested", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%t", o.Suggested) }},
	"available": {"Available", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%t", o.IsAvailable) }},
	"public":    {"Public", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%t", o.IsPublic) }},
	"image_req": {"Image Req.", func(o responses.InstanceServiceOffering) string { return fmt.Sprintf("%t", o.HasImageRequirement) }},
	"price_up": {"Price Up", func(o responses.InstanceServiceOffering) string {
		if o.HourlyPrice != nil {
			return fmt.Sprintf("%d", o.HourlyPrice.Up)
		}
		return ""
	}},
	"price_down": {"Price Down", func(o responses.InstanceServiceOffering) string {
		if o.HourlyPrice != nil {
			return fmt.Sprintf("%d", o.HourlyPrice.Down)
		}
		return ""
	}},
	"nodisc_up": {"NoDisc Up", func(o responses.InstanceServiceOffering) string {
		if o.HourlyPriceNoDiscount != nil {
			return fmt.Sprintf("%d", o.HourlyPriceNoDiscount.Up)
		}
		return ""
	}},
	"nodisc_down": {"NoDisc Down", func(o responses.InstanceServiceOffering) string {
		if o.HourlyPriceNoDiscount != nil {
			return fmt.Sprintf("%d", o.HourlyPriceNoDiscount.Down)
		}
		return ""
	}},
	"description": {"Description", func(o responses.InstanceServiceOffering) string {
		if o.Description != nil {
			return *o.Description
		}
		return ""
	}},
}

var defaultColumns = []string{"id", "name", "category", "cpu", "memory", "storage", "suggested", "available"}

var instanceServiceOfferingListCmd = &cobra.Command{
	Use:     "service-offering-list",
	Aliases: []string{"offering", "service-offering", "service-offerings"},
	Short:   "List available service offerings for instances in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}

		// Validate columns if provided
		userColsStr, _ := cmd.Flags().GetString("columns")
		if userColsStr != "" {
			userCols := SplitAndTrim(userColsStr)
			invalid := []string{}
			for _, col := range userCols {
				if _, ok := validColumns[col]; !ok {
					invalid = append(invalid, col)
				}
			}
			if len(invalid) > 0 {
				return fmt.Errorf("invalid column(s): %v", invalid)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &soListOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListInstanceServiceOfferings(zoneID)
		if err != nil {
			slog.Error("failed to list instance service offerings", "error", err, "zoneID", zoneID)
			return fmt.Errorf("failed to list instance service offerings: %w", err)
		}

		// Filter by availability if flag is set
		if soListOpt.Available {
			filtered := make([]responses.InstanceServiceOffering, 0, len(resp.Data))
			for _, offering := range resp.Data {
				if offering.IsAvailable {
					filtered = append(filtered, offering)
				}
			}
			resp.Data = filtered
		}

		// Parse columns
		selectedColumns := defaultColumns
		if soListOpt.Columns != "" {
			selectedColumns = SplitAndTrim(soListOpt.Columns)
		}
		renderInstanceServiceOfferings(resp, selectedColumns)
		return nil
	},
}

func renderInstanceServiceOfferings(resp *responses.InstanceServiceOfferingListResponse, columns []string) {
	headers := []string{}
	for _, col := range columns {
		headers = append(headers, validColumns[col].Header)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, offering := range resp.Data {
		row := []string{}
		for _, col := range columns {
			row = append(row, validColumns[col].Value(offering))
		}
		table.Append(row)
	}
	table.Render()
}

func init() {
	InstanceCmd.AddCommand(instanceServiceOfferingListCmd)
	_ = cli.BindFlagsFromStruct(instanceServiceOfferingListCmd, &soListOpt)
}
