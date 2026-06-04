---
name: filescom-messages
description: |
  A Message is a part of Files.com's project management features and represents a message posted by a user to a project.
---

# filescom-messages

A Message is a part of Files.com's project management features and represents a message posted by a user to a project.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli messages list`

List Messages.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--project-id` | int64 | Project for which to return messages. **Required.** |

### `files-cli messages find`

Show Message.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message ID. **Required.** |

### `files-cli messages create`

Create Message.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--project-id` | int64 | Project to which the message should be attached. **Required.** |
| `--subject` | string | Message subject. **Required.** |
| `--body` | string | Message body. **Required.** |

### `files-cli messages update`

Update Message.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message ID. **Required.** |
| `--project-id` | int64 | Project to which the message should be attached. **Required.** |
| `--subject` | string | Message subject. **Required.** |
| `--body` | string | Message body. **Required.** |

### `files-cli messages delete`

Delete Message.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message ID. **Required.** |

