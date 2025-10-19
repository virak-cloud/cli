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

// listOptions holds the options for the `cluster list` command.
// These options are populated from command-line flags.
type listOptions struct {
	// ZoneID is the ID of the zone to list clusters from.
	// This is optional if a default zone is set in the config.
	ZoneID string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
}

var listOpts listOptions

// kubernetesClusterListCmd represents the `cluster list` command.
// It lists all Kubernetes clusters in a specified zone.
var kubernetesClusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all kubernetes clusters",
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
		clusters, err := httpClient.GetKubernetesClusters(zoneID)
		if err != nil {
			slog.Error("failed to get kubernetes clusters", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Status", "Version", "Worker Size"})

		for _, cluster := range clusters.Data {
			table.Append([]string{cluster.ID, cluster.Name, cluster.Status, cluster.KubernetesVersion.Version, fmt.Sprintf("%d", cluster.ClusterSize)})
		}
		table.Render()

		return nil
	},
}

// init registers the `cluster list` command with the parent `cluster` command
// and binds the flags for the `listOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterListCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterListCmd, &listOpts)
}
