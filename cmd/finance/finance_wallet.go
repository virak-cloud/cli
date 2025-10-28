package finance

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

type walletOptions struct{}

var walletOpt walletOptions

var financeWalletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Show wallet balances",
	Long:  `Displays the current balance and related information for your wallet.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &walletOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		resp, err := httpClient.GetWallet()
		if err != nil {
			slog.Error("failed to get wallet", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("wallet information retrieved successfully")
		fmt.Println("Wallet information retrieved successfully.")
		presenter.RenderWallet(*resp)
		return nil
	},
}

func init() {
	FinanceCmd.AddCommand(financeWalletCmd)

	_ = cli.BindFlagsFromStruct(financeWalletCmd, &walletOpt)
}