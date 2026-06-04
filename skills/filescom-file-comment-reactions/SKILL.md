---
name: filescom-file-comment-reactions
description: |
  A FileCommentReaction is a reaction that is attached to a comment on a file.
---

# filescom-file-comment-reactions

A FileCommentReaction is a reaction that is attached to a comment on a file.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli file-comment-reactions create`

Create File Comment Reaction.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--file-comment-id` | int64 | ID of file comment to attach reaction to. **Required.** |
| `--emoji` | string | Emoji to react with. **Required.** |

### `files-cli file-comment-reactions delete`

Delete File Comment Reaction.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | File Comment Reaction ID. **Required.** |

