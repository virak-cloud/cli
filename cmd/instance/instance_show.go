package instance

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"github.com/virak-cloud/cli/pkg/http/responses"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// showOptions holds the options for the `instance show` command.
// These options are populated from command-line flags or through the interactive mode.
type showOptions struct {
	// ZoneID is the ID of the zone where the instance is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// InstanceID is the ID of the instance to be shown.
	InstanceID  string `flag:"instanceId" usage:"Instance ID"`
	// Interactive specifies whether to run the command in interactive mode.
	Interactive bool   `flag:"interactive" usage:"Run interactive instance selection workflow"`
}

var showOpt showOptions

// instanceShowCmd represents the `instance show` command.
// It shows the details of a specific virtual machine instance.
// The command can be run in two modes:
// - Non-interactive: The instance ID is provided as a flag.
// - Interactive: The command prompts the user to select an instance to show from a list of available instances.
var instanceShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of an instance",
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

		if err := cli.LoadFromCobraFlags(cmd, &showOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		instanceID := showOpt.InstanceID

		if showOpt.Interactive {
			listResp, err := httpClient.ListInstances(zoneID)
			if err != nil {
				return fmt.Errorf("could not fetch instance list: %w", err)
			}
			if len(listResp.Data) == 0 {
				fmt.Println("No instances found in this zone.")
				slog.Error("No instances found in this zone.")
				return nil
			}
			fmt.Println("Select an instance to show details:")
			for i, inst := range listResp.Data {
				fmt.Printf("[%d] %s (%s)\n", i+1, inst.Name, inst.ID)
			}
			var selection int
			for {
				fmt.Print("Enter number: ")
				buf := make([]byte, 100)
				n, _ := os.Stdin.Read(buf)
				s := strings.TrimSpace(string(buf[:n]))
				idx, err := strconv.Atoi(s)
				if err == nil && idx > 0 && idx <= len(listResp.Data) {
					selection = idx
					break
				}
				fmt.Println("Invalid selection. Please enter a valid number.")
			}
			instanceID = listResp.Data[selection-1].ID
		}

		resp, err := httpClient.ShowInstance(zoneID, instanceID)
		if err != nil {
			return fmt.Errorf("could not fetch instance details: %w", err)
		}
		if resp.Data.ID == "" {
			return fmt.Errorf("instance not found")
		}
		renderInstanceDetails(resp.Data)
		return nil
	},
}

// renderInstanceDetails renders a table with the details of an instance.
func renderInstanceDetails(inst responses.Instance) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.Append([]string{"ID", inst.ID})
	table.Append([]string{"Name", inst.Name})
	table.Append([]string{"Status", inst.Status})
	table.Append([]string{"Instance Status", inst.InstanceStatus})
	table.Append([]string{"Zone ID", inst.ZoneID})
	table.Append([]string{"Created At", fmt.Sprintf("%d", inst.CreatedAt)})
	table.Append([]string{"Updated At", fmt.Sprintf("%d", inst.UpdatedAt)})
	table.Append([]string{"Username", inst.Username})
	table.Append([]string{"Password", inst.Password})
	if inst.VMImage != nil {
		table.Append([]string{"VM Image Name", inst.VMImage.Name})
		table.Append([]string{"VM Image OS", inst.VMImage.OSName})
	}
	if inst.ServiceOffering != nil {
		table.Append([]string{"Service Offering", inst.ServiceOffering.Name})
	}
	table.Render()
}

// init registers the `instance show` command with the parent `instance` command
// and binds the flags for the `showOptions` struct.
func init() {
	InstanceCmd.AddCommand(instanceShowCmd)
	_ = cli.BindFlagsFromStruct(instanceShowCmd, &showOpt)
}
