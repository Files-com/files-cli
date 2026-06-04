---
name: filescom-histories
description: |
  An Action is a single record in our history log.
---

# filescom-histories

An Action is a single record in our history log. File and Login actions on Files.com are recorded and can be queried via History API.

This set of endpoints only provides access to recent actions (actions created within 24 hours).  In most cases,
you would be better served to create a Webhook, which sends actions to your server, rather than poll this endpoint.

The HistoryExport set of endpoints provides a richer ability to query and filter, as well as search the entire
lifetime of your history log.

Note: Failed Logins are no longer logged in this logging mechanism.
The `failedlogin` value is still listed in the `action` documentation for legacy reasons.
Use API or other protocol logs (SFTP, FTP, etc.) for detailed and better information about login failures.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli histories list-for-file [path]`

List history for specific file.

| Flag | Type | Description |
| --- | --- | --- |
| `--start-at` | datetime | Leave blank or set to a date/time to filter earlier entries. |
| `--end-at` | datetime | Leave blank or set to a date/time to filter later entries. |
| `--display` | string | Display format. Leave blank or set to `full` or `parent`. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `--path` | string | Path to operate on. **Required.** |

### `files-cli histories list-for-folder [path]`

List history for specific folder.

| Flag | Type | Description |
| --- | --- | --- |
| `--start-at` | datetime | Leave blank or set to a date/time to filter earlier entries. |
| `--end-at` | datetime | Leave blank or set to a date/time to filter later entries. |
| `--display` | string | Display format. Leave blank or set to `full` or `parent`. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `--path` | string | Path to operate on. **Required.** |

### `files-cli histories list-for-user`

List history for specific user.

| Flag | Type | Description |
| --- | --- | --- |
| `--start-at` | datetime | Leave blank or set to a date/time to filter earlier entries. |
| `--end-at` | datetime | Leave blank or set to a date/time to filter later entries. |
| `--display` | string | Display format. Leave blank or set to `full` or `parent`. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `--user-id` | int64 | User ID. **Required.** |

### `files-cli histories list-logins`

List site login history.

| Flag | Type | Description |
| --- | --- | --- |
| `--start-at` | datetime | Leave blank or set to a date/time to filter earlier entries. |
| `--end-at` | datetime | Leave blank or set to a date/time to filter later entries. |
| `--display` | string | Display format. Leave blank or set to `full` or `parent`. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |

### `files-cli histories list`

List site full action history.

| Flag | Type | Description |
| --- | --- | --- |
| `--start-at` | datetime | Leave blank or set to a date/time to filter earlier entries. |
| `--end-at` | datetime | Leave blank or set to a date/time to filter later entries. |
| `--display` | string | Display format. Leave blank or set to `full` or `parent`. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `user_id`, `folder` or `path`. Valid field combinations are `[  ]`, `[ path ]`, `[ path ]` or `[ path ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. |

