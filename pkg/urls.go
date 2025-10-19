// Package pkg defines centralized API endpoint URLs and configuration for the Virak Cloud API.
//
// This package contains all API endpoint templates used by the HTTP client to communicate
// with the Virak Cloud REST API. The endpoints are organized by service category and
// use template strings that can be formatted with specific parameters.
//
// URL Configuration:
//   - BaseUrl: Default base URL for API requests (configurable at build time)
//   - LoginUrl: URL for generating CLI authentication tokens
//
// Build-time Configuration:
//   Both BaseUrl and LoginUrl can be overridden during build using ldflags:
//     go build -ldflags "-X github.com/virak-cloud/cli/pkg.BaseUrl=https://custom.example.com"
//
// Endpoint Categories:
//   - Zone Management: Datacenter and zone information
//   - Instance Management: Virtual machine lifecycle operations
//   - Network Management: Network configuration, firewall, and load balancing
//   - DNS Management: Domain and DNS record operations
//   - Kubernetes Management: Container orchestration operations
//   - Object Storage: Bucket and storage operations
//   - User Management: Account and authentication operations
//
// URL Template Format:
//   All endpoint constants use printf-style formatting with placeholders:
//     fmt.Sprintf(EndpointTemplate, BaseUrl, zoneID, resourceID, ...)
//
// Thread Safety:
//   All constants and variables in this package are immutable and safe for
//   concurrent access by multiple goroutines.
package pkg

// Global API configuration that can be overridden at build time via -ldflags.
//
// These variables provide the foundation for all API communication. They can be
// customized during the build process to support different environments or
// deployment configurations.
var (
	// LoginUrl is the URL for creating a new CLI authentication token.
	// The %s placeholder should be replaced with a unique token name.
	//
	// Template: "https://panel.virakcloud.com/web-service/create?name=cli-token-%s&abilities=*"
	//
	// Parameters:
	//   - %s: Token name (typically "cli-token-{timestamp}" or similar)
	//
	// Example:
	//   tokenURL := fmt.Sprintf(LoginUrl, "cli-token-2023-001")
	//   // Result: https://panel.virakcloud.com/web-service/create?name=cli-token-2023-001&abilities=*
	LoginUrl = "https://panel.virakcloud.com/web-service/create?name=cli-token-%s&abilities=*"

	// BaseUrl is the base URL for all Virak Cloud API endpoints.
	// This URL serves as the prefix for all API endpoint templates.
	//
	// Default: "https://public-api.virakcloud.com"
	//
	// Environment Configuration:
	//   - Development: Can be set to "http://localhost:1410"
	//   - Staging: Can be set to "https://staging-api.virakcloud.com"
	//   - Production: Default value
	//
	// Build-time Override:
	//   go build -ldflags "-X github.com/virak-cloud/cli/pkg.BaseUrl=https://custom.example.com"
	BaseUrl = "https://public-api.virakcloud.com"
)

