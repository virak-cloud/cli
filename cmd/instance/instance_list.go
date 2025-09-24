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

type listOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	Columns     string `flag:"columns" usage:"Comma-separated list of columns to display"`
	ListColumns bool   `flag:"list-columns" usage:"Show all valid columns for instance list output"`
}

var listOpt listOptions

var validListColumns = map[string]struct {
	Header string
	Value  func(instance responses.Instance) string
}{
	"id": {"ID", func(i responses.Instance) string { return i.ID }},

	"customer_id": {"Customer ID", func(i responses.Instance) string { return i.CustomerID }},
	"name":        {"Name", func(i responses.Instance) string { return i.Name }},
	"zone_id":     {"Zone ID", func(i responses.Instance) string { return i.ZoneID }},
	"created":     {"Created", func(i responses.Instance) string { return fmt.Sprintf("%t", i.Created) }},
	"template_id": {"Template ID", func(i responses.Instance) string {
		if i.TemplateID != nil {
			return *i.TemplateID
		}
		return ""
	}},
	"status":          {"Status", func(i responses.Instance) string { return i.Status }},
	"instance_status": {"Instance Status", func(i responses.Instance) string { return i.InstanceStatus }},
	"password":        {"Password", func(i responses.Instance) string { return i.Password }},
	"username":        {"Username", func(i responses.Instance) string { return i.Username }},
	"created_at":      {"Created At", func(i responses.Instance) string { return fmt.Sprintf("%d", i.CreatedAt) }},
	"updated_at":      {"Updated At", func(i responses.Instance) string { return fmt.Sprintf("%d", i.UpdatedAt) }},
	"disk_offering_id": {"Disk Offering ID", func(i responses.Instance) string {
		if i.DiskOfferingID != nil {
			return *i.DiskOfferingID
		}
		return ""
	}},
	"service_offering_id": {"Service Offering ID", func(i responses.Instance) string { return i.ServiceOfferingID }},
	"kubernetes_cluster_id": {"K8s Cluster ID", func(i responses.Instance) string {
		if i.KubernetesClusterID != nil {
			return *i.KubernetesClusterID
		}
		return ""
	}},
	// vm_image fields
	"vm_image.id": {"VM Image ID", func(i responses.Instance) string {
		if i.VMImage != nil {
			return i.VMImage.ID
		}
		return ""
	}},
	"vm_image.name": {"VM Image Name", func(i responses.Instance) string {
		if i.VMImage != nil {
			return i.VMImage.Name
		}
		return ""
	}},
	"vm_image.type": {"VM Image Type", func(i responses.Instance) string {
		if i.VMImage != nil {
			return i.VMImage.Type
		}
		return ""
	}},
	"vm_image.os_type": {"VM Image OS Type", func(i responses.Instance) string {
		if i.VMImage != nil {
			return i.VMImage.OSType
		}
		return ""
	}},
	"vm_image.os_name": {"VM Image OS Name", func(i responses.Instance) string {
		if i.VMImage != nil {
			return i.VMImage.OSName
		}
		return ""
	}},
	"vm_image.os_version": {"VM Image OS Version", func(i responses.Instance) string {
		if i.VMImage != nil {
			return i.VMImage.OSVersion
		}
		return ""
	}},
	"vm_image.category": {"VM Image Category", func(i responses.Instance) string {
		if i.VMImage != nil {
			return i.VMImage.Category
		}
		return ""
	}},
	// zone fields
	"zone.id": {"Zone ID", func(i responses.Instance) string {
		if i.Zone != nil {
			return i.Zone.ID
		}
		return ""
	}},
	"zone.name": {"Zone Name", func(i responses.Instance) string {
		if i.Zone != nil {
			return i.Zone.Name
		}
		return ""
	}},
	"zone.location": {"Zone Location", func(i responses.Instance) string {
		if i.Zone != nil {
			return i.Zone.Location
		}
		return ""
	}},
	"zone.is_public": {"Zone Public", func(i responses.Instance) string {
		if i.Zone != nil {
			return fmt.Sprintf("%t", i.Zone.IsPublic)
		}
		return ""
	}},
	"zone.is_featured": {"Zone Featured", func(i responses.Instance) string {
		if i.Zone != nil {
			return fmt.Sprintf("%t", i.Zone.IsFeatured)
		}
		return ""
	}},
	"zone.is_ready": {"Zone Ready", func(i responses.Instance) string {
		if i.Zone != nil {
			return fmt.Sprintf("%t", i.Zone.IsReady)
		}
		return ""
	}},
	// service_offering fields
	"service_offering.id": {"Service Offering ID", func(i responses.Instance) string {
		if i.ServiceOffering != nil {
			return i.ServiceOffering.ID
		}
		return ""
	}},
	"service_offering.name": {"Service Offering Name", func(i responses.Instance) string {
		if i.ServiceOffering != nil {
			return i.ServiceOffering.Name
		}
		return ""
	}},
	"service_offering.category": {"Service Offering Category", func(i responses.Instance) string {
		if i.ServiceOffering != nil {
			return i.ServiceOffering.Category
		}
		return ""
	}},
	"service_offering.is_available": {"Service Offering Available", func(i responses.Instance) string {
		if i.ServiceOffering != nil {
			return fmt.Sprintf("%t", i.ServiceOffering.IsAvailable)
		}
		return ""
	}},
	"service_offering.is_public": {"Service Offering Public", func(i responses.Instance) string {
		if i.ServiceOffering != nil {
			return fmt.Sprintf("%t", i.ServiceOffering.IsPublic)
		}
		return ""
	}},
	"service_offering.suggested": {"Service Offering Suggested", func(i responses.Instance) string {
		if i.ServiceOffering != nil {
			return fmt.Sprintf("%t", i.ServiceOffering.Suggested)
		}
		return ""
	}},
	"service_offering.description": {"Service Offering Description", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.Description != nil {
			return *i.ServiceOffering.Description
		}
		return ""
	}},
	"service_offering.hourly_price.up": {"SO Price Up", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.HourlyPrice != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.HourlyPrice.Up)
		}
		return ""
	}},
	"service_offering.hourly_price.down": {"SO Price Down", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.HourlyPrice != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.HourlyPrice.Down)
		}
		return ""
	}},
	"service_offering.hourly_price_no_discount.up": {"SO NoDisc Up", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.HourlyPriceNoDiscount != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.HourlyPriceNoDiscount.Up)
		}
		return ""
	}},
	"service_offering.hourly_price_no_discount.down": {"SO NoDisc Down", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.HourlyPriceNoDiscount != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.HourlyPriceNoDiscount.Down)
		}
		return ""
	}},
	"service_offering.hardware.cpu_core": {"SO CPU Core", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.Hardware != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.Hardware.CPUCore)
		}
		return ""
	}},
	"service_offering.hardware.memory_mb": {"SO Memory MB", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.Hardware != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.Hardware.MemoryMB)
		}
		return ""
	}},
	"service_offering.hardware.cpu_speed_MHz": {"SO CPU Speed MHz", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.Hardware != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.Hardware.CPUSpeedMHz)
		}
		return ""
	}},
	"service_offering.hardware.root_disk_size_gB": {"SO Root Disk GB", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.Hardware != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.Hardware.RootDiskSizeGB)
		}
		return ""
	}},
	"service_offering.hardware.network_rate": {"SO Network Rate", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.Hardware != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.Hardware.NetworkRate)
		}
		return ""
	}},
	"service_offering.hardware.disk_iops": {"SO Disk IOPS", func(i responses.Instance) string {
		if i.ServiceOffering != nil && i.ServiceOffering.Hardware != nil {
			return fmt.Sprintf("%d", i.ServiceOffering.Hardware.DiskIOPS)
		}
		return ""
	}},
}

