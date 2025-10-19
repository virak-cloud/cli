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

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// snapshotListOptions holds the options for the `instance snapshot list` command.
// These options are populated from command-line flags or through the interactive mode.
type snapshotListOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to list snapshots for.
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Interactively select instance"`
}

var snapshotListOpt snapshotListOptions

// instanceSnapshotListCmd represents the `instance snapshot list` command.
// It lists all snapshots of a virtual machine instance.
// The command can be run in two modes:
// - Non-interactive: The instance ID is provided as a flag.
// - Interactive: The command prompts the user to select an instance to list snapshots for.
var instanceSnapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List snapshots of an instance",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("instanceId"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &snapshotListOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		instanceID := snapshotListOpt.InstanceID

		if snapshotListOpt.Interactive {
			reader := bufio.NewReader(os.Stdin)
			instancesResp, err := httpClient.ListInstances(zoneID)
			if err != nil || len(instancesResp.Data) == 0 {
				slog.Error("could not fetch instances or no instances found in this zone")
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
			renderInstanceSnapshots(instancesResp.Data[instIdx].Snapshot)
			return nil
		}

		// Non-interactive mode
		instancesResp, err := httpClient.ListInstances(zoneID)
		if err != nil || len(instancesResp.Data) == 0 {
			return fmt.Errorf("could not fetch instances or no instances found in this zone")
		}
		var found bool
		for _, inst := range instancesResp.Data {
			if inst.ID == instanceID {
				renderInstanceSnapshots(inst.Snapshot)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("instance not found in this zone")
		}
		return nil
	},
}

// renderInstanceSnapshots renders a table of instance snapshots.
func renderInstanceSnapshots(snapshots []responses.InstanceSnapshot) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Status", "CreatedAt", "Current", "ParentID"})
	for _, snap := range snapshots {
		created := fmt.Sprintf("%d", snap.CreatedAt)
		current := fmt.Sprintf("%t", snap.Current)
		parent := ""
		if snap.ParentID != nil {
			parent = *snap.ParentID
		}
		table.Append([]string{
			snap.ID,
			snap.Name,
			snap.Status,
			created,
			current,
			parent,
		})
	}
	table.Render()
}

// init registers the `instance snapshot list` command with the parent `instance snapshot` command
// and binds the flags for the `snapshotListOptions` struct.
func init() {
	instanceSnapshotCmd.AddCommand(instanceSnapshotListCmd)
	_ = cli.BindFlagsFromStruct(instanceSnapshotListCmd, &snapshotListOpt)
}
