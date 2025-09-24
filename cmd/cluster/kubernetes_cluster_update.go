package cluster

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type updateOptions struct {
	DefaultZone bool   `flag:"default-zone" usage:"Use default.zoneId from config"`
	ZoneID      string `flag:"zoneId" usage:"Zone ID"`
	ClusterID   string `flag:"clusterId" usage:"Cluster ID"`
	Name        string `flag:"name" usage:"New cluster name"`
	Description string `flag:"description" usage:"New cluster description"`
}

var updateOpts updateOptions

var kubernetesClusterUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a kubernetes cluster",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("clusterId"),
			cli.IsUlid("clusterId"),
			cli.Required("name"),
			cli.MinLength("name", 3),
			cli.MaxLength("name", 50),
			cli.Required("description"),
			cli.MinLength("description", 3),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &updateOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		cluster, err := httpClient.UpdateKubernetesClusterDetails(zoneID, updateOpts.ClusterID, updateOpts.Name, updateOpts.Description)
		if err != nil {
			slog.Error("failed to update kubernetes cluster", "error", err)
			fmt.Println("Error: failed to update kubernetes cluster.")
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("kubernetes cluster updated successfully")
		fmt.Printf("Kubernetes cluster updated successfully!\nID: %s\nName: %s\nDescription: %s\n", cluster.Data.ID, cluster.Data.Name, cluster.Data.Description)

		return nil
	},
}

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterUpdateCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterUpdateCmd, &updateOpts)
}