// API endpoint templates organized by service category.
//
// Each constant defines a printf-style template for constructing API endpoints.
// The templates use placeholders that should be replaced with actual values:
//   - %s: String values (base URL, zone ID, resource ID, etc.)
//   - %d: Numeric values (port numbers, sizes, etc.)
//
// Common Parameters:
//   - Base URL: First parameter in all endpoints
//   - Zone ID: Second parameter for zone-specific operations
//   - Resource ID: Varies by endpoint type (instance, network, domain, etc.)
//
// Example Usage:
//   url := fmt.Sprintf(InstanceList, BaseUrl, "zone-123")
//   // Result: https://public-api.virakcloud.com/zone/zone-123/instance
const (
	// ZoneList is the endpoint for listing all zones.
	ZoneList               string = "%s/zones"
	// ZoneActiveServicesList is the endpoint for listing active services in a zone.
	ZoneActiveServicesList string = "%s/zone/%s"
	// ZoneNetworkList is the endpoint for listing networks in a zone.
	ZoneNetworkList        string = "%s/zone/%s/network"
	// ZoneResourcesList is the endpoint for listing resources in a zone.
	ZoneResourcesList      string = "%s/zone/%s/resources"

	// BucketList is the endpoint for listing all buckets in a zone.
	BucketList       string = "%s/zone/%s/object-storage/buckets"
	// BucketCreate is the endpoint for creating a new bucket.
	BucketCreate     string = "%s/zone/%s/object-storage/buckets"
	// BucketShow is the endpoint for showing a bucket.
	BucketShow       string = "%s/zone/%s/object-storage/buckets/%s"
	// BucketUpdate is the endpoint for updating a bucket.
	BucketUpdate     string = "%s/zone/%s/object-storage/buckets/%s"
	// BucketDelete is the endpoint for deleting a bucket.
	BucketDelete     string = "%s/zone/%s/object-storage/buckets/%s"
	// BucketsEventList is the endpoint for listing all bucket events in a zone.
	BucketsEventList string = "%s/zone/%s/object-storage/events"
	// BucketEventList is the endpoint for listing events for a specific bucket.
	BucketEventList  string = "%s/zone/%s/object-storage/buckets/%s/events"

	// DomainListURL is the endpoint for listing all domains.
	DomainListURL   string = "%s/dns/domains"
	// DomainCreateURL is the endpoint for creating a new domain.
	DomainCreateURL string = "%s/dns/domains"
	// DomainShowURL is the endpoint for showing a domain.
	DomainShowURL   string = "%s/dns/domains/%s"
	// DomainDeleteURL is the endpoint for deleting a domain.
	DomainDeleteURL string = "%s/dns/domains/%s"

	// DNSEventsURL is the endpoint for listing all DNS events.
	DNSEventsURL string = "%s/dns/events"

	// RecordListURL is the endpoint for listing all records for a domain.
	RecordListURL   string = "%s/dns/domains/%s/records"
	// RecordCreateURL is the endpoint for creating a new record.
	RecordCreateURL string = "%s/dns/domains/%s/records"
	// RecordUpdateURL is the endpoint for updating a record.
	RecordUpdateURL string = "%s/dns/domains/%s/records/%s/%s/%s"
	// RecordDeleteURL is the endpoint for deleting a record.
	RecordDeleteURL string = "%s/dns/domains/%s/records/%s/%s/%s"

	// KubernetesClusterList is the endpoint for listing all Kubernetes clusters in a zone.
	KubernetesClusterList          string = "%s/zone/%s/kubernetes"
	// KubernetesClusterCreate is the endpoint for creating a new Kubernetes cluster.
	KubernetesClusterCreate        string = "%s/zone/%s/kubernetes"
	// KubernetesClusterShow is the endpoint for showing a Kubernetes cluster.
	KubernetesClusterShow          string = "%s/zone/%s/kubernetes/%s"
	// KubernetesClusterUpdate is the endpoint for updating a Kubernetes cluster.
	KubernetesClusterUpdate        string = "%s/zone/%s/kubernetes/%s"
	// KubernetesClusterDelete is the endpoint for deleting a Kubernetes cluster.
	KubernetesClusterDelete        string = "%s/zone/%s/kubernetes/%s"
	// KubernetesClusterStart is the endpoint for starting a Kubernetes cluster.
	KubernetesClusterStart         string = "%s/zone/%s/kubernetes/%s/start"
	// KubernetesClusterStop is the endpoint for stopping a Kubernetes cluster.
	KubernetesClusterStop          string = "%s/zone/%s/kubernetes/%s/stop"
	// KubernetesClusterScale is the endpoint for scaling a Kubernetes cluster.
	KubernetesClusterScale         string = "%s/zone/%s/kubernetes/%s/scale"
	// KubernetesClusterEvents is the endpoint for listing events for a specific Kubernetes cluster.
	KubernetesClusterEvents        string = "%s/zone/%s/kubernetes/%s/events"
	// KubernetesServiceEvents is the endpoint for listing all Kubernetes service events in a zone.
	KubernetesServiceEvents        string = "%s/zone/%s/kubernetes/events"
	// KubernetesVersionsList is the endpoint for listing all available Kubernetes versions in a zone.
	KubernetesVersionsList         string = "%s/zone/%s/kubernetes/versions"
	// KubernetesServiceOfferingsList is the endpoint for listing all available Kubernetes service offerings in a zone.
	KubernetesServiceOfferingsList string = "%s/zone/%s/kubernetes/service-offerings"

	// UserSSHKeyList is the endpoint for listing all SSH keys for a user.
	UserSSHKeyList       string = "%s/user/ssh-key"
	// UserSSHKeyCreate is the endpoint for creating a new SSH key.
	UserSSHKeyCreate     string = "%s/user/ssh-key"
	// UserSSHKeyDelete is the endpoint for deleting an SSH key.
	UserSSHKeyDelete     string = "%s/user/ssh-key/%s"
	// UserBalance is the endpoint for getting the user's balance.
	UserBalance          string = "%s/user/finance/wallet"
	// UserPaymentList is the endpoint for listing all payments for a user.
	UserPaymentList      string = "%s/user/finance/payments"
	// UserCostDocumentList is the endpoint for listing all cost documents for a user.
	UserCostDocumentList string = "%s/user/finance/documents"
	// UserTokenAbilities is the endpoint for getting the abilities of a token.
	UserTokenAbilities   string = "%s/user/token-abilities"
	// UserTokenValidate is the endpoint for validating a token.
	UserTokenValidate    string = "%s/user/token"

	// NetworkCreateL3 is the endpoint for creating a new L3 network.
	NetworkCreateL3           string = "%s/zone/%s/network/l3"
	// NetworkCreateL2 is the endpoint for creating a new L2 network.
	NetworkCreateL2           string = "%s/zone/%s/network/l2"
	// NetworkList is the endpoint for listing all networks in a zone.
	NetworkList               string = "%s/zone/%s/network"
	// NetworkShow is the endpoint for showing a network.
	NetworkShow               string = "%s/zone/%s/network/%s"
	// NetworkDelete is the endpoint for deleting a network.
	NetworkDelete             string = "%s/zone/%s/network/%s"
	// NetworkInstanceConnect is the endpoint for connecting an instance to a network.
	NetworkInstanceConnect    string = "%s/zone/%s/network/%s/instance/connect"
	// NetworkInstanceDisconnect is the endpoint for disconnecting an instance from a network.
	NetworkInstanceDisconnect string = "%s/zone/%s/network/%s/instance/disconnect"
	// NetworkInstanceList is the endpoint for listing all instances connected to a network.
	NetworkInstanceList       string = "%s/zone/%s/network/%s/instance"
	// NetworkFirewallIPv4List is the endpoint for listing all IPv4 firewall rules for a network.
	NetworkFirewallIPv4List   string = "%s/zone/%s/network/%s/firewall/ipv4"
	// NetworkFirewallIPv4Create is the endpoint for creating a new IPv4 firewall rule.
	NetworkFirewallIPv4Create string = "%s/zone/%s/network/%s/firewall/ipv4"
	// NetworkFirewallIPv4Delete is the endpoint for deleting an IPv4 firewall rule.
	NetworkFirewallIPv4Delete string = "%s/zone/%s/network/%s/firewall/ipv4/%s"
	// NetworkFirewallIPv6List is the endpoint for listing all IPv6 firewall rules for a network.
	NetworkFirewallIPv6List   string = "%s/zone/%s/network/%s/firewall/ipv6"
	// NetworkFirewallIPv6Create is the endpoint for creating a new IPv6 firewall rule.
	NetworkFirewallIPv6Create string = "%s/zone/%s/network/%s/firewall/ipv6"
	// NetworkFirewallIPv6Delete is the endpoint for deleting an IPv6 firewall rule.
	NetworkFirewallIPv6Delete string = "%s/zone/%s/network/%s/firewall/ipv6/%s"

	// NetworkPublicIpList is the endpoint for listing all public IPs for a network.
	NetworkPublicIpList             string = "%s/zone/%s/network/%s/public-ip"
	// NetworkPublicIpAssociate is the endpoint for associating a new public IP with a network.
	NetworkPublicIpAssociate        string = "%s/zone/%s/network/%s/public-ip"
	// NetworkPublicIpDisassociate is the endpoint for disassociating a public IP from a network.
	NetworkPublicIpDisassociate     string = "%s/zone/%s/network/%s/public-ip/%s"
	// NetworkPublicIpStaticNatEnable is the endpoint for enabling static NAT for a public IP.
	NetworkPublicIpStaticNatEnable  string = "%s/zone/%s/network/%s/public-ip/%s/static-nat"
	// NetworkPublicIpStaticNatDisable is the endpoint for disabling static NAT for a public IP.
	NetworkPublicIpStaticNatDisable string = "%s/zone/%s/network/%s/public-ip/%s/static-nat"

	// NetworkVpnShowURL is the endpoint for showing VPN details for a network.
	NetworkVpnShowURL    string = "%s/zone/%s/network/%s/vpn"
	// NetworkVpnEnableURL is the endpoint for enabling VPN for a network.
	NetworkVpnEnableURL  string = "%s/zone/%s/network/%s/vpn/enable"
	// NetworkVpnDisableURL is the endpoint for disabling VPN for a network.
	NetworkVpnDisableURL string = "%s/zone/%s/network/%s/vpn/disable"
	// NetworkVpnUpdateURL is the endpoint for updating VPN credentials for a network.
	NetworkVpnUpdateURL  string = "%s/zone/%s/network/%s/vpn"

	// NetworkLoadBalancerList is the endpoint for listing all load balancer rules for a network.
	NetworkLoadBalancerList         string = "%s/zone/%s/network/%s/load-balancer"
	// NetworkLoadBalancerRuleCreate is the endpoint for creating a new load balancer rule.
	NetworkLoadBalancerRuleCreate   string = "%s/zone/%s/network/%s/load-balancer/rule"
	// NetworkLoadBalancerRuleDelete is the endpoint for deleting a load balancer rule.
	NetworkLoadBalancerRuleDelete   string = "%s/zone/%s/network/%s/load-balancer/rule/%s"
	// NetworkLoadBalancerRuleAssign is the endpoint for assigning instances to a load balancer rule.
	NetworkLoadBalancerRuleAssign   string = "%s/zone/%s/network/%s/load-balancer/rule/%s/assign"
	// NetworkLoadBalancerRuleDeassign is the endpoint for de-assigning an instance from a load balancer rule.
	NetworkLoadBalancerRuleDeassign string = "%s/zone/%s/network/%s/load-balancer/rule/%s/de-assign"
	// NetworkHaproxyLive is the endpoint for getting live HAProxy statistics.
	NetworkHaproxyLive              string = "%s/zone/%s/network/%s/ha/live"
	// NetworkHaproxyLog is the endpoint for getting HAProxy logs.
	NetworkHaproxyLog               string = "%s/zone/%s/network/%s/ha/log"

	// NetworkServiceOfferingList is the endpoint for listing all available network service offerings in a zone.
	NetworkServiceOfferingList string = "%s/zone/%s/network/service-offering"

	// InstanceList is the endpoint for listing all instances in a zone.
	InstanceList                string = "%s/zone/%s/instance"
	// InstanceServiceOfferingList is the endpoint for listing all available instance service offerings in a zone.
	InstanceServiceOfferingList string = "%s/zone/%s/instance/service-offerings"
	// InstanceVMImageList is the endpoint for listing all available VM images in a zone.
	InstanceVMImageList         string = "%s/zone/%s/instance/vm-images"
	// InstanceCreate is the endpoint for creating a new instance.
	InstanceCreate              string = "%s/zone/%s/instance"
	// InstanceRebuild is the endpoint for rebuilding an instance.
	InstanceRebuild             string = "%s/zone/%s/instance/%s/rebuild"
	// InstanceStart is the endpoint for starting an instance.
	InstanceStart               string = "%s/zone/%s/instance/%s/start"
	// InstanceStop is the endpoint for stopping an instance.
	InstanceStop                string = "%s/zone/%s/instance/%s/stop"
	// InstanceReboot is the endpoint for rebooting an instance.
	InstanceReboot              string = "%s/zone/%s/instance/%s/reboot"
	// InstanceDelete is the endpoint for deleting an instance.
	InstanceDelete              string = "%s/zone/%s/instance/%s"
	// InstanceShow is the endpoint for showing an instance.
	InstanceShow                string = "%s/zone/%s/instance/%s"

	// InstanceMetricsURL is the endpoint for getting instance metrics.
	InstanceMetricsURL                   string = "%s/zone/%s/instance/%s/metrics"
	// InstanceSnapshotCreateURL is the endpoint for creating a new snapshot.
	InstanceSnapshotCreateURL            string = "%s/zone/%s/instance/%s/snapshot"
	// InstanceSnapshotDeleteURL is the endpoint for deleting a snapshot.
	InstanceSnapshotDeleteURL            string = "%s/zone/%s/instance/%s/snapshot/%s"
	// InstanceSnapshotRevertURL is the endpoint for reverting an instance to a snapshot.
	InstanceSnapshotRevertURL            string = "%s/zone/%s/instance/%s/snapshot/%s/revert"
	// InstanceVolumeServiceOfferingListURL is the endpoint for listing all available instance volume service offerings in a zone.
	InstanceVolumeServiceOfferingListURL string = "%s/zone/%s/instance/volumes/service-offering"
	// InstanceVolumeListURL is the endpoint for listing all volumes in a zone.
	InstanceVolumeListURL                string = "%s/zone/%s/instance/volumes"
	// InstanceVolumeCreateURL is the endpoint for creating a new volume.
	InstanceVolumeCreateURL              string = "%s/zone/%s/instance/volumes"
	// InstanceVolumeDeleteURL is the endpoint for deleting a volume.
	InstanceVolumeDeleteURL              string = "%s/zone/%s/instance/volumes/%s"
	// InstanceVolumeDetachURL is the endpoint for detaching a volume from an instance.
	InstanceVolumeDetachURL              string = "%s/zone/%s/instance/volumes/%s/detach/%s"
	// InstanceVolumeAttachURL is the endpoint for attaching a volume to an instance.
	InstanceVolumeAttachURL              string = "%s/zone/%s/instance/volumes/%s/attach/%s"
)
