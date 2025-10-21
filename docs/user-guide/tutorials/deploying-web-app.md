# Deploying a Web Application

This tutorial walks you through deploying a web application on Virak Cloud using the Virak CLI.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Overview](#overview)
- [Step 1: Create Network Infrastructure](#step-1-create-network-infrastructure)
- [Step 2: Deploy Web Server](#step-2-deploy-web-server)
- [Step 3: Configure Firewall](#step-3-configure-firewall)
- [Step 4: Set Up Load Balancer](#step-4-set-up-load-balancer)
- [Step 5: Configure DNS](#step-5-configure-dns)
- [Step 6: Deploy Application](#step-6-deploy-application)
- [Step 7: Monitor and Scale](#step-7-monitor-and-scale)
- [Cleanup](#cleanup)

## Prerequisites

Before starting this tutorial, make sure you have:

1. **Virak CLI installed and configured** - See [Installation Guide](../installation.md)
2. **Authenticated with Virak Cloud** - See [Authentication Guide](../authentication.md)
3. **Basic knowledge of web servers and networking**
4. **A domain name (optional)** - For DNS configuration

## Overview

In this tutorial, we'll deploy a scalable web application with the following architecture:

```
Internet → Load Balancer → Web Servers → Database
                ↓
              DNS (example.com)
```

We'll create:
- A private network for our application
- Multiple web server instances
- A load balancer to distribute traffic
- Firewall rules for security
- DNS configuration for domain access

## Step 1: Create Network Infrastructure

First, let's create a private network for our application:

```bash
# Create a private L3 network
virak-cli network create l3 \
  --name web-app-network \
  --display-text "Web Application Network"

# Note the network ID
NETWORK_ID="network-def456"

# Verify network creation
virak-cli network list
```

## Step 2: Deploy Web Server

Now let's deploy our web server instances:

```bash
# List available VM images
virak-cli instance vm-image list

# List available service offerings
virak-cli instance service-offering list

# Choose appropriate image and offering, then create instances
# For this example, we'll use Ubuntu 20.04 and a small instance

# Create first web server
virak-cli instance create \
  --name web-server-1 \
  --vm-image-id img-ubuntu-20.04 \
  --service-offering-id offering-small-1cpu-2gb

# Note the instance ID
WEB_SERVER_1="instance-abc123"

# Create second web server for high availability
virak-cli instance create \
  --name web-server-2 \
  --vm-image-id img-ubuntu-20.04 \
  --service-offering-id offering-small-1cpu-2gb

# Note the instance ID
WEB_SERVER_2="instance-def456"

# Start both instances
virak-cli instance start --id $WEB_SERVER_1
virak-cli instance start --id $WEB_SERVER_2

# Connect instances to our network
virak-cli network instance connect \
  --instance-id $WEB_SERVER_1 \
  --network-id $NETWORK_ID

virak-cli network instance connect \
  --instance-id $WEB_SERVER_2 \
  --network-id $NETWORK_ID

# Wait for instances to be ready
sleep 30

# Check instance status
virak-cli instance show --id $WEB_SERVER_1
virak-cli instance show --id $WEB_SERVER_2
```

## Step 3: Configure Firewall

Let's configure firewall rules to allow web traffic:

```bash
# Allow HTTP traffic (port 80)
virak-cli network firewall ipv4 create \
  --network-id $NETWORK_ID \
  --protocol tcp \
  --start-port 80 \
  --end-port 80 \
  --cidr "0.0.0.0/0"

# Allow HTTPS traffic (port 443)
virak-cli network firewall ipv4 create \
  --network-id $NETWORK_ID \
  --protocol tcp \
  --start-port 443 \
  --end-port 443 \
  --cidr "0.0.0.0/0"

# Allow SSH access for management (restrict to your IP)
virak-cli network firewall ipv4 create \
  --network-id $NETWORK_ID \
  --protocol tcp \
  --start-port 22 \
  --end-port 22 \
  --cidr "YOUR_IP_ADDRESS/32"

# List firewall rules
virak-cli network firewall ipv4 list --network-id $NETWORK_ID
```

## Step 4: Set Up Load Balancer

Now let's set up a load balancer to distribute traffic between our web servers:

```bash
# Create a load balancer
virak-cli network lb create \
  --name web-app-lb \
  --network-id $NETWORK_ID

# Note the load balancer ID
LB_ID="lb-ghi789"

# Assign web servers to the load balancer
virak-cli network lb assign \
  --lb-id $LB_ID \
  --instance-id $WEB_SERVER_1

virak-cli network lb assign \
  --lb-id $LB_ID \
  --instance-id $WEB_SERVER_2

# Get load balancer details including public IP
virak-cli network lb show --id $LB_ID

# Note the public IP address
LB_IP=$(virak-cli network lb show --id $LB_ID | grep "Public IP" | awk '{print $3}')
echo "Load Balancer IP: $LB_IP"

# Check HAProxy status (if using HAProxy)
virak-cli network lb haproxy live --lb-id $LB_ID
```

## Step 5: Configure DNS

If you have a domain name, let's configure DNS to point to our load balancer:

```bash
# Create a DNS domain (if you don't have one)
virak-cli dns domain create --name example.com

# Note the domain ID
DOMAIN_ID="domain-jkl012"

# Create an A record pointing to your load balancer
virak-cli dns record create \
  --domain-name example.com \
  --name @ \
  --type A \
  --content $LB_IP \
  --ttl 300

# Create a www subdomain
virak-cli dns record create \
  --domain-name example.com \
  --name www \
  --type A \
  --content $LB_IP \
  --ttl 300

# List DNS records
virak-cli dns record list --domain-name example.com
```

## Step 6: Deploy Application

Now let's deploy our web application to the servers:

```bash
# Get instance details for SSH access
virak-cli instance show --id $WEB_SERVER_1
virak-cli instance show --id $WEB_SERVER_2

# Connect to the first instance via SSH (use the IP from the show command)
# ssh root@<instance-ip>

# Once connected, install and configure your web application
# Example for a simple Nginx web server:

# Update system
apt update && apt upgrade -y

# Install Nginx
apt install nginx -y

# Create a simple web page
cat > /var/www/html/index.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>My Web App</title>
</head>
<body>
    <h1>Hello from Web Server 1!</h1>
    <p>This page is served by Virak Cloud.</p>
</body>
</html>
EOF

# Configure Nginx
cat > /etc/nginx/sites-available/default << EOF
server {
    listen 80 default_server;
    listen [::]:80 default_server;

    root /var/www/html;
    index index.html;

    server_name _;

    location / {
        try_files \$uri \$uri/ =404;
    }

    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;
}
EOF

# Restart Nginx
systemctl restart nginx
systemctl enable nginx

# Exit SSH
exit

# Repeat similar steps for the second server
# ssh root@<second-instance-ip>
# ... (same installation steps, but change the content to "Hello from Web Server 2!")
```

## Step 7: Monitor and Scale

Let's monitor our deployment and set up scaling:

```bash
# Check load balancer status
virak-cli network lb haproxy live --lb-id $LB_ID

# View HAProxy logs
virak-cli network lb haproxy log --lb-id $LB_ID

# Check instance metrics
virak-cli instance metrics --id $WEB_SERVER_1
virak-cli instance metrics --id $WEB_SERVER_2

# Create a snapshot for backup
virak-cli instance snapshot create \
  --instance-id $WEB_SERVER_1 \
  --name "web-server-1-backup-$(date +%Y%m%d)"

# If you need to scale up, create more instances
virak-cli instance create \
  --name web-server-3 \
  --vm-image-id img-ubuntu-20.04 \
  --service-offering-id offering-small-1cpu-2gb

# Note the new instance ID
WEB_SERVER_3="instance-mno345"

# Start and connect to network
virak-cli instance start --id $WEB_SERVER_3
virak-cli network instance connect \
  --instance-id $WEB_SERVER_3 \
  --network-id $NETWORK_ID

# Add to load balancer
virak-cli network lb assign \
  --lb-id $LB_ID \
  --instance-id $WEB_SERVER_3
```

## Testing Your Deployment

```bash
# Test the load balancer
curl http://$LB_IP

# Test with domain name (if DNS is configured)
curl http://example.com

# Test multiple times to see load balancing
for i in {1..10}; do
  echo "Request $i:"
  curl -s http://$LB_IP | grep "Web Server"
  echo "---"
done
```

## Cleanup

When you're done with the tutorial, clean up your resources:

```bash
# Remove instances from load balancer
virak-cli network lb deassign \
  --lb-id $LB_ID \
  --instance-id $WEB_SERVER_1

virak-cli network lb deassign \
  --lb-id $LB_ID \
  --instance-id $WEB_SERVER_2

virak-cli network lb deassign \
  --lb-id $LB_ID \
  --instance-id $WEB_SERVER_3

# Delete load balancer
virak-cli network lb delete --id $LB_ID

# Stop and delete instances
virak-cli instance stop --id $WEB_SERVER_1
virak-cli instance delete --id $WEB_SERVER_1

virak-cli instance stop --id $WEB_SERVER_2
virak-cli instance delete --id $WEB_SERVER_2

virak-cli instance stop --id $WEB_SERVER_3
virak-cli instance delete --id $WEB_SERVER_3

# Delete firewall rules
RULE_IDS=$(virak-cli network firewall ipv4 list --network-id $NETWORK_ID | grep -E '^\|' | grep -v ID | awk '{print $2}')
for rule_id in $RULE_IDS; do
  virak-cli network firewall ipv4 delete \
    --network-id $NETWORK_ID \
    --id $rule_id
done

# Delete network
virak-cli network delete --id $NETWORK_ID

# Delete DNS records (if created)
virak-cli dns record delete --domain-name example.com --name @
virak-cli dns record delete --domain-name example.com --name www

# Delete domain (if created)
virak-cli dns domain delete --name example.com
```

## Next Steps

Congratulations! You've successfully deployed a scalable web application on Virak Cloud. Here are some next steps to consider:

1. **Set up Monitoring**
   - Implement application monitoring
   - Set up alerting for critical metrics

2. **Implement CI/CD**
   - Automate deployment using Virak CLI in your pipeline
   - Set up blue-green deployments

3. **Add Database**
   - Deploy a database instance
   - Configure connection pooling

4. **Implement SSL/TLS**
   - Obtain SSL certificates
   - Configure HTTPS on your load balancer

5. **Backup Strategy**
   - Set up automated backups
   - Implement disaster recovery

## Automation Script

Here's a complete script to automate the deployment:

```bash
#!/bin/bash
# deploy-web-app.sh

set -e

# Configuration
APP_NAME="web-app"
INSTANCE_COUNT=2
VM_IMAGE="img-ubuntu-20.04"
SERVICE_OFFERING="offering-small-1cpu-2gb"

# Create network
echo "Creating network..."
NETWORK_OUTPUT=$(virak-cli network create l3 --name $APP_NAME-network --display-text "$APP_NAME Network")
NETWORK_ID=$(echo $NETWORK_OUTPUT | grep -o 'network-[a-zA-Z0-9]*')

# Create instances
echo "Creating instances..."
INSTANCE_IDS=()
for i in $(seq 1 $INSTANCE_COUNT); do
  INSTANCE_OUTPUT=$(virak-cli instance create \
    --name $APP_NAME-server-$i \
    --vm-image-id $VM_IMAGE \
    --service-offering-id $SERVICE_OFFERING)
  INSTANCE_ID=$(echo $INSTANCE_OUTPUT | grep -o 'instance-[a-zA-Z0-9]*')
  INSTANCE_IDS+=($INSTANCE_ID)
  
  virak-cli instance start --id $INSTANCE_ID
  virak-cli network instance connect --instance-id $INSTANCE_ID --network-id $NETWORK_ID
done

# Configure firewall
echo "Configuring firewall..."
virak-cli network firewall ipv4 create --network-id $NETWORK_ID --protocol tcp --start-port 80 --end-port 80 --cidr "0.0.0.0/0"
virak-cli network firewall ipv4 create --network-id $NETWORK_ID --protocol tcp --start-port 443 --end-port 443 --cidr "0.0.0.0/0"

# Create load balancer
echo "Creating load balancer..."
LB_OUTPUT=$(virak-cli network lb create --name $APP_NAME-lb --network-id $NETWORK_ID)
LB_ID=$(echo $LB_OUTPUT | grep -o 'lb-[a-zA-Z0-9]*')

# Assign instances to load balancer
for instance_id in "${INSTANCE_IDS[@]}"; do
  virak-cli network lb assign --lb-id $LB_ID --instance-id $instance_id
done

# Get load balancer IP
LB_IP=$(virak-cli network lb show --id $LB_ID | grep "Public IP" | awk '{print $3}')

echo "Deployment complete!"
echo "Load Balancer IP: $LB_IP"
echo "Network ID: $NETWORK_ID"
echo "Load Balancer ID: $LB_ID"
echo "Instance IDs: ${INSTANCE_IDS[*]}"
```

This script provides a starting point for automating your web application deployments on Virak Cloud.