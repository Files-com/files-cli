---
name: filescom-permissions
description: |
  A Permission object represents a grant of access permission on a specific Path to a User or Group.
---

# filescom-permissions

A Permission object represents a grant of access permission on a specific Path to a User or Group.

They can be optionally recursive or nonrecursive into the subfolders of that path.

A Permission may be applied to a User *or* a Group, but not both at once.

The following table sets forth the available Permission types:

| Permission | Access Level Granted | Automatically Also Includes/Implies Permissions |
| --- | ----------- | --------------------- |
| `admin` | Able to manage Folder Behaviors, Permissions, and Notifications for the folder.  Also grants all other permissions. | `bundle`, `full`, `writeonly`, `readonly`, `list`, `history` |
| `bundle` | Able to share files and folders via a Bundle (share link). | `readonly`, `list` |
| `full` | Able to read, write, move, delete, and rename files and folders. Also grants the ability to overwrite files upon upload. | `writeonly`, `readonly`, `list` |
| `history` | Able to view the history of files and folders and to create email notifications for themselves. | `list` |
| `list` | Able to list files and folders, but not download. | none |
| `readonly` | Able to list, preview, and download files and folders. | `list` |
| `readonly_site_admin` | Able to behave as a read-only Site Admin on a given child site. Only applies to child sites. | `readonly`, `list`, `history` |
| `site_admin` | Able to behave as a Site Admin on a given child site. Only applies to child sites. | `bundle`, `full`, `writeonly`, `readonly`, `list`, `history` |
| `writeonly` | Able to upload files, create folders and list subfolders the user has write permission to. | none |

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli permissions list [path]`

List Permissions.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `group_id`, `path`, `user_id`, `partner_id` or `id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `group_id`, `partner_id` or `user_id`. Valid field combinations are `[ group_id, path ]`, `[ partner_id, path ]`, `[ user_id, path ]`, `[ user_id, group_id ]`, `[ user_id, group_id, path ]`, `[ user_id, group_id, partner_id ]` or `[ user_id, group_id, partner_id, path ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. |
| `--path` | string | Permission path.  If provided, will scope all permissions(including upward) to this path. |
| `--include-groups` | bool | If searching by user or group, also include user's permissions that are inherited from its groups? |
| `--group-id` | string | (no description) |
| `--partner-id` | string | (no description) |
| `--user-id` | string | (no description) |

### `files-cli permissions create [path]`

Create Permission.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Folder path **Required.** |
| `--group-id` | int64 | Group ID. Provide `group_name` or `group_id` |
| `--group-ids` | string | Group IDs when the permission requires multiple groups. If sent as a string, it should be comma-delimited. |
| `--permission` | string | Permission type.  Can be `admin`, `full`, `readonly`, `writeonly`, `list`, or `history` |
| `--recursive` | bool | Apply to subfolders recursively? |
| `--partner-id` | int64 | Partner ID if this Permission belongs to a partner. |
| `--user-id` | int64 | User ID.  Provide `username` or `user_id` |
| `--username` | string | User username.  Provide `username` or `user_id` |
| `--group-name` | string | Group name.  Provide `group_name` or `group_id` |
| `--site-id` | int64 | Site ID. If not provided, will default to current site. Used when creating a permission for a child site. |

### `files-cli permissions delete`

Delete Permission.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Permission ID. **Required.** |

