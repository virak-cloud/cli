# Command Reference

This section provides comprehensive documentation for all Virak CLI commands.

## Table of Contents

- [Command Structure](#command-structure)
- [Global Flags](#global-flags)
- [Output Formats](#output-formats)
- [Command Categories](#command-categories)

## Command Structure

Virak CLI follows a hierarchical command structure:

```
virak-cli [global-flags] <command> [command-flags] [subcommand] [subcommand-flags] [arguments]
```

### Example

```bash
virak-cli --zone-id zone-123 instance create --name my-instance --vm-image-id img-123
```

## Global Flags

These flags can be used with any command:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | Path to configuration file | `~/.virak-cli.yaml` |
| `--zone-id` | `-z` | Zone ID for the operation | From config |
| `--output` | `-o` | Output format (table, json, yaml) | `table` |
| `--debug` | `-d` | Enable debug logging | `false` |
| `--help` | `-h` | Show help for command | |
| `--version` | `-v` | Show version information | |
| `--api-url` | | API base URL | From config |
| `--timeout` | | Request timeout | From config |
| `--insecure` | | Skip SSL verification | `false` |

## Output Formats

### Table Format (default)

Human-readable table format:

```bash
virak-cli instance list
```

```
+------------+----------------+----------+--------+
|    ID      |      NAME      |  STATUS  |  ZONE  |
+------------+----------------+----------+--------+
| instance-1 | my-instance    | RUNNING  | zone-1 |
| instance-2 | another-instance| STOPPED  | zone-1 |
+------------+----------------+----------+--------+
```

### JSON Format

Machine-readable JSON format:

```bash
virak-cli --output json instance list
```

```json
{
  "data": [
    {
      "id": "instance-1",
      "name": "my-instance",
      "status": "RUNNING",
      "zoneId": "zone-1"
    }
  ],
  "metadata": {
    "totalCount": 1,
    "page": 1
  }
}
```

### YAML Format

Human-readable YAML format:

```bash
virak-cli --output yaml instance list
```

```yaml
data:
- id: instance-1
  name: my-instance
  status: RUNNING
  zoneId: zone-1
metadata:
  totalCount: 1
  page: 1
```

## Command Categories

### Authentication Commands

Commands for managing authentication:

- [`login`](../user-guide/authentication.md) - Authenticate with Virak Cloud
- [`logout`](../user-guide/authentication.md) - Log out from Virak Cloud

### Bucket Commands

Object storage management commands:

- [Bucket Reference](bucket.md) - Complete bucket command documentation

### Cluster Commands

Kubernetes cluster management commands:

- [Cluster Reference](cluster.md) - Complete cluster command documentation

### DNS Commands

DNS domain and record management commands:

- [DNS Reference](dns.md) - Complete DNS command documentation

### Instance Commands

Virtual machine lifecycle management commands:

- [Instance Reference](instance.md) - Complete instance command documentation

### Network Commands

Network configuration and security commands:

- [Network Reference](network.md) - Complete network command documentation

### Zone Commands

Zone and resource management commands:

- [Zone Reference](zone.md) - Complete zone command documentation

## Common Patterns

### Resource Identification

Most commands require a resource ID. Resource IDs follow these patterns:

- Instances: `instance-xxxxxxxx`
- Networks: `network-xxxxxxxx`
- Zones: `zone-xxxxxxxx`
- Buckets: `bucket-name` (user-defined)
- Clusters: `cluster-xxxxxxxx`

### Zone Specification

Most commands require a zone ID. You can specify it in several ways:

1. **Configuration file** (recommended):
   ```yaml
   default:
     zoneId: "zone-12345"
   ```

2. **Environment variable**:
   ```bash
   export VIRAK_ZONE_ID="zone-12345"
   ```

3. **Command-line flag**:
   ```bash
   virak-cli --zone-id zone-12345 instance list
   ```

### Filtering and Pagination

List commands support filtering and pagination:

```bash
# Filter by status
virak-cli instance list --filter status=RUNNING

# Paginate results
virak-cli instance list --page 2 --limit 10

# Sort results
virak-cli instance list --sort name --order desc
```

### Resource Creation

Create commands typically follow this pattern:

```bash
virak-cli <resource> create \
  --name <resource-name> \
  --<parameter> <value> \
  --<parameter> <value>
```

### Resource Management

Common resource management commands:

```bash
# List resources
virak-cli <resource> list

# Show resource details
virak-cli <resource> show --id <resource-id>

# Update resource
virak-cli <resource> update --id <resource-id> --parameter <value>

# Delete resource
virak-cli <resource> delete --id <resource-id>
```

## Error Handling

### Common Error Codes

| Error Code | Description | Solution |
|------------|-------------|----------|
| `401` | Authentication failed | Re-authenticate with `virak-cli login` |
| `403` | Permission denied | Check token permissions |
| `404` | Resource not found | Verify resource ID exists |
| `409` | Resource conflict | Resource already exists |
| `422` | Validation error | Check command parameters |
| `500` | Server error | Try again later |

### Error Output Format

Errors are displayed in a consistent format:

```bash
Error: Resource not found (404)
Details: Instance 'instance-123' does not exist in zone 'zone-456'
```

In JSON output:

```json
{
  "error": {
    "code": 404,
    "message": "Resource not found",
    "details": "Instance 'instance-123' does not exist in zone 'zone-456'"
  }
}
```

## Tips and Best Practices

1. **Use Tab Completion**
   ```bash
   # Enable bash completion (if supported)
   source <(virak-cli completion bash)
   ```

2. **Save Common Commands**
   ```bash
   # Create aliases for common commands
   alias vlist='virak-cli instance list'
   alias vshow='virak-cli instance show'
   ```

3. **Use Scripts for Automation**
   ```bash
   #!/bin/bash
   # Create multiple instances
   for i in {1..3}; do
     virak-cli instance create --name "web-$i" --vm-image-id img-123
   done
   ```

4. **Monitor Resource Usage**
   ```bash
   # Watch instance status
   watch virak-cli instance list
   ```

For more detailed information about specific commands, see the individual command reference pages.