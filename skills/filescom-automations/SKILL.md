---
name: filescom-automations
description: |
  An Automation is an automated process of controlling workflows on your Files.com site.
---

# filescom-automations

An Automation is an automated process of controlling workflows on your Files.com site.

Automations are different from Behaviors because Behaviors are associated with a current folder, while Automations apply across your entire site.

Automations are never removed when folders are removed, while Behaviors are removed when the associated folder is removed.

## Path Matching

The `path` attribute specifies which folders this automation applies to.
It gets combined with the `source` attribute to determine which files are actually affected by the automation.
Note that the `path` attribute supports globs, and only refers to _folders_.
It's the `source` attribute, which also supports globs, combined with the `path` attribute that determines which files are affected, and automations only operate on the files themselves.
Additionally, paths in Automations can refer to folders which don't yet exist.

### Path Globs

Although Automations may have a `path` specified, it can be a glob (which includes wildcards), which affects multiple folders.

`*` matches any folder at that level of the path, but not subfolders.  For example, `path/to/*` matches `path/to/folder1` and `path/to/folder2`, but not `path/to/folder1/subfolder`.

`**` matches subfolders recursively.  For example, `path/to/**` matches `path/to/folder1`, `path/to/folder1/subfolder`, `path/to/folder2`, `path/to/folder2/subfolder`, etc.

`?` matches any one character.

Use square brackets `[]` to match any character from a set. This works like a regular expression, including negation using `^`.

Curly brackets `{}` can be used to denote parts of a pattern which will accept a number of alternatives, separated by commas `,`.
These alternatives can either be literal text or include special characters including nested curly brackets.
For example `{Mon,Tue,Wed,Thu,Fri}` would match abbreviated weekdays, and `202{3-{0[7-9],1?},4-0[1-6]}-*` would match dates from `2023-07-01` through `2024-06-30`.

To match any of the special characters literally, precede it with a backslash and enclose that pair with square brackets. For example to match a literal `?`, use `[\?]`.

Globs are supported on `path`, `source`, and `exclude_pattern` fields. Globs are not supported on remote paths of any kind or for any field.

By default, Copy and Move automations that use globs will implicitly replicate matched folder structures at the destination.  If you want to flatten the folder structure, set `flatten_destination_structure` to `true`.

## Automation Triggers

Automations can be triggered in the following ways:

* `custom_schedule` : The automation will run according to the custom schedule parameters for `days_of_week` (0-based) and `times_of_day`. A time zone may be specified via `time_zone` in Rails TimeZone name format.
* `daily` : The automation will run once in a picked `interval`. You can specify `recurring_day` to select day number inside a picked `interval` it should be run on.
* `webhook` : the automation will run when a request is sent to the corresponding webhook URL.
* `action` : The automation will run when a specific action happens, e.g. a file is created or downloaded.

Future enhancements will allow Automations to be triggered by an incoming email, or by other services.

Currently, all Automation types support all triggers, with the following exceptions: `Create Folder` and  `Run Remote Server Sync` are not supported by the `action` trigger.

Automations can be triggered manually if trigger is not set to `action`.

## Destinations

The `destinations` parameter is a list of paths where files will be copied, moved, or created. It may include formatting parameters to dynamically determine the destination at runtime.

### Relative vs. Absolute Paths

In order to specify a relative path, it must start with either `./` or `../`. All other paths are considered absolute. In general, leading slashes should never be used on Files.com paths, including here. Paths are interpreted as absolute in all contexts, even without a leading slash.

### Files vs. Folders

If the destination path ends with a `/`, the filename from the source path will be preserved and put into the folder of this name.  If the destination path does not end with a `/`, it will be interpreted as a filename and will override the source file's filename entirely.

### Formatting Parameters

**Action-Triggered Automations**

* `%tf` : The name of the file that triggered the automation.
* `%tp` : The path of the file that triggered the automation.
* `%td` : The directory of the file that triggered the automation.
* `%tb` : The name of the file (without extension) that triggered the automation.
* `%te` : The extension of the file that triggered the automation.

For example, if the triggering file is at `path/to/file.txt`, then the automation destination `path/to/dest/incoming-%tf` will result in the actual destination being `path/to/dest/incoming-file.txt`.

**Parent Folders**

To reference the parent folder of a source file, use `%p1`, `%p2`, `%p3`, etc. for the first, second, third, etc. parent folder, respectively.

To reference path components from the root downward, use `%P1`, `%P2`, `%P3`, etc. for the first, second, third, etc. path component, respectively.

