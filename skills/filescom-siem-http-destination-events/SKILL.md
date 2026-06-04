---
name: filescom-siem-http-destination-events
description: |
  A SiemHttpDestinationEvent is a log record for SIEM publishing failures and recoveries.
---

# filescom-siem-http-destination-events

A SiemHttpDestinationEvent is a log record for SIEM publishing failures and recoveries.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli siem-http-destination-events list`

List SIEM HTTP Destination Events.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status` or `siem_http_destination_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `siem_http_destination_id` or `status`. Valid field combinations are `[ siem_http_destination_id, created_at ]`, `[ status, created_at ]`, `[ siem_http_destination_id, status ]` or `[ siem_http_destination_id, status, created_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |

### `files-cli siem-http-destination-events find`

Show SIEM HTTP Destination Event.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Siem Http Destination Event ID. **Required.** |

