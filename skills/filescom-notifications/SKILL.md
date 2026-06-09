---
name: filescom-notifications
description: |
  A Notification is our feature that sends E-Mails when specific actions occur in the folder.
---

# filescom-notifications

A Notification is our feature that sends E-Mails when specific actions occur in the folder.

Emails are sent in batches, with email frequency options of every 5 minutes, every 15 minutes, hourly, or daily. They will include a list of the matching actions within the configured notification period, limited to the first 100.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli notifications list [path]`

List Notifications.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `path`, `user_id` or `group_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `user_id`, `workspace_id` or `group_id`. Valid field combinations are `[ workspace_id, path ]`, `[ workspace_id, user_id ]`, `[ workspace_id, group_id ]` or `[ workspace_id, user_id, path ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. |
| `--path` | string | Show notifications for this Path. |
| `--include-ancestors` | bool | If `include_ancestors` is `true` and `path` is specified, include notifications for any parent paths. Ignored if `path` is not specified. |
| `--group-id` | string | (no description) |

### `files-cli notifications find`

Show Notification.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Notification ID. **Required.** |

### `files-cli notifications create [path]`

Create Notification.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | The id of the user to notify. Provide `user_id`, `username` or `group_id`. |
| `--notify-on-copy` | bool | If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads. |
| `--notify-on-delete` | bool | Trigger on files deleted in this path? |
| `--notify-on-download` | bool | Trigger on files downloaded in this path? |
| `--notify-on-move` | bool | Trigger on files moved to this path? |
| `--notify-on-upload` | bool | Trigger on files created/uploaded/updated/changed in this path? |
| `--notify-user-actions` | bool | If `true` actions initiated by the user will still result in a notification |
| `--recursive` | bool | If `true`, enable notifications for each subfolder in this path |
| `--send-interval` | string | The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`. |
| `--subject` | string | Custom subject line to use for notification emails |
| `--message` | string | Custom message to include in notification emails |
| `--triggering-filenames` | []string | Array of filenames (possibly with wildcards) to scope trigger |
| `--triggering-group-ids` | []int64 | If set, will only notify on actions made by a member of one of the specified groups |
| `--triggering-user-ids` | []int64 | If set, will only notify on actions made one of the specified users |
| `--trigger-by-share-recipients` | bool | Notify when actions are performed by a share recipient? |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |
| `--group-id` | int64 | The ID of the group to notify.  Provide `user_id`, `username` or `group_id`. |
| `--group-ids` | string | Group IDs when the notification requires multiple groups. If sent as a string, it should be comma-delimited. |
| `--path` | string | Path |
| `--username` | string | The username of the user to notify.  Provide `user_id`, `username` or `group_id`. |

### `files-cli notifications update`

Update Notification.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Notification ID. **Required.** |
| `--notify-on-copy` | bool | If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads. |
| `--notify-on-delete` | bool | Trigger on files deleted in this path? |
| `--notify-on-download` | bool | Trigger on files downloaded in this path? |
| `--notify-on-move` | bool | Trigger on files moved to this path? |
| `--notify-on-upload` | bool | Trigger on files created/uploaded/updated/changed in this path? |
| `--notify-user-actions` | bool | If `true` actions initiated by the user will still result in a notification |
| `--recursive` | bool | If `true`, enable notifications for each subfolder in this path |
| `--send-interval` | string | The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`. |
| `--subject` | string | Custom subject line to use for notification emails |
| `--message` | string | Custom message to include in notification emails |
| `--triggering-filenames` | []string | Array of filenames (possibly with wildcards) to scope trigger |
| `--triggering-group-ids` | []int64 | If set, will only notify on actions made by a member of one of the specified groups |
| `--triggering-user-ids` | []int64 | If set, will only notify on actions made one of the specified users |
| `--trigger-by-share-recipients` | bool | Notify when actions are performed by a share recipient? |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli notifications delete`

Delete Notification.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Notification ID. **Required.** |

