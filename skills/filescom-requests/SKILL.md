---
name: filescom-requests
description: |
  A Request is a file that *should* be uploaded by a specific user or group.
---

# filescom-requests

A Request is a file that *should* be uploaded by a specific user or group.

Requests can either be manually created and managed, or managed automatically by an Automation.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli requests list [path]`

List Requests.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `--mine` | bool | Only show requests of the current user?  (Defaults to true if current user is not a site admin.) |
| `--path` | string | Path to show requests for.  If omitted, shows all paths. Send `/` to represent the root directory. |

### `files-cli requests get-folder [path]`

List Requests.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `--mine` | bool | Only show requests of the current user?  (Defaults to true if current user is not a site admin.) |
| `--path` | string | Path to show requests for.  If omitted, shows all paths. Send `/` to represent the root directory. **Required.** |

### `files-cli requests create [path]`

Create Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Folder path on which to request the file. **Required.** |
| `--destination` | string | Destination filename (without extension) to request. **Required.** |
| `--user-ids` | string | A list of user IDs to request the file from. If sent as a string, it should be comma-delimited. |
| `--group-ids` | string | A list of group IDs to request the file from. If sent as a string, it should be comma-delimited. |

### `files-cli requests delete`

Delete Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Request ID. **Required.** |

