package dns

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"
)

type recordCreateOptions struct {
	Domain   string `flag:"domain" usage:"Domain name"`
	Record   string `flag:"record" usage:"DNS record name (e.g., www, mail)"`
	Type     string `flag:"type" usage:"DNS record type (A, AAAA, CNAME, MX, TXT, NS, SOA, SRV, CAA, TLSA)"`
	Content  string `flag:"content" usage:"DNS record content/value (IP address, hostname, text, etc.)"`
	TTL      int    `flag:"ttl" usage:"Time To Live in seconds (default 3600)"`
	Priority int    `flag:"priority" usage:"MX/SRV record priority (lower values have higher priority)"`
	Weight   int    `flag:"weight" usage:"SRV record weight for load balancing"`
	Port     int    `flag:"port" usage:"SRV record port number"`
	Flags    int    `flag:"flags" usage:"CAA record flags (0 or 128 for critical)"`
	Tag      string `flag:"tag" usage:"CAA record tag (issue, issuewild, iodef, etc.)"`
	License  int    `flag:"license" usage:"TLSA record certificate usage (0=PKIX-TA, 1=PKIX-EE, 2=DANE-TA, 3=DANE-EE)"`
	Choicer  int    `flag:"choicer" usage:"TLSA record selector (0=full cert, 1=subject public key)"`
	Match    int    `flag:"match" usage:"TLSA record matching type (0=exact, 1=SHA-256, 2=SHA-512)"`
}

var recordCreateOpts recordCreateOptions

var recordCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new record for a domain",
	Long:  "Create a new record for a domain",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(false)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("domain"),
			cli.Required("record"),
			cli.Required("type"),
			cli.Required("content"),
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

		if err := cli.LoadFromCobraFlags(cmd, &recordCreateOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.CreateRecord(recordCreateOpts.Domain, recordCreateOpts.Record, recordCreateOpts.Type, recordCreateOpts.Content, recordCreateOpts.TTL, recordCreateOpts.Priority, recordCreateOpts.Weight, recordCreateOpts.Port, recordCreateOpts.Flags, recordCreateOpts.Tag, recordCreateOpts.License, recordCreateOpts.Choicer, recordCreateOpts.Match)
		if err != nil {
			slog.Error("failed to create record", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("record created successfully")
		fmt.Println("Success")
		return nil
	},
}

func init() {
	recordCmd.AddCommand(recordCreateCmd)
	_ = cli.BindFlagsFromStruct(recordCreateCmd, &recordCreateOpts)
}
