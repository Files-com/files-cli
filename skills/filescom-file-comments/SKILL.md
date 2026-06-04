---
name: filescom-file-comments
description: |
  A FileComment is a comment attached to a file by a user.
---

# filescom-file-comments

A FileComment is a comment attached to a file by a user.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli file-comments list-for [path]`

List File Comments by Path.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |

### `files-cli file-comments create [path]`

Create File Comment.

| Flag | Type | Description |
| --- | --- | --- |
| `--body` | string | Comment body. **Required.** |
| `--path` | string | File path. **Required.** |

### `files-cli file-comments update`

Update File Comment.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | File Comment ID. **Required.** |
| `--body` | string | Comment body. **Required.** |

### `files-cli file-comments delete`

Delete File Comment.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | File Comment ID. **Required.** |

