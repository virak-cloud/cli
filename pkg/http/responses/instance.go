// Package responses defines data structures for API responses in the Virak Cloud CLI.
//
// This file contains response type definitions for virtual machine instance operations,
// including instance metadata, service offerings, VM images, volumes, snapshots,
// and metrics. These structures provide type-safe access to API response data.
package responses

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// InstanceListResponse represents the response from listing instances.
//
// Contains a collection of instance objects returned by the ListInstances API call.
type InstanceListResponse struct {
	Data []Instance `json:"data"`
}

// Instance represents a virtual machine instance with all its configuration and status.
//
// This struct provides comprehensive information about an instance including
// hardware specifications, network configuration, storage, and operational status.
type Instance struct {
	// ID is the unique identifier for the instance.
	ID string `json:"id"`
	// CustomerID is the unique identifier for the customer who owns the instance.
	CustomerID string `json:"customer_id"`
	// Name is the human-readable name assigned to the instance.
	Name string `json:"name"`
	// ZoneID is the unique identifier of the zone where the instance is located.
	ZoneID string `json:"zone_id"`
	// Created indicates whether the instance has been fully provisioned.
	Created bool `json:"created"`
	// TemplateID is the ID of the template used (if applicable).
	TemplateID *string `json:"template_id"`
	// VMImage contains information about the operating system image.
	VMImage *InstanceVMImage `json:"vm_image"`
	// Zone contains zone information where the instance is deployed.
	Zone *InstanceZone `json:"zone"`
	// ServiceOffering contains the hardware specifications for the instance.
	ServiceOffering *InstanceServiceOffering `json:"service_offering"`
	// DiskOfferingID is the ID of the disk offering used (if applicable).
	DiskOfferingID *string `json:"disk_offering_id"`
	// ServiceOfferingID is the ID of the service offering used.
	ServiceOfferingID string `json:"service_offering_id"`
	// Status is the current operational status of the instance.
	Status string `json:"status"`
	// InstanceStatus provides detailed status information.
	InstanceStatus string `json:"instance_status"`
	// Password is the initial password for instance access.
	Password string `json:"password"`
	// Username is the default username for instance access.
	Username string `json:"username"`
	// CreatedAt is the timestamp when the instance was created (Unix epoch).
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt is the timestamp when the instance was last updated (Unix epoch).
	UpdatedAt int64 `json:"updated_at"`
	// KubernetesClusterID is the ID of the Kubernetes cluster (if part of one).
	KubernetesClusterID *string `json:"kubernetes_cluster_id"`
	// Metadata contains additional instance metadata.
	Metadata []interface{} `json:"metadata"`
	// DataVolumes contains information about attached data volumes.
	DataVolumes []interface{} `json:"data_volumes"`
	// Snapshot contains information about instance snapshots.
	Snapshot []InstanceSnapshot `json:"snapshot"`
}

type InstanceServiceOfferingListResponse struct {
	Data []InstanceServiceOffering `json:"data"`
}

type InstanceServiceOffering struct {
	ID                    string                           `json:"id"`
	Name                  string                           `json:"name"`
	Category              string                           `json:"category"`
	Suggested             bool                             `json:"suggested"`
	Hardware              *InstanceServiceOfferingHardware `json:"hardware"`
	IsAvailable           bool                             `json:"is_available"`
	HasImageRequirement   bool                             `json:"has_image_requirement"`
	IsPublic              bool                             `json:"is_public"`
	HourlyPrice           *InstanceServiceOfferingPrice    `json:"hourly_price"`
	HourlyPriceNoDiscount *InstanceServiceOfferingPrice    `json:"hourly_price_no_discount"`
	Description           *string                          `json:"description"`
}

type InstanceServiceOfferingPrice struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}

type InstanceServiceOfferingHardware struct {
	CPUCore        int `json:"cpu_core"`
	MemoryMB       int `json:"memory_mb"`
	RootDiskSizeGB int `json:"root_disk_size_gB"`
	CPUSpeedMHz    int `json:"cpu_speed_MHz"`
	NetworkRate    int `json:"network_rate"`
	DiskIOPS       int `json:"disk_iops"`
}

type InstanceVMImageListResponse struct {
	Data []InstanceVMImage `json:"data"`
}

