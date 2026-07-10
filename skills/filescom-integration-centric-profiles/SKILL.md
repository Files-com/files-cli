---
name: filescom-integration-centric-profiles
description: |
  An Integration Centric Profile defines the Remote Server integrations a user is expected to add and connect during integration-centric onboarding.
---

# filescom-integration-centric-profiles

An Integration Centric Profile defines the Remote Server integrations a user is expected to add and connect during integration-centric onboarding.

Use this to automate setup guidance for users who need access to multiple business systems without sending long manual instructions. Common scenarios include ongoing access to systems such as SharePoint, bridging Google, Microsoft, and Box environments after M&A activity, and migrations where users connect legacy EFSS accounts during transition work.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli integration-centric-profiles list`

List Integration Centric Profiles.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |

### `files-cli integration-centric-profiles find`

Show Integration Centric Profile.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Integration Centric Profile ID. **Required.** |

### `files-cli integration-centric-profiles create`

Create Integration Centric Profile.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Profile name **Required.** |
| `--expected-remote-servers` | []object | Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`. **Required.** |
| `--workspace-id` | int64 | Workspace ID |
| `--use-for-all-users` | bool | Whether this profile applies to all users in the Workspace by default |

### `files-cli integration-centric-profiles update`

Update Integration Centric Profile.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Integration Centric Profile ID. **Required.** |
| `--name` | string | Profile name |
| `--workspace-id` | int64 | Workspace ID |
| `--expected-remote-servers` | []object | Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`. |
| `--use-for-all-users` | bool | Whether this profile applies to all users in the Workspace by default |

### `files-cli integration-centric-profiles delete`

Delete Integration Centric Profile.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Integration Centric Profile ID. **Required.** |

