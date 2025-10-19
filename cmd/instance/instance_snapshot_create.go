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

// instanceSnapshotCmd is the parent command for all instance snapshot related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var instanceSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage instance snapshots",
}

// snapshotCreateOptions holds the options for the `instance snapshot create` command.
// These options are populated from command-line flags or through the interactive mode.
type snapshotCreateOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to create a snapshot of.
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
	// Name is the name of the snapshot.
	Name        string `flag:"name" usage:"Snapshot name"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Prompt for required fields interactively"`
}

var snapshotCreateOpt snapshotCreateOptions

// instanceSnapshotCreateCmd represents the `instance snapshot create` command.
// It creates a new snapshot of a virtual machine instance.
// The command can be run in two modes:
// - Non-interactive: The instance ID and snapshot name are provided as flags.
// - Interactive: The command prompts the user to select an instance and provide a name for the snapshot.
var instanceSnapshotCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a snapshot of an instance",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("instanceId"),
				cli.Required("name"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &snapshotCreateOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		instanceID := snapshotCreateOpt.InstanceID
		name := snapshotCreateOpt.Name

		if snapshotCreateOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)

			// Fetch instances in zone
			instanceListResp, err := httpClient.ListInstances(zoneID)
			if err != nil || instanceListResp == nil || len(instanceListResp.Data) == 0 {
				slog.Error("failed to fetch instances", "error", err)
				return fmt.Errorf("could not fetch instances or no instances found in this zone")
			}

			// Filter to only UP instances
			var upInstances []responses.Instance
			for _, inst := range instanceListResp.Data {
				if strings.ToUpper(inst.Status) == "UP" {
					upInstances = append(upInstances, inst)
				}
			}
			if len(upInstances) == 0 {
				return fmt.Errorf("no instances with status 'UP' found in this zone. Cannot create snapshot")
			}

			// Present selection menu for UP instances only
			fmt.Println("Select an instance to snapshot (only 'UP' status):")
			for i, inst := range upInstances {
				fmt.Printf("%d) %s (ID: %s, Status: %s)\n", i+1, inst.Name, inst.ID, inst.Status)
			}
			var instIdx int
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(upInstances) {
					instIdx = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			selected := upInstances[instIdx]
			instanceID = selected.ID

			fmt.Printf("You have selected: %s (ID: %s)\n", selected.Name, selected.ID)

			// Prompt for snapshot name
			if name == "" {
				fmt.Print("Enter Snapshot Name: ")
				input, _ := reader.ReadString('\n')
				name = strings.TrimSpace(input)
			}
			if name == "" {
				return fmt.Errorf("snapshot name is required")
			}
		}

		// Check for WAITING snapshot before proceeding
		instanceDetail, err := httpClient.ShowInstance(zoneID, instanceID)
		if err != nil {
			slog.Error("failed to fetch instance details", "error", err)
			return fmt.Errorf("could not fetch instance details")
		}
		for _, snap := range instanceDetail.Data.Snapshot {
			if strings.ToUpper(snap.Status) == "WAITING" {
				return fmt.Errorf("there is already a snapshot in WAITING status for this instance. Please wait until it completes")
			}
		}

		resp, err := httpClient.CreateInstanceSnapshot(zoneID, instanceID, name)
		if err != nil {
			slog.Error("failed to create snapshot", "error", err, "zoneID", zoneID, "instanceID", instanceID, "name", name)
			return fmt.Errorf("failed to create snapshot: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Snapshot created successfully.")
		} else {
			fmt.Println("Snapshot creation failed.")
		}
		return nil
	},
}

// init registers the `instance snapshot create` command with the parent `instance snapshot` command
// and binds the flags for the `snapshotCreateOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceSnapshotCmd)
	instanceSnapshotCmd.AddCommand(instanceSnapshotCreateCmd)
	_ = cli.BindFlagsFromStruct(instanceSnapshotCreateCmd, &snapshotCreateOpt)
}
