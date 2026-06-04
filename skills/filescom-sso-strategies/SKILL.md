---
name: filescom-sso-strategies
description: |
  An SSOStrategy is a way for users to sign in via another identity provider, such as Okta or Auth0.
---

# filescom-sso-strategies

An SSOStrategy is a way for users to sign in via another identity provider, such as Okta or Auth0.

It is rare that you will need to use API endpoints for managing these, and we recommend instead managing these via the web interface.
Nevertheless, we share the API documentation here.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli sso-strategies list`

List SSO Strategies.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |

### `files-cli sso-strategies find`

Show SSO Strategy.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sso Strategy ID. **Required.** |

### `files-cli sso-strategies sync`

Synchronize provisioning data with the SSO remote server.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sso Strategy ID. **Required.** |

