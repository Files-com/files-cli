---
name: filescom-inbox-registrations
description: |
  An InboxRegistration is created when a user fills out the form to access the inbox.
---

# filescom-inbox-registrations

An InboxRegistration is created when a user fills out the form to access the inbox.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli inbox-registrations list`

List Inbox Registrations.

| Flag | Type | Description |
| --- | --- | --- |
| `--folder-behavior-id` | int64 | ID of the associated Inbox. This is required if the user is not a site admin. |

