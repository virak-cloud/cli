package user

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

type sshKeyDeleteOptions struct {
	ID string `flag:"id" usage:"ID of the SSH key to delete"`
}

var sshKeyDeleteOpt sshKeyDeleteOptions

var userSSHKeyDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an SSH key",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("id"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &sshKeyDeleteOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		if _, err := httpClient.DeleteUserSSHKey(sshKeyDeleteOpt.ID); err != nil {
			slog.Error("failed to delete SSH key", "error", err, "id", sshKeyDeleteOpt.ID)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("SSH key deleted successfully", "id", sshKeyDeleteOpt.ID)
		fmt.Println("SSH key deleted successfully.")
		return nil
	},
}

func init() {
	UserSSHKeyCmd.AddCommand(userSSHKeyDeleteCmd)
	_ = cli.BindFlagsFromStruct(userSSHKeyDeleteCmd, &sshKeyDeleteOpt)
}