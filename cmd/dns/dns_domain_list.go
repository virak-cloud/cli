package dns

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

// domainListCmd represents the `dns domain list` command.
// It lists all DNS domains for the current user.
var domainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all domains",
	Long:  "List all domains",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cli.Preflight(false)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		resp, err := httpClient.GetDomains()
		if err != nil {
			slog.Error("failed to get domains", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		renderDomainList(resp)
		return nil
	},
}

// renderDomainList renders a table of DNS domains.
func renderDomainList(resp *responses.DomainList) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "Status"})
	for _, domain := range resp.Data {
		table.Append([]string{domain.Domain, domain.Status})
	}
	table.Render()
}

// init registers the `dns domain list` command with the parent `dns domain` command.
func init() {
	domainCmd.AddCommand(domainListCmd)
}
