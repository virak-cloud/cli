# Managing Finance and Billing with Virak CLI

Virak Cloud exposes wallet balances, invoices, payments, and granular expenses through the CLI so you can reconcile usage without leaving your terminal.

## Prerequisites

- Authenticated Virak CLI session with billing access
- Organization ID with active services in the selected zone

## Wallet Snapshot

```sh
virak-cli finance wallet
```

Returns balance, credit limit, and billing contact information. Use this in monitoring jobs to alert when credit drops below your threshold.

## Cost Documents

```sh
virak-cli finance documents --year 2024
```

The `--year` flag is required. Results include invoices, pro-forma statements, and any other documents generated during that fiscal year. Pair the CLI output with your accounting system’s imports to automate reconciliations.

## Payments

```sh
virak-cli finance payments
```

This lists successful and pending transactions. For detailed inspection, use the panel to download the associated receipt.

## Expenses with Filters

`finance expenses` mirrors the expenses report described in the Virak docs. It requires both a product type and product ID to avoid ambiguous queries.

Valid product types (from `cmd/finance/finance_expenses.go`):

- Instance
- InstanceNetworkSecondaryIpAddressV4
- BucketSize
- NetworkInternetPublicAddressV4
- KubernetesNode
- NetworkDevice
- NetworkTraffic
- SupportOfferings
- BucketUploadTraffic
- BucketDownloadTraffic

Example: investigate a specific instance’s hourly costs.

```sh
virak-cli finance expenses \
  --product-type Instance \
  --product-id inst-01HF... \
  --start-date 2025-01-01 \
  --end-date 2025-01-31 \
  --type HourlyUsage \
  --page 1
```

Additional flags:

- `--start-date` / `--end-date` in `YYYY-MM-DD`
- `--type` filters by expense category returned by the API
- `--page` paginates long histories

The CLI validates product types before calling the API, so typos are caught locally.

## Budgeting Tips

- Run `finance wallet` on a schedule and send the JSON payload to your monitoring stack.
- Use `finance documents --year $(date +%Y)` at month end to ensure invoices were issued for every project.
- Track expensive instances by combining `finance expenses` with metadata from `virak-cli instance show`.
- Follow the storage optimization tips in `/virak-cloud/docs` (lifecycle rules, compression) to cut down on `BucketUploadTraffic` and `BucketDownloadTraffic` charges.

## Troubleshooting

- **Missing invoices**: rerun `finance documents --year ...` and confirm the command used the correct year. Invoices generate per calendar year.
- **Unexpected spikes**: query expenses per product ID to identify runaway instances, NAT traffic, or support plans.
- **Permission errors**: ensure the token used for the CLI has access to the billing APIs; contact an account administrator if needed.

For deeper integrations, reference:

- `GET /api/external/finance/wallet`
- `GET /api/external/finance/documents?year=...`
- `GET /api/external/finance/payments`
- `GET /api/external/finance/expenses`

The Virak CLI calls these endpoints directly, so you can replicate requests in custom tooling when required.