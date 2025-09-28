package dns

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

type deleteOptions struct {
	Domain string `flag:"domain" usage:"Domain name to delete"`
}

var deleteOpts deleteOptions

var domainDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a domain",
	Long:  "Delete a domain",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd, cli.Required("domain"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &deleteOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.DeleteDomain(deleteOpts.Domain)
		if err != nil {
			slog.Error("failed to delete domain", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("domain deleted successfully")
		fmt.Println("Domain Delete request submitted successfully.")
		return nil
	},
}

func init() {
	domainCmd.AddCommand(domainDeleteCmd)
	_ = cli.BindFlagsFromStruct(domainDeleteCmd, &deleteOpts)
}
