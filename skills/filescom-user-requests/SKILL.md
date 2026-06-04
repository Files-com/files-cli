---
name: filescom-user-requests
description: |
  A UserRequest is an operation that allows anonymous users to place a request for access on the login screen to the site administrator.
---

# filescom-user-requests

A UserRequest is an operation that allows anonymous users to place a request for access on the login screen to the site administrator.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli user-requests list`

List User Requests.

No flags beyond the global ones.

### `files-cli user-requests find`

Show User Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Request ID. **Required.** |

### `files-cli user-requests create`

Create User Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Name of user requested **Required.** |
| `--email` | string | Email of user requested **Required.** |
| `--details` | string | Details of the user request **Required.** |
| `--company` | string | Company of the user requested |

### `files-cli user-requests delete`

Delete User Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Request ID. **Required.** |

