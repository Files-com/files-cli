---
name: filescom-user-additional-email-recipients
description: |
  Files.com User Additional Email Recipients via files-cli.
---

# filescom-user-additional-email-recipients

Files.com User Additional Email Recipients management via files-cli.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli user-additional-email-recipients list`

List User Additional Email Recipients.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `email`, `user_id` or `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `email` and `workspace_id`. Valid field combinations are `[ workspace_id, email ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `email`. |

### `files-cli user-additional-email-recipients find`

Show User Additional Email Recipient.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Additional Email Recipient ID. **Required.** |

### `files-cli user-additional-email-recipients create`

Create User Additional Email Recipient.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--email` | string | Additional email recipient address **Required.** |

### `files-cli user-additional-email-recipients update`

Update User Additional Email Recipient.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Additional Email Recipient ID. **Required.** |
| `--email` | string | Additional email recipient address |

### `files-cli user-additional-email-recipients delete`

Delete User Additional Email Recipient.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Additional Email Recipient ID. **Required.** |

