package cluster

import (
	"github.com/spf13/cobra"
)

var KubernetesClusterCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"clusters", "k8s", "kubernetes"},
	Short:   "Manage your kubernetes clusters",
	Long:    `Manage your kubernetes clusters`,
}

func init() {

}
