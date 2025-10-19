// Package http provides Kubernetes cluster management operations for the Virak Cloud API.
//
// This file contains client methods for managing Kubernetes clusters, including
// cluster lifecycle operations, scaling, version management, and service offerings.
// All operations are performed within the context of a specific zone.
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

// GetKubernetesClusters retrieves all Kubernetes clusters in the specified zone.
//
// This method returns a list of all Kubernetes clusters accessible to the user
// within the specified zone, including their current status, configuration,
// and resource allocation details.
//
// Parameters:
//   - zoneID: The unique identifier of the zone to query for clusters
//
// Returns:
//   - *responses.KubernetesClusterListResponse: List of clusters with metadata
//   - error: API error or network failure
//
// Example:
//   clusters, err := client.GetKubernetesClusters("zone-123")
//   if err != nil {
//       return fmt.Errorf("failed to list clusters: %w", err)
//   }
//   for _, cluster := range clusters.Data {
//       fmt.Printf("Cluster: %s (Status: %s)\n", cluster.Name, cluster.Status)
//   }
func (client *Client) GetKubernetesClusters(zoneID string) (*responses.KubernetesClusterListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterList, urls.BaseUrl, zoneID)
	var result responses.KubernetesClusterListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetKubernetesCluster retrieves detailed information about a specific Kubernetes cluster.
