package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"net/http"
)

func (client *Client) ListInstances(zoneId string) (*responses.InstanceListResponse, error) {
	var result responses.InstanceListResponse
	url := fmt.Sprintf(urls.InstanceList, urls.BaseUrl, zoneId)

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) ListInstanceServiceOfferings(zoneId string) (*responses.InstanceServiceOfferingListResponse, error) {
	var result responses.InstanceServiceOfferingListResponse
	url := fmt.Sprintf(urls.InstanceServiceOfferingList, urls.BaseUrl, zoneId)

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) ListInstanceVMImages(zoneId string) (*responses.InstanceVMImageListResponse, error) {
	var result responses.InstanceVMImageListResponse
	url := fmt.Sprintf(urls.InstanceVMImageList, urls.BaseUrl, zoneId)

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) CreateInstance(zoneId, serviceOfferingId, vmImageId string, networkIds []string, name string) (*responses.InstanceCreateResponse, error) {
	var result responses.InstanceCreateResponse
	url := fmt.Sprintf(urls.InstanceCreate, urls.BaseUrl, zoneId)
	body, err := json.Marshal(map[string]interface{}{
		"service_offering_id": serviceOfferingId,
		"vm_image_id":         vmImageId,
		"network_ids":         networkIds,
		"name":                name,
	})
	if err != nil {
		return nil, err
	}
	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) RebuildInstance(zoneId, instanceId, vmImageId string) (*responses.InstanceCreateResponse, error) {
	var result responses.InstanceCreateResponse
	url := fmt.Sprintf(urls.InstanceRebuild, urls.BaseUrl, zoneId, instanceId)
	body, err := json.Marshal(map[string]string{
		"vm_image_id": vmImageId,
	})
	if err != nil {
		return nil, err
	}
	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) StartInstance(zoneId, instanceId string) (*responses.InstanceCreateResponse, error) {
	var result responses.InstanceCreateResponse
	url := fmt.Sprintf(urls.InstanceStart, urls.BaseUrl, zoneId, instanceId)
	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) StopInstance(zoneId, instanceId string, forced bool) (*responses.InstanceCreateResponse, error) {
	var result responses.InstanceCreateResponse
	url := fmt.Sprintf(urls.InstanceStop, urls.BaseUrl, zoneId, instanceId)
	body, err := json.Marshal(map[string]bool{
		"forced": forced,
	})
	if err != nil {
		return nil, err
	}
	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) RebootInstance(zoneId, instanceId string) (*responses.InstanceCreateResponse, error) {
	var result responses.InstanceCreateResponse
	url := fmt.Sprintf(urls.InstanceReboot, urls.BaseUrl, zoneId, instanceId)
	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) DeleteInstance(zoneId, instanceId, name string) (*responses.InstanceCreateResponse, error) {
	var result responses.InstanceCreateResponse
	url := fmt.Sprintf(urls.InstanceDelete, urls.BaseUrl, zoneId, instanceId)
	body, err := json.Marshal(map[string]string{
		"name": name,
	})
	if err != nil {
		return nil, err
	}
	err = client.handleRequest(http.MethodDelete, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Instance Metrics
func (client *Client) GetInstanceMetrics(zoneId, instanceId string, metrics []string, time int, aggregator string) (*responses.InstanceMetricsResponse, error) {
	var result responses.InstanceMetricsResponse
	url := fmt.Sprintf(urls.InstanceMetricsURL, urls.BaseUrl, zoneId, instanceId)
	body, err := json.Marshal(map[string]interface{}{
		"metrics":    metrics,
		"time":       time,
		"aggregator": aggregator,
	})
	if err != nil {
		return nil, err
	}
	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Snapshot Create
func (client *Client) CreateInstanceSnapshot(zoneId, instanceId, name string) (*responses.InstanceSnapshotCreateResponse, error) {
	var result responses.InstanceSnapshotCreateResponse
	url := fmt.Sprintf(urls.InstanceSnapshotCreateURL, urls.BaseUrl, zoneId, instanceId)
	body, err := json.Marshal(map[string]string{
		"instance_id": instanceId,
		"name":        name,
	})
	if err != nil {
		return nil, err
	}
	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Snapshot Delete
func (client *Client) DeleteInstanceSnapshot(zoneId, instanceId, snapshotId string) (*responses.InstanceSnapshotActionResponse, error) {
	var result responses.InstanceSnapshotActionResponse
	url := fmt.Sprintf(urls.InstanceSnapshotDeleteURL, urls.BaseUrl, zoneId, instanceId, snapshotId)
	err := client.handleRequest(http.MethodDelete, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Snapshot Revert
func (client *Client) RevertInstanceSnapshot(zoneId, instanceId, snapshotId string) (*responses.InstanceSnapshotActionResponse, error) {
	var result responses.InstanceSnapshotActionResponse
	url := fmt.Sprintf(urls.InstanceSnapshotRevertURL, urls.BaseUrl, zoneId, instanceId, snapshotId)

	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Volume Service Offering List
func (client *Client) ListInstanceVolumeServiceOfferings(zoneId string) (*responses.InstanceVolumeServiceOfferingListResponse, error) {
	var result responses.InstanceVolumeServiceOfferingListResponse
	url := fmt.Sprintf(urls.InstanceVolumeServiceOfferingListURL, urls.BaseUrl, zoneId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Volume List
func (client *Client) ListInstanceVolumes(zoneId string) (*responses.InstanceVolumeListResponse, error) {
	var result responses.InstanceVolumeListResponse
	url := fmt.Sprintf(urls.InstanceVolumeListURL, urls.BaseUrl, zoneId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Volume Create
func (client *Client) CreateInstanceVolume(zoneId, serviceOfferingId string, size int, name string) (*responses.InstanceVolumeCreateResponse, error) {
	var result responses.InstanceVolumeCreateResponse
	url := fmt.Sprintf(urls.InstanceVolumeCreateURL, urls.BaseUrl, zoneId)
	body, err := json.Marshal(map[string]interface{}{
		"service_offering_id": serviceOfferingId,
		"size":                size,
		"name":                name,
	})
	if err != nil {
		return nil, err
	}
	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Volume Delete
func (client *Client) DeleteInstanceVolume(zoneId, volumeId string) (*responses.InstanceVolumeActionResponse, error) {
	var result responses.InstanceVolumeActionResponse
	url := fmt.Sprintf(urls.InstanceVolumeDeleteURL, urls.BaseUrl, zoneId, volumeId)
	err := client.handleRequest(http.MethodDelete, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Volume Detach
func (client *Client) DetachInstanceVolume(zoneId, volumeId, instanceId string) (*responses.InstanceVolumeActionResponse, error) {
	var result responses.InstanceVolumeActionResponse
	url := fmt.Sprintf(urls.InstanceVolumeDetachURL, urls.BaseUrl, zoneId, volumeId, instanceId)

	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Volume Attach
func (client *Client) AttachInstanceVolume(zoneId, volumeId, instanceId string) (*responses.InstanceVolumeActionResponse, error) {
	var result responses.InstanceVolumeActionResponse
	url := fmt.Sprintf(urls.InstanceVolumeAttachURL, urls.BaseUrl, zoneId, volumeId, instanceId)
	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Show Instance
func (client *Client) ShowInstance(zoneId, instanceId string) (*responses.InstanceShowResponse, error) {
	var result responses.InstanceShowResponse
	url := fmt.Sprintf(urls.InstanceShow, urls.BaseUrl, zoneId, instanceId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
