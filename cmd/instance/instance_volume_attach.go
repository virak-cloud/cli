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

type volumeAttachOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	VolumeID    string `flag:"volumeId" usage:"Volume ID"`
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
	Interactive bool   `flag:"interactive" usage:"Interactively select volume and instance"`
}

var volumeAttachOpt volumeAttachOptions

var instanceVolumeAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a volume to an instance",
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

		if err := cli.LoadFromCobraFlags(cmd, &volumeAttachOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		volumeID := volumeAttachOpt.VolumeID
		instanceID := volumeAttachOpt.InstanceID

		if volumeAttachOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)

			// Select Volume
			volumesResp, err := httpClient.ListInstanceVolumes(zoneID)
			if err != nil || len(volumesResp.Data) == 0 {
				return fmt.Errorf("no volumes found or error fetching volumes")
			}
			fmt.Println("Select a volume to attach:")
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
			fmt.Println("Select an instance:")
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

			// Check for snapshots
			snapShotCount := len(instancesResp.Data[instChoice].Snapshot)
			if snapShotCount != 0 {
				return fmt.Errorf("cannot attach a volume to an instance that has existing snapshots")
			}
		}

		if volumeID == "" || instanceID == "" {
			return fmt.Errorf("--volumeId and --instanceId flags are required")
		}

		resp, err := httpClient.AttachInstanceVolume(zoneID, volumeID, instanceID)
		if err != nil {
			slog.Error("failed to attach volume", "error", err, "zoneId", zoneID, "volumeId", volumeID, "instanceId", instanceID)
			return fmt.Errorf("failed to attach volume: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Volume attached successfully.")
		} else {
			fmt.Println("Volume attach failed.")
		}
		return nil
	},
}

func init() {
	instanceVolumeCmd.AddCommand(instanceVolumeAttachCmd)
	_ = cli.BindFlagsFromStruct(instanceVolumeAttachCmd, &volumeAttachOpt)
}
