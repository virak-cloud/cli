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

type showOptions struct {
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	ClusterID string `flag:"clusterId" usage:"Cluster ID"`
}

var showOpts showOptions

var kubernetesClusterShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a kubernetes cluster",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,

			cli.Required("clusterId"),
			cli.IsUlid("clusterId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &showOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		cluster, err := httpClient.GetKubernetesCluster(zoneID, showOpts.ClusterID)
		if err != nil {
			slog.Error("failed to get kubernetes cluster", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Status", "Version", "Size", "Created At"})
		table.Append([]string{cluster.Data.ID, cluster.Data.Name, cluster.Data.Status, cluster.Data.KubernetesVersion.Version, fmt.Sprintf("%d", cluster.Data.ClusterSize), fmt.Sprintf("%d", cluster.Data.CreatedAt)})
		table.Render()

		if cluster.Data.Status == "Failed" && cluster.Data.FailedReason != "" {
			fmt.Printf("\nCluster failed: %s\n", cluster.Data.FailedReason)
		}

		return nil
	},
}

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterShowCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterShowCmd, &showOpts)
}
