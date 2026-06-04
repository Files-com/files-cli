---
name: filescom-bandwidth-snapshots
description: |
  A BandwidthSnapshot is information about site bandwidth usage for the given date.
---

# filescom-bandwidth-snapshots

A BandwidthSnapshot is information about site bandwidth usage for the given date.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli bandwidth-snapshots list`

List Bandwidth Snapshots.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `logged_at`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `logged_at`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `logged_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `logged_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `logged_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `logged_at`. |

