package cluster

import (
	"fmt"
	"log/slog"
	"os"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type versionsListOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
}

var versionsListOpts versionsListOptions

var kubernetesVersionsListCmd = &cobra.Command{
	Use:   "versions-list",
	Short: "List available kubernetes versions",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		httpClient := http.NewClient(token)
		versions, err := httpClient.GetKubernetesVersions(zoneID)
		if err != nil {
			slog.Error("failed to get kubernetes versions", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Version", "Enabled"})

		for _, v := range versions.Data {
			table.Append([]string{v.ID, v.Version, fmt.Sprintf("%t", v.Enabled)})
		}
		table.Render()

		return nil
	},
}

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesVersionsListCmd)
	_ = cli.BindFlagsFromStruct(kubernetesVersionsListCmd, &versionsListOpts)
}
