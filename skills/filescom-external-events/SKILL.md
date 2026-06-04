---
name: filescom-external-events
description: |
  An ExternalEvent is a log that is sent to the cloud from a client application such as the Files.com CLI.
---

# filescom-external-events

An ExternalEvent is a log that is sent to the cloud from a client application such as the Files.com CLI.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli external-events list`

List External Events.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status` or `event_type`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `status`. Valid field combinations are `[ status, created_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |

### `files-cli external-events find`

Show External Event.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | External Event ID. **Required.** |

### `files-cli external-events create`

Create External Event.

| Flag | Type | Description |
| --- | --- | --- |
| `--status` | enum | Status of event. One of: `success`, `failure`, `partial_failure`, `in_progress`, `skipped`. **Required.** |
| `--body` | string | Event body **Required.** |

