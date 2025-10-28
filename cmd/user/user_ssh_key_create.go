package user

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	httpc "github.com/virak-cloud/cli/pkg/http"
)

type createOptions struct {
	Name      string `flag:"name" usage:"Name of the SSH key"`
	PublicKey string `flag:"public-key" usage:"Public SSH key"`
}

var createOpt createOptions

var sshKeyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new SSH key",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		// Declarative validation rules for the rest
		return cli.Validate(cmd,
			cli.Required("name"),
			cli.Required("public-key"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &createOpt); err != nil {
			return err
		}

		httpClient := httpc.NewClient(token)
		if _, err := httpClient.AddUserSSHKey(createOpt.Name, createOpt.PublicKey); err != nil {
			slog.Error("failed to create SSH key", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("SSH key created successfully", "name", createOpt.Name)
		fmt.Println("SSH key created successfully.")
		return nil
	},
}

func init() {
	UserSSHKeyCmd.AddCommand(sshKeyCreateCmd)

	_ = cli.BindFlagsFromStruct(sshKeyCreateCmd, &createOpt)
}