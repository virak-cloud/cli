package dns

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

// createOptions holds the options for the `dns domain create` command.
// These options are populated from command-line flags.
type createOptions struct {
	// Domain is the name of the domain to be created.
	Domain string `flag:"domain" usage:"Domain name to create"`
}

var createOpts createOptions

// domainCreateCmd represents the `dns domain create` command.
// It creates a new DNS domain.
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

// init registers the `dns domain create` command with the parent `dns domain` command
// and binds the flags for the `createOptions` struct.
func init() {
	domainCmd.AddCommand(domainCreateCmd)
	_ = cli.BindFlagsFromStruct(domainCreateCmd, &createOpts)
}
