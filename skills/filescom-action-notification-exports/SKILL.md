---
name: filescom-action-notification-exports
description: |
  An ActionNotificationExport is an operation that provides access to outgoing webhook logs.
---

# filescom-action-notification-exports

An ActionNotificationExport is an operation that provides access to outgoing webhook logs. Querying webhook logs is a little different than other APIs.

All queries against the archive must be submitted as Exports.  (Even our Web UI creates an Export behind the scenes.)

In any query field in this API, you may specify multiple values separated by commas.  That means that commas
cannot be searched for themselves, and neither can single quotation marks.

Use the following steps to complete an export:

1. Initiate the export by using the Create Action Notification Export endpoint. Non Site Admins must query by folder or path.
2. Using the `id` from the response to step 1, poll the Show Action Notification Export endpoint. Check the `status` field until it is `ready`.
3. You can download the results of the export as a CSV file using the `results_url` field in the response from step 2. If you want to page through the records in JSON format, use the List Action Notification Export Results endpoint, passing the `id` that you got in step 1 as the `action_notification_export_id` parameter. Check the `X-Files-Cursor-Next` header to see if there are more records available, and resubmit the same request with a `cursor` parameter to fetch the next page of results.  Unlike most API Endpoints, this endpoint does not provide `X-Files-Cursor-Prev` cursors allowing reverse pagination through the results.  This is due to limitations in Amazon Athena, the underlying data lake for these records.

If you intend to use this API for high volume or automated use, please contact us with more information
about your use case.

## Example Queries

* History for a folder: `{ "query_folder": "path/to/folder" }`
* History for a range of time: `{ "start_at": "2021-03-18 12:00:00", "end_at": "2021-03-19 12:00:00" }`
* History of all notifications that used GET or POST: `{ "query_request_method": "GET,POST" }`

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli action-notification-exports find`

Show Action Notification Export.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Action Notification Export ID. **Required.** |

### `files-cli action-notification-exports create`

Create Action Notification Export.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--start-at` | datetime | Start date/time of export range. |
| `--end-at` | datetime | End date/time of export range. |
| `--query-message` | string | Error message associated with the request, if any. |
| `--query-request-method` | string | The HTTP request method used by the webhook. |
| `--query-request-url` | string | The target webhook URL. |
| `--query-status` | string | The HTTP status returned from the server in response to the webhook request. |
| `--query-success` | bool | true if the webhook request succeeded (i.e. returned a 200 or 204 response status). false otherwise. |
| `--query-path` | string | Return notifications that were triggered by actions on this specific path. |
| `--query-folder` | string | Return notifications that were triggered by actions in this folder. |

