// Package http provides DNS management operations for the Virak Cloud API.
//
// This file contains client methods for managing DNS domains and records, including
// domain creation/deletion, DNS record management (A, CNAME, MX, SRV, CAA, TLSA),
// and DNS event monitoring. These operations allow users to manage their DNS
// infrastructure through the Virak Cloud platform.
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	s "strings"

	urls "github.com/virak-cloud/cli/pkg"

	"github.com/virak-cloud/cli/pkg/http/responses"
)

// GetDomains retrieves a list of all DNS domains.
//
// It sends a GET request to the `/dns` endpoint and returns a `DomainList`
// struct containing a list of domains.
//
// Returns:
//   - `*responses.DomainList`: A struct containing the list of domains.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) GetDomains() (*responses.DomainList, error) {
	url := fmt.Sprintf(urls.DomainListURL, urls.BaseUrl)
	var domainList responses.DomainList
	err := client.handleRequest(http.MethodGet, url, nil, &domainList)
	if err != nil {
		return nil, err
	}
	return &domainList, nil
}

// CreateDomain creates a new DNS domain.
//
// It sends a POST request to the `/dns` endpoint with the domain name.
//
// Parameters:
//   - `domain`: The name of the domain to create.
//
// Returns:
//   - `*responses.DnsMessage`: A message confirming the creation of the domain.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) CreateDomain(domain string) (*responses.DnsMessage, error) {
	url := fmt.Sprintf(urls.DomainCreateURL, urls.BaseUrl)
	body := []byte(fmt.Sprintf(`{"domain": "%s"}`, domain))
	var message responses.DnsMessage
	err := client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetDomain retrieves a specific DNS domain.
//
// It sends a GET request to the `/dns/{domain}` endpoint and returns a
// `DomainShow` struct containing the details of the domain.
//
// Parameters:
//   - `domain`: The name of the domain to retrieve.
//
// Returns:
//   - `*responses.DomainShow`: A struct containing the details of the domain.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) GetDomain(domain string) (*responses.DomainShow, error) {
	url := fmt.Sprintf(urls.DomainShowURL, urls.BaseUrl, domain)
	var domainShow responses.DomainShow

	err := client.handleRequest(http.MethodGet, url, nil, &domainShow)

	if err != nil {
		if s.Contains(err.Error(), "cannot unmarshal array into Go struct field DomainShow.data of type responses.Domain") {
			return &domainShow, nil
		}

		return nil, err
	}

	return &domainShow, nil
}

// DeleteDomain deletes a specific DNS domain.
//
// It sends a DELETE request to the `/dns/{domain}` endpoint.
//
// Parameters:
//   - `domain`: The name of the domain to delete.
//
// Returns:
//   - `*responses.DnsMessage`: A message confirming the deletion of the domain.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) DeleteDomain(domain string) (*responses.DnsMessage, error) {
	url := fmt.Sprintf(urls.DomainDeleteURL, urls.BaseUrl, domain)
	var message responses.DnsMessage
	err := client.handleRequest(http.MethodDelete, url, nil, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetRecords retrieves a list of all DNS records for a specific domain.
//
// It sends a GET request to the `/dns/{domain}/records` endpoint and returns
// a `RecordList` struct containing a list of records.
//
// Parameters:
//   - `domain`: The name of the domain to retrieve records for.
//
// Returns:
//   - `*responses.RecordList`: A struct containing the list of records.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) GetRecords(domain string) (*responses.RecordList, error) {
	url := fmt.Sprintf(urls.RecordListURL, urls.BaseUrl, domain)
	var recordList responses.RecordList
	err := client.handleRequest(http.MethodGet, url, nil, &recordList)
	if err != nil {
		return nil, err
	}
	return &recordList, nil
}

// CreateRecord creates a new DNS record for a specific domain.
//
// It sends a POST request to the `/dns/{domain}/records` endpoint with the
// details of the record to create.
//
// Parameters:
//   - `domain`: The name of the domain to create the record in.
//   - `record`: The name of the record.
//   - `recordType`: The type of the record (e.g., "A", "CNAME", "MX").
//   - `content`: The content of the record.
//   - `ttl`: The TTL of the record.
//   - `priority`: The priority of the record (for MX and SRV records).
//   - `weight`: The weight of the record (for SRV records).
//   - `port`: The port of the record (for SRV records).
//   - `flags`: The flags of the record (for CAA records).
//   - `tag`: The tag of the record (for CAA records).
//   - `license`: The license of the record (for TLSA records).
//   - `choicer`: The choicer of the record (for TLSA records).
//   - `match`: The match of the record (for TLSA records).
//
// Returns:
//   - `*responses.DnsMessage`: A message confirming the creation of the record.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) CreateRecord(domain, record, recordType, content string, ttl, priority, weight, port, flags int, tag string, license, choicer, match int) (*responses.DnsMessage, error) {
	url := fmt.Sprintf(urls.RecordCreateURL, urls.BaseUrl, domain)

	// Base fields that are always included
	bodyMap := map[string]interface{}{
		"record":  record,
		"type":    recordType,
		"ttl":     ttl,
		"content": content,
	}

	// Add type-specific fields
	switch recordType {
	case "MX", "SRV":
		bodyMap["priority"] = priority
		if recordType == "SRV" {
			bodyMap["weight"] = weight
			bodyMap["port"] = port
		}
	case "CAA":
		bodyMap["flags"] = flags
		bodyMap["tag"] = tag
	case "TLSA":
		bodyMap["license"] = license
		bodyMap["choicer"] = choicer
		bodyMap["match"] = match
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}
	var message responses.DnsMessage
	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// UpdateRecord updates a specific DNS record for a specific domain.
//
// It sends a PUT request to the `/dns/{domain}/records/{record}/{type}/{content}`
// endpoint with the updated details of the record.
//
// Parameters:
//   - `domain`: The name of the domain the record belongs to.
//   - `record`: The name of the record to update.
//   - `recordType`: The type of the record to update.
//   - `contentId`: The content of the record to update.
//   - `newContent`: The new content of the record.
//   - `newTTL`: The new TTL of the record.
//   - `priority`: The new priority of the record (for MX and SRV records).
//   - `weight`: The new weight of the record (for SRV records).
//   - `port`: The new port of the record (for SRV records).
//   - `flags`: The new flags of the record (for CAA records).
//   - `tag`: The new tag of the record (for CAA records).
//   - `license`: The new license of the record (for TLSA records).
//   - `choicer`: The new choicer of the record (for TLSA records).
//   - `match`: The new match of the record (for TLSA records).
//
// Returns:
//   - `*responses.DnsMessage`: A message confirming the update of the record.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) UpdateRecord(domain, record, recordType, contentId, newContent string, newTTL, priority, weight, port, flags int, tag string, license, choicer, match int) (*responses.DnsMessage, error) {
	url := fmt.Sprintf(urls.RecordUpdateURL, urls.BaseUrl, domain, record, recordType, contentId)

	// Base fields that are always included
	bodyMap := map[string]interface{}{
		"ttl":     newTTL,
		"content": newContent,
	}

	// Add type-specific fields
	switch recordType {
	case "MX", "SRV":
		bodyMap["priority"] = priority
		if recordType == "SRV" {
			bodyMap["weight"] = weight
			bodyMap["port"] = port
		}
	case "CAA":
		bodyMap["flags"] = flags
		bodyMap["tag"] = tag
	case "TLSA":
		bodyMap["license"] = license
		bodyMap["choicer"] = choicer
		bodyMap["match"] = match
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	var message responses.DnsMessage
	err = client.handleRequest(http.MethodPut, url, bytes.NewBuffer(body), &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

// DeleteRecord deletes a specific DNS record for a specific domain.
//
// It sends a DELETE request to the `/dns/{domain}/records/{record}/{type}/{content}`
// endpoint.
//
// Parameters:
//   - `domain`: The name of the domain the record belongs to.
//   - `record`: The name of the record to delete.
//   - `recordType`: The type of the record to delete.
//   - `contentId`: The content of the record to delete.
//
// Returns:
//   - `*responses.DnsMessage`: A message confirming the deletion of the record.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) DeleteRecord(domain, record, recordType, contentId string) (*responses.DnsMessage, error) {
	url := fmt.Sprintf(urls.RecordDeleteURL, urls.BaseUrl, domain, record, recordType, contentId)
	var message responses.DnsMessage
	err := client.handleRequest(http.MethodDelete, url, nil, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetDNSEvents retrieves a list of all DNS events.
//
// It sends a GET request to the `/dns/events` endpoint and returns a
// `DNSEventsResponse` struct containing a list of events.
//
// Returns:
//   - `*responses.DNSEventsResponse`: A struct containing the list of DNS events.
//   - `error`: An error if the request fails or the response cannot be decoded.
func (client *Client) GetDNSEvents() (*responses.DNSEventsResponse, error) {
	url := fmt.Sprintf(urls.DNSEventsURL, urls.BaseUrl)
	var dnsEvents responses.DNSEventsResponse
	err := client.handleRequest(http.MethodGet, url, nil, &dnsEvents)
	if err != nil {
		return nil, err
	}
	return &dnsEvents, nil
}
