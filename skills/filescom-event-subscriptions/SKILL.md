---
name: filescom-event-subscriptions
description: |
  An EventSubscription selects EventRecords for an EventChannel and sends them to one or more EventTargets.
---

# filescom-event-subscriptions

An EventSubscription selects EventRecords for an EventChannel and sends them to one or more EventTargets.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

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
| `--subject` | string | Custom subject line to use for notification emails. |
| `--message` | string | Custom message to include in notification emails. |
| `--message-only` | bool | If true, notification email bodies contain only the custom message, omitting event details and the review button. Requires a custom message, defaults to false, and does not affect non-email targets. |
| `--enabled` | bool | Whether this Event Subscription can dispatch events. |
| `--event-types` | []string | Event type strings matched by this subscription. Blank means all event types. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1` |
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
| `--subject` | string | Custom subject line to use for notification emails. |
| `--message` | string | Custom message to include in notification emails. |
| `--message-only` | bool | If true, notification email bodies contain only the custom message, omitting event details and the review button. Requires a custom message, defaults to false, and does not affect non-email targets. |
| `--enabled` | bool | Whether this Event Subscription can dispatch events. |
| `--event-types` | []string | Event type strings matched by this subscription. Blank means all event types. Valid values: `automation_run.canceled.v1`, `automation_run.failure.v1`, `automation_run.failure_will_retry.v1`, `automation_run.partial_failure.v1`, `automation_run.partial_failure_will_retry.v1`, `automation_run.skipped.v1`, `automation_run.success.v1`, `expectation_evaluation.invalid.v1`, `expectation_evaluation.late.v1`, `expectation_evaluation.missing.v1`, `expectation_evaluation.success.v1`, `expectation_incident.acknowledged.v1`, `expectation_incident.open.v1`, `expectation_incident.resolved.v1`, `expectation_incident.snoozed.v1`, `external_event.client_log.failure.v1`, `external_event.client_log.partial_failure.v1`, `external_event.client_log.skipped.v1`, `external_event.client_log.success.v1`, `pending_work_event.failure.v1`, `pending_work_event.partial_failure.v1`, `pending_work_event.skipped.v1`, `pending_work_event.success.v1`, `siem_http_destination_event.failure.v1`, `siem_http_destination_event.partial_failure.v1`, `siem_http_destination_event.skipped.v1`, `siem_http_destination_event.success.v1`, `sso_event.ldap_login.failure.v1`, `sso_event.ldap_login.partial_failure.v1`, `sso_event.ldap_login.skipped.v1`, `sso_event.ldap_login.success.v1`, `sso_event.ldap_sync.failure.v1`, `sso_event.ldap_sync.partial_failure.v1`, `sso_event.ldap_sync.skipped.v1`, `sso_event.ldap_sync.success.v1`, `sso_event.saml_login.failure.v1`, `sso_event.saml_login.partial_failure.v1`, `sso_event.saml_login.skipped.v1`, `sso_event.saml_login.success.v1`, `sync_run.failure.v1`, `sync_run.partial_failure.v1`, `sync_run.skipped.v1`, `sync_run.success.v1`, `user_security_event.lockout.v1` |
| `--filter` | object | Structured event payload filter. |
| `--delivery-policy` | object | Event Subscription delivery policy. |
| `--event-target-ids` | []int64 | Event Target IDs this subscription sends to. |

### `files-cli event-subscriptions delete`

Delete Event Subscription.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Event Subscription ID. **Required.** |

