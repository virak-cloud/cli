# Managing Instances with Virak CLI

Virak CLI exposes the entire compute surface of Virak Cloud so you can script repeatable infrastructure changes instead of clicking through the panel.

## Prerequisites

- Virak Cloud account with an active API token
- `virak-cli login` completed (interactive OAuth or `--token`)
- Default zone set in `~/.virak-cli.yaml` or passed with `--zoneId`
- At least one network in the target zone (see the Networks guide)

Example configuration file:

```yaml
auth:
  token: "vrk-..."
default:
  zoneId: "zone-xyz"
  zoneName: "Tehran-1"
```

Environment variables such as `VIRAK_CLI_TOKEN` and `VIRAK_CLI_ZONE_ID` can override the file for CI pipelines.

## Discover Instance Resources

### Service offerings

```sh
virak-cli instance service-offering-list --available --zoneId zone-xyz
```

Use `--columns id,name,cpu,memory,price_up` to print only the metrics you need. Pair this with the official sizing guidance for workloads (web, database, batch) from the Virak Cloud docs.

### VM images

```sh
virak-cli instance vm-image-list --zoneId zone-xyz
```

Images expose OS name, version, category, and availability so you can match panel choices exactly.

### Networks and connectivity

```sh
virak-cli network list --zoneId zone-xyz
```

Collect the network IDs you plan to attach. Layer‑3 networks are required when instances must reach managed Kubernetes clusters or VPN gateways as described in the networking best practices.

## List and Inspect Instances

```sh
virak-cli instance list --zoneId zone-xyz
```

Show full metadata, including service offering, image, and attached volumes:

```sh
virak-cli instance show --instanceId inst-abc123 --zoneId zone-xyz
```

Interactive mode (`--interactive`) lets you pick from the current zone without remembering IDs.

## Creating Instances

### Non-interactive workflow

```sh
virak-cli instance create \
  --zoneId zone-xyz \
  --name "app-01" \
  --service-offering-id "so-highcpu" \
  --vm-image-id "img-ubuntu-22" \
  --network-ids '["net-public","net-private"]'
```

Flags must match the struct tags in `cmd/instance/instance_create.go`:

- `--name` human readable instance name
- `--service-offering-id` ULID from `instance service-offering-list`
- `--vm-image-id` ULID from `instance vm-image-list`
- `--network-ids` JSON array of network IDs
- `--zoneId` optional when `default.zoneId` is configured

### Interactive wizard

```sh
virak-cli instance create --interactive --zoneId zone-xyz
```

The CLI fetches offerings, filters to available plans, prompts for VM images, networks, and a name, and confirms before submitting the API call. This mirrors the panel workflow and uses the same validation checks outlined in the official compute guide.

### Best practices

- Reuse the firewall templates from the networking documentation to allow SSH (`22/tcp`) and application ports.
- Track service offerings and VM images via source control to keep reproducible environments aligned with change management.

## Lifecycle Operations

```sh
virak-cli instance start --instanceId inst-abc123 --zoneId zone-xyz
virak-cli instance stop --instanceId inst-abc123 --zoneId zone-xyz
virak-cli instance reboot --instanceId inst-abc123 --zoneId zone-xyz
virak-cli instance rebuild --instanceId inst-abc123 --zoneId zone-xyz
virak-cli instance delete --instanceId inst-abc123 --zoneId zone-xyz
```

Rebuild reinstalls the original VM image and matches the destructive behavior documented in the panel. Always snapshot or back up attached volumes before running it.

## Snapshots

```sh
virak-cli instance snapshot create \
  --instanceId inst-abc123 \
  --name "pre-release-2024-11" \
  --zoneId zone-xyz

virak-cli instance snapshot list --instanceId inst-abc123 --zoneId zone-xyz
virak-cli instance snapshot revert \
  --instanceId inst-abc123 \
  --snapshotId snap-01HF... \
  --zoneId zone-xyz
virak-cli instance snapshot delete \
  --instanceId inst-abc123 \
  --snapshotId snap-01HF... \
  --zoneId zone-xyz
```

Interactive mode is available for `list`, `create`, and `revert`. The CLI refuses to revert if a snapshot is already in `WAITING`, mirroring the protection described in the compute reference docs.

## Volume Management

List available offerings (useful for performance tiers and pricing):

```sh
virak-cli instance volume service-offering-list --zoneId zone-xyz
```

Create, attach, detach, and delete volumes:

```sh
virak-cli instance volume create \
  --zoneId zone-xyz \
  --serviceOfferingId vol-tier1 \
  --size 100 \
  --name "db-data"

virak-cli instance volume attach \
  --zoneId zone-xyz \
  --volumeId vol-01HF... \
  --instanceId inst-abc123

virak-cli instance volume detach \
  --zoneId zone-xyz \
  --volumeId vol-01HF... \
  --instanceId inst-abc123

virak-cli instance volume delete \
  --zoneId zone-xyz \
  --volumeId vol-01HF...
```

Interactive prompts help when you do not remember IDs. The attach command validates that the target instance has no pending snapshots, matching the safety rules outlined in the Virak Cloud storage guide.

## Monitoring and Console Access

```sh
virak-cli instance metrics \
  --instanceId inst-abc123 \
  --metrics memoryusedkbs,cpuused \
  --time 6 \
  --zoneId zone-xyz

virak-cli instance console --instanceId inst-abc123 --zoneId zone-xyz
```

Metrics stream the same data shown in the panel’s charts. The console command prints a signed URL so you can open the built-in VNC viewer when SSH is unavailable.

## Automation Patterns

- Pipe `virak-cli instance list --zoneId zone-xyz` through `jq` to locate stale instances, then call `instance delete`.
- Combine `instance service-offering-list --available` with CI environment variables to ensure every deployment uses currently in-stock hardware.
- Before scaling application tiers, script snapshot creation and verify success using `instance snapshot list` to align with the backup recommendations from the official docs.

## Troubleshooting

- `instance create` fails: confirm the service offering is available in the selected zone and that networks accept additional attachments.
- `instance start` never completes: check `virak-cli instance show --instanceId ...` and review recent service events in the zone (see the Zones guide).
- Console URL expired: rerun `instance console` to mint a fresh token.

For API-level details, refer to `POST /api/external/instances` and related endpoints in the Virak Cloud public API documentation. The CLI wraps these calls directly, so every flag maps to a documented request body.