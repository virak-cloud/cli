package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	urls "virak-cli/pkg"

	"virak-cli/pkg/http/responses"
)

func (client *Client) GetObjectStorageBuckets(zoneId string) (*responses.ObjectStorageBucketsResponse, error) {

	var result responses.ObjectStorageBucketsResponse
	url := fmt.Sprintf(urls.BucketList, urls.BaseUrl, zoneId)
	err := client.handleRequest(http.MethodGet, url, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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
