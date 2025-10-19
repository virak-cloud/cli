package cluster

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// deleteOptions holds the options for the `cluster delete` command.
// These options are populated from command-line flags.
type deleteOptions struct {
	// ZoneID is the ID of the zone where the cluster is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// ClusterID is the ID of the cluster to be deleted.
	ClusterID string `flag:"clusterId" usage:"Cluster ID"`
}

var deleteOpts deleteOptions

// kubernetesClusterDeleteCmd represents the `cluster delete` command.
// It deletes a Kubernetes cluster in a specified zone.
var kubernetesClusterDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a kubernetes cluster",
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

		if err := cli.LoadFromCobraFlags(cmd, &deleteOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.DeleteKubernetesCluster(zoneID, deleteOpts.ClusterID)
		if err != nil {
			slog.Error("failed to delete kubernetes cluster", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("kubernetes cluster deleted successfully")
		fmt.Println("Kubernetes cluster deleted successfully")

		return nil
	},
}

// init registers the `cluster delete` command with the parent `cluster` command
// and binds the flags for the `deleteOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterDeleteCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterDeleteCmd, &deleteOpts)
}
