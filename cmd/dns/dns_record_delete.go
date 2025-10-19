package dns

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

// recordDeleteOptions holds the options for the `dns record delete` command.
// These options are populated from command-line flags.
type recordDeleteOptions struct {
	// Domain is the name of the domain to delete the record from.
	Domain    string `flag:"domain" usage:"Domain name"`
	// Record is the name of the DNS record to delete.
	Record    string `flag:"record" usage:"Record name"`
	// Type is the type of the DNS record to delete.
	Type      string `flag:"type" usage:"Record type"`
	// ContentID is the ID of the content to delete.
	ContentID string `flag:"content-id" usage:"Content ID"`
}

var recordDeleteOpts recordDeleteOptions

// recordDeleteCmd represents the `dns record delete` command.
// It deletes a DNS record from a domain.
var recordDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a record for a domain",
	Long:  "Delete a record for a domain",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("domain"),
			cli.Required("record"),
			cli.Required("type"),
			cli.Required("content-id"),
			cli.OneOf("type", "A", "AAAA", "CNAME", "MX", "TXT", "NS", "SOA", "SRV", "CAA", "TLSA"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !internal.IsValidULID(recordDeleteOpts.ContentID) {
			return fmt.Errorf("error: content-id must be a valid ULID")
		}

		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &recordDeleteOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.DeleteRecord(recordDeleteOpts.Domain, recordDeleteOpts.Record, recordDeleteOpts.Type, recordDeleteOpts.ContentID)
		if err != nil {
			slog.Error("failed to delete record", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("record deleted successfully")
		fmt.Println("Success")
		return nil
	},
}

// init registers the `dns record delete` command with the parent `dns record` command
// and binds the flags for the `recordDeleteOptions` struct.
func init() {
	recordCmd.AddCommand(recordDeleteCmd)
	_ = cli.BindFlagsFromStruct(recordDeleteCmd, &recordDeleteOpts)
}
