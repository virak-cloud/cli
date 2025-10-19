package cluster

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// updateOptions holds the options for the `cluster update` command.
// These options are populated from command-line flags.
type updateOptions struct {
	// ZoneID is the ID of the zone where the cluster is located.
	// This is optional if a default zone is set in the config.
	ZoneID      string `flag:"zoneId" usage:"Zone ID (optional if default.zoneId is set in config)"`
	// ClusterID is the ID of the cluster to be updated.
	ClusterID   string `flag:"clusterId" usage:"Cluster ID"`
	// Name is the new name for the cluster.
	Name        string `flag:"name" usage:"New cluster name"`
	// Description is the new description for the cluster.
	Description string `flag:"description" usage:"New cluster description"`
}

var updateOpts updateOptions

// kubernetesClusterUpdateCmd represents the `cluster update` command.
// It updates the name and description of a Kubernetes cluster.
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

// init registers the `cluster update` command with the parent `cluster` command
// and binds the flags for the `updateOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterUpdateCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterUpdateCmd, &updateOpts)
}
