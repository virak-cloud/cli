# User Guide

Welcome to the Virak CLI User Guide! This section contains everything you need to effectively use the Virak CLI to manage your Virak Cloud resources.

## Table of Contents

- [Getting Started](#getting-started)
- [Core Concepts](#core-concepts)
- [Common Workflows](#common-workflows)
- [Additional Resources](#additional-resources)

## Getting Started

If you're new to Virak CLI, we recommend following these guides in order:

1. [Installation Guide](installation.md) - Install Virak CLI on your system
2. [Authentication](authentication.md) - Set up authentication with Virak Cloud
3. [Getting Started Tutorial](getting-started.md) - Learn the basics with hands-on examples

## Core Concepts

### Zones

Virak Cloud resources are organized into zones, which represent different geographical locations or data centers. Most commands require a zone ID, which you can obtain using:

```bash
virak-cli zone list
```

You can set a default zone in your configuration file (`~/.virak-cli.yaml`) to avoid specifying it for every command:

```yaml
default:
  zoneId: "your-default-zone-id"
  zoneName: "your-default-zone-name"
```

### Authentication

Virak CLI uses token-based authentication. You can authenticate either via OAuth (browser-based) or by providing a token directly. See the [Authentication Guide](authentication.md) for detailed instructions.

### Resource Management

Most resources in Virak Cloud follow a standard lifecycle:

1. **Create** - Create a new resource
2. **List** - List all resources of a type
3. **Show** - Get detailed information about a specific resource
4. **Update** - Modify an existing resource
5. **Delete** - Remove a resource

## Common Workflows

### Basic Instance Management

```bash
# List available VM images
virak-cli instance vm-image list

# List service offerings
virak-cli instance service-offering list

# Create an instance
virak-cli instance create --name my-instance --vm-image-id <image-id> --service-offering-id <offering-id>

# Start the instance
virak-cli instance start --id <instance-id>

# Check instance status
virak-cli instance show --id <instance-id>
```

### Network Configuration

```bash
# Create a network
virak-cli network create l3 --name my-network --display-text "My Network"

# List networks
virak-cli network list

# Connect an instance to a network
virak-cli network instance connect --instance-id <instance-id> --network-id <network-id>

# Create firewall rules
virak-cli network firewall ipv4 create --network-id <network-id> --protocol tcp --start-port 80 --end-port 80
```

### Storage Management

```bash
# Create a bucket
virak-cli bucket create --name my-bucket --policy Private

# List buckets
virak-cli bucket list

# Upload files (using your preferred S3-compatible tool)
s3cmd --host=https://s3.virakcloud.com mb s3://my-bucket
```

## Additional Resources

- [Command Reference](../reference/) - Detailed documentation for all commands
- [Tutorials](tutorials/) - Step-by-step guides for common tasks
- [Troubleshooting](../troubleshooting/) - Solutions to common issues
- [Configuration Guide](configuration.md) - Advanced configuration options

## Getting Help

If you need help with Virak CLI:

1. Check the [troubleshooting guides](../troubleshooting/)
2. Use the built-in help: `virak-cli --help` or `virak-cli <command> --help`
3. Search the [GitHub Issues](https://github.com/virak-cloud/cli/issues)
4. Create a new issue if you can't find a solution