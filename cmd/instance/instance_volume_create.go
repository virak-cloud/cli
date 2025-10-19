package instance

import (
	"bufio"
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// instanceVolumeCmd is the parent command for all instance volume related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var instanceVolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Manage instance volumes",
}

// volumeCreateOptions holds the options for the `instance volume create` command.
// These options are populated from command-line flags or through the interactive mode.
type volumeCreateOptions struct {
	// ZoneID is the ID of the zone where the volume will be created.
	// This is optional if a default zone is set in the config.
	ZoneID            string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// ServiceOfferingID is the ID of the service offering for the volume.
	ServiceOfferingID string `flag:"serviceOfferingId" usage:"Service Offering ID"`
	// Size is the size of the volume in GB.
	Size              int    `flag:"size" usage:"Volume size (GB)"`
	// Name is the name of the volume.
	Name              string `flag:"name" usage:"Volume name"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive       bool   `flag:"interactive" usage:"Interactively select service offering, size, and name"`
}

var volumeCreateOpt volumeCreateOptions

// instanceVolumeCreateCmd represents the `instance volume create` command.
// It creates a new data volume for an instance.
// The command can be run in two modes:
// - Non-interactive: The service offering ID, size, and name are provided as flags.
// - Interactive: The command prompts the user to select a service offering, and to provide a size and name for the volume.
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

// init registers the `instance volume create` command with the parent `instance volume` command
// and binds the flags for the `volumeCreateOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceVolumeCmd)
	instanceVolumeCmd.AddCommand(instanceVolumeCreateCmd)
	_ = cli.BindFlagsFromStruct(instanceVolumeCreateCmd, &volumeCreateOpt)
}
