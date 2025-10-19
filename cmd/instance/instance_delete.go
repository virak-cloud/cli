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

// deleteOptions holds the options for the `instance delete` command.
// These options are populated from command-line flags or through the interactive mode.
type deleteOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to be deleted.
	InstanceID  string `flag:"instance-id" usage:"ID of the instance to delete"`
	// Name is the name of the instance to be deleted.
	Name        string `flag:"name" usage:"Name of the instance to delete"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Run interactive instance deletion workflow"`
}

var deleteOpt deleteOptions

// instanceDeleteCmd represents the `instance delete` command.
// It deletes a virtual machine instance in a specified zone.
// The command can be run in two modes:
// - Non-interactive: The instance ID and name are provided as flags.
// - Interactive: The command prompts the user to select an instance to delete from a list of available instances.
var instanceDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an instance",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}

		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return cli.Validate(cmd,
				cli.Required("instance-id"),
				cli.Required("name"),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &deleteOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)

		if deleteOpt.Interactive {
			// Interactive deletion flow
			instanceListResp, err := httpClient.ListInstances(zoneID)
			if err != nil || instanceListResp == nil || len(instanceListResp.Data) == 0 {
				slog.Error("failed to fetch instances", "error", err)
				return fmt.Errorf("could not fetch instances or no instances found in this zone")
			}

			fmt.Println("Select an instance to delete:")
			for i, inst := range instanceListResp.Data {
				fmt.Printf("%d) %s (ID: %s, Status: %s)\n", i+1, inst.Name, inst.ID, inst.Status)
			}
			var instIdx int
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Enter number: ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				idx, err := strconv.Atoi(input)
				if err == nil && idx > 0 && idx <= len(instanceListResp.Data) {
					instIdx = idx - 1
					break
				}
				fmt.Println("Invalid selection. Try again.")
			}
			selected := instanceListResp.Data[instIdx]

			fmt.Printf("You have selected: %s (ID: %s)\n", selected.Name, selected.ID)
			fmt.Print("Proceed with deletion? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Aborted.")
				return nil
			}

			resp, err := httpClient.DeleteInstance(zoneID, selected.ID, selected.Name)
			if err != nil {
				slog.Error("failed to delete instance", "error", err, "zoneID", zoneID, "instanceID", selected.ID)
				return fmt.Errorf("failed to delete instance: %w", err)
			}
			if resp.Data.Success {
				fmt.Println("Instance deletion request accepted. Your instance will be deleted soon.")
			} else {
				fmt.Println("Instance deletion failed.")
			}
			return nil
		}

		// Non-interactive mode
		resp, err := httpClient.DeleteInstance(zoneID, deleteOpt.InstanceID, deleteOpt.Name)
		if err != nil {
			slog.Error("failed to delete instance", "error", err, "zoneID", zoneID, "instanceID", deleteOpt.InstanceID)
			return fmt.Errorf("failed to delete instance: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Instance deletion request accepted. Your instance will be deleted soon.")
		} else {
			fmt.Println("Instance deletion failed.")
		}
		return nil
	},
}

// init registers the `instance delete` command with the parent `instance` command
// and binds the flags for the `deleteOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceDeleteCmd)
	_ = cli.BindFlagsFromStruct(instanceDeleteCmd, &deleteOpt)
}
