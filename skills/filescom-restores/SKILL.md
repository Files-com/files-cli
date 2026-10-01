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
  This restores that folder and its contents without matching sibling names that begin with the same text.
* To restore all deleted items, specify an empty string (`''`) in the prefix field or omit the field from the request. A workspace-scoped restore includes only that workspace.
* With `restore_deleted_permissions=true`, an in-place restore also restores user and group permissions on folders returned by the restore and their subfolders, when those permissions were deleted within 24 hours after the folder's deletion, including both endpoints. Permissions removed before the folder was deleted or more than 24 hours afterward remain deleted. Permissions removed independently during that window are included.
* If a folder already exists at the original path, restored content goes into it and eligible permissions apply to it. Permissions for deleted users or groups remain deleted. Restoring into a new restoration folder does not restore permissions.

## Restoring deleted users (`restoration_type=users`)
* Use `earliest_date` to select users deleted on or after that date/time. Restore all matching users by omitting `prefix` (or using `''`).
* To restore specific deleted users by username, use `prefix` as a case-insensitive username prefix.
  Example: A prefix of `john` will match `john`, `johnny`, `John.Doe`, etc.
* For each restored user, we also restore associated authentication and access records removed as part of deleting that user:
  - Permissions (when `restore_deleted_permissions=true`)
  - Two-factor authentication methods
  - SFTP/SSH keys
  - API keys
* Records deleted independently of the user remain deleted, even if they were deleted after `earliest_date`.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

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
| `--earliest-date` | datetime | Restore files or users deleted on or after this date/time. Don't set this earlier than you need. Can not be greater than 365 days prior to the restore request. **Required.** |
| `--prefix` | string | Prefix of the files/folders to restore, or a case-insensitive username prefix for a user restore. A trailing slash selects that folder and its contents only; without it, the prefix matches any path beginning with that text. Do not use a leading slash. To restore all deleted items of the selected restoration type within the selected site or workspace scope, specify an empty string (`''`) or omit the field. |
| `--restoration-type` | enum | Type of restoration to perform. `files` restores deleted filesystem items. `users` restores deleted users and associated access/authentication records removed as part of deleting those users. One of: `files`, `users`. |
| `--restore-deleted-permissions` | bool | If true, a user restore restores permissions removed as part of deleting the selected users. An in-place file restore restores user and group permissions on returned folders and their subfolders that were deleted from the time of the folder's deletion through 24 hours afterward, including permissions removed independently within that window. File restores into a new folder do not restore permissions. |
| `--restore-in-place` | bool | If true, we will restore the files in place (into their original paths). If false, we will create a new restoration folder in the root and restore files there. |
| `--update-timestamps` | bool | If true, we will update the last modified timestamp of restored files to today's date. If false, we might trigger File Expiration to delete the file again. |
| `--workspace-id` | int64 | Workspace ID for a workspace-scoped restore. `0` means the default site-wide scope. |

