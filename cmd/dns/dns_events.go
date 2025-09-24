package dns

import (
	"fmt"
	"log/slog"
	"os"
	"virak-cli/internal/cli"
	"virak-cli/pkg/http"
	"virak-cli/pkg/http/responses"

	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

var dnsEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Get DNS events",
	Long:  `Get DNS events.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cli.Preflight(false)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		res, err := httpClient.GetDNSEvents()
		if err != nil {
			slog.Error("failed to get dns events", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		if len(res.Data) == 0 {
			slog.Info("No dns events found.")
			return nil
		}
		renderDNSEvents(res)
		return nil
	},
}

func renderDNSEvents(eventsResponse *responses.DNSEventsResponse) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Content", "Created At"})

	for _, event := range eventsResponse.Data {
		table.Append([]string{event.Type, event.Content, fmt.Sprintf("%d", event.CreatedAt)})
	}
	table.Render()
}

func init() {
	DnsCmd.AddCommand(dnsEventsCmd)
}
