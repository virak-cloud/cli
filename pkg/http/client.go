package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/virak-cloud/cli/pkg/http/responses"
)

type Client struct {
	HttpClient *http.Client
	Token      string
	BaseURL    string
}

func NewClient(token string) *Client {
	return &Client{
		HttpClient: &http.Client{},
		Token:      token,
		BaseURL:    "http://localhost:1410",
	}
}

// handleRequest is a generic helper to execute HTTP requests and decode responses.
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

func (client *Client) Request(method, url string, body []byte) ([]byte, error) {
	var responseBody []byte
	err := client.handleRequest(method, url, bytes.NewBuffer(body), &responseBody)
	return responseBody, err
}
