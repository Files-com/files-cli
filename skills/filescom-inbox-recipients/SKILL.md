---
name: filescom-inbox-recipients
description: |
  An InboxRecipient is a person who has had an inbox shared with them via email invitation.
---

# filescom-inbox-recipients

An InboxRecipient is a person who has had an inbox shared with them via email invitation.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli inbox-recipients list`

List Inbox Recipients.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `has_registrations`. |
| `--inbox-id` | int64 | List recipients for the inbox with this ID. **Required.** |

### `files-cli inbox-recipients create`

Create Inbox Recipient.

| Flag | Type | Description |
| --- | --- | --- |
| `--inbox-id` | int64 | Inbox to share. **Required.** |
| `--recipient` | string | Email address to share this inbox with. **Required.** |
| `--name` | string | Name of recipient. |
| `--company` | string | Company of recipient. |
| `--note` | string | Note to include in email. |
| `--share-after-create` | bool | Set to true to share the link with the recipient upon creation. |

