# System Architecture

This document provides an overview of the Virak CLI system architecture, including its components, design patterns, and data flow.

## Table of Contents

- [Overview](#overview)
- [Architecture Principles](#architecture-principles)
- [Component Architecture](#component-architecture)
- [Command Structure](#command-structure)
- [Data Flow](#data-flow)
- [HTTP Client Architecture](#http-client-architecture)
- [Configuration Management](#configuration-management)
- [Error Handling](#error-handling)
- [Output Formatting](#output-formatting)
- [Plugin System](#plugin-system)

## Overview

The Virak CLI is a command-line interface for interacting with the Virak Cloud API. It's built using Go and follows a modular architecture that separates concerns and promotes maintainability.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Virak CLI                           │
├─────────────────────────────────────────────────────────────┤
│  Command Layer (cmd/)                                       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │   instance  │ │   network   │ │    zone     │           │
│  │   commands  │ │  commands   │ │  commands   │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
├─────────────────────────────────────────────────────────────┤
│  Service Layer (internal/)                                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │     cli     │ │   logger    │ │  presenter  │           │
│  │  utilities  │ │  utilities  │ │  utilities  │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
├─────────────────────────────────────────────────────────────┤
│  HTTP Client Layer (pkg/http/)                              │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │    client   │ │  responses  │ │    urls     │           │
│  │   wrapper   │ │  types      │ │  constants  │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
├─────────────────────────────────────────────────────────────┤
│  Configuration Layer                                        │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │   config    │ │   auth      │ │   zone      │           │
│  │  management │ │  management │ │  management │           │
│  └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

## Architecture Principles

### 1. Separation of Concerns

Each component has a specific responsibility:
- **Commands**: Handle user interaction and command parsing
- **Services**: Implement business logic
- **HTTP Client**: Handle API communication
- **Configuration**: Manage settings and state

### 2. Dependency Injection

Dependencies are injected rather than hard-coded:
```go
// HTTP client is injected into services
func NewInstanceService(client *http.Client) *InstanceService {
    return &InstanceService{client: client}
}
```

### 3. Interface-Based Design

Components interact through interfaces to promote testability:
```go
type HTTPClient interface {
    GetInstances(zoneID string) (*InstanceListResponse, error)
    CreateInstance(zoneID string, opts CreateOptions) (*InstanceResponse, error)
}
```

### 4. Error Handling

Consistent error handling with structured error types:
```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details"`
}
```

## Component Architecture

### Command Layer (cmd/)

The command layer contains all CLI commands organized by functionality:

```
cmd/
├── instance/          # Instance management commands
│   ├── instance.go
│   ├── instance_create.go
│   ├── instance_list.go
│   └── ...
├── network/           # Network management commands
│   ├── network.go
│   ├── create/
│   ├── firewall/
│   └── ...
├── zone/              # Zone management commands
├── bucket/            # Bucket management commands
├── cluster/           # Kubernetes cluster commands
├── dns/               # DNS management commands
├── login.go           # Authentication command
├── logout.go          # Logout command
└── root.go            # Root command
```

#### Command Structure Pattern

Each command follows this pattern:

```go
// Command options struct
type createOptions struct {
    ZoneID    string `flag:"zoneId" usage:"Zone ID"`
    Name      string `flag:"name" usage:"Instance name"`
    // ... other options
}

// Command variable
var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new instance",
    PreRunE: func(cmd *cobra.Command, args []string) error {
        // Validation logic
        return cli.Validate(cmd, cli.Required("name"))
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        // Command execution logic
        return executeCreate(cmd)
    },
}

// Initialization
func init() {
    InstanceCmd.AddCommand(createCmd)
    _ = cli.BindFlagsFromStruct(createCmd, &createOpts)
}
```

### Service Layer (internal/)

The service layer contains business logic and utilities:

```
internal/
├── cli/               # CLI utilities and validation
│   ├── bind.go
│   ├── validate.go
│   ├── preflight.go
│   └── ...
├── logger/            # Logging utilities
│   └── logger.go
└── presenter/         # Output formatting
    ├── instance.go
    ├── network.go
    └── ...
```

#### CLI Utilities

The CLI utilities provide common functionality:

```go
// Validation functions
func Validate(cmd *cobra.Command, rules ...ValidationRule) error

// Flag binding
func BindFlagsFromStruct(cmd *cobra.Command, opts interface{}) error

// Preflight checks
func Preflight(requireZone bool) func(*cobra.Command, []string) error
```

### HTTP Client Layer (pkg/http/)

The HTTP client layer handles all API communication:

```
pkg/http/
├── client.go          # Main HTTP client
├── responses/         # Response type definitions
│   ├── instance.go
│   ├── network.go
│   └── ...
└── urls.go            # API endpoint URLs
```

#### HTTP Client Design

```go
type Client struct {
    token    string
    baseURL string
    client   *http.Client
}

func NewClient(token string) *Client {
    return &Client{
        token:    token,
        baseURL: "https://api.virakcloud.com",
        client:   &http.Client{Timeout: 30 * time.Second},
    }
}

func (c *Client) GetInstances(zoneID string) (*responses.InstanceListResponse, error) {
    url := fmt.Sprintf("%s/v1/zones/%s/instances", c.baseURL, zoneID)
    // Implementation...
}
```

## Command Structure

### Command Hierarchy

The CLI uses Cobra for command-line parsing:

```
virak-cli
├── login
├── logout
├── bucket
│   ├── create
│   ├── delete
│   ├── list
│   └── show
├── cluster
│   ├── create
│   ├── delete
│   ├── list
│   └── ...
├── dns
│   ├── domain
│   │   ├── create
│   │   ├── delete
│   │   └── ...
│   └── record
│       ├── create
│       ├── delete
│       └── ...
├── instance
│   ├── create
│   ├── delete
│   ├── list
│   └── ...
├── network
│   ├── create
│   │   ├── l2
│   │   └── l3
│   ├── delete
│   ├── list
│   └── ...
└── zone
    ├── list
    ├── networks
    ├── resources
    └── services
```

### Command Lifecycle

1. **Initialization**: Command is registered with Cobra
2. **PreRun**: Validation and preflight checks
3. **Run**: Command execution
4. **Output**: Results are formatted and displayed

```go
var instanceCreateCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new instance",
    PreRunE: func(cmd *cobra.Command, args []string) error {
        // 1. Preflight checks (authentication, zone)
        if err := cli.Preflight(true)(cmd, args); err != nil {
            return err
        }
        
        // 2. Input validation
        return cli.Validate(cmd,
            cli.Required("name"),
            cli.Required("service-offering-id"),
        )
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        // 3. Execute command
        return executeInstanceCreate(cmd)
    },
}
```

## Data Flow

### Typical Command Execution Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   User      │    │   Command   │    │   Service   │    │  HTTP API   │
│   Input     │───▶│   Layer     │───▶│   Layer     │───▶│   Layer     │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                           │                     │                     │
                           ▼                     ▼                     ▼
                    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
                    │ Validation  │    │ Business    │    │ API Request │
                    │ & Preflight │    │ Logic       │    │ & Response  │
                    └─────────────┘    └─────────────┘    └─────────────┘
                           │                     │                     │
                           ▼                     ▼                     ▼
                    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
                    │   Error     │    │   Data      │    │   Error     │
                    │  Handling   │    │ Processing  │    │  Handling   │
                    └─────────────┘    └─────────────┘    └─────────────┘
                           │                     │                     │
                           ▼                     ▼                     ▼
                    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
                    │   Output    │    │   Result    │    │   Result    │
                    │ Formatting  │    │ Formatting  │    │ Processing  │
                    └─────────────┘    └─────────────┘    └─────────────┘
```

### Example: Instance Creation

1. **User Input**:
   ```bash
   virak-cli instance create --name my-instance --service-offering-id abc123
   ```

2. **Command Layer**:
   ```go
   // Parse flags and validate
   if err := cli.Validate(cmd, cli.Required("name")); err != nil {
       return err
   }
   ```

3. **Service Layer**:
   ```go
   // Business logic
   httpClient := http.NewClient(token)
   instance, err := httpClient.CreateInstance(zoneID, options)
   ```

4. **HTTP Client Layer**:
   ```go
   // API request
   resp, err := c.client.Post(url, body, headers)
   ```

5. **Response Processing**:
   ```go
   // Parse and format response
   presenter.RenderInstanceCreateSuccess(instance)
   ```

## HTTP Client Architecture

### Client Design

The HTTP client is designed to be:
- **Reusable**: Single client instance for all API calls
- **Configurable**: Customizable timeouts, retries, and headers
- **Testable**: Interface-based design for mocking

```go
type Client struct {
    token     string
    baseURL   string
    userAgent string
    timeout   time.Duration
    client    *http.Client
}

// Interface for testing
type HTTPClient interface {
    GetInstances(zoneID string) (*responses.InstanceListResponse, error)
    CreateInstance(zoneID string, opts CreateOptions) (*responses.InstanceResponse, error)
}
```

### Request/Response Pattern

All API methods follow this pattern:

```go
func (c *Client) GetInstances(zoneID string) (*responses.InstanceListResponse, error) {
    // 1. Build URL
    url := c.buildURL("zones", zoneID, "instances")
    
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
    var result responses.InstanceListResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &result, nil
}
```

## Configuration Management

### Configuration Structure

Configuration is managed through a YAML file:

```yaml
auth:
  token: "your-api-token"
  
default:
  zoneId: "zone-123"
  zoneName: "Tehran-1"
  
api:
  baseURL: "https://api.virakcloud.com"
  timeout: 30s
  
output:
  format: "table"
  debug: false
```

### Configuration Loading

```go
type Config struct {
    Auth    AuthConfig    `yaml:"auth"`
    Default DefaultConfig `yaml:"default"`
    API     APIConfig     `yaml:"api"`
    Output  OutputConfig  `yaml:"output"`
}

func LoadConfig() (*Config, error) {
    configPath := getConfigPath()
    
    // Check if config exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        return &Config{}, nil // Return default config
    }
    
    // Load and parse config
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    return &config, nil
}
```

## Error Handling

### Error Types

The CLI uses structured error types:

```go
// API errors from the server
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details"`
}

// Client-side errors
type ValidationError struct {
    Field   string
    Message string
}

type ConfigError struct {
    Path    string
    Message string
}
```

### Error Handling Pattern

```go
func (c *Client) GetInstances(zoneID string) (*responses.InstanceListResponse, error) {
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Handle HTTP errors
    if resp.StatusCode != http.StatusOK {
        return nil, c.handleHTTPError(resp)
    }
    
    // Parse response
    var result responses.InstanceListResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &result, nil
}

func (c *Client) handleHTTPError(resp *http.Response) error {
    var apiErr APIError
    if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
        return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    }
    
    return fmt.Errorf("API error %d: %s - %s", apiErr.Code, apiErr.Message, apiErr.Details)
}
```

## Output Formatting

### Presenter Pattern

The presenter pattern handles output formatting:

```go
type Presenter interface {
    RenderTable(data interface{})
    RenderJSON(data interface{})
    RenderYAML(data interface{})
}

type InstancePresenter struct{}

func (p *InstancePresenter) RenderInstanceList(instances []responses.Instance) {
    switch getOutputFormat() {
    case "json":
        p.renderJSON(instances)
    case "yaml":
        p.renderYAML(instances)
    default:
        p.renderTable(instances)
    }
}

func (p *InstancePresenter) renderTable(instances []responses.Instance) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"ID", "Name", "Status", "Created"})
    
    for _, instance := range instances {
        table.Append([]string{
            instance.ID,
            instance.Name,
            instance.Status,
            formatTime(instance.CreatedAt),
        })
    }
    
    table.Render()
}
```

### Output Formats

The CLI supports multiple output formats:

- **Table**: Human-readable table format (default)
- **JSON**: Machine-readable JSON format
- **YAML**: Human-readable YAML format

```bash
# Table output (default)
virak-cli instance list

# JSON output
virak-cli --output json instance list

# YAML output
virak-cli --output yaml instance list
```

## Plugin System

### Command Registration

Commands are registered through initialization functions:

```go
// In cmd/instance/instance.go
func init() {
    rootCmd.AddCommand(InstanceCmd)
}

// In cmd/instance/instance_create.go
func init() {
    InstanceCmd.AddCommand(instanceCreateCmd)
}
```

### Modular Design

Each command package is self-contained:

```
cmd/instance/
├── instance.go          # Package initialization
├── instance_create.go   # Create command
├── instance_list.go     # List command
├── instance_show.go     # Show command
└── doc.go              # Package documentation
```

This design allows:
- Easy addition of new commands
- Independent testing of command packages
- Clear separation of concerns
- Maintainable codebase

## Future Enhancements

### Planned Improvements

1. **Plugin System**: Dynamic loading of command plugins
2. **Middleware**: Request/response middleware for cross-cutting concerns
3. **Caching**: Response caching for improved performance
4. **Streaming**: Support for streaming responses
5. **Async Operations**: Background operation support

### Extensibility

The architecture is designed to be extensible:

- New commands can be added without modifying existing code
- Output formatters can be added for new formats
- HTTP client can be extended with new features
- Configuration can be extended with new options

---

This architecture document provides a high-level overview of the Virak CLI system. For more detailed information about specific components, see the relevant API documentation.