package cluster

import (
	"github.com/spf13/cobra"
)

// KubernetesClusterCmd is the parent command for all Kubernetes cluster related commands.
// It doesn't have a run function of its own, as it relies on subcommands for its functionality.
var KubernetesClusterCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"clusters", "k8s", "kubernetes"},
	Short:   "Manage your kubernetes clusters",
	Long:    `Manage your kubernetes clusters`,
}

func init() {

}
