---
name: filescom-expectation-incidents
description: |
  An ExpectationIncident groups ongoing failure behavior for an Expectation over time.
---

# filescom-expectation-incidents

An ExpectationIncident groups ongoing failure behavior for an Expectation over time.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli expectation-incidents list`

List Expectation Incidents.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `expectation_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `expectation_id` and `workspace_id`. Valid field combinations are `[ workspace_id, expectation_id ]`. |

### `files-cli expectation-incidents find`

Show Expectation Incident.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation Incident ID. **Required.** |

### `files-cli expectation-incidents resolve`

Resolve an expectation incident.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation Incident ID. **Required.** |

### `files-cli expectation-incidents snooze`

Snooze an expectation incident until a specified time.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation Incident ID. **Required.** |
| `--snoozed-until` | datetime | Time until which the incident should remain snoozed. **Required.** |

### `files-cli expectation-incidents acknowledge`

Acknowledge an expectation incident.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation Incident ID. **Required.** |

