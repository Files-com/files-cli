---
name: filescom-holiday-calendars
description: |
  A Holiday Calendar defines site-wide holiday dates and optional partial-day windows that scheduled resources skip.
---

# filescom-holiday-calendars

A Holiday Calendar defines site-wide holiday dates and optional partial-day windows that scheduled resources skip.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli holiday-calendars list`

List Holiday Calendars.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |

### `files-cli holiday-calendars find`

Show Holiday Calendar.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Holiday Calendar ID. **Required.** |

### `files-cli holiday-calendars create`

Create Holiday Calendar.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Holiday Calendar name. **Required.** |

### `files-cli holiday-calendars update`

Update Holiday Calendar.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Holiday Calendar ID. **Required.** |
| `--name` | string | Holiday Calendar name. |

### `files-cli holiday-calendars delete`

Delete Holiday Calendar.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Holiday Calendar ID. **Required.** |

