# Zone Commands Reference

This page provides comprehensive documentation for zone commands in Virak CLI.

## Table of Contents

- [Overview](#overview)
- [Commands](#commands)
  - [zone list](#zone-list)
  - [zone networks](#zone-networks)
  - [zone resources](#zone-resources)
  - [zone services](#zone-services)
- [Zone Management](#zone-management)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

Zone commands allow you to manage zones in Virak Cloud. Zones are geographical regions where your resources are deployed. These commands provide information about available zones, their resources, networks, and services.

Zone commands help you:
- List available zones
- Set a default zone for operations
- View zone resources and capacity
- Check available networks in a zone
- Verify active services in a zone

## Commands

### zone list

List all available zones and optionally set a default zone.

#### Syntax

```bash
virak-cli zone list
```

#### Examples

```bash
# List all available zones
virak-cli zone list

# The command will prompt to set a default zone
```

#### Output

```
[1] Zone Name: Tehran-1, ID: zone-123
[2] Zone Name: Mashhad-1, ID: zone-456
[3] Zone Name: Shiraz-1, ID: zone-789

Do you want to set a default zone? (y/N): y
Enter the number of the default zone: 1
Default zone set to: Tehran-1
```

### zone networks

List all networks for a specific zone.

#### Syntax

```bash
virak-cli zone networks [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID to use (optional if default.zoneId is set in config) |

#### Examples

```bash
# List networks for default zone
virak-cli zone networks

# List networks for specific zone
virak-cli zone networks --zoneId zone-123

# List networks in JSON format
virak-cli --output json zone networks --zoneId zone-123
```

#### Output

```
Networks for Zone:
[1] Name: web-network, ID: net-123, Status: ACTIVE
    Offering: L2 Basic (l2-basic), Type: L2, Rate: 100 Mbps
[2] Name: backend-network, ID: net-456, Status: ACTIVE
    Offering: L3 Isolated (l3-isolated), Type: Isolated, Rate: 200 Mbps
```

### zone resources

List resource usage and limits for a specific zone.

#### Syntax

```bash
virak-cli zone resources [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID to use (optional if default.zoneId is set in config) |

#### Examples

```bash
# List resources for default zone
virak-cli zone resources

# List resources for specific zone
virak-cli zone resources --zoneId zone-123

# List resources in JSON format
virak-cli --output json zone resources --zoneId zone-123
```

#### Output

```
Resources for Zone:
  Memory: 8192/16384 (Megabyte)
  CPU Number: 4/8 (Core)
  Data Volume: 500/1000 (Gigabyte)
  VM Limit: 5/10
```

### zone services

Show active services for a specific zone.

#### Syntax

```bash
virak-cli zone services [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID to use (optional if default.zoneId is set in config) |

#### Examples

```bash
# Show services for default zone
virak-cli zone services

# Show services for specific zone
virak-cli zone services --zoneId zone-123

# Show services in JSON format
virak-cli --output json zone services --zoneId zone-123
```

#### Output

```
Active Services for Zone:
  Instance: true
  DataVolume: true
  Network: true
  ObjectStorage: true
  K8s: true
```

## Zone Management

### Setting a Default Zone

Setting a default zone simplifies command usage by eliminating the need to specify the zone ID for every command.

#### Methods to Set Default Zone

1. **Interactive Method** (Recommended):
   ```bash
   virak-cli zone list
   # Follow the prompts to set a default zone
   ```

2. **Configuration File**:
   ```yaml
   # ~/.virak-cli.yaml
   default:
     zoneId: "zone-123"
     zoneName: "Tehran-1"
   ```

3. **Environment Variable**:
   ```bash
   export VIRAK_ZONE_ID="zone-123"
   ```

#### Benefits of Default Zone

- Simplifies command syntax
- Reduces repetitive zone ID specification
- Minimizes errors from incorrect zone IDs
- Improves workflow efficiency

### Zone Selection Criteria

When selecting a zone, consider:

1. **Geographic Location**: Choose zones closest to your users
2. **Resource Availability**: Check if required services are available
3. **Resource Capacity**: Verify sufficient resources are available
4. **Compliance Requirements**: Consider data residency regulations
5. **Latency Requirements**: Evaluate network latency to your users

## Examples

### Complete Zone Setup Workflow

```bash
# 1. List all available zones
virak-cli zone list

# 2. Set a default zone (interactive)
virak-cli zone list
# Follow prompts to set default zone

# 3. Verify default zone is set
virak-cli zone resources

# 4. Check available services in the zone
virak-cli zone services

# 5. List available networks
virak-cli zone networks

# 6. Check resource usage
virak-cli zone resources

# 7. Create resources in the default zone
virak-cli instance create \
  --name my-instance \
  --service-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --vm-image-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --network-ids '["net-123"]'
```

### Multi-Zone Management

```bash
#!/bin/bash
# multi-zone-management.sh

# Function to check zone status
check_zone_status() {
  local zone_id=$1
  echo "Checking zone: $zone_id"
  
  # Check services
  echo "Services:"
  virak-cli zone services --zoneId "$zone_id"
  
  # Check resources
  echo "Resources:"
  virak-cli zone resources --zoneId "$zone_id"
  
  # Check networks
  echo "Networks:"
  virak-cli zone networks --zoneId "$zone_id"
  
  echo "----------------------------------------"
}

# Function to find best zone for deployment
find_best_zone() {
  local min_cpu=$1
  local min_memory=$2
  
  zones=$(virak-cli --output json zone list | jq -r '.data[].id')
  
  for zone in $zones; do
    resources=$(virak-cli --output json zone resources --zoneId "$zone")
    
    cpu_available=$(echo "$resources" | jq '.data.instanceResourceCollected.cPUNumber.total')
    memory_available=$(echo "$resources" | jq '.data.instanceResourceCollected.memory.total')
    
    if [ "$cpu_available" -ge "$min_cpu" ] && [ "$memory_available" -ge "$min_memory" ]; then
      echo "Zone $zone meets requirements:"
      echo "  CPU: $cpu_available cores"
      echo "  Memory: $memory_available MB"
    fi
  done
}

# Check all zones
zones=("zone-123" "zone-456" "zone-789")
for zone in "${zones[@]}"; do
  check_zone_status "$zone"
done

# Find best zone for deployment (requires 4 CPU, 8GB RAM)
find_best_zone 4 8192
```

### Zone Resource Monitoring

```bash
#!/bin/bash
# zone-monitoring.sh

# Function to monitor zone resources
monitor_zone_resources() {
  local zone_id=$1
  local threshold=${2:-80} # Default threshold 80%
  
  resources=$(virak-cli --output json zone resources --zoneId "$zone_id")
  
  # Calculate resource usage percentages
  memory_used=$(echo "$resources" | jq '.data.instanceResourceCollected.memory.collected')
  memory_total=$(echo "$resources" | jq '.data.instanceResourceCollected.memory.total')
  memory_usage=$((memory_used * 100 / memory_total))
  
  cpu_used=$(echo "$resources" | jq '.data.instanceResourceCollected.cPUNumber.collected')
  cpu_total=$(echo "$resources" | jq '.data.instanceResourceCollected.cPUNumber.total')
  cpu_usage=$((cpu_used * 100 / cpu_total))
  
  volume_used=$(echo "$resources" | jq '.data.instanceResourceCollected.dataVolume.collected')
  volume_total=$(echo "$resources" | jq '.data.instanceResourceCollected.dataVolume.total')
  volume_usage=$((volume_used * 100 / volume_total))
  
  echo "Zone $zone_id Resource Usage:"
  echo "  Memory: $memory_usage% ($memory_used/$memory_total MB)"
  echo "  CPU: $cpu_usage% ($cpu_used/$cpu_total cores)"
  echo "  Volume: $volume_usage% ($volume_used/$volume_total GB)"
  
  # Check thresholds
  if [ "$memory_usage" -gt "$threshold" ]; then
    echo "  WARNING: Memory usage above $threshold%"
  fi
  
  if [ "$cpu_usage" -gt "$threshold" ]; then
    echo "  WARNING: CPU usage above $threshold%"
  fi
  
  if [ "$volume_usage" -gt "$threshold" ]; then
    echo "  WARNING: Volume usage above $threshold%"
  fi
}

# Monitor all zones
zones=$(virak-cli --output json zone list | jq -r '.data[].id')
for zone in $zones; do
  monitor_zone_resources "$zone" 70
  echo ""
done
```

### Zone Service Verification

```bash
#!/bin/bash
# zone-service-check.sh

# Function to verify required services
verify_zone_services() {
  local zone_id=$1
  shift
  local required_services=("$@")
  
  echo "Verifying services in zone $zone_id..."
  
  services=$(virak-cli --output json zone services --zoneId "$zone_id")
  
  all_available=true
  for service in "${required_services[@]}"; do
    service_available=$(echo "$services" | jq -r ".data.$service")
    
    if [ "$service_available" = "true" ]; then
      echo "  ✓ $service: Available"
    else
      echo "  ✗ $service: Not Available"
      all_available=false
    fi
  done
  
  if [ "$all_available" = true ]; then
    echo "All required services are available in zone $zone_id"
    return 0
  else
    echo "Some required services are not available in zone $zone_id"
    return 1
  fi
}

# Check if zone supports web application deployment
verify_zone_services "zone-123" "Instance" "Network" "ObjectStorage"

# Check if zone supports Kubernetes deployment
verify_zone_services "zone-123" "Instance" "Network" "K8s"
```

## Best Practices

### Zone Selection

1. **Choose Nearest Zone**: Select zones closest to your users for better latency
2. **Check Service Availability**: Verify required services are available before deployment
3. **Monitor Resource Usage**: Regularly check resource utilization
4. **Plan for Growth**: Consider future resource requirements

```bash
# Check zone suitability before deployment
virak-cli zone services --zoneId zone-123
virak-cli zone resources --zoneId zone-123
```

### Default Zone Management

1. **Set Default Zone**: Always set a default zone to simplify operations
2. **Document Default Zone**: Keep track of which zone is set as default
3. **Verify Default Zone**: Periodically verify the default zone is correct
4. **Use Zone-Specific Commands**: Override default when working with other zones

```bash
# Set and verify default zone
virak-cli zone list
virak-cli zone resources  # Should show default zone resources

# Override default when needed
virak-cli zone resources --zoneId zone-456
```

### Resource Monitoring

1. **Regular Monitoring**: Monitor resource usage regularly
2. **Set Alerts**: Configure alerts for resource thresholds
3. **Plan Capacity**: Plan for capacity upgrades before hitting limits
4. **Optimize Usage**: Optimize resource usage to reduce costs

```bash
# Monitor resource usage
virak-cli zone resources --zoneId zone-123

# Create monitoring script
#!/bin/bash
while true; do
  virak-cli zone resources --zoneId zone-123
  sleep 300  # Check every 5 minutes
done
```

### Multi-Zone Deployments

1. **Distribute Resources**: Distribute resources across multiple zones
2. **Implement Failover**: Set up failover between zones
3. **Data Replication**: Replicate data across zones
4. **Test Failover**: Regularly test zone failover procedures

```bash
# Deploy across multiple zones
zones=("zone-123" "zone-456")
for zone in "${zones[@]}"; do
  virak-cli instance create \
    --zoneId "$zone" \
    --name "app-$(echo $zone | cut -d'-' -f2)" \
    --service-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
    --vm-image-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
    --network-ids '["net-123"]'
done
```

### Automation

```bash
#!/bin/bash
# zone-automation.sh

# Function to find zone with most available resources
find_best_zone_for_deployment() {
  local best_zone=""
  local best_score=0
  
  zones=$(virak-cli --output json zone list | jq -r '.data[].id')
  
  for zone in $zones; do
    resources=$(virak-cli --output json zone resources --zoneId "$zone")
    services=$(virak-cli --output json zone services --zoneId "$zone")
    
    # Check if all required services are available
    if [ "$(echo "$services" | jq -r '.data.Instance')" = "true" ] && \
       [ "$(echo "$services" | jq -r '.data.Network')" = "true" ]; then
      
      # Calculate availability score
      memory_available=$(echo "$resources" | jq '.data.instanceResourceCollected.memory.total')
      cpu_available=$(echo "$resources" | jq '.data.instanceResourceCollected.cPUNumber.total')
      
      score=$((memory_available + cpu_available * 1024))
      
      if [ "$score" -gt "$best_score" ]; then
        best_score=$score
        best_zone=$zone
      fi
    fi
  done
  
  echo "$best_zone"
}

# Function to deploy to best available zone
deploy_to_best_zone() {
  local instance_name=$1
  
  best_zone=$(find_best_zone_for_deployment)
  
  if [ -n "$best_zone" ]; then
    echo "Deploying $instance_name to zone $best_zone"
    
    virak-cli instance create \
      --zoneId "$best_zone" \
      --name "$instance_name" \
      --service-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
      --vm-image-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
      --network-ids '["net-123"]'
  else
    echo "No suitable zone found for deployment"
    return 1
  fi
}

# Deploy instance to best available zone
deploy_to_best_zone "auto-deployed-instance"
```

## Troubleshooting

### Common Issues

1. **Default Zone Not Set**
   ```
   Error: No default zone configured
   ```
   Solution: Set a default zone using `virak-cli zone list`

2. **Zone Not Found**
   ```
   Error: Zone 'zone-xyz' not found
   ```
   Solution: Check available zones with `virak-cli zone list`

3. **Service Not Available**
   ```
   Error: Service not available in zone
   ```
   Solution: Check available services with `virak-cli zone services`

4. **Insufficient Resources**
   ```
   Error: Insufficient resources in zone
   ```
   Solution: Check resource usage with `virak-cli zone resources`

### Debug Commands

```bash
# Check available zones
virak-cli zone list

# Verify default zone
virak-cli zone resources

# Check zone services
virak-cli zone services --zoneId <zone-id>

# Check zone resources
virak-cli zone resources --zoneId <zone-id>

# Check zone networks
virak-cli zone networks --zoneId <zone-id>

# Enable debug logging
virak-cli --debug zone list
```

For more troubleshooting help, see our [Troubleshooting Guide](../troubleshooting/common-issues.md).