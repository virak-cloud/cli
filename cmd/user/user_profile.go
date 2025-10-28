package user

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

type profileOptions struct{}

var profileOpt profileOptions

var userProfileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show user profile information",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &profileOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		resp, err := httpClient.GetUserProfile()
		if err != nil {
			slog.Error("failed to get user profile", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("user profile retrieved successfully")
		fmt.Println("User profile retrieved successfully.")

		// Render profile in table format
		presenter.RenderUserProfile(resp)

		return nil
	},
}

func init() {
	UserCmd.AddCommand(userProfileCmd)
	_ = cli.BindFlagsFromStruct(userProfileCmd, &profileOpt)
}
