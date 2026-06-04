---
name: filescom-group-users
description: |
  A GroupUser is a record about membership of a User within a Group.
---

# filescom-group-users

A GroupUser is a record about membership of a User within a Group.

## Creating GroupUsers
GroupUsers can be created via the normal `create` action. When using the `update` action, if the
GroupUser record does not exist for the given user/group IDs it will be created.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli group-users list`

List Group Users.

| Flag | Type | Description |
| --- | --- | --- |
| `--group-id` | int64 | Group ID.  If provided, will return group_users of this group. |
| `--user-id` | int64 | User ID.  If provided, will return group_users of this user. |

### `files-cli group-users create`

Create Group User.

| Flag | Type | Description |
| --- | --- | --- |
| `--group-id` | int64 | Group ID to add user to. **Required.** |
| `--user-id` | int64 | User ID to add to group. **Required.** |
| `--admin` | bool | Is the user a group administrator? |

### `files-cli group-users update`

Update Group User.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Group User ID. **Required.** |
| `--group-id` | int64 | Group ID to add user to. **Required.** |
| `--user-id` | int64 | User ID to add to group. **Required.** |
| `--admin` | bool | Is the user a group administrator? |

### `files-cli group-users delete`

Delete Group User.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Group User ID. **Required.** |
| `--group-id` | int64 | Group ID from which to remove user. **Required.** |
| `--user-id` | int64 | User ID to remove from group. **Required.** |

