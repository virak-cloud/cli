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

type rebootOptions struct {
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	InstanceID  string `flag:"instance-id" usage:"ID of the instance to reboot"`
	Interactive bool   `flag:"interactive" usage:"Run interactive instance reboot workflow"`
}

var rebootOpt rebootOptions

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

func init() {
	InstanceCmd.AddCommand(instanceRebootCmd)
	_ = cli.BindFlagsFromStruct(instanceRebootCmd, &rebootOpt)
}