var defaultListColumns = []string{"id", "name", "status", "created_at"}

var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instances in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Don't require zone for --list-columns
		listCols, _ := cmd.Flags().GetBool("list-columns")
		if listCols {
			return cli.Preflight(false)(cmd, args)
		}

		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}

		// Validate columns if provided
		userColsStr, _ := cmd.Flags().GetString("columns")
		if userColsStr != "" {
			userCols := SplitAndTrim(userColsStr)
			var invalid []string
			for _, col := range userCols {
				if _, ok := validListColumns[col]; !ok {
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
		if err := cli.LoadFromCobraFlags(cmd, &listOpt); err != nil {
			return err
		}

		if listOpt.ListColumns {
			fmt.Println("Valid columns:")
			for k := range validListColumns {
				fmt.Printf("  %s\n", k)
			}
			return nil
		}

		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		instancesResponse, err := httpClient.ListInstances(zoneID)
		if err != nil {
			slog.Error("failed to list instances", "error", err, "zoneID", zoneID)
			return fmt.Errorf("failed to list instances: %w", err)
		}

		selectedColumns := defaultListColumns
		if listOpt.Columns != "" {
			selectedColumns = SplitAndTrim(listOpt.Columns)
		}
		renderInstances(instancesResponse, selectedColumns)
		return nil
	},
}

func renderInstances(instancesResponse *responses.InstanceListResponse, columns []string) {
	var headers []string
	for _, col := range columns {
		headers = append(headers, validListColumns[col].Header)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, instance := range instancesResponse.Data {
		var row []string
		for _, col := range columns {
			row = append(row, validListColumns[col].Value(instance))
		}
		table.Append(row)
	}
	table.Render()
}

func init() {
	InstanceCmd.AddCommand(instanceListCmd)
	_ = cli.BindFlagsFromStruct(instanceListCmd, &listOpt)
}
