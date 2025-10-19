// Package http provides zone and datacenter management operations for the Virak Cloud API.
//
// This file contains client methods for retrieving zone information, including
// available datacenters, their locations, capabilities, active services,
// resource availability, and network configurations. Zone information is
// essential for planning resource deployment and understanding service availability.
package http

import (
	"fmt"
	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"net/http"
)

// GetZoneList retrieves a list of all available zones (datacenters).
//
// It sends a GET request to the `/zones` endpoint and returns a `DataCenter`
// struct containing a list of zones.
//
// Returns:
//   - `*responses.DataCenter`: A struct containing the list of zones.
//   - `error`: An error if the request fails or the response cannot be decoded.
//
// Example:
//
//	client := http.NewClient("your-auth-token")
//	zones, err := client.GetZoneList()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, zone := range zones.Data.Zones {
//	    fmt.Printf("Zone: %s\n", zone.Name)
//	}
func (client *Client) GetZoneList() (*responses.DataCenter, error) {
	var result responses.DataCenter
	url := fmt.Sprintf(urls.ZoneList, urls.BaseUrl)

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetZoneActiveServices retrieves a list of active services for a specific zone.
//
// It sends a GET request to the `/zones/{zoneID}/services` endpoint and returns
// a `ZoneActiveServicesResponse` struct containing a list of active services.
//
// Parameters:
//   - `zoneID`: The ID of the zone to retrieve active services for.
//
// Returns:
//   - `*responses.ZoneActiveServicesResponse`: A struct containing the list of active services.
//   - `error`: An error if the request fails or the response cannot be decoded.
//
// Example:
//
//	client := http.NewClient("your-auth-token")
//	services, err := client.GetZoneActiveServices("zone-123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, service := range services.Data.Services {
//	    fmt.Printf("Service: %s\n", service.Name)
//	}
func (client *Client) GetZoneActiveServices(zoneID string) (*responses.ZoneActiveServicesResponse, error) {
	var result responses.ZoneActiveServicesResponse
	url := fmt.Sprintf(urls.ZoneActiveServicesList, urls.BaseUrl, zoneID)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetZoneCustomerResource retrieves a list of customer resources for a specific zone.
//
// It sends a GET request to the `/zones/{zoneID}/resources` endpoint and returns
// a `CustomerResourceResponse` struct containing a list of customer resources.
//
// Parameters:
//   - `zoneID`: The ID of the zone to retrieve customer resources for.
//
// Returns:
//   - `*responses.CustomerResourceResponse`: A struct containing the list of customer resources.
//   - `error`: An error if the request fails or the response cannot be decoded.
//
// Example:
//
//	client := http.NewClient("your-auth-token")
//	resources, err := client.GetZoneCustomerResource("zone-123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, resource := range resources.Data.Resources {
//	    fmt.Printf("Resource: %s\n", resource.Name)
//	}
func (client *Client) GetZoneCustomerResource(zoneID string) (*responses.CustomerResourceResponse, error) {
	var result responses.CustomerResourceResponse
	url := fmt.Sprintf(urls.ZoneResourcesList, urls.BaseUrl, zoneID)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetZoneNetworks retrieves a list of networks for a specific zone.
//
// It sends a GET request to the `/zones/{zoneID}/networks` endpoint and returns
// a `ZoneNetworksResponse` struct containing a list of networks.
//
// Parameters:
//   - `zoneID`: The ID of the zone to retrieve networks for.
//
// Returns:
//   - `*responses.ZoneNetworksResponse`: A struct containing the list of networks.
//   - `error`: An error if the request fails or the response cannot be decoded.
//
// Example:
//
//	client := http.NewClient("your-auth-token")
//	networks, err := client.GetZoneNetworks("zone-123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, network := range networks.Data.Networks {
//	    fmt.Printf("Network: %s\n", network.Name)
//	}
func (client *Client) GetZoneNetworks(zoneID string) (*responses.ZoneNetworksResponse, error) {
	var result responses.ZoneNetworksResponse
	url := fmt.Sprintf(urls.ZoneNetworkList, urls.BaseUrl, zoneID)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