//
// This method returns comprehensive details about a single cluster, including
// its configuration, node status, network settings, and operational metrics.
//
// Parameters:
//   - zoneID: The unique identifier of the zone containing the cluster
//   - clusterID: The unique identifier of the cluster to retrieve
//
// Returns:
//   - *responses.KubernetesClusterResponse: Detailed cluster information
//   - error: API error if the cluster doesn't exist or access is denied
//
// Example:
//   cluster, err := client.GetKubernetesCluster("zone-123", "cluster-456")
//   if err != nil {
//       return fmt.Errorf("failed to get cluster: %w", err)
//   }
//   fmt.Printf("Cluster %s has %d nodes\n", cluster.Data.Name, cluster.Data.NodeCount)
func (client *Client) GetKubernetesCluster(zoneID string, clusterID string) (*responses.KubernetesClusterResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterShow, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesClusterResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateKubernetesCluster creates a new Kubernetes cluster with the specified configuration.
//
// This method provisions a new Kubernetes cluster with customizable settings including
// high availability, private registry configuration, and node sizing. The cluster
// creation is an asynchronous operation that returns immediately with a job ID.
//
// Required Parameters:
//   - zoneID: The zone where the cluster will be created
//   - name: Human-readable name for the cluster
//   - versionID: Kubernetes version identifier (from GetKubernetesVersions)
//   - offeringID: Service offering ID defining node specifications
//   - sshKey: SSH key ID for node access
//   - networkID: Network ID for cluster connectivity
//   - ha: Whether high availability should be enabled
//   - size: Number of worker nodes in the cluster
//
// Optional Parameters:
//   - description: Optional cluster description
//   - privateRegistryUsername: Username for private container registry
//   - privateRegistryPassword: Password for private container registry
//   - privateRegistryURL: URL of private container registry
//   - haControllerNodes: Number of controller nodes (if HA enabled)
//   - haExternalLBIP: External load balancer IP (if HA enabled)
//
// Returns:
//   - *responses.KubernetesMessage: Cluster creation response with job ID
//   - error: Validation error or API failure
//
// Example:
//   resp, err := client.CreateKubernetesCluster(
//       "zone-123", "my-cluster", "v1.28.0", "offering-456",
//       "ssh-key-789", "network-abc", false, 3, "Production cluster",
//       "", "", "", 0, "",
//   )
//   if err != nil {
//       return fmt.Errorf("failed to create cluster: %w", err)
//   }
//   fmt.Printf("Cluster creation job: %s\n", resp.Data.JobID)
func (client *Client) CreateKubernetesCluster(zoneID, name, versionID, offeringID, sshKey, networkID string, ha bool, size int, description string, privateRegistryUsername, privateRegistryPassword, privateRegistryURL string, haControllerNodes int, haExternalLBIP string) (*responses.KubernetesMessage, error) {
	url := fmt.Sprintf(urls.KubernetesClusterCreate, urls.BaseUrl, zoneID)
	var result responses.KubernetesMessage

	body := map[string]interface{}{
		"name":                  name,
		"kubernetes_version_id": versionID,
		"service_offering_id":   offeringID,
		"ha_enabled":            ha,
		"sshkey_id":             sshKey,
		"network_id":            networkID,
		"cluster_size":          size,
	}

	if description != "" {
		body["description"] = description
	}

	if privateRegistryUsername != "" && privateRegistryPassword != "" && privateRegistryURL != "" {
		body["private_registry"] = map[string]string{
			"username": privateRegistryUsername,
			"password": privateRegistryPassword,
			"url":      privateRegistryURL,
		}
	}

	if ha && haControllerNodes > 0 && haExternalLBIP != "" {
		body["ha_config"] = map[string]interface{}{
			"controller_nodes":         haControllerNodes,
			"external_loadbalancer_ip": haExternalLBIP,
		}
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateKubernetesClusterDetails updates the name and description of an existing cluster.
//
// Parameters:
//   - zoneID: Zone containing the cluster
//   - clusterID: Cluster to update
//   - name: New cluster name
//   - description: New cluster description
//
// Returns:
//   - *responses.KubernetesClusterResponse: Updated cluster information
//   - error: API error if update fails
func (client *Client) UpdateKubernetesClusterDetails(zoneID string, clusterID string, name string, description string) (*responses.KubernetesClusterResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterUpdate, urls.BaseUrl, zoneID, clusterID)

	var result responses.KubernetesClusterResponse

	body := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	jsonBody, _ := json.Marshal(body)

	err := client.handleRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteKubernetesCluster permanently removes a Kubernetes cluster.
//
// This operation is irreversible and will delete all nodes, storage,
// and associated resources. Ensure all important data is backed up
// before calling this method.
//
// Parameters:
//   - zoneID: Zone containing the cluster
//   - clusterID: Cluster to delete
//
// Returns:
//   - *responses.KubernetesMessage: Deletion confirmation with job ID
//   - error: API error if deletion fails
func (client *Client) DeleteKubernetesCluster(zoneID string, clusterID string) (*responses.KubernetesMessage, error) {
	url := fmt.Sprintf(urls.KubernetesClusterDelete, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesMessage

	err := client.handleRequest(http.MethodDelete, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// StartKubernetesCluster starts a stopped Kubernetes cluster.
//
// Parameters:
//   - zoneID: Zone containing the cluster
//   - clusterID: Cluster to start
//
// Returns:
//   - *responses.KubernetesClusterResponse: Updated cluster status
//   - error: API error if start fails
func (client *Client) StartKubernetesCluster(zoneID string, clusterID string) (*responses.KubernetesClusterResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterStart, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesClusterResponse

	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// StopKubernetesCluster stops a running Kubernetes cluster.
//
// Parameters:
//   - zoneID: Zone containing the cluster
//   - clusterID: Cluster to stop
//
// Returns:
//   - *responses.KubernetesClusterResponse: Updated cluster status
//   - error: API error if stop fails
func (client *Client) StopKubernetesCluster(zoneID string, clusterID string) (*responses.KubernetesClusterResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterStop, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesClusterResponse

	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ScaleKubernetesCluster scales a cluster manually or enables auto-scaling.
//
// Parameters:
//   - zoneID: Zone containing the cluster
//   - clusterID: Cluster to scale
//   - autoScaling: Enable auto-scaling if true, manual scaling if false
//   - clusterSize: Fixed cluster size (if autoScaling is false)
//   - minClusterSize: Minimum cluster size (if autoScaling is true)
//   - maxClusterSize: Maximum cluster size (if autoScaling is true)
//
// Returns:
//   - *responses.KubernetesClusterResponse: Updated cluster configuration
//   - error: API error if scaling fails
func (client *Client) ScaleKubernetesCluster(zoneID string, clusterID string, autoScaling bool, clusterSize int, minClusterSize int, maxClusterSize int) (*responses.KubernetesClusterResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterScale, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesClusterResponse

	body := map[string]interface{}{
		"auto_scaling": autoScaling,
	}
	if !autoScaling {
		body["cluster_size"] = clusterSize
	} else {
		body["min_cluster_size"] = minClusterSize
		body["max_cluster_size"] = maxClusterSize
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
// GetKubernetesVersions retrieves available Kubernetes versions for the zone.
//
// Returns a list of supported Kubernetes versions that can be used for
// cluster creation and upgrades.
//
// Parameters:
//   - zoneID: Zone to query for available versions
//
// Returns:
//   - *responses.KubernetesVersionsListResponse: List of available versions
//   - error: API error if request fails
func (client *Client) GetKubernetesVersions(zoneID string) (*responses.KubernetesVersionsListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesVersionsList, urls.BaseUrl, zoneID)
	var result responses.KubernetesVersionsListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetKubernetesServiceEvents retrieves Kubernetes service events for the zone.
//
// Returns audit logs and events for all Kubernetes operations within the zone.
//
// Parameters:
//   - zoneID: Zone to query for events
//
// Returns:
//   - *responses.KubernetesEventsListResponse: List of service events
//   - error: API error if request fails
func (client *Client) GetKubernetesServiceEvents(zoneID string) (*responses.KubernetesEventsListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesServiceEvents, urls.BaseUrl, zoneID)
	var result responses.KubernetesEventsListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetKubernetesClusterEvents retrieves events for a specific cluster.
//
// Returns detailed logs and events for operations performed on the
// specified cluster.
//
// Parameters:
//   - zoneID: Zone containing the cluster
//   - clusterID: Cluster to query for events
//
// Returns:
//   - *responses.KubernetesEventsListResponse: List of cluster events
//   - error: API error if request fails
func (client *Client) GetKubernetesClusterEvents(zoneID string, clusterID string) (*responses.KubernetesEventsListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterEvents, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesEventsListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetKubernetesServiceOfferings retrieves available Kubernetes service offerings.
//
// Returns available node configurations with pricing and specifications
// that can be used for cluster creation.
//
// Parameters:
//   - zoneID: Zone to query for service offerings
//
// Returns:
//   - *responses.KubernetesServiceOfferingsListResponse: List of available offerings
//   - error: API error if request fails
func (client *Client) GetKubernetesServiceOfferings(zoneID string) (*responses.KubernetesServiceOfferingsListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesServiceOfferingsList, urls.BaseUrl, zoneID)
	var result responses.KubernetesServiceOfferingsListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
