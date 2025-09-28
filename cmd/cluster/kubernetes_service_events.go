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

type serviceEventsOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
}

var serviceEventsOpts serviceEventsOptions

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

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesServiceEventsCmd)
	_ = cli.BindFlagsFromStruct(kubernetesServiceEventsCmd, &serviceEventsOpts)
}
