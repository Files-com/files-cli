---
name: filescom-restores
description: |
  A Restore kicks off a process to restore deleted data for your Site.
---

# filescom-restores

A Restore kicks off a process to restore deleted data for your Site. We are only able to restore deleted items up less than 365 days old. This is available to Site Administrators and to Workspace Administrators when the restore is scoped to their workspace.

This process does not validate the existence of prior files and folders at the creation time of the Restore. If a file or directory is incorrect, does not exist, or is older than 365 days, the restore process will still show to have completed successfully even if no items were restored.

While regular expressions and user-supplied wildcards are not supported, the system automatically applies a wildcard at the end of the prefix and is case-insensitive.
Example: A prefix of `test` will match `test1/`, `test2/`, `testing.txt`, `Testfile1.mp4`, `Testing.pdf`, etc.

## Restore types

Restore supports multiple restoration types, controlled by the `restoration_type` field:
* `files` (default): Restore deleted files/folders (and optionally file permissions) by path prefix.
* `users`: Restore deleted users (and optionally user permissions) by username prefix.

## Restoring deleted files/folders (`restoration_type=files`)
* To restore a specific file, specify the path to the file in the prefix field. `Example: path/to/my/deleted_file.txt`
* To restore a directory, specify the directory path ending with an `/` in the prefix field. Example: `path/to/my/deleted_directory/`
* To restore all deleted items, specify an empty string (`''`) in the prefix field or omit the field from the request.

## Restoring deleted users (`restoration_type=users`)
* Restore all deleted users since `earliest_date` by omitting prefix (or using `''`).
* To restore specific deleted users by username, use `prefix` as a case-insensitive username prefix.
  Example: A prefix of `john` will match `john`, `johnny`, `John.Doe`, etc.
* When restoring users, we also restore associated authentication and access records deleted since `earliest_date`, including:
  - Permissions (when `restore_deleted_permissions=true`)
  - Two-factor authentication methods
  - Public keys
  - API keys

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli restores list`

List Restores.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `restoration_type`. |

### `files-cli restores create`

Create Restore.

| Flag | Type | Description |
| --- | --- | --- |
| `--earliest-date` | datetime | Restore all files deleted after this date/time. Don't set this earlier than you need. Can not be greater than 365 days prior to the restore request. **Required.** |
| `--prefix` | string | Prefix of the files/folders to restore. To restore a folder, add a trailing slash to the folder name. Do not use a leading slash. To restore all deleted items, specify an empty string (`''`) in the prefix field or omit the field from the request. |
| `--restoration-type` | enum | Type of restoration to perform. `files` restores deleted filesystem items. `users` restores deleted users and associated access/authentication records. One of: `files`, `users`. |
| `--restore-deleted-permissions` | bool | If true, we will also restore any Permissions that match the same path prefix from the same dates. |
| `--restore-in-place` | bool | If true, we will restore the files in place (into their original paths). If false, we will create a new restoration folder in the root and restore files there. |
| `--update-timestamps` | bool | If true, we will update the last modified timestamp of restored files to today's date. If false, we might trigger File Expiration to delete the file again. |
| `--workspace-id` | int64 | Workspace ID for a workspace-scoped restore. `0` means the default site-wide scope. |

