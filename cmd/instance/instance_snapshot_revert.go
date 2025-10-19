package instance

import (
	"bufio"
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// snapshotRevertOptions holds the options for the `instance snapshot revert` command.
// These options are populated from command-line flags or through the interactive mode.
type snapshotRevertOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to revert.
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
	// SnapshotID is the ID of the snapshot to revert to.
	SnapshotID  string `flag:"snapshotId" usage:"Snapshot ID"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Interactively select instance and snapshot"`
}

var snapshotRevertOpt snapshotRevertOptions

// instanceSnapshotRevertCmd represents the `instance snapshot revert` command.
// It reverts a virtual machine instance to a previous snapshot.
// The command can be run in two modes:
// - Non-interactive: The instance ID and snapshot ID are provided as flags.
// - Interactive: The command prompts the user to select an instance and then a snapshot to revert to.
var instanceSnapshotRevertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Revert an instance to a snapshot",
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

		if err := cli.LoadFromCobraFlags(cmd, &snapshotRevertOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		instanceID := snapshotRevertOpt.InstanceID
		snapshotID := snapshotRevertOpt.SnapshotID

		if snapshotRevertOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)
			instancesResp, err := httpClient.ListInstances(zoneID)
			if err != nil || len(instancesResp.Data) == 0 {
				return fmt.Errorf("could not fetch instances or no instances found in this zone")
			}
			fmt.Println("Select an instance:")
			for i, inst := range instancesResp.Data {
				fmt.Printf("%d) %s (ID: %s, Status: %s)\n", i+1, inst.Name, inst.ID, inst.Status)
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

			// Filter READY snapshots
			var readySnapshots []responses.InstanceSnapshot
			for _, snap := range instancesResp.Data[instIdx].Snapshot {
				if snap.Status == "READY" {
					readySnapshots = append(readySnapshots, snap)
				}
			}
			if len(readySnapshots) == 0 {
				return fmt.Errorf("no READY snapshots available for this instance")
			}
			fmt.Println("Select a snapshot to revert to:")
			for i, snap := range readySnapshots {
				fmt.Printf("%d) %s (ID: %s, CreatedAt: %d)\n", i+1, snap.Name, snap.ID, snap.CreatedAt)
			}
			var snapIdx int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(readySnapshots) {
					snapIdx = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			snapshotID = readySnapshots[snapIdx].ID
		}

		resp, err := httpClient.RevertInstanceSnapshot(zoneID, instanceID, snapshotID)
		if err != nil {
			slog.Error("failed to revert snapshot", "error", err, "zoneId", zoneID, "instanceId", instanceID, "snapshotId", snapshotID)
			return fmt.Errorf("failed to revert snapshot: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Instance going to revert to snapshot. This may take a few minutes.")
		} else {
			fmt.Println("Instance revert to snapshot failed.")
		}
		return nil
	},
}

// init registers the `instance snapshot revert` command with the parent `instance snapshot` command
// and binds the flags for the `snapshotRevertOptions` struct.
func init() {
	instanceSnapshotCmd.AddCommand(instanceSnapshotRevertCmd)
	_ = cli.BindFlagsFromStruct(instanceSnapshotRevertCmd, &snapshotRevertOpt)
}
