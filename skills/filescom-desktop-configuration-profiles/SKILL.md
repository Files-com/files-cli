---
name: filescom-desktop-configuration-profiles
description: |
  A Desktop Configuration Profile centrally defines desktop mount point mappings for users in a Site or Workspace.
---

# filescom-desktop-configuration-profiles

A Desktop Configuration Profile centrally defines desktop mount point mappings for users in a Site or Workspace.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli desktop-configuration-profiles list`

List Desktop Configuration Profiles.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |

### `files-cli desktop-configuration-profiles find`

Show Desktop Configuration Profile.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Desktop Configuration Profile ID. **Required.** |

### `files-cli desktop-configuration-profiles create`

Create Desktop Configuration Profile.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Profile name **Required.** |
| `--mount-mappings` | object | Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there. **Required.** |
| `--workspace-id` | int64 | Workspace ID |
| `--use-for-all-users` | bool | Whether this profile applies to all users in the Workspace by default |
| `--disable-drive-mounting` | bool | Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile |

### `files-cli desktop-configuration-profiles update`

Update Desktop Configuration Profile.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Desktop Configuration Profile ID. **Required.** |
| `--name` | string | Profile name |
| `--workspace-id` | int64 | Workspace ID |
| `--mount-mappings` | object | Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there. |
| `--use-for-all-users` | bool | Whether this profile applies to all users in the Workspace by default |
| `--disable-drive-mounting` | bool | Whether the desktop app should hide drive mounting, prevent new drive mounts, and unmount active drive mounts for users with this profile |

### `files-cli desktop-configuration-profiles delete`

Delete Desktop Configuration Profile.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Desktop Configuration Profile ID. **Required.** |

