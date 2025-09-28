package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"net/http"
)

// User-related API methods for the Client

// ValidateUserToken checks if the user's token is valid (204 No Content on success).
func (client *Client) ValidateUserToken() error {
	url := fmt.Sprintf(urls.UserTokenValidate, urls.BaseUrl)
	_, err := client.Request("GET", url, nil)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}
	return nil

}

// GetUserTokenAbilities fetches the abilities associated with the user's token.
func (client *Client) GetUserTokenAbilities() (*responses.UserTokenAbilitiesResponse, error) {

	var result responses.UserTokenAbilitiesResponse
	url := fmt.Sprintf(urls.UserTokenAbilities, urls.BaseUrl)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

// SSH Key Management

// ListUserSSHKeys fetches the list of user's SSH keys.
func (client *Client) ListUserSSHKeys() (*responses.UserSSHKeyListResponse, error) {
	var result responses.UserSSHKeyListResponse
	url := fmt.Sprintf(urls.UserSSHKeyList, urls.BaseUrl)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddUserSSHKey adds a new SSH key for the user. Pass body as JSON: {"name":"My SSH Key","ssh_key":"ssh-rsa ..."}
func (client *Client) AddUserSSHKey(sshKeyName string, sshkey string) (*responses.AddUserSSHKeyResponse, error) {

	var result responses.AddUserSSHKeyResponse
	url := fmt.Sprintf(urls.UserSSHKeyCreate, urls.BaseUrl)
	body, err := json.Marshal(map[string]string{"name": sshKeyName, "ssh_key": sshkey})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ssh: %w", err)
	}

	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

// DeleteUserSSHKey deletes a user's SSH key by its ID.
func (client *Client) DeleteUserSSHKey(sshKeyId string) (*responses.DeleteUserSSHKeyResponse, error) {
	var result responses.DeleteUserSSHKeyResponse
	url := fmt.Sprintf(urls.UserSSHKeyDelete, urls.BaseUrl, sshKeyId)
	err := client.handleRequest(http.MethodDelete, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

// Finance

// GetWalletsBalance fetches the user's wallet balance.
func (client *Client) GetWalletsBalance() (*responses.WalletsBalanceResponse, error) {
	var result responses.WalletsBalanceResponse
	url := fmt.Sprintf(urls.UserBalance, urls.BaseUrl)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

// GetCostDocumentsYearly fetches yearly cost documents. Pass body as JSON: {"year":1402}
func (client *Client) GetCostDocumentsYearly(year int) (*responses.CostDocumentsYearlyResponse, error) {
	var result responses.CostDocumentsYearlyResponse
	body, err := json.Marshal(map[string]int{"year": year})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal year: %w", err)
	}
	url := fmt.Sprintf(urls.UserCostDocumentList, urls.BaseUrl)
	err = client.handleRequest(http.MethodGet, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPaymentList fetches the user's payment list.
func (client *Client) GetPaymentList() (*responses.PaymentListResponse, error) {
	var result responses.PaymentListResponse
	url := fmt.Sprintf(urls.UserPaymentList, urls.BaseUrl)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetTokenAbilities() (*responses.UserTokenAbilitiesResponse, error) {

	var result responses.UserTokenAbilitiesResponse
	url := fmt.Sprintf(urls.UserTokenAbilities, urls.BaseUrl)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
