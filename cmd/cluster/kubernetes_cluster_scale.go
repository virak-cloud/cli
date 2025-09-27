package cluster

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type scaleOptions struct {
	ZoneID         string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	ClusterID      string `flag:"clusterId" usage:"Cluster ID"`
	AutoScaling    bool   `flag:"auto-scaling" usage:"Enable auto scaling"`
	ClusterSize    int    `flag:"cluster-size" usage:"Cluster size (required if auto-scaling is false)"`
	MinClusterSize int    `flag:"min-cluster-size" usage:"Minimum cluster size (required if auto-scaling is true)"`
	MaxClusterSize int    `flag:"max-cluster-size" usage:"Maximum cluster size (required if auto-scaling is true)"`
}

var scaleOpts scaleOptions

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

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterScaleCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterScaleCmd, &scaleOpts)
}
