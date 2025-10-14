package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http/responses"
)

func (client *Client) CreateL3Network(zoneId, networkOfferingId, name, gateway, netmask string) (*responses.NetworkCreateResponse, error) {
	var result responses.NetworkCreateResponse
	url := fmt.Sprintf(urls.NetworkCreateL3, urls.BaseUrl, zoneId)
	body, err := json.Marshal(map[string]string{
		"network_offering_id": networkOfferingId,
		"name":                name,
		"gateway":             gateway,
		"netmask":             netmask,
	})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) CreateL2Network(zoneId, networkOfferingId, name string) (*responses.NetworkCreateResponse, error) {
	var result responses.NetworkCreateResponse
	url := fmt.Sprintf(urls.NetworkCreateL2, urls.BaseUrl, zoneId)
	body, err := json.Marshal(map[string]string{
		"network_offering_id": networkOfferingId,
		"name":                name,
	})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) ListNetworks(zoneId string) (*responses.NetworkListResponse, error) {
	var result responses.NetworkListResponse
	url := fmt.Sprintf(urls.NetworkList, urls.BaseUrl, zoneId)

	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) ShowNetwork(zoneId, networkId string) (*responses.NetworkShowResponse, error) {
	var result responses.NetworkShowResponse
	url := fmt.Sprintf(urls.NetworkShow, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) DeleteNetwork(zoneId, networkId string) (*responses.NetworkDeleteResponse, error) {
	var result responses.NetworkDeleteResponse
	url := fmt.Sprintf(urls.NetworkDelete, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodDelete, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) ConnectInstanceToNetwork(zoneId, networkId, instanceId string) (*responses.InstanceNetworkActionResponse, error) {
	var result responses.InstanceNetworkActionResponse
	url := fmt.Sprintf(urls.NetworkInstanceConnect, urls.BaseUrl, zoneId, networkId)
	body, err := json.Marshal(map[string]string{
		"instance_id": instanceId,
	})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) DisconnectInstanceFromNetwork(zoneId, networkId, instanceId, instanceNetworkId string) (*responses.InstanceNetworkActionResponse, error) {
	var result responses.InstanceNetworkActionResponse
	url := fmt.Sprintf(urls.NetworkInstanceDisconnect, urls.BaseUrl, zoneId, networkId)
	body, err := json.Marshal(map[string]string{
		"instance_id":         instanceId,
		"instance_network_id": instanceNetworkId,
	})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) ListNetworkInstances(zoneId, networkId string, instanceId string) (*responses.InstanceNetworkListResponse, error) {
	var result responses.InstanceNetworkListResponse
	url := fmt.Sprintf(urls.NetworkInstanceList, urls.BaseUrl, zoneId, networkId)
	body, err := json.Marshal(map[string]string{
		"instance_id": instanceId,
	})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodGet, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// List IPv4 Firewall Rules
func (client *Client) ListIPv4FirewallRules(zoneId, networkId string) (*responses.IPv4FirewallRuleListResponse, error) {
	var result responses.IPv4FirewallRuleListResponse
	url := fmt.Sprintf(urls.NetworkFirewallIPv4List, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create IPv4 Firewall Rule
func (client *Client) CreateIPv4FirewallRule(zoneId, networkId string, body map[string]interface{}) (*responses.IPv4FirewallRuleActionResponse, error) {
	var result responses.IPv4FirewallRuleActionResponse
	url := fmt.Sprintf(urls.NetworkFirewallIPv4Create, urls.BaseUrl, zoneId, networkId)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete IPv4 Firewall Rule
func (client *Client) DeleteIPv4FirewallRule(zoneId, networkId, ruleId string) (*responses.IPv4FirewallRuleActionResponse, error) {
	var result responses.IPv4FirewallRuleActionResponse
	url := fmt.Sprintf(urls.NetworkFirewallIPv4Delete, urls.BaseUrl, zoneId, networkId, ruleId)
	if err := client.handleRequest(http.MethodDelete, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// List IPv6 Firewall Rules
func (client *Client) ListIPv6FirewallRules(zoneId, networkId string) (*responses.IPv6FirewallRuleListResponse, error) {
	var result responses.IPv6FirewallRuleListResponse
	url := fmt.Sprintf(urls.NetworkFirewallIPv6List, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create IPv6 Firewall Rule
func (client *Client) CreateIPv6FirewallRule(zoneId, networkId string, body map[string]interface{}) (*responses.IPv6FirewallRuleActionResponse, error) {
	var result responses.IPv6FirewallRuleActionResponse
	url := fmt.Sprintf(urls.NetworkFirewallIPv6Create, urls.BaseUrl, zoneId, networkId)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete IPv6 Firewall Rule
func (client *Client) DeleteIPv6FirewallRule(zoneId, networkId, ruleId string) (*responses.IPv6FirewallRuleActionResponse, error) {
	var result responses.IPv6FirewallRuleActionResponse
	url := fmt.Sprintf(urls.NetworkFirewallIPv6Delete, urls.BaseUrl, zoneId, networkId, ruleId)
	if err := client.handleRequest(http.MethodDelete, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Public IP: List
func (client *Client) ListNetworkPublicIps(zoneId, networkId string) (*responses.NetworkPublicIpListResponse, error) {
	var result responses.NetworkPublicIpListResponse
	url := fmt.Sprintf(urls.NetworkPublicIpList, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Public IP: Associate
func (client *Client) AssociateNetworkPublicIp(zoneId, networkId string) (*responses.NetworkPublicIpActionResponse, error) {
	var result responses.NetworkPublicIpActionResponse
	url := fmt.Sprintf(urls.NetworkPublicIpAssociate, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodPost, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Public IP: Disassociate
func (client *Client) DisassociateNetworkPublicIp(zoneId, networkId, networkPublicIpId string) (*responses.NetworkPublicIpActionResponse, error) {
	var result responses.NetworkPublicIpActionResponse
	url := fmt.Sprintf(urls.NetworkPublicIpDisassociate, urls.BaseUrl, zoneId, networkId, networkPublicIpId)
	if err := client.handleRequest(http.MethodDelete, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Public IP: Enable Static NAT
func (client *Client) EnableNetworkPublicIpStaticNat(zoneId, networkId, networkPublicIpId, instanceId string) (*responses.NetworkPublicIpActionResponse, error) {
	var result responses.NetworkPublicIpActionResponse
	url := fmt.Sprintf(urls.NetworkPublicIpStaticNatEnable, urls.BaseUrl, zoneId, networkId, networkPublicIpId)
	body, err := json.Marshal(map[string]string{"instance_id": instanceId})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Public IP: Disable Static NAT
func (client *Client) DisableNetworkPublicIpStaticNat(zoneId, networkId, networkPublicIpId string) (*responses.NetworkPublicIpActionResponse, error) {
	var result responses.NetworkPublicIpActionResponse
	url := fmt.Sprintf(urls.NetworkPublicIpStaticNatDisable, urls.BaseUrl, zoneId, networkId, networkPublicIpId)
	if err := client.handleRequest(http.MethodDelete, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetNetworkVpnDetails(zoneId, networkId string) (*responses.NetworkVpnDetailResponse, error) {
	var result responses.NetworkVpnDetailResponse
	url := fmt.Sprintf(urls.NetworkVpnShowURL, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) EnableNetworkVpn(zoneId, networkId string) (*responses.NetworkVpnSuccessResponse, error) {
	var result responses.NetworkVpnSuccessResponse
	url := fmt.Sprintf(urls.NetworkVpnEnableURL, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodPost, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) DisableNetworkVpn(zoneId, networkId string) (*responses.NetworkVpnSuccessResponse, error) {
	var result responses.NetworkVpnSuccessResponse
	url := fmt.Sprintf(urls.NetworkVpnDisableURL, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodPost, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) UpdateNetworkVpnCredentials(zoneId, networkId string) (*responses.NetworkVpnSuccessResponse, error) {
	var result responses.NetworkVpnSuccessResponse
	url := fmt.Sprintf(urls.NetworkVpnUpdateURL, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodPut, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) ListLoadBalancerRules(zoneId, networkId string) (*responses.LoadBalancerRuleListResponse, error) {
	var result responses.LoadBalancerRuleListResponse
	url := fmt.Sprintf(urls.NetworkLoadBalancerList, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) CreateLoadBalancerRule(zoneId, networkId, publicIpId, name, algorithm string, publicPort, privatePort int) (*responses.SuccessResponse, error) {
	var result responses.SuccessResponse
	url := fmt.Sprintf(urls.NetworkLoadBalancerRuleCreate, urls.BaseUrl, zoneId, networkId)
	body, err := json.Marshal(map[string]interface{}{
		"public_ip_id": publicIpId,
		"name":         name,
		"algorithm":    algorithm,
		"public_port":  publicPort,
		"private_port": privatePort,
	})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) DeleteLoadBalancerRule(zoneId, networkId, ruleId string) (*responses.SuccessResponse, error) {
	var result responses.SuccessResponse
	url := fmt.Sprintf(urls.NetworkLoadBalancerRuleDelete, urls.BaseUrl, zoneId, networkId, ruleId)
	if err := client.handleRequest(http.MethodDelete, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) AssignLoadBalancerRule(zoneId, networkId, ruleId string, instanceNetworkIds []string) (*responses.SuccessResponse, error) {
	var result responses.SuccessResponse
	url := fmt.Sprintf(urls.NetworkLoadBalancerRuleAssign, urls.BaseUrl, zoneId, networkId, ruleId)
	body, err := json.Marshal(map[string]interface{}{
		"instance_network_ids": instanceNetworkIds,
	})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) DeassignLoadBalancerRule(zoneId, networkId, ruleId, instanceNetworkId string) (*responses.SuccessResponse, error) {
	var result responses.SuccessResponse
	url := fmt.Sprintf(urls.NetworkLoadBalancerRuleDeassign, urls.BaseUrl, zoneId, networkId, ruleId)
	body, err := json.Marshal(map[string]interface{}{
		"instance_network_id": instanceNetworkId,
	})
	if err != nil {
		return nil, err
	}
	if err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetHaproxyLive(zoneId, networkId string) (*responses.HaproxyLiveResponse, error) {
	var result responses.HaproxyLiveResponse
	url := fmt.Sprintf(urls.NetworkHaproxyLive, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetHaproxyLog(zoneId, networkId string) (*responses.HaproxyLogResponse, error) {
	var result responses.HaproxyLogResponse
	url := fmt.Sprintf(urls.NetworkHaproxyLog, urls.BaseUrl, zoneId, networkId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) ListNetworkServiceOfferings(zoneId string) (*responses.NetworkServiceOfferingListResponse, error) {
	var result responses.NetworkServiceOfferingListResponse
	url := fmt.Sprintf(urls.NetworkServiceOfferingList, urls.BaseUrl, zoneId)
	if err := client.handleRequest(http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetL2NetworkServiceOfferings(zoneId string) (*responses.NetworkServiceOfferingListResponse, error) {
	allOfferings, err := client.ListNetworkServiceOfferings(zoneId)
	if err != nil {
		return nil, err
	}

	var l2Offerings []responses.NetworkOffering
	for _, offering := range allOfferings.Data {
		if offering.Type == "L2" {
			l2Offerings = append(l2Offerings, offering)
		}
	}

	return &responses.NetworkServiceOfferingListResponse{
		Data: l2Offerings,
	}, nil
}

func (client *Client) GetL3NetworkServiceOfferings(zoneId string) (*responses.NetworkServiceOfferingListResponse, error) {
	allOfferings, err := client.ListNetworkServiceOfferings(zoneId)
	if err != nil {
		return nil, err
	}

	var l3Offerings []responses.NetworkOffering
	for _, offering := range allOfferings.Data {
		if offering.Type == "Isolated" {
			l3Offerings = append(l3Offerings, offering)
		}
	}

	return &responses.NetworkServiceOfferingListResponse{
		Data: l3Offerings,
	}, nil
}
