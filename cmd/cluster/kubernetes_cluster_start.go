package cluster

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type startOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID"`
	ClusterID   string `flag:"clusterId" usage:"Cluster ID"`
}

var startOpts startOptions

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

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterStartCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterStartCmd, &startOpts)
}
