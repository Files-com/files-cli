---
name: filescom-chat-sessions
description: |
  A ChatSession represents one conversation with the Files.com AI Assistant.
---

# filescom-chat-sessions

A ChatSession represents one conversation with the Files.com AI Assistant.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli chat-sessions list`

List Chat Sessions.

| Flag | Type | Description |
| --- | --- | --- |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `ai_task_id`. |

### `files-cli chat-sessions find`

Show Chat Session.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | string | Chat Session ID. **Required.** |

