package user

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

type tokenAbilitiesOptions struct{}

var tokenAbilitiesOpt tokenAbilitiesOptions

var userTokenAbilitiesCmd = &cobra.Command{
	Use:   "abilities",
	Short: "List token permissions and scopes",
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

		if err := cli.LoadFromCobraFlags(cmd, &tokenAbilitiesOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		resp, err := httpClient.GetUserTokenAbilities()
		if err != nil {
			slog.Error("failed to fetch token abilities", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("token abilities retrieved successfully")
		presenter.RenderTokenAbilities(resp.Abilities)
		return nil
	},
}

func init() {
	TokenCmd.AddCommand(userTokenAbilitiesCmd)
	_ = cli.BindFlagsFromStruct(userTokenAbilitiesCmd, &tokenAbilitiesOpt)
}
