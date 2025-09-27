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

type listOptions struct {
	ZoneID string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
}

var listOpts listOptions

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

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterListCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterListCmd, &listOpts)
}
