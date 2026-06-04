---
name: filescom-projects
description: |
  A Project is associated with a folder and add project management features to that folder.
---

# filescom-projects

A Project is associated with a folder and add project management features to that folder.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli projects list`

List Projects.

No flags beyond the global ones.

### `files-cli projects find`

Show Project.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Project ID. **Required.** |

### `files-cli projects create`

Create Project.

| Flag | Type | Description |
| --- | --- | --- |
| `--global-access` | string | Global permissions.  Can be: `none`, `anyone_with_read`, `anyone_with_full`. **Required.** |

### `files-cli projects update`

Update Project.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Project ID. **Required.** |
| `--global-access` | string | Global permissions.  Can be: `none`, `anyone_with_read`, `anyone_with_full`. **Required.** |

### `files-cli projects delete`

Delete Project.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Project ID. **Required.** |

