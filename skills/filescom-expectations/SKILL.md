---
name: filescom-expectations
description: |
  Expectations let your Files.com site define what “correct” file delivery looks like, continuously evaluate whether it happened, and keep history when it did not.
---

# filescom-expectations

Expectations let your Files.com site define what “correct” file delivery looks like, continuously evaluate whether it happened, and keep history when it did not.

Expectations are meant to answer operational questions like:

* Did the expected file arrive?
* Was it on time?
* Did it meet the required shape and count rules?
* Is there an active issue someone needs to acknowledge?

Expectations are different from Automations and Syncs. Automations and Syncs act on files; Expectations monitor whether expected files arrived on time, in the right place, and in the right shape. In practice, Expectations are the sensor and Automations are the actuator.

An Expectation combines four concepts:

1. **Scope**: where to look for candidate files, using `path`, `source`, and optional `exclude_pattern`.
2. **Trigger / timing**: when a window opens and how long it stays eligible, using `trigger`, schedule fields, `lookback_interval`, `late_acceptance_interval`, `inactivity_interval`, and `max_open_interval`.
3. **Criteria**: what must be true for the window to succeed, using the structured `criteria` JSON document.
4. **Outcome history**: what happened over time, exposed through `ExpectationEvaluation` history and `ExpectationIncident` lifecycle records.

## Scope and matching

Expectations reuse the familiar Files.com path-plus-glob model.

The `path` field identifies the folder scope, while `source` identifies which files within that scope are candidates. `exclude_pattern` removes files from consideration.

Like Automations, these fields support glob-style matching. Expectations treat those matches as one logical candidate set for each window. A single Expectation does not implicitly fan out into separate per-customer or per-folder evaluations just because the path contains wildcards.

## Expectation windows

Expectations are evaluated in windows.

Each window is persisted as an `ExpectationEvaluation` record. A window opens, remains `open` while evidence can still arrive, and then closes into a terminal result such as `success`, `late`, `missing`, or `invalid`.

An Expectation has only one open window at a time.

## Trigger modes

Expectations can open windows in three ways:

* `daily`: run on a recurring daily/weekly/monthly/quarterly/yearly cadence using `interval` and `recurring_day`.
* `custom_schedule`: run on specific weekdays and times using `schedule_days_of_week`, `schedule_times_of_day`, and optional `schedule_time_zone` / `holiday_region`.
* `manual`: an operator explicitly opens the window.

Schedule-driven expectations define an on-time deadline and may optionally remain eligible to close as `late` during `late_acceptance_interval`.

Manual expectations have no concept of `late`; they open when triggered and close based on inactivity or hard-stop timing.

## Success criteria

The `criteria` field is a structured JSON object describing what counts as success for the window.

Criteria v1 can express things like:

* file count constraints
* total byte constraints
* allowed extensions
* filename regex validation
* forbidden files
* required named or globbed files with their own per-file constraints

Criteria v2 adds `content_validation`, which runs a customer-authored Files Transform Script in either `per_file` or `whole_batch` mode. Per-file scripts receive the file contents parsed by FTS as `payload`. Whole-batch scripts receive an array of file objects containing `path`, `name`, `size`, `last_modified_at`, and each file's parsed `payload`.

A content-validation script returns `true` or `{ success: true }` to pass. It returns `false` or `{ success: false, errors: [...] }` to fail. Error entries may be strings or structured objects with values such as `message`, `field`, `row`, `expected`, and `actual`; these details are preserved in readable form in the Evaluation's `criteria_errors`. Script, parsing, download, and size-limit errors also fail the criterion. Each file is limited to 100 MB, and whole-batch mode additionally limits the combined raw input to 100 MB.

Required file rule keys may also include standard strftime-style date/time tokens like `%Y`, `%m`, and `%d`. Those tokens are resolved at evaluation time using a stable window anchor: schedule-driven expectations use the window's `deadline_at`, while manual and upload expectations use the window's `opened_at`.

