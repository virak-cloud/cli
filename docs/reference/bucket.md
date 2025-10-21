# Bucket Commands Reference

This page provides comprehensive documentation for bucket (object storage) commands in Virak CLI.

## Table of Contents

- [Overview](#overview)
- [Commands](#commands)
  - [bucket create](#bucket-create)
  - [bucket delete](#bucket-delete)
  - [bucket events](#bucket-events)
  - [bucket list](#bucket-list)
  - [bucket show](#bucket-show)
  - [bucket update](#bucket-update)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

Bucket commands allow you to manage object storage buckets in Virak Cloud. Buckets provide S3-compatible object storage for storing files, media, backups, and other data.

All bucket commands require a zone ID, which can be provided via the `--zoneId` flag or by setting a default zone in the configuration.

## Commands

### bucket create

Create a new object storage bucket.

#### Syntax

```bash
virak-cli bucket create [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--name` | Bucket name (must be unique across the zone) |
| `--policy` | Bucket policy (Private, Public) |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Create a private bucket
virak-cli bucket create --name my-private-bucket --policy Private

# Create a public bucket
virak-cli bucket create --name my-public-bucket --policy Public

# Create bucket in specific zone
virak-cli bucket create --name my-bucket --policy Private --zoneId zone-456
```

#### Output

```
Bucket 'my-bucket' created successfully.
Bucket ID: my-bucket
Policy: Private
Zone: zone-123
```

### bucket delete

Delete an object storage bucket.

**Warning:** This action is irreversible. All data in the bucket will be permanently deleted.

#### Syntax

```bash
virak-cli bucket delete [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--bucketId` | Bucket name/ID to delete |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--force` | Skip confirmation prompt |

#### Examples

```bash
# Delete a bucket (with confirmation)
virak-cli bucket delete --bucketId my-bucket

# Delete a bucket without confirmation
virak-cli bucket delete --bucketId my-bucket --force

# Delete bucket in specific zone
virak-cli bucket delete --bucketId my-bucket --zoneId zone-456
```

#### Output

```
Bucket 'my-bucket' deleted successfully.
```

### bucket events

View events related to a bucket.

#### Syntax

```bash
virak-cli bucket events [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--bucketId` | Bucket name/ID |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--limit` | Maximum number of events to return (default: 50) |
| `--page` | Page number for pagination (default: 1) |

#### Examples

```bash
# View bucket events
virak-cli bucket events --bucketId my-bucket

# View last 100 events
virak-cli bucket events --bucketId my-bucket --limit 100

# View second page of events
virak-cli bucket events --bucketId my-bucket --page 2
```

#### Output

```
+------------+----------------------+----------------+----------------+
|    ID      |        TYPE          |   DESCRIPTION  |    TIMESTAMP    |
+------------+----------------------+----------------+----------------+
| event-123  | BUCKET_CREATED       | Bucket created | 2023-10-19...  |
| event-124  | OBJECT_UPLOADED      | File uploaded  | 2023-10-19...  |
+------------+----------------------+----------------+----------------+
```

### bucket list

List all buckets in a zone.

#### Syntax

```bash
virak-cli bucket list [flags]
```

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--filter` | Filter results (format: key=value) |
| `--limit` | Maximum number of buckets to return (default: 50) |
| `--page` | Page number for pagination (default: 1) |
| `--sort` | Sort field (name, created, policy) |
| `--order` | Sort order (asc, desc) |

#### Examples

```bash
# List all buckets
virak-cli bucket list

# List buckets in specific zone
virak-cli bucket list --zoneId zone-456

# List only public buckets
virak-cli bucket list --filter policy=Public

# List buckets sorted by name
virak-cli bucket list --sort name --order asc

# List buckets in JSON format
virak-cli --output json bucket list
```

#### Output

```
+------------+----------------+----------+----------------+
|    NAME    |     POLICY     |  STATUS  |    CREATED     |
+------------+----------------+----------+----------------+
| my-bucket  |     Private    |  ACTIVE  | 2023-10-19...  |
| public-bkt |     Public     |  ACTIVE  | 2023-10-18...  |
+------------+----------------+----------+----------------+
```

### bucket show

Show detailed information about a bucket.

#### Syntax

```bash
virak-cli bucket show [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--bucketId` | Bucket name/ID |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |

#### Examples

```bash
# Show bucket details
virak-cli bucket show --bucketId my-bucket

# Show bucket details in JSON format
virak-cli --output json bucket show --bucketId my-bucket

# Show bucket details in specific zone
virak-cli bucket show --bucketId my-bucket --zoneId zone-456
```

#### Output

```
Bucket Details:
  Name: my-bucket
  ID: my-bucket
  Policy: Private
  Status: ACTIVE
  Zone: zone-123
  Created: 2023-10-19 10:30:00 UTC
  Updated: 2023-10-19 10:30:00 UTC
  Size: 1.2 GB
  Object Count: 42
  Access Key: AKIAIOSFODNN7EXAMPLE
  Secret Key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
  Endpoint: https://s3.virakcloud.com
```

### bucket update

Update bucket properties.

#### Syntax

```bash
virak-cli bucket update [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--bucketId` | Bucket name/ID |

#### Optional Flags

| Flag | Description |
|------|-------------|
| `--zoneId` | Zone ID (overrides default zone) |
| `--policy` | New bucket policy (Private, Public) |

#### Examples

```bash
# Change bucket policy to public
virak-cli bucket update --bucketId my-bucket --policy Public

# Change bucket policy to private
virak-cli bucket update --bucketId my-bucket --policy Private
```

#### Output

```
Bucket 'my-bucket' updated successfully.
Policy: Public
```

## Examples

### Complete Bucket Lifecycle

```bash
# 1. Create a bucket
virak-cli bucket create --name my-app-storage --policy Private

# 2. List buckets to verify
virak-cli bucket list

# 3. Show bucket details
virak-cli bucket show --bucketId my-app-storage

# 4. Update bucket policy
virak-cli bucket update --bucketId my-app-storage --policy Public

# 5. View bucket events
virak-cli bucket events --bucketId my-app-storage

# 6. Delete bucket (when done)
virak-cli bucket delete --bucketId my-app-storage --force
```

### Using with S3-Compatible Tools

```bash
# Get bucket credentials
virak-cli bucket show --bucketId my-bucket

# Use with s3cmd
s3cmd --host=https://s3.virakcloud.com \
  --access-key=AKIAIOSFODNN7EXAMPLE \
  --secret-key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY \
  mb s3://my-bucket

# Use with AWS CLI
aws s3 --endpoint-url=https://s3.virakcloud.com \
  mb s3://my-bucket

# Upload files
aws s3 --endpoint-url=https://s3.virakcloud.com \
  cp local-file.txt s3://my-bucket/
```

### Managing Multiple Buckets

```bash
# Create multiple buckets
for bucket in data backups logs; do
  virak-cli bucket create --name "my-app-$bucket" --policy Private
done

# List all buckets
virak-cli bucket list

# Update policy for logs bucket
virak-cli bucket update --bucketId my-app-logs --policy Public

# Delete all buckets
for bucket in data backups logs; do
  virak-cli bucket delete --bucketId "my-app-$bucket" --force
done
```

## Best Practices

### Bucket Naming

1. **Use Unique Names**: Bucket names must be unique across the entire zone
2. **Use DNS-Compatible Names**: Use lowercase letters, numbers, and hyphens
3. **Avoid Special Characters**: Don't use underscores, spaces, or special characters
4. **Use Descriptive Names**: Include project or environment information

```bash
# Good names
my-app-prod-data
company-backups-2023
user-media-storage

# Bad names
My App Data (spaces)
bucket_123 (underscores)
special@chars# (special characters)
```

### Security

1. **Use Private Policy by Default**: Only make buckets public when necessary
2. **Rotate Credentials**: Regularly update access keys and secrets
3. **Monitor Events**: Keep track of bucket events for security auditing
4. **Use IAM Policies**: Implement fine-grained access control when available

```bash
# Monitor bucket events
virak-cli bucket events --bucketId sensitive-data --limit 100

# Check bucket policy
virak-cli bucket show --bucketId my-bucket --output json | jq '.policy'
```

### Performance

1. **Choose Appropriate Zones**: Select zones closest to your users
2. **Monitor Bucket Size**: Keep track of storage usage
3. **Use Lifecycle Policies**: Implement automatic cleanup when available

```bash
# Check bucket size
virak-cli bucket show --bucketId my-bucket | grep Size

# List buckets by size
virak-cli bucket list --sort created --order desc
```

### Automation

```bash
#!/bin/bash
# backup-buckets.sh

# List all buckets
BUCKETS=$(virak-cli --output json bucket list | jq -r '.data[].name')

# Create backup of each bucket
for bucket in $BUCKETS; do
  echo "Backing up bucket: $bucket"
  
  # Get bucket details
  virak-cli bucket show --bucketId "$bucket" > "backup-$bucket-$(date +%Y%m%d).json"
  
  # Export bucket data using s3cmd
  s3cmd --host=https://s3.virakcloud.com \
    --access-key=$(virak-cli bucket show --bucketId "$bucket" | grep AccessKey | awk '{print $2}') \
    --secret-key=$(virak-cli bucket show --bucketId "$bucket" | grep SecretKey | awk '{print $2}') \
    sync s3://$bucket/ "backup-$bucket/"
done
```

## Troubleshooting

### Common Issues

1. **Bucket Name Already Exists**
   ```
   Error: Bucket name 'my-bucket' already exists
   ```
   Solution: Choose a different name or delete the existing bucket

2. **Bucket Not Found**
   ```
   Error: Bucket 'my-bucket' not found
   ```
   Solution: Check bucket name and zone ID

3. **Permission Denied**
   ```
   Error: Access denied to bucket 'my-bucket'
   ```
   Solution: Check your token permissions and bucket policy

4. **Bucket Not Empty**
   ```
   Error: Cannot delete non-empty bucket
   ```
   Solution: Empty the bucket first using S3-compatible tools

### Debug Commands

```bash
# Enable debug logging
virak-cli --debug bucket list

# Check bucket status
virak-cli bucket show --bucketId my-bucket

# Verify zone
virak-cli zone list

# Test API connection
curl -H "Authorization: Bearer $VIRAK_TOKEN" \
  "$VIRAK_API_URL/v1/buckets"
```

For more troubleshooting help, see our [Troubleshooting Guide](../troubleshooting/common-issues.md).