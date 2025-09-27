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
	"virak-cli/pkg/http/responses"

	"github.com/spf13/cobra"
)

type startOptions struct {
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	InstanceID  string `flag:"instance-id" usage:"ID of the instance to start"`
	Interactive bool   `flag:"interactive" usage:"Run interactive instance start workflow"`
}

var startOpt startOptions

var instanceStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a stopped instance",
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

		if err := cli.LoadFromCobraFlags(cmd, &startOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		instanceID := startOpt.InstanceID

		if startOpt.Interactive {
			instanceListResp, err := httpClient.ListInstances(zoneID)
			if err != nil || instanceListResp == nil || len(instanceListResp.Data) == 0 {
				slog.Error("failed to fetch instances", "error", err)
				return fmt.Errorf("could not fetch instances or no instances found in this zone")
			}

			// Filter: only instances with status == "DOWN" are selectable
			var selectable []responses.Instance
			for _, inst := range instanceListResp.Data {
				if strings.ToUpper(inst.Status) == "DOWN" {
					selectable = append(selectable, inst)
				}
			}
			if len(selectable) == 0 {
				fmt.Println("No stopped instances (status: DOWN) available to start in this zone.")
				return nil
			}

			fmt.Println("Select an instance to start:")
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
			fmt.Print("Proceed with start? (y/n): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		resp, err := httpClient.StartInstance(zoneID, instanceID)
		if err != nil {
			slog.Error("failed to start instance", "error", err, "zoneID", zoneID, "instanceID", instanceID)
			return fmt.Errorf("failed to start instance: %w", err)
		}
		if resp.Data.Success {
			fmt.Println("Instance start request accepted.")
		} else {
			fmt.Println("Instance start failed.")
		}
		return nil
	},
}

func init() {
	InstanceCmd.AddCommand(instanceStartCmd)
	_ = cli.BindFlagsFromStruct(instanceStartCmd, &startOpt)
}
