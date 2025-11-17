# Managing DNS with Virak CLI

Virak Cloud DNS offers the same domain and record features as the panel plus the record-type coverage described in [dns guide in virakcloud docs](https://docs.virakcloud.com) (A, AAAA, CNAME, MX, TXT, NS, PTR, SRV, TLSA, CAA, SOA). This guide shows the exact command syntax implemented under `virak-cli`.

## Prerequisites

- Authenticated Virak CLI session
- Domain registered with Virak Cloud or delegated to Virak’s nameservers
- API token with DNS permissions

## Domain Operations

```sh
virak-cli dns domain list

virak-cli dns domain create \
  --domain "example.com"

virak-cli dns domain show \
  --domain "example.com"

virak-cli dns domain delete \
  --domain "example.com"
```

`domain show` returns `pending` when delegation has not finished, matching the behavior documented in the panel.

## Record Management

The record subcommands expect both the domain and record metadata. Each flag aligns with the struct tags in `cmd/dns/dns_record_*.go`.

### List records

```sh
virak-cli dns record list \
  --domain "example.com"
```

### Create records

```sh
# Root A record
virak-cli dns record create \
  --domain "example.com" \
  --record "@" \
  --type A \
  --content "203.0.113.10" \
  --ttl 3600

# AAAA
virak-cli dns record create \
  --domain "example.com" \
  --record "@" \
  --type AAAA \
  --content "2001:db8::10" \
  --ttl 3600

# CNAME
virak-cli dns record create \
  --domain "example.com" \
  --record "www" \
  --type CNAME \
  --content "example.com" \
  --ttl 3600

# MX
virak-cli dns record create \
  --domain "example.com" \
  --record "@" \
  --type MX \
  --content "mail.example.com" \
  --priority 10 \
  --ttl 3600

# TXT (SPF)
virak-cli dns record create \
  --domain "example.com" \
  --record "@" \
  --type TXT \
  --content "\"v=spf1 include:_spf.google.com ~all\"" \
  --ttl 3600

# SRV
virak-cli dns record create \
  --domain "example.com" \
  --record "_sip._tcp" \
  --type SRV \
  --content "sip.example.com" \
  --priority 10 \
  --weight 5 \
  --port 5060 \
  --ttl 3600

# TLSA
virak-cli dns record create \
  --domain "example.com" \
  --record "_443._tcp" \
  --type TLSA \
  --license 3 \
  --choicer 1 \
  --match 1 \
  --content "certificate-hash" \
  --ttl 3600

# CAA
virak-cli dns record create \
  --domain "example.com" \
  --record "@" \
  --type CAA \
  --flags 0 \
  --tag issue \
  --content "letsencrypt.org" \
  --ttl 3600
```

Key flags:

- `--record` is the label (`@`, `www`, `_sip._tcp`, etc.)
- `--type` must be one of the supported types enforced by the CLI
- `--content` is the payload (IP, hostname, text, or hash)
- `--priority`, `--weight`, `--port` apply to MX/SRV
- `--flags` and `--tag` apply to CAA
- `--license`, `--choicer`, `--match` apply to TLSA

### Update records

Record updates reference a specific content ULID. Retrieve it via `record list`, then:

```sh
virak-cli dns record update \
  --domain "example.com" \
  --record "@" \
  --type A \
  --contentId 01HF... \
  --content "203.0.113.20" \
  --ttl 300
```

Use this pattern to rotate MX targets, change TXT values, or refresh TLSA hashes.

### Delete records

```sh
virak-cli dns record delete \
  --domain "example.com" \
  --record "@" \
  --type A \
  --content-id 01HF...
```

The CLI validates the ULID before hitting the API, preventing accidental deletions.

## Events and Auditing

```sh
virak-cli dns events
```

Events capture domain additions, record mutations, and propagation actions. They are the CLI equivalent of the history tab in the Virak panel.

## Recipes

### Basic web + mail setup

```sh
virak-cli dns record create --domain example.com --record "@" --type A --content "203.0.113.10"
virak-cli dns record create --domain example.com --record "www" --type CNAME --content "example.com"
virak-cli dns record create --domain example.com --record "@" --type MX --content "mail.example.com" --priority 10
virak-cli dns record create --domain example.com --record "@" --type TXT --content "\"v=spf1 include:_spf.google.com ~all\""
```

### Delegating subdomains

```sh
virak-cli dns record create --domain example.com --record "api" --type A --content "203.0.113.50"
virak-cli dns record create --domain example.com --record "blog" --type CNAME --content "hosted.example.net"
```

### TLS hardening

```sh
virak-cli dns record create --domain example.com --record "@" --type CAA --flags 0 --tag issue --content "letsencrypt.org"
virak-cli dns record create --domain example.com --record "_443._tcp" --type TLSA --license 3 --choicer 1 --match 1 --content "<sha256>"
```

## Best Practices

- Follow the record-type reference from `/virak-cloud/docs` to choose the correct resource for each use case (A/AAAA for direct IPs, CNAME for aliases, SRV for SIP/VoIP, PTR for reverse mapping, TLSA for DANE).
- Use high TTLs (86400) for NS records, medium TTLs (3600) for stable services, and low TTLs (300) when planning migrations.
- Pair DNS changes with `virak-cli dns events` to verify that automation executed as expected.
- Store critical domain metadata (nameservers, registrar reference, contact emails) in the same repository alongside these CLI commands to keep a single source of truth.

## Troubleshooting

- `Domain is pending`: update your registrar to use the Virak Cloud nameservers and wait for propagation.
- `Record already exists`: list records to find the `contentId`, then update instead of creating duplicates.
- Email deliverability issues: ensure SPF, DKIM, and reverse PTR records match, leveraging the TXT and PTR examples above.
- TLS handshake failures: refresh TLSA hashes immediately after reissuing certificates.

For full REST API details, see the DNS section of the Virak Cloud public API documentation (`/api/external/domains` and `/api/external/domains/{id}/records`). The CLI commands described here are thin wrappers around those endpoints.