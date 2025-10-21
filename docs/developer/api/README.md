# API Documentation

This section provides detailed documentation about the Virak CLI API interfaces, HTTP client implementation, and response structures.

## Table of Contents

- [HTTP Client](#http-client)
- [API Interfaces](#api-interfaces)
- [Response Structures](#response-structures)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)

## HTTP Client

The HTTP client is the core component responsible for communicating with the Virak Cloud API. It provides a simple, consistent interface for making API requests.

### Client Overview

```go
type Client struct {
    token     string
    baseURL   string
    userAgent string
    timeout   time.Duration
    client    *http.Client
}
```

### Creating a Client

```go
// Create a new HTTP client with authentication token
client := http.NewClient("your-api-token")

// The client will automatically handle:
// - Authentication headers
// - Request timeouts
// - Error handling
// - Response parsing
```

### Client Configuration

The client can be configured with:

- **Base URL**: API base URL (default: https://api.virakcloud.com)
- **Timeout**: Request timeout (default: 30 seconds)
- **User Agent**: Custom user agent string
- **Retry Policy**: Automatic retry for failed requests

### Request Pattern

All API methods follow a consistent pattern:

```go
func (c *Client) MethodName(params) (*ResponseType, error) {
    // 1. Build URL
    url := fmt.Sprintf("%s/v1/endpoint", c.baseURL)
    
    // 2. Create request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // 3. Set headers
    c.setHeaders(req)
    
    // 4. Execute request
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // 5. Handle response
    if resp.StatusCode != http.StatusOK {
        return nil, c.handleErrorResponse(resp)
    }
    
    // 6. Parse response
    var result ResponseType
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &result, nil
}
```

### Example: Instance API

```go
// Get instances in a zone
func (c *Client) GetInstances(zoneID string) (*responses.InstanceListResponse, error) {
    url := fmt.Sprintf("%s/v1/zones/%s/instances", c.baseURL, zoneID)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    c.setHeaders(req)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, c.handleErrorResponse(resp)
    }
    
    var result responses.InstanceListResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &result, nil
}

// Create an instance
func (c *Client) CreateInstance(zoneID, offeringID, imageID string, networkIDs []string, name string) (*responses.InstanceResponse, error) {
    url := fmt.Sprintf("%s/v1/zones/%s/instances", c.baseURL, zoneID)
    
    requestBody := map[string]interface{}{
        "serviceOfferingId": offeringID,
        "vmImageId":         imageID,
        "networkIds":        networkIDs,
        "name":              name,
    }
    
    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    c.setHeaders(req)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return nil, c.handleErrorResponse(resp)
    }
    
    var result responses.InstanceResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &result, nil
}
```

## API Interfaces

The API is organized into logical interfaces based on resource types:

### Instance Interface

```go
type InstanceAPI interface {
    // List instances
    GetInstances(zoneID string) (*responses.InstanceListResponse, error)
    
    // Get instance details
    GetInstance(zoneID, instanceID string) (*responses.InstanceResponse, error)
    
    // Create instance
    CreateInstance(zoneID, offeringID, imageID string, networkIDs []string, name string) (*responses.InstanceResponse, error)
    
    // Update instance
    UpdateInstance(zoneID, instanceID string, opts UpdateOptions) (*responses.InstanceResponse, error)
    
    // Delete instance
    DeleteInstance(zoneID, instanceID string) error
    
    // Instance actions
    StartInstance(zoneID, instanceID string) error
    StopInstance(zoneID, instanceID string) error
    RebootInstance(zoneID, instanceID string) error
    
    // Instance snapshots
    CreateSnapshot(zoneID, instanceID, name, description string) (*responses.SnapshotResponse, error)
    ListSnapshots(zoneID, instanceID string) (*responses.SnapshotListResponse, error)
    DeleteSnapshot(zoneID, instanceID, snapshotID string) error
    RevertSnapshot(zoneID, instanceID, snapshotID string) error
    
    // Instance volumes
    AttachVolume(zoneID, instanceID, volumeID string) error
    DetachVolume(zoneID, instanceID, volumeID string) error
    ListVolumes(zoneID, instanceID string) (*responses.VolumeListResponse, error)
    
    // Instance metrics
    GetMetrics(zoneID, instanceID string) (*responses.MetricsResponse, error)
}
```

### Network Interface

```go
type NetworkAPI interface {
    // List networks
    ListNetworks(zoneID string) (*responses.NetworkListResponse, error)
    
    // Get network details
    GetNetwork(zoneID, networkID string) (*responses.NetworkResponse, error)
    
    // Create networks
    CreateL2Network(zoneID, offeringID, name string) (*responses.NetworkResponse, error)
    CreateL3Network(zoneID, offeringID, name, gateway, netmask string) (*responses.NetworkResponse, error)
    
    // Delete network
    DeleteNetwork(zoneID, networkID string) error
    
    // Network service offerings
    GetL2NetworkServiceOfferings(zoneID string) (*responses.NetworkServiceOfferingListResponse, error)
    GetL3NetworkServiceOfferings(zoneID string) (*responses.NetworkServiceOfferingListResponse, error)
    ListNetworkServiceOfferings(zoneID string) (*responses.NetworkServiceOfferingListResponse, error)
    
    // Network instances
    ConnectInstance(zoneID, networkID, instanceID string) error
    DisconnectInstance(zoneID, networkID, instanceID string) error
    ListNetworkInstances(zoneID, networkID string) (*responses.NetworkInstanceListResponse, error)
    
    // Firewall rules
    CreateFirewallRule(zoneID, networkID string, rule FirewallRule) (*responses.FirewallRuleResponse, error)
    ListFirewallRules(zoneID, networkID string) (*responses.FirewallRuleListResponse, error)
    DeleteFirewallRule(zoneID, networkID, ruleID string) error
    
    // Load balancers
    CreateLoadBalancer(zoneID, networkID string, lb LoadBalancer) (*responses.LoadBalancerResponse, error)
    ListLoadBalancers(zoneID, networkID string) (*responses.LoadBalancerListResponse, error)
    DeleteLoadBalancer(zoneID, networkID, lbID string) error
    
    // Public IPs
    ListPublicIPs(zoneID string) (*responses.PublicIPListResponse, error)
    AssociatePublicIP(zoneID, networkID, instanceID string) (*responses.PublicIPResponse, error)
    DisassociatePublicIP(zoneID, networkID, instanceID string) error
    
    // VPN
    EnableVPN(zoneID, networkID string) (*responses.VPNResponse, error)
    DisableVPN(zoneID, networkID string) error
    GetVPN(zoneID, networkID string) (*responses.VPNResponse, error)
    UpdateVPN(zoneID, networkID string, opts VPNOptions) error
}
```

### Zone Interface

```go
type ZoneAPI interface {
    // List zones
    ListZones() (*responses.ZoneListResponse, error)
    
    // Get zone details
    GetZone(zoneID string) (*responses.ZoneResponse, error)
    
    // Zone resources
    GetZoneResources(zoneID string) (*responses.ZoneResourcesResponse, error)
    
    // Zone services
    GetZoneServices(zoneID string) (*responses.ZoneServicesResponse, error)
    
    // Zone networks
    GetZoneNetworks(zoneID string) (*responses.NetworkListResponse, error)
}
```

### Bucket Interface

```go
type BucketAPI interface {
    // List buckets
    ListBuckets(zoneID string) (*responses.BucketListResponse, error)
    
    // Get bucket details
    GetBucket(zoneID, bucketID string) (*responses.BucketResponse, error)
    
    // Create bucket
    CreateBucket(zoneID, name, policy string) (*responses.BucketResponse, error)
    
    // Update bucket
    UpdateBucket(zoneID, bucketID, policy string) (*responses.BucketResponse, error)
    
    // Delete bucket
    DeleteBucket(zoneID, bucketID string) error
    
    // Bucket events
    GetBucketEvents(zoneID, bucketID string, limit, page int) (*responses.BucketEventListResponse, error)
}
```

### DNS Interface

```go
type DNSAPI interface {
    // List domains
    ListDomains() (*responses.DomainListResponse, error)
    
    // Get domain details
    GetDomain(domainName string) (*responses.DomainResponse, error)
    
    // Create domain
    CreateDomain(domainName string) (*responses.DomainResponse, error)
    
    // Delete domain
    DeleteDomain(domainName string) error
    
    // DNS records
    CreateRecord(domainName, recordName, recordType, content string, ttl, priority, weight, port, flags, tag, license, choicer, match int) (*responses.DNSRecordResponse, error)
    ListRecords(domainName string) (*responses.DNSRecordListResponse, error)
    UpdateRecord(domainName, recordName, recordType, content string, ttl, priority, weight, port, flags, tag, license, choicer, match int) (*responses.DNSRecordResponse, error)
    DeleteRecord(domainName, recordName, recordType string) error
    
    // DNS events
    GetDNSEvents() (*responses.DNSEventListResponse, error)
}
```

### Cluster Interface

```go
type ClusterAPI interface {
    // List clusters
    ListClusters(zoneID string) (*responses.ClusterListResponse, error)
    
    // Get cluster details
    GetCluster(zoneID, clusterID string) (*responses.ClusterResponse, error)
    
    // Create cluster
    CreateCluster(zoneID, name, versionID, offeringID, sshKeyID, networkID string, haEnabled bool, clusterSize int, description, privateRegistryUsername, privateRegistryPassword, privateRegistryURL string, haConfigControllerNodes int, haConfigExternalLBIP string) (*responses.ClusterResponse, error)
    
    // Update cluster
    UpdateCluster(zoneID, clusterID, name, description string) (*responses.ClusterResponse, error)
    
    // Delete cluster
    DeleteCluster(zoneID, clusterID string) error
    
    // Cluster actions
    StartCluster(zoneID, clusterID string) error
    StopCluster(zoneID, clusterID string) error
    ScaleCluster(zoneID, clusterID string, nodePoolSize int) error
    
    // Cluster service events
    GetClusterServiceEvents(zoneID, clusterID string) (*responses.ClusterServiceEventListResponse, error)
    
    // Cluster versions
    ListClusterVersions(zoneID string) (*responses.ClusterVersionListResponse, error)
    
    // Cluster service offerings
    ListClusterServiceOfferings(zoneID string) (*responses.ClusterServiceOfferingListResponse, error)
}
```

## Response Structures

All API responses follow a consistent structure:

### Base Response Structure

```go
type BaseResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Errors  interface{} `json:"errors,omitempty"`
}
```

### List Response Structure

```go
type ListResponse struct {
    Data     []interface{} `json:"data"`
    Metadata Metadata       `json:"metadata"`
}

type Metadata struct {
    TotalCount int `json:"totalCount"`
    Page       int `json:"page"`
    PageSize   int `json:"pageSize"`
}
```

### Instance Response Structures

```go
type InstanceResponse struct {
    Data Instance `json:"data"`
}

type InstanceListResponse struct {
    Data     []Instance `json:"data"`
    Metadata Metadata   `json:"metadata"`
}

type Instance struct {
    ID                   string                `json:"id"`
    Name                 string                `json:"name"`
    Status               string                `json:"status"`
    InstanceStatus       string                `json:"instanceStatus"`
    ZoneID               string                `json:"zoneId"`
    CreatedAt            int64                 `json:"createdAt"`
    UpdatedAt            int64                 `json:"updatedAt"`
    Username             string                `json:"username"`
    Password             string                `json:"password"`
    ServiceOfferingID    string                `json:"serviceOfferingId"`
    TemplateID           *string               `json:"templateId,omitempty"`
    DiskOfferingID       *string               `json:"diskOfferingId,omitempty"`
    KubernetesClusterID  *string               `json:"kubernetesClusterId,omitempty"`
    VMImage              *VMImage              `json:"vmImage,omitempty"`
    ServiceOffering      *ServiceOffering      `json:"serviceOffering,omitempty"`
    Zone                 *Zone                 `json:"zone,omitempty"`
}
```

### Network Response Structures

```go
type NetworkResponse struct {
    Data Network `json:"data"`
}

type NetworkListResponse struct {
    Data     []Network `json:"data"`
    Metadata Metadata  `json:"metadata"`
}

type Network struct {
    ID               string           `json:"id"`
    Name             string           `json:"name"`
    Status           string           `json:"status"`
    ZoneID           string           `json:"zoneId"`
    NetworkOffering  NetworkOffering  `json:"networkOffering"`
    CreatedAt        int64            `json:"createdAt"`
    UpdatedAt        int64            `json:"updatedAt"`
}

type NetworkOffering struct {
    ID                       string  `json:"id"`
    Name                     string  `json:"name"`
    DisplayName              string  `json:"displayName"`
    Type                     string  `json:"type"`
    InternetProtocol         string  `json:"internetProtocol"`
    Description              string  `json:"description"`
    NetworkRate              int     `json:"networkRate"`
    TrafficTransferPlan      int     `json:"trafficTransferPlan"`
    HourlyStartedPrice      float64 `json:"hourlyStartedPrice"`
    TrafficTransferOverprice float64 `json:"trafficTransferOverprice"`
}
```

### Error Response Structure

```go
type ErrorResponse struct {
    Error APIError `json:"error"`
}

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details"`
}
```

## Authentication

The API uses token-based authentication:

### Authentication Header

```go
func (c *Client) setHeaders(req *http.Request) {
    req.Header.Set("Authorization", "Bearer "+c.token)
    req.Header.Set("User-Agent", c.userAgent)
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
}
```

### Token Management

```go
// Set authentication token
func (c *Client) SetToken(token string) {
    c.token = token
}

// Get current token
func (c *Client) GetToken() string {
    return c.token
}

// Check if authenticated
func (c *Client) IsAuthenticated() bool {
    return c.token != ""
}
```

## Error Handling

The HTTP client provides comprehensive error handling:

### Error Types

```go
// Network errors
type NetworkError struct {
    URL     string
    Message string
    Err     error
}

// API errors
type APIError struct {
    Code    int
    Message string
    Details string
}

// Validation errors
type ValidationError struct {
    Field   string
    Message string
}

// Authentication errors
type AuthenticationError struct {
    Message string
}
```

### Error Handling Pattern

```go
func (c *Client) handleErrorResponse(resp *http.Response) error {
    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read error response: %w", err)
    }
    
    // Try to parse as API error
    var apiErr APIError
    if err := json.Unmarshal(body, &apiErr); err == nil {
        return &apiErr
    }
    
    // Return generic HTTP error
    return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
}
```

### Error Response Examples

```json
// Authentication error
{
  "error": {
    "code": 401,
    "message": "Authentication failed",
    "details": "Invalid or expired token"
  }
}

// Validation error
{
  "error": {
    "code": 422,
    "message": "Validation error",
    "details": "Name is required"
  }
}

// Not found error
{
  "error": {
    "code": 404,
    "message": "Resource not found",
    "details": "Instance 'instance-123' does not exist"
  }
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

### Rate Limit Headers

```go
type RateLimitInfo struct {
    Limit     int `json:"limit"`
    Remaining int `json:"remaining"`
    Reset     int `json:"reset"`
}

func (c *Client) parseRateLimitHeaders(resp *http.Response) RateLimitInfo {
    info := RateLimitInfo{}
    
    if limit := resp.Header.Get("X-RateLimit-Limit"); limit != "" {
        info.Limit, _ = strconv.Atoi(limit)
    }
    
    if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
        info.Remaining, _ = strconv.Atoi(remaining)
    }
    
    if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
        info.Reset, _ = strconv.Atoi(reset)
    }
    
    return info
}
```

### Rate Limit Handling

```go
func (c *Client) handleRateLimit(resp *http.Response) error {
    if resp.StatusCode == http.StatusTooManyRequests {
        retryAfter := resp.Header.Get("Retry-After")
        if retryAfter != "" {
            if seconds, err := strconv.Atoi(retryAfter); err == nil {
                return &RateLimitError{
                    RetryAfter: time.Duration(seconds) * time.Second,
                }
            }
        }
        
        return &RateLimitError{
            RetryAfter: time.Minute, // Default retry after
        }
    }
    
    return nil
}
```

---

For more detailed information about specific API endpoints, see the individual API interface documentation.