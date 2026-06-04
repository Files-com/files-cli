---
name: filescom-pending-work-events
description: |
  A PendingWorkEvent is a log record for pending file work failures, such as GPG Encryption and GPG Decryption.
---

# filescom-pending-work-events

A PendingWorkEvent is a log record for pending file work failures, such as GPG Encryption and GPG Decryption.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli pending-work-events list`

List Pending Work Events.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status` or `folder_behavior_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `folder_behavior_id` or `status`. Valid field combinations are `[ folder_behavior_id, created_at ]`, `[ status, created_at ]`, `[ folder_behavior_id, status ]` or `[ folder_behavior_id, status, created_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |

### `files-cli pending-work-events find`

Show Pending Work Event.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Pending Work Event ID. **Required.** |

