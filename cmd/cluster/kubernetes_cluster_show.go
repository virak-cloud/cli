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

// showOptions holds the options for the `cluster show` command.
// These options are populated from command-line flags.
type showOptions struct {
	// ZoneID is the ID of the zone where the cluster is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// ClusterID is the ID of the cluster to be shown.
	ClusterID string `flag:"clusterId" usage:"Cluster ID"`
}

var showOpts showOptions

// kubernetesClusterShowCmd represents the `cluster show` command.
// It shows the details of a specific Kubernetes cluster.
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

// init registers the `cluster show` command with the parent `cluster` command
// and binds the flags for the `showOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterShowCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterShowCmd, &showOpts)
}
