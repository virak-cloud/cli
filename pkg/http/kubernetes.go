package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	urls "virak-cli/pkg"
	"virak-cli/pkg/http/responses"
)

func (client *Client) GetKubernetesClusters(zoneID string) (*responses.KubernetesClusterListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterList, urls.BaseUrl, zoneID)
	var result responses.KubernetesClusterListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetKubernetesCluster(zoneID string, clusterID string) (*responses.KubernetesClusterResponse, error) {

	url := fmt.Sprintf(urls.KubernetesClusterShow, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesClusterResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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

func (client *Client) DeleteKubernetesCluster(zoneID string, clusterID string) (*responses.KubernetesMessage, error) {
	url := fmt.Sprintf(urls.KubernetesClusterDelete, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesMessage

	err := client.handleRequest(http.MethodDelete, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) StartKubernetesCluster(zoneID string, clusterID string) (*responses.KubernetesClusterResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterStart, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesClusterResponse

	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) StopKubernetesCluster(zoneID string, clusterID string) (*responses.KubernetesClusterResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterStop, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesClusterResponse

	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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
func (client *Client) GetKubernetesVersions(zoneID string) (*responses.KubernetesVersionsListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesVersionsList, urls.BaseUrl, zoneID)
	var result responses.KubernetesVersionsListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetKubernetesServiceEvents(zoneID string) (*responses.KubernetesEventsListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesServiceEvents, urls.BaseUrl, zoneID)
	var result responses.KubernetesEventsListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (client *Client) GetKubernetesClusterEvents(zoneID string, clusterID string) (*responses.KubernetesEventsListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesClusterEvents, urls.BaseUrl, zoneID, clusterID)
	var result responses.KubernetesEventsListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (client *Client) GetKubernetesServiceOfferings(zoneID string) (*responses.KubernetesServiceOfferingsListResponse, error) {
	url := fmt.Sprintf(urls.KubernetesServiceOfferingsList, urls.BaseUrl, zoneID)
	var result responses.KubernetesServiceOfferingsListResponse

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
