---
name: filescom-remote-mount-backends
description: |
  A Remote Mount Backend is used to provide high availability for a Remote Server Mount Folder Behavior.
---

# filescom-remote-mount-backends

A Remote Mount Backend is used to provide high availability for a Remote Server Mount Folder Behavior.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli remote-mount-backends list`

List Remote Mount Backends.

| Flag | Type | Description |
| --- | --- | --- |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `remote_server_mount_id`. |

### `files-cli remote-mount-backends find`

Show Remote Mount Backend.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Mount Backend ID. **Required.** |

### `files-cli remote-mount-backends create`

Create Remote Mount Backend.

| Flag | Type | Description |
| --- | --- | --- |
| `--enabled` | bool | True if this backend is enabled. |
| `--fall` | int64 | Number of consecutive failures before considering the backend unhealthy. |
| `--health-check-enabled` | bool | True if health checks are enabled for this backend. |
| `--health-check-type` | enum | Type of health check to perform. One of: `active`, `passive`. |
| `--interval` | int64 | Interval in seconds between health checks. |
| `--min-free-cpu` | double | Minimum free CPU percentage required for this backend to be considered healthy. |
| `--min-free-mem` | double | Minimum free memory percentage required for this backend to be considered healthy. |
| `--priority` | int64 | Priority of this backend. |
| `--remote-path` | string | Path on the remote server to treat as the root of this mount. |
| `--rise` | int64 | Number of consecutive successes before considering the backend healthy. |
| `--canary-file-path` | string | Path to the canary file used for health checks. **Required.** |
| `--remote-server-mount-id` | int64 | The mount ID of the Remote Server Mount that this backend is associated with. **Required.** |
| `--remote-server-id` | int64 | The remote server that this backend is associated with. **Required.** |

### `files-cli remote-mount-backends reset-status`

Reset backend status to healthy.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Mount Backend ID. **Required.** |

### `files-cli remote-mount-backends update`

Update Remote Mount Backend.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Mount Backend ID. **Required.** |
| `--enabled` | bool | True if this backend is enabled. |
| `--fall` | int64 | Number of consecutive failures before considering the backend unhealthy. |
| `--health-check-enabled` | bool | True if health checks are enabled for this backend. |
| `--health-check-type` | enum | Type of health check to perform. One of: `active`, `passive`. |
| `--interval` | int64 | Interval in seconds between health checks. |
| `--min-free-cpu` | double | Minimum free CPU percentage required for this backend to be considered healthy. |
| `--min-free-mem` | double | Minimum free memory percentage required for this backend to be considered healthy. |
| `--priority` | int64 | Priority of this backend. |
| `--remote-path` | string | Path on the remote server to treat as the root of this mount. |
| `--rise` | int64 | Number of consecutive successes before considering the backend healthy. |
| `--canary-file-path` | string | Path to the canary file used for health checks. |
| `--remote-server-id` | int64 | The remote server that this backend is associated with. |

### `files-cli remote-mount-backends delete`

Delete Remote Mount Backend.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Mount Backend ID. **Required.** |

