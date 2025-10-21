# Network Commands Reference

This page provides comprehensive documentation for network commands in Virak CLI.

## Table of Contents

- [Overview](#overview)
- [Commands](#commands)
  - [network create](#network-create)
    - [network create l2](#network-create-l2)
    - [network create l3](#network-create-l3)
  - [network delete](#network-delete)
  - [network list](#network-list)
  - [network service-offering](#network-service-offering)
  - [network show](#network-show)
  - [network firewall](#network-firewall)
    - [network firewall ipv4](#network-firewall-ipv4)
    - [network firewall ipv6](#network-firewall-ipv6)
  - [network instance](#network-instance)
  - [network loadbalance](#network-loadbalance)
  - [network public-ip](#network-public-ip)
  - [network vpn](#network-vpn)
- [Network Types](#network-types)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

Network commands allow you to manage network resources in Virak Cloud. These commands provide complete control over your network infrastructure, including creating networks, managing firewall rules, configuring load balancers, and setting up VPN connections.

All network commands require a zone ID, which can be provided via the `--zoneId` flag or by setting a default zone in the configuration.

## Commands

### network create

Create a new network in a zone. This command has subcommands for creating different types of networks.

#### Syntax

```bash
virak-cli network create <type> [flags]
```

### network create l2

Create a new Layer 2 (L2) network.

#### Syntax

```bash
virak-cli network create l2 [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--network-offering-id` | Network offering ID (ULID format) |
| `--name` | Network name |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Create an L2 network
virak-cli network create l2 \
  --name my-l2-network \
  --network-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Create L2 network in specific zone
virak-cli network create l2 \
  --name my-l2-network \
  --network-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --zoneId zone-456
```

#### Output

```
L2 network created successfully
```

### network create l3

Create a new Layer 3 (L3) network.

#### Syntax

```bash
virak-cli network create l3 [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--network-offering-id` | Network offering ID (ULID format) |
| `--name` | Network name |
| `--gateway` | Gateway IP address |
| `--netmask` | Netmask |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Create an L3 network
virak-cli network create l3 \
  --name my-l3-network \
  --network-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --gateway 192.168.1.1 \
  --netmask 255.255.255.0

# Create L3 network with different subnet
virak-cli network create l3 \
  --name backend-network \
  --network-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --gateway 10.0.0.1 \
  --netmask 255.255.255.0
```

#### Output

```
L3 network created successfully
```

### network delete

Delete a network.

**Warning:** This action is irreversible. All resources in the network will be permanently deleted.

#### Syntax

```bash
virak-cli network delete [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--id` | Network ID to delete (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Delete a network
virak-cli network delete --id 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Delete network in specific zone
virak-cli network delete --id 01H8X9V2X8Y7Z6W5V4U3T2S1R --zoneId zone-456
```

#### Output

```
Network '01H8X9V2X8Y7Z6W5V4U3T2S1R' deleted successfully.
```

### network list

List all networks in a zone.

#### Syntax

```bash
virak-cli network list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# List all networks
virak-cli network list

# List networks in specific zone
virak-cli network list --zoneId zone-456

# List networks in JSON format
virak-cli --output json network list
```

#### Output

```
+----------------------+----------------+----------+--------+
|          ID          |      NAME      |  TYPE    | STATUS |
+----------------------+----------------+----------+--------+
| 01H8X9V2X8Y7Z6W5V4U3 | my-l2-network  |    L2    | ACTIVE |
| 01H8X9V2X8Y7Z6W5V4U4 | my-l3-network  |    L3    | ACTIVE |
+----------------------+----------------+----------+--------+
```

### network service-offering

List available network service offerings.

#### Syntax

```bash
virak-cli network service-offering [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--type` | Filter by service offering type: l2, l3, or all |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# List all network service offerings
virak-cli network service-offering --type all

# List only L2 service offerings
virak-cli network service-offering --type l2

# List only L3 service offerings
virak-cli network service-offering --type l3

# List service offerings in JSON format
virak-cli --output json network service-offering --type all
```

#### Output

```
+----------------------+----------------+----------------+--------+-----------+------------+-----------+--------+----------+----------------+
|          ID          |      NAME      |  DISPLAY NAME  | PRICE  | OVERPRICE | PLAN(GiB)  | RATE(Mbps)|  TYPE  | PROTOCOL |      DESC      |
+----------------------+----------------+----------------+--------+-----------+------------+-----------+--------+----------+----------------+
| 01H8X9V2X8Y7Z6W5V4U3 |   l2-basic     |   L2 Basic     |  10.00 |    5.00   |    100     |    100    |   L2   |   IPv4   | Basic L2 net   |
| 01H8X9V2X8Y7Z6W5V4U4 |   l3-isolated  |  L3 Isolated   |  15.00 |    7.50   |    200     |    200    | Isolated|   IPv4   | Isolated L3 net|
+----------------------+----------------+----------------+--------+-----------+------------+-----------+--------+----------+----------------+
```

### network show

Show detailed information about a network.

#### Syntax

```bash
virak-cli network show [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--networkId` | Network ID to show (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Show network details
virak-cli network show --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Show network details in JSON format
virak-cli --output json network show --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
Network Details:
  ID: 01H8X9V2X8Y7Z6W5V4U3T2S1R
  Name: my-l2-network
  Type: L2
  Status: ACTIVE
  Zone: zone-123
  Created: 2023-10-19 10:30:00 UTC
  Updated: 2023-10-19 10:30:00 UTC
  Gateway: N/A
  Netmask: N/A
  CIDR: N/A
```

### network firewall

Manage network firewall rules. This command has subcommands for IPv4 and IPv6 firewall rules.

#### Syntax

```bash
virak-cli network firewall <version> <action> [flags]
```

### network firewall ipv4

Manage IPv4 firewall rules.

#### Subcommands

- `create`: Create a new IPv4 firewall rule
- `delete`: Delete an IPv4 firewall rule
- `list`: List all IPv4 firewall rules

#### Examples

```bash
# Create an IPv4 firewall rule
virak-cli network firewall ipv4 create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --protocol tcp \
  --startPort 80 \
  --endPort 80 \
  --cidr 0.0.0.0/0

# List IPv4 firewall rules
virak-cli network firewall ipv4 list --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Delete an IPv4 firewall rule
virak-cli network firewall ipv4 delete \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --ruleId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### network firewall ipv6

Manage IPv6 firewall rules.

#### Subcommands

- `create`: Create a new IPv6 firewall rule
- `delete`: Delete an IPv6 firewall rule
- `list`: List all IPv6 firewall rules

#### Examples

```bash
# Create an IPv6 firewall rule
virak-cli network firewall ipv6 create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --protocol tcp \
  --startPort 443 \
  --endPort 443 \
  --cidr ::/0

# List IPv6 firewall rules
virak-cli network firewall ipv6 list --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Delete an IPv6 firewall rule
virak-cli network firewall ipv6 delete \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --ruleId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### network instance

Manage network instance connections.

#### Subcommands

- `connect`: Connect an instance to a network
- `disconnect`: Disconnect an instance from a network
- `list`: List all instances in a network

#### Examples

```bash
# Connect an instance to a network
virak-cli network instance connect \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Disconnect an instance from a network
virak-cli network instance disconnect \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# List instances in a network
virak-cli network instance list --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### network loadbalance

Manage network load balancer rules.

#### Subcommands

- `create`: Create a new load balancer
- `delete`: Delete a load balancer
- `list`: List all load balancers
- `assign`: Assign a load balancer
- `deassign`: Deassign a load balancer
- `haproxy`: Manage HAProxy load balancers

#### Examples

```bash
# Create a load balancer
virak-cli network loadbalance create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --name my-lb \
  --algorithm round-robin

# List load balancers
virak-cli network loadbalance list --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# View HAProxy stats
virak-cli network loadbalance haproxy live \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --lbId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### network public-ip

Manage public IP addresses.

#### Subcommands

- `list`: List all public IPs
- `associate`: Associate a public IP
- `disassociate`: Disassociate a public IP
- `staticnat`: Manage static NAT

#### Examples

```bash
# List public IPs
virak-cli network public-ip list --zoneId zone-123

# Associate a public IP
virak-cli network public-ip associate \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Enable static NAT
virak-cli network public-ip staticnat enable \
  --ipId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### network vpn

Manage VPN connections.

#### Subcommands

- `enable`: Enable VPN
- `disable`: Disable VPN
- `show`: Show VPN details
- `update`: Update VPN configuration

#### Examples

```bash
# Enable VPN for a network
virak-cli network vpn enable --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Show VPN details
virak-cli network vpn show --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Update VPN configuration
virak-cli network vpn update \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --ipsecPsk "new-secret-key"

# Disable VPN
virak-cli network vpn disable --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

## Network Types

### Layer 2 (L2) Networks

Layer 2 networks operate at the data link layer and provide:
- MAC address-based communication
- Broadcast domain support
- Simple network topology
- Lower overhead compared to L3 networks

**Use cases:**
- Simple application deployments
- Development environments
- When you need broadcast capabilities

### Layer 3 (L3) Networks

Layer 3 networks operate at the network layer and provide:
- IP routing capabilities
- Network isolation
- Custom gateway configuration
- Subnet management

**Use cases:**
- Multi-tier applications
- Production environments
- When you need network isolation
- Complex network topologies

## Examples

### Complete Network Setup

```bash
# 1. List available network service offerings
virak-cli network service-offering --type all

# 2. Create an L2 network for web servers
virak-cli network create l2 \
  --name web-network \
  --network-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 3. Create an L3 network for backend services
virak-cli network create l3 \
  --name backend-network \
  --network-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --gateway 10.0.0.1 \
  --netmask 255.255.255.0

# 4. List networks to verify
virak-cli network list

# 5. Show network details
virak-cli network show --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 6. Configure firewall rules
virak-cli network firewall ipv4 create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --protocol tcp \
  --startPort 80 \
  --endPort 80 \
  --cidr 0.0.0.0/0

# 7. Connect instances to networks
virak-cli network instance connect \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### Firewall Configuration

```bash
# Create firewall rules for web server
virak-cli network firewall ipv4 create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --protocol tcp \
  --startPort 80 \
  --endPort 80 \
  --cidr 0.0.0.0/0 \
  --description "HTTP access"

virak-cli network firewall ipv4 create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --protocol tcp \
  --startPort 443 \
  --endPort 443 \
  --cidr 0.0.0.0/0 \
  --description "HTTPS access"

# Create firewall rule for SSH access (restricted)
virak-cli network firewall ipv4 create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --protocol tcp \
  --startPort 22 \
  --endPort 22 \
  --cidr 192.168.1.0/24 \
  --description "SSH access from office"

# List firewall rules
virak-cli network firewall ipv4 list --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### Load Balancer Setup

```bash
# Create a load balancer
virak-cli network loadbalance create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --name web-lb \
  --algorithm round-robin \
  --protocol tcp \
  --port 80

# Assign instances to load balancer
virak-cli network loadbalance assign \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --lbId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# View load balancer status
virak-cli network loadbalance list --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# View HAProxy statistics
virak-cli network loadbalance haproxy live \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --lbId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### Public IP Management

```bash
# List available public IPs
virak-cli network public-ip list --zoneId zone-123

# Associate a public IP with an instance
virak-cli network public-ip associate \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Enable static NAT
virak-cli network public-ip staticnat enable \
  --ipId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Disable static NAT
virak-cli network public-ip staticnat disable \
  --ipId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Disassociate public IP
virak-cli network public-ip disassociate \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

## Best Practices

### Network Planning

1. **Choose Appropriate Network Types**: Use L2 for simple deployments, L3 for complex architectures
2. **Plan IP Addressing**: Design your IP address scheme carefully
3. **Use Descriptive Names**: Use clear, descriptive names for networks
4. **Document Network Topology**: Keep documentation of your network design

```bash
# Check available service offerings before creating networks
virak-cli network service-offering --type all

# Plan your network addressing
# L2 networks: No IP planning needed
# L3 networks: Plan gateway and subnet carefully
```

### Security

1. **Implement Firewall Rules**: Configure appropriate firewall rules for security
2. **Use Network Isolation**: Separate different application tiers
3. **Limit Access**: Restrict access to sensitive services
4. **Regular Audits**: Regularly review firewall rules and access

```bash
# Create restrictive firewall rules
virak-cli network firewall ipv4 create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --protocol tcp \
  --startPort 22 \
  --endPort 22 \
  --cidr 192.168.1.0/24

# Review firewall rules
virak-cli network firewall ipv4 list --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### Performance

1. **Choose Appropriate Service Offerings**: Select service offerings based on performance requirements
2. **Monitor Network Usage**: Track network bandwidth and performance
3. **Optimize Load Balancer Configuration**: Configure load balancers for optimal performance
4. **Use Content Delivery Networks**: Offload static content to CDNs

```bash
# Monitor network performance
virak-cli network show --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Check load balancer performance
virak-cli network loadbalance haproxy live \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --lbId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### High Availability

1. **Use Multiple Networks**: Distribute services across multiple networks
2. **Configure Load Balancers**: Use load balancers for high availability
3. **Implement Redundancy**: Create redundant network paths
4. **Monitor Network Health**: Set up monitoring for network components

```bash
# Create load balancer for high availability
virak-cli network loadbalance create \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --name ha-lb \
  --algorithm round-robin

# Assign multiple instances
virak-cli network loadbalance assign \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --lbId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### Automation

```bash
#!/bin/bash
# network-automation.sh

# Function to create a complete network setup
setup_network_stack() {
  local env_name=$1
  local zone_id=$2
  
  # Create L2 network for web servers
  web_network_id=$(virak-cli --output json network create l2 \
    --name "$env_name-web-network" \
    --network-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
    --zoneId "$zone_id" | jq -r '.data.id')
  
  # Create L3 network for backend services
  backend_network_id=$(virak-cli --output json network create l3 \
    --name "$env_name-backend-network" \
    --network-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
    --gateway "10.0.0.1" \
    --netmask "255.255.255.0" \
    --zoneId "$zone_id" | jq -r '.data.id')
  
  # Configure firewall rules
  virak-cli network firewall ipv4 create \
    --networkId "$web_network_id" \
    --protocol tcp \
    --startPort 80 \
    --endPort 80 \
    --cidr "0.0.0.0/0"
  
  virak-cli network firewall ipv4 create \
    --networkId "$web_network_id" \
    --protocol tcp \
    --startPort 443 \
    --endPort 443 \
    --cidr "0.0.0.0/0"
  
  echo "Network stack created for $env_name"
  echo "Web Network ID: $web_network_id"
  echo "Backend Network ID: $backend_network_id"
}

# Function to backup network configuration
backup_network_config() {
  local backup_dir="network-backup-$(date +%Y%m%d)"
  mkdir -p "$backup_dir"
  
  # Export network configurations
  networks=$(virak-cli --output json network list | jq -r '.data[].id')
  
  for network in $networks; do
    virak-cli --output json network show --networkId "$network" > "$backup_dir/network-$network.json"
    
    # Export firewall rules
    virak-cli --output json network firewall ipv4 list --networkId "$network" > "$backup_dir/firewall-ipv4-$network.json" 2>/dev/null || true
    virak-cli --output json network firewall ipv6 list --networkId "$network" > "$backup_dir/firewall-ipv6-$network.json" 2>/dev/null || true
  done
  
  echo "Network configuration backed up to $backup_dir"
}

# Function to cleanup unused networks
cleanup_unused_networks() {
  networks=$(virak-cli --output json network list | jq -r '.data[] | select(.instance_count == 0) | .id')
  
  for network in $networks; do
    echo "Deleting unused network: $network"
    virak-cli network delete --id "$network"
  done
  
  echo "Cleanup completed"
}
```

## Troubleshooting

### Common Issues

1. **Network Creation Fails**
   ```
   Error: Network creation failed
   ```
   Solution: Check service offering availability and zone resources

2. **Instance Cannot Connect to Network**
   ```
   Error: Instance connection failed
   ```
   Solution: Check network status and instance configuration

3. **Firewall Rules Not Working**
   ```
   Error: Firewall rule not applied
   ```
   Solution: Verify rule syntax and network configuration

4. **Load Balancer Not Responding**
   ```
   Error: Load balancer not accessible
   ```
   Solution: Check load balancer configuration and health of backend instances

### Debug Commands

```bash
# Enable debug logging
virak-cli --debug network list

# Check network status
virak-cli network show --networkId <network-id>

# Verify firewall rules
virak-cli network firewall ipv4 list --networkId <network-id>

# Check load balancer status
virak-cli network loadbalance list --networkId <network-id>

# Test network connectivity
ping <instance-ip>
traceroute <instance-ip>
```

For more troubleshooting help, see our [Troubleshooting Guide](../troubleshooting/common-issues.md).