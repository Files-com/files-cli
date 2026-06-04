---
name: filescom-api-keys
description: |
  An APIKey is a key that allows programmatic access to your Site.
---

# filescom-api-keys

An APIKey is a key that allows programmatic access to your Site.

API keys confer all the permissions of the user who owns them.
If an API key is created without a user owner, it is considered a site-wide API key, which has full permissions to do anything on the Site.

We recommend registering API keys to service users wherever possible and then using User or Group Permissions to restrict that API Key appropriately.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli api-keys list`

List API Keys.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`. |
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
| `--path` | string | Folder path restriction for `office_integration` permission set API keys. |
| `--permission-set` | enum | Permissions for this API Key. It must be full for site-wide API Keys.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. One of: `none`, `full`, `desktop_app`, `sync_app`, `office_integration`, `mobile_app`. |

### `files-cli api-keys update-current`

Update current API key.  (Requires current API connection to be using an API key.).

| Flag | Type | Description |
| --- | --- | --- |
| `--expires-at` | datetime | API Key expiration date |
| `--name` | string | Internal name for the API Key.  For your use. |
| `--permission-set` | enum | Permissions for this API Key. It must be full for site-wide API Keys.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. One of: `none`, `full`, `desktop_app`, `sync_app`, `office_integration`, `mobile_app`. |

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

