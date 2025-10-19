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

// stopOptions holds the options for the `instance stop` command.
// These options are populated from command-line flags or through the interactive mode.
type stopOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to be stopped.
	InstanceID  string `flag:"instance-id" usage:"ID of the instance to stop"`
	// Forced specifies whether to force stop the instance.
	Forced      bool   `flag:"forced" usage:"Force stop the instance"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Run interactive instance stop workflow"`
}

var stopOpt stopOptions

// instanceStopCmd represents the `instance stop` command.
// It stops a running virtual machine instance.
// The command can be run in two modes:
// - Non-interactive: The instance ID is provided as a flag.
// - Interactive: The command prompts the user to select a running instance to stop.
var instanceStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running instance",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}

		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("instance-id"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &stopOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		instanceID := stopOpt.InstanceID

		if stopOpt.Interactive {
			instanceListResp, err := httpClient.ListInstances(zoneID)
			if err != nil || instanceListResp == nil || len(instanceListResp.Data) == 0 {
				slog.Error("failed to fetch instances", "error", err)
				return fmt.Errorf("could not fetch instances or no instances found in this zone")
			}

			// Filter: only instances with status == "UP" are selectable
			var selectable []responses.Instance
			for _, inst := range instanceListResp.Data {
				if strings.ToUpper(inst.Status) == "UP" {
					selectable = append(selectable, inst)
				}
			}
			if len(selectable) == 0 {
				fmt.Println("No running instances (status: UP) available to stop in this zone.")
				return nil
			}

			fmt.Println("Select an instance to stop:")
			for i, inst := range selectable {
				fmt.Printf("%d) %s (ID: %s, Status: %s)\n", i+1, inst.Name, inst.ID, inst.Status)
			}
			var instIdx int
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(selectable) {
					instIdx = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			selected := selectable[instIdx]
			instanceID = selected.ID

			fmt.Printf("You have selected: %s (ID: %s)\n", selected.Name, selected.ID)
			fmt.Print("Proceed with stop? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		resp, err := httpClient.StopInstance(zoneID, instanceID, stopOpt.Forced)
		if err != nil {
			slog.Error("failed to stop instance", "error", err, "zoneID", zoneID, "instanceID", instanceID)
			return fmt.Errorf("failed to stop instance: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Instance stop request accepted.")
		} else {
			fmt.Println("Instance stop failed.")
		}
		return nil
	},
}

// init registers the `instance stop` command with the parent `instance` command
// and binds the flags for the `stopOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceStopCmd)
	_ = cli.BindFlagsFromStruct(instanceStopCmd, &stopOpt)
}
