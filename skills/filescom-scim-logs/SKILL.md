---
name: filescom-scim-logs
description: |
  A SCIM log entry represents a single SCIM request made to the system.
---

# filescom-scim-logs

A SCIM log entry represents a single SCIM request made to the system. It includes the request made and response provided to the SCIM client.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli scim-logs list`

List Scim Logs.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |

### `files-cli scim-logs find`

Show Scim Log.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Scim Log ID. **Required.** |

