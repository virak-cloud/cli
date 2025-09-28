package dns

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

type createOptions struct {
	Domain string `flag:"domain" usage:"Domain name to create"`
}

var createOpts createOptions

var domainCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "Create a new domain",
	Long:    "Create a new domain",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd, cli.Required("domain"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &createOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.CreateDomain(createOpts.Domain)
		if err != nil {
			slog.Error("failed to create domain", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("domain created successfully")
		fmt.Println("Domain creation initiated successfully, consider to setup nameservers if not done already.")
		return nil
	},
}

func init() {
	domainCmd.AddCommand(domainCreateCmd)
	_ = cli.BindFlagsFromStruct(domainCreateCmd, &createOpts)
}
