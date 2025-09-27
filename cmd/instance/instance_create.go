package instance

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"
	"virak-cli/pkg/http/responses"
)

type instanceCreateOptions struct {
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	ServiceOfferingID string `flag:"service-offering-id" usage:"ID of the service offering"`
	VMImageID         string `flag:"vm-image-id" usage:"ID of the VM image"`
	NetworkIDsRaw     string `flag:"network-ids" usage:"JSON array of network IDs, e.g. '[\"id1\",\"id2\"]'"`
	Name              string `flag:"name" usage:"Name of the instance"`
	Interactive       bool   `flag:"interactive" usage:"Run interactive instance creation workflow"`
}

var createOpt instanceCreateOptions

var instanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new instance in a zone",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}

		// In non-interactive mode, these flags are required.
		// In interactive mode, we prompt for them, so we don't need to validate them here.
		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("service-offering-id"),
				cli.Required("vm-image-id"),
				cli.Required("network-ids"),
				cli.Required("name"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &createOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)

		if createOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)

			// Service Offering Selection
			soResp, err := httpClient.ListInstanceServiceOfferings(zoneID)
			if err != nil || soResp == nil || len(soResp.Data) == 0 {
				slog.Error("failed to fetch service offerings", "error", err)
				return fmt.Errorf("could not fetch service offerings")
			}
			var activeOfferings []responses.InstanceServiceOffering
			for _, so := range soResp.Data {
				if so.IsAvailable {
					activeOfferings = append(activeOfferings, so)
				}
			}
			if len(activeOfferings) == 0 {
				fmt.Println("No active service offerings available.")
				return nil
			}
			fmt.Println("Select a Service Offering:")
			for i, so := range activeOfferings {
				hourly := "N/A"
				if so.HourlyPrice != nil {
					hourly = fmt.Sprintf("%d", so.HourlyPrice.Up)
				}
				fmt.Printf("%d) %s (ID: %s, Hourly: %s IRR)\n", i+1, so.Name, so.ID, hourly)
			}
			var soIdx int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(activeOfferings) {
					soIdx = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			createOpt.ServiceOfferingID = activeOfferings[soIdx].ID

			// VM Image Selection
			imgResp, err := httpClient.ListInstanceVMImages(zoneID)
			if err != nil {
				slog.Error("failed to fetch VM images", "error", err)
				return fmt.Errorf("could not fetch VM images")
			}

			if imgResp == nil || len(imgResp.Data) == 0 {
				slog.Error("no VM images returned or empty list")
				return fmt.Errorf("no VM images available in this zone")
			}
			fmt.Println("Select a VM Image:")
			for i, img := range imgResp.Data {
				fmt.Printf("%d) %s (ID: %s)\n", i+1, img.Name, img.ID)
			}
			var imgIdx int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(imgResp.Data) {
					imgIdx = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			createOpt.VMImageID = imgResp.Data[imgIdx].ID

			// Network Selection
			netResp, err := httpClient.ListNetworks(zoneID)
			if err != nil || netResp == nil || len(netResp.Data) == 0 {
				slog.Error("failed to fetch networks", "error", err)
				return fmt.Errorf("could not fetch networks")
			}
			fmt.Println("Select one or more Networks (comma separated numbers):")
			for i, net := range netResp.Data {
				fmt.Printf("%d) %s (ID: %s)\n", i+1, net.Name, net.ID)
			}
			var networkIds []string
			for {
				fmt.Print("Enter numbers: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				parts := strings.Split(input, ",")
				valid := true
				ids := []string{}
				for _, p := range parts {
					p = strings.TrimSpace(p)
					idx, err := strconv.Atoi(p)
					if err != nil || idx < 1 || idx > len(netResp.Data) {
						valid = false
						break
					}
					ids = append(ids, netResp.Data[idx-1].ID)
				}
				if valid && len(ids) > 0 {
					networkIds = ids
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			networkIdsBytes, _ := json.Marshal(networkIds)
			createOpt.NetworkIDsRaw = string(networkIdsBytes)

			// Instance Name Input
			for {
				fmt.Print("Enter instance name: ")
				input, _ := reader.ReadString('\n')
				createOpt.Name = strings.TrimSpace(input)
				if createOpt.Name != "" {
					break
				}
				fmt.Println("Name cannot be empty.")
			}

			// Confirmation
			fmt.Println("\nSummary:")
			fmt.Printf("Service Offering: %s\n", activeOfferings[soIdx].Name)
			fmt.Printf("VM Image: %s\n", imgResp.Data[imgIdx].Name)
			fmt.Printf("Networks: ")
			for i, nid := range networkIds {
				for _, net := range netResp.Data {
					if net.ID == nid {
						if i > 0 {
							fmt.Print(", ")
						}
						fmt.Print(net.Name)
					}
				}
			}
			fmt.Println()
			fmt.Printf("Name: %s\n", createOpt.Name)
			fmt.Print("Proceed with creation? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		var networkIds []string
		if err := json.Unmarshal([]byte(createOpt.NetworkIDsRaw), &networkIds); err != nil {
			slog.Error("invalid --network-ids format", "error", err)
			return fmt.Errorf("--network-ids must be a JSON array of strings, e.g. '[\"id1\",\"id2\"]'")
		}

		resp, err := httpClient.CreateInstance(zoneID, createOpt.ServiceOfferingID, createOpt.VMImageID, networkIds, createOpt.Name)
		if err != nil {
			slog.Error("failed to create instance", "error", err, "zoneID", zoneID)

			// Try to parse API error response for 'errors' field
			var apiErr struct {
				Errors map[string][]string `json:"errors"`
			}
			if errBody, ok := err.(interface{ Error() string }); ok {
				errStr := errBody.Error()
				if json.Unmarshal([]byte(errStr), &apiErr) == nil && len(apiErr.Errors) > 0 {
					fmt.Println("API Error:")
					for field, msgs := range apiErr.Errors {
						for _, msg := range msgs {
							fmt.Printf("- %s: %s\n", field, msg)
						}
					}
					return fmt.Errorf("API returned errors")
				}
			}
			return fmt.Errorf("failed to create instance: %w", err)
		}
		if resp != nil && resp.Data.Success {
			fmt.Println("Instance creation request accepted. Your instance will be created soon.")
			fmt.Println("Please check the instance list to see when it becomes active.")
		} else {
			b, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println(string(b))
		}
		return nil
	},
}

func init() {
	InstanceCmd.AddCommand(instanceCreateCmd)
	_ = cli.BindFlagsFromStruct(instanceCreateCmd, &createOpt)
}
