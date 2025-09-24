package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	s "strings"

	urls "virak-cli/pkg"

	"virak-cli/pkg/http/responses"
)

func (client *Client) GetDomains() (*responses.DomainList, error) {
	url := fmt.Sprintf(urls.DomainListURL, urls.BaseUrl)
	var domainList responses.DomainList
	err := client.handleRequest(http.MethodGet, url, nil, &domainList)
	if err != nil {
		return nil, err
	}
	return &domainList, nil
}

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

func (client *Client) DeleteDomain(domain string) (*responses.DnsMessage, error) {
	url := fmt.Sprintf(urls.DomainDeleteURL, urls.BaseUrl, domain)
	var message responses.DnsMessage
	err := client.handleRequest(http.MethodDelete, url, nil, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (client *Client) GetRecords(domain string) (*responses.RecordList, error) {
	url := fmt.Sprintf(urls.RecordListURL, urls.BaseUrl, domain)
	var recordList responses.RecordList
	err := client.handleRequest(http.MethodGet, url, nil, &recordList)
	if err != nil {
		return nil, err
	}
	return &recordList, nil
}

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

func (client *Client) DeleteRecord(domain, record, recordType, contentId string) (*responses.DnsMessage, error) {
	url := fmt.Sprintf(urls.RecordDeleteURL, urls.BaseUrl, domain, record, recordType, contentId)
	var message responses.DnsMessage
	err := client.handleRequest(http.MethodDelete, url, nil, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (client *Client) GetDNSEvents() (*responses.DNSEventsResponse, error) {
	url := fmt.Sprintf(urls.DNSEventsURL, urls.BaseUrl)
	var dnsEvents responses.DNSEventsResponse
	err := client.handleRequest(http.MethodGet, url, nil, &dnsEvents)
	if err != nil {
		return nil, err
	}
	return &dnsEvents, nil
}
