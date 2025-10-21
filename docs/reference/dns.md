# DNS Commands Reference

This page provides comprehensive documentation for DNS commands in Virak CLI.

## Table of Contents

- [Overview](#overview)
- [Commands](#commands)
  - [dns domain create](#dns-domain-create)
  - [dns domain delete](#dns-domain-delete)
  - [dns domain list](#dns-domain-list)
  - [dns domain show](#dns-domain-show)
  - [dns events](#dns-events)
  - [dns record create](#dns-record-create)
  - [dns record delete](#dns-record-delete)
  - [dns record list](#dns-record-list)
  - [dns record update](#dns-record-update)
- [DNS Record Types](#dns-record-types)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

DNS commands allow you to manage DNS domains and records in Virak Cloud. These commands provide complete control over your DNS configuration, including domain management and various types of DNS records.

The DNS commands are organized into two main categories:
- Domain commands: Manage DNS domains
- Record commands: Manage DNS records within domains

## Commands

### dns domain create

Create a new DNS domain.

#### Syntax

```bash
virak-cli dns domain create [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--domain` | Domain name to create |

#### Examples

```bash
# Create a new domain
virak-cli dns domain create --domain example.com

# Create a subdomain
virak-cli dns domain create --domain app.example.com
```

#### Output

```
Domain creation initiated successfully, consider to setup nameservers if not done already.
```

### dns domain delete

Delete a DNS domain.

**Warning:** This action is irreversible. All DNS records in the domain will be permanently deleted.

#### Syntax

```bash
virak-cli dns domain delete [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--domain` | Domain name to delete |

#### Examples

```bash
# Delete a domain
virak-cli dns domain delete --domain example.com
```

#### Output

```
Domain 'example.com' deleted successfully.
```

### dns domain list

List all DNS domains.

#### Syntax

```bash
virak-cli dns domain list
```

#### Examples

```bash
# List all domains
virak-cli dns domain list

# List domains in JSON format
virak-cli --output json dns domain list
```

#### Output

```
+------------------+---------+
|      Domain      | Status  |
+------------------+---------+
|    example.com   |  ACTIVE |
|  app.example.com |  ACTIVE |
|  api.example.com |  ACTIVE |
+------------------+---------+
```

### dns domain show

Show detailed information about a DNS domain.

#### Syntax

```bash
virak-cli dns domain show [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--domain` | Domain name to show |

#### Examples

```bash
# Show domain details
virak-cli dns domain show --domain example.com

# Show domain details in JSON format
virak-cli --output json dns domain show --domain example.com
```

#### Output

```
Domain Details:
  Name: example.com
  Status: ACTIVE
  Created: 2023-10-19 10:30:00 UTC
  Updated: 2023-10-19 10:30:00 UTC
  Nameservers:
    - ns1.virakcloud.com
    - ns2.virakcloud.com
  Record Count: 5
```

### dns events

View DNS-related events.

#### Syntax

```bash
virak-cli dns events
```

#### Examples

```bash
# View DNS events
virak-cli dns events
```

#### Output

```
+------------+----------------------+----------------+----------------+
|    ID      |        TYPE          |   DESCRIPTION  |    TIMESTAMP    |
+------------+----------------------+----------------+----------------+
| event-123  | DOMAIN_CREATED       | Domain created | 2023-10-19...  |
| event-124  | RECORD_CREATED       | Record created | 2023-10-19...  |
+------------+----------------------+----------------+----------------+
```

### dns record create

Create a new DNS record for a domain.

#### Syntax

```bash
virak-cli dns record create [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--domain` | Domain name |
| `--record` | DNS record name (e.g., www, mail) |
| `--type` | DNS record type |
| `--content` | DNS record content/value |

#### Optional Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--ttl` | Time To Live in seconds | `3600` |
| `--priority` | MX/SRV record priority | |
| `--weight` | SRV record weight | |
| `--port` | SRV record port number | |
| `--flags` | CAA record flags | |
| `--tag` | CAA record tag | |
| `--license` | TLSA record certificate usage | |
| `--choicer` | TLSA record selector | |
| `--match` | TLSA record matching type | |

#### Examples

```bash
# Create an A record
virak-cli dns record create --domain example.com --record www --type A --content 192.0.2.1

# Create an AAAA record
virak-cli dns record create --domain example.com --record www --type AAAA --content 2001:db8::1

# Create a CNAME record
virak-cli dns record create --domain example.com --record blog --type CNAME --content example.com

# Create an MX record
virak-cli dns record create --domain example.com --record @ --type MX --content mail.example.com --priority 10

# Create a TXT record
virak-cli dns record create --domain example.com --record _dmarc --type TXT --content "v=DMARC1; p=quarantine"

# Create an SRV record
virak-cli dns record create --domain example.com --record _sip._tcp --type SRV --content sip.example.com --priority 10 --weight 60 --port 5060

# Create a CAA record
virak-cli dns record create --domain example.com --record @ --type CAA --flags 0 --tag issue --content "letsencrypt.org"

# Create a TLSA record
virak-cli dns record create --domain example.com --record _443._tcp --type TLSA --license 3 --choicer 1 --match 1 --content "abcdef1234567890"
```

#### Output

```
Success
```

### dns record delete

Delete a DNS record.

#### Syntax

```bash
virak-cli dns record delete [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--domain` | Domain name |
| `--record` | DNS record name |
| `--type` | DNS record type |

#### Examples

```bash
# Delete an A record
virak-cli dns record delete --domain example.com --record www --type A

# Delete a TXT record
virak-cli dns record delete --domain example.com --record _dmarc --type TXT
```

#### Output

```
Record 'www' (A) deleted successfully from domain 'example.com'.
```

### dns record list

List all DNS records for a domain.

#### Syntax

```bash
virak-cli dns record list [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--domain` | Domain name |

#### Examples

```bash
# List all records for a domain
virak-cli dns record list --domain example.com

# List records in JSON format
virak-cli --output json dns record list --domain example.com
```

#### Output

```
+-------+-------+--------+------------------+--------+
| NAME  | TYPE  | TTL    |     CONTENT      | STATUS |
+-------+-------+--------+------------------+--------+
| @     | SOA   | 3600   | ns1.virakcloud.com| ACTIVE |
| @     | NS    | 3600   | ns1.virakcloud.com| ACTIVE |
| @     | NS    | 3600   | ns2.virakcloud.com| ACTIVE |
| www   | A     | 3600   | 192.0.2.1       | ACTIVE |
| mail  | MX    | 3600   | mail.example.com | ACTIVE |
+-------+-------+--------+------------------+--------+
```

### dns record update

Update an existing DNS record.

#### Syntax

```bash
virak-cli dns record update [flags]
```

#### Required Flags

| Flag | Description |
|------|-------------|
| `--domain` | Domain name |
| `--record` | DNS record name |
| `--type` | DNS record type |
| `--content` | New DNS record content/value |

#### Optional Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--ttl` | New Time To Live in seconds | `3600` |
| `--priority` | New MX/SRV record priority | |
| `--weight` | New SRV record weight | |
| `--port` | New SRV record port number | |
| `--flags` | New CAA record flags | |
| `--tag` | New CAA record tag | |
| `--license` | New TLSA record certificate usage | |
| `--choicer` | New TLSA record selector | |
| `--match` | New TLSA record matching type | |

#### Examples

```bash
# Update an A record
virak-cli dns record update --domain example.com --record www --type A --content 192.0.2.100

# Update MX record priority
virak-cli dns record update --domain example.com --record @ --type MX --content mail.example.com --priority 20

# Update TXT record content
virak-cli dns record update --domain example.com --record _dmarc --type TXT --content "v=DMARC1; p=reject"
```

#### Output

```
Record 'www' (A) updated successfully in domain 'example.com'.
```

## DNS Record Types

### A Record
Maps a domain name to an IPv4 address.

```bash
virak-cli dns record create --domain example.com --record www --type A --content 192.0.2.1
```

### AAAA Record
Maps a domain name to an IPv6 address.

```bash
virak-cli dns record create --domain example.com --record www --type AAAA --content 2001:db8::1
```

### CNAME Record
Maps a domain name to another domain name (alias).

```bash
virak-cli dns record create --domain example.com --record blog --type CNAME --content example.com
```

### MX Record
Specifies mail server for the domain.

```bash
virak-cli dns record create --domain example.com --record @ --type MX --content mail.example.com --priority 10
```

### TXT Record
Stores arbitrary text data.

```bash
virak-cli dns record create --domain example.com --record _dmarc --type TXT --content "v=DMARC1; p=quarantine"
```

### NS Record
Delegates a domain to a set of authoritative name servers.

```bash
virak-cli dns record create --domain example.com --record @ --type NS --content ns1.virakcloud.com
```

### SOA Record
Start of Authority record for the domain.

```bash
virak-cli dns record create --domain example.com --record @ --type SOA --content ns1.virakcloud.com
```

### SRV Record
Specifies location of services.

```bash
virak-cli dns record create --domain example.com --record _sip._tcp --type SRV --content sip.example.com --priority 10 --weight 60 --port 5060
```

### CAA Record
Certificate Authority Authorization record.

```bash
virak-cli dns record create --domain example.com --record @ --type CAA --flags 0 --tag issue --content "letsencrypt.org"
```

### TLSA Record
TLSA certificate association record.

```bash
virak-cli dns record create --domain example.com --record _443._tcp --type TLSA --license 3 --choicer 1 --match 1 --content "abcdef1234567890"
```

## Examples

### Complete Domain Setup

```bash
# 1. Create a new domain
virak-cli dns domain create --domain example.com

# 2. List domains to verify
virak-cli dns domain list

# 3. Create basic records
virak-cli dns record create --domain example.com --record www --type A --content 192.0.2.1
virak-cli dns record create --domain example.com --record mail --type A --content 192.0.2.2

# 4. Create MX record
virak-cli dns record create --domain example.com --record @ --type MX --content mail.example.com --priority 10

# 5. Create SPF record
virak-cli dns record create --domain example.com --record @ --type TXT --content "v=spf1 mx -all"

# 6. Create DMARC record
virak-cli dns record create --domain example.com --record _dmarc --type TXT --content "v=DMARC1; p=quarantine"

# 7. List all records
virak-cli dns record list --domain example.com
```

### Subdomain Management

```bash
# Create a subdomain
virak-cli dns domain create --domain api.example.com

# Add records to subdomain
virak-cli dns record create --domain api.example.com --record @ --type A --content 192.0.2.10
virak-cli dns record create --domain api.example.com --record v1 --type CNAME --content api.example.com
```

### Email Configuration

```bash
# Create email-related records
virak-cli dns record create --domain example.com --record @ --type MX --content mail1.example.com --priority 10
virak-cli dns record create --domain example.com --record @ --type MX --content mail2.example.com --priority 20

# Create SPF record
virak-cli dns record create --domain example.com --record @ --type TXT --content "v=spf1 include:_spf.google.com ~all"

# Create DKIM record
virak-cli dns record create --domain example.com --record k1._domainkey --type TXT --content "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQ..."

# Create DMARC record
virak-cli dns record create --domain example.com --record _dmarc --type TXT --content "v=DMARC1; p=reject; rua=mailto:dmarc@example.com"
```

### Web Application Configuration

```bash
# Configure web application with load balancer
virak-cli dns record create --domain example.com --record @ --type A --content 192.0.2.100
virak-cli dns record create --domain example.com --record @ --type A --content 192.0.2.101

# Configure CDN
virak-cli dns record create --domain example.com --record cdn --type CNAME --content cdn.provider.com

# Configure API subdomain
virak-cli dns record create --domain example.com --record api --type CNAME --content api.example.com

# Configure development environment
virak-cli dns record create --domain example.com --record dev --type CNAME --content dev.example.com
```

## Best Practices

### Domain Management

1. **Use Descriptive Names**: Choose domain names that reflect your project or service
2. **Plan Subdomains**: Organize services using subdomains (api, www, blog, etc.)
3. **Monitor Domain Status**: Regularly check domain status and renewal dates
4. **Backup DNS Configuration**: Keep backups of your DNS records

```bash
# Export DNS configuration
virak-cli --output json dns record list --domain example.com > dns-backup-$(date +%Y%m%d).json
```

### Record Management

1. **Use Appropriate TTL**: Set TTL values based on how often records change
2. **Test Changes**: Test DNS changes in a staging environment first
3. **Document Records**: Keep documentation of your DNS configuration
4. **Use CNAME for Aliases**: Use CNAME records instead of multiple A records when possible

```bash
# Short TTL for frequently changing records
virak-cli dns record create --domain example.com --record test --type A --content 192.0.2.1 --ttl 300

# Long TTL for stable records
virak-cli dns record create --domain example.com --record www --type A --content 192.0.2.1 --ttl 86400
```

### Security

1. **Implement SPF/DKIM/DMARC**: Configure email authentication records
2. **Use CAA Records**: Specify authorized certificate authorities
3. **Limit Zone Transfers**: Restrict DNS zone transfers to authorized servers
4. **Monitor for Changes**: Set up alerts for DNS record changes

```bash
# Configure email security
virak-cli dns record create --domain example.com --record @ --type TXT --content "v=spf1 mx -all"
virak-cli dns record create --domain example.com --record _dmarc --type TXT --content "v=DMARC1; p=reject"
virak-cli dns record create --domain example.com --record @ --type CAA --flags 0 --tag issue --content "letsencrypt.org"
```

### Performance

1. **Use Geographically Distributed Servers**: Configure multiple A records for load balancing
2. **Optimize TTL Values**: Balance between performance and flexibility
3. **Use CDN Services**: Offload static content to CDN providers
4. **Monitor DNS Performance**: Track DNS query response times

```bash
# Configure load balancing
virak-cli dns record create --domain example.com --record @ --type A --content 192.0.2.100
virak-cli dns record create --domain example.com --record @ --type A --content 192.0.2.101
virak-cli dns record create --domain example.com --record @ --type A --content 192.0.2.102
```

### Automation

```bash
#!/bin/bash
# dns-management.sh

# Function to create domain with basic records
setup_domain() {
  local domain=$1
  local ip=$2
  
  # Create domain
  virak-cli dns domain create --domain "$domain"
  
  # Create basic records
  virak-cli dns record create --domain "$domain" --record @ --type A --content "$ip"
  virak-cli dns record create --domain "$domain" --record www --type A --content "$ip"
  
  # Create email records
  virak-cli dns record create --domain "$domain" --record @ --type MX --content "mail.$domain" --priority 10
  virak-cli dns record create --domain "$domain" --record @ --type TXT --content "v=spf1 mx -all"
  
  echo "Domain $domain setup complete"
}

# Function to backup all domains
backup_all_domains() {
  local backup_dir="dns-backup-$(date +%Y%m%d)"
  mkdir -p "$backup_dir"
  
  # Get all domains
  domains=$(virak-cli --output json dns domain list | jq -r '.data[].domain')
  
  # Backup each domain
  for domain in $domains; do
    virak-cli --output json dns record list --domain "$domain" > "$backup_dir/$domain.json"
  done
  
  echo "All domains backed up to $backup_dir"
}
```

## Troubleshooting

### Common Issues

1. **Domain Creation Fails**
   ```
   Error: Domain already exists
   ```
   Solution: Choose a different domain name or delete the existing domain

2. **Record Creation Fails**
   ```
   Error: Invalid record content
   ```
   Solution: Check record format and content based on record type

3. **DNS Propagation Delay**
   ```
   DNS record not resolving
   ```
   Solution: Wait for DNS propagation (can take up to 48 hours)

4. **Invalid TTL Value**
   ```
   Error: TTL must be between 300 and 86400
   ```
   Solution: Use TTL value within the allowed range

### Debug Commands

```bash
# Check domain status
virak-cli dns domain show --domain example.com

# List all records
virak-cli dns record list --domain example.com

# Check DNS resolution
dig example.com
dig www.example.com

# Test DNS propagation
dig +trace example.com
```

For more troubleshooting help, see our [Troubleshooting Guide](../troubleshooting/common-issues.md).