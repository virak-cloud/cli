package responses

import "encoding/json"

type NetworkCreateResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type NetworkListResponse struct {
	Data []Network `json:"data"`
}

type Network struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Status          string            `json:"status"`
	NetworkOffering NetworkOffering   `json:"network_offering"`
	InstanceNetwork []InstanceNetwork `json:"instance_network"`
}

type IPConfig struct {
	Gateway string `json:"gateway"`
	Netmask string `json:"netmask"`
}

type IPConfigV6 struct {
	GatewayV6      *string `json:"gateway_v6"`
	NetmaskV6      *string `json:"netmask_v6"`
	RouteGatewayV6 *string `json:"route_gateway_v6"`
	RouteSubnetV6  *string `json:"route_subnet_v6"`
}

type IPConfigOrArray IPConfig

func (i *IPConfigOrArray) UnmarshalJSON(data []byte) error {
	if string(data) == "[]" {
		*i = IPConfigOrArray{}
		return nil
	}
	return json.Unmarshal(data, (*IPConfig)(i))
}

type IPConfigV6OrArray IPConfigV6

func (i *IPConfigV6OrArray) UnmarshalJSON(data []byte) error {
	if string(data) == "[]" {
		*i = IPConfigV6OrArray{}
		return nil
	}
	return json.Unmarshal(data, (*IPConfigV6)(i))
}

type NetworkOffering struct {
	ID                       string  `json:"id"`
	Name                     string  `json:"name"`
	DisplayName              string  `json:"displayname"`
	DisplayNameFA            string  `json:"displayname_fa"`
	HourlyStartedPrice       float64 `json:"hourly_started_price"`
	TrafficTransferOverprice float64 `json:"traffic_transfer_overprice"`
	TrafficTransferPlan      int     `json:"traffic_transfer_plan"`
	NetworkRate              int     `json:"networkrate"`
	Type                     string  `json:"type"`
	Description              string  `json:"description"`
	InternetProtocol         string  `json:"internet_protocol"`
}

type NetworkShowResponse struct {
	Data Network `json:"data"`
}

type NetworkDeleteResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type InstanceNetworkActionResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type InstanceNetworkListResponse struct {
	Data []InstanceNetwork `json:"data"`
}

type InstanceNetwork struct {
	ID              string                 `json:"id"`
	InstanceID      string                 `json:"instance_id"`
	InstanceName    string                 `json:"instance_name"`
	IPAddress       string                 `json:"ipaddress"`
	IPAddressV6     string                 `json:"ipaddress_v6"`
	MACAddress      string                 `json:"macaddress"`
	IsDefault       bool                   `json:"is_default"`
	CreatedAt       int64                  `json:"created_at"`
	Network         NetworkSummary         `json:"network"`
	NetworkOffering NetworkOfferingSummary `json:"network_offering"`
	SecondaryIPs    []SecondaryIP          `json:"secondary_ips"`
}

type SecondaryIP struct {
	ID        string `json:"id"`
	IPAddress string `json:"ipaddress"`
	CreatedAt int64  `json:"created_at"`
}

type NetworkSummary struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	IPConfig   IPConfigOrArray   `json:"ip_config"`
	IPConfigV6 IPConfigV6OrArray `json:"ip_config_v6"`
}

type NetworkOfferingSummary struct {
	ID                     string  `json:"id"`
	Name                   string  `json:"name"`
	DisplayName            string  `json:"displayname"`
	DisplayNameFA          string  `json:"displayname_fa"`
	HourlyStartedPrice     float64 `json:"hourly_started_price"`
	TrafficPricePerGig     float64 `json:"traffic_price_per_gig"`
	TrafficTransferFreeGig int     `json:"traffic_transfer_free_gig"`
	NetworkRate            int     `json:"networkrate"`
	Type                   string  `json:"type"`
	Description            string  `json:"description"`
	InternetProtocol       string  `json:"internet_protocol"`
}

