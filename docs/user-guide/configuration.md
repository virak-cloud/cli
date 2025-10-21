# Configuration Guide

This guide explains how to configure Virak CLI for optimal use.

## Table of Contents

- [Configuration File](#configuration-file)
- [Configuration Options](#configuration-options)
- [Environment Variables](#environment-variables)
- [Command-Line Flags](#command-line-flags)
- [Advanced Configuration](#advanced-configuration)
- [Configuration Examples](#configuration-examples)

## Configuration File

Virak CLI uses a YAML configuration file located at:

- **Linux/macOS**: `~/.virak-cli.yaml`
- **Windows**: `%USERPROFILE%\.virak-cli.yaml`

The configuration file is automatically created when you first authenticate with `virak-cli login`.

### Basic Configuration Structure

```yaml
auth:
  token: "your-api-token-here"
default:
  zoneId: "zone-12345"
  zoneName: "tehran-1"
api:
  baseUrl: "https://public-api.virakcloud.com"
  timeout: "30s"
logging:
  level: "info"
  format: "text"
```

## Configuration Options

### Authentication Settings

```yaml
auth:
  token: "your-api-token-here"          # API authentication token
  expiresIn: "24h"                      # Token expiration (if applicable)
  refreshToken: "refresh-token"         # Refresh token (if applicable)
```

### Default Settings

```yaml
default:
  zoneId: "zone-12345"                  # Default zone ID
  zoneName: "tehran-1"                  # Default zone name
  outputFormat: "table"                 # Default output format (table, json, yaml)
```

### API Settings

```yaml
api:
  baseUrl: "https://public-api.virakcloud.com"  # API base URL
  timeout: "30s"                                # Request timeout
  retries: 3                                    # Number of retries
  retryDelay: "1s"                              # Delay between retries
```

### Logging Settings

```yaml
logging:
  level: "info"                         # Log level (debug, info, warn, error)
  format: "text"                        # Log format (text, json)
  file: ""                              # Log file path (empty for stdout)
```

### Proxy Settings

```yaml
proxy:
  http: "http://proxy.example.com:8080"  # HTTP proxy URL
  https: "https://proxy.example.com:8080" # HTTPS proxy URL
  noProxy: "localhost,127.0.0.1"         # No proxy for these hosts
```

## Environment Variables

Environment variables override configuration file settings:

### Authentication

```bash
export VIRAK_TOKEN="your-api-token-here"
export VIRAK_REFRESH_TOKEN="refresh-token"
```

### Default Zone

```bash
export VIRAK_ZONE_ID="zone-12345"
export VIRAK_ZONE_NAME="tehran-1"
```

### API Configuration

```bash
export VIRAK_API_URL="https://public-api.virakcloud.com"
export VIRAK_TIMEOUT="30s"
export VIRAK_RETRIES="3"
```

### Logging

```bash
export VIRAK_LOG_LEVEL="info"
export VIRAK_LOG_FORMAT="text"
export VIRAK_LOG_FILE="/var/log/virak-cli.log"
```

### Proxy

```bash
export HTTP_PROXY="http://proxy.example.com:8080"
export HTTPS_PROXY="https://proxy.example.com:8080"
export NO_PROXY="localhost,127.0.0.1"
```

## Command-Line Flags

Command-line flags override both environment variables and configuration file settings:

```bash
# Specify zone
virak-cli --zone-id zone-12345 instance list

# Specify output format
virak-cli --output json instance list

# Specify API URL
virak-cli --api-url https://api.virakcloud.com instance list

# Specify timeout
virak-cli --timeout 60s instance list

# Enable debug logging
virak-cli --debug instance list

# Use insecure connection (not recommended)
virak-cli --insecure instance list
```

## Advanced Configuration

### Multiple Profiles

You can use multiple configuration files for different environments:

```bash
# Use specific configuration file
virak-cli --config ~/.virak-cli-dev.yaml instance list
virak-cli --config ~/.virak-cli-prod.yaml instance list
```

### Conditional Configuration

Use shell scripts to manage multiple configurations:

```bash
#!/bin/bash
# virak-env.sh

if [ "$ENV" = "production" ]; then
    export VIRAK_ZONE_ID="zone-prod-123"
    export VIRAK_TOKEN="$PROD_TOKEN"
elif [ "$ENV" = "staging" ]; then
    export VIRAK_ZONE_ID="zone-staging-456"
    export VIRAK_TOKEN="$STAGING_TOKEN"
else
    export VIRAK_ZONE_ID="zone-dev-789"
    export VIRAK_TOKEN="$DEV_TOKEN"
fi
```

### Configuration Validation

Validate your configuration:

```bash
# Test configuration
virak-cli zone list

# Check current configuration
virak-cli --debug zone list 2>&1 | grep -E "(Using|Config)"

# Test API connectivity
curl -H "Authorization: Bearer $VIRAK_TOKEN" \
  "$VIRAK_API_URL/v1/zones"
```

## Configuration Examples

### Development Environment

```yaml
# ~/.virak-cli-dev.yaml
auth:
  token: "dev-token-here"
default:
  zoneId: "zone-dev-789"
  zoneName: "dev-zone"
  outputFormat: "json"
api:
  baseUrl: "https://dev-api.virakcloud.com"
  timeout: "60s"
logging:
  level: "debug"
  format: "json"
```

### Production Environment

```yaml
# ~/.virak-cli-prod.yaml
auth:
  token: "prod-token-here"
default:
  zoneId: "zone-prod-123"
  zoneName: "production"
  outputFormat: "table"
api:
  baseUrl: "https://public-api.virakcloud.com"
  timeout: "30s"
logging:
  level: "warn"
  format: "text"
```

### CI/CD Environment

```bash
#!/bin/bash
# ci-setup.sh

# Use environment variables for CI/CD
export VIRAK_TOKEN="${CI_VIRAK_TOKEN}"
export VIRAK_ZONE_ID="${CI_ZONE_ID}"
export VIRAK_LOG_LEVEL="info"
export VIRAK_OUTPUT_FORMAT="json"

# Validate configuration
if [ -z "$VIRAK_TOKEN" ]; then
    echo "Error: VIRAK_TOKEN not set"
    exit 1
fi

# Test connection
virak-cli zone list > /dev/null
if [ $? -ne 0 ]; then
    echo "Error: Failed to connect to Virak Cloud API"
    exit 1
fi

echo "Configuration validated successfully"
```

### Corporate Proxy Environment

```yaml
# ~/.virak-cli.yaml
auth:
  token: "corporate-token"
default:
  zoneId: "zone-corp-123"
api:
  baseUrl: "https://public-api.virakcloud.com"
  timeout: "60s"
proxy:
  http: "http://corporate-proxy:8080"
  https: "https://corporate-proxy:8080"
  noProxy: "localhost,127.0.0.1,.local"
logging:
  level: "info"
```

## Configuration Best Practices

### Security

1. **File Permissions**
   ```bash
   chmod 600 ~/.virak-cli.yaml
   ```

2. **Use Environment Variables for Sensitive Data**
   ```bash
   export VIRAK_TOKEN="sensitive-token"
   ```

3. **Avoid Hardcoding Tokens in Scripts**
   ```bash
   # Bad
   virak-cli --token "hardcoded-token" instance list
   
   # Good
   virak-cli --token "$VIRAK_TOKEN" instance list
   ```

### Performance

1. **Adjust Timeout for Large Operations**
   ```yaml
   api:
     timeout: "120s"  # For large deployments
   ```

2. **Use Appropriate Output Format**
   ```bash
   # For scripts
   virak-cli --output json instance list
   
   # For human reading
   virak-cli --output table instance list
   ```

### Reliability

1. **Configure Retries**
   ```yaml
   api:
     retries: 5
     retryDelay: "2s"
   ```

2. **Set Default Zone**
   ```yaml
   default:
     zoneId: "your-preferred-zone"
   ```

## Troubleshooting Configuration

### Common Issues

1. **Configuration File Not Found**
   ```bash
   # Create default configuration
   virak-cli login
   ```

2. **Invalid YAML Syntax**
   ```bash
   # Validate YAML
   python -c "import yaml; yaml.safe_load(open('~/.virak-cli.yaml'))"
   ```

3. **Permission Denied**
   ```bash
   # Fix permissions
   chmod 600 ~/.virak-cli.yaml
   chown $USER:$USER ~/.virak-cli.yaml
   ```

4. **Environment Variables Not Working**
   ```bash
   # Check if variables are set
   env | grep VIRAK_
   
   # Export variables in current shell
   export VIRAK_TOKEN="your-token"
   ```

### Debug Configuration

```bash
# Enable debug logging
virak-cli --debug zone list

# Check effective configuration
virak-cli --debug zone list 2>&1 | grep -E "(Using|Config|Token)"

# Test API connection
virak-cli --debug zone list 2>&1 | grep -E "(Request|Response|URL)"
```

For more troubleshooting help, see our [Troubleshooting Guide](../troubleshooting/common-issues.md).