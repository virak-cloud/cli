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

// showOptions holds the options for the `dns domain show` command.
// These options are populated from command-line flags.
type showOptions struct {
	// Domain is the name of the domain to be shown.
	Domain string `flag:"domain" usage:"Domain name to show"`
}

var showOpts showOptions

// domainShowCmd represents the `dns domain show` command.
// It shows the details of a specific DNS domain.
var domainShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a domain",
	Long:  "Show a domain",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd, cli.Required("domain"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &showOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.GetDomain(showOpts.Domain)
		if err != nil {
			slog.Error("failed to get domain", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		renderDomainShow(resp)
		return nil
	},
}

// renderDomainShow renders a table with the details of a DNS domain.
func renderDomainShow(resp *responses.DomainShow) {
	if resp.Data.Domain == "" || resp.Data.Status == "" {
		fmt.Println("Domain is in pending, please check later")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "Status"})
	table.Append([]string{resp.Data.Domain, resp.Data.Status})
	table.Render()
}

// init registers the `dns domain show` command with the parent `dns domain` command
// and binds the flags for the `showOptions` struct.
func init() {
	domainCmd.AddCommand(domainShowCmd)
	_ = cli.BindFlagsFromStruct(domainShowCmd, &showOpts)
}