For example, if the source file is at `accounts/file.txt`, then the automation destination `path/to/dest/%p1/some_file_name.txt` will result in the actual destination being `path/to/dest/accounts/some_file_name.txt`.

If the source file is at `partner/app/team/inbound/file.txt`, then the automation destination `path/to/dest/%P1/%P3/%P4/file.txt` will result in the actual destination being `path/to/dest/partner/team/inbound/file.txt`.

**Source File Name**

To reference the name of the source file being processed, use the following tokens:

* `%Ff` : The name of the source file, with extension.
* `%Fb` : The name of the source file, without extension.
* `%Fe` : The extension of the source file.
* `%Fl` : The name of the source file, with extension, converted to lowercase.
* `%Fn` : The name of the source file, without non-alphanumeric characters, with extension.
* `%Fp` : The name of the source file, with extension, spaces removed, lowercase, non-ASCII normalized.

For example, if the source file is `Daily Report.xlsx` and the destination is `archive/%Y-%m-%d/%Fb.xlsx`, the resolved destination will be `archive/2024-01-15/Daily Report.xlsx`.

**Dates and Times**

* `%Y`  : The current year (4 digits)
* `%m`  : The current month (2 digits)
* `%B`  : The current month (full name)
* `%d`  : The current day (2 digits)
* `%H`  : The current hour (2 digits, 24-hour clock)
* `%M`  : The current minute (2 digits)
* `%S`  : The current second (2 digits)
* `%z`  : UTC Time Zone (e.g. -0900)

For example, if the current date is June 23, 2023 and the source file is named `daily_sales.csv`, then the following automation destination `path/to/dest/%Y/%m/%d/` will result in the actual destination being `path/to/dest/2023/06/23/daily_sales.csv`.

### Replacing Text

To replace text in the source filename, use the `destination_replace_from` and `destination_replace_to` parameters. This will perform a simple text replacement on the source filename before inserting it into the destination path.

For example, if the `destination_replace_from` is `incoming` and the `destination_replace_to` is `outgoing`, then `path/to/incoming.txt` will translate to `path/to/outgoing.txt`.


## Automation Types

There are several types of automations:  Create Folder, Copy File, Move File, Delete File and, Run Remote Server Sync.


### Create Folder

Creates the folder with named by `destinations` in the path named by `path`.
Destination may include formatting parameters to insert the date/time into the destination name.

Example Use case:  Our business files sales tax for each division in 11 states every quarter.
I want to create the folders where those sales tax forms and data will be collected.

I could create a Create Folder automation as follows:

* Trigger: `daily`
* Interval: `quarter_end`
* Path: `AccountingAndTax/SalesTax/State/*/`
* Destinations: `%Y/Quarter-ending-%m-%d`

Note this assumes you have folders in `AccountingAndTax/SalesTax/State/` already created for each state, e.g. `AccountingAndTax/SalesTax/State/CA/`.


### Delete File

Deletes the file with path matching `source` (wildcards allowed) in the path named by `path`.


### Copy File

Copies files in the folder named by `path` to the path specified in `destinations`.
The automation will only fire on files matching the `source` (wildcards allowed). In the case of an action-triggered automation, it will only operate on the actual file that triggered the automation.
If the parameter `limit` exists, the automation will only copy the newest `limit` files in each matching folder.


### Move File

Moves files in the folder named by `path` to the path specified in `destinations`.
The automation will only fire on files matching the `source` (wildcards allowed). In the case of an action-triggered automation, it will only operate on the actual file that triggered the automation.
If the parameter `limit` exists, the automation will only move the newest `limit` files in each matching folder.
Note that for a move with multiple destinations, all but one destination is treated as a copy.


### Run Remote Server Sync

The Run Remote Server Sync automation runs the remote server syncs specified by the `sync_ids`.

Typically when this automation is used, the remote server syncs in question are set to the manual
scheduling mode (`manual` to `true` via the API) to disable the built in sync scheduler.


### Import File

Retrieves files from one or more URLs and saves the results under the path specified in `destinations`.

The URLs to retrieve are specified as a JSON array in the `import_urls` property.

```json
[
  {
    "name": "response.json",
    "url": "https://example.com/api",
    "method": "post",
    "headers": {
      "Content-Type": "application/json"
    },
    "content": { "trigger-file": "%tp" }
  }
]
```

The recognized keys are:

