---
name: filescom-history-exports
description: |
  A History Export is a resource on the API that is used to export historical action (history) logs.
---

# filescom-history-exports

A History Export is a resource on the API that is used to export historical action (history) logs.

All queries against the archive must be submitted as Exports.  (Even our Web UI creates an Export behind
the scenes.)

We use Amazon Athena behind the scenes for processing these queries, and as such, have powerful
search capabilities.  We've done our best to expose search capabilities via this History Export API.

In any query field in this API, you may specify multiple values separated by commas.  That means that commas
cannot be searched for themselves, and neither can single quotation marks.

We do not currently partition data by date on the backend, so all queries result in a full scan of the entire
data lake.  This means that all queries will take about the same amount of time to complete, and we incur about
the same cost per query internally.  We don't typically bill our customers for these queries, assuming
usage is occasional and manual.

If you intend to use this API for high volume or automated use, please contact us with more information
about your use case.  We may decide to change the backend data schema to match your use case more closely, and
we may also need to charge an additional cost per query.

## Example History Queries

* History for a user: `{ "query_user_id": 123 }`
* History for a range of time: `{ "start_at": "2021-03-18 12:00:00", "end_at": "2021-03-19 12:00:00" }`
* History of logins and failed logins: `{ "query_action": "login,failedlogin" }`
* A Complex query: `{ "query_folder": "uploads", "query_action": "create,copy,move", "start_at": "2021-03-18 12:00:00", "end_at": "2021-03-19 12:00:00" }`

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli history-exports find`

Show History Export.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | History Export ID. **Required.** |

### `files-cli history-exports create`

Create History Export.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--start-at` | datetime | Start date/time of export range. |
| `--end-at` | datetime | End date/time of export range. |
| `--query-action` | string | Filter results by this this action type. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`, `archived_delete` |
| `--query-interface` | string | Filter results by this this interface type. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`, `inbound_email`, `remote`, `inbound_s3` |
| `--query-user-id` | string | Return results that are actions performed by the user indicated by this User ID |
| `--query-file-id` | string | Return results that are file actions related to the file indicated by this File ID |
| `--query-parent-id` | string | Return results that are file actions inside the parent folder specified by this folder ID |
| `--query-path` | string | Return results that are file actions related to paths matching this pattern. |
| `--query-folder` | string | Return results that are file actions related to files or folders inside folder paths matching this pattern. |
| `--query-src` | string | Return results that are file moves originating from paths matching this pattern. |
| `--query-destination` | string | Return results that are file moves with paths matching this pattern as destination. |
| `--query-ip` | string | Filter results by this IP address. |
| `--query-username` | string | Filter results by this username. |
| `--query-failure-type` | string | If searching for Histories about login failures, this parameter restricts results to failures of this specific type.  Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`, `insecure_ftp`, `insecure_cipher`, `rate_limited` |
| `--query-target-id` | string | If searching for Histories about specific objects (such as Users, or API Keys), this parameter restricts results to objects that match this ID. |
| `--query-target-name` | string | If searching for Histories about Users, Groups or other objects with names, this parameter restricts results to objects with this name/username. |
| `--query-target-permission` | string | If searching for Histories about Permissions, this parameter restricts results to permissions of this level. |
| `--query-target-user-id` | string | If searching for Histories about API keys, this parameter restricts results to API keys created by/for this user ID. |
| `--query-target-username` | string | If searching for Histories about API keys, this parameter restricts results to API keys created by/for this username. |
| `--query-target-platform` | string | If searching for Histories about API keys, this parameter restricts results to API keys associated with this platform. |
| `--query-target-permission-set` | string | If searching for Histories about API keys, this parameter restricts results to API keys with this permission set. |

