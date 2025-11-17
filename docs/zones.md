# Managing Zones with Virak CLI

Zones represent the physical regions documented at [zone in virakcloud docs](https://docs.virakcloud.com) and surfaced through the `/api/external/zones` endpoints. 

## List and Select Zones

```sh
virak-cli zone list
```

The CLI prints each zone’s name and ULID. After listing, you can set a default zone directly from the interactive prompt or edit the config file:

```yaml
default:
  zoneId: "zone-xyz"
  zoneName: "Tehran-1"
```

Any command that accepts `--zoneId` will use this default when the flag is omitted.

## Inspect Zone Networks

```sh
virak-cli zone networks \
  --zoneId zone-xyz
```

This command fetches `GET /api/external/zones/{id}/networks`, revealing every L2 and L3 network offering in the region. Use it to validate capacity before provisioning new workloads or to discover shared networks available to your organization.

## Check Available Resources

```sh
virak-cli zone resources \
  --zoneId zone-xyz
```

Output includes collected vs. total memory, CPU cores, data volume, and VM limits for the current customer in that zone. These metrics mirror the capacity widgets in the Virak panel and help you preflight automation runs.

## List Active Services

```sh
virak-cli zone services \
  --zoneId zone-xyz
```

You will see booleans for each service family (Instances, DataVolume, Network, ObjectStorage, Kubernetes). Consult this before attempting to create clusters or specialized networks; if a feature is disabled, switch zones or contact support.

## How Zones Affect Other Commands

- Every compute, network, bucket, and cluster command loads `zoneId` from context. Override it with `--zoneId` when you need to work in multiple regions during the same session.
- CI jobs should export `VIRAK_CLI_ZONE_ID` to avoid editing config files.
- When following `/virak-cloud/docs` tutorials (e.g., Kubernetes or Kubernetes + DNS workflows), always verify that the documented services are enabled in your chosen zone via `zone services`.

## Best Practices

- Co-locate latency-sensitive workloads (instances, load balancers, Kubernetes) within the same zone to minimize cross-region bandwidth charges.
- When designing DR, use `zone list` to choose a secondary region that offers the same service set and instance types.
- Cache the output of `zone resources` during deployments to detect quota exhaustion early and trigger scale-up requests with Virak support if needed.

## API Reference

All zone commands call the public API endpoints listed in the Virak documentation:

- `GET /api/external/zones`
- `GET /api/external/zones/{id}/networks`
- `GET /api/external/zones/{id}/resources`
- `GET /api/external/zones/{id}/services`

Understanding these endpoints helps when you need to extend automation beyond the CLI.