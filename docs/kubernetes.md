# Managing Kubernetes Clusters with Virak CLI

Virak Cloud’s managed Kubernetes service gives you HA masters, autoscaling workers, and native integration with networks, load balancers, and object storage. The Virak CLI mirrors every workflow documented in  [kubernetes guide in virakcloud docs](https://docs.virakcloud.com) while exposing additional automation hooks.

## Prerequisites

- Authenticated Virak CLI session
- SSH key uploaded via `virak-cli user ssh-key create`
- Layer‑3 network already provisioned in the target zone
- Default zone configured or passed with `--zoneId`

## Discover Versions and Offerings

```sh
virak-cli cluster versions-list --zoneId zone-xyz
virak-cli cluster offering --zoneId zone-xyz
virak-cli user ssh-key list
virak-cli network list --zoneId zone-xyz
```

Use these commands to gather the IDs used by create/scale/update operations. The output matches the “Kubernetes Versions” and “Node Profiles” tables from the official docs.

## Create Clusters

```sh
virak-cli cluster create \
  --zoneId zone-xyz \
  --name "prod-gke" \
  --versionId ver-1-29 \
  --offeringId so-highcpu \
  --sshKeyId ssh-01HF... \
  --networkId net-l3-prod \
  --size 5 \
  --ha \
  --description "Production workload" \
  --privateRegistryUsername "ghcr-user" \
  --privateRegistryPassword "token" \
  --privateRegistryUrl "ghcr.io/org"
```

Flag summary (see `cmd/cluster/kubernetes_cluster_create.go`):

- `--versionId` Kubernetes release ULID (from `cluster versions-list`)
- `--offeringId` node flavor ULID (from `cluster offering`)
- `--sshKeyId` Virak Cloud SSH key ID
- `--networkId` L3 network ID
- `--size` total nodes (masters + workers). When `--ha` is set, at least three controllers are provisioned automatically.
- Optional container registry credentials allow private images to pull successfully.

Autoscaling is controlled after creation via `cluster scale` (see below).

## Inventory and Inspection

```sh
virak-cli cluster list --zoneId zone-xyz
virak-cli cluster show --zoneId zone-xyz --clusterId clu-01HF...
```

`cluster show` returns node details, load balancer addresses, and kubeconfig download hints, mirroring what the panel displays.

## Lifecycle Operations

```sh
virak-cli cluster start --zoneId zone-xyz --clusterId clu-01HF...
virak-cli cluster stop --zoneId zone-xyz --clusterId clu-01HF...
virak-cli cluster delete --zoneId zone-xyz --clusterId clu-01HF...

virak-cli cluster update \
  --zoneId zone-xyz \
  --clusterId clu-01HF... \
  --name "prod-gke-v2" \
  --description "Renamed to match project"

virak-cli cluster scale \
  --zoneId zone-xyz \
  --clusterId clu-01HF... \
  --auto-scaling \
  --min-cluster-size 3 \
  --max-cluster-size 10
```

When `--auto-scaling` is disabled, pass `--cluster-size` to pin the worker count. The CLI enforces logical constraints (min ≤ max, positive integers).

Service events help diagnose provisioning issues:

```sh
virak-cli cluster service-events --zoneId zone-xyz
```

## Accessing the Cluster

After creation, download the kubeconfig from the panel or via the API. Export it for `kubectl`:

```sh
export KUBECONFIG=/path/to/virak-prod.yaml
kubectl get nodes
kubectl get pods --all-namespaces
```

These steps mirror the official “Connect with kubectl” section. Keep credentials secure; rotate them according to your org’s policies.

## Networking and Security

- Clusters must live inside a dedicated L3 network. Use the Networks guide to create it and attach firewall rules (allow 6443/tcp from your management IPs, allow internal pod/service CIDRs).
- Public services should sit behind Virak load balancers. Assign worker instances to a load balancer (`virak-cli network lb assign ...`) and terminate TLS using the same patterns referenced in `/virak-cloud/docs`.
- Object storage and container registries reside in the same private network, so configure security groups/firewalls accordingly.

## Integration Examples

### CI/CD bootstrap

```sh
virak-cli cluster create ... --name "ci-${GIT_SHA}" --size 3
virak-cli cluster show --clusterId "$CLUSTER_ID" --zoneId zone-xyz > cluster.json
jq -r '.data.kubeconfig' cluster.json > kubeconfig
kubectl --kubeconfig kubeconfig apply -f manifests/
```

### Scheduled upgrades

1. List available versions: `virak-cli cluster versions-list`
2. Validate release notes in `/virak-cloud/docs`
3. Recreate with the new `--versionId` or use the panel for in-place upgrades when supported.

## Troubleshooting

- `cluster create` fails: verify the L3 network is not already reserved by another cluster and that the SSH key exists.
- Nodes stuck in `REGISTERING`: check `cluster service-events` and firewall rules; ensure outbound internet access if you pull public container images.
- `kubectl` cannot connect: confirm the kubeconfig endpoint is reachable (port 6443) and that your VPN or bastion host allows it.
- Autoscaling idle: inspect workloads to ensure pods request sufficient CPU/memory; the Virak autoscaler follows the same logic as upstream Cluster Autoscaler.

For API details, consult `POST /api/external/kubernetes/clusters`, `POST /api/external/kubernetes/clusters/{id}/scale`, and related endpoints in the Virak Cloud documentation. The CLI commands shown above are direct wrappers, so flags map to request payloads one-to-one.