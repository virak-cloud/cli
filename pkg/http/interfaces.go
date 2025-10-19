package http

import "github.com/virak-cloud/cli/pkg/http/responses"

// InstanceManager defines the interface for all instance-related operations
// in the Virak Cloud API. This interface provides a contract for instance
// management functionality and enables easy testing and mocking.
//
// Implementations of this interface should be thread-safe and handle
// authentication, error handling, and request validation internally.
type InstanceManager interface {
	// Lifecycle Management
	
	// ListInstances retrieves all instances in the specified zone.
	// Returns instances with their current status and configuration.
	ListInstances(zoneID string) (*responses.InstanceListResponse, error)
	
	// ShowInstance retrieves detailed information about a specific instance.
	// Returns comprehensive instance data including network configuration.
	ShowInstance(zoneID, instanceID string) (*responses.InstanceShowResponse, error)
	
	// CreateInstance creates a new instance with the specified configuration.
	// Returns the creation response with the new instance ID.
	CreateInstance(zoneID, offeringID, imageID string, networkIDs []string, name string) (*responses.InstanceCreateResponse, error)
	
	// StartInstance starts a stopped instance.
	// Returns the operation result with updated instance status.
	StartInstance(zoneID, instanceID string) (*responses.InstanceCreateResponse, error)
	
	// StopInstance stops a running instance.
	// The forced parameter determines whether to force immediate shutdown.
	StopInstance(zoneID, instanceID string, forced bool) (*responses.InstanceCreateResponse, error)
	
	// RebootInstance restarts a running instance.
	// Returns the operation result with updated instance status.
	RebootInstance(zoneID, instanceID string) (*responses.InstanceCreateResponse, error)
	
	// DeleteInstance permanently removes an instance.
	// The name parameter must match the instance name for confirmation.
	DeleteInstance(zoneID, instanceID, name string) (*responses.InstanceCreateResponse, error)
	
	// RebuildInstance rebuilds an instance with a new VM image.
	// This preserves the instance ID but replaces the operating system.
	RebuildInstance(zoneID, instanceID, imageID string) (*responses.InstanceCreateResponse, error)
	
	// Monitoring and Metrics
	
	// GetInstanceMetrics retrieves performance metrics for an instance.
	// The metrics parameter specifies which metrics to collect.
	GetInstanceMetrics(zoneID, instanceID string, metrics []string, time int, aggregator string) (*responses.InstanceMetricsResponse, error)
	
	// Snapshot Management
	
	// CreateInstanceSnapshot creates a point-in-time snapshot of an instance.
	// Snapshots can be used for backup or to create new instances.
	CreateInstanceSnapshot(zoneID, instanceID, name string) (*responses.InstanceSnapshotCreateResponse, error)
	
	// DeleteInstanceSnapshot removes a snapshot.
	// This permanently deletes the snapshot and frees storage space.
	DeleteInstanceSnapshot(zoneID, instanceID, snapshotID string) (*responses.InstanceSnapshotActionResponse, error)
	
	// RevertInstanceSnapshot restores an instance to a previous snapshot state.
	// This replaces the current instance state with the snapshot state.
	RevertInstanceSnapshot(zoneID, instanceID, snapshotID string) (*responses.InstanceSnapshotActionResponse, error)
	
	// Volume Management
	
	// ListInstanceVolumes retrieves all volumes attached to an instance.
	// Returns both root and additional volumes with their specifications.
	ListInstanceVolumes(zoneID string) (*responses.InstanceVolumeListResponse, error)
	
	// CreateInstanceVolume creates a new storage volume for an instance.
	// The size parameter is specified in GB.
	CreateInstanceVolume(zoneID, offeringID string, size int, name string) (*responses.InstanceVolumeCreateResponse, error)
	
	// DeleteInstanceVolume removes a volume from an instance.
	// This permanently deletes the volume and all data on it.
	DeleteInstanceVolume(zoneID, volumeID string) (*responses.InstanceVolumeActionResponse, error)
	
	// AttachInstanceVolume attaches a volume to an instance.
	// The volume becomes accessible to the instance after attachment.
	AttachInstanceVolume(zoneID, volumeID, instanceID string) (*responses.InstanceVolumeActionResponse, error)
	
	// DetachInstanceVolume detaches a volume from an instance.
	// The volume remains available but is no longer accessible to the instance.
	DetachInstanceVolume(zoneID, volumeID, instanceID string) (*responses.InstanceVolumeActionResponse, error)
	
	// Resource Information
	
	// ListInstanceServiceOfferings retrieves available instance configurations.
	// Returns CPU, memory, and storage options with pricing information.
	ListInstanceServiceOfferings(zoneID string) (*responses.InstanceServiceOfferingListResponse, error)
	
	// ListInstanceVMImages retrieves available operating system images.
	// Returns both public and private images with their specifications.
	ListInstanceVMImages(zoneID string) (*responses.InstanceVMImageListResponse, error)
	
	// ListInstanceVolumeServiceOfferings retrieves available volume configurations.
	// Returns storage types and performance options with pricing.
	ListInstanceVolumeServiceOfferings(zoneID string) (*responses.InstanceVolumeServiceOfferingListResponse, error)
}

// NetworkManager defines the interface for all network-related operations
// in the Virak Cloud API. This interface provides a contract for network
// management functionality including L2/L3 networks, firewall rules, and
// load balancers.
type NetworkManager interface {
	// Network Management
	
	// ListNetworks retrieves all networks in the specified zone.
	// Returns both L2 and L3 networks with their configurations.
	ListNetworks(zoneID string) (*responses.NetworkListResponse, error)
	
	// ShowNetwork retrieves detailed information about a specific network.
	// Returns network configuration, connected instances, and status.
	ShowNetwork(zoneID, networkID string) (*responses.NetworkShowResponse, error)
	
	// CreateL2Network creates a new Layer 2 network.
	// L2 networks operate at the data link layer for switching within a broadcast domain.
	CreateL2Network(zoneID, offeringID, name string) (*responses.NetworkCreateResponse, error)
	
	// CreateL3Network creates a new Layer 3 network.
	// L3 networks operate at the network layer with routing between broadcast domains.
	CreateL3Network(zoneID, offeringID, name, gateway, netmask string) (*responses.NetworkCreateResponse, error)
	
	// DeleteNetwork permanently removes a network.
	// All connected instances must be disconnected before deletion.
	DeleteNetwork(zoneID, networkID string) (*responses.NetworkDeleteResponse, error)
	
	// Instance Network Management
	
	// ConnectInstanceToNetwork attaches an instance to a network.
	// The instance will receive network connectivity through this connection.
	ConnectInstanceToNetwork(zoneID, networkID, instanceID string) (*responses.InstanceNetworkActionResponse, error)
	
	// DisconnectInstanceFromNetwork removes an instance from a network.
	// The instance will lose network connectivity through this connection.
	DisconnectInstanceFromNetwork(zoneID, networkID, instanceID, instanceNetworkID string) (*responses.InstanceNetworkActionResponse, error)
	
	// ListNetworkInstances retrieves all instances connected to a network.
	// Returns instance network connections with IP addresses and status.
	ListNetworkInstances(zoneID, networkID string, instanceID string) (*responses.InstanceNetworkListResponse, error)
	
	// Firewall Management
	
	// ListIPv4FirewallRules retrieves all IPv4 firewall rules for a network.
	// Returns rules with protocol, ports, and source/destination configurations.
	ListIPv4FirewallRules(zoneID, networkID string) (*responses.IPv4FirewallRuleListResponse, error)
	
	// CreateIPv4FirewallRule creates a new IPv4 firewall rule.
	// The body parameter contains rule configuration (protocol, ports, etc.).
	CreateIPv4FirewallRule(zoneID, networkID string, body map[string]interface{}) (*responses.IPv4FirewallRuleActionResponse, error)
	
	// DeleteIPv4FirewallRule removes an IPv4 firewall rule.
	// This permanently deletes the rule and stops filtering traffic.
	DeleteIPv4FirewallRule(zoneID, networkID, ruleID string) (*responses.IPv4FirewallRuleActionResponse, error)
	
	// ListIPv6FirewallRules retrieves all IPv6 firewall rules for a network.
	// Returns rules with protocol, ports, and source/destination configurations.
	ListIPv6FirewallRules(zoneID, networkID string) (*responses.IPv6FirewallRuleListResponse, error)
	
	// CreateIPv6FirewallRule creates a new IPv6 firewall rule.
	// The body parameter contains rule configuration (protocol, ports, etc.).
	CreateIPv6FirewallRule(zoneID, networkID string, body map[string]interface{}) (*responses.IPv6FirewallRuleActionResponse, error)
	
	// DeleteIPv6FirewallRule removes an IPv6 firewall rule.
	// This permanently deletes the rule and stops filtering traffic.
	DeleteIPv6FirewallRule(zoneID, networkID, ruleID string) (*responses.IPv6FirewallRuleActionResponse, error)
	
	// Public IP Management
	
	// ListNetworkPublicIps retrieves all public IPs for a network.
	// Returns both assigned and unassigned public IP addresses.
	ListNetworkPublicIps(zoneID, networkID string) (*responses.NetworkPublicIpListResponse, error)
	
	// AssociateNetworkPublicIp assigns a new public IP to a network.
	// The public IP can be used for NAT and load balancing.
	AssociateNetworkPublicIp(zoneID, networkID string) (*responses.NetworkPublicIpActionResponse, error)
	
	// DisassociateNetworkPublicIp removes a public IP from a network.
	// The public IP is released and becomes unavailable.
	DisassociateNetworkPublicIp(zoneID, networkID, networkPublicIPID string) (*responses.NetworkPublicIpActionResponse, error)
	
	// EnableNetworkPublicIpStaticNat configures static NAT for a public IP.
	// Routes all traffic for the public IP to the specified instance.
	EnableNetworkPublicIpStaticNat(zoneID, networkID, networkPublicIPID, instanceID string) (*responses.NetworkPublicIpActionResponse, error)
	
	// DisableNetworkPublicIpStaticNat removes static NAT configuration.
	// The public IP is no longer routed to the instance.
	DisableNetworkPublicIpStaticNat(zoneID, networkID, networkPublicIPID string) (*responses.NetworkPublicIpActionResponse, error)
	
	// VPN Management
	
	// GetNetworkVpnDetails retrieves VPN configuration for a network.
	// Returns connection details including IP, credentials, and status.
	GetNetworkVpnDetails(zoneID, networkID string) (*responses.NetworkVpnDetailResponse, error)
	
	// EnableNetworkVpn activates VPN service for a network.
	// Generates VPN credentials and configuration.
	EnableNetworkVpn(zoneID, networkID string) (*responses.NetworkVpnSuccessResponse, error)
	
	// DisableNetworkVpn deactivates VPN service for a network.
	// Terminates VPN connections and releases resources.
	DisableNetworkVpn(zoneID, networkID string) (*responses.NetworkVpnSuccessResponse, error)
	
	// UpdateNetworkVpnCredentials regenerates VPN credentials.
	// Creates new username/password for VPN access.
	UpdateNetworkVpnCredentials(zoneID, networkID string) (*responses.NetworkVpnSuccessResponse, error)
	
	// Load Balancer Management
	
	// ListLoadBalancerRules retrieves all load balancer rules for a network.
	// Returns rules with algorithm, ports, and assigned instances.
	ListLoadBalancerRules(zoneID, networkID string) (*responses.LoadBalancerRuleListResponse, error)
	
	// CreateLoadBalancerRule creates a new load balancer rule.
	// Configures traffic distribution across multiple instances.
	CreateLoadBalancerRule(zoneID, networkID, publicIPID, name, algorithm string, publicPort, privatePort int) (*responses.SuccessResponse, error)
	
	// DeleteLoadBalancerRule removes a load balancer rule.
	// Stops load balancing for the specified rule.
	DeleteLoadBalancerRule(zoneID, networkID, ruleID string) (*responses.SuccessResponse, error)
	
	// AssignLoadBalancerRule adds instances to a load balancer rule.
	// Traffic will be distributed to all assigned instances.
	AssignLoadBalancerRule(zoneID, networkID, ruleID string, instanceNetworkIDs []string) (*responses.SuccessResponse, error)
	
	// DeassignLoadBalancerRule removes an instance from a load balancer rule.
	// Traffic will no longer be distributed to the specified instance.
	DeassignLoadBalancerRule(zoneID, networkID, ruleID, instanceNetworkID string) (*responses.SuccessResponse, error)
	
	// HAProxy Management
	
	// GetHaproxyLive retrieves live HAProxy statistics.
	// Returns current load balancer performance metrics.
	GetHaproxyLive(zoneID, networkID string) (*responses.HaproxyLiveResponse, error)
	
	// GetHaproxyLog retrieves HAProxy access logs.
	// Returns recent log entries for troubleshooting.
	GetHaproxyLog(zoneID, networkID string) (*responses.HaproxyLogResponse, error)
	
	// Service Offerings
	
	// ListNetworkServiceOfferings retrieves available network configurations.
	// Returns L2 and L3 network offerings with pricing.
	ListNetworkServiceOfferings(zoneID string) (*responses.NetworkServiceOfferingListResponse, error)
	
	// GetL2NetworkServiceOfferings retrieves only L2 network offerings.
	// Filters offerings to return only Layer 2 network types.
	GetL2NetworkServiceOfferings(zoneID string) (*responses.NetworkServiceOfferingListResponse, error)
	
	// GetL3NetworkServiceOfferings retrieves only L3 network offerings.
	// Filters offerings to return only Layer 3 network types.
	GetL3NetworkServiceOfferings(zoneID string) (*responses.NetworkServiceOfferingListResponse, error)
}