type InstanceVMImage struct {
	ID                   string                              `json:"id"`
	Type                 string                              `json:"type"`
	Name                 string                              `json:"name"`
	IsAvailable          bool                                `json:"is_available"`
	DisplayText          string                              `json:"display_text"`
	NameOrginal          string                              `json:"name_orginal"`
	ReadyToUseApp        bool                                `json:"ready_to_use_app"`
	ReadyToUseAppName    *string                             `json:"ready_to_use_app_name"`
	ReadyToUseAppVersion *string                             `json:"ready_to_use_app_version"`
	OSType               string                              `json:"os_type"`
	OSName               string                              `json:"os_name"`
	OSVersion            string                              `json:"os_version"`
	HardwareRequirement  *InstanceVMImageHardwareRequirement `json:"hardware_requirement"`
	Category             string                              `json:"category"`
}

type InstanceVMImageHardwareRequirement struct {
	CPUNumber    IntString `json:"cpunumber"`
	CPUSpeed     int       `json:"cpuspeed"`
	Memory       IntString `json:"memory"`
	RootDiskSize IntString `json:"rootdisksize"`
}

type InstanceZone struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	IsPublic   bool   `json:"is_public"`
	IsFeatured bool   `json:"is_featured"`
	IsReady    bool   `json:"is_ready"`
}

type InstanceCreateResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

// Instance Metrics

// InstanceMetricsResponse represents the metrics response for an instance.
type InstanceMetricsResponse struct {
	Data []InstanceMetricColumn `json:"data"`
}

type InstanceMetricColumn struct {
	Column string                `json:"column"`
	Values []InstanceMetricValue `json:"values"`
}

type InstanceMetricValue struct {
	Value float64 `json:"value"`
	Time  string  `json:"time"`
}

// Snapshot Create

type InstanceSnapshotCreateResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

// Snapshot Delete & Revert

type InstanceSnapshotActionResponse struct {
	Data InstanceSnapshotActionResult `json:"data"`
}

type InstanceSnapshotActionResult struct {
	Success bool `json:"success"`
}

// Volume Service Offering List

type InstanceVolumeServiceOfferingListResponse struct {
	Data []InstanceVolumeServiceOffering `json:"data"`
}

type InstanceVolumeServiceOffering struct {
	ID          string `json:"id"`
	Size        string `json:"size"`
	Price       string `json:"price"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	IsFeatured  bool   `json:"is_featured"`
}

// Volume List

type InstanceVolumeListResponse struct {
	Data []InstanceVolume `json:"data"`
}

type InstanceVolume struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Size   int    `json:"size"`
	Status string `json:"status"`
}

// Volume Create

type InstanceVolumeCreateResponse struct {
	Data InstanceVolume `json:"data"`
}

// Volume Delete, Attach, Detach

type InstanceVolumeActionResponse struct {
	Data InstanceVolumeActionResult `json:"data"`
}

type InstanceVolumeActionResult struct {
	Success bool `json:"success"`
}

// Custom type to handle int or string JSON values
// Use for fields that may be string or int in JSON
// Example: "cpunumber": "4" or "cpunumber": 4
// Example: "memory": "2048" or "memory": 2048
type IntString int

func (i *IntString) UnmarshalJSON(b []byte) error {
	var intVal int
	if err := json.Unmarshal(b, &intVal); err == nil {
		*i = IntString(intVal)
		return nil
	}
	var strVal string
	if err := json.Unmarshal(b, &strVal); err == nil {
		val, err := strconv.Atoi(strVal)
		if err != nil {
			return err
		}
		*i = IntString(val)
		return nil
	}
	return fmt.Errorf("IntString: value is not int or string: %s", string(b))
}

type InstanceSnapshotListResponse struct {
	Data []InstanceSnapshot `json:"data"`
}

type InstanceSnapshot struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	CreatedAt int64   `json:"created_at"`
	Current   bool    `json:"current"`
	ParentID  *string `json:"parent_id"`
}

type InstanceShowResponse struct {
	Data Instance `json:"data"`
}

type VolumeListResponse struct {
	Data []Volume `json:"data"`
}

type Volume struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	ZoneID string `json:"zone_id"`
	Size   int    `json:"size"`
	Status string `json:"status"`
}
