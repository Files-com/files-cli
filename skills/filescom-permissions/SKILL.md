---
name: filescom-permissions
description: |
  A Permission object represents a grant of access permission on a specific Path to a User or Group. Use this skill to answer access questions — who can access a folder, and what a given user or group can access — and to grant or revoke folder access. For a full site-wide permissions report, use a `permission_audit` scheduled export (see `recipe-generating-reports`).
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

## Limitations and considerations

**Higher permission levels include the lower ones.** The Permission types table above lists, for each level, the permissions it automatically implies — `full` already includes `writeonly`, `readonly`, and `list`, and `admin` includes everything. Grant the lowest level that covers the need rather than stacking redundant permissions. The `--permission` flag on `create` accepts `admin`, `full`, `readonly`, `writeonly`, `list`, or `history`.

**A permission is granted to a user, a group, or a partner — provide exactly one.** Supply one of `--user-id`/`--username`, `--group-id`/`--group-name`, or `--partner-id`. Use `--group-ids` (comma-delimited) only when a single permission should require membership in *multiple* groups at once.

**Effective access is the union of direct and group-inherited permissions.** A user's access to a path combines the permissions granted to that user directly with the permissions granted to every group the user belongs to. Adding a permission to a group immediately propagates to all of its members. When answering "what can this user access?", pass `--include-groups=true` so group-inherited permissions are included.

**To revoke access, delete the permission by its ID.** There is no "remove this user from this path" call. List the permissions for the path (or user), find the matching `id`, and `permissions delete --id=…`. See Common patterns below.

**Permissions move with the folder.** Moving a folder carries its permissions (and folder settings) to the new location.

**Permissions are scoped per Workspace.** Under Workspaces, folder permissions apply only within a Workspace's root, and Workspace Administrators manage them for their own Workspace. A user on `workspace_id=0` is granted access to folders in other Workspaces by adding a Permission record — see the Workspace constraints in `CONTEXT.md`.

## Common patterns

### Answer "who can access this folder?"

List the permissions that apply to the path. The `--path` filter scopes results to that path *and* the ancestor-folder permissions that also grant access to it:

    files-cli permissions list --path="/finance/reports" --format=json

### Answer "what can this user access?"

Filter by the user and include the permissions inherited from the user's groups:

    files-cli permissions list --user-id=USER_ID --include-groups=true --format=json

For a group, filter by the group instead:

    files-cli permissions list --group-id=GROUP_ID --format=json

### Grant access

    files-cli permissions create --path="/finance/reports" --user-id=USER_ID --permission=readonly --recursive=true --format=json

### Revoke access

Permissions are removed by ID — find the permission first, then delete it:

    files-cli permissions list --path="/finance/reports" --format=json
    files-cli permissions delete --id=PERMISSION_ID --format=json

### Produce a full permissions audit

To answer a specific access question, use `permissions list` as shown above — that data comes straight back. For a recurring site-wide report of every permission delivered to a person, create a `permission_audit` scheduled export:

    files-cli scheduled-exports create --name="Permission audit" --export-type=permission_audit --user-id=SITE_ADMIN_USER_ID --export-options='{"group_by":"user"}' --trigger=daily --format=json

`group_by` accepts `user` (one row per user per path they can reach) or `path` (one row per path per user or group granted access). **The report is emailed to the Site Admin in `--user-id` and is not readable back through the API** — a `ScheduledExport` has no `results_url` or results endpoint — so don't create one to obtain the data yourself; use `permissions list` for that. See `recipe-generating-reports`.
