# Cluster Commands Reference

This page provides comprehensive documentation for Kubernetes cluster commands in Virak CLI.

## Table of Contents

- [Overview](#overview)
- [Commands](#commands)
  - [cluster create](#cluster-create)
  - [cluster delete](#cluster-delete)
  - [cluster list](#cluster-list)
  - [cluster scale](#cluster-scale)
  - [cluster show](#cluster-show)
  - [cluster start](#cluster-start)
  - [cluster stop](#cluster-stop)
  - [cluster update](#cluster-update)
  - [cluster service events](#cluster-service-events)
  - [cluster service-offerings list](#cluster-service-offerings-list)
  - [cluster versions list](#cluster-versions-list)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

Cluster commands allow you to manage Kubernetes clusters in Virak Cloud. These commands provide complete lifecycle management for Kubernetes clusters, including creation, configuration, scaling, and deletion.

All cluster commands require a zone ID, which can be provided via the `--zoneId` flag or by setting a default zone in the configuration.

## Commands

### cluster create

Create a new Kubernetes cluster.

#### Syntax

```bash
virak-cli cluster create [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--name` | Cluster name |
| `--versionId` | Kubernetes version ID (ULID format) |
| `--offeringId` | Service offering ID for cluster nodes (ULID format) |
| `--sshKeyId` | SSH key ID to be installed on nodes (ULID format) |
| `--networkId` | Network ID to connect the cluster to (ULID format) |

#### Optional Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--zoneId` | Zone ID (overrides default zone) | |
| `--ha` | Enable high availability | `false` |
| `--size` | Cluster size (number of nodes) | `1` |
| `--description` | Cluster description | |
| `--privateRegistryUsername` | Private container registry username | |
| `--privateRegistryPassword` | Private container registry password | |
| `--privateRegistryUrl` | Private container registry URL | |
| `--haControllerNodes` | HA config controller nodes | |
| `--haExternalLBIP` | HA config external load balancer IP | |

#### Examples

```bash
# Create a basic cluster
virak-cli cluster create \
  --name my-cluster \
  --versionId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --offeringId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --sshKeyId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Create a high-availability cluster
virak-cli cluster create \
  --name prod-cluster \
  --versionId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --offeringId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --sshKeyId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --ha \
  --size 3 \
  --haControllerNodes 3

# Create cluster with private registry
virak-cli cluster create \
  --name private-registry-cluster \
  --versionId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --offeringId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --sshKeyId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --privateRegistryUsername myuser \
  --privateRegistryPassword mypass \
  --privateRegistryUrl registry.example.com
```

#### Output

```
kubernetes cluster created successfully. Check status with 'virak-cli cluster list'
```

### cluster delete

Delete a Kubernetes cluster.

**Warning:** This action is irreversible. All data in the cluster will be permanently deleted.

#### Syntax

```bash
virak-cli cluster delete [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--clusterId` | Cluster ID to delete (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Delete a cluster
virak-cli cluster delete --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Delete cluster in specific zone
virak-cli cluster delete --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R --zoneId zone-456
```

#### Output

```
Cluster '01H8X9V2X8Y7Z6W5V4U3T2S1R' deleted successfully.
```

### cluster list

List all Kubernetes clusters in a zone.

#### Syntax

```bash
virak-cli cluster list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# List all clusters
virak-cli cluster list

# List clusters in specific zone
virak-cli cluster list --zoneId zone-456

# List clusters in JSON format
virak-cli --output json cluster list
```

#### Output

```
+----------------------+-------------+---------+---------+--------------+
|          ID          |     NAME    | STATUS  | VERSION | WORKER SIZE  |
+----------------------+-------------+---------+---------+--------------+
| 01H8X9V2X8Y7Z6W5V4U3 | my-cluster  | RUNNING | 1.28.0  |      3       |
| 01H8X9V2X8Y7Z6W5V4U4 | prod-cluster| RUNNING | 1.27.3  |      5       |
+----------------------+-------------+---------+---------+--------------+
```

### cluster scale

Scale a Kubernetes cluster by changing the number of worker nodes.

#### Syntax

```bash
virak-cli cluster scale [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--clusterId` | Cluster ID to scale (ULID format) |
| `--nodePoolSize` | New node pool size |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Scale cluster to 5 nodes
virak-cli cluster scale --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R --nodePoolSize 5

# Scale down to 2 nodes
virak-cli cluster scale --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R --nodePoolSize 2
```

#### Output

```
Cluster '01H8X9V2X8Y7Z6W5V4U3T2S1R' scaled to 5 nodes successfully.
```

### cluster show

Show detailed information about a Kubernetes cluster.

#### Syntax

```bash
virak-cli cluster show [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--clusterId` | Cluster ID (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Show cluster details
virak-cli cluster show --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Show cluster details in JSON format
virak-cli --output json cluster show --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
+----------------------+-------------+---------+---------+----------+----------------+
|          ID          |     NAME    | STATUS  | VERSION |   SIZE   |    CREATED AT  |
+----------------------+-------------+---------+---------+----------+----------------+
| 01H8X9V2X8Y7Z6W5V4U3 | my-cluster  | RUNNING | 1.28.0  |    3     |  1697712000     |
+----------------------+-------------+---------+---------+----------+----------------+
```

### cluster start

Start a stopped Kubernetes cluster.

#### Syntax

```bash
virak-cli cluster start [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--clusterId` | Cluster ID to start (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Start a stopped cluster
virak-cli cluster start --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
Cluster '01H8X9V2X8Y7Z6W5V4U3T2S1R' started successfully.
```

### cluster stop

Stop a running Kubernetes cluster.

#### Syntax

```bash
virak-cli cluster stop [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--clusterId` | Cluster ID to stop (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Stop a running cluster
virak-cli cluster stop --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
Cluster '01H8X9V2X8Y7Z6W5V4U3T2S1R' stopped successfully.
```

### cluster update

Update a Kubernetes cluster configuration.

#### Syntax

```bash
virak-cli cluster update [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--clusterId` | Cluster ID to update (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--name` | New cluster name |
| `--description` | New cluster description |

#### Examples

```bash
# Update cluster name and description
virak-cli cluster update \
  --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --name updated-cluster \
  --description "Updated cluster description"
```

#### Output

```
Cluster '01H8X9V2X8Y7Z6W5V4U3T2S1R' updated successfully.
```

### cluster service events

View events related to Kubernetes cluster services.

#### Syntax

```bash
virak-cli cluster service events [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--clusterId` | Cluster ID (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# View cluster service events
virak-cli cluster service events --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
+------------+----------------------+----------------+----------------+
|    ID      |        TYPE          |   DESCRIPTION  |    TIMESTAMP    |
+------------+----------------------+----------------+----------------+
| event-123  | SERVICE_CREATED      | Service created| 2023-10-19...  |
| event-124  | SERVICE_UPDATED      | Service updated| 2023-10-19...  |
+------------+----------------------+----------------+----------------+
```

### cluster service-offerings list

List available Kubernetes cluster service offerings.

#### Syntax

```bash
virak-cli cluster service-offerings list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# List available service offerings
virak-cli cluster service-offerings list

# List service offerings in JSON format
virak-cli --output json cluster service-offerings list
```

#### Output

```
+----------------------+------------------+----------+--------+
|          ID          |       NAME       |  CPU     |  RAM   |
+----------------------+------------------+----------+--------+
| 01H8X9V2X8Y7Z6W5V4U3 |    small-k8s     |   2      |  4GB   |
| 01H8X9V2X8Y7Z6W5V4U4 |   medium-k8s     |   4      |  8GB   |
| 01H8X9V2X8Y7Z6W5V4U5 |    large-k8s     |   8      | 16GB   |
+----------------------+------------------+----------+--------+
```

### cluster versions list

List available Kubernetes versions.

#### Syntax

```bash
virak-cli cluster versions list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# List available Kubernetes versions
virak-cli cluster versions list

# List versions in JSON format
virak-cli --output json cluster versions list
```

#### Output

```
+----------------------+----------+---------+
|          ID          | VERSION  | STATUS  |
+----------------------+----------+---------+
| 01H8X9V2X8Y7Z6W5V4U3 |  1.28.0  | STABLE  |
| 01H8X9V2X8Y7Z6W5V4U4 |  1.27.3  | STABLE  |
| 01H8X9V2X8Y7Z6W5V4U5 |  1.29.0  | BETA    |
+----------------------+----------+---------+
```

## Examples

### Complete Cluster Lifecycle

```bash
# 1. List available versions and offerings
virak-cli cluster versions list
virak-cli cluster service-offerings list

# 2. Create a new cluster
virak-cli cluster create \
  --name my-app-cluster \
  --versionId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --offeringId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --sshKeyId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 3. List clusters to verify
virak-cli cluster list

# 4. Show cluster details
virak-cli cluster show --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 5. Scale the cluster
virak-cli cluster scale --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R --nodePoolSize 5

# 6. View service events
virak-cli cluster service events --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 7. Stop the cluster
virak-cli cluster stop --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 8. Start the cluster
virak-cli cluster start --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 9. Delete the cluster
virak-cli cluster delete --clusterId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### Creating a High-Availability Cluster

```bash
# Create a HA cluster with 3 control plane nodes
virak-cli cluster create \
  --name prod-ha-cluster \
  --versionId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --offeringId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --sshKeyId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --ha \
  --size 3 \
  --haControllerNodes 3 \
  --haExternalLBIP 192.168.1.100
```

## Best Practices

### Cluster Planning

1. **Choose Appropriate Sizes**: Select service offerings based on your workload requirements
2. **Plan for Growth**: Start with enough capacity to handle expected growth
3. **Use High Availability**: Enable HA for production workloads
4. **Select Stable Versions**: Use stable Kubernetes versions for production

```bash
# Check available offerings before creating
virak-cli cluster service-offerings list

# Use stable versions
virak-cli cluster versions list | grep STABLE
```

### Security

1. **Use SSH Keys**: Always provide SSH keys for node access
2. **Private Networks**: Place clusters in appropriate network segments
3. **Private Registries**: Use private registries for sensitive images
4. **Regular Updates**: Keep clusters updated with security patches

```bash
# Create cluster with private registry
virak-cli cluster create \
  --name secure-cluster \
  --versionId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --offeringId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --sshKeyId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --networkId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --privateRegistryUsername myuser \
  --privateRegistryPassword mypass \
  --privateRegistryUrl registry.example.com
```

### Monitoring

1. **Check Cluster Status**: Regularly monitor cluster status
2. **Watch Service Events**: Keep track of service events
3. **Monitor Resource Usage**: Track CPU and memory usage
4. **Set Up Alerts**: Configure alerts for critical events

```bash
# Monitor cluster status
watch virak-cli cluster list

# Check service events
virak-cli cluster service events --clusterId <cluster-id>
```

### Automation

```bash
#!/bin/bash
# cluster-management.sh

# Function to create cluster
create_cluster() {
  local name=$1
  local version=$2
  local offering=$3
  local ssh_key=$4
  local network=$5
  
  virak-cli cluster create \
    --name "$name" \
    --versionId "$version" \
    --offeringId "$offering" \
    --sshKeyId "$ssh_key" \
    --networkId "$network"
  
  echo "Cluster $name creation initiated"
}

# Function to scale cluster
scale_cluster() {
  local cluster_id=$1
  local size=$2
  
  virak-cli cluster scale --clusterId "$cluster_id" --nodePoolSize "$size"
  echo "Cluster $cluster_id scaled to $size nodes"
}

# Function to backup cluster config
backup_cluster() {
  local cluster_id=$1
  
  virak-cli --output json cluster show --clusterId "$cluster_id" > "cluster-backup-$cluster_id-$(date +%Y%m%d).json"
  echo "Cluster configuration backed up"
}
```

## Troubleshooting

### Common Issues

1. **Cluster Creation Fails**
   ```
   Error: Cluster creation failed
   ```
   Solution: Check service offering, version, network, and SSH key IDs

2. **Cluster Stuck in Provisioning**
   ```
   Status: PROVISIONING
   ```
   Solution: Check service events and wait for provisioning to complete

3. **Scaling Fails**
   ```
   Error: Failed to scale cluster
   ```
   Solution: Check available resources and service offering limits

4. **Cluster Not Accessible**
   ```
   Error: Cluster not reachable
   ```
   Solution: Check network configuration and security groups

### Debug Commands

```bash
# Enable debug logging
virak-cli --debug cluster list

# Check cluster status
virak-cli cluster show --clusterId <cluster-id>

# View service events
virak-cli cluster service events --clusterId <cluster-id>

# Verify available resources
virak-cli cluster service-offerings list
virak-cli zone resources
```

For more troubleshooting help, see our [Troubleshooting Guide](../troubleshooting/common-issues.md).