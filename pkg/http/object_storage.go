// Package http provides object storage management operations for the Virak Cloud API.
//
// This file contains client methods for managing object storage buckets, including
// bucket creation, deletion, configuration, and event monitoring. Object storage
// provides scalable storage for unstructured data with S3-compatible API access.
// All operations are performed within the context of a specific zone.
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	urls "github.com/virak-cloud/cli/pkg"
	"net/http"

	"github.com/virak-cloud/cli/pkg/http/responses"
)

// GetObjectStorageBuckets retrieves all object storage buckets in the specified zone.
//
// Returns a list of all buckets accessible to the user within the specified zone,
// including their configuration, access policies, and usage statistics.
//
// Parameters:
//   - zoneId: Zone to query for buckets
//
// Returns:
//   - *responses.ObjectStorageBucketsResponse: List of buckets with metadata
//   - error: API error or network failure
func (client *Client) GetObjectStorageBuckets(zoneId string) (*responses.ObjectStorageBucketsResponse, error) {
	var result responses.ObjectStorageBucketsResponse
	url := fmt.Sprintf(urls.BucketList, urls.BaseUrl, zoneId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateObjectStorageBucket creates a new object storage bucket.
//
// Creates a new bucket with the specified name and access policy. Buckets provide
// scalable storage for objects with S3-compatible access. Bucket names must be
// globally unique across the entire platform.
//
// Parameters:
//   - zoneId: Zone where the bucket will be created
//   - name: Globally unique bucket name
//   - policy: Access policy (e.g., "private", "public-read", "public-read-write")
//
// Returns:
//   - *responses.ObjectStorageBucketCreationResponse: Bucket creation details with access credentials
//   - error: Validation error or API failure
func (client *Client) CreateObjectStorageBucket(zoneId string, name string, policy string) (*responses.ObjectStorageBucketCreationResponse, error) {
	var result responses.ObjectStorageBucketCreationResponse
	url := fmt.Sprintf(urls.BucketCreate, urls.BaseUrl, zoneId)

	body, err := json.Marshal(map[string]string{
		"name":   name,
		"policy": policy,
	})
	if err != nil {
		return nil, err
	}

	err = client.handleRequest(http.MethodPost, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetObjectStorageBucket(zoneId string, bucketId string) (*responses.ObjectStorageBucketResponse, error) {

	var result responses.ObjectStorageBucketResponse
	url := fmt.Sprintf(urls.BucketShow, urls.BaseUrl, zoneId, bucketId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

func (client *Client) UpdateObjectStorageBucket(zoneId string, bucketId string, policy string) (*responses.ObjectStorageBucketResponse, error) {

	var result responses.ObjectStorageBucketResponse
	url := fmt.Sprintf(urls.BucketUpdate, urls.BaseUrl, zoneId, bucketId)

	body, err := json.Marshal(map[string]string{
		"policy": policy,
	})
	if err != nil {
		return nil, err
	}

	err = client.handleRequest(http.MethodPut, url, bytes.NewBuffer(body), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

func (client *Client) DeleteObjectStorageBucket(zoneId, bucketId string) error {

	url := fmt.Sprintf(urls.BucketDelete, urls.BaseUrl, zoneId, bucketId)
	err := client.handleRequest(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}
	return nil

}

func (client *Client) GetObjectStorageEvents(zoneId string) (*responses.ObjectStorageEventsResponse, error) {
	var result responses.ObjectStorageEventsResponse
	url := fmt.Sprintf(urls.BucketsEventList, urls.BaseUrl, zoneId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *Client) GetObjectStorageBucketEvents(zoneId string, bucketId string) (*responses.ObjectStorageEventsResponse, error) {
	var result responses.ObjectStorageEventsResponse
	url := fmt.Sprintf(urls.BucketEventList, urls.BaseUrl, zoneId, bucketId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
