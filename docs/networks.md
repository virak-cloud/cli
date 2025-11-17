# Managing Networks with Virak CLI

Virak Cloud networking combines layer‑2 isolation, layer‑3 routed networks, HAProxy load balancing, VPN gateways, and granular firewall rules. 

## Prerequisites

- Authenticated Virak CLI session (`virak-cli login`)
- Default zone configured or supplied with `--zoneId`
- Appropriate service entitlements in the target zone (check with the Zones guide)

## Discover Network Offerings

```sh
virak-cli network service-offering \
  --type all \
  --zoneId zone-xyz
```

Use `--type l2` for flat VLAN networks or `--type l3` for routed networks. The CLI prints bandwidth limits, transfer plans, hourly rates, and localization data just like the Virak Cloud panel.

## Create and Manage Networks

### Layer‑2 networks

```sh
virak-cli network create l2 \
  --zoneId zone-xyz \
  --network-offering-id netoff-l2 \
  --name "db-backend"
```

The CLI validates that `network-offering-id` refers to an L2 offering before submitting the request.

### Layer‑3 networks

```sh
virak-cli network create l3 \
  --zoneId zone-xyz \
  --network-offering-id netoff-l3 \
  --name "k8s-mesh" \
  --gateway "10.25.0.1" \
  --netmask "255.255.255.0"
```

L3 networks power features such as Kubernetes clusters, VPN, and port forwarding. Choose CIDR blocks that do not overlap with on-prem networks to simplify hybrid routing.

### Inventory and connections

```sh
virak-cli network list --zoneId zone-xyz
virak-cli network show --zoneId zone-xyz --networkId net-abc123
virak-cli network instance list --zoneId zone-xyz --networkId net-abc123

virak-cli network instance connect \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --instanceId inst-01HF...

virak-cli network instance disconnect \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --instanceId inst-01HF...
```

When a network is no longer needed, delete it:

```sh
virak-cli network delete --zoneId zone-xyz --networkId net-abc123
```

## Firewall Rules

The CLI exposes separate subcommands for IPv4 and IPv6. They implement the same schema the official docs describe: traffic type, protocol, source/destination, and optional port or ICMP parameters.

```sh
virak-cli network firewall ipv4 create \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --trafficType Ingress \
  --protocolType TCP \
  --ipSource 0.0.0.0/0 \
  --ipDestination 10.25.0.0/24 \
  --portStart 80 \
  --portEnd 80

virak-cli network firewall ipv6 create \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --trafficType Ingress \
  --protocolType TCP \
  --ipSource ::/0 \
  --ipDestination 2001:db8::/64 \
  --portStart 443 \
  --portEnd 443
```

List or delete rules with `network firewall ipv4|ipv6 list` and `network firewall ipv4|ipv6 delete`. Follow the recommendations from `/virak-cloud/docs`:

- Allow HTTP/HTTPS from `0.0.0.0/0` or your CDN ranges
- Limit SSH (port 22) to trusted IPs
- Add an explicit deny rule (priority 1000) as a safety net

## Load Balancers

```sh
virak-cli network lb list --zoneId zone-xyz

virak-cli network lb create \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --publicIpId pip-01HF... \
  --name "frontend" \
  --algorithm roundrobin \
  --publicPort 80 \
  --privatePort 8080
```

Assign instances by passing instance network IDs:

```sh
virak-cli network lb assign \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --ruleId lbr-01HF... \
  --instanceNetworkIds instnet-1,instnet-2
```

Monitor HAProxy from the CLI:

```sh
virak-cli network lb haproxy live \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --ruleId lbr-01HF...

virak-cli network lb haproxy log \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --ruleId lbr-01HF...
```

## Public IPs and NAT

```sh
virak-cli network public-ip list \
  --zoneId zone-xyz \
  --networkId net-abc123

virak-cli network public-ip associate \
  --zoneId zone-xyz \
  --networkId net-abc123

virak-cli network public-ip disassociate \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --publicIpId pip-01HF...
```

Static NAT ties a public IP to an instance:

```sh
virak-cli network public-ip staticnat enable \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --networkPublicIpId pip-01HF... \
  --instanceId inst-01HF...

virak-cli network public-ip staticnat disable \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --networkPublicIpId pip-01HF...
```

## Port Forwarding

```sh
virak-cli network port-forward create \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --protocol TCP \
  --publicPort 2222 \
  --privatePort 22 \
  --privateIp 10.25.0.10

virak-cli network port-forward list \
  --zoneId zone-xyz \
  --networkId net-abc123

virak-cli network port-forward delete \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --ruleId pf-01HF...
```

## VPN

```sh
virak-cli network vpn show --zoneId zone-xyz --networkId net-abc123
virak-cli network vpn enable --zoneId zone-xyz --networkId net-abc123
virak-cli network vpn disable --zoneId zone-xyz --networkId net-abc123
virak-cli network vpn update \
  --zoneId zone-xyz \
  --networkId net-abc123 \
  --psk "new-shared-secret"
```

The VPN configuration you download from the panel can be distributed to remote clients, matching the workflow described in the official VPN guide.

## Best Practices

- Tag each network by purpose in the `--name` field (`app-public`, `db-private`, `shared-services`) and keep a table in version control.
- Combine `network service-offering` output with the firewall template from `/virak-cloud/docs` to ensure compliance before provisioning workloads.
- Automate drift detection: periodically run `virak-cli network firewall ipv4 list` and compare against your golden rule set.
- When exposing public endpoints, pair load balancers with DNS records (see the DNS guide) and TLS certificates managed through Virak Cloud or your preferred ACM.

## Troubleshooting

- Network creation fails: confirm the offering ID belongs to the selected zone and that quota allows more networks.
- Firewall rules not applied: check for overlapping rules with higher priority; the CLI prints rule IDs so you can delete or update the offending entry.
- Load balancer traffic not flowing: verify backend instances are connected to the same network and that health checks succeed via `network lb haproxy live`.
- Public IP association stuck: ensure the network supports internet connectivity and that you have available quotas for elastic IPs.

For API-level insight, reference `POST /api/external/networks`, `POST /api/external/networks/{id}/firewall/rules`, and related endpoints in the Virak Cloud public API docs. The Virak CLI commands above are thin wrappers around those endpoints.