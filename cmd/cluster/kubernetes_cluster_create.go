package cluster

import (
	"fmt"
	"log/slog"

	"virak-cli/internal/cli"
	"virak-cli/pkg/http"

	"github.com/spf13/cobra"
)

type createOptions struct {
	ZoneID                  string `flag:"zoneId" usage:"Zone ID to use (optional if default.zoneId is set in config)"`
	Name                    string `flag:"name" usage:"Cluster name"`
	VersionID               string `flag:"versionId" usage:"Kubernetes version ID"`
	OfferingID              string `flag:"offeringId" usage:"Service offering ID"`
	SSHKeyID                string `flag:"sshKeyId" usage:"SSH key ID"`
	NetworkID               string `flag:"networkId" usage:"Network ID"`
	HAEnabled               bool   `flag:"ha" usage:"Enable high availability"`
	ClusterSize             int    `flag:"size" usage:"Cluster size" default:"1"`
	Description             string `flag:"description" usage:"Cluster description"`
	PrivateRegistryUsername string `flag:"privateRegistryUsername" usage:"Private registry username"`
	PrivateRegistryPassword string `flag:"privateRegistryPassword" usage:"Private registry password"`
	PrivateRegistryURL      string `flag:"privateRegistryUrl" usage:"Private registry URL"`
	HAConfigControllerNodes int    `flag:"haControllerNodes" usage:"HA config controller nodes"`
	HAConfigExternalLBIP    string `flag:"haExternalLBIP" usage:"HA config external load balancer IP"`
}

var createOpts createOptions

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

func init() {
	KubernetesClusterCmd.AddCommand(kubernetesClusterCreateCmd)
	_ = cli.BindFlagsFromStruct(kubernetesClusterCreateCmd, &createOpts)
}
