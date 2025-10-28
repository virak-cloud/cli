package user

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

type tokenValidateOptions struct{}

var tokenValidateOpt tokenValidateOptions

var userTokenValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate user token",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &tokenValidateOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		if err := httpClient.ValidateUserToken(); err != nil {
			slog.Error("failed to validate user token", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("user token validated successfully")
		fmt.Println("Token is valid")
		// Note: Expiration information is not provided by the validation endpoint
		// Token expiration details may be available through other means
		return nil
	},
}

func init() {
	TokenCmd.AddCommand(userTokenValidateCmd)
	_ = cli.BindFlagsFromStruct(userTokenValidateCmd, &tokenValidateOpt)
}
