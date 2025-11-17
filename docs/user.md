# Managing User Accounts with Virak CLI

Use these commands to authenticate the CLI, inspect your profile, manage SSH keys, and validate API tokens.

## Authenticate Once

```sh
virak-cli login
```

- Without flags, the CLI opens the Virak OAuth page in your browser and stores the returned token in `~/.virak-cli.yaml`.
- Provide `--token` to paste an API token directly (useful for headless servers).
- The config file also stores your default zone so you do not need to pass `--zoneId` for every command.

## View Profile Details

```sh
virak-cli user profile
```

This fetches `/api/external/user/profile` and prints name, email, organization, and feature entitlements. Keep this output handy for audits or support tickets.

## SSH Key Management

List keys:

```sh
virak-cli user ssh-key list
```

Add a new key (generate one if needed):

```sh
ssh-keygen -t ed25519 -C "ops@example.com"

virak-cli user ssh-key create \
  --name "ops-laptop" \
  --public-key "$(cat ~/.ssh/id_ed25519.pub)"
```

Delete a key:

```sh
virak-cli user ssh-key delete --ssh-key-id ssh-01HF...
```

The CLI validates the ULID before hitting the API, reducing accidental removals. Removing unused keys follows the security guidance in the Virak documentation.

## Token Visibility

```sh
virak-cli user token abilities
virak-cli user token validate
```

- `token abilities` lists every scope attached to the token currently stored in the config or `VIRAK_CLI_TOKEN`.
- `token validate` confirms whether the token is still active. Run this in CI pipelines before starting long operations to fail fast when access is revoked.

## Best Practices

- Store tokens as secrets in your CI/CD platform rather than embedding them in source control.
- Rotate SSH keys regularly and remove keys from the CLI as soon as devices are decommissioned.
- Combine `user token abilities` with infrastructure as code reviews to ensure automation uses the minimum required scope set, as recommended in `/virak-cloud/docs`.
- When collaborating across teams, use separate tokens per automation context so revoking one does not break unrelated workflows.

## Troubleshooting

- Browser does not open during `login`: pass `--token` with a token generated from the Virak panel.
- `user token validate` fails: refresh the token via the panel or ask an organization admin for new credentials.
- SSH access denied: confirm the key appears in `user ssh-key list`, then redeploy the instance or ensure the key was attached during instance creation.

User commands map directly to `GET/POST /api/external/user/*` endpoints. If you need custom integrations, inspect network traffic with the `--debug` flag or call those endpoints directly using the token you validated above.
