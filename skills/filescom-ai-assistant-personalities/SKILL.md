---
name: filescom-ai-assistant-personalities
description: |
  An AI Assistant Personality defines a system prompt used to customize the in-app AI Assistant.
---

# filescom-ai-assistant-personalities

An AI Assistant Personality defines a system prompt used to customize the in-app AI Assistant.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli ai-assistant-personalities list`

List Ai Assistant Personalities.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |

### `files-cli ai-assistant-personalities find`

Show Ai Assistant Personality.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Ai Assistant Personality ID. **Required.** |

### `files-cli ai-assistant-personalities create`

Create Ai Assistant Personality.

| Flag | Type | Description |
| --- | --- | --- |
| `--apply-to-all-workspaces` | bool | If true, this default-workspace personality can apply to users in all workspaces. |
| `--name` | string | AI Assistant Personality name. **Required.** |
| `--system-prompt` | string | System prompt injected into the in-app AI Assistant. **Required.** |
| `--use-by-default` | bool | Whether this personality is the default personality for the Workspace. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli ai-assistant-personalities update`

Update Ai Assistant Personality.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Ai Assistant Personality ID. **Required.** |
| `--apply-to-all-workspaces` | bool | If true, this default-workspace personality can apply to users in all workspaces. |
| `--name` | string | AI Assistant Personality name. |
| `--system-prompt` | string | System prompt injected into the in-app AI Assistant. |
| `--use-by-default` | bool | Whether this personality is the default personality for the Workspace. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli ai-assistant-personalities delete`

Delete Ai Assistant Personality.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Ai Assistant Personality ID. **Required.** |

