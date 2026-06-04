---
name: filescom-metadata-categories
description: |
  A MetadataCategory defines a reusable set of Custom Metadata rules that can be assigned to folders via a folder behavior.
---

# filescom-metadata-categories

A MetadataCategory defines a reusable set of Custom Metadata rules that can be assigned to folders
via a folder behavior. Each category specifies named metadata keys with optional allowed-value
constraints, and a set of default columns to display in the UI.

If a key's `allowed_values` array is empty, it is treated as a free-form text field.
If the array is non-empty, the key is constrained to those values in the Web UI.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli metadata-categories list`

List Metadata Categories.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |

### `files-cli metadata-categories find`

Show Metadata Category.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Metadata Category ID. **Required.** |

### `files-cli metadata-categories list-for [path]`

List Metadata Categories by Path.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |

### `files-cli metadata-categories create`

Create Metadata Category.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Name of the metadata category. **Required.** |
| `--default-columns` | []string | Metadata keys that should appear as columns in the UI by default. |

### `files-cli metadata-categories update`

Update Metadata Category.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Metadata Category ID. **Required.** |
| `--name` | string | Name of the metadata category. |
| `--default-columns` | []string | Metadata keys that should appear as columns in the UI by default. |

### `files-cli metadata-categories delete`

Delete Metadata Category.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Metadata Category ID. **Required.** |

