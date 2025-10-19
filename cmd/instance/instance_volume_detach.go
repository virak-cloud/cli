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

// volumeDetachOptions holds the options for the `instance volume detach` command.
// These options are populated from command-line flags or through the interactive mode.
type volumeDetachOptions struct {
	// ZoneID is the ID of the zone where the volume and instance are located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// VolumeID is the ID of the volume to be detached.
	VolumeID    string `flag:"volumeId" usage:"Volume ID"`
	// InstanceID is the ID of the instance to detach the volume from.
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Interactively select volume and instance to detach"`
}

var volumeDetachOpt volumeDetachOptions

// instanceVolumeDetachCmd represents the `instance volume detach` command.
// It detaches a volume from a virtual machine instance.
// The command can be run in two modes:
// - Non-interactive: The volume ID and instance ID are provided as flags.
// - Interactive: The command prompts the user to select a volume and an instance to detach it from.
var instanceVolumeDetachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach a volume from an instance",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("volumeId"),
				cli.Required("instanceId"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &volumeDetachOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		volumeID := volumeDetachOpt.VolumeID
		instanceID := volumeDetachOpt.InstanceID

		if volumeDetachOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)

			// Select Volume
			volumesResp, err := httpClient.ListInstanceVolumes(zoneID)
			if err != nil || len(volumesResp.Data) == 0 {
				return fmt.Errorf("no volumes found or error fetching volumes")
			}
			fmt.Println("Select a volume to detach:")
			for i, v := range volumesResp.Data {
				fmt.Printf("%d) %s (ID: %s, Size: %dGB, Status: %s)\n", i+1, v.Name, v.ID, v.Size, v.Status)
			}
			var volChoice int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(volumesResp.Data) {
					volChoice = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			volumeID = volumesResp.Data[volChoice].ID

			// Select Instance
			instancesResp, err := httpClient.ListInstances(zoneID)
			if err != nil || len(instancesResp.Data) == 0 {
				return fmt.Errorf("no instances found or error fetching instances")
			}
			fmt.Println("Select an instance to detach from:")
			for i, inst := range instancesResp.Data {
				fmt.Printf("%d) %s (ID: %s, Status: %s)\n", i+1, inst.Name, inst.ID, inst.Status)
			}
			var instChoice int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(instancesResp.Data) {
					instChoice = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			instanceID = instancesResp.Data[instChoice].ID
		}

		if volumeID == "" || instanceID == "" {
			return fmt.Errorf("--volumeId and --instanceId flags are required")
		}

		resp, err := httpClient.DetachInstanceVolume(zoneID, volumeID, instanceID)
		if err != nil {
			slog.Error("failed to detach volume", "error", err, "zoneId", zoneID, "volumeId", volumeID, "instanceId", instanceID)
			return fmt.Errorf("failed to detach volume: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Volume detached successfully.")
		} else {
			fmt.Println("Volume detach failed.")
		}
		return nil
	},
}

// init registers the `instance volume detach` command with the parent `instance volume` command
// and binds the flags for the `volumeDetachOptions` struct.
func init() {
	instanceVolumeCmd.AddCommand(instanceVolumeDetachCmd)
	_ = cli.BindFlagsFromStruct(instanceVolumeDetachCmd, &volumeDetachOpt)
}
