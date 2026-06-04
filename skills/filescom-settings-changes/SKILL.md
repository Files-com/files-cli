---
name: filescom-settings-changes
description: |
  A SettingsChange is any change in your site caused by any user.
---

# filescom-settings-changes

A SettingsChange is any change in your site caused by any user.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli settings-changes list`

List Settings Changes.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `api_key_id` or `user_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `api_key_id` and `user_id`. |

