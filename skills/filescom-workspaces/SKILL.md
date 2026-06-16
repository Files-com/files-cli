---
name: filescom-workspaces
description: |
  A Workspace is a lightweight way to organize related resources inside a single Files.com Site.
---

# filescom-workspaces

A Workspace is a lightweight way to organize related resources inside a single Files.com Site.

Customers commonly group resources by project, department, client, or region. Workspaces provide a built-in structure for that grouping, so the UI can operate within a clear “workspace context” and admins can delegate management for a subset of resources without requiring full site-level isolation.

Every Site has an implicit Default workspace (ID 0). Resources that are not explicitly assigned to a named workspace are considered part of the Default workspace.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli workspaces list`

List Workspaces.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `name`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |

### `files-cli workspaces find`

Show Workspace.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Workspace ID. **Required.** |

### `files-cli workspaces create`

Create Workspace.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Workspace name |

### `files-cli workspaces update`

Update Workspace.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Workspace ID. **Required.** |
| `--name` | string | Workspace name |

### `files-cli workspaces delete`

Delete Workspace.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Workspace ID. **Required.** |

## Limitations and considerations

**A Workspace is just an organizational container — `id` and `name`.** Resources (users, folders, automations, remote servers, and so on) belong to a Workspace by carrying its ID; the Default workspace is ID `0`; any other workspace ID is a named workspace. Users in a named Workspace are confined to that Workspace, while a Site Administrator is not confined to any Workspace. A user belonging to a Workspace is Workspace Admin for it if they have the `workspace_admin` flag enabled.

**Resources cannot be moved between Workspaces.** There is no operation that reassigns a folder, automation, remote server, or any other resource from one Workspace to another. Do not try to. The only supported Workspace reassignment is a Site Administrator returning a *user* to the Default workspace by changing that user's `workspace_id` from a non-zero value to `0` with `files-cli users update --id=USER_ID --workspace-id=0`.

**Grant cross-Workspace access with Permissions, not by reassigning a Workspace.** A user on `workspace_id=0` is not confined to the Default Workspace. To give them access to resources that belong to other Workspaces, add a Permission record for that user on the relevant path (`files-cli permissions create --path=PATH --user-id=USER_ID --permission=LEVEL`) rather than changing their Workspace ID assignment, where PATH is prefixed like `_/Workspaces/$WORKSPACE_ID/$FOLDER_PATH`. If $FOLDER_PATH is empty, the permission is for the Workspace's root folder. If `admin` permission is granted on a Workspace's root folder, this grants access as Workspace Admin to the entire Workspace. See also the `filescom-permissions` skill.

**Deleting a Workspace deletes everything in it.** Removing a Workspace removes all resources that belong to it, such as users, groups, partners, folders, automations, remote servers, syncs, and notifications.

