---
name: filescom-automation-runs
description: |
  An AutomationRun is a record with information about a single execution of a given Automation.
---

# filescom-automation-runs

An AutomationRun is a record with information about a single execution of a given Automation.

Automation Runs can be retried upon `failure` or `partial_failure` by setting the `retry_on_failure_interval_in_minutes` and `retry_on_failure_number_of_attempts` settings on the associated Automation.

When retries occur, a new AutomationRun will be created for each retry.  The property `retry_at` on the original run, if set, represents when that run will be retried next.  The properties `retried_at` and `retried_in_run_id` will be set in the original run that failed, at the time of retry.  The property `retry_of_run_id` will be set in the new run.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli automation-runs list`

List Automation Runs.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `automation_id`, `created_at` or `status`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `status`, `workspace_id` or `automation_id`. Valid field combinations are `[ workspace_id, status ]`, `[ automation_id, status ]`, `[ workspace_id, automation_id ]` or `[ workspace_id, automation_id, status ]`. |
| `--automation-id` | int64 | ID of the associated Automation. **Required.** |

### `files-cli automation-runs find`

Show Automation Run.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Automation Run ID. **Required.** |

