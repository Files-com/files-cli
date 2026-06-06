---
name: filescom-scheduled-exports
description: |
  A Scheduled Export defines a recurring schedule for generating one of the built-in CSV exports and e-mailing it to a Site Admin recipient.
---

# filescom-scheduled-exports

A Scheduled Export defines a recurring schedule for generating one of the built-in CSV exports and e-mailing it to a Site Admin recipient.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli scheduled-exports list`

List Scheduled Exports.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `export_type` or `disabled`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled` and `export_type`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `export_type`. |

### `files-cli scheduled-exports find`

Show Scheduled Export.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Scheduled Export ID. **Required.** |

### `files-cli scheduled-exports create`

Create Scheduled Export.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Name for this scheduled export. **Required.** |
| `--export-type` | string | Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit **Required.** |
| `--export-options` | object | Report-specific options. `permission_audit` supports `group_by` with `user` or `path`. |
| `--user-id` | int64 | Site Admin user who receives the completed export e-mail. |
| `--disabled` | bool | If true, this scheduled export will not run. |
| `--trigger` | enum | Schedule trigger type: `daily` or `custom_schedule`. One of: `daily`, `custom_schedule`. |
| `--interval` | string | If trigger is `daily`, this specifies how often to run the scheduled export. |
| `--recurring-day` | int64 | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `--schedule-times-of-day` | []string | Times of day in HH:MM format for schedule-driven exports. |
| `--schedule-time-zone` | string | Time zone used by the scheduled export. |
| `--holiday-region` | string | Optional holiday region used by schedule-driven exports. |

### `files-cli scheduled-exports update`

Update Scheduled Export.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Scheduled Export ID. **Required.** |
| `--name` | string | Name for this scheduled export. |
| `--export-type` | string | Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit |
| `--export-options` | object | Report-specific options. `permission_audit` supports `group_by` with `user` or `path`. |
| `--user-id` | int64 | Site Admin user who receives the completed export e-mail. |
| `--disabled` | bool | If true, this scheduled export will not run. |
| `--trigger` | enum | Schedule trigger type: `daily` or `custom_schedule`. One of: `daily`, `custom_schedule`. |
| `--interval` | string | If trigger is `daily`, this specifies how often to run the scheduled export. |
| `--recurring-day` | int64 | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `--schedule-times-of-day` | []string | Times of day in HH:MM format for schedule-driven exports. |
| `--schedule-time-zone` | string | Time zone used by the scheduled export. |
| `--holiday-region` | string | Optional holiday region used by schedule-driven exports. |

### `files-cli scheduled-exports delete`

Delete Scheduled Export.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Scheduled Export ID. **Required.** |

