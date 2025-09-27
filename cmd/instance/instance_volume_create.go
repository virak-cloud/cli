package instance

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

var instanceVolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Manage instance volumes",
}

type volumeCreateOptions struct {
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	ServiceOfferingID string `flag:"serviceOfferingId" usage:"Service Offering ID"`
	Size              int    `flag:"size" usage:"Volume size (GB)"`
	Name              string `flag:"name" usage:"Volume name"`
	Interactive       bool   `flag:"interactive" usage:"Interactively select service offering, size, and name"`
}

var volumeCreateOpt volumeCreateOptions

var instanceVolumeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new data volume",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("serviceOfferingId"),
				cli.Required("size"),
				cli.Required("name"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &volumeCreateOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		serviceOfferingID := volumeCreateOpt.ServiceOfferingID
		size := volumeCreateOpt.Size
		name := volumeCreateOpt.Name

		if volumeCreateOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)
			serviceOfferingsResp, err := httpClient.ListInstanceVolumeServiceOfferings(zoneID)
			if err != nil || len(serviceOfferingsResp.Data) == 0 {
				return fmt.Errorf("no volume service offerings found or error fetching offerings")
			}
			fmt.Println("Select a volume service offering:")
			for i, so := range serviceOfferingsResp.Data {
				fmt.Printf("%d) %s (ID: %s)\n", i+1, so.Name, so.ID)
			}
			var soChoice int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(serviceOfferingsResp.Data) {
					soChoice = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			serviceOfferingID = serviceOfferingsResp.Data[soChoice].ID

			if size == 0 {
				for {
					fmt.Print("Enter volume size (GB): ")
					input, _ := reader.ReadString('\n')
					input = strings.TrimSpace(input)
					s, err := strconv.Atoi(input)
					if err == nil && s > 0 {
						size = s
						break
					}
					fmt.Println("Invalid size. Must be a number greater than 0.")
				}
			}

			if name == "" {
				for {
					fmt.Print("Enter volume name: ")
					input, _ := reader.ReadString('\n')
					name = strings.TrimSpace(input)
					if name != "" {
						break
					}
					fmt.Println("Volume name cannot be empty.")
				}
			}
		}

		if serviceOfferingID == "" || size <= 0 || name == "" {
			return fmt.Errorf("--serviceOfferingId, --size (>0), and --name flags are required in non-interactive mode")
		}

		resp, err := httpClient.CreateInstanceVolume(zoneID, serviceOfferingID, size, name)
		if err != nil {
			slog.Error("failed to create volume", "error", err, "zoneID", zoneID, "serviceOfferingID", serviceOfferingID, "size", size, "name", name)
			return fmt.Errorf("failed to create volume: %w", err)
		}
		fmt.Printf("Volume created: ID=%s, Name=%s, Size=%d, Status=%s\n",
			resp.Data.ID, resp.Data.Name, resp.Data.Size, resp.Data.Status)
		return nil
	},
}

func init() {
	InstanceCmd.AddCommand(instanceVolumeCmd)
	instanceVolumeCmd.AddCommand(instanceVolumeCreateCmd)
	_ = cli.BindFlagsFromStruct(instanceVolumeCreateCmd, &volumeCreateOpt)
}
