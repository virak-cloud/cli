package cluster

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

type stopOptions struct {
	ZoneID    string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	ClusterID string `flag:"clusterId" usage:"Cluster ID"`
}

var stopOpts stopOptions

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

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterStopCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterStopCmd, &stopOpts)
}