// ZoneManager defines the interface for zone and datacenter operations.
// This interface provides access to zone information and resources.
type ZoneManager interface {
	// GetZoneList retrieves all available zones (datacenters).
	// Returns zones with their locations, capabilities, and status.
	GetZoneList() (*responses.DataCenter, error)
	
	// GetZoneActiveServices retrieves active services for a specific zone.
	// Returns available services and their current status.
	GetZoneActiveServices(zoneID string) (*responses.ZoneActiveServicesResponse, error)
	
	// GetZoneCustomerResource retrieves customer resources for a zone.
	// Returns resource usage and allocation information.
	GetZoneCustomerResource(zoneID string) (*responses.CustomerResourceResponse, error)
	
	// GetZoneNetworks retrieves networks for a specific zone.
	// Returns all networks accessible in the specified zone.
	GetZoneNetworks(zoneID string) (*responses.ZoneNetworksResponse, error)
}

// DNSManager defines the interface for DNS management operations.
// This interface provides domain and DNS record management.
type DNSManager interface {
	// Domain Management
	
	// GetDomains retrieves all DNS domains.
	// Returns domains with their configuration and status.
	GetDomains() (*responses.DomainList, error)
	
	// CreateDomain creates a new DNS domain.
	// The domain parameter is the domain name to create.
	CreateDomain(domain string) (*responses.DnsMessage, error)
	
	// GetDomain retrieves a specific DNS domain.
	// Returns domain configuration and DNS records.
	GetDomain(domain string) (*responses.DomainShow, error)
	
	// DeleteDomain removes a DNS domain.
	// This permanently deletes the domain and all its records.
	DeleteDomain(domain string) (*responses.DnsMessage, error)
	
	// Record Management
	
	// GetRecords retrieves all DNS records for a domain.
	// Returns A, CNAME, MX, and other record types.
	GetRecords(domain string) (*responses.RecordList, error)
	
	// CreateRecord creates a new DNS record.
	// Supports A, CNAME, MX, SRV, CAA, and TLSA record types.
	CreateRecord(domain, record, recordType, content string, ttl, priority, weight, port, flags int, tag string, license, choicer, match int) (*responses.DnsMessage, error)
	
	// UpdateRecord modifies an existing DNS record.
	// Updates record content, TTL, and type-specific fields.
	UpdateRecord(domain, record, recordType, contentID, newContent string, newTTL, priority, weight, port, flags int, tag string, license, choicer, match int) (*responses.DnsMessage, error)
	
	// DeleteRecord removes a DNS record.
	// This permanently deletes the specified record.
	DeleteRecord(domain, record, recordType, contentID string) (*responses.DnsMessage, error)
	
	// Events
	
	// GetDNSEvents retrieves DNS events and audit logs.
	// Returns recent DNS operations and changes.
	GetDNSEvents() (*responses.DNSEventsResponse, error)
}

// Ensure Client implements all interfaces
var _ InstanceManager = (*Client)(nil)
var _ NetworkManager = (*Client)(nil)
var _ ZoneManager = (*Client)(nil)
var _ DNSManager = (*Client)(nil)