package dns

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

type recordListOptions struct {
	Domain string `flag:"domain" usage:"Domain name"`
}

var recordListOpts recordListOptions

var recordListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all records for a domain",
	Long:  "List all records for a domain",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd, cli.Required("domain"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &recordListOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.GetRecords(recordListOpts.Domain)
		if err != nil {
			slog.Error("failed to get records", "error", err)
			return fmt.Errorf("error: %w", err)
		}
		renderRecordList(resp)
		return nil
	},
}

func renderRecordList(resp *responses.RecordList) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "TTL", "Status", "Protected", "Content"})
	for _, record := range resp.Data {
		var contents []string
		for _, content := range record.Content {
			contents = append(contents, content.ContentRaw)
		}
		contentStr := strings.Join(contents, ", ")
		protectedStr := "No"
		if record.IsProtected {
			protectedStr = "Yes"
		}
		table.Append([]string{record.Name, record.Type, fmt.Sprintf("%d", record.TTL), record.Status, protectedStr, contentStr})
	}
	table.Render()
}

func init() {
	recordCmd.AddCommand(recordListCmd)
	_ = cli.BindFlagsFromStruct(recordListCmd, &recordListOpts)
}
