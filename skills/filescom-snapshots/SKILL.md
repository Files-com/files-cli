---
name: filescom-snapshots
description: |
  Snapshots allow you to create a read-only archive of files at a specific point in time.
---

# filescom-snapshots

Snapshots allow you to create a read-only archive of files at a specific point in time. You can define a snapshot, add files to it, and then finalize it. Once finalized, the snapshot’s contents are immutable.

Each snapshot may have an expiration date. When the expiration date is reached, the snapshot is automatically deleted from the Files.com platform.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli snapshots list`

List Snapshots.

No flags beyond the global ones.

### `files-cli snapshots find`

Show Snapshot.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Snapshot ID. **Required.** |

### `files-cli snapshots create`

Create Snapshot.

| Flag | Type | Description |
| --- | --- | --- |
| `--expires-at` | datetime | When the snapshot expires. |
| `--name` | string | A name for the snapshot. |
| `--paths` | []string | An array of paths to add to the snapshot. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli snapshots finalize`

Finalize Snapshot.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Snapshot ID. **Required.** |

### `files-cli snapshots update`

Update Snapshot.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Snapshot ID. **Required.** |
| `--expires-at` | datetime | When the snapshot expires. |
| `--name` | string | A name for the snapshot. |
| `--paths` | []string | An array of paths to add to the snapshot. |

### `files-cli snapshots delete`

Delete Snapshot.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Snapshot ID. **Required.** |

