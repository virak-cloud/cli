package dns

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

type recordUpdateOptions struct {
	Domain    string `flag:"domain" usage:"Domain name"`
	Record    string `flag:"record" usage:"DNS record name (e.g., www, mail)"`
	Type      string `flag:"type" usage:"DNS record type (A, AAAA, CNAME, MX, TXT, NS, SOA, SRV, CAA, TLSA)"`
	ContentID string `flag:"contentId" usage:"Content ID of the record to update"`
	Content   string `flag:"content" usage:"DNS record content/value (IP address, hostname, text, etc.)"`
	TTL       int    `flag:"ttl" usage:"Time To Live in seconds (default 3600)"`
	Priority  int    `flag:"priority" usage:"MX/SRV record priority (lower values have higher priority)"`
	Weight    int    `flag:"weight" usage:"SRV record weight for load balancing"`
	Port      int    `flag:"port" usage:"SRV record port number"`
	Flags     int    `flag:"flags" usage:"CAA record flags (0 or 128 for critical)"`
	Tag       string `flag:"tag" usage:"CAA record tag (issue, issuewild, iodef, etc.)"`
	License   int    `flag:"license" usage:"TLSA record certificate usage (0=PKIX-TA, 1=PKIX-EE, 2=DANE-TA, 3=DANE-EE)"`
	Choicer   int    `flag:"choicer" usage:"TLSA record selector (0=full cert, 1=subject public key)"`
	Match     int    `flag:"match" usage:"TLSA record matching type (0=exact, 1=SHA-256, 2=SHA-512)"`
}

var recordUpdateOpts recordUpdateOptions

var recordUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a record for a domain",
	Long:  "Update a record for a domain",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("domain"),
			cli.Required("record"),
			cli.Required("type"),
			cli.Required("contentId"),
			cli.IsUlid("contentId"),
			cli.Required("content"),
			// cli.Required("ttl"),
			cli.RequiredIf("priority", func(v cli.Values) bool { return v.GetString("type") == "MX" || v.GetString("type") == "SRV" }),
			cli.RequiredIf("weight", func(v cli.Values) bool { return v.GetString("type") == "SRV" }),
			cli.RequiredIf("port", func(v cli.Values) bool { return v.GetString("type") == "SRV" }),
			cli.RequiredIf("flags", func(v cli.Values) bool { return v.GetString("type") == "CAA" }),
			cli.RequiredIf("tag", func(v cli.Values) bool { return v.GetString("type") == "CAA" }),
			cli.RequiredIf("license", func(v cli.Values) bool { return v.GetString("type") == "TLSA" }),
			cli.RequiredIf("choicer", func(v cli.Values) bool { return v.GetString("type") == "TLSA" }),
			cli.RequiredIf("match", func(v cli.Values) bool { return v.GetString("type") == "TLSA" }),
			cli.OneOf("type", "A", "AAAA", "CNAME", "MX", "TXT", "NS", "SOA", "SRV", "CAA", "TLSA"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		token := cli.TokenFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &recordUpdateOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.UpdateRecord(recordUpdateOpts.Domain, recordUpdateOpts.Record, recordUpdateOpts.Type, recordUpdateOpts.ContentID, recordUpdateOpts.Content, recordUpdateOpts.TTL, recordUpdateOpts.Priority, recordUpdateOpts.Weight, recordUpdateOpts.Port, recordUpdateOpts.Flags, recordUpdateOpts.Tag, recordUpdateOpts.License, recordUpdateOpts.Choicer, recordUpdateOpts.Match)
		if err != nil {
			slog.Error("failed to update record", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("record updated successfully")
		fmt.Println("Record updated successfully")
		return nil
	},
}

func init() {
	recordCmd.AddCommand(recordUpdateCmd)
	_ = cli.BindFlagsFromStruct(recordUpdateCmd, &recordUpdateOpts)
}
