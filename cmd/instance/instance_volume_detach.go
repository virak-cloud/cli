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

type volumeDetachOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use"`
	VolumeID    string `flag:"volumeId" usage:"Volume ID"`
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
	Interactive bool   `flag:"interactive" usage:"Interactively select volume and instance to detach"`
}

var volumeDetachOpt volumeDetachOptions

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

func init() {
	instanceVolumeCmd.AddCommand(instanceVolumeDetachCmd)
	_ = cli.BindFlagsFromStruct(instanceVolumeDetachCmd, &volumeDetachOpt)
}
