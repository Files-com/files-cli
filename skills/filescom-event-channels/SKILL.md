---
name: filescom-event-channels
description: |
  An EventChannel is a named grouping of EventSubscriptions.
---

# filescom-event-channels

An EventChannel is a named grouping of EventSubscriptions.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli event-channels list`

List Event Channels.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `enabled`, `default_channel` or `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `enabled`, `default_channel` or `workspace_id`. Valid field combinations are `[ workspace_id, enabled ]` and `[ workspace_id, default_channel ]`. |

### `files-cli event-channels find`

Show Event Channel.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Channel ID. **Required.** |

### `files-cli event-channels create`

Create Event Channel.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Event Channel name. **Required.** |
| `--workspace-id` | int64 | Workspace ID. 0 means the default workspace. |
| `--description` | string | Event Channel description. |
| `--enabled` | bool | Whether this Event Channel can dispatch events. |
| `--default-channel` | bool | Whether this Event Channel is the default destination for newly published events. |

### `files-cli event-channels update`

Update Event Channel.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Channel ID. **Required.** |
| `--name` | string | Event Channel name. |
| `--workspace-id` | int64 | Workspace ID. 0 means the default workspace. |
| `--description` | string | Event Channel description. |
| `--enabled` | bool | Whether this Event Channel can dispatch events. |
| `--default-channel` | bool | Whether this Event Channel is the default destination for newly published events. |

### `files-cli event-channels delete`

Delete Event Channel.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Channel ID. **Required.** |

