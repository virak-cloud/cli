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

// rebootOptions holds the options for the `instance reboot` command.
// These options are populated from command-line flags or through the interactive mode.
type rebootOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to be rebooted.
	InstanceID  string `flag:"instance-id" usage:"ID of the instance to reboot"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Run interactive instance reboot workflow"`
}

var rebootOpt rebootOptions

// instanceRebootCmd represents the `instance reboot` command.
// It reboots a virtual machine instance in a specified zone.
// The command can be run in two modes:
// - Non-interactive: The instance ID is provided as a flag.
// - Interactive: The command prompts the user to select an instance to reboot from a list of available instances.
var instanceRebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot a running instance",
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

		if err := cli.LoadFromCobraFlags(cmd, &rebootOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)

		instanceID := rebootOpt.InstanceID

		if rebootOpt.Interactive {
			instanceListResp, err := httpClient.ListInstances(zoneID)
			if err != nil || instanceListResp == nil || len(instanceListResp.Data) == 0 {
				slog.Error("failed to fetch instances", "error", err)
				return fmt.Errorf("could not fetch instances or no instances found in this zone")
			}

			fmt.Println("Select an instance to reboot:")
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
			instanceID = selected.ID

			fmt.Printf("You have selected: %s (ID: %s)\n", selected.Name, selected.ID)
			fmt.Print("Proceed with reboot? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		resp, err := httpClient.RebootInstance(zoneID, instanceID)
		if err != nil {
			slog.Error("failed to reboot instance", "error", err, "zoneID", zoneID, "instanceID", instanceID)
			return fmt.Errorf("failed to reboot instance: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Instance reboot request accepted.")
		} else {
			fmt.Println("Instance reboot failed.")
		}
		return nil
	},
}

// init registers the `instance reboot` command with the parent `instance` command
// and binds the flags for the `rebootOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceRebootCmd)
	_ = cli.BindFlagsFromStruct(instanceRebootCmd, &rebootOpt)
}