// IPv4 Firewall Rule
// Used for both list and create responses
// Example: see API docs

type IPv4FirewallRuleListResponse struct {
	Data []IPv4FirewallRule `json:"data"`
}

type IPv4FirewallRule struct {
	ID                string  `json:"id"`
	NetworkPublicIPID *string `json:"network_public_ip_id"`
	Protocol          string  `json:"protocol"`
	TrafficType       string  `json:"traffic_type"`
	IPSource          string  `json:"ip_source"`
	IPDestination     string  `json:"ip_destination"`
	PortStart         *string `json:"port_start"`
	PortEnd           *string `json:"port_end"`
	ICMPCode          *int    `json:"icmp_code"`
	ICMPType          *int    `json:"icmp_type"`
	Status            string  `json:"status"`
	CreatedAt         int64   `json:"created_at"`
}

type IPv4FirewallRuleActionResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type IPv6FirewallRuleListResponse struct {
	Data []IPv6FirewallRule `json:"data"`
}

type IPv6FirewallRule struct {
	ID            string  `json:"id"`
	Protocol      string  `json:"protocol"`
	TrafficType   string  `json:"traffic_type"`
	IPSource      string  `json:"ip_source"`
	IPDestination string  `json:"ip_destination"`
	PortStart     *string `json:"port_start"`
	PortEnd       *string `json:"port_end"`
	ICMPCode      *int    `json:"icmp_code"`
	ICMPType      *int    `json:"icmp_type"`
	Status        string  `json:"status"`
	CreatedAt     int64   `json:"created_at"`
}

type IPv6FirewallRuleActionResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

// Public IP responses

// List response
type NetworkPublicIpListResponse struct {
	Data []NetworkPublicIp `json:"data"`
}

// Public IP struct
type NetworkPublicIp struct {
	ID              string   `json:"id"`
	NetworkID       string   `json:"network_id"`
	IpAddress       string   `json:"ipaddress"`
	IsSourceNat     bool     `json:"is_sourcenat"`
	CreatedAt       int64    `json:"created_at"`
	StaticNatEnable bool     `json:"staticnat_enable"`
	StaticNat       []string `json:"staticnat"`
}

// Action response (associate, disassociate, static nat)
type NetworkPublicIpActionResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type NetworkVpnDetailResponse struct {
	Data struct {
		IPAddress    string `json:"ipaddress"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		PresharedKey string `json:"presharedkey"`
		Status       string `json:"status"`
	} `json:"data"`
}

type NetworkVpnSuccessResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type LoadBalancerRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Algorithm   string `json:"algorithm"`
	PublicPort  int    `json:"public_port"`
	PrivatePort int    `json:"private_port"`
	Status      string `json:"status"`
}

type LoadBalancerRuleListResponse struct {
	Data []LoadBalancerRule `json:"data"`
}

type SuccessResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}

type HaproxyLiveResponse struct {
	Data struct {
		UpdatedAt int64              `json:"updated_at"`
		Rules     []LoadBalancerRule `json:"rules"`
	} `json:"data"`
}

type HaproxyLogResponse struct {
	Data []interface{} `json:"data"`
}

type NetworkServiceOfferingListResponse struct {
	Data []NetworkOffering `json:"data"`
}

// Port Forwarding responses

type PortForwardRule struct {
	ID         string `json:"id"`
	NetworkID  string `json:"network_id"`
	Protocol   string `json:"protocol"`
	PublicPort int    `json:"public_port"`
	PrivatePort int   `json:"private_port"`
	PrivateIP  string `json:"private_ip"`
	Status     string `json:"status"`
	CreatedAt  int64  `json:"created_at"`
}

type PortForwardListResponse struct {
	Data []PortForwardRule `json:"data"`
}

type PortForwardActionResponse struct {
	Data struct {
		Success bool `json:"success"`
	} `json:"data"`
}