## History and incidents

The Expectation itself stores summary state like `last_evaluated_at`, `last_success_at`, `last_failure_at`, and `last_result`.

For deeper inspection:

* `ExpectationEvaluation` history shows each open or closed window and the evidence captured for it.
* `ExpectationIncident` records track ongoing failure situations over time, including acknowledge, snooze, and resolve actions.

Manual windows do not open incidents in v1. Schedule-driven failures can open incidents, and later qualifying success can resolve them.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli expectations list`

List Expectations.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `disabled`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled` and `workspace_id`. Valid field combinations are `[ workspace_id, disabled ]`. |

### `files-cli expectations find`

Show Expectation.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation ID. **Required.** |

### `files-cli expectations create [path]`

Create Expectation.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Expectation name. |
| `--description` | string | Expectation description. |
| `--path` | string | Path scope for the expectation. Supports workspace-relative presentation. |
| `--source` | string | Source glob used to select candidate files. |
| `--exclude-pattern` | string | Optional source exclusion glob. |
| `--disabled` | bool | If true, the expectation is disabled. |
| `--trigger` | enum | How this expectation opens windows. One of: `manual`, `upload`, `daily`, `custom_schedule`. |
| `--interval` | string | If trigger is `daily`, this specifies how often to run the expectation. |
| `--recurring-day` | int64 | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `--schedule-times-of-day` | []string | Times of day in HH:MM format for schedule-driven expectations. |
| `--schedule-time-zone` | string | Time zone used by the expectation schedule. |
| `--holiday-region` | string | Optional holiday region used by schedule-driven expectations. |
| `--lookback-interval` | int64 | How many seconds before the due boundary the window starts. |
| `--late-acceptance-interval` | int64 | How many seconds a schedule-driven window may remain eligible to close as late. |
| `--inactivity-interval` | int64 | How many quiet seconds are required before final closure. |
| `--max-open-interval` | int64 | Hard-stop duration in seconds for unscheduled expectations. |
| `--criteria` | object | Versioned success criteria definition for the expectation, including optional Files Transform Script content validation in criteria v2. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli expectations trigger-evaluation`

Manually open an Expectation window.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation ID. **Required.** |

### `files-cli expectations update [path]`

Update Expectation.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation ID. **Required.** |
| `--name` | string | Expectation name. |
| `--description` | string | Expectation description. |
| `--path` | string | Path scope for the expectation. Supports workspace-relative presentation. |
| `--source` | string | Source glob used to select candidate files. |
| `--exclude-pattern` | string | Optional source exclusion glob. |
| `--disabled` | bool | If true, the expectation is disabled. |
| `--trigger` | enum | How this expectation opens windows. One of: `manual`, `upload`, `daily`, `custom_schedule`. |
| `--interval` | string | If trigger is `daily`, this specifies how often to run the expectation. |
| `--recurring-day` | int64 | If trigger is `daily`, this selects the day number inside the chosen interval. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`, the 0-based weekdays used by the schedule. |
| `--schedule-times-of-day` | []string | Times of day in HH:MM format for schedule-driven expectations. |
| `--schedule-time-zone` | string | Time zone used by the expectation schedule. |
| `--holiday-region` | string | Optional holiday region used by schedule-driven expectations. |
| `--lookback-interval` | int64 | How many seconds before the due boundary the window starts. |
| `--late-acceptance-interval` | int64 | How many seconds a schedule-driven window may remain eligible to close as late. |
| `--inactivity-interval` | int64 | How many quiet seconds are required before final closure. |
| `--max-open-interval` | int64 | Hard-stop duration in seconds for unscheduled expectations. |
| `--criteria` | object | Versioned success criteria definition for the expectation, including optional Files Transform Script content validation in criteria v2. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli expectations delete`

Delete Expectation.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation ID. **Required.** |