* `name`: The file name which will be used to save the returned content. Required. `%` tokens will be replaced as described under Formatting Parameters.
* `url`: The URL which will be requested. Required.
* `method`: The HTTP method to be used for the request. May be either `get` or `post` (case insensitive). Defaults to `get`.
* `headers`: Optional headers to be included in the request. `%` tokens in the values will be replaced as described under Formatting Parameters.
* `content`: Optional body to send for POST request. If supplied as a string, `%` tokens will be expanded. If supplied as a JSON Object, `%` tokens will be expanded for top-level values. Other JSON types will be sent as-is.


### Help us build the future of Automations

Do you have an idea for something that would work well as a Files.com Automation?  Let us know!
We are actively improving the types of automations offered on our platform.


## Retrying Failures

Automations will automatically retry individual action steps up to 3 times, with pauses between retries that increase from 15 seconds to 1 minute.  If individual action steps fail after our 3rd attempt, that action will fail.  If every action step in an Automation Run fails, that automation run will move to a `failure` status.  If at least one step succeeds and one step fails, that automation run will move to a `partial_failure` status.

Automation Runs can be retried automatically when they enter a `failure` or `partial_failure` status as described above.  A retry will re-run the automation from scratch, including the "planning" phase, which expands globs (wildcards) and identifies which files to transfer or skip.

Retrying of entire Automation Runs must be explicitly enabled by setting the `retry_on_failure_interval_in_minutes` and `retry_on_failure_number_of_attempts` values on the Automation.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli automations list`

List Automations.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name`, `automation`, `last_modified_at` or `disabled`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `disabled`, `last_modified_at`, `workspace_id` or `automation`. Valid field combinations are `[ disabled, last_modified_at ]`, `[ workspace_id, disabled ]`, `[ disabled, automation ]`, `[ workspace_id, last_modified_at ]`, `[ automation, last_modified_at ]`, `[ workspace_id, automation ]`, `[ workspace_id, disabled, last_modified_at ]`, `[ disabled, automation, last_modified_at ]`, `[ workspace_id, disabled, automation ]`, `[ workspace_id, automation, last_modified_at ]` or `[ workspace_id, disabled, automation, last_modified_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `last_modified_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `last_modified_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `last_modified_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `last_modified_at`. |

### `files-cli automations find`

Show Automation.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Automation ID. **Required.** |

### `files-cli automations create [path]`

Create Automation.

