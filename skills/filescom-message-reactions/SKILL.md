---
name: filescom-message-reactions
description: |
  A MessageReaction is a reaction emoji made by a user on a message.
---

# filescom-message-reactions

A MessageReaction is a reaction emoji made by a user on a message.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli message-reactions list`

List Message Reactions.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--message-id` | int64 | Message to return reactions for. **Required.** |

### `files-cli message-reactions find`

Show Message Reaction.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message Reaction ID. **Required.** |

### `files-cli message-reactions create`

Create Message Reaction.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--emoji` | string | Emoji to react with. **Required.** |

### `files-cli message-reactions delete`

Delete Message Reaction.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Message Reaction ID. **Required.** |

