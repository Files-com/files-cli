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

No flags beyond the global ones.

### `files-cli chat-sessions find`

Show Chat Session.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Chat Session ID. **Required.** |