| Flag | Type | Description |
| --- | --- | --- |
| `--source` | string | Source path/glob.  See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts. |
| `--destinations` | []string | A list of destination paths. Use a trailing slash for folder destinations and omit it for file destinations. |
| `--destination-replace-from` | string | If set, this string in the destination path will be replaced with the value in `destination_replace_to`. |
| `--destination-replace-to` | string | If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here. |
| `--interval` | string | How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `--path` | string | Path on which this Automation runs.  Supports globs, except on remote mounts. |
| `--legacy-sync-ids` | string | A list of legacy sync IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `--sync-ids` | string | A list of sync IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `--user-ids` | string | A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `--group-ids` | string | A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`. A list of days of the week to run this automation. 0 is Sunday, 1 is Monday, etc. |
| `--schedule-times-of-day` | []string | Times of day to run in HH:MM format (24-hour). Required for `custom_schedule` triggers. Optional for `daily` triggers - if not set, runs at midnight UTC. |
| `--schedule-time-zone` | string | Time zone for scheduled times. Optional for both `custom_schedule` and `daily` triggers. If not set, times are interpreted as UTC. |
| `--holiday-region` | string | Skip automation on holidays in this region. Optional for both `custom_schedule` and `daily` triggers. |
| `--always-overwrite-size-matching-files` | bool | Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer.  If this flag is `true` we will overwrite the destination file always.  Note that this may cause large amounts of wasted transfer usage.  This setting has no effect unless `overwrite_files` is also set to `true`. |
| `--always-serialize-jobs` | bool | Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance. |
| `--description` | string | Description for the this Automation. |
| `--disabled` | bool | If true, this automation will not run. |
| `--exclude-pattern` | string | If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts. |
| `--import-urls` | []object | List of URLs to be imported and names to be used. |
| `--flatten-destination-structure` | bool | Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination.  If this flag is `true`, the source folder structure will be flattened in the destination.  This is useful for copying or moving files from multiple folders into a single destination folder. |
| `--ignore-locked-folders` | bool | If true, the Lock Folders behavior will be disregarded for automated actions. |
| `--legacy-folder-matching` | bool | DEPRECATED: If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files.  This behavior no longer works and should not be used. |
| `--name` | string | Name for this automation. |
| `--overwrite-files` | bool | If true, existing files will be overwritten with new files on Move/Copy automations.  Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file.  Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what. |
| `--path-time-zone` | string | Timezone to use when rendering timestamps in paths. |
| `--retry-on-failure-interval-in-minutes` | int64 | If the Automation fails, retry at this interval (in minutes).  Acceptable values are 5 through 1440 (one day).  Set to null to disable. |
| `--retry-on-failure-number-of-attempts` | int64 | If the Automation fails, retry at most this many times.  Maximum allowed value: 10.  Set to null to disable. |
| `--trigger` | enum | How this automation is triggered to run. One of: `manual`, `daily`, `custom_schedule`, `webhook`, `email`, `action`. |
| `--trigger-actions` | []string | If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy |
| `--value` | object | A Hash of attributes specific to the automation type. |
| `--recurring-day` | int64 | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `--automation` | enum | Automation type. One of: `create_folder`, `delete_file`, `copy_file`, `move_file`, `as2_send`, `run_sync`, `import_file`, `v2`. **Required.** |
| `--workspace-id` | int64 | Workspace ID |

### `files-cli automations manual-run`

Manually Run Automation.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Automation ID. **Required.** |

### `files-cli automations update [path]`

Update Automation.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Automation ID. **Required.** |
| `--source` | string | Source path/glob.  See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts. |
| `--destinations` | []string | A list of destination paths. Use a trailing slash for folder destinations and omit it for file destinations. |
| `--destination-replace-from` | string | If set, this string in the destination path will be replaced with the value in `destination_replace_to`. |
| `--destination-replace-to` | string | If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here. |
| `--interval` | string | How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end` |
| `--path` | string | Path on which this Automation runs.  Supports globs, except on remote mounts. |
| `--legacy-sync-ids` | string | A list of legacy sync IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `--sync-ids` | string | A list of sync IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `--user-ids` | string | A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `--group-ids` | string | A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited. |
| `--schedule-days-of-week` | []int64 | If trigger is `custom_schedule`. A list of days of the week to run this automation. 0 is Sunday, 1 is Monday, etc. |
| `--schedule-times-of-day` | []string | Times of day to run in HH:MM format (24-hour). Required for `custom_schedule` triggers. Optional for `daily` triggers - if not set, runs at midnight UTC. |
| `--schedule-time-zone` | string | Time zone for scheduled times. Optional for both `custom_schedule` and `daily` triggers. If not set, times are interpreted as UTC. |
| `--holiday-region` | string | Skip automation on holidays in this region. Optional for both `custom_schedule` and `daily` triggers. |
| `--always-overwrite-size-matching-files` | bool | Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer.  If this flag is `true` we will overwrite the destination file always.  Note that this may cause large amounts of wasted transfer usage.  This setting has no effect unless `overwrite_files` is also set to `true`. |
| `--always-serialize-jobs` | bool | Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance. |
| `--description` | string | Description for the this Automation. |
| `--disabled` | bool | If true, this automation will not run. |
| `--exclude-pattern` | string | If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts. |
| `--import-urls` | []object | List of URLs to be imported and names to be used. |
| `--flatten-destination-structure` | bool | Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination.  If this flag is `true`, the source folder structure will be flattened in the destination.  This is useful for copying or moving files from multiple folders into a single destination folder. |
| `--ignore-locked-folders` | bool | If true, the Lock Folders behavior will be disregarded for automated actions. |
| `--legacy-folder-matching` | bool | DEPRECATED: If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files.  This behavior no longer works and should not be used. |
| `--name` | string | Name for this automation. |
| `--overwrite-files` | bool | If true, existing files will be overwritten with new files on Move/Copy automations.  Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file.  Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what. |
| `--path-time-zone` | string | Timezone to use when rendering timestamps in paths. |
| `--retry-on-failure-interval-in-minutes` | int64 | If the Automation fails, retry at this interval (in minutes).  Acceptable values are 5 through 1440 (one day).  Set to null to disable. |
| `--retry-on-failure-number-of-attempts` | int64 | If the Automation fails, retry at most this many times.  Maximum allowed value: 10.  Set to null to disable. |
| `--trigger` | enum | How this automation is triggered to run. One of: `manual`, `daily`, `custom_schedule`, `webhook`, `email`, `action`. |
| `--trigger-actions` | []string | If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, copy, move, archived_delete, update, read, destroy |
| `--value` | object | A Hash of attributes specific to the automation type. |
| `--recurring-day` | int64 | If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`. |
| `--automation` | enum | Automation type. One of: `create_folder`, `delete_file`, `copy_file`, `move_file`, `as2_send`, `run_sync`, `import_file`, `v2`. |

### `files-cli automations delete`

Delete Automation.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Automation ID. **Required.** |

