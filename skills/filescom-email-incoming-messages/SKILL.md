---
name: filescom-email-incoming-messages
description: |
  An EmailIncomingMessage is an audit log designed to track status of email to inbox uploads.
---

# filescom-email-incoming-messages

An EmailIncomingMessage is an audit log designed to track status of email to inbox uploads.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli email-incoming-messages list`

List Email Incoming Messages.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `workspace_id`, `sender`, `status` or `inbox_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `inbox_id`, `sender`, `status` or `workspace_id`. Valid field combinations are `[ inbox_id, created_at ]`, `[ sender, created_at ]`, `[ status, created_at ]`, `[ workspace_id, created_at ]`, `[ inbox_id, status ]`, `[ workspace_id, inbox_id ]`, `[ status, sender ]`, `[ workspace_id, sender ]`, `[ workspace_id, status ]`, `[ inbox_id, status, created_at ]`, `[ workspace_id, inbox_id, created_at ]`, `[ status, sender, created_at ]`, `[ workspace_id, sender, created_at ]`, `[ workspace_id, status, created_at ]`, `[ inbox_id, status, sender ]`, `[ workspace_id, inbox_id, status ]`, `[ workspace_id, status, sender ]`, `[ inbox_id, status, sender, created_at ]`, `[ workspace_id, inbox_id, status, created_at ]`, `[ workspace_id, status, sender, created_at ]`, `[ workspace_id, inbox_id, status, sender ]` or `[ workspace_id, inbox_id, status, sender, created_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `sender`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |

