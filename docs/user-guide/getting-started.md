# Getting Started Tutorial

This tutorial will guide you through the basics of using Virak CLI to manage your Virak Cloud resources.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Step 1: Installation and Authentication](#step-1-installation-and-authentication)
- [Step 2: Exploring Available Resources](#step-2-exploring-available-resources)
- [Step 3: Creating Your First Instance](#step-3-creating-your-first-instance)
- [Step 4: Managing Networks](#step-4-managing-networks)
- [Step 5: Working with Storage](#step-5-working-with-storage)
- [Step 6: Cleaning Up Resources](#step-6-cleaning-up-resources)
- [Next Steps](#next-steps)

## Prerequisites

Before starting this tutorial, make sure you have:

1. **Virak CLI installed** - See [Installation Guide](installation.md)
2. **Virak Cloud account** - Sign up at [Virak Cloud](https://virakcloud.com/)
3. **Basic command-line knowledge** - Familiarity with terminal/shell commands

## Step 1: Installation and Authentication

First, verify that Virak CLI is installed and authenticate with your account:

```bash
# Check if Virak CLI is installed
virak-cli --version

# Authenticate with Virak Cloud
virak-cli login
```

This will open your browser for OAuth authentication. Follow the prompts to log in and authorize the CLI.

After authentication, verify it worked by listing available zones:

```bash
# List all available zones
virak-cli zone list
```

You should see output similar to:

```
+-----------+----------------+-------------+----------------+
|   ID      |      NAME      |    STATUS   |   LOCATION     |
+-----------+----------------+-------------+----------------+
| zone-123  | tehran-1       |     UP      |    Tehran      |
| zone-456  | frankfurt-1    |     UP      |   Frankfurt    |
+-----------+----------------+-------------+----------------+
```

Set a default zone to avoid specifying it for every command:

```bash
# Set default zone (replace with your preferred zone)
export VIRAK_ZONE_ID="zone-123"
export VIRAK_ZONE_NAME="tehran-1"
```

Or add it to your configuration file (`~/.virak-cli.yaml`):

```yaml
default:
  zoneId: "zone-123"
  zoneName: "tehran-1"
```

## Step 2: Exploring Available Resources

Let's explore what resources are available in your chosen zone:

```bash
# List available VM images
virak-cli instance vm-image list

# List available service offerings (instance sizes)
virak-cli instance service-offering list

# List existing instances
virak-cli instance list

# List existing networks
virak-cli network list
```

### Understanding VM Images

VM images are pre-configured operating system templates. Common types include:
- **Linux distributions**: Ubuntu, CentOS, Debian
- **Windows Server**: Various Windows Server versions
- **Custom images**: Your own uploaded images

### Understanding Service Offerings

Service offerings define the resources allocated to your instance:
- **CPU**: Number of virtual CPUs
- **Memory**: RAM allocation
- **Disk**: Storage size and type

Choose a service offering based on your workload requirements.

## Step 3: Creating Your First Instance

Now let's create a virtual machine instance:

```bash
# First, note the IDs from the previous commands
# For example:
# VM Image ID: img-ubuntu-20.04
# Service Offering ID: offering-small-1cpu-2gb

# Create an instance
virak-cli instance create \
  --name my-first-instance \
  --vm-image-id img-ubuntu-20.04 \
  --service-offering-id offering-small-1cpu-2gb
```

Wait for the instance to be created, then check its status:

```bash
# Get the instance ID from the create command output
INSTANCE_ID="instance-abc123"

# Check instance status
virak-cli instance show --id $INSTANCE_ID

# List all instances to see yours
virak-cli instance list
```

Start the instance:

```bash
# Start the instance
virak-cli instance start --id $INSTANCE_ID

# Wait a moment, then check status again
virak-cli instance show --id $INSTANCE_ID
```

The instance should now be running. You can access it using SSH if you've configured SSH keys or the console through the Virak Cloud Panel.

### Instance Lifecycle Management

Try these common instance operations:

```bash
# Reboot the instance
virak-cli instance reboot --id $INSTANCE_ID

# Stop the instance
virak-cli instance stop --id $INSTANCE_ID

# Start it again
virak-cli instance start --id $INSTANCE_ID
```

## Step 4: Managing Networks

Instances need network connectivity. Let's create and configure a network:

```bash
# Create a new L3 network
virak-cli network create l3 \
  --name my-network \
  --display-text "My Tutorial Network"

# Note the network ID from the output
NETWORK_ID="network-def456"

# List networks to see your new network
virak-cli network list

# Connect your instance to the network
virak-cli network instance connect \
  --instance-id $INSTANCE_ID \
  --network-id $NETWORK_ID

# List instances in the network
virak-cli network instance list --network-id $NETWORK_ID
```

### Configuring Firewall Rules

Secure your instance with firewall rules:

```bash
# Allow SSH access (port 22)
virak-cli network firewall ipv4 create \
  --network-id $NETWORK_ID \
  --protocol tcp \
  --start-port 22 \
  --end-port 22

# Allow HTTP access (port 80)
virak-cli network firewall ipv4 create \
  --network-id $NETWORK_ID \
  --protocol tcp \
  --start-port 80 \
  --end-port 80

# Allow HTTPS access (port 443)
virak-cli network firewall ipv4 create \
  --network-id $NETWORK_ID \
  --protocol tcp \
  --start-port 443 \
  --end-port 443

# List firewall rules
virak-cli network firewall ipv4 list --network-id $NETWORK_ID
```

## Step 5: Working with Storage

Virak Cloud provides object storage through buckets. Let's create and manage storage:

```bash
# Create a new bucket
virak-cli bucket create \
  --name my-tutorial-bucket \
  --policy Private

# List all buckets
virak-cli bucket list

# Get bucket details
virak-cli bucket show --bucketId my-tutorial-bucket
```

### Using with S3-Compatible Tools

Your bucket is compatible with S3 tools. Here's how to use it with common tools:

```bash
# Using s3cmd (install first: pip install s3cmd)
s3cmd --host=https://s3.virakcloud.com \
  --access-key=$(virak-cli bucket show --bucketId my-tutorial-bucket | grep AccessKey) \
  --secret-key=$(virak-cli bucket show --bucketId my-tutorial-bucket | grep SecretKey) \
  mb s3://my-tutorial-bucket

# Using AWS CLI
aws s3 --endpoint-url=https://s3.virakcloud.com mb s3://my-tutorial-bucket
```

### Managing Instance Volumes

You can also add additional storage to your instances:

```bash
# List volume service offerings
virak-cli instance volume service-offering list

# Create a new volume
virak-cli instance volume create \
  --name my-data-volume \
  --service-offering-id volume-ssd-10gb

# Note the volume ID
VOLUME_ID="volume-ghi789"

# Attach the volume to your instance
virak-cli instance volume attach \
  --instance-id $INSTANCE_ID \
  --volume-id $VOLUME_ID

# List volumes
virak-cli instance volume list
```

## Step 6: Cleaning Up Resources

When you're done with the tutorial, clean up your resources to avoid charges:

```bash
# Detach and delete the volume
virak-cli instance volume detach \
  --instance-id $INSTANCE_ID \
  --volume-id $VOLUME_ID

virak-cli instance volume delete --id $VOLUME_ID

# Stop and delete the instance
virak-cli instance stop --id $INSTANCE_ID
virak-cli instance delete --id $INSTANCE_ID

# Delete firewall rules
virak-cli network firewall ipv4 delete \
  --network-id $NETWORK_ID \
  --id $(virak-cli network firewall ipv4 list --network-id $NETWORK_ID | grep -E '^\|' | grep tcp | awk '{print $2}' | head -1)

# Disconnect instance from network
virak-cli network instance disconnect \
  --instance-id $INSTANCE_ID \
  --network-id $NETWORK_ID

# Delete the network
virak-cli network delete --id $NETWORK_ID

# Delete the bucket
virak-cli bucket delete --bucketId my-tutorial-bucket
```

## Next Steps

Congratulations! You've completed the basic tutorial. Here are some suggested next steps:

1. **Explore Advanced Features**
   - [Deploy a Web Application](tutorials/deploying-web-app.md)
   - [Set Up Containers](tutorials/setting-up-containers.md)
   - [Managing Networks](tutorials/managing-networks.md)

2. **Learn About Kubernetes**
   - [Kubernetes Cluster Management](../reference/cluster.md)
   - Deploy containerized applications

3. **Explore DNS Management**
   - [DNS Commands Reference](../reference/dns.md)
   - Configure custom domains

4. **Advanced Configuration**
   - [Configuration Guide](configuration.md)
   - [Command Reference](../reference/)

5. **Automation**
   - Use Virak CLI in scripts
   - Integrate with CI/CD pipelines

## Common Tasks Quick Reference

```bash
# List all resources in a zone
virak-cli zone resources --zone-id $ZONE_ID

# Create multiple instances quickly
for i in {1..3}; do
  virak-cli instance create \
    --name "web-server-$i" \
    --vm-image-id img-ubuntu-20.04 \
    --service-offering-id offering-small-1cpu-2gb
done

# Monitor instance metrics
virak-cli instance metrics --id $INSTANCE_ID

# Create snapshots for backup
virak-cli instance snapshot create \
  --instance-id $INSTANCE_ID \
  --name "backup-$(date +%Y%m%d)"

# List all your snapshots
virak-cli instance snapshot list
```

## Tips and Best Practices

1. **Use Descriptive Names** - Give your resources meaningful names for easy identification
2. **Set Default Zone** - Configure a default zone to avoid repetitive typing
3. **Use Environment Variables** - For automation, use environment variables for sensitive data
4. **Regular Backups** - Create snapshots of important instances regularly
5. **Monitor Resources** - Keep track of your resource usage to avoid unexpected charges
6. **Security First** - Configure firewall rules properly and use SSH keys for authentication

If you run into any issues, check our [Troubleshooting Guide](../troubleshooting/common-issues.md) or [Authentication Issues](../troubleshooting/authentication.md).