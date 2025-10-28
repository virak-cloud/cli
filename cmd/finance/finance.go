package finance

import (
	"github.com/spf13/cobra"
)

// FinanceCmd is the root command for managing financial operations.
var FinanceCmd = &cobra.Command{
	Use:   "finance",
	Short: "Manage financial operations",
	Long:  `The finance command allows you to manage your financial operations including wallet balances, cost documents, payments, and expenses.`,
}

func init() {
	// Add all subcommands
	FinanceCmd.AddCommand(financeWalletCmd)
	FinanceCmd.AddCommand(financeDocumentsCmd)
	FinanceCmd.AddCommand(financePaymentsCmd)
	FinanceCmd.AddCommand(financeExpensesCmd)
}