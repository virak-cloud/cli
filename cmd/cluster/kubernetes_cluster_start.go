package cluster

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// startOptions holds the options for the `cluster start` command.
// These options are populated from command-line flags.
type startOptions struct {
	// ZoneID is the ID of the zone where the cluster is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// ClusterID is the ID of the cluster to be started.
	ClusterID string `flag:"clusterId" usage:"Cluster ID"`
}

var startOpts startOptions

// kubernetesClusterStartCmd represents the `cluster start` command.
// It starts a Kubernetes cluster in a specified zone.
var kubernetesClusterStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a kubernetes cluster",
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

		if err := cli.LoadFromCobraFlags(cmd, &startOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.StartKubernetesCluster(zoneID, startOpts.ClusterID)
		if err != nil {
			slog.Error("failed to start kubernetes cluster", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("kubernetes cluster started successfully")
		fmt.Println("Success")

		return nil
	},
}

// init registers the `cluster start` command with the parent `cluster` command
// and binds the flags for the `startOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterStartCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterStartCmd, &startOpts)
}
