---
name: filescom-action-notification-export-results
description: |
  An ActionNotificationExportResult is a single record containing webhook log information that can be obtained as a part of Action Notification Export request.
---

# filescom-action-notification-export-results

An ActionNotificationExportResult is a single record containing webhook log information that can be obtained as a part of Action Notification Export request.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli action-notification-export-results list`

List Action Notification Export Results.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--action-notification-export-id` | int64 | ID of the associated action notification export. **Required.** |

