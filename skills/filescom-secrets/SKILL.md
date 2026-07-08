---
name: filescom-secrets
description: |
  A Secret stores named, typed secret material for later use by features that reference the Secret by ID.
---

# filescom-secrets

A Secret stores named, typed secret material for later use by features that reference the Secret by ID.

Secret values are encrypted at rest and are write-only. API responses include metadata and configured value field names, but never include the stored secret values.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli secrets list`

List Secrets.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `secret_type`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `name` or `secret_type`. Valid field combinations are `[ workspace_id, name ]`, `[ workspace_id, secret_type ]`, `[ secret_type, name ]` or `[ workspace_id, secret_type, name ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |

### `files-cli secrets find`

Show Secret.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Secret ID. **Required.** |

### `files-cli secrets create`

Create Secret.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Secret name. **Required.** |
| `--description` | string | Internal description for your reference. |
| `--secret-type` | enum | Secret type. One of: `basic`, `token`, `headers`, `certificate`, `key_value`. **Required.** |
| `--metadata` | object | Non-secret metadata for the Secret type. |
| `--workspace-id` | int64 | Workspace ID. 0 means the default workspace. |

### `files-cli secrets update`

Update Secret.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Secret ID. **Required.** |
| `--name` | string | Secret name. |
| `--description` | string | Internal description for your reference. |
| `--secret-type` | enum | Secret type. One of: `basic`, `token`, `headers`, `certificate`, `key_value`. |
| `--metadata` | object | Non-secret metadata for the Secret type. |

### `files-cli secrets delete`

Delete Secret.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Secret ID. **Required.** |

