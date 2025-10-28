package user

import (
	"fmt"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/internal/presenter"
	"log/slog"

	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

type sshKeyListOptions struct{}

var sshKeyListOpt sshKeyListOptions

var userSSHKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List SSH keys",
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &sshKeyListOpt); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		resp, err := httpClient.ListUserSSHKeys()
		if err != nil {
			slog.Error("failed to list SSH keys", "error", err)
			return fmt.Errorf("could not list SSH keys: %w", err)
		}

		slog.Info("successfully retrieved SSH keys", "count", len(resp.UserData))
		presenter.RenderSSHKeyList(resp.UserData)
		fmt.Println("SSH keys listed successfully.")

		return nil
	},
}

func init() {
	UserSSHKeyCmd.AddCommand(userSSHKeyListCmd)
	_ = cli.BindFlagsFromStruct(userSSHKeyListCmd, &sshKeyListOpt)
}