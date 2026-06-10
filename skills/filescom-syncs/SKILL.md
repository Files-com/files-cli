---
name: filescom-syncs
description: |
  A Sync represents a file synchronization job between two locations (local-remote, remote-remote, local-child_site, etc).
---

# filescom-syncs

A Sync represents a file synchronization job between two locations (local-remote, remote-remote, local-child_site, etc). 
It can be scheduled, run manually, or triggered by custom logic. 
Syncs track their runs, status, and configuration.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli syncs list`

List Syncs.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id` or `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `disabled`, `src_remote_server_id` or `dest_remote_server_id`. Valid field combinations are `[ workspace_id, disabled ]`, `[ workspace_id, src_remote_server_id ]`, `[ workspace_id, dest_remote_server_id ]`, `[ disabled, src_remote_server_id ]`, `[ disabled, dest_remote_server_id ]`, `[ workspace_id, disabled, src_remote_server_id ]` or `[ workspace_id, disabled, dest_remote_server_id ]`. |

### `files-cli syncs find`

Show Sync.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sync ID. **Required.** |

### `files-cli syncs create`

Create Sync.

| Flag | Type | Description |
| --- | --- | --- |
| `--delete-empty-folders` | bool | Delete empty folders after sync? |
| `--description` | string | Description for this sync job |
| `--dest-path` | string | Absolute destination path for the sync |
| `--dest-remote-server-id` | int64 | Remote server ID for the destination (if remote) |
| `--disabled` | bool | Is this sync disabled? |
| `--exclude-patterns` | []string | Array of glob patterns to exclude |
| `--holiday-region` | string | Skip sync if there is a formal, observed holiday for this region. |
| `--include-patterns` | []string | Array of glob patterns to include |
| `--interval` | string | If trigger is `daily`, this specifies how often to run this sync.  One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `--keep-after-copy` | bool | Keep files after copying? |
| `--name` | string | Name for this sync job |
| `--recurring-day` | int64 | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc. |
| `--schedule-time-zone` | string | Time zone for scheduled times. If not set, times are interpreted as UTC. |
| `--schedule-times-of-day` | []string | Times of day to run in HH:MM format. For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date. |
| `--src-path` | string | Absolute source path for the sync |
| `--src-remote-server-id` | int64 | Remote server ID for the source (if remote) |
| `--sync-interval-minutes` | int64 | Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty. |
| `--trigger` | enum | Trigger type: daily, custom_schedule, or manual. One of: `daily`, `custom_schedule`, `manual`. |
| `--trigger-file` | string | Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync. |
| `--always-write-trigger-file` | bool | If true, the trigger file will be sent at the end of a successful sync even when no files were transferred. |
| `--workspace-id` | int64 | Workspace ID this sync belongs to |

### `files-cli syncs dry-run`

Dry Run Sync.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sync ID. **Required.** |

### `files-cli syncs manual-run`

Manually Run Sync.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sync ID. **Required.** |

### `files-cli syncs update`

Update Sync.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sync ID. **Required.** |
| `--delete-empty-folders` | bool | Delete empty folders after sync? |
| `--description` | string | Description for this sync job |
| `--dest-path` | string | Absolute destination path for the sync |
| `--dest-remote-server-id` | int64 | Remote server ID for the destination (if remote) |
| `--disabled` | bool | Is this sync disabled? |
| `--exclude-patterns` | []string | Array of glob patterns to exclude |
| `--holiday-region` | string | Skip sync if there is a formal, observed holiday for this region. |
| `--include-patterns` | []string | Array of glob patterns to include |
| `--interval` | string | If trigger is `daily`, this specifies how often to run this sync.  One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `--keep-after-copy` | bool | Keep files after copying? |
| `--name` | string | Name for this sync job |
| `--recurring-day` | int64 | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc. |
| `--schedule-time-zone` | string | Time zone for scheduled times. If not set, times are interpreted as UTC. |
| `--schedule-times-of-day` | []string | Times of day to run in HH:MM format. For `custom_schedule`, run at these times on specified days of week. For `daily`, run at these times on the scheduled interval date. |
| `--src-path` | string | Absolute source path for the sync |
| `--src-remote-server-id` | int64 | Remote server ID for the source (if remote) |
| `--sync-interval-minutes` | int64 | Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty. |
| `--trigger` | enum | Trigger type: daily, custom_schedule, or manual. One of: `daily`, `custom_schedule`, `manual`. |
| `--trigger-file` | string | Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync. |
| `--always-write-trigger-file` | bool | If true, the trigger file will be sent at the end of a successful sync even when no files were transferred. |

### `files-cli syncs delete`

Delete Sync.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sync ID. **Required.** |

