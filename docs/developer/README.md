# Developer Documentation

Welcome to the Virak CLI developer documentation. This section provides comprehensive information for developers who want to contribute to the Virak CLI project or understand its internal architecture.

## Table of Contents

- [Getting Started](#getting-started)
- [Architecture](#architecture)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [Development Setup](#development-setup)
- [Code Style](#code-style)
- [Testing](#testing)
- [Release Process](#release-process)

## Getting Started

If you're new to the Virak CLI project, start with these guides:

1. [Development Setup](development-setup.md) - Set up your development environment
2. [Architecture Overview](architecture.md) - Understand the system architecture
3. [API Documentation](api/README.md) - Learn about the API interfaces
4. [Contributing Guide](contributing.md) - Learn how to contribute

## Architecture

The Virak CLI is built with the following architectural principles:

- **Modular Design**: Commands are organized into separate packages
- **HTTP Client**: Centralized HTTP client for API communication
- **Configuration Management**: Flexible configuration system
- **Error Handling**: Consistent error handling across commands
- **Output Formatting**: Multiple output formats (table, JSON, YAML)

### Core Components

- **Command Structure**: Cobra-based command hierarchy
- **HTTP Client**: Custom HTTP client with authentication
- **Response Handling**: Structured response processing
- **Configuration**: YAML-based configuration management
- **Logging**: Structured logging with slog

## API Documentation

The API documentation provides detailed information about:

- [HTTP Client](api/http-client.md) - HTTP client implementation
- [API Interfaces](api/interfaces.md) - API interface specifications
- [Response Structures](api/responses.md) - Response type definitions

## Contributing

We welcome contributions from the community! Please read our [Contributing Guide](contributing.md) to understand:

- How to set up your development environment
- Our code style and conventions
- The pull request process
- Code review guidelines

### Quick Contribution Checklist

- [ ] Fork the repository
- [ ] Create a feature branch
- [ ] Make your changes
- [ ] Add tests for new functionality
- [ ] Ensure all tests pass
- [ ] Submit a pull request

## Development Setup

To set up your development environment:

1. **Prerequisites**
   - Go 1.24.4 or later
   - Git
   - Make (optional, for build scripts)

2. **Setup Steps**
   ```bash
   # Clone the repository
   git clone https://github.com/virak-cloud/cli.git
   cd cli
   
   # Install dependencies
   go mod download
   
   # Build the project
   go build -o virak-cli main.go
   ```

3. **Verification**
   ```bash
   # Test the build
   ./virak-cli --version
   
   # Run tests
   go test ./...
   ```

For detailed setup instructions, see [Development Setup](development-setup.md).

## Code Style

We follow the Go standard code style with additional conventions:

- Use `gofmt` for formatting
- Use `goimports` for import management
- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Write meaningful commit messages
- Add comments for public functions and types

### Example Code Style

```go
// Package example provides example functionality.
//
// This package demonstrates the expected code style and conventions
// used throughout the Virak CLI project.
package example

import (
    "context"
    "fmt"
    
    "github.com/virak-cloud/cli/internal/cli"
    "github.com/virak-cloud/cli/pkg/http"
)

// ExampleService provides example functionality.
// It demonstrates the expected structure for service implementations.
type ExampleService struct {
    client *http.Client
}

// NewExampleService creates a new ExampleService instance.
//
// Parameters:
//   - client: HTTP client for API communication
//
// Returns:
//   - *ExampleService: A new service instance
func NewExampleService(client *http.Client) *ExampleService {
    return &ExampleService{
        client: client,
    }
}

// DoSomething performs an example operation.
//
// This method demonstrates the expected pattern for implementing
// business logic in service methods.
//
// Parameters:
//   - ctx: Context for the operation
//   - param: Example parameter
//
// Returns:
//   - string: Result of the operation
//   - error: Error if the operation fails
func (s *ExampleService) DoSomething(ctx context.Context, param string) (string, error) {
    if param == "" {
        return "", fmt.Errorf("parameter cannot be empty")
    }
    
    // Implementation here...
    return "result", nil
}
```

## Testing

Testing is an essential part of the Virak CLI project:

### Test Structure

```
cmd/example/
├── example.go          # Implementation
├── example_test.go     # Unit tests
└── doc.go             # Documentation
```

### Writing Tests

```go
package example

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestExampleService_DoSomething(t *testing.T) {
    // Setup
    client := &mockHTTPClient{}
    service := NewExampleService(client)
    
    // Test cases
    tests := []struct {
        name    string
        param   string
        want    string
        wantErr bool
    }{
        {
            name:    "success case",
            param:   "test",
            want:    "result",
            wantErr: false,
        },
        {
            name:    "empty parameter",
            param:   "",
            want:    "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            if tt.param != "" {
                client.On("DoSomething", tt.param).Return(tt.want, nil)
            }
            
            // Execute
            got, err := service.DoSomething(context.Background(), tt.param)
            
            // Assert
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
            
            // Verify mocks
            client.AssertExpectations(t)
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./cmd/example

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Release Process

The Virak CLI uses automated releases through GitHub Actions:

1. **Versioning**: We follow semantic versioning (semver)
2. **Release Triggers**: Tags trigger the release process
3. **Build Process**: GoReleaser handles the build and release
4. **Distribution**: Releases are published to GitHub Releases

### Creating a Release

```bash
# Ensure all changes are committed and pushed
git add .
git commit -m "feat: add new feature"
git push origin main

# Tag the release
git tag v1.2.3
git push origin v1.2.3
```

The release workflow will automatically:
- Build binaries for all platforms
- Create archives and packages
- Generate a changelog
- Publish the release

## Resources

### External Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [GoReleaser](https://goreleaser.com/)
- [GitHub Actions](https://docs.github.com/en/actions)

### Internal Resources

- [Source Code](https://github.com/virak-cloud/cli)
- [Issues](https://github.com/virak-cloud/cli/issues)
- [Discussions](https://github.com/virak-cloud/cli/discussions)
- [Releases](https://github.com/virak-cloud/cli/releases)

## Getting Help

If you need help with development:

1. Check the [existing issues](https://github.com/virak-cloud/cli/issues)
2. Search the [discussions](https://github.com/virak-cloud/cli/discussions)
3. Create a new issue for bugs or feature requests
4. Join our community discussions

## Code of Conduct

Please read and follow our [Code of Conduct](https://github.com/virak-cloud/cli/blob/main/CODE_OF_CONDUCT.md) when participating in this project.

---

Thank you for contributing to Virak CLI!