---
name: filescom-message-comments
description: |
  A MessageComment is a comment made by a user on a message.
---

# filescom-message-comments

A MessageComment is a comment made by a user on a message.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli message-comments list`

List Message Comments.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `--message-id` | int64 | Message comment to return comments for. **Required.** |

### `files-cli message-comments find`

Show Message Comment.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message Comment ID. **Required.** |

### `files-cli message-comments create`

Create Message Comment.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--body` | string | Comment body. **Required.** |

### `files-cli message-comments update`

Update Message Comment.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message Comment ID. **Required.** |
| `--body` | string | Comment body. **Required.** |

### `files-cli message-comments delete`

Delete Message Comment.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message Comment ID. **Required.** |

