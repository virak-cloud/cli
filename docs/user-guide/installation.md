# Installation Guide

This guide provides detailed instructions for installing Virak CLI on various platforms.

## Table of Contents

- [System Requirements](#system-requirements)
- [Installation Methods](#installation-methods)
  - [From Releases](#from-releases)
  - [From Source](#from-source)
  - [Using Go Install](#using-go-install)
- [Platform-Specific Instructions](#platform-specific-instructions)
  - [Linux](#linux)
  - [macOS](#macos)
  - [Windows](#windows)
- [Verification](#verification)
- [Uninstallation](#uninstallation)
- [Troubleshooting](#troubleshooting)

## System Requirements

- **Go 1.24.4 or later** (required only when building from source)
- **Operating Systems**: Linux, macOS, Windows
- **Architectures**: amd64, arm64

## Installation Methods

### From Releases (Recommended)

The easiest way to install Virak CLI is by downloading the pre-compiled binary from [GitHub Releases](https://github.com/virak-cloud/cli/releases).

#### Available Binaries

- **Linux (amd64)**: `virak-cli-linux-amd64`
- **Linux (arm64)**: `virak-cli-linux-arm64`
- **macOS (amd64)**: `virak-cli-darwin-amd64`
- **macOS (arm64)**: `virak-cli-darwin-arm64`
- **Windows (amd64)**: `virak-cli-windows-amd64.exe`
- **Windows (arm64)**: `virak-cli-windows-arm64.exe`

### From Source

If you prefer to build from source or need to modify the CLI:

```bash
# Clone the repository
git clone https://github.com/virak-cloud/cli.git
cd cli

# Build the binary
go build -o virak-cli main.go

# Move to your PATH
sudo mv virak-cli /usr/local/bin/  # Linux/macOS
# or add to PATH on Windows
```

### Using Go Install

If you have Go installed, you can use `go install`:

```bash
go install github.com/virak-cloud/cli@latest
```

Note: This will install the binary to `$GOPATH/bin` (usually `~/go/bin`), which should be in your PATH.

## Platform-Specific Instructions

### Linux

#### Method 1: Binary Installation

```bash
# Download the latest binary
curl -L https://github.com/virak-cloud/cli/releases/latest/download/virak-cli-linux-amd64 -o virak-cli

# Make it executable
chmod +x virak-cli

# Move to a directory in your PATH
sudo mv virak-cli /usr/local/bin/

# Verify installation
virak-cli --version
```

#### Method 2: Package Installation (Debian/Ubuntu)

```bash
# Download the .deb package
wget https://github.com/virak-cloud/cli/releases/latest/download/virak-cli_latest_amd64.deb

# Install the package
sudo dpkg -i virak-cli_latest_amd64.deb

# Verify installation
virak-cli --version
```

#### Method 3: Using Package Manager (if available)

```bash
# Using apt (if repository is configured)
sudo apt update
sudo apt install virak-cli

# Using yay (AUR)
yay -S virak-cli-bin
```

### macOS

#### Method 1: Binary Installation

```bash
# Download the latest binary
curl -L https://github.com/virak-cloud/cli/releases/latest/download/virak-cli-darwin-amd64 -o virak-cli

# For Apple Silicon (M1/M2)
# curl -L https://github.com/virak-cloud/cli/releases/latest/download/virak-cli-darwin-arm64 -o virak-cli

# Make it executable
chmod +x virak-cli

# Move to a directory in your PATH
sudo mv virak-cli /usr/local/bin/

# Verify installation
virak-cli --version
```

#### Method 2: Using Homebrew

```bash
# Add tap (if not already added)
brew tap virak-cloud/cli

# Install
brew install virak-cli

# Verify installation
virak-cli --version
```

#### Method 3: Using MacPorts

```bash
sudo port selfupdate
sudo port install virak-cli
```

### Windows

#### Method 1: Binary Installation

1. Download the appropriate binary from [GitHub Releases](https://github.com/virak-cloud/cli/releases):
   - `virak-cli-windows-amd64.exe` for Intel/AMD processors
   - `virak-cli-windows-arm64.exe` for ARM processors

2. Create a folder for the CLI, for example `C:\Program Files\virak-cli\`

3. Move the downloaded `.exe` file to that folder

4. Add the folder to your system's PATH:
   - Open System Properties (search for "environment variables" in Start menu)
   - Click "Environment Variables"
   - Under "System variables", find "Path" and click "Edit"
   - Click "New" and add `C:\Program Files\virak-cli\`
   - Click "OK" to save

5. Open a new Command Prompt or PowerShell window and verify:
   ```cmd
   virak-cli --version
   ```

#### Method 2: Using Chocolatey

```powershell
# Install Chocolatey (if not already installed)
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install Virak CLI
choco install virak-cli

# Verify installation
virak-cli --version
```

#### Method 3: Using Scoop

```powershell
# Install Scoop (if not already installed)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# Add bucket and install
scoop bucket add virak-cloud https://github.com/virak-cloud/scoop-bucket.git
scoop install virak-cli

# Verify installation
virak-cli --version
```

## Verification

After installation, verify that Virak CLI is working correctly:

```bash
# Check version
virak-cli --version

# Check help
virak-cli --help

# Test connection (requires authentication)
virak-cli zone list
```

If you see version information and help text, the installation was successful.

## Uninstallation

### Linux/macOS

```bash
# Remove binary
sudo rm /usr/local/bin/virak-cli

# Remove configuration (optional)
rm -rf ~/.virak-cli.yaml
```

### Windows

```cmd
# Remove binary
del "C:\Program Files\virak-cli\virak-cli.exe"

# Remove from PATH (manual step)
# Remove C:\Program Files\virak-cli from your PATH environment variable

# Remove configuration (optional)
del %USERPROFILE%\.virak-cli.yaml
```

## Troubleshooting

### "command not found" Error

If you get a "command not found" error after installation:

1. **Verify the binary is in your PATH**:
   ```bash
   which virak-cli  # Linux/macOS
   where virak-cli  # Windows
   ```

2. **Add to PATH if needed**:
   - Linux/macOS: `export PATH=$PATH:/path/to/virak-cli`
   - Windows: Add the directory to your PATH environment variable

3. **Restart your terminal** to apply PATH changes

### Permission Denied

If you get a permission denied error:

```bash
# Make the binary executable (Linux/macOS)
chmod +x virak-cli
```

### SSL Certificate Issues

If you encounter SSL certificate errors:

1. **Update your system certificates**
2. **Use the insecure flag (not recommended for production)**:
   ```bash
   virak-cli --insecure zone list
   ```

### Proxy Issues

If you're behind a proxy:

```bash
# Set proxy environment variables
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=https://proxy.example.com:8080

# Or configure in ~/.virak-cli.yaml
```

### Version Conflicts

If you have multiple versions installed:

```bash
# Find all instances
which -a virak-cli  # Linux/macOS
where virak-cli     # Windows

# Remove old versions
```

For more troubleshooting help, see our [Troubleshooting Guide](../troubleshooting/common-issues.md).