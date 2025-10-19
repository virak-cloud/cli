// Package pkg provides shared packages and core utilities for the Virak CLI.
//
// The pkg package serves as the foundation for the Virak Cloud CLI application,
// containing essential components such as API endpoint definitions, HTTP client
// implementations, response type definitions, and shared utilities.
//
// Package Structure:
//   - urls.go: Centralized API endpoint definitions and URL constants
//   - http/: Complete HTTP client implementation for Virak Cloud API
//     - client.go: Core HTTP client with authentication and request handling
//     - interfaces.go: Service interfaces defining API contracts
//     - kubernetes.go: Kubernetes cluster management operations
//     - instance.go: Virtual machine lifecycle management
//     - network.go: Network configuration and security management
//     - dns.go: Domain and DNS record management
//     - object_storage.go: Object storage and bucket operations
//     - user.go: User account and SSH key management
//     - datacenter.go: Zone and datacenter information
//     - responses/: Response type definitions for API operations
//
// Key Features:
//   - Centralized API endpoint management with configurable base URLs
//   - Type-safe HTTP client with automatic authentication
//   - Comprehensive error handling and response parsing
//   - Interface-based design for easy testing and mocking
//   - Support for all Virak Cloud API services
//
// Usage:
//   The pkg package is primarily used by the CLI commands to interact with
//   the Virak Cloud API. End users typically interact with the CLI rather than
//   directly using this package.
//
//   // Example of creating an HTTP client
//   client := http.NewClient("api-token")
//   client.BaseURL = "https://public-api.virakcloud.com"
//
//   // Example of listing instances
//   instances, err := client.ListInstances("zone-123")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Found %d instances\n", len(instances.Data))
//
// Configuration:
//   - BaseUrl: Default API endpoint (can be overridden at build time)
//   - LoginUrl: Token creation URL (can be overridden at build time)
//   - API endpoints: Template strings for all API operations
//
// Thread Safety:
//   All HTTP client implementations are safe for concurrent use by multiple
//   goroutines. The client maintains immutable state after creation.
//
// For detailed API documentation, see the subpackages and individual function
// documentation within each package.
package pkg
