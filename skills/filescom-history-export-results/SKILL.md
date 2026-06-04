---
name: filescom-history-export-results
description: |
  A HistoryExportResult is a single record with historical information about actions that happened in your site.
---

# filescom-history-export-results

A HistoryExportResult is a single record with historical information about actions that happened in your site. Results can be obtained as a part of History Export request.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli history-export-results list`

List History Export Results.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--history-export-id` | int64 | ID of the associated history export. **Required.** |

