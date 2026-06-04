---
name: filescom-bundle-downloads
description: |
  A BundleDownload is a record of the download action that happened in the bundle.
---

# filescom-bundle-downloads

A BundleDownload is a record of the download action that happened in the bundle.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli bundle-downloads list`

List Share Link Downloads.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `created_at`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `--bundle-id` | int64 | Bundle ID |
| `--bundle-registration-id` | int64 | BundleRegistration ID |

