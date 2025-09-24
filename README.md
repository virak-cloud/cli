# Virak CLI

A command-line interface for interacting with the [Virak Cloud API](https://api-docs.virakcloud.com/), built with the Go programming language.

[![VirakCloud](https://img.shields.io/badge/VirakCloud-Site-blue)](https://virakcloud.com)
[![Public API](https://img.shields.io/badge/Public_API-Docs-green)](https://api-docs.virakcloud.com/)
[![Service Document](https://img.shields.io/badge/Service_Document-Guides-orange)](https://docs.virakcloud.com/en/guides/)
[![Panel](https://img.shields.io/badge/Panel-Dashboard-red)](https://panel.virakcloud.com/)

## Table of Contents

- [About The Project](#about-the-project)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
    - [From Releases](#from-releases)
      - [Installing Binaries](#installing-binaries)
      - [Installing from Archives](#installing-from-archives)
      - [Installing on Debian/Ubuntu](#installing-on-debianubuntu)
      - [Installing on Windows](#installing-on-windows)
    - [From Source](#from-source)
- [Usage](#usage)
- [Authentication](#authentication)
- [Commands](#commands)
  - [Authentication](#authentication-1)
  - [Bucket (Object Storage)](#bucket-object-storage)
  - [DNS](#dns)
  - [Instance (VM)](#instance-vm)
  - [Kubernetes Clusters](#kubernetes-clusters)
  - [Network](#network)
  - [Zone](#zone)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Development](#development)
  - [Building](#building)
  - [Code Quality](#code-quality)
  - [Releasing](#releasing)
- [Contributing](#contributing)
- [License](#license)

## About The Project

The Virak CLI is a command-line interface that allows you to manage your Virak Cloud resources directly from your terminal.

## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

* Go 1.24.4 or later

### Installation

#### From Releases

Download the latest release from [GitHub Releases](https://github.com/virak-cloud/cli/releases).

Choose the appropriate binary for your platform:

- **Linux (amd64)**: `virak-cli-linux-amd64`
- **Linux (arm64)**: `virak-cli-linux-arm64`
- **macOS (amd64)**: `virak-cli-darwin-amd64`
- **macOS (arm64)**: `virak-cli-darwin-arm64`
- **Windows (amd64)**: `virak-cli-windows-amd64.exe`
- **Windows (arm64)**: `virak-cli-windows-arm64.exe`

##### Installing Binaries

1. Download the binary file for your platform from the releases page.
2. Make the binary executable (skip this step on Windows):
   ```sh
   chmod +x virak-cli-*
   ```
3. Move the binary to a directory in your PATH, for example:
   ```sh
   sudo mv virak-cli-* /usr/local/bin/virak-cli
   ```
4. Verify the installation:
   ```sh
   virak-cli --version
   ```

##### Installing from Archives

Alternatively, download the compressed archive (.tar.gz for Linux/macOS, .zip for Windows), extract it, and follow the binary installation steps above.

##### Installing on Debian/Ubuntu

For Debian-based systems, download the `.deb` package and install with:
```sh
sudo dpkg -i virak-cli_*.deb
```

##### Installing on Windows

1. Download the `.exe` file for your Windows architecture from the releases page.
2. Create a folder for the CLI, for example `C:\Program Files\virak-cli\`.
3. Move the downloaded `.exe` file to that folder.
4. Add the folder to your system's PATH:
   - Open System Properties (search for "environment variables" in Start menu).
   - Click "Environment Variables".
   - Under "System variables", find "Path" and click "Edit".
   - Click "New" and add `C:\Program Files\virak-cli\`.
   - Click "OK" to save.
5. Open a new Command Prompt or PowerShell window and verify:
   ```cmd
   virak-cli --version
   ```

Alternatively, you can use the Windows archive (.zip) - extract it and follow the same steps.

> **Note:** Releases are automated using [GoReleaser](https://goreleaser.com/), ensuring consistent and verified builds across platforms.

#### From Source

1. Clone the repo
   ```sh
   git clone https://github.com/virak-cloud/cli.git
   ```
2. Build the project
   ```sh
   go build -o virak-cli main.go
   ```

Alternatively, you can use `go install`:
```sh
go install github.com/virak-cloud/cli@latest
```

## Usage

Use the `virak-cli` command to interact with your Virak Cloud resources.

```sh
virak-cli [command]
```

### Authentication

Before using the CLI, you need to authenticate:

```sh
virak-cli login
```

This will open your browser for OAuth authentication. Alternatively, you can provide a token directly:

```sh
virak-cli login --token YOUR_TOKEN
```

## Commands

The following commands are available:

### Authentication
* `virak-cli login`: Authenticate with Virak Cloud
* `virak-cli logout`: Log out from Virak Cloud

### Bucket (Object Storage)
* `virak-cli bucket create`: Create a new bucket
* `virak-cli bucket delete`: Delete a bucket
* `virak-cli bucket events`: View bucket events
* `virak-cli bucket list`: List all buckets
* `virak-cli bucket show`: Show details of a bucket
* `virak-cli bucket update`: Update a bucket

### DNS
* `virak-cli dns domain create`: Create a new domain
* `virak-cli dns domain delete`: Delete a domain
* `virak-cli dns domain list`: List all domains
* `virak-cli dns domain show`: Show details of a domain
* `virak-cli dns events`: View DNS events
* `virak-cli dns record create`: Create a new DNS record
* `virak-cli dns record delete`: Delete a DNS record
* `virak-cli dns record list`: List all DNS records
* `virak-cli dns record update`: Update a DNS record

### Instance (VM)
* `virak-cli instance create`: Create a new instance
* `virak-cli instance delete`: Delete an instance
* `virak-cli instance list`: List all instances
* `virak-cli instance metrics`: View instance metrics
* `virak-cli instance reboot`: Reboot an instance
* `virak-cli instance rebuild`: Rebuild an instance
* `virak-cli instance service-offering list`: List instance service offerings
* `virak-cli instance show`: Show details of an instance
* `virak-cli instance snapshot create`: Create an instance snapshot
* `virak-cli instance snapshot delete`: Delete an instance snapshot
* `virak-cli instance snapshot list`: List instance snapshots
* `virak-cli instance snapshot revert`: Revert to an instance snapshot
* `virak-cli instance start`: Start an instance
* `virak-cli instance stop`: Stop an instance
* `virak-cli instance vm-image list`: List VM images
* `virak-cli instance volume attach`: Attach a volume to an instance
* `virak-cli instance volume create`: Create a volume for an instance
* `virak-cli instance volume delete`: Delete an instance volume
* `virak-cli instance volume detach`: Detach a volume from an instance
* `virak-cli instance volume list`: List instance volumes
* `virak-cli instance volume service-offering list`: List volume service offerings

### Kubernetes Clusters
* `virak-cli cluster create`: Create a new cluster
* `virak-cli cluster delete`: Delete a cluster
* `virak-cli cluster list`: List all clusters
* `virak-cli cluster scale`: Scale a cluster
* `virak-cli cluster show`: Show details of a cluster
* `virak-cli cluster start`: Start a cluster
* `virak-cli cluster stop`: Stop a cluster
* `virak-cli cluster update`: Update a cluster
* `virak-cli cluster update-service-offering`: Update service offering for a cluster
* `virak-cli cluster upgrade`: Upgrade a cluster
* `virak-cli cluster service events`: View Kubernetes service events
* `virak-cli cluster versions list`: List all Kubernetes versions

### Network
* `virak-cli network create`: Create a new network
* `virak-cli network create l2`: Create a new L2 network
* `virak-cli network create l3`: Create a new L3 network
* `virak-cli network delete`: Delete a network
* `virak-cli network firewall ipv4 create`: Create an IPv4 firewall rule
* `virak-cli network firewall ipv4 delete`: Delete an IPv4 firewall rule
* `virak-cli network firewall ipv4 list`: List all IPv4 firewall rules
* `virak-cli network firewall ipv6 create`: Create an IPv6 firewall rule
* `virak-cli network firewall ipv6 delete`: Delete an IPv6 firewall rule
* `virak-cli network firewall ipv6 list`: List all IPv6 firewall rules
* `virak-cli network instance connect`: Connect an instance to a network
* `virak-cli network instance disconnect`: Disconnect an instance from a network
* `virak-cli network instance list`: List all instances in a network
* `virak-cli network lb assign`: Assign a load balancer
* `virak-cli network lb create`: Create a new load balancer
* `virak-cli network lb deassign`: Deassign a load balancer
* `virak-cli network lb delete`: Delete a load balancer
* `virak-cli network lb haproxy live`: View live HAProxy stats
* `virak-cli network lb haproxy log`: View HAProxy logs
* `virak-cli network lb list`: List all load balancers
* `virak-cli network list`: List all networks
* `virak-cli network public-ip associate`: Associate a public IP
* `virak-cli network public-ip disassociate`: Disassociate a public IP
* `virak-cli network public-ip list`: List all public IPs
* `virak-cli network public-ip staticnat disable`: Disable static NAT
* `virak-cli network public-ip staticnat enable`: Enable static NAT
* `virak-cli network service-offering list`: List all network service offerings
* `virak-cli network show`: Show details of a network
* `virak-cli network vpn disable`: Disable VPN
* `virak-cli network vpn enable`: Enable VPN
* `virak-cli network vpn show`: Show VPN details
* `virak-cli network vpn update`: Update VPN

### Zone
* `virak-cli zone list`: List all zones
* `virak-cli zone networks`: Manage zone networks
* `virak-cli zone resources`: Manage zone resources
* `virak-cli zone services`: Manage zone services

## Project Structure

```
virak-cli/
├── cmd/                          # CLI command implementations
│   ├── bucket/                   # Bucket (Object Storage) commands
│   ├── cluster/                  # Kubernetes cluster commands
│   ├── dns/                      # DNS management commands
│   ├── instance/                 # VM instance commands
│   ├── network/                  # Network management commands
│   ├── zone/                     # Zone management commands
│   ├── login.go                  # Authentication login
│   ├── logout.go                 # Authentication logout
│   └── root.go                   # Root command
├── internal/                     # Internal packages
│   ├── cli/                      # CLI utilities and validation
│   ├── logger/                   # Logging utilities
│   └── presenter/                # Output formatting
├── pkg/                          # Reusable packages
│   ├── http/                     # HTTP client and API calls
│   ├── responses/                # API response structures
│   └── urls/                     # API endpoint URLs
├── main.go                       # Application entry point
├── go.mod                        # Go module file
├── go.sum                        # Go dependencies
├── LICENSE                       # License file
└── README.md                     # This file
```

## Configuration

The CLI can be configured via a configuration file located at `~/.virak-cli.yaml`.

The configuration file can be used to store your API token and other settings.

Example configuration:
```yaml
auth:
  token: "your-api-token"
default:
  zoneId: "your-default-zone-id"
  zoneName: "your-default-zone-name"
```

## Development

### Building

```sh
go build -o virak-cli main.go
```
x

### Code Quality

Run linters and formatters:
```sh
gofmt -w .
goimports -w .
golint ./...
go vet ./...
```

### Releasing

This project uses [GoReleaser](https://goreleaser.com/) for automated releases. Releases are triggered via GitHub Actions on tag pushes.

To create a new release:

1. Ensure all changes are committed and pushed to the `master` branch.
2. Tag the commit with a semantic version (e.g., `v1.2.3`):
   ```sh
   git tag v1.2.3
   git push origin v1.2.3
   ```
3. The [release workflow](.github/workflows/release.yml) will automatically:
   - Build binaries for Linux, Windows, and macOS (AMD64 and ARM64).
   - Create archives and Debian packages.
   - Generate a changelog from GitHub commits.
   - Publish the release to GitHub Releases.

The release process is configured in [.goreleaser.yml](.goreleaser.yml) and uses environment variables for API URLs.

## Contributing

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
