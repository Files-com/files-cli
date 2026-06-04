---
name: filescom-message-comment-reactions
description: |
  A MessageCommentReaction is a reaction emoji made by a user on a message comment.
---

# filescom-message-comment-reactions

A MessageCommentReaction is a reaction emoji made by a user on a message comment.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli message-comment-reactions list`

List Message Comment Reactions.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `--message-comment-id` | int64 | Message comment to return reactions for. **Required.** |

### `files-cli message-comment-reactions find`

Show Message Comment Reaction.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message Comment Reaction ID. **Required.** |

### `files-cli message-comment-reactions create`

Create Message Comment Reaction.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--emoji` | string | Emoji to react with. **Required.** |

### `files-cli message-comment-reactions delete`

Delete Message Comment Reaction.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message Comment Reaction ID. **Required.** |

