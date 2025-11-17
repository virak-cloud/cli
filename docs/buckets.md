# Managing Object Storage (Buckets) with Virak CLI

Virak Cloud Object Storage is S3 compatible. Everything you can do from the panel or AWS CLI is available through `virak-cli` plus the standard AWS SDK ecosystem. This guide aligns with the [object-storage guide in virakcloud docs](https://docs.virakcloud.com/en/guides/storage/object-storage) object storage tutorial.

## Prerequisites

- Authenticated Virak CLI session (`virak-cli login`)
- Default zone configured or provided with `--zoneId`
- AWS CLI installed for object-level operations

Configure AWS CLI once per workstation:

```sh
aws configure
# Access Key: vrk-...
# Secret Key: copy from the panel
# Region: us-east-1
# Output: json
```

Virak Cloud requires the custom endpoint:

```sh
export AWS_ENDPOINT_URL=https://s3.virakcloud.com
# or pass --endpoint-url https://s3.virakcloud.com on every command
```

## Bucket Lifecycle with Virak CLI

### List buckets

```sh
virak-cli bucket list --zoneId zone-xyz
```

The CLI prints bucket IDs, names, policy, and region so you can feed the IDs into the rest of the commands. To inspect a single bucket:

```sh
virak-cli bucket show \
  --bucketId buk-01HF... \
  --zoneId zone-xyz
```

### Create buckets

```sh
virak-cli bucket create \
  --zoneId zone-xyz \
  --name "media-prod" \
  --policy Public
```

Policies accept `Private` or `Public`, matching the same options exposed in the panel. Creation is asynchronous; follow up with `bucket list` or `bucket show`.

### Update policies

```sh
virak-cli bucket update \
  --zoneId zone-xyz \
  --bucketId buk-01HF... \
  --policy Private
```

Use this to flip between public and private access before layering additional AWS ACL or policy rules.

### Delete buckets

```sh
virak-cli bucket delete \
  --zoneId zone-xyz \
  --bucketId buk-01HF...
```

Virak Cloud enforces the same empty-bucket requirement as AWS S3. Remove objects first (see the next section).

### Events and auditing

```sh
virak-cli bucket events \
  --zoneId zone-xyz \
  --bucketId buk-01HF...

virak-cli bucket events --zoneId zone-xyz
```

Without `--bucketId` the CLI returns all bucket events for the zone. These entries mirror the activity stream documented in the Virak panel.

## Object Operations with AWS CLI

Because Virak Cloud exposes a strict S3 API, reuse the official workflow:

```sh
# Upload
aws s3 cp ./logo.webp s3://media-prod/assets/logo.webp --endpoint-url https://s3.virakcloud.com

# Download
aws s3 cp s3://media-prod/assets/logo.webp ./logo.webp --endpoint-url https://s3.virakcloud.com

# List objects
aws s3 ls s3://media-prod/assets/ --endpoint-url https://s3.virakcloud.com

# Sync directories
aws s3 sync ./dist s3://media-prod/web --endpoint-url https://s3.virakcloud.com
```

The same applies to advanced features:

- Multipart uploads (`aws s3 cp large.iso ...`)
- Presigned URLs (`aws s3 presign s3://media-prod/private/file.zip --expires-in 900 --endpoint-url ...`)
- Versioning (`aws s3api put-bucket-versioning --bucket media-prod --versioning-configuration Status=Enabled`)
- Static website hosting (`aws s3 website s3://media-prod/ --index-document index.html --endpoint-url ...`)

All commands and examples come directly from the Virak Cloud documentation snippet for AWS CLI configuration.

## Typical Workflows

### Backups

```sh
aws s3 sync /var/backups s3://company-backups/daily --endpoint-url https://s3.virakcloud.com
```

Pair this with `virak-cli bucket events` to confirm lifecycle policies executed as expected.

### Static website hosting

1. `virak-cli bucket create --policy Public`
2. `aws s3 sync ./build s3://my-site/ --endpoint-url https://s3.virakcloud.com`
3. `aws s3 website s3://my-site/ --index-document index.html --error-document 404.html --endpoint-url https://s3.virakcloud.com`
4. Use the DNS guide to point `www.example.com` at the bucket endpoint.

### Log archiving

```sh
aws s3 cp /var/log/app.log s3://company-logs/$(date +%Y-%m-%d)/app.log --endpoint-url https://s3.virakcloud.com
```

Virak Cloud exposes dedicated metrics for upload/download traffic (see Finance guide) so you can keep cost visibility.

## Best Practices

- Keep `virak-cli bucket show` outputs in infrastructure runbooks so teams always know the canonical IDs.
- Use `Private` policy and layer ACLs or bucket policies only when you need public access.

- For automation, export `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_ENDPOINT_URL` inside CI jobs instead of calling `aws configure`.

## Troubleshooting

- `AccessDenied`: confirm your API token still has object storage permissions and that the AWS keys were copied correctly.
- `NoSuchBucket`: verify the bucket was created in the same zone whose endpoint you are targeting.
- Slow transfers: enable multipart uploads by default for files larger than 100 MB, as recommended in the official guide.
- `bucket delete` stuck: run `aws s3 rm s3://bucket --recursive --endpoint-url ...` to empty it, then retry.

For REST API specifics, see `GET/PUT/DELETE /api/external/buckets` in the public Virak Cloud API reference. The CLI maps its flags directly to those endpoints while AWS CLI handles object-level operations over the S3-compatible API.