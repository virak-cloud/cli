package http

import (
	"fmt"
	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"net/http"
)

// GetZoneList fetches the list of zones and returns a ZoneListResponse struct.
func (client *Client) GetZoneList() (*responses.DataCenter, error) {
	var result responses.DataCenter
	url := fmt.Sprintf(urls.ZoneList, urls.BaseUrl)

	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetZoneActiveServices fetches the active services for a specific zone.
func (client *Client) GetZoneActiveServices(zoneID string) (*responses.ZoneActiveServicesResponse, error) {
	var result responses.ZoneActiveServicesResponse
	url := fmt.Sprintf(urls.ZoneActiveServicesList, urls.BaseUrl, zoneID)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetZoneCustomerResource fetches the customer resource for a specific zone.
func (client *Client) GetZoneCustomerResource(zoneID string) (*responses.CustomerResourceResponse, error) {
	var result responses.CustomerResourceResponse
	url := fmt.Sprintf(urls.ZoneResourcesList, urls.BaseUrl, zoneID)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetZoneNetworks fetches the networks for a specific zone.
func (client *Client) GetZoneNetworks(zoneID string) (*responses.ZoneNetworksResponse, error) {
	var result responses.ZoneNetworksResponse
	url := fmt.Sprintf(urls.ZoneNetworkList, urls.BaseUrl, zoneID)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
