package finance

import (
	"fmt"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

var validProductTypes = []string{
	"Instance",
	"InstanceNetworkSecondaryIpAddressV4",
	"BucketSize",
	"NetworkInternetPublicAddressV4",
	"KubernetesNode",
	"NetworkDevice",
	"NetworkTraffic",
	"SupportOfferings",
	"BucketUploadTraffic",
	"BucketDownloadTraffic",
}

type expensesOptions struct {
	startDate   string
	endDate     string
	expenseType string
	productType string
	productID   string
	page        uint
}

var expensesOpt expensesOptions

var financeExpensesCmd = &cobra.Command{
	Use:   "expenses",
	Short: "List expenses with filtering",
	Long:  `Displays your expenses with optional filtering by date range and type.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cli.Preflight(true)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		// Validate required parameters
		if expensesOpt.productType == "" {
			return fmt.Errorf("product-type is required")
		}
		if expensesOpt.productID == "" {
			return fmt.Errorf("product-id is required")
		}

		// Validate product type
		valid := false
		for _, validType := range validProductTypes {
			if expensesOpt.productType == validType {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid product-type '%s', must be one of: %v", expensesOpt.productType, validProductTypes)
		}

		// Build filters map
		filters := make(map[string]string)
		if expensesOpt.startDate != "" {
			filters["start_date"] = expensesOpt.startDate
		}
		if expensesOpt.endDate != "" {
			filters["end_date"] = expensesOpt.endDate
		}
		if expensesOpt.expenseType != "" {
			filters["type"] = expensesOpt.expenseType
		}
		if expensesOpt.page > 0 {
			filters["page"] = fmt.Sprintf("%d", expensesOpt.page)
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListExpensesWithRequiredParams(expensesOpt.productType, expensesOpt.productID, filters)
		if err != nil {
			return fmt.Errorf("could not fetch expenses: %w", err)
		}

		presenter.RenderExpenses(resp.Data)
		return nil
	},
}

func init() {
	FinanceCmd.AddCommand(financeExpensesCmd)
	financeExpensesCmd.Flags().StringVar(&expensesOpt.startDate, "start-date", "", "Start date for filtering expenses (YYYY-MM-DD format)")
	financeExpensesCmd.Flags().StringVar(&expensesOpt.endDate, "end-date", "", "End date for filtering expenses (YYYY-MM-DD format)")
	financeExpensesCmd.Flags().StringVar(&expensesOpt.expenseType, "type", "", "Type of expenses to filter by")
	financeExpensesCmd.Flags().StringVar(&expensesOpt.productType, "product-type", "", "Product type (required)")
	financeExpensesCmd.Flags().StringVar(&expensesOpt.productID, "product-id", "", "Product ID (required)")
	financeExpensesCmd.Flags().UintVar(&expensesOpt.page, "page", 1, "Page number for pagination (default 1)")
	_ = cli.BindFlagsFromStruct(financeExpensesCmd, &expensesOpt)
}
