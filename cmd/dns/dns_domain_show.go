package dns

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"
	"virak-cli/pkg/http/responses"
)

type showOptions struct {
	Domain string `flag:"domain" usage:"Domain name to show"`
}

var showOpts showOptions

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

func init() {
	domainCmd.AddCommand(domainShowCmd)
	_ = cli.BindFlagsFromStruct(domainShowCmd, &showOpts)
}
