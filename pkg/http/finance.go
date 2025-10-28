package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	urls "github.com/virak-cloud/cli/pkg"
	"github.com/virak-cloud/cli/pkg/http/responses"
	"net/http"
)

// Finance-related API methods for the Client

// GetWallet fetches the user's wallet balance.
func (client *Client) GetWallet() (*responses.WalletsBalanceResponse, error) {
	var result responses.WalletsBalanceResponse
	url := fmt.Sprintf(urls.UserBalance, urls.BaseUrl)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDocuments fetches cost documents for a specific year.
func (client *Client) ListDocuments(year int) (*responses.CostDocumentsYearlyResponse, error) {
	var result responses.CostDocumentsYearlyResponse
	// DEBUG: Log the current approach
	fmt.Printf("DEBUG: Current approach - Using POST method with year=%d\n", year)
	body, err := json.Marshal(map[string]int{"year": year})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal year: %w", err)
	}
	url := fmt.Sprintf(urls.UserCostDocumentList, urls.BaseUrl)
	fmt.Printf("DEBUG: URL: %s\n", url)
	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDocumentsGET fetches cost documents for a specific year using GET method.
func (client *Client) ListDocumentsGET(year int) (*responses.CostDocumentsYearlyResponse, error) {
	var result responses.CostDocumentsYearlyResponse
	
	url := fmt.Sprintf("%s/user/finance/documents?year=%d", urls.BaseUrl, year)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListPayments fetches the user's payment history.
func (client *Client) ListPayments() (*responses.PaymentListResponse, error) {
	var result responses.PaymentListResponse
	url := fmt.Sprintf(urls.UserPaymentList, urls.BaseUrl)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListExpenses fetches the user's expenses with optional filtering.
func (client *Client) ListExpenses(filters map[string]string) (*responses.ExpensesListResponse, error) {
	var result responses.ExpensesListResponse
	
	// DEBUG: Log the current approach
	fmt.Printf("DEBUG: Current approach - Using GET method with filters=%+v\n", filters)
	
	// Build query string from filters
	url := fmt.Sprintf("%s/user/finance/expenses", urls.BaseUrl)
	
	// Add filters as query parameters if any
	if len(filters) > 0 {
		queryParams := ""
		for key, value := range filters {
			if queryParams == "" {
				queryParams = "?"
			} else {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("%s=%s", key, value)
		}
		url += queryParams
	}
	
	fmt.Printf("DEBUG: URL: %s\n", url)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListExpensesWithRequiredParams fetches the user's expenses with required parameters.
func (client *Client) ListExpensesWithRequiredParams(productType, productID string, filters map[string]string) (*responses.ExpensesListResponse, error) {
	var result responses.ExpensesListResponse
	
	// DEBUG: Log the new approach with required params
	fmt.Printf("DEBUG: New approach - Using GET method with product_type=%s, product_id=%s, filters=%+v\n", productType, productID, filters)
	
	// Build query string with required parameters
	url := fmt.Sprintf("%s/user/finance/expenses?product_type=%s&product_id=%s", urls.BaseUrl, productType, productID)
	
	// Add additional filters as query parameters if any
	if len(filters) > 0 {
		for key, value := range filters {
			url += fmt.Sprintf("&%s=%s", key, value)
		}
	}
	
	fmt.Printf("DEBUG: URL: %s\n", url)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}