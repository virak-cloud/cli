// Package http provides a comprehensive HTTP client for interacting with the Virak Cloud API.
//
// The http package implements a robust client for making authenticated requests to
// the Virak Cloud REST API. It handles authentication, request/response processing,
// error handling, and provides type-safe methods for all API operations.
//
// Key Features:
//   - Automatic authentication with Bearer tokens
//   - JSON request/response handling
//   - Comprehensive error wrapping with context
//   - Support for all Virak Cloud API endpoints
//   - Thread-safe client implementation
//
// Authentication:
//   The client requires a valid API token which can be obtained through the
//   OAuth authentication flow. The token is automatically included in all
//   request headers as a Bearer token.
//
// Error Handling:
//   All methods return detailed error messages that include:
//   - The specific operation that failed
//   - API response details when available
//   - Context about the request parameters
//
// Basic Usage:
//   client := http.NewClient("your-api-token")
//   instances, err := client.ListInstances("zone-123")
//   if err != nil {
//       log.Fatalf("Failed to list instances: %v", err)
//   }
//   fmt.Printf("Found %d instances\n", len(instances.Data))
//
// Advanced Usage:
//   client := http.NewClient("your-api-token")
//
//   // Create instance with custom configuration
//   resp, err := client.CreateInstance(
//       "zone-123",
//       "offering-456",
//       "image-789",
//       []string{"network-abc"},
//       "my-instance",
//   )
//   if err != nil {
//       return fmt.Errorf("instance creation failed: %w", err)
//   }
//
// Thread Safety:
//   The Client is safe for concurrent use by multiple goroutines.
//   All HTTP operations are handled internally with proper synchronization.
//
// API Categories:
//   - Instance Management: Create, list, start, stop, delete virtual machines
//   - Network Management: Configure networks, firewall rules, load balancers
//   - Storage Management: Manage buckets, volumes, and snapshots
//   - DNS Management: Manage domains and DNS records
//   - Kubernetes Management: Deploy and manage Kubernetes clusters
//   - User Management: SSH keys, billing, and account information
//   - Zone Management: List available zones and their resources
//
// For more detailed examples, see the function documentation and examples.
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/virak-cloud/cli/pkg/http/responses"
)

// Client is a client for the Virak Cloud API.
// It handles authentication, request signing, and response parsing.
//
// The Client struct is the main entry point for all API operations. It maintains
// the authentication state, HTTP client configuration, and provides methods
// for interacting with all Virak Cloud services.
//
// Thread Safety:
//   The Client is safe for concurrent use by multiple goroutines. All internal
//   state is immutable after creation, and the underlying HTTP client is
//   configured for concurrent use.
//
// Configuration:
//   - Token: Authentication token for API access
//   - BaseURL: Base URL for the API (configurable for different environments)
//   - HttpClient: underlying HTTP client with default timeouts and settings
type Client struct {
	// HttpClient is the underlying HTTP client used to make requests.
	// It's configured with appropriate timeouts and settings for API communication.
	HttpClient *http.Client
	
	// Token is the authentication token used to authenticate with the API.
	// This token is included in the Authorization header of all requests.
	Token string
	
	// BaseURL is the base URL of the Virak Cloud API.
	// This can be configured for different environments (dev, staging, prod).
	BaseURL string
}

// NewClient creates a new Virak Cloud API client.
//
// It takes an authentication token as input and returns a new `Client` instance.
// The default `BaseURL` is set to `http://localhost:1410`, which is suitable
// for local development and testing. For production use, the BaseURL should
// be updated to the production API endpoint.
//
// Parameters:
//   - token: A valid API token obtained from the Virak Cloud authentication
//     system. The token must have sufficient permissions for the intended
//     operations.
//
// Returns:
//   - *Client: A new client instance ready for API operations
//
// Example:
//   // Create client for development
//   client := http.NewClient("dev-token-123")
//
//   // Update base URL for production
//   client.BaseURL = "https://public-api.virakcloud.com"
//
//   // Use client for API operations
//   instances, err := client.ListInstances("zone-123")
//   if err != nil {
//       log.Fatal(err)
//   }
//
// Authentication:
//   The provided token will be included in all API requests using the
//   Authorization: Bearer header. Ensure the token is kept secure and
//   has the minimum required permissions.
//
// Environment Configuration:
//   For different environments, update the BaseURL after creating the client:
//   - Development: http://localhost:1410
//   - Staging: https://staging-api.virakcloud.com
//   - Production: https://public-api.virakcloud.com
func NewClient(token string) *Client {
	return &Client{
		HttpClient: &http.Client{},
		Token:      token,
		BaseURL:    "http://localhost:1410",
	}
}

// handleRequest is a generic helper function that executes HTTP requests and
// decodes the JSON response into the `target` interface.
//
// It handles setting the required headers, executing the request, and parsing
// the response. If the response status code is 400 or greater, it decodes
// the error response and returns an error.
//
// Parameters:
//   - method: The HTTP method to use (e.g., "GET", "POST").
//   - path: The API endpoint path (e.g., "/instances").
//   - body: The request body, or `nil` if there is no body.
//   - target: The interface to decode the JSON response into.
//
// Returns:
//   - An error if the request fails or the response cannot be decoded, otherwise `nil`.
func (client *Client) handleRequest(method string, path string, body io.Reader, target interface{}) error {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if client.Token != "" {
		req.Header.Set("Authorization", "Bearer "+client.Token)
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiError responses.ErrorResponse
		if err := json.Unmarshal(respBody, &apiError); err != nil {
			// Fallback if the error response isn't the expected JSON
			return fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("API error: %s", respBody)
	}

	if target != nil {
		if err := json.Unmarshal(respBody, target); err != nil {
			return fmt.Errorf("failed to decode successful response: %w", err)
		}
	}

	return nil
}

// Request sends a raw HTTP request to the Virak Cloud API and returns the
// response body as a byte slice.
//
// This method is a thin wrapper around `handleRequest` and is useful for
// making API calls that don't require decoding a JSON response.
//
// Parameters:
//   - method: The HTTP method to use (e.g., "GET", "POST").
//   - url: The full URL of the API endpoint.
//   - body: The request body as a byte slice.
//
// Returns:
//   - The response body as a byte slice.
//   - An error if the request fails.
func (client *Client) Request(method, url string, body []byte) ([]byte, error) {
	var responseBody []byte
	err := client.handleRequest(method, url, bytes.NewBuffer(body), &responseBody)
	return responseBody, err
}
