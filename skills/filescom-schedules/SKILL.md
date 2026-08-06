---
name: filescom-schedules
description: |
  A Schedule is a named, reusable weekday-and-time schedule shared by scheduled resources across a Site.
---

# filescom-schedules

A Schedule is a named, reusable weekday-and-time schedule shared by scheduled resources across a Site.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli schedules list`

List Schedules.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`. |

### `files-cli schedules find`

Show Schedule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Schedule ID. **Required.** |

### `files-cli schedules create`

Create Schedule.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Schedule name. **Required.** |
| `--schedule-days-of-week` | []int64 | 0-based weekdays used by the Schedule. 0 is Sunday. **Required.** |
| `--schedule-times-of-day` | []string | Times of day in HH:MM format (24-hour). **Required.** |
| `--schedule-time-zone` | string | Time zone for scheduled times. If not set, times are interpreted as UTC. |
| `--holiday-region` | string | Optional holiday region on which linked resources do not run. |

### `files-cli schedules update`

Update Schedule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Schedule ID. **Required.** |
| `--name` | string | Schedule name. |
| `--schedule-days-of-week` | []int64 | 0-based weekdays used by the Schedule. 0 is Sunday. |
| `--schedule-times-of-day` | []string | Times of day in HH:MM format (24-hour). |
| `--schedule-time-zone` | string | Time zone for scheduled times. If not set, times are interpreted as UTC. |
| `--holiday-region` | string | Optional holiday region on which linked resources do not run. |

### `files-cli schedules delete`

Delete Schedule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Schedule ID. **Required.** |

