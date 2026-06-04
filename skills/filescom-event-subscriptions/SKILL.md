---
name: filescom-event-subscriptions
description: |
  An EventSubscription selects EventRecords for an EventChannel and sends them to one or more EventTargets.
---

# filescom-event-subscriptions

An EventSubscription selects EventRecords for an EventChannel and sends them to one or more EventTargets.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli event-subscriptions list`

List Event Subscriptions.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`, `enabled`, `event_channel_id` or `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `enabled`, `event_channel_id` or `workspace_id`. Valid field combinations are `[ enabled, event_channel_id ]`, `[ workspace_id, enabled ]` or `[ workspace_id, enabled, event_channel_id ]`. |

### `files-cli event-subscriptions find`

Show Event Subscription.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Subscription ID. **Required.** |

### `files-cli event-subscriptions create`

Create Event Subscription.

| Flag | Type | Description |
| --- | --- | --- |
| `--event-channel-id` | int64 | Event Channel ID |
| `--workspace-id` | int64 | Workspace ID. 0 means the default workspace or site-wide. |
| `--apply-to-all-workspaces` | bool | If true, this default-workspace subscription applies to events from all workspaces. |
| `--name` | string | Event Subscription name. **Required.** |
| `--enabled` | bool | Whether this Event Subscription can dispatch events. |
| `--event-types` | []string | Event type strings matched by this subscription. Blank means all event types. |
| `--filter` | object | Structured event payload filter. |
| `--delivery-policy` | object | Event Subscription delivery policy. |
| `--event-target-ids` | []int64 | Event Target IDs this subscription sends to. |

### `files-cli event-subscriptions update`

Update Event Subscription.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Subscription ID. **Required.** |
| `--event-channel-id` | int64 | Event Channel ID |
| `--workspace-id` | int64 | Workspace ID. 0 means the default workspace or site-wide. |
| `--apply-to-all-workspaces` | bool | If true, this default-workspace subscription applies to events from all workspaces. |
| `--name` | string | Event Subscription name. |
| `--enabled` | bool | Whether this Event Subscription can dispatch events. |
| `--event-types` | []string | Event type strings matched by this subscription. Blank means all event types. |
| `--filter` | object | Structured event payload filter. |
| `--delivery-policy` | object | Event Subscription delivery policy. |
| `--event-target-ids` | []int64 | Event Target IDs this subscription sends to. |

### `files-cli event-subscriptions delete`

Delete Event Subscription.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Subscription ID. **Required.** |

