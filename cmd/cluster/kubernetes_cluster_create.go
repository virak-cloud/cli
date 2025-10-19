package cluster

import (
	"fmt"
	"log/slog"

	"github.com/virak-cloud/cli/internal/cli"
	"github.com/virak-cloud/cli/pkg/http"

	"github.com/spf13/cobra"
)

// createOptions holds the options for the `cluster create` command.
// These options are populated from command-line flags.
type createOptions struct {
	// ZoneID is the ID of the zone where the cluster will be created.
	// This is optional if a default zone is set in the config.
	ZoneID                  string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	// Name is the name of the Kubernetes cluster.
	Name                    string `flag:"name" usage:"Cluster name"`
	// VersionID is the ID of the Kubernetes version to use.
	VersionID               string `flag:"versionId" usage:"Kubernetes version ID"`
	// OfferingID is the ID of the service offering for the cluster nodes.
	OfferingID              string `flag:"offeringId" usage:"Service offering ID"`
	// SSHKeyID is the ID of the SSH key to be installed on the nodes.
	SSHKeyID                string `flag:"sshKeyId" usage:"SSH key ID"`
	// NetworkID is the ID of the network to which the cluster will be connected.
	NetworkID               string `flag:"networkId" usage:"Network ID"`
	// HAEnabled specifies whether to enable high availability for the cluster.
	HAEnabled               bool   `flag:"ha" usage:"Enable high availability"`
	// ClusterSize is the number of nodes in the cluster.
	ClusterSize             int    `flag:"size" usage:"Cluster size" default:"1"`
	// Description is a description of the cluster.
	Description             string `flag:"description" usage:"Cluster description"`
	// PrivateRegistryUsername is the username for a private container registry.
	PrivateRegistryUsername string `flag:"privateRegistryUsername" usage:"Private registry username"`
	// PrivateRegistryPassword is the password for a private container registry.
	PrivateRegistryPassword string `flag:"privateRegistryPassword" usage:"Private registry password"`
	// PrivateRegistryURL is the URL of a private container registry.
	PrivateRegistryURL      string `flag:"privateRegistryUrl" usage:"Private registry URL"`
	// HAConfigControllerNodes is the number of controller nodes for a high-availability cluster.
	HAConfigControllerNodes int    `flag:"haControllerNodes" usage:"HA config controller nodes"`
	// HAConfigExternalLBIP is the external load balancer IP for a high-availability cluster.
	HAConfigExternalLBIP    string `flag:"haExternalLBIP" usage:"HA config external load balancer IP"`
}

var createOpts createOptions

// kubernetesClusterCreateCmd represents the `cluster create` command.
// It creates a new Kubernetes cluster in a specified zone.
var kubernetesClusterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new kubernetes cluster",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.Preflight(true)(cmd, args); err != nil {
			return err
		}
		return cli.Validate(cmd,
			cli.Required("name"),
			cli.Required("versionId"),
			cli.IsUlid("versionId"),
			cli.Required("offeringId"),
			cli.IsUlid("offeringId"),
			cli.Required("sshKeyId"),
			cli.IsUlid("sshKeyId"),
			cli.Required("networkId"),
			cli.IsUlid("networkId"),
		)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		token := cli.TokenFromContext(cmd.Context())
		zoneID := cli.ZoneIDFromContext(cmd.Context())

		if err := cli.LoadFromCobraFlags(cmd, &createOpts); err != nil {
			return err
		}

		httpClient := http.NewClient(token)
		_, err := httpClient.CreateKubernetesCluster(zoneID, createOpts.Name, createOpts.VersionID, createOpts.OfferingID, createOpts.SSHKeyID, createOpts.NetworkID, createOpts.HAEnabled, createOpts.ClusterSize, createOpts.Description, createOpts.PrivateRegistryUsername, createOpts.PrivateRegistryPassword, createOpts.PrivateRegistryURL, createOpts.HAConfigControllerNodes, createOpts.HAConfigExternalLBIP)
		if err != nil {
			slog.Error("failed to create kubernetes cluster", "error", err)
			return fmt.Errorf("error: %w", err)
		}

		slog.Info("kubernetes cluster created successfully")
		fmt.Println("kubernetes cluster created successfully. Check status with 'virak-cli cluster list'")

		return nil
	},
}

// init registers the `cluster create` command with the parent `cluster` command
// and binds the flags for the `createOptions` struct.
func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterCreateCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterCreateCmd, &createOpts)
}
