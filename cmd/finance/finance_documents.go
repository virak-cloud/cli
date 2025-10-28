package finance

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

type documentsOptions struct {
	Year int `flag:"year" usage:"Year to fetch cost documents for"`
}

var documentsOpt documentsOptions

var financeDocumentsCmd = &cobra.Command{
	Use:   "documents",
	Short: "List cost documents by year",
	Long:  `Displays cost documents for a specific year. Requires the --year flag.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd, cli.Required("year"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &documentsOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		resp, err := httpClient.ListDocumentsGET(documentsOpt.Year)
		if err != nil {
			slog.Error("failed to get cost documents", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("cost documents retrieved successfully", "year", documentsOpt.Year)
		fmt.Println("Cost documents retrieved successfully.")
		presenter.RenderCostDocuments(resp.Data)
		return nil
	},
}

func init() {
	FinanceCmd.AddCommand(financeDocumentsCmd)

	_ = cli.BindFlagsFromStruct(financeDocumentsCmd, &documentsOpt)
}