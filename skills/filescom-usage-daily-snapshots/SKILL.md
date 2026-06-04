---
name: filescom-usage-daily-snapshots
description: |
  A UsageDailySnapshot is a detailed site-wide usage report for the day.
---

# filescom-usage-daily-snapshots

A UsageDailySnapshot is a detailed site-wide usage report for the day. Includes information about users, API, and storage.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli usage-daily-snapshots list`

List Usage Daily Snapshots.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `date`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `date` and `usage_snapshot_id`. Valid field combinations are `[ usage_snapshot_id, date ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `date`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `date`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `date`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `date`. |

