---
name: filescom-key-lifecycle-rules
description: |
  A KeyLifecycleRule represents a rule that applies to GPG keys and SSH keys (also called User Public Keys) based on their inactivity or age.
---

# filescom-key-lifecycle-rules

A KeyLifecycleRule represents a rule that applies to GPG keys and SSH keys (also called User Public Keys) based on their inactivity or age.

Keys that have been unused for the specified number of days will be deleted. SSH keys can also be configured to expire after a specified number of days. SSH key expiration applies only to User Public Keys used for inbound SFTP/SSH login, not Remote Server outbound SSH keys.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli key-lifecycle-rules list`

List Key Lifecycle Rules.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `key_type`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `key_type` and `workspace_id`. Valid field combinations are `[ workspace_id, key_type ]`. |

### `files-cli key-lifecycle-rules find`

Show Key Lifecycle Rule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Key Lifecycle Rule ID. **Required.** |

### `files-cli key-lifecycle-rules create`

Create Key Lifecycle Rule.

| Flag | Type | Description |
| --- | --- | --- |
| `--apply-to-all-workspaces` | bool | If true, a default-workspace rule also applies to keys in all workspaces. |
| `--expiration-days` | int64 | Number of days after creation before an SSH key expires. Applies only to SSH keys. |
| `--key-type` | enum | Key type for which the rule will apply (gpg or ssh). One of: `gpg`, `ssh`. |
| `--inactivity-days` | int64 | Number of days of inactivity before the rule applies. |
| `--name` | string | Key Lifecycle Rule name |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli key-lifecycle-rules update`

Update Key Lifecycle Rule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Key Lifecycle Rule ID. **Required.** |
| `--apply-to-all-workspaces` | bool | If true, a default-workspace rule also applies to keys in all workspaces. |
| `--expiration-days` | int64 | Number of days after creation before an SSH key expires. Applies only to SSH keys. |
| `--key-type` | enum | Key type for which the rule will apply (gpg or ssh). One of: `gpg`, `ssh`. |
| `--inactivity-days` | int64 | Number of days of inactivity before the rule applies. |
| `--name` | string | Key Lifecycle Rule name |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli key-lifecycle-rules delete`

Delete Key Lifecycle Rule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Key Lifecycle Rule ID. **Required.** |

