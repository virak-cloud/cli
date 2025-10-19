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

// volumeDeleteOptions holds the options for the `instance volume delete` command.
// These options are populated from command-line flags or through the interactive mode.
type volumeDeleteOptions struct {
	// ZoneID is the ID of the zone where the volume is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// VolumeID is the ID of the volume to be deleted.
	VolumeID    string `flag:"volumeId" usage:"Volume ID"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Interactively select volume to delete"`
}

var volumeDeleteOpt volumeDeleteOptions

// instanceVolumeDeleteCmd represents the `instance volume delete` command.
// It deletes a data volume.
// The command can be run in two modes:
// - Non-interactive: The volume ID is provided as a flag.
// - Interactive: The command prompts the user to select a volume to delete.
var instanceVolumeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a volume",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("volumeId"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &volumeDeleteOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		volumeID := volumeDeleteOpt.VolumeID

		if volumeDeleteOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)
			volumesResp, err := httpClient.ListInstanceVolumes(zoneID)
			if err != nil || len(volumesResp.Data) == 0 {
				return fmt.Errorf("no volumes found or error fetching volumes")
			}
			fmt.Println("Select a volume to delete:")
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
		}

		if volumeID == "" {
			return fmt.Errorf("--volumeId flag is required")
		}

		resp, err := httpClient.DeleteInstanceVolume(zoneID, volumeID)
		if err != nil {
			slog.Error("failed to delete volume", "error", err, "zoneID", zoneID, "volumeID", volumeID)
			return fmt.Errorf("failed to delete volume: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Volume deleted successfully.")
		} else {
			fmt.Println("Volume delete failed.")
		}
		return nil
	},
}

// init registers the `instance volume delete` command with the parent `instance volume` command
// and binds the flags for the `volumeDeleteOptions` struct.
func init() {
	instanceVolumeCmd.AddCommand(instanceVolumeDeleteCmd)
	_ = cli.BindFlagsFromStruct(instanceVolumeDeleteCmd, &volumeDeleteOpt)
}
