---
name: filescom-event-targets
description: |
  An EventTarget is a delivery destination for EventRecords.
---

# filescom-event-targets

An EventTarget is a delivery destination for EventRecords.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli event-targets list`

List Event Targets.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `enabled` or `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `enabled`, `target_type` or `workspace_id`. Valid field combinations are `[ enabled, target_type ]`, `[ workspace_id, enabled ]` or `[ workspace_id, enabled, target_type ]`. |

### `files-cli event-targets find`

Show Event Target.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Target ID. **Required.** |

### `files-cli event-targets create`

Create Event Target.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Event Target name. **Required.** |
| `--workspace-id` | int64 | Workspace ID. 0 means the default workspace or site-wide. |
| `--apply-to-all-workspaces` | bool | If true, this default-workspace target can receive events from all workspaces. |
| `--target-type` | enum | Event Target type. One of: `email`, `webhook`, `slack_webhook`, `teams_webhook`, `amazon_sns`, `google_pubsub`. **Required.** |
| `--enabled` | bool | Whether this Event Target can receive events. |
| `--config` | object | Event Target configuration. **Required.** |
| `--delivery-policy` | object | Event Target delivery policy. Email targets support batch_interval in seconds, between 600 and 86400. |

### `files-cli event-targets update`

Update Event Target.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Target ID. **Required.** |
| `--name` | string | Event Target name. |
| `--workspace-id` | int64 | Workspace ID. 0 means the default workspace or site-wide. |
| `--apply-to-all-workspaces` | bool | If true, this default-workspace target can receive events from all workspaces. |
| `--target-type` | enum | Event Target type. One of: `email`, `webhook`, `slack_webhook`, `teams_webhook`, `amazon_sns`, `google_pubsub`. |
| `--enabled` | bool | Whether this Event Target can receive events. |
| `--config` | object | Event Target configuration. |
| `--delivery-policy` | object | Event Target delivery policy. Email targets support batch_interval in seconds, between 600 and 86400. |

### `files-cli event-targets delete`

Delete Event Target.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Target ID. **Required.** |

