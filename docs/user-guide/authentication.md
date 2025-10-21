# Authentication Guide

This guide explains how to authenticate with Virak Cloud using the Virak CLI.

## Table of Contents

- [Overview](#overview)
- [Authentication Methods](#authentication-methods)
- [OAuth Authentication](#oauth-authentication)
- [Token Authentication](#token-authentication)
- [Configuration File](#configuration-file)
- [Managing Authentication](#managing-authentication)
- [Security Best Practices](#security-best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

Virak CLI uses token-based authentication to communicate with the Virak Cloud API. You can authenticate using one of two methods:

1. **OAuth Authentication** - Browser-based authentication flow
2. **Token Authentication** - Direct token provision

Both methods result in a token being stored in your configuration file for subsequent API calls.

## Authentication Methods

### OAuth Authentication (Recommended)

OAuth authentication is the most secure method as it doesn't require you to handle tokens directly.

```bash
virak-cli login
```

This command will:
1. Open your default browser
2. Redirect you to the Virak Cloud login page
3. Prompt you to authorize the CLI application
4. Automatically retrieve and store the token

### Token Authentication

If you already have a token or prefer to provide it directly:

```bash
virak-cli login --token YOUR_TOKEN_HERE
```

You can obtain a token from the [Virak Cloud Panel](https://panel.virakcloud.com/) or by contacting support.

## Configuration File

Authentication information is stored in `~/.virak-cli.yaml` (or `%USERPROFILE%\.virak-cli.yaml` on Windows).

### Example Configuration

```yaml
auth:
  token: "your-api-token-here"
default:
  zoneId: "zone-12345"
  zoneName: "tehran-1"
```

### Configuration Locations

The CLI looks for configuration files in the following order:

1. `./.virak-cli.yaml` (current directory)
2. `$HOME/.virak-cli.yaml` (user home directory)
3. Environment variables

### Environment Variables

You can use environment variables instead of a configuration file:

```bash
export VIRAK_TOKEN="your-token-here"
export VIRAK_ZONE_ID="zone-12345"
export VIRAK_ZONE_NAME="tehran-1"
```

## Managing Authentication

### Checking Authentication Status

```bash
# Test authentication by listing zones
virak-cli zone list

# If successful, you're authenticated
# If you get an error, you need to authenticate
```

### Logging Out

To remove stored credentials:

```bash
virak-cli logout
```

This will:
1. Remove the token from your configuration file
2. Revoke the token on the server (if supported)
3. Clear any cached authentication data

### Refreshing Tokens

If your token expires or becomes invalid:

```bash
# Re-authenticate
virak-cli login

# Or provide a new token
virak-cli login --token NEW_TOKEN_HERE
```

## Security Best Practices

### Token Security

1. **Never share your token** - Treat it like a password
2. **Use environment variables** in CI/CD environments instead of configuration files
3. **Rotate tokens regularly** - Generate new tokens periodically
4. **Use minimal permissions** - Create tokens with only the permissions you need

### Configuration File Security

```bash
# Set appropriate permissions on your config file
chmod 600 ~/.virak-cli.yaml  # Linux/macOS
```

### Token Scopes

When creating tokens in the Virak Cloud Panel, consider the following scopes:

- **Read-only** - For monitoring and listing resources
- **Read-write** - For full resource management
- **Admin** - For administrative operations

### Multi-Factor Authentication

If your Virak Cloud account has MFA enabled:
1. OAuth authentication will handle MFA prompts automatically
2. Token authentication requires you to create tokens with MFA already enabled

## Advanced Authentication

### Service Account Authentication

For automated workflows and CI/CD:

1. Create a service account in the Virak Cloud Panel
2. Generate a service account token
3. Use the token in your automation scripts:

```bash
# In CI/CD pipeline
export VIRAK_TOKEN="service-account-token"
virak-cli instance list
```

### Temporary Tokens

For temporary access:

```bash
# Create a temporary token (if supported)
virak-cli login --temporary --duration 1h

# Token will automatically expire after specified duration
```

### Multiple Accounts

To manage multiple Virak Cloud accounts:

```bash
# Method 1: Use environment variables
export VIRAK_TOKEN="account1-token"
virak-cli --config account1.yaml zone list

export VIRAK_TOKEN="account2-token"
virak-cli --config account2.yaml zone list

# Method 2: Use different config files
virak-cli --config ~/.virak-cli-account1.yaml zone list
virak-cli --config ~/.virak-cli-account2.yaml zone list
```

## Troubleshooting

### "Authentication Failed" Error

**Causes:**
- Invalid or expired token
- Network connectivity issues
- Incorrect API endpoint

**Solutions:**
1. Re-authenticate: `virak-cli login`
2. Check network connectivity
3. Verify API endpoint in configuration

### "Token Expired" Error

**Solutions:**
1. Re-authenticate: `virak-cli login`
2. Use a new token: `virak-cli login --token NEW_TOKEN`

### "Permission Denied" Error

**Causes:**
- Token lacks required permissions
- Resource not accessible to your account

**Solutions:**
1. Check token permissions in Virak Cloud Panel
2. Ensure you're accessing the correct zone
3. Contact administrator if needed

### Browser Not Opening (OAuth)

**Solutions:**
1. Use token authentication instead: `virak-cli login --token TOKEN`
2. Manually visit the OAuth URL (shown in error message)
3. Check if browser is installed and accessible

### Configuration File Issues

**Debug configuration:**

```bash
# Check current configuration
cat ~/.virak-cli.yaml

# Test with a clean configuration
mv ~/.virak-cli.yaml ~/.virak-cli.yaml.backup
virak-cli login
```

### Proxy Issues

If you're behind a proxy:

```bash
# Set proxy environment variables
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=https://proxy.example.com:8080

# Or configure in ~/.virak-cli.yaml
proxy:
  http: "http://proxy.example.com:8080"
  https: "https://proxy.example.com:8080"
```

### SSL Certificate Issues

If you encounter SSL certificate errors:

```bash
# For testing only (not recommended for production)
virak-cli --insecure zone list

# Or configure custom CA bundle
export SSL_CERT_FILE=/path/to/ca-bundle.crt
```

## Examples

### Basic Authentication Workflow

```bash
# First-time setup
virak-cli login

# Verify authentication
virak-cli zone list

# Use the CLI
virak-cli instance list
```

### CI/CD Authentication

```bash
# In your CI/CD pipeline
#!/bin/bash

# Set token from environment variable
export VIRAK_TOKEN="${VIRAK_API_TOKEN}"

# Use CLI in your scripts
virak-cli instance create --name "ci-instance" --vm-image-id "${VM_IMAGE}" --service-offering-id "${SERVICE_OFFERING}"

# Clean up
virak-cli instance delete --id "${INSTANCE_ID}"
```

### Multi-Environment Setup

```bash
# Development
export VIRAK_TOKEN="dev-token"
export VIRAK_ZONE_ID="dev-zone"
virak-cli --config dev.yaml instance list

# Production
export VIRAK_TOKEN="prod-token"
export VIRAK_ZONE_ID="prod-zone"
virak-cli --config prod.yaml instance list
```

For more troubleshooting help, see our [Authentication Troubleshooting Guide](../troubleshooting/authentication.md).