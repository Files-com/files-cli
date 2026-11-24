---
name: filescom-automation-versions
description: |
  An AutomationVersion is an immutable revision of an Automation v2 graph.
---

# filescom-automation-versions

An AutomationVersion is an immutable revision of an Automation v2 graph. Each saved graph change creates a revision. Restoring a revision saves its definition as a new current revision and keeps existing runs and Run Automation nodes pinned to their selected revisions.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli automation-versions list`

List Automation Versions.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `--automation-id` | int64 | ID of the Automation that owns this revision. **Required.** |

### `files-cli automation-versions find`

Show Automation Version.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Automation Version ID. **Required.** |

### `files-cli automation-versions restore`

Restore an Automation v2 definition as a new current version.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Automation Version ID. **Required.** |

