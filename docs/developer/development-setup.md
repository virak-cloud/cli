# Development Setup Guide

This guide will help you set up your development environment for contributing to the Virak CLI project.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Project Setup](#project-setup)
- [Development Workflow](#development-workflow)
- [Building](#building)
- [Testing](#testing)
- [Debugging](#debugging)
- [IDE Configuration](#ide-configuration)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before you begin, ensure you have the following installed:

### Required Software

1. **Go** (version 1.24.4 or later)
   - Download from [golang.org](https://golang.org/dl/)
   - Follow the installation instructions for your operating system
   - Verify installation: `go version`

2. **Git**
   - Download from [git-scm.com](https://git-scm.com/)
   - Configure Git:
     ```bash
     git config --global user.name "Your Name"
     git config --global user.email "your-email@example.com"
     ```

3. **Make** (optional, for build scripts)
   - Linux: Usually pre-installed
   - macOS: Install Xcode Command Line Tools: `xcode-select --install`
   - Windows: Install through WSL or MinGW

### Recommended Tools

1. **Visual Studio Code** (or your preferred IDE)
   - Install Go extension
   - Install other helpful extensions

2. **Docker** (optional, for testing in containers)
   - Download from [docker.com](https://docker.com/)

3. **GitHub CLI** (optional, for GitHub interactions)
   - Install from [cli.github.com](https://cli.github.com/)

## Installation

### 1. Clone the Repository

```bash
# Fork the repository on GitHub first
# Then clone your fork
git clone https://github.com/your-username/cli.git
cd cli

# Add upstream remote
git remote add upstream https://github.com/virak-cloud/cli.git
```

### 2. Install Go Dependencies

```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

### 3. Install Development Tools

```bash
# Install Go formatting tools
go install golang.org/x/tools/cmd/goimports@latest
go install golang.org/x/tools/cmd/gofmt@latest

# Install linting tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install testing tools
go install github.com/stretchr/testify/assert@latest

# Install godoc for documentation
go install golang.org/x/tools/cmd/godoc@latest
```

### 4. Verify Installation

```bash
# Build the project
go build -o virak-cli main.go

# Test the build
./virak-cli --version

# Run tests
go test ./...
```

## Project Setup

### Environment Variables

Create a `.env` file in the project root for local development:

```bash
# .env
VIRAK_API_URL=https://api.virakcloud.com
VIRAK_TOKEN=your-test-token
VIRAK_ZONE_ID=zone-123
```

Add `.env` to your `.gitignore` file:

```bash
# .gitignore
.env
*.env
```

### Configuration File

Create a local configuration file for testing:

```bash
# ~/.virak-cli-dev.yaml
auth:
  token: "your-test-token"
  
default:
  zoneId: "zone-123"
  zoneName: "Test Zone"
  
api:
  baseURL: "https://api.virakcloud.com"
  timeout: 30s
  
output:
  format: "json"
  debug: true
```

### IDE Configuration

#### VS Code

Create `.vscode/settings.json`:

```json
{
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "go.gopath": "",
  "go.goroot": "",
  "go.formatOnSave": "file",
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.testOnSave": false,
  "go.coverOnSave": false,
  "go.coverageDecorator": {
    "type": "gutter",
    "coveredHighlightColor": "rgba(64,128,64,0.5)",
    "uncoveredHighlightColor": "rgba(128,64,64,0.25)"
  },
  "files.associations": {
    "*.go": "go"
  },
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

Create `.vscode/launch.json` for debugging:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch virak-cli",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/main.go",
      "env": {
        "VIRAK_CONFIG": "${workspaceFolder}/.virak-cli-dev.yaml"
      },
      "args": []
    },
    {
      "name": "Launch virak-cli with args",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/main.go",
      "env": {
        "VIRAK_CONFIG": "${workspaceFolder}/.virak-cli-dev.yaml"
      },
      "args": ["instance", "list"]
    }
  ]
}
```

Create `.vscode/tasks.json` for build tasks:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Build",
      "type": "shell",
      "command": "go",
      "args": ["build", "-o", "virak-cli", "main.go"],
      "group": {
        "kind": "build",
        "isDefault": true
      },
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "Test",
      "type": "shell",
      "command": "go",
      "args": ["test", "./..."],
      "group": {
        "kind": "test",
        "isDefault": true
      },
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "Lint",
      "type": "shell",
      "command": "golangci-lint",
      "args": ["run"],
      "group": "build",
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    }
  ]
}
```

## Development Workflow

### 1. Create a Feature Branch

```bash
# Sync with upstream
git checkout main
git pull upstream main

# Create a new branch
git checkout -b feature/your-feature-name
```

### 2. Make Changes

1. **Write Code**
   - Follow the project's code style
   - Add comments for public functions
   - Write tests for new functionality

2. **Format Code**
   ```bash
   # Format code
   gofmt -w .
   goimports -w .
   ```

3. **Run Linter**
   ```bash
   # Run linter
   golangci-lint run
   ```

4. **Run Tests**
   ```bash
   # Run all tests
   go test ./...
   
   # Run tests with coverage
   go test -cover ./...
   ```

### 3. Test Your Changes

```bash
# Build the project
go build -o virak-cli main.go

# Test your changes
./virak-cli --help
./virak-cli instance list --zoneId zone-123
```

### 4. Commit Changes

```bash
# Stage changes
git add .

# Commit changes
git commit -m "feat(instance): add new feature"

# Push to your fork
git push origin feature/your-feature-name
```

### 5. Create a Pull Request

1. Go to your fork on GitHub
2. Click "New Pull Request"
3. Select your feature branch
4. Fill out the PR template
5. Submit the pull request

## Building

### Build for Current Platform

```bash
# Simple build
go build -o virak-cli main.go

# Build with version info
go build -ldflags "-X main.version=1.2.3" -o virak-cli main.go
```

### Build for Multiple Platforms

```bash
# Build for Linux AMD64
GOOS=linux GOARCH=amd64 go build -o virak-cli-linux-amd64 main.go

# Build for macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o virak-cli-darwin-amd64 main.go

# Build for Windows AMD64
GOOS=windows GOARCH=amd64 go build -o virak-cli-windows-amd64.exe main.go

# Build for Linux ARM64
GOOS=linux GOARCH=arm64 go build -o virak-cli-linux-arm64 main.go

# Build for macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o virak-cli-darwin-arm64 main.go
```

### Build Script

Create a build script for convenience:

```bash
#!/bin/bash
# build.sh

set -e

VERSION=${1:-"dev"}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD)

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

# Create build directory
mkdir -p build

# Build for multiple platforms
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o build/virak-cli-linux-amd64 main.go

echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o build/virak-cli-darwin-amd64 main.go

echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o build/virak-cli-windows-amd64.exe main.go

echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o build/virak-cli-linux-arm64 main.go

echo "Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o build/virak-cli-darwin-arm64 main.go

echo "Build complete!"
ls -la build/
```

Make it executable:

```bash
chmod +x build.sh
```

Run the build script:

```bash
./build.sh 1.2.3
```

## Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test -v ./cmd/instance

# Run a specific test
go test -v ./cmd/instance -run TestInstanceCreate
```

### Generate Coverage Report

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

### Run Benchmarks

```bash
# Run benchmarks
go test -bench=. ./...

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./...
```

### Run Race Condition Tests

```bash
# Run tests with race detection
go test -race ./...
```

### Integration Tests

Create integration tests for end-to-end testing:

```bash
# Run integration tests
go test -tags=integration ./...

# Run integration tests with specific environment
VIRAK_API_URL=https://api-test.virakcloud.com go test -tags=integration ./...
```

## Debugging

### Debug with VS Code

1. Set breakpoints in your code
2. Use the "Launch virak-cli" configuration in `.vscode/launch.json`
3. Press F5 to start debugging

### Debug with Delve

1. Install Delve:
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   ```

2. Debug your application:
   ```bash
   # Debug main.go
   dlv debug main.go
   
   # Debug with arguments
   dlv debug main.go -- instance list
   
   # Debug tests
   dlv test ./cmd/instance
   ```

### Debug with Print Statements

Add debug logging to your code:

```go
import "log/slog"

func SomeFunction() {
    slog.Debug("Starting SomeFunction", "param1", value1)
    
    // Your code here
    
    slog.Info("Operation completed", "result", result)
}
```

Enable debug logging:

```bash
# Run with debug logging
./virak-cli --debug instance list
```

## IDE Configuration

### GoLand (JetBrains)

1. **Open Project**
   - File → Open → Select the project directory
   - GoLand will detect the Go module

2. **Configure Go SDK**
   - File → Settings → Go → GOROOT
   - Select your Go installation

3. **Configure Code Style**
   - File → Settings → Editor → Code Style → Go
   - Configure formatting preferences

4. **Configure Run/Debug Configurations**
   - Run → Edit Configurations
   - Add new Go configuration
   - Set File to `main.go`
   - Set Program arguments as needed

### Vim/Neovim

1. **Install vim-go**
   ```bash
   :Plug 'fatih/vim-go'
   ```

2. **Configure vim-go**
   ```vim
   " .vimrc
   let g:go_fmt_command = "goimports"
   let g:go_fmt_autosave = 1
   let g:go_lint_on_save = 1
   let g:go_test_autosave = 0
   let g:go_highlight_types = 1
   let g:go_highlight_fields = 1
   let g:go_highlight_functions = 1
   let g:go_highlight_methods = 1
   let g:go_highlight_structs = 1
   let g:go_highlight_interfaces = 1
   let g:go_highlight_build_constraints = 1
   let g:go_highlight_extra_types = 1
   ```

3. **Install binaries**
   ```vim
   :GoInstallBinaries
   ```

### Emacs

1. **Install go-mode**
   ```elisp
   ;; .emacs.d/init.el
   (require 'go-mode)
   (require 'go-eldoc)
   
   ;; Format on save
   (add-hook 'before-save-hook 'gofmt-before-save)
   
   ;; Run tests
   (global-set-key (kbd "C-c t") 'go-test-current-test)
   ```

## Troubleshooting

### Common Issues

1. **Go Not Found**
   ```
   bash: go: command not found
   ```
   Solution: Install Go and add it to your PATH

2. **Module Download Fails**
   ```
   go: module download failed
   ```
   Solution: Check your internet connection and try:
   ```bash
   go clean -modcache
   go mod download
   ```

3. **Build Fails with Import Errors**
   ```
   import "github.com/virak-cloud/cli/pkg/http": cannot find package
   ```
   Solution: Run `go mod tidy` to fix dependencies

4. **Tests Fail with Connection Errors**
   ```
   Error: dial tcp: connection refused
   ```
   Solution: Check your API configuration and network connection

5. **Linting Fails**
   ```
   golangci-lint run: exit status 1
   ```
   Solution: Fix the linting issues reported by golangci-lint

### Getting Help

If you're having trouble with development setup:

1. **Check the Documentation**
   - Read this guide carefully
   - Check the [Contributing Guide](contributing.md)

2. **Search Existing Issues**
   - Check [GitHub Issues](https://github.com/virak-cloud/cli/issues)
   - Search for similar problems

3. **Create a New Issue**
   - If you can't find a solution, create a new issue
   - Provide detailed information about your problem
   - Include error messages and steps to reproduce

4. **Ask for Help**
   - Join our community discussions
   - Ask questions in GitHub Discussions

### Development Tips

1. **Keep Your Branch Updated**
   ```bash
   git checkout main
   git pull upstream main
   git checkout your-branch
   git rebase main
   ```

2. **Use meaningful commit messages**
   ```bash
   # Good
   git commit -m "feat(instance): add interactive instance creation"
   
   # Bad
   git commit -m "fixed stuff"
   ```

3. **Write tests for your code**
   - Aim for high test coverage
   - Test both success and failure cases
   - Use table-driven tests for multiple scenarios

4. **Keep your code clean**
   - Follow Go best practices
   - Use meaningful variable names
   - Add comments for complex logic

5. **Use the tools available**
   - Use `go fmt` and `goimports` for formatting
   - Use `golangci-lint` for linting
   - Use `go test` for testing

---

Now you're ready to start contributing to Virak CLI! Happy coding!