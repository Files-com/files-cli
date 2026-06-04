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

