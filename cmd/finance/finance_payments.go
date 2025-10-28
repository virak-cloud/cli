package finance

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

type paymentsOptions struct{}

var paymentsOpt paymentsOptions

var financePaymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "List payment history",
	Long:  `Displays your payment history and transaction records.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Centralized login + zone handling (required for this command)
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &paymentsOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		resp, err := httpClient.ListPayments()
		if err != nil {
			slog.Error("failed to list payments", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("payment history retrieved successfully", "count", len(resp.Data))
		fmt.Println("Payment history retrieved successfully.")
		presenter.RenderPayments(resp.Data)
		return nil
	},
}

func init() {
	FinanceCmd.AddCommand(financePaymentsCmd)
	_ = cli.BindFlagsFromStruct(financePaymentsCmd, &paymentsOpt)
}