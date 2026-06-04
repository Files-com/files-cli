---
name: filescom-bundle-recipients
description: |
  A BundleRecipient is a person who has had a bundle shared with them via email invite.
---

# filescom-bundle-recipients

A BundleRecipient is a person who has had a bundle shared with them via email invite. A Bundle can be re-shared
with a Bundle Recipient by sending a create request with the inbox_id, recipient email address,
and share_after_create => true.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli bundle-recipients list`

List Share Link Recipients.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `has_registrations`. |
| `--bundle-id` | int64 | List recipients for the bundle with this ID. **Required.** |

### `files-cli bundle-recipients create`

Create Share Link Recipient.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--bundle-id` | int64 | Bundle to share. **Required.** |
| `--recipient` | string | Email addresses to share this bundle with. **Required.** |
| `--name` | string | Name of recipient. |
| `--company` | string | Company of recipient. |
| `--note` | string | Note to include in email. |
| `--share-after-create` | bool | Set to true to share the link with the recipient upon creation. |

