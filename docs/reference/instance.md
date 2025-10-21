# Instance Commands Reference

This page provides comprehensive documentation for instance (virtual machine) commands in Virak CLI.

## Table of Contents

- [Overview](#overview)
- [Commands](#commands)
  - [instance create](#instance-create)
  - [instance delete](#instance-delete)
  - [instance list](#instance-list)
  - [instance metrics](#instance-metrics)
  - [instance reboot](#instance-reboot)
  - [instance rebuild](#instance-rebuild)
  - [instance service-offering list](#instance-service-offering-list)
  - [instance show](#instance-show)
  - [instance snapshot create](#instance-snapshot-create)
  - [instance snapshot delete](#instance-snapshot-delete)
  - [instance snapshot list](#instance-snapshot-list)
  - [instance snapshot revert](#instance-snapshot-revert)
  - [instance start](#instance-start)
  - [instance stop](#instance-stop)
  - [instance vm-image list](#instance-vm-image-list)
  - [instance volume attach](#instance-volume-attach)
  - [instance volume create](#instance-volume-create)
  - [instance volume delete](#instance-volume-delete)
  - [instance volume detach](#instance-volume-detach)
  - [instance volume list](#instance-volume-list)
  - [instance volume service-offering list](#instance-volume-service-offering-list)
- [Instance States](#instance-states)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

Instance commands allow you to manage virtual machine instances in Virak Cloud. These commands provide complete lifecycle management for VM instances, including creation, configuration, monitoring, and deletion.

All instance commands require a zone ID, which can be provided via the `--zoneId` flag or by setting a default zone in the configuration.

## Commands

### instance create

Create a new virtual machine instance.

#### Syntax

```bash
virak-cli instance create [flags]
```

#### Required Flags (Non-interactive Mode)

| Flag | Description |
|------|-------------|
| `--service-offering-id` | ID of the service offering (ULID format) |
| `--vm-image-id` | ID of the VM image (ULID format) |
| `--network-ids` | JSON array of network IDs, e.g. '["id1","id2"]' |
| `--name` | Name of the instance |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance creation workflow |

#### Examples

```bash
# Create instance in non-interactive mode
virak-cli instance create \
  --name my-instance \
  --service-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --vm-image-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --network-ids '["01H8X9V2X8Y7Z6W5V4U3T2S1R"]'

# Create instance in interactive mode
virak-cli instance create --interactive
```

#### Output

```
Instance creation request accepted. Your instance will be created soon.
Please check the instance list to see when it becomes active.
```

### instance delete

Delete a virtual machine instance.

**Warning:** This action is irreversible. All data on the instance will be permanently deleted.

#### Syntax

```bash
virak-cli instance delete [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID to delete (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Delete an instance
virak-cli instance delete --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Delete instance in interactive mode
virak-cli instance delete --interactive
```

#### Output

```
Instance '01H8X9V2X8Y7Z6W5V4U3T2S1R' deleted successfully.
```

### instance list

List all virtual machine instances in a zone.

#### Syntax

```bash
virak-cli instance list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--columns` | Comma-separated list of columns to display |
| `--list-columns` | Show all valid columns for instance list output |

#### Examples

```bash
# List all instances with default columns
virak-cli instance list

# List instances with custom columns
virak-cli instance list --columns "id,name,status,service_offering.name,vm_image.os_name"

# Show all available columns
virak-cli instance list --list-columns

# List instances in JSON format
virak-cli --output json instance list
```

#### Output

```
+----------------------+-------------+---------+----------------+
|          ID          |     NAME    | STATUS  |    CREATED     |
+----------------------+-------------+---------+----------------+
| 01H8X9V2X8Y7Z6W5V4U3 | my-instance | RUNNING |  1697712000    |
| 01H8X9V2X8Y7Z6W5V4U4 | web-server  | STOPPED |  1697708000    |
+----------------------+-------------+---------+----------------+
```

### instance metrics

View metrics for a virtual machine instance.

#### Syntax

```bash
virak-cli instance metrics [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# View instance metrics
virak-cli instance metrics --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# View metrics in JSON format
virak-cli --output json instance metrics --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
+----------------+--------+-----------+-----------+
|     METRIC     |  VALUE |   UNIT    | TIMESTAMP |
+----------------+--------+-----------+-----------+
| CPU Usage      |   45   |    %      | 1697712000|
| Memory Usage   |  2048  |    MB     | 1697712000|
| Disk Usage     |  10240 |    MB     | 1697712000|
| Network In     |  1250  |   KB/s    | 1697712000|
| Network Out    |   890  |   KB/s    | 1697712000|
+----------------+--------+-----------+-----------+
```

### instance reboot

Reboot a virtual machine instance.

#### Syntax

```bash
virak-cli instance reboot [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID to reboot (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Reboot an instance
virak-cli instance reboot --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Reboot instance in interactive mode
virak-cli instance reboot --interactive
```

#### Output

```
Instance '01H8X9V2X8Y7Z6W5V4U3T2S1R' reboot initiated successfully.
```

### instance rebuild

Rebuild a virtual machine instance from a VM image.

#### Syntax

```bash
virak-cli instance rebuild [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID to rebuild (ULID format) |
| `--vmImageId` | VM image ID to rebuild from (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Rebuild instance with new image
virak-cli instance rebuild \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --vmImageId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
Instance '01H8X9V2X8Y7Z6W5V4U3T2S1R' rebuild initiated successfully.
```

### instance service-offering list

List available instance service offerings.

#### Syntax

```bash
virak-cli instance service-offering list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# List available service offerings
virak-cli instance service-offering list

# List service offerings in JSON format
virak-cli --output json instance service-offering list
```

#### Output

```
+----------------------+------------------+----------+--------+---------+
|          ID          |       NAME       |   CPU    |  RAM   |  PRICE  |
+----------------------+------------------+----------+--------+---------+
| 01H8X9V2X8Y7Z6W5V4U3 |     small-vm     |    1     |  2GB   |  1000   |
| 01H8X9V2X8Y7Z6W5V4U4 |    medium-vm     |    2     |  4GB   |  2000   |
| 01H8X9V2X8Y7Z6W5V4U5 |     large-vm     |    4     |  8GB   |  4000   |
+----------------------+------------------+----------+--------+---------+
```

### instance show

Show detailed information about a virtual machine instance.

#### Syntax

```bash
virak-cli instance show [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Show instance details
virak-cli instance show --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Show instance details in interactive mode
virak-cli instance show --interactive

# Show instance details in JSON format
virak-cli --output json instance show --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
+----------------------+----------------+
|        FIELD         |      VALUE      |
+----------------------+----------------+
| ID                   | 01H8X9V2X8Y7Z6W5V4U3 |
| Name                 | my-instance     |
| Status               | RUNNING         |
| Instance Status      | RUNNING         |
| Zone ID              | zone-123        |
| Created At           | 1697712000      |
| Updated At           | 1697712500      |
| Username             | admin           |
| Password             | password123     |
| VM Image Name        | ubuntu-22.04    |
| VM Image OS          | Ubuntu 22.04    |
| Service Offering     | medium-vm       |
+----------------------+----------------+
```

### instance snapshot create

Create a snapshot of a virtual machine instance.

#### Syntax

```bash
virak-cli instance snapshot create [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID to snapshot (ULID format) |
| `--name` | Snapshot name |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--description` | Snapshot description |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Create a snapshot
virak-cli instance snapshot create \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --name my-snapshot \
  --description "Snapshot before update"

# Create snapshot in interactive mode
virak-cli instance snapshot create --interactive
```

#### Output

```
Snapshot 'my-snapshot' created successfully for instance '01H8X9V2X8Y7Z6W5V4U3T2S1R'.
```

### instance snapshot delete

Delete a snapshot of a virtual machine instance.

#### Syntax

```bash
virak-cli instance snapshot delete [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |
| `--snapshotId` | Snapshot ID to delete (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Delete a snapshot
virak-cli instance snapshot delete \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --snapshotId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Delete snapshot in interactive mode
virak-cli instance snapshot delete --interactive
```

#### Output

```
Snapshot '01H8X9V2X8Y7Z6W5V4U3T2S1R' deleted successfully.
```

### instance snapshot list

List all snapshots for a virtual machine instance.

#### Syntax

```bash
virak-cli instance snapshot list [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# List snapshots for an instance
virak-cli instance snapshot list --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# List snapshots in interactive mode
virak-cli instance snapshot list --interactive

# List snapshots in JSON format
virak-cli --output json instance snapshot list --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
+----------------------+--------------+----------------+----------------+
|          ID          |     NAME     |     STATUS     |    CREATED     |
+----------------------+--------------+----------------+----------------+
| 01H8X9V2X8Y7Z6W5V4U3 | my-snapshot  |    COMPLETED   |  1697712000    |
| 01H8X9V2X8Y7Z6W5V4U4 | backup-2023  |    COMPLETED   |  1697708000    |
+----------------------+--------------+----------------+----------------+
```

### instance snapshot revert

Revert a virtual machine instance to a snapshot.

#### Syntax

```bash
virak-cli instance snapshot revert [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |
| `--snapshotId` | Snapshot ID to revert to (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Revert to a snapshot
virak-cli instance snapshot revert \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --snapshotId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Revert in interactive mode
virak-cli instance snapshot revert --interactive
```

#### Output

```
Instance '01H8X9V2X8Y7Z6W5V4U3T2S1R' reverted to snapshot '01H8X9V2X8Y7Z6W5V4U3T2S1R' successfully.
```

### instance start

Start a stopped virtual machine instance.

#### Syntax

```bash
virak-cli instance start [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID to start (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Start an instance
virak-cli instance start --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Start instance in interactive mode
virak-cli instance start --interactive
```

#### Output

```
Instance '01H8X9V2X8Y7Z6W5V4U3T2S1R' start initiated successfully.
```

### instance stop

Stop a running virtual machine instance.

#### Syntax

```bash
virak-cli instance stop [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID to stop (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Stop an instance
virak-cli instance stop --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Stop instance in interactive mode
virak-cli instance stop --interactive
```

#### Output

```
Instance '01H8X9V2X8Y7Z6W5V4U3T2S1R' stop initiated successfully.
```

### instance vm-image list

List available VM images.

#### Syntax

```bash
virak-cli instance vm-image list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# List available VM images
virak-cli instance vm-image list

# List VM images in JSON format
virak-cli --output json instance vm-image list
```

#### Output

```
+----------------------+------------------+----------+---------+
|          ID          |       NAME       |   TYPE   |  OS TYPE |
+----------------------+------------------+----------+---------+
| 01H8X9V2X8Y7Z6W5V4U3 |   ubuntu-22.04   |   OS     |  Linux  |
| 01H8X9V2X8Y7Z6W5V4U4 |  centos-8-stream  |   OS     |  Linux  |
| 01H8X9V2X8Y7Z6W5V4U5 |   windows-2022   |   OS     | Windows |
+----------------------+------------------+----------+---------+
```

### instance volume attach

Attach a volume to a virtual machine instance.

#### Syntax

```bash
virak-cli instance volume attach [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |
| `--volumeId` | Volume ID to attach (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Attach a volume to an instance
virak-cli instance volume attach \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --volumeId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Attach volume in interactive mode
virak-cli instance volume attach --interactive
```

#### Output

```
Volume '01H8X9V2X8Y7Z6W5V4U3T2S1R' attached to instance '01H8X9V2X8Y7Z6W5V4U3T2S1R' successfully.
```

### instance volume create

Create a new volume for a virtual machine instance.

#### Syntax

```bash
virak-cli instance volume create [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |
| `--name` | Volume name |
| `--serviceOfferingId` | Volume service offering ID (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Create a volume for an instance
virak-cli instance volume create \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --name data-volume \
  --serviceOfferingId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Create volume in interactive mode
virak-cli instance volume create --interactive
```

#### Output

```
Volume 'data-volume' created successfully for instance '01H8X9V2X8Y7Z6W5V4U3T2S1R'.
```

### instance volume delete

Delete a volume from a virtual machine instance.

#### Syntax

```bash
virak-cli instance volume delete [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |
| `--volumeId` | Volume ID to delete (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Delete a volume
virak-cli instance volume delete \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --volumeId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Delete volume in interactive mode
virak-cli instance volume delete --interactive
```

#### Output

```
Volume '01H8X9V2X8Y7Z6W5V4U3T2S1R' deleted successfully.
```

### instance volume detach

Detach a volume from a virtual machine instance.

#### Syntax

```bash
virak-cli instance volume detach [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |
| `--volumeId` | Volume ID to detach (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# Detach a volume from an instance
virak-cli instance volume detach \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --volumeId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# Detach volume in interactive mode
virak-cli instance volume detach --interactive
```

#### Output

```
Volume '01H8X9V2X8Y7Z6W5V4U3T2S1R' detached from instance '01H8X9V2X8Y7Z6W5V4U3T2S1R' successfully.
```

### instance volume list

List all volumes for a virtual machine instance.

#### Syntax

```bash
virak-cli instance volume list [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--instanceId` | Instance ID (ULID format) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--interactive` | Run interactive instance selection workflow |

#### Examples

```bash
# List volumes for an instance
virak-cli instance volume list --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# List volumes in interactive mode
virak-cli instance volume list --interactive

# List volumes in JSON format
virak-cli --output json instance volume list --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

#### Output

```
+----------------------+-------------+--------+----------+
|          ID          |     NAME    | SIZE   |  STATUS  |
+----------------------+-------------+--------+----------+
| 01H8X9V2X8Y7Z6W5V4U3 | data-volume |  100GB |  ATTACHED|
| 01H8X9V2X8Y7Z6W5V4U4 | backup-vol  |  500GB |  DETACHED|
+----------------------+-------------+--------+----------+
```

### instance volume service-offering list

List available volume service offerings.

#### Syntax

```bash
virak-cli instance volume service-offering list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# List available volume service offerings
virak-cli instance volume service-offering list

# List volume service offerings in JSON format
virak-cli --output json instance volume service-offering list
```

#### Output

```
+----------------------+------------------+--------+---------+
|          ID          |       NAME       |  SIZE  |  PRICE  |
+----------------------+------------------+--------+---------+
| 01H8X9V2X8Y7Z6W5V4U3 |    small-vol     |  100GB |   500   |
| 01H8X9V2X8Y7Z6W5V4U4 |   medium-vol     |  500GB |  2000   |
| 01H8X9V2X8Y7Z6W5V4U5 |    large-vol     |  1TB   |  4000   |
+----------------------+------------------+--------+---------+
```

## Instance States

Instances can be in one of the following states:

| State | Description |
|-------|-------------|
| `CREATING` | Instance is being provisioned |
| `RUNNING` | Instance is running and operational |
| `STOPPING` | Instance is in the process of stopping |
| `STOPPED` | Instance is stopped |
| `STARTING` | Instance is in the process of starting |
| `REBOOTING` | Instance is rebooting |
| `DELETING` | Instance is being deleted |
| `ERROR` | Instance is in an error state |
| `REBUILDING` | Instance is being rebuilt |
| `MAINTENANCE` | Instance is under maintenance |

## Examples

### Complete Instance Lifecycle

```bash
# 1. List available service offerings and VM images
virak-cli instance service-offering list
virak-cli instance vm-image list

# 2. Create a new instance
virak-cli instance create \
  --name my-web-server \
  --service-offering-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --vm-image-id 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --network-ids '["01H8X9V2X8Y7Z6W5V4U3T2S1R"]'

# 3. List instances to verify
virak-cli instance list

# 4. Show instance details
virak-cli instance show --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 5. View instance metrics
virak-cli instance metrics --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 6. Create a snapshot
virak-cli instance snapshot create \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --name initial-setup

# 7. Create and attach a volume
virak-cli instance volume create \
  --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R \
  --name data-storage \
  --serviceOfferingId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 8. Stop the instance
virak-cli instance stop --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 9. Start the instance
virak-cli instance start --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R

# 10. Delete the instance
virak-cli instance delete --instanceId 01H8X9V2X8Y7Z6W5V4U3T2S1R
```

### Interactive Instance Management

```bash
# Create instance interactively
virak-cli instance create --interactive

# Show instance details interactively
virak-cli instance show --interactive

# Manage snapshots interactively
virak-cli instance snapshot create --interactive
virak-cli instance snapshot list --interactive
virak-cli instance snapshot revert --interactive

# Manage volumes interactively
virak-cli instance volume create --interactive
virak-cli instance volume attach --interactive
virak-cli instance volume detach --interactive
```

### Batch Operations

```bash
#!/bin/bash
# batch-instance-operations.sh

# Function to create multiple instances
create_instances() {
  local base_name=$1
  local count=$2
  local offering_id=$3
  local image_id=$4
  local network_ids=$5
  
  for i in $(seq 1 $count); do
    virak-cli instance create \
      --name "$base_name-$i" \
      --service-offering-id "$offering_id" \
      --vm-image-id "$image_id" \
      --network-ids "$network_ids"
    
    echo "Instance $base_name-$i creation initiated"
  done
}

# Function to stop all instances
stop_all_instances() {
  instances=$(virak-cli --output json instance list | jq -r '.data[].id')
  
  for instance in $instances; do
    virak-cli instance stop --instanceId "$instance"
    echo "Instance $instance stop initiated"
  done
}

# Function to create snapshots for all instances
snapshot_all_instances() {
  instances=$(virak-cli --output json instance list | jq -r '.data[].id')
  
  for instance in $instances; do
    snapshot_name="backup-$(date +%Y%m%d)-$(echo $instance | cut -c1-8)"
    virak-cli instance snapshot create \
      --instanceId "$instance" \
      --name "$snapshot_name"
    
    echo "Snapshot $snapshot_name created for instance $instance"
  done
}
```

## Best Practices

### Instance Planning

1. **Choose Appropriate Service Offerings**: Select service offerings based on your workload requirements
2. **Plan for Growth**: Start with enough capacity to handle expected growth
3. **Use Appropriate VM Images**: Choose images that match your application requirements
4. **Consider High Availability**: Use multiple instances for critical applications

```bash
# Check available offerings before creating
virak-cli instance service-offering list

# Use appropriate images
virak-cli instance vm-image list
```

### Security

1. **Use Strong Credentials**: Change default passwords after instance creation
2. **Keep Software Updated**: Regularly update operating systems and applications
3. **Use Security Groups**: Configure appropriate firewall rules
4. **Monitor Access**: Keep track of who accesses your instances

```bash
# Check instance credentials
virak-cli instance show --instanceId <instance-id>

# Monitor instance access
virak-cli instance metrics --instanceId <instance-id>
```

### Backup and Recovery

1. **Regular Snapshots**: Create snapshots before making changes
2. **Test Restores**: Regularly test snapshot restoration
3. **Document Recovery Procedures**: Keep clear documentation of recovery steps
4. **Monitor Snapshot Status**: Ensure snapshots complete successfully

```bash
# Create regular snapshots
virak-cli instance snapshot create \
  --instanceId <instance-id> \
  --name "backup-$(date +%Y%m%d)"

# List snapshots
virak-cli instance snapshot list --instanceId <instance-id>
```

### Performance Monitoring

1. **Monitor Resource Usage**: Track CPU, memory, and disk usage
2. **Set Up Alerts**: Configure alerts for critical metrics
3. **Optimize Resources**: Right-size instances based on actual usage
4. **Regular Maintenance**: Schedule regular maintenance windows

```bash
# Monitor instance metrics
virak-cli instance metrics --instanceId <instance-id>

# Check instance status
virak-cli instance list --columns "id,name,status,service_offering.name"
```

### Cost Management

1. **Choose Cost-Effective Offerings**: Select service offerings that balance cost and performance
2. **Stop Unused Instances**: Stop instances when not in use
3. **Use Snapshots Wisely**: Clean up old snapshots regularly
4. **Monitor Costs**: Regularly review your resource usage and costs

```bash
# List instances with cost information
virak-cli instance list --columns "id,name,status,service_offering.hourly_price.up"

# Stop unused instances
virak-cli instance stop --instanceId <unused-instance-id>
```

### Automation

```bash
#!/bin/bash
# instance-automation.sh

# Function to backup instances
backup_instances() {
  local backup_dir="instance-backups-$(date +%Y%m%d)"
  mkdir -p "$backup_dir"
  
  instances=$(virak-cli --output json instance list | jq -r '.data[].id')
  
  for instance in $instances; do
    # Create snapshot
    snapshot_name="auto-backup-$(date +%Y%m%d%H%M%S)"
    virak-cli instance snapshot create \
      --instanceId "$instance" \
      --name "$snapshot_name"
    
    # Export instance configuration
    virak-cli --output json instance show --instanceId "$instance" > "$backup_dir/$instance.json"
    
    echo "Instance $instance backed up"
  done
  
  echo "All instances backed up to $backup_dir"
}

# Function to monitor instance health
monitor_instances() {
  instances=$(virak-cli --output json instance list | jq -r '.data[].id')
  
  for instance in $instances; do
    status=$(virak-cli --output json instance show --instanceId "$instance" | jq -r '.data.status')
    
    if [ "$status" != "RUNNING" ]; then
      echo "WARNING: Instance $instance is not running (status: $status)"
    fi
  done
}

# Function to clean up old snapshots
cleanup_old_snapshots() {
  local days_to_keep=${1:-30}
  local cutoff_date=$(date -d "$days_to_keep days ago" +%s)
  
  instances=$(virak-cli --output json instance list | jq -r '.data[].id')
  
  for instance in $instances; do
    snapshots=$(virak-cli --output json instance snapshot list --instanceId "$instance" | jq -r '.data[] | select(.created_at < '$cutoff_date') | .id')
    
    for snapshot in $snapshots; do
      virak-cli instance snapshot delete \
        --instanceId "$instance" \
        --snapshotId "$snapshot"
      
      echo "Deleted old snapshot $snapshot for instance $instance"
    done
  done
}
```

## Troubleshooting

### Common Issues

1. **Instance Creation Fails**
   ```
   Error: Instance creation failed
   ```
   Solution: Check service offering, VM image, and network availability

2. **Instance Stuck in Creating State**
   ```
   Status: CREATING
   ```
   Solution: Check resource availability and wait for provisioning to complete

3. **Instance Not Accessible**
   ```
   Error: Connection timeout
   ```
   Solution: Check network configuration and security groups

4. **Snapshot Creation Fails**
   ```
   Error: Snapshot creation failed
   ```
   Solution: Ensure instance is in a stable state and has sufficient disk space

### Debug Commands

```bash
# Enable debug logging
virak-cli --debug instance list

# Check instance status
virak-cli instance show --instanceId <instance-id>

# View instance metrics
virak-cli instance metrics --instanceId <instance-id>

# Check available resources
virak-cli instance service-offering list
virak-cli zone resources
```

For more troubleshooting help, see our [Troubleshooting Guide](../troubleshooting/common-issues.md).