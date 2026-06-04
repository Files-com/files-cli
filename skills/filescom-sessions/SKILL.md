---
name: filescom-sessions
description: |
  A Session is an operation that allows you to make further API calls using our REST API or SDKs as a specific user.
---

# filescom-sessions

A Session is an operation that allows you to make further API calls using our REST API or SDKs as a specific user.
This is the only way to use the API if you know a username/password but not an API key.

Sessions in the API and SDKs use the exact same mechanism (and work the same) as sessions in the web interface.

After creating a session, the Session object returned will include plenty of relevant information about the current user, often used to customize the interface or enable further automation.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli sessions create`

Create user session (log in).

| Flag | Type | Description |
| --- | --- | --- |
| `--username` | string | Username to sign in as |
| `--password` | string | Password for sign in |
| `--otp` | string | If this user has a 2FA device, provide its OTP or code here. |
| `--partial-session-id` | string | Identifier for a partially-completed login |

### `files-cli sessions delete`

Delete user session (log out).

No flags beyond the global ones.

