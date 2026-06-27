---
name: filescom-ai-tasks
description: |
  An AI Task defines a Files.com AI prompt that can run on a schedule or in response to file actions.
---

# filescom-ai-tasks

An AI Task defines a Files.com AI prompt that can run on a schedule or in response to file actions.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli ai-tasks list`

List Ai Tasks.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `id`, `disabled` or `updated_at`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled`, `trigger` or `workspace_id`. Valid field combinations are `[ workspace_id, disabled ]`. |

### `files-cli ai-tasks find`

Show Ai Task.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Ai Task ID. **Required.** |

### `files-cli ai-tasks create [path]`

Create Ai Task.

| Flag | Type | Description |
| --- | --- | --- |
| `--description` | string | AI Task description. |
| `--disabled` | bool | If true, this AI Task will not run. |
| `--holiday-region` | string | Optional holiday region used by scheduled AI Tasks. |
| `--interval` | string | If trigger is `daily`, this specifies how often to run the AI Task. |
| `--name` | string | AI Task name. **Required.** |
| `--path` | string | Path scope used for action-triggered AI Tasks. |
| `--prompt` | string | Prompt sent when this AI Task is invoked. **Required.** |
| `--recurring-day` | int64 | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `--schedule-time-zone` | string | Time zone used by the AI Task schedule. |
| `--schedule-times-of-day` | []string | Times of day in HH:MM format for scheduled AI Tasks. |
| `--source` | string | Source glob used with `path` for action-triggered AI Tasks. |
| `--trigger` | enum | How this AI Task is triggered. One of: `manual`, `daily`, `custom_schedule`, `action`. |
| `--trigger-actions` | []string | If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli ai-tasks manual-run`

Manually Run AI Task.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Ai Task ID. **Required.** |

### `files-cli ai-tasks update [path]`

Update Ai Task.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Ai Task ID. **Required.** |
| `--description` | string | AI Task description. |
| `--disabled` | bool | If true, this AI Task will not run. |
| `--holiday-region` | string | Optional holiday region used by scheduled AI Tasks. |
| `--interval` | string | If trigger is `daily`, this specifies how often to run the AI Task. |
| `--name` | string | AI Task name. |
| `--path` | string | Path scope used for action-triggered AI Tasks. |
| `--prompt` | string | Prompt sent when this AI Task is invoked. |
| `--recurring-day` | int64 | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `--schedule-time-zone` | string | Time zone used by the AI Task schedule. |
| `--schedule-times-of-day` | []string | Times of day in HH:MM format for scheduled AI Tasks. |
| `--source` | string | Source glob used with `path` for action-triggered AI Tasks. |
| `--trigger` | enum | How this AI Task is triggered. One of: `manual`, `daily`, `custom_schedule`, `action`. |
| `--trigger-actions` | []string | If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli ai-tasks delete`

Delete Ai Task.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Ai Task ID. **Required.** |

