# Contributing Guide

Thank you for your interest in contributing to Virak CLI! This guide will help you get started with contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Style](#code-style)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Pull Request Process](#pull-request-process)
- [Code Review](#code-review)
- [Release Process](#release-process)
- [Community](#community)

## Getting Started

### Prerequisites

Before you start contributing, make sure you have:

- Go 1.24.4 or later installed
- Git installed and configured
- A GitHub account
- Basic knowledge of Go and command-line tools

### First Steps

1. **Fork the Repository**
   ```bash
   # Fork the repository on GitHub
   # Then clone your fork
   git clone https://github.com/your-username/cli.git
   cd cli
   ```

2. **Add Upstream Remote**
   ```bash
   git remote add upstream https://github.com/virak-cloud/cli.git
   ```

3. **Create a Development Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Set Up Your Development Environment**
   ```bash
   # Install dependencies
   go mod download
   
   # Build the project
   go build -o virak-cli main.go
   
   # Run tests to ensure everything works
   go test ./...
   ```

## Development Setup

### Environment Configuration

1. **Install Go**
   - Download and install Go from [golang.org](https://golang.org/dl/)
   - Ensure Go is in your PATH
   - Verify installation: `go version`

2. **Install Development Tools**
   ```bash
   # Install Go formatting tools
   go install golang.org/x/tools/cmd/goimports@latest
   go install golang.org/x/tools/cmd/gofmt@latest
   
   # Install linting tools
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Install testing tools
   go install github.com/stretchr/testify/assert@latest
   ```

3. **Configure Git**
   ```bash
   git config --global user.name "Your Name"
   git config --global user.email "your-email@example.com"
   ```

### Project Structure

Understanding the project structure is essential for contributing:

```
virak-cli/
├── cmd/                    # CLI commands
│   ├── instance/           # Instance management commands
│   ├── network/            # Network management commands
│   ├── zone/               # Zone management commands
│   └── ...
├── internal/               # Internal packages
│   ├── cli/                # CLI utilities
│   ├── logger/             # Logging utilities
│   └── presenter/          # Output formatting
├── pkg/                    # Public packages
│   ├── http/               # HTTP client
│   └── responses/          # API response types
├── docs/                   # Documentation
│   ├── user-guide/         # User documentation
│   ├── reference/          # Command reference
│   └── developer/          # Developer documentation
├── scripts/                # Build and utility scripts
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies
└── README.md               # Project README
```

### Building and Testing

1. **Build the Project**
   ```bash
   # Build for current platform
   go build -o virak-cli main.go
   
   # Build for multiple platforms
   GOOS=linux GOARCH=amd64 go build -o virak-cli-linux-amd64 main.go
   GOOS=darwin GOARCH=amd64 go build -o virak-cli-darwin-amd64 main.go
   GOOS=windows GOARCH=amd64 go build -o virak-cli-windows-amd64.exe main.go
   ```

2. **Run Tests**
   ```bash
   # Run all tests
   go test ./...
   
   # Run tests with coverage
   go test -cover ./...
   
   # Run tests for a specific package
   go test ./cmd/instance
   
   # Run tests with verbose output
   go test -v ./...
   ```

3. **Check Code Quality**
   ```bash
   # Format code
   gofmt -w .
   goimports -w .
   
   # Run linter
   golangci-lint run
   
   # Check for vulnerabilities
   go list -json -m all | nancy sleuth
   ```

## Code Style

We follow the Go standard code style with additional conventions:

### Formatting

- Use `gofmt` for all code formatting
- Use `goimports` for import management
- Maximum line length: 120 characters
- Use tabs for indentation (Go standard)

### Naming Conventions

- **Packages**: Use short, lowercase names (e.g., `instance`, `network`)
- **Constants**: Use `UPPER_SNAKE_CASE` for exported constants
- **Variables**: Use `camelCase` for local variables
- **Functions**: Use `PascalCase` for exported functions, `camelCase` for unexported
- **Types**: Use `PascalCase` for all types

### Example Code Style

```go
// Package instance provides commands for managing virtual machine instances.
package instance

import (
    "context"
    "fmt"
    
    "github.com/spf13/cobra"
    "github.com/virak-cloud/cli/internal/cli"
    "github.com/virak-cloud/cli/pkg/http"
)

// DefaultInstanceTTL is the default TTL for instance operations.
const DefaultInstanceTTL = 3600

// InstanceService provides functionality for managing instances.
type InstanceService struct {
    client *http.Client
}

// NewInstanceService creates a new InstanceService.
//
// Parameters:
//   - client: HTTP client for API communication
//
// Returns:
//   - *InstanceService: A new service instance
func NewInstanceService(client *http.Client) *InstanceService {
    return &InstanceService{
        client: client,
    }
}

// CreateInstance creates a new virtual machine instance.
//
// This method creates a new instance with the specified parameters.
// It validates the input and makes an API request to create the instance.
//
// Parameters:
//   - ctx: Context for the operation
//   - zoneID: Zone ID where the instance will be created
//   - opts: Instance creation options
//
// Returns:
//   - *responses.Instance: The created instance
//   - error: Error if the operation fails
func (s *InstanceService) CreateInstance(ctx context.Context, zoneID string, opts CreateOptions) (*responses.Instance, error) {
    if err := opts.Validate(); err != nil {
        return nil, fmt.Errorf("invalid options: %w", err)
    }
    
    // Implementation here...
    return nil, nil
}
```

### Comments

- Comment all exported functions, types, and constants
- Use godoc format for comments
- Include parameter and return value descriptions
- Add usage examples for complex functions

### Error Handling

- Use structured error types
- Include context in error messages
- Wrap errors with `fmt.Errorf` and `%w` verb
- Handle errors appropriately at each layer

```go
// ValidationError represents a validation error.
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// ValidateOptions validates instance creation options.
func ValidateOptions(opts CreateOptions) error {
    if opts.Name == "" {
        return &ValidationError{
            Field:   "name",
            Message: "name is required",
        }
    }
    
    if opts.ServiceOfferingID == "" {
        return &ValidationError{
            Field:   "serviceOfferingId",
            Message: "service offering ID is required",
        }
    }
    
    return nil
}
```

## Testing

Testing is an essential part of the Virak CLI project:

### Test Structure

```
cmd/instance/
├── instance.go          # Implementation
├── instance_test.go     # Unit tests
├── example_test.go      # Example tests
└── doc.go              # Documentation
```

### Writing Tests

1. **Unit Tests**
   ```go
   func TestInstanceService_CreateInstance(t *testing.T) {
       tests := []struct {
           name    string
           zoneID  string
           opts    CreateOptions
           want    *responses.Instance
           wantErr bool
       }{
           {
               name:   "success case",
               zoneID: "zone-123",
               opts: CreateOptions{
                   Name:              "test-instance",
                   ServiceOfferingID: "offering-123",
                   VMImageID:         "image-123",
                   NetworkIDs:        []string{"net-123"},
               },
               want: &responses.Instance{
                   Name: "test-instance",
               },
               wantErr: false,
           },
           {
               name:   "validation error",
               zoneID: "zone-123",
               opts:   CreateOptions{}, // Empty options
               want:   nil,
               wantErr: true,
           },
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               client := &mockHTTPClient{}
               service := NewInstanceService(client)
               
               if tt.wantErr {
                   client.On("CreateInstance", mock.Anything, mock.Anything).
                       Return(nil, fmt.Errorf("validation error"))
               } else {
                   client.On("CreateInstance", mock.Anything, mock.Anything).
                       Return(tt.want, nil)
               }
               
               got, err := service.CreateInstance(context.Background(), tt.zoneID, tt.opts)
               
               if tt.wantErr {
                   assert.Error(t, err)
               } else {
                   assert.NoError(t, err)
                   assert.Equal(t, tt.want.Name, got.Name)
               }
               
               client.AssertExpectations(t)
           })
       }
   }
   ```

2. **Integration Tests**
   ```go
   func TestInstanceCommand_Integration(t *testing.T) {
       if testing.Short() {
           t.Skip("Skipping integration test in short mode")
       }
       
       // Setup test environment
       testZone := os.Getenv("TEST_ZONE_ID")
       if testZone == "" {
           t.Skip("TEST_ZONE_ID not set")
       }
       
       // Test actual command execution
       cmd := NewInstanceCreateCmd()
       cmd.SetArgs([]string{
           "--name", "test-instance",
           "--service-offering-id", "offering-123",
           "--vm-image-id", "image-123",
           "--network-ids", `["net-123"]`,
           "--zoneId", testZone,
       })
       
       err := cmd.Execute()
       assert.NoError(t, err)
   }
   ```

3. **Example Tests**
   ```go
   func ExampleInstanceService_CreateInstance() {
       client := &mockHTTPClient{}
       service := NewInstanceService(client)
       
       client.On("CreateInstance", mock.Anything, mock.Anything).
           Return(&responses.Instance{Name: "test-instance"}, nil)
       
       instance, err := service.CreateInstance(
           context.Background(),
           "zone-123",
           CreateOptions{Name: "test-instance"},
       )
       
       if err != nil {
           fmt.Printf("Error: %v\n", err)
           return
       }
       
       fmt.Printf("Created instance: %s\n", instance.Name)
       // Output: Created instance: test-instance
   }
   ```

### Test Coverage

- Aim for at least 80% test coverage
- Focus on testing business logic
- Test error paths as well as success paths
- Use table-driven tests for multiple test cases

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for a specific package
go test -v ./cmd/instance

# Run benchmarks
go test -bench=. ./...

# Run race condition tests
go test -race ./...
```

## Submitting Changes

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Maintenance tasks

#### Examples

```
feat(instance): add interactive instance creation

Add interactive mode for instance creation that prompts users for
required parameters instead of requiring all flags.

Closes #123

fix(network): handle empty network list

Fix panic when listing networks with empty response.
Add proper validation for network list response.

docs(api): update HTTP client documentation

Update HTTP client documentation with new authentication methods
and error handling patterns.

```

### Pre-commit Hooks

Set up pre-commit hooks to ensure code quality:

```bash
# Install pre-commit
pip install pre-commit

# Install pre-commit hooks
pre-commit install

# Run pre-commit on all files
pre-commit run --all-files
```

Example `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files

  - repo: https://github.com/tekwizely/pre-commit-golang
    rev: v1.0.0-rc.1
    hooks:
      - id: go-fmt
      - id: go-vet-mod
      - id: go-mod-tidy
      - id: golangci-lint
```

## Pull Request Process

### Creating a Pull Request

1. **Update Your Fork**
   ```bash
   git checkout main
   git pull upstream main
   git checkout your-feature-branch
   git rebase upstream/main
   ```

2. **Run Tests**
   ```bash
   go test ./...
   go vet ./...
   gofmt -w .
   goimports -w .
   golangci-lint run
   ```

3. **Create Pull Request**
   - Go to your fork on GitHub
   - Click "New Pull Request"
   - Select your feature branch
   - Fill out the PR template

### Pull Request Template

```markdown
## Description
Brief description of the changes made.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests performed
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests pass locally
- [ ] No breaking changes (or documented)
```

### Pull Request Review Process

1. **Automated Checks**
   - CI/CD pipeline runs tests
   - Code quality checks
   - Security scans

2. **Code Review**
   - At least one maintainer must review
   - Review focuses on:
     - Code quality and style
     - Functionality and correctness
     - Test coverage
     - Documentation

3. **Approval and Merge**
   - PR must be approved by a maintainer
   - All checks must pass
   - PR is merged to main branch

## Code Review

### Review Guidelines

When reviewing code, focus on:

1. **Functionality**
   - Does the code work as intended?
   - Are there edge cases not handled?
   - Is the logic correct?

2. **Code Quality**
   - Is the code readable and maintainable?
   - Does it follow project conventions?
   - Are there opportunities for refactoring?

3. **Testing**
   - Are tests comprehensive?
   - Do tests cover edge cases?
   - Are tests well-written?

4. **Documentation**
   - Is the code well-documented?
   - Are comments clear and useful?
   - Is user documentation updated?

### Review Process

1. **Automated Review**
   - CI/CD runs automated checks
   - Code quality tools run
   - Security scans performed

2. **Peer Review**
   - Reviewer adds comments and suggestions
   - Author addresses feedback
   - Reviewer approves changes

3. **Final Review**
   - Maintainer performs final review
   - Ensures all requirements are met
   - Approves for merge

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Steps

1. **Prepare Release**
   ```bash
   # Update version in go.mod
   # Update CHANGELOG.md
   # Create release branch
   git checkout -b release/v1.2.3
   ```

2. **Test Release**
   ```bash
   # Run full test suite
   go test ./...
   
   # Build release binaries
   goreleaser release --snapshot --clean
   ```

3. **Tag Release**
   ```bash
   # Merge release branch to main
   git checkout main
   git merge release/v1.2.3
   
   # Create tag
   git tag v1.2.3
   git push origin v1.2.3
   ```

4. **Publish Release**
   - GitHub Actions automatically builds and publishes
   - Release notes are generated
   - Binaries are uploaded to GitHub Releases

## Community

### Getting Help

- **GitHub Issues**: Report bugs or request features
- **GitHub Discussions**: Ask questions or discuss ideas
- **Documentation**: Check existing documentation

### Communication Channels

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For general questions and discussions
- **Email**: For security issues or private matters

### Code of Conduct

We are committed to providing a welcoming and inclusive environment. Please read and follow our [Code of Conduct](https://github.com/virak-cloud/cli/blob/main/CODE_OF_CONDUCT.md).

---

Thank you for contributing to Virak CLI! Your contributions help make this project better for everyone.