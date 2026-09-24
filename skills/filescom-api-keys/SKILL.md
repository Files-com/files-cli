---
name: filescom-api-keys
description: |
  An APIKey is a key that allows programmatic access to your Site.
---

# filescom-api-keys

An APIKey is a key that allows programmatic access to your Site.

API keys use the owning user's permissions, narrowed by the key's permission set, workspace scope, and any folder path restriction.
If an API key is created without a user owner, it is considered a site-wide API key. Site-wide API keys with the `files_only` permission set are restricted to file-user permissions and workspace scoping.

Set `path` when creating a key to limit file and folder access to that folder and its descendants. Except for `office_integration` keys, the path does not need to exist when the key is created. The restriction applies to the files, folders, and file_actions endpoints, on every path those requests access, including both source and destination paths for copy and move operations. Requests to those endpoints outside the restriction are denied with `not-authorized/api-key-is-path-restricted`. It never grants additional access to the owning user. A key with the `files_only` permission set can use only those endpoints, plus `GET /file_migrations/{id}` to follow a copy, move, zip, or unzip it started and `GET /api_key` to read its own record, so a `files_only` key with a `path` is confined to that folder; every other endpoint denies a `files_only` key with `not-authorized/api-key-only-for-file-operations`.

We recommend registering API keys to service users wherever possible and then using User or Group Permissions to restrict that API Key appropriately.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli api-keys list`

List API Keys.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id` and `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `aws_style_credentials` and `expires_at`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `expires_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `expires_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `expires_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `expires_at`. |

### `files-cli api-keys find-current`

Show information about current API key.  (Requires current API connection to be using an API key.).

No flags beyond the global ones.

### `files-cli api-keys find`

Show API Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Api Key ID. **Required.** |

### `files-cli api-keys create [path]`

Create API Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--description` | string | User-supplied description of API key. |
| `--expires-at` | datetime | API Key expiration date |
| `--name` | string | Internal name for the API Key.  For your use. **Required.** |
| `--aws-style-credentials` | bool | If `true`, this API key will be usable with AWS-compatible endpoints, such as our Inbound S3-compatible endpoint. |
| `--path` | string | Restricts the file and folder operations made with this key, meaning the files, folders, and file_actions endpoints, to the specified folder and its descendants, including copy and move destinations. Other endpoints do not apply the path restriction; use the `files_only` permission set to confine a key to the endpoints that do. Does not grant access beyond the owning user's permissions. Optional except for `office_integration` keys, which require a path the owning user can read. |
| `--permission-set` | enum | Permissions for this API Key. Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Keys with the `files_only` permission set can use only the files, folders, and file_actions endpoints, where they perform file operations as a full-access file user in the key's workspace scope, along with `GET /file_migrations/{id}` and `GET /api_key`. They cannot use site admin, workspace admin, folder admin, group admin, partner admin, or billing privileges from the owning user, and every other endpoint denies them with `not-authorized/api-key-only-for-file-operations`. One of: `none`, `full`, `desktop_app`, `sync_app`, `office_integration`, `mobile_app`, `files_only`. |
| `--workspace-id` | int64 | Workspace ID for this API Key. `0` means the default workspace. |

### `files-cli api-keys update-current`

Update current API key.  (Requires current API connection to be using an API key.).

| Flag | Type | Description |
| --- | --- | --- |
| `--expires-at` | datetime | API Key expiration date |
| `--name` | string | Internal name for the API Key.  For your use. |
| `--permission-set` | enum | Permissions for this API Key. Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Keys with the `files_only` permission set can use only the files, folders, and file_actions endpoints, where they perform file operations as a full-access file user in the key's workspace scope, along with `GET /file_migrations/{id}` and `GET /api_key`. They cannot use site admin, workspace admin, folder admin, group admin, partner admin, or billing privileges from the owning user, and every other endpoint denies them with `not-authorized/api-key-only-for-file-operations`. One of: `none`, `full`, `desktop_app`, `sync_app`, `office_integration`, `mobile_app`, `files_only`. |

### `files-cli api-keys update`

Update API Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Api Key ID. **Required.** |
| `--description` | string | User-supplied description of API key. |
| `--expires-at` | datetime | API Key expiration date |
| `--name` | string | Internal name for the API Key.  For your use. |

### `files-cli api-keys delete-current`

Delete current API key.  (Requires current API connection to be using an API key.).

No flags beyond the global ones.

### `files-cli api-keys delete`

Delete API Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Api Key ID. **Required.** |

