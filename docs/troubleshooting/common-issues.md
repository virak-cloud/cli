# Common Issues and Solutions

This guide provides solutions to the most common issues encountered while using Virak CLI.

## Table of Contents

- [Authentication Issues](#authentication-issues)
- [Configuration Issues](#configuration-issues)
- [Network Connectivity Issues](#network-connectivity-issues)
- [Resource Management Issues](#resource-management-issues)
- [Command Execution Issues](#command-execution-issues)
- [Output Formatting Issues](#output-formatting-issues)
- [Performance Issues](#performance-issues)
- [Zone-Specific Issues](#zone-specific-issues)

## Authentication Issues

### Issue: "Authentication failed" or "Invalid token"

**Symptoms:**
```
Error: Authentication failed (401)
Details: Invalid or expired token
```

**Diagnosis:**
1. Check if token is expired
2. Verify token is correctly configured
3. Test API connectivity

**Solutions:**

#### Solution 1: Re-authenticate
```bash
# Logout first
virak-cli logout

# Login again
virak-cli login
```

#### Solution 2: Check Token Configuration
```bash
# Check configuration file
cat ~/.virak-cli.yaml

# Check token with debug mode
virak-cli --debug zone list
```

#### Solution 3: Use Token Flag
```bash
# Use token flag for single command
virak-cli --token YOUR_TOKEN zone list

# Set token environment variable
export VIRAK_TOKEN=YOUR_TOKEN
virak-cli zone list
```

### Issue: "Permission denied"

**Symptoms:**
```
Error: Permission denied (403)
Details: Insufficient permissions for this resource
```

**Diagnosis:**
1. Check token permissions in Virak Cloud dashboard
2. Verify resource is accessible
3. Test with different resources

**Solutions:**

#### Solution 1: Check Token Permissions
1. Log in to Virak Cloud dashboard
2. Navigate to API tokens
3. Verify token has required permissions

#### Solution 2: Re-authenticate with Different Account
```bash
# Logout and login with different account
virak-cli logout
virak-cli login
```

#### Solution 3: Contact Support
If you believe you should have access, contact Virak Cloud support.

## Configuration Issues

### Issue: "No default zone configured"

**Symptoms:**
```
Error: No default zone configured
```

**Diagnosis:**
1. Check if default zone is set
2. Verify configuration file exists
3. Test with explicit zone ID

**Solutions:**

#### Solution 1: Set Default Zone
```bash
# List available zones
virak-cli zone list

# Follow prompts to set default zone
```

#### Solution 2: Use Zone Flag
```bash
# Use zone flag for single command
virak-cli --zoneId zone-123 instance list
```

#### Solution 3: Create Configuration File
```bash
# Create configuration file
cat > ~/.virak-cli.yaml << EOF
auth:
  token: "YOUR_TOKEN"
default:
  zoneId: "zone-123"
  zoneName: "Your Zone Name"
EOF
```

### Issue: "Configuration file not found"

**Symptoms:**
```
Error: Configuration file not found
```

**Diagnosis:**
1. Check if configuration file exists
2. Verify file location
3. Test with environment variables

**Solutions:**

#### Solution 1: Create Configuration File
```bash
# Create default configuration
mkdir -p ~/.virak-cli
cat > ~/.virak-cli.yaml << EOF
auth:
  token: "YOUR_TOKEN"
default:
  zoneId: "zone-123"
  zoneName: "Your Zone Name"
EOF
```

#### Solution 2: Use Environment Variables
```bash
# Set environment variables
export VIRAK_TOKEN="YOUR_TOKEN"
export VIRAK_ZONE_ID="zone-123"

# Test with environment variables
virak-cli instance list
```

#### Solution 3: Use Custom Config File
```bash
# Use custom config file
export VIRAK_CONFIG="/path/to/config.yaml"
virak-cli instance list
```

## Network Connectivity Issues

### Issue: "Connection timeout" or "Network unreachable"

**Symptoms:**
```
Error: Connection timeout
Error: Network unreachable
Error: dial tcp: lookup api.virakcloud.com: no such host
```

**Diagnosis:**
1. Check internet connection
2. Verify DNS resolution
3. Test API endpoint accessibility

**Solutions:**

#### Solution 1: Check Internet Connection
```bash
# Check basic connectivity
ping 8.8.8.8

# Check DNS resolution
nslookup api.virakcloud.com

# Check API endpoint
curl -I https://api.virakcloud.com
```

#### Solution 2: Check Firewall Settings
```bash
# Check if port 443 is blocked
telnet api.virakcloud.com 443

# Or use nc
nc -zv api.virakcloud.com 443
```

#### Solution 3: Use Proxy or VPN
If you're behind a corporate firewall:
```bash
# Use HTTP proxy
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=https://proxy.example.com:8080

# Use SOCKS proxy
export ALL_PROXY=socks5://proxy.example.com:1080
```

#### Solution 4: Use Insecure Flag (Not Recommended)
```bash
# Skip SSL verification (not recommended for production)
virak-cli --insecure instance list
```

### Issue: "SSL certificate verification failed"

**Symptoms:**
```
Error: SSL certificate verification failed
Error: x509: certificate signed by unknown authority
```

**Diagnosis:**
1. Check system time
2. Verify certificate store
3. Test with different certificate settings

**Solutions:**

#### Solution 1: Check System Time
```bash
# Check system time
date

# Sync time if needed
sudo ntpdate pool.ntp.org
# or
sudo timedatectl set-ntp true
```

#### Solution 2: Update Certificate Store
```bash
# Ubuntu/Debian
sudo apt-get update && sudo apt-get install ca-certificates

# CentOS/RHEL
sudo yum update ca-certificates

# macOS
# Update through system preferences or software update
```

#### Solution 3: Use Custom CA Certificate
```bash
# Add custom CA to system store
sudo cp custom-ca.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```

## Resource Management Issues

### Issue: "Instance creation failed"

**Symptoms:**
```
Error: Instance creation failed
Details: Insufficient resources
Error: Service offering not available
Error: VM image not found
```

**Diagnosis:**
1. Check zone resources
2. Verify service offerings
3. Confirm VM image availability

**Solutions:**

#### Solution 1: Check Zone Resources
```bash
# Check available resources
virak-cli zone resources --zoneId zone-123

# Look for:
# - Memory: Available/Total
# - CPU Number: Available/Total
# - Data Volume: Available/Total
# - VM Limit: Available/Total
```

#### Solution 2: Check Available Offerings
```bash
# List available service offerings
virak-cli instance service-offering list --zoneId zone-123

# Look for offerings with IsAvailable: true
```

#### Solution 3: Check Available VM Images
```bash
# List available VM images
virak-cli instance vm-image list --zoneId zone-123
```

#### Solution 4: Try Different Parameters
```bash
# Create instance with different parameters
virak-cli instance create \
  --name test-instance \
  --service-offering-id DIFFERENT_OFFERING_ID \
  --vm-image-id DIFFERENT_IMAGE_ID \
  --network-ids '["net-123"]'
```

### Issue: "Instance stuck in creating state"

**Symptoms:**
```
Status: CREATING
```

**Diagnosis:**
1. Check how long it's been in creating state
2. Verify zone resources
3. Check for system issues

**Solutions:**

#### Solution 1: Wait Longer
Instance creation can take time, especially for larger instances:
```bash
# Monitor instance status
watch virak-cli instance show --instanceId instance-123
```

#### Solution 2: Check Zone Resources
```bash
# Check if resources are available
virak-cli zone resources --zoneId zone-123
```

#### Solution 3: Contact Support
If instance stays in creating state for more than 30 minutes:
1. Collect instance details
2. Contact Virak Cloud support
3. Provide instance ID and zone ID

## Command Execution Issues

### Issue: "Command not found"

**Symptoms:**
```
bash: virak-cli: command not found
```

**Diagnosis:**
1. Check if CLI is installed
2. Verify PATH configuration
3. Test binary execution

**Solutions:**

#### Solution 1: Install CLI
```bash
# Clone and build
git clone https://github.com/virak-cloud/cli.git
cd cli
go build -o virak-cli main.go

# Move to PATH
sudo mv virak-cli /usr/local/bin/
```

#### Solution 2: Update PATH
```bash
# Add to PATH (temporary)
export PATH=$PATH:/path/to/virak-cli

# Add to PATH (permanent)
echo 'export PATH=$PATH:/path/to/virak-cli' >> ~/.bashrc
source ~/.bashrc
```

#### Solution 3: Use Full Path
```bash
# Use full path to binary
/path/to/virak-cli instance list
```

### Issue: "Unknown command" or "Invalid flag"

**Symptoms:**
```
Error: unknown command "invalid-command"
Error: unknown flag --invalid-flag
```

**Diagnosis:**
1. Check command syntax
2. Verify command exists
3. Check flag names

**Solutions:**

#### Solution 1: Check Available Commands
```bash
# List all commands
virak-cli --help

# List subcommands
virak-cli instance --help
```

#### Solution 2: Check Command Syntax
```bash
# Check command help
virak-cli instance create --help

# Look for correct flag names
virak-cli instance create --help | grep "name"
```

#### Solution 3: Use Command Completion
```bash
# Enable bash completion
source <(virak-cli completion bash)

# Or install permanently
virak-cli completion bash > /etc/bash_completion.d/virak-cli
```

## Output Formatting Issues

### Issue: "Invalid output format"

**Symptoms:**
```
Error: invalid output format 'invalid-format'
```

**Diagnosis:**
1. Check available output formats
2. Verify format syntax
3. Test with different formats

**Solutions:**

#### Solution 1: Use Valid Output Formats
```bash
# Available formats: table, json, yaml
virak-cli --output json instance list
virak-cli --output yaml instance list
virak-cli --output table instance list
```

#### Solution 2: Check Format Help
```bash
# Check output format options
virak-cli --help | grep -A 5 "output"
```

#### Solution 3: Use Environment Variable
```bash
# Set default output format
export VIRAK_OUTPUT_FORMAT=json
virak-cli instance list
```

### Issue: "Table formatting issues"

**Symptoms:**
```
Table output is misaligned
Table output is truncated
```

**Diagnosis:**
1. Check terminal width
2. Verify output content
3. Test with different output formats

**Solutions:**

#### Solution 1: Adjust Terminal Width
```bash
# Set terminal width
export COLUMNS=120

# Or resize terminal window
```

#### Solution 2: Use JSON Output
```bash
# Use JSON for better data handling
virak-cli --output json instance list | jq '.'
```

#### Solution 3: Redirect to File
```bash
# Redirect output to file
virak-cli instance list > output.txt
```

## Performance Issues

### Issue: "Slow command execution"

**Symptoms:**
- Commands take a long time to execute
- API responses are slow

**Diagnosis:**
1. Check network latency
2. Verify API performance
3. Test with debug mode

**Solutions:**

#### Solution 1: Check Network Latency
```bash
# Check ping time to API
ping api.virakcloud.com

# Check HTTP response time
curl -w "@curl-format.txt" -o /dev/null -s https://api.virakcloud.com
```

#### Solution 2: Use Debug Mode
```bash
# Use debug mode to identify bottlenecks
virak-cli --debug instance list
```

#### Solution 3: Use Different Zone
```bash
# Try a different zone with better performance
virak-cli --zoneId zone-456 instance list
```

## Zone-Specific Issues

### Issue: "Zone not found"

**Symptoms:**
```
Error: Zone not found (404)
Details: Zone 'zone-xyz' does not exist
```

**Diagnosis:**
1. Check zone ID
2. Verify zone availability
3. Test with different zones

**Solutions:**

#### Solution 1: List Available Zones
```bash
# List all available zones
virak-cli zone list
```

#### Solution 2: Use Correct Zone ID
```bash
# Use zone ID from list
virak-cli --zoneId zone-123 instance list
```

#### Solution 3: Set Default Zone
```bash
# Set default zone
virak-cli zone list
# Follow prompts to set default zone
```

### Issue: "Insufficient resources in zone"

**Symptoms:**
```
Error: Insufficient resources in zone
```

**Diagnosis:**
1. Check zone resources
2. Verify resource usage
3. Test with different zones

**Solutions:**

#### Solution 1: Check Zone Resources
```bash
# Check zone resources
virak-cli zone resources --zoneId zone-123
```

#### Solution 2: Try Different Zone
```bash
# Try a different zone
virak-cli --zoneId zone-456 instance list
```

#### Solution 3: Wait for Resources
```bash
# Monitor resources over time
watch virak-cli zone resources --zoneId zone-123
```

---

For more help, see our [Troubleshooting Guide](README.md) or [User Guide](../user-guide/README.md).