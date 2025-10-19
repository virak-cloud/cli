package cluster

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// versionsListOptions holds the options for the `cluster versions-list` command.
// These options are populated from command-line flags.
type versionsListOptions struct {
	// ZoneID is the ID of the zone to list Kubernetes versions from.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
}

var versionsListOpts versionsListOptions

// kubernetesVersionsListCmd represents the `cluster versions-list` command.
// It lists all available Kubernetes versions in a specified zone.
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

// init registers the `cluster versions-list` command with the parent `cluster` command
// and binds the flags for the `versionsListOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesVersionsListCmd)
	_ = cli.BindFlagsFromStruct(kubernetesVersionsListCmd, &versionsListOpts)
}
