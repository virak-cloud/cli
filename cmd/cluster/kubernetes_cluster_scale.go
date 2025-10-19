package cluster

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// scaleOptions holds the options for the `cluster scale` command.
// These options are populated from command-line flags.
type scaleOptions struct {
	// ZoneID is the ID of the zone where the cluster is located.
	// This is optional if a default zone is set in the config.
	ZoneID         string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// ClusterID is the ID of the cluster to be scaled.
	ClusterID      string `flag:"clusterId" usage:"Cluster ID"`
	// AutoScaling specifies whether to enable auto-scaling for the cluster.
	AutoScaling    bool   `flag:"auto-scaling" usage:"Enable auto scaling"`
	// ClusterSize is the target number of nodes in the cluster.
	// This is required if auto-scaling is false.
	ClusterSize    int    `flag:"cluster-size" usage:"Cluster size (required if auto-scaling is false)"`
	// MinClusterSize is the minimum number of nodes in the cluster when auto-scaling is enabled.
	MinClusterSize int    `flag:"min-cluster-size" usage:"Minimum cluster size (required if auto-scaling is true)"`
	// MaxClusterSize is the maximum number of nodes in the cluster when auto-scaling is enabled.
	MaxClusterSize int    `flag:"max-cluster-size" usage:"Maximum cluster size (required if auto-scaling is true)"`
}

var scaleOpts scaleOptions

// kubernetesClusterScaleCmd represents the `cluster scale` command.
// It scales a Kubernetes cluster up or down, or enables/disables auto-scaling.
var kubernetesClusterScaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale a kubernetes cluster",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("clusterId"),
			cli.IsUlid("clusterId"),
			cli.RequiredIf("cluster-size", func(v cli.Values) bool { return v.GetBool("auto-scaling") == false }),
			cli.RequiredIf("min-cluster-size", func(v cli.Values) bool { return v.GetBool("auto-scaling") == true }),
			cli.RequiredIf("max-cluster-size", func(v cli.Values) bool { return v.GetBool("auto-scaling") == true }),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &scaleOpts); err != nil {
			return err
		}

		// Conditional validation
		if !scaleOpts.AutoScaling {
			if scaleOpts.ClusterSize <= 0 {
				return fmt.Errorf("cluster-size is required and must be greater than 0 when auto-scaling is false")
			}
		} else {
			if scaleOpts.MinClusterSize <= 0 {
				return fmt.Errorf("min-cluster-size is required and must be greater than 0 when auto-scaling is true")
			}
			if scaleOpts.MaxClusterSize <= 0 {
				return fmt.Errorf("max-cluster-size is required and must be greater than 0 when auto-scaling is true")
			}
			if scaleOpts.MinClusterSize > scaleOpts.MaxClusterSize {
				return fmt.Errorf("min-cluster-size must be less than or equal to max-cluster-size")
			}
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.ScaleKubernetesCluster(zoneID, scaleOpts.ClusterID, scaleOpts.AutoScaling, scaleOpts.ClusterSize, scaleOpts.MinClusterSize, scaleOpts.MaxClusterSize)
		if err != nil {
			slog.Error("failed to scale kubernetes cluster", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("kubernetes cluster scaled successfully")
		fmt.Println("kubernetes cluster scaled successfully")

		return nil
	},
}

// init registers the `cluster scale` command with the parent `cluster` command
// and binds the flags for the `scaleOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterScaleCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterScaleCmd, &scaleOpts)
}
