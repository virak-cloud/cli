package dns

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal"
	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

type recordDeleteOptions struct {
	Domain    string `flag:"domain" usage:"Domain name"`
	Record    string `flag:"record" usage:"Record name"`
	Type      string `flag:"type" usage:"Record type"`
	ContentID string `flag:"content-id" usage:"Content ID"`
}

var recordDeleteOpts recordDeleteOptions

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

func init() {
	recordCmd.AddCommand(recordDeleteCmd)
	_ = cli.BindFlagsFromStruct(recordDeleteCmd, &recordDeleteOpts)
}
