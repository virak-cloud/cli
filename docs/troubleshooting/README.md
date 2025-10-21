# Troubleshooting Guide

This guide provides solutions to common issues you might encounter while using Virak CLI.

## Table of Contents

- [Getting Help](#getting-help)
- [Common Issues](#common-issues)
  - [Authentication Issues](#authentication-issues)
  - [Configuration Issues](#configuration-issues)
  - [Network Issues](#network-issues)
  - [Instance Issues](#instance-issues)
  - [Bucket Issues](#bucket-issues)
  - [DNS Issues](#dns-issues)
  - [Network Issues](#network-issues-1)
  - [Zone Issues](#zone-issues)
  - [Build Issues](#build-issues)
- [Debug Mode](#debug-mode)
- [Logging](#logging)
- [Error Codes](#error-codes)
- [Reporting Issues](#reporting-issues)

## Getting Help

If you're experiencing issues with Virak CLI:

1. **Check this guide** - Look for your issue in the sections below
2. **Enable debug mode** - Use `--debug` flag for more detailed error messages
3. **Check logs** - Review application logs for additional context
4. **Search existing issues** - Check [GitHub Issues](https://github.com/virak-cloud/cli/issues) for similar problems
5. **Create a new issue** - If you can't find a solution, create a new issue with details

## Common Issues

### Authentication Issues

#### Issue: "Authentication failed" or "Invalid token"

**Symptoms:**
```
Error: Authentication failed (401)
Details: Invalid or expired token
```

**Causes:**
- Token has expired
- Token is invalid
- Token is not properly configured

**Solutions:**

1. **Re-authenticate:**
   ```bash
   virak-cli login
   ```

2. **Check token in config file:**
   ```bash
   # View config file
   cat ~/.virak-cli.yaml
   
   # Or use debug mode to check
   virak-cli --debug zone list
   ```

3. **Use token flag:**
   ```bash
   virak-cli --token YOUR_TOKEN zone list
   ```

#### Issue: "Permission denied"

**Symptoms:**
```
Error: Permission denied (403)
Details: Insufficient permissions for this resource
```

**Causes:**
- Token doesn't have required permissions
- Resource is not accessible with current token

**Solutions:**

1. **Check token permissions** in the Virak Cloud dashboard
2. **Re-authenticate with appropriate permissions:**
   ```bash
   virak-cli logout
   virak-cli login
   ```
3. **Contact support** if you believe you should have access

### Configuration Issues

#### Issue: "No default zone configured"

**Symptoms:**
```
Error: No default zone configured
```

**Causes:**
- Default zone is not set in configuration
- Configuration file is missing or corrupted

**Solutions:**

1. **Set default zone:**
   ```bash
   virak-cli zone list
   # Follow prompts to set default zone
   ```

2. **Use zone flag:**
   ```bash
   virak-cli --zoneId zone-123 instance list
   ```

3. **Check configuration file:**
   ```bash
   # View config file
   cat ~/.virak-cli.yaml
   ```

4. **Create configuration file:**
   ```bash
   # Create default config
   cat > ~/.virak-cli.yaml << EOF
   default:
     zoneId: "zone-123"
     zoneName: "Your Zone Name"
   EOF
   ```

#### Issue: "Configuration file not found"

**Symptoms:**
```
Error: Configuration file not found
```

**Causes:**
- Configuration file doesn't exist
- Configuration file is in wrong location

**Solutions:**

1. **Create configuration file:**
   ```bash
   # Create default config
   cat > ~/.virak-cli.yaml << EOF
   auth:
     token: "YOUR_TOKEN"
   default:
     zoneId: "zone-123"
     zoneName: "Your Zone Name"
   EOF
   ```

2. **Use environment variables:**
   ```bash
   export VIRAK_TOKEN="YOUR_TOKEN"
   export VIRAK_ZONE_ID="zone-123"
   ```

### Network Issues

#### Issue: "Connection timeout" or "Network unreachable"

**Symptoms:**
```
Error: Connection timeout
Error: Network unreachable
```

**Causes:**
- Network connectivity issues
- Firewall blocking connections
- API server is down

**Solutions:**

1. **Check internet connection:**
   ```bash
   ping api.virakcloud.com
   ```

2. **Check firewall settings**
3. **Use VPN or proxy if required**
4. **Try again later** - server might be temporarily unavailable

#### Issue: "SSL certificate verification failed"

**Symptoms:**
```
Error: SSL certificate verification failed
```

**Causes:**
- SSL certificate is invalid or expired
- System time is incorrect
- Corporate proxy interfering with SSL

**Solutions:**

1. **Check system time:**
   ```bash
   date
   ```

2. **Use insecure flag (not recommended for production):**
   ```bash
   virak-cli --insecure instance list
   ```

3. **Update certificate store:**
   - On Ubuntu/Debian: `sudo apt-get update && sudo apt-get install ca-certificates`
   - On CentOS/RHEL: `sudo yum update ca-certificates`
   - On macOS: Update through system updates

### Instance Issues

#### Issue: "Instance creation failed"

**Symptoms:**
```
Error: Instance creation failed
Details: Insufficient resources
```

**Causes:**
- Insufficient resources in zone
- Invalid parameters
- Service offering not available

**Solutions:**

1. **Check zone resources:**
   ```bash
   virak-cli zone resources --zoneId zone-123
   ```

2. **Check available service offerings:**
   ```bash
   virak-cli instance service-offering list --zoneId zone-123
   ```

3. **Check available VM images:**
   ```bash
   virak-cli instance vm-image list --zoneId zone-123
   ```

4. **Try different parameters:**
   ```bash
   virak-cli instance create \
     --name test-instance \
     --service-offering-id DIFFERENT_OFFERING_ID \
     --vm-image-id DIFFERENT_IMAGE_ID \
     --network-ids '["net-123"]'
   ```

#### Issue: "Instance not found"

**Symptoms:**
```
Error: Instance not found (404)
Details: Instance 'instance-123' does not exist
```

**Causes:**
- Instance ID is incorrect
- Instance is in different zone
- Instance was deleted

**Solutions:**

1. **List instances to verify:**
   ```bash
   virak-cli instance list
   ```

2. **Check instance ID:**
   ```bash
   virak-cli instance show --instanceId instance-123
   ```

3. **Check if you're in the right zone:**
   ```bash
   virak-cli --zoneId zone-123 instance list
   ```

#### Issue: "Instance stuck in creating state"

**Symptoms:**
```
Status: CREATING
```

**Causes:**
- Resource allocation issues
- Network configuration problems
- System overload

**Solutions:**

1. **Wait longer** - instance creation can take time
2. **Check zone resources:**
   ```bash
   virak-cli zone resources --zoneId zone-123
   ```
3. **Contact support** if instance stays in creating state for too long

### Bucket Issues

#### Issue: "Bucket name already exists"

**Symptoms:**
```
Error: Bucket name 'my-bucket' already exists
```

**Causes:**
- Bucket name is already in use
- Bucket name is reserved

**Solutions:**

1. **Choose a different name:**
   ```bash
   virak-cli bucket create --name my-unique-bucket --policy Private
   ```

2. **List existing buckets:**
   ```bash
   virak-cli bucket list
   ```

3. **Use more specific naming:**
   ```bash
   virak-cli bucket create --name my-project-bucket-$(date +%Y%m%d) --policy Private
   ```

#### Issue: "Cannot delete non-empty bucket"

**Symptoms:**
```
Error: Cannot delete non-empty bucket
```

**Causes:**
- Bucket contains objects
- Bucket is not empty

**Solutions:**

1. **Empty the bucket first** using S3-compatible tools:
   ```bash
   # Using AWS CLI
   aws s3 --endpoint-url=https://s3.virakcloud.com rm s3://my-bucket --recursive
   
   # Using s3cmd
   s3cmd --host=https://s3.virakcloud.com del s3://my-bucket --recursive
   ```

2. **Delete bucket after emptying:**
   ```bash
   virak-cli bucket delete --bucketId my-bucket
   ```

### DNS Issues

#### Issue: "Domain already exists"

**Symptoms:**
```
Error: Domain 'example.com' already exists
```

**Causes:**
- Domain is already registered
- Domain is in use by another account

**Solutions:**

1. **Choose a different domain:**
   ```bash
   virak-cli dns domain create --domain my-example.com
   ```

2. **Check existing domains:**
   ```bash
   virak-cli dns domain list
   ```

3. **Use subdomains:**
   ```bash
   virak-cli dns domain create --domain app.example.com
   ```

#### Issue: "Record creation failed"

**Symptoms:**
```
Error: Record creation failed
Details: Invalid record content
```

**Causes:**
- Invalid record content
- Incorrect record type
- Missing required fields

**Solutions:**

1. **Check record format:**
   ```bash
   # A record
   virak-cli dns record create --domain example.com --record www --type A --content 192.0.2.1
   
   # CNAME record
   virak-cli dns record create --domain example.com --record blog --type CNAME --content example.com
   
   # MX record (requires priority)
   virak-cli dns record create --domain example.com --record @ --type MX --content mail.example.com --priority 10
   ```

2. **Use appropriate record types:**
   - A: IPv4 address
   - AAAA: IPv6 address
   - CNAME: Domain name
   - MX: Mail server (requires priority)
   - TXT: Text content
   - NS: Name server
   - SOA: Start of authority

### Network Issues

#### Issue: "Network creation failed"

**Symptoms:**
```
Error: Network creation failed
Details: Invalid network offering
```

**Causes:**
- Network offering ID is invalid
- Network offering is not available
- Zone doesn't support network type

**Solutions:**

1. **Check available network offerings:**
   ```bash
   virak-cli network service-offering --type all
   ```

2. **Use correct offering ID:**
   ```bash
   # For L2 network
   virak-cli network service-offering --type l2
   
   # For L3 network
   virak-cli network service-offering --type l3
   ```

3. **Create network with correct offering:**
   ```bash
   virak-cli network create l2 \
     --name my-network \
     --network-offering-id VALID_OFFERING_ID
   ```

#### Issue: "Cannot connect instance to network"

**Symptoms:**
```
Error: Cannot connect instance to network
```

**Causes:**
- Instance is already connected to maximum networks
- Network is not available
- Instance is in wrong state

**Solutions:**

1. **Check instance status:**
   ```bash
   virak-cli instance show --instanceId instance-123
   ```

2. **Check network status:**
   ```bash
   virak-cli network show --networkId net-123
   ```

3. **Disconnect from other networks if needed:**
   ```bash
   virak-cli network instance disconnect --networkId net-123 --instanceId instance-123
   ```

### Zone Issues

#### Issue: "Zone not found"

**Symptoms:**
```
Error: Zone not found (404)
Details: Zone 'zone-xyz' does not exist
```

**Causes:**
- Zone ID is incorrect
- Zone is not available

**Solutions:**

1. **List available zones:**
   ```bash
   virak-cli zone list
   ```

2. **Use correct zone ID:**
   ```bash
   virak-cli --zoneId zone-123 instance list
   ```

3. **Set default zone:**
   ```bash
   virak-cli zone list
   # Follow prompts to set default zone
   ```

#### Issue: "Insufficient resources in zone"

**Symptoms:**
```
Error: Insufficient resources in zone
```

**Causes:**
- Zone has reached resource limits
- No available resources

**Solutions:**

1. **Check zone resources:**
   ```bash
   virak-cli zone resources --zoneId zone-123
   ```

2. **Try a different zone:**
   ```bash
   virak-cli zone list
   # Choose a different zone with available resources
   ```

3. **Wait for resources to become available**
4. **Contact support** to increase resource limits

### Build Issues

#### Issue: "Build failed with import errors"

**Symptoms:**
```
import "github.com/virak-cloud/cli/pkg/http": cannot find package
```

**Causes:**
- Dependencies not properly installed
- Go module issues

**Solutions:**

1. **Clean and reinstall dependencies:**
   ```bash
   go clean -modcache
   go mod download
   go mod tidy
   ```

2. **Check Go version:**
   ```bash
   go version
   # Ensure you're using Go 1.24.4 or later
   ```

3. **Update Go modules:**
   ```bash
   go mod download
   go mod tidy
   ```

#### Issue: "Tests fail with connection errors"

**Symptoms:**
```
Error: dial tcp: connection refused
```

**Causes:**
- Tests trying to connect to API
- No test environment configured

**Solutions:**

1. **Set up test environment variables:**
   ```bash
   export VIRAK_API_URL=https://api-test.virakcloud.com
   export VIRAK_TOKEN=test-token
   ```

2. **Use mock tests for unit testing:**
   ```go
   // Example mock test
   func TestInstanceService_CreateInstance(t *testing.T) {
       client := &mockHTTPClient{}
       service := NewInstanceService(client)
       
       client.On("CreateInstance", mock.Anything, mock.Anything).
           Return(&responses.Instance{}, nil)
       
       result, err := service.CreateInstance(context.Background(), "zone-123", CreateOptions{})
       
       assert.NoError(t, err)
       assert.NotNil(t, result)
   }
   ```

3. **Skip integration tests in CI:**
   ```bash
   go test -short ./...
   ```

## Debug Mode

Enable debug mode for more detailed error messages:

```bash
# Enable debug for a single command
virak-cli --debug instance list

# Enable debug for all commands
export VIRAK_DEBUG=true
virak-cli instance list
```

Debug mode provides:
- Detailed HTTP request/response information
- Stack traces for errors
- Configuration details
- Timing information

## Logging

Virak CLI uses structured logging with slog:

### Enable Logging

```bash
# Enable debug logging
virak-cli --debug instance list

# Set log level
export VIRAK_LOG_LEVEL=debug
virak-cli instance list
```

### Log Levels

- `debug`: Detailed debugging information
- `info`: General information (default)
- `warn`: Warning messages
- `error`: Error messages only

### Log Output

Logs are written to stderr by default. You can redirect them:

```bash
# Save logs to file
virak-cli --debug instance list 2> debug.log

# View logs in real-time
virak-cli --debug instance list 2>&1 | tee debug.log
```

## Error Codes

| Error Code | Description | Solution |
|------------|-------------|----------|
| `400` | Bad Request | Check command parameters |
| `401` | Unauthorized | Re-authenticate with `virak-cli login` |
| `403` | Forbidden | Check permissions |
| `404` | Not Found | Verify resource ID exists |
| `409` | Conflict | Resource already exists |
| `422` | Unprocessable Entity | Check command parameters |
| `429` | Too Many Requests | Wait and retry |
| `500` | Internal Server Error | Try again later |
| `502` | Bad Gateway | Try again later |
| `503` | Service Unavailable | Try again later |

## Reporting Issues

If you can't resolve your issue with the solutions above:

1. **Gather Information**
   - Error message
   - Command used
   - Debug output (if applicable)
   - System information (OS, Go version, CLI version)

2. **Search Existing Issues**
   - Check [GitHub Issues](https://github.com/virak-cloud/cli/issues)
   - Search for your error message

3. **Create a New Issue**
   - Go to [GitHub Issues](https://github.com/virak-cloud/cli/issues/new)
   - Use the issue template
   - Provide detailed information
   - Include steps to reproduce

4. **Include Debug Output**
   ```bash
   virak-cli --debug your-command > debug.log 2>&1
   # Attach debug.log to your issue
   ```

### Issue Template

```markdown
## Description
Brief description of the issue.

## Steps to Reproduce
1. Run command: `virak-cli ...`
2. Expected result: ...
3. Actual result: ...

## Environment
- OS: [e.g., Ubuntu 22.04]
- Go version: [e.g., 1.24.4]
- CLI version: [e.g., 1.2.3]

## Error Message
```
Paste error message here
```

## Debug Output
```
Paste debug output here (optional)
```

## Additional Context
Any additional information that might be helpful.
```

---

For more help, see our [Developer Documentation](../developer/README.md) or [User Guide](../user-guide/README.md).