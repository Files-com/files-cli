---
name: filescom-bundle-notifications
description: |
  A BundleNotification is an E-Mail sent out to users when certain actions are performed on or within a shared set of files and folders.
---

# filescom-bundle-notifications

A BundleNotification is an E-Mail sent out to users when certain actions are performed on or within a shared set of files and folders.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli bundle-notifications list`

List Share Link Notifications.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `bundle_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `bundle_id`. |
| `--bundle-id` | int64 | Bundle ID |

### `files-cli bundle-notifications find`

Show Share Link Notification.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Bundle Notification ID. **Required.** |

### `files-cli bundle-notifications create`

Create Share Link Notification.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--bundle-id` | int64 | Bundle ID to notify on **Required.** |
| `--notify-user-id` | int64 | The id of the user to notify. |
| `--notify-on-registration` | bool | Triggers bundle notification when a registration action occurs for it. |
| `--notify-on-upload` | bool | Triggers bundle notification when a upload action occurs for it. |

### `files-cli bundle-notifications update`

Update Share Link Notification.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Bundle Notification ID. **Required.** |
| `--notify-on-registration` | bool | Triggers bundle notification when a registration action occurs for it. |
| `--notify-on-upload` | bool | Triggers bundle notification when a upload action occurs for it. |

### `files-cli bundle-notifications delete`

Delete Share Link Notification.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Bundle Notification ID. **Required.** |

