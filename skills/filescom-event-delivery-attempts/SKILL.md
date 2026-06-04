---
name: filescom-event-delivery-attempts
description: |
  An EventDeliveryAttempt records delivery state for an EventRecord and EventTarget.
---

# filescom-event-delivery-attempts

An EventDeliveryAttempt records delivery state for an EventRecord and EventTarget.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli event-delivery-attempts list`

List Event Delivery Attempts.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `status`, `event_record_id`, `event_target_id` or `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `status`, `workspace_id`, `event_record_id` or `event_target_id`. Valid field combinations are `[ workspace_id, status ]`, `[ workspace_id, event_record_id ]` or `[ workspace_id, event_target_id ]`. |

### `files-cli event-delivery-attempts find`

Show Event Delivery Attempt.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Delivery Attempt ID. **Required.** |

