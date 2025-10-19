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

// snapshotDeleteOptions holds the options for the `instance snapshot delete` command.
// These options are populated from command-line flags or through the interactive mode.
type snapshotDeleteOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance the snapshot belongs to.
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
	// SnapshotID is the ID of the snapshot to be deleted.
	SnapshotID  string `flag:"snapshotId" usage:"Snapshot ID"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Interactively select instance and snapshot"`
}

var snapshotDeleteOpt snapshotDeleteOptions

// instanceSnapshotDeleteCmd represents the `instance snapshot delete` command.
// It deletes a snapshot of a virtual machine instance.
// The command can be run in two modes:
// - Non-interactive: The instance ID and snapshot ID are provided as flags.
// - Interactive: The command prompts the user to select an instance and then a snapshot to delete.
var instanceSnapshotDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a snapshot of an instance",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("instanceId"),
				cli.Required("snapshotId"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &snapshotDeleteOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		instanceID := snapshotDeleteOpt.InstanceID
		snapshotID := snapshotDeleteOpt.SnapshotID

		if snapshotDeleteOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)

			// Fetch instances
			instancesResp, err := httpClient.ListInstances(zoneID)
			if err != nil || len(instancesResp.Data) == 0 {
				return fmt.Errorf("could not fetch instances or no instances found in this zone")
			}
			fmt.Println("Select an instance:")
			for i, inst := range instancesResp.Data {
				fmt.Printf("%d) %s (ID: %s)\n", i+1, inst.Name, inst.ID)
			}
			var instIdx int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(instancesResp.Data) {
					instIdx = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			instanceID = instancesResp.Data[instIdx].ID

			// Fetch snapshots for selected instance
			snapshotsResp, err := httpClient.ShowInstance(zoneID, instanceID)
			if err != nil || len(snapshotsResp.Data.Snapshot) == 0 {
				return fmt.Errorf("could not fetch snapshots or no snapshots found for this instance")
			}
			fmt.Println("Select a snapshot to delete:")
			for i, snap := range snapshotsResp.Data.Snapshot {
				fmt.Printf("%d) %s (ID: %s, Status: %s)\n", i+1, snap.Name, snap.ID, snap.Status)
			}
			var snapIdx int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(snapshotsResp.Data.Snapshot) {
					snapIdx = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			snapshotID = snapshotsResp.Data.Snapshot[snapIdx].ID
		}

		resp, err := httpClient.DeleteInstanceSnapshot(zoneID, instanceID, snapshotID)
		if err != nil {
			slog.Error("failed to delete snapshot", "error", err, "zoneId", zoneID, "instanceId", instanceID, "snapshotId", snapshotID)
			return fmt.Errorf("failed to delete snapshot: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Snapshot deleted successfully.")
		} else {
			fmt.Println("Snapshot delete failed.")
		}
		return nil
	},
}

// init registers the `instance snapshot delete` command with the parent `instance snapshot` command
// and binds the flags for the `snapshotDeleteOptions` struct.
func init() {
	instanceSnapshotCmd.AddCommand(instanceSnapshotDeleteCmd)
	_ = cli.BindFlagsFromStruct(instanceSnapshotDeleteCmd, &snapshotDeleteOpt)
}
