package cluster

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// stopOptions holds the options for the `cluster stop` command.
// These options are populated from command-line flags.
type stopOptions struct {
	// ZoneID is the ID of the zone where the cluster is located.
	// This is optional if a default zone is set in the config.
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// ClusterID is the ID of the cluster to be stopped.
	ClusterID string `flag:"clusterId" usage:"Cluster ID"`
}

var stopOpts stopOptions

// kubernetesClusterStopCmd represents the `cluster stop` command.
// It stops a Kubernetes cluster in a specified zone.
var kubernetesClusterStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a kubernetes cluster",
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

		if err := cli.LoadFromCobraFlags(cmd, &stopOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.StopKubernetesCluster(zoneID, stopOpts.ClusterID)
		if err != nil {
			slog.Error("failed to stop kubernetes cluster", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("kubernetes cluster stopped successfully")
		fmt.Println("Success")

		return nil
	},
}

// init registers the `cluster stop` command with the parent `cluster` command
// and binds the flags for the `stopOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterStopCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterStopCmd, &stopOpts)
}
