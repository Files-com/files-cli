---
name: filescom-behaviors
description: |
  A Behavior is an API resource for what are also known as Folder Settings.
---

# filescom-behaviors

A Behavior is an API resource for what are also known as Folder Settings.  Every behavior is associated with a folder.

Depending on the behavior, it may also operate on child folders.  It may be overridable at the child folder level or maybe can be added to at the child folder level.  The exact options for each behavior type are explained in the table below.

Each behavior type also has a recursion mode in the behavior type documentation. `always` means the behavior is always recursive, `never` means it is never recursive, and `sometimes` means callers may choose the value of the `recursive` field.

Additionally, some behaviors are visible to non-admins, and others are even settable by non-admins.  All the details are below.

Each behavior uses a different format for storing its settings value.  Next to each behavior type is an example value.  Our API and SDKs currently require that the value for behaviors be sent as raw JSON within the `value` field.  Our SDK generator and API documentation generator doesn't fully keep up with this requirement, so if you need any help finding the exact syntax to use for your language or use case, just reach out.

Note: Append Timestamp behavior removed. Check [Override Upload Filename](#override-upload-filename-behaviors) behavior which have even more functionality to modify name on upload.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli behaviors list`

List Behaviors.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `behavior`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `clickwrap_id`, `form_field_set_id`, `impacts_ui`, `remote_server_id` or `behavior`. Valid field combinations are `[ impacts_ui, behavior ]`. |

### `files-cli behaviors find`

Show Behavior.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Behavior ID. **Required.** |

### `files-cli behaviors list-for [path]`

List Behaviors by Path.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `behavior`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `impacts_ui` and `behavior`. Valid field combinations are `[ impacts_ui, behavior ]`. |
| `--path` | string | Path to operate on. **Required.** |
| `--ancestor-behaviors` | bool | If `true`, behaviors above this path are shown. |

### `files-cli behaviors create [path]`

Create Behavior.

| Flag | Type | Description |
| --- | --- | --- |
| `--value` | object | This field stores a hash of data specific to the type of behavior. See The Behavior Types section for example values for each type of behavior. |
| `--attachment-file` | file | Certain behaviors may require a file, for instance, the `watermark` behavior requires a watermark image. Attach that file here. |
| `--disable-parent-folder-behavior` | bool | If `true`, the parent folder's behavior will be disabled for this folder and its children. This is the main mechanism for canceling out a `recursive` behavior higher in the folder tree. |
| `--recursive` | bool | Whether the behavior should apply to child folders. This is only configurable for behavior types whose recursion mode is `sometimes`; `always` behaviors stay recursive and `never` behaviors stay non-recursive. |
| `--name` | string | Name for this behavior. |
| `--description` | string | Description for this behavior. |
| `--path` | string | Path where this behavior should apply. **Required.** |
| `--behavior` | string | Behavior type. **Required.** |

### `files-cli behaviors webhook-test`

Test Webhook.

| Flag | Type | Description |
| --- | --- | --- |
| `--url` | string | URL for testing the webhook. **Required.** |
| `--method` | string | HTTP request method (GET or POST). |
| `--encoding` | string | Encoding type for the webhook payload. Can be JSON, XML, or RAW (form data). |
| `--headers` | object | Additional request headers to send via HTTP. |
| `--body` | object | Additional body parameters to include in the webhook payload. |

### `files-cli behaviors update`

Update Behavior.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Behavior ID. **Required.** |
| `--value` | object | This field stores a hash of data specific to the type of behavior. See The Behavior Types section for example values for each type of behavior. |
| `--attachment-file` | file | Certain behaviors may require a file, for instance, the `watermark` behavior requires a watermark image. Attach that file here. |
| `--disable-parent-folder-behavior` | bool | If `true`, the parent folder's behavior will be disabled for this folder and its children. This is the main mechanism for canceling out a `recursive` behavior higher in the folder tree. |
| `--recursive` | bool | Whether the behavior should apply to child folders. This is only configurable for behavior types whose recursion mode is `sometimes`; `always` behaviors stay recursive and `never` behaviors stay non-recursive. |
| `--name` | string | Name for this behavior. |
| `--description` | string | Description for this behavior. |
| `--attachment-delete` | bool | If `true`, delete the file stored in `attachment`. |

### `files-cli behaviors delete`

Delete Behavior.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Behavior ID. **Required.** |

