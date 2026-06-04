---
name: filescom-user-security-events
description: |
  A UserSecurityEvent is a log record for user security activity such as user lockouts.
---

# filescom-user-security-events

A UserSecurityEvent is a log record for user security activity such as user lockouts.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli user-security-events list`

List User Security Events.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at` and `user_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `user_id`. Valid field combinations are `[ user_id, created_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |

### `files-cli user-security-events find`

Show User Security Event.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Security Event ID. **Required.** |

