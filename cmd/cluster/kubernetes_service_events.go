package cluster

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// serviceEventsOptions holds the options for the `cluster service-events` command.
// These options are populated from command-line flags.
type serviceEventsOptions struct {
	// ZoneID is the ID of the zone to list service events from.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
}

var serviceEventsOpts serviceEventsOptions

// kubernetesServiceEventsCmd represents the `cluster service-events` command.
// It lists all Kubernetes service events in a specified zone.
var kubernetesServiceEventsCmd = &cobra.Command{
	Use:   "service-events",
	Short: "List kubernetes service events",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		events, err := httpClient.GetKubernetesServiceEvents(zoneID)
		if err != nil {
			slog.Error("failed to get kubernetes service events", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Message", "Timestamp"})

		for _, e := range events.Data {
			table.Append([]string{e.ID, e.Message, fmt.Sprintf("%d", e.Timestamp)})
		}
		table.Render()

		return nil
	},
}

// init registers the `cluster service-events` command with the parent `cluster` command
// and binds the flags for the `serviceEventsOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesServiceEventsCmd)
	_ = cli.BindFlagsFromStruct(kubernetesServiceEventsCmd, &serviceEventsOpts)
}
