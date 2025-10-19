// Package http provides virtual machine instance management operations for the Virak Cloud API.
//
// This file contains client methods for managing virtual machine instances, including
// lifecycle operations (create, start, stop, delete), snapshot management, volume
// operations, and performance monitoring. All operations are performed within
// the context of a specific zone.
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"net/http"
)

// ListInstances retrieves all virtual machine instances in the specified zone.
//
// Returns a comprehensive list of all instances accessible to the user within
// the specified zone, including their current status, configuration details,
// network information, and resource allocation.
//
// Parameters:
//   - zoneId: The unique identifier of the zone to query for instances
//
// Returns:
//   - *responses.InstanceListResponse: List of instances with detailed metadata
//   - error: API error or network failure
//
// Example:
//   instances, err := client.ListInstances("zone-123")
//   if err != nil {
//       return fmt.Errorf("failed to list instances: %w", err)
//   }
//   for _, instance := range instances.Data {
//       fmt.Printf("Instance: %s (Status: %s)\n", instance.Name, instance.Status)
//   }
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

// CreateInstance creates a new virtual machine instance with the specified configuration.
//
// This method provisions a new VM with the specified resources, operating system,
// and network connectivity. The instance creation is typically an asynchronous
// operation that returns immediately with a job ID for tracking progress.
//
// Parameters:
//   - zoneId: Zone where the instance will be created
//   - serviceOfferingId: Service offering defining CPU, memory, and storage specifications
//   - vmImageId: VM image ID for the operating system template
//   - networkIds: List of network IDs to connect the instance to
//   - name: Human-readable name for the instance
//
// Returns:
//   - *responses.InstanceCreateResponse: Creation response with instance details and job ID
//   - error: Validation error or API failure
//
// Example:
//   resp, err := client.CreateInstance(
//       "zone-123", "offering-456", "image-789",
//       []string{"network-abc"}, "my-instance",
//   )
//   if err != nil {
//       return fmt.Errorf("failed to create instance: %w", err)
//   }
//   fmt.Printf("Instance created with ID: %s\n", resp.Data.ID)
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

// StartInstance starts a stopped virtual machine instance.
//
// Powers on the instance and makes it available for network connections.
// This operation may take some time as the instance boots up.
//
// Parameters:
//   - zoneId: Zone containing the instance
//   - instanceId: Instance to start
//
// Returns:
//   - *responses.InstanceCreateResponse: Updated instance status with job ID
//   - error: API error if start fails
func (client *Client) StartInstance(zoneId, instanceId string) (*responses.InstanceCreateResponse, error) {
	var result responses.InstanceCreateResponse
	url := fmt.Sprintf(urls.InstanceStart, urls.BaseUrl, zoneId, instanceId)
	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// StopInstance stops a running virtual machine instance.
//
// Gracefully shuts down the instance. If forced is true, performs an immediate
// power off which may result in data loss (similar to pulling the power cord).
//
// Parameters:
//   - zoneId: Zone containing the instance
//   - instanceId: Instance to stop
//   - forced: Whether to force immediate shutdown (true) or graceful shutdown (false)
//
// Returns:
//   - *responses.InstanceCreateResponse: Updated instance status with job ID
//   - error: API error if stop fails
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

// RebootInstance restarts a running virtual machine instance.
//
// Performs a graceful restart of the instance, equivalent to a reboot command
// from within the operating system.
//
// Parameters:
//   - zoneId: Zone containing the instance
//   - instanceId: Instance to reboot
//
// Returns:
//   - *responses.InstanceCreateResponse: Updated instance status with job ID
//   - error: API error if reboot fails
func (client *Client) RebootInstance(zoneId, instanceId string) (*responses.InstanceCreateResponse, error) {
	var result responses.InstanceCreateResponse
	url := fmt.Sprintf(urls.InstanceReboot, urls.BaseUrl, zoneId, instanceId)
	err := client.handleRequest(http.MethodPost, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteInstance permanently removes a virtual machine instance.
//
// This operation is irreversible and will delete all data, storage volumes,
// and network configurations associated with the instance. The name parameter
// must match the instance name exactly as a safety confirmation.
//
// Parameters:
//   - zoneId: Zone containing the instance
//   - instanceId: Instance to delete
//   - name: Instance name (must match exactly for safety confirmation)
//
// Returns:
//   - *responses.InstanceCreateResponse: Deletion confirmation with job ID
//   - error: API error if deletion fails or name doesn't match
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

// ShowInstance retrieves detailed information about a specific instance.
//
// Returns comprehensive details about a single instance including hardware
// specifications, network configuration, storage volumes, security groups,
// and current operational status.
//
// Parameters:
//   - zoneId: Zone containing the instance
//   - instanceId: Instance to retrieve details for
//
// Returns:
//   - *responses.InstanceShowResponse: Detailed instance information
//   - error: API error if instance doesn't exist or access is denied
//
// Example:
//   instance, err := client.ShowInstance("zone-123", "instance-456")
//   if err != nil {
//       return fmt.Errorf("failed to get instance details: %w", err)
//   }
//   fmt.Printf("Instance %s: %s vCPU, %s RAM\n",
//       instance.Data.Name, instance.Data.ServiceOffering.Hardware.Cpu,
//       instance.Data.ServiceOffering.Hardware.Memory)
func (client *Client) ShowInstance(zoneId, instanceId string) (*responses.InstanceShowResponse, error) {
	var result responses.InstanceShowResponse
	url := fmt.Sprintf(urls.InstanceShow, urls.BaseUrl, zoneId, instanceId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
