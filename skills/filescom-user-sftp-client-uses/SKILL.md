---
name: filescom-user-sftp-client-uses
description: |
  A UserSftpClientUse is a way to see the exact set of SFTP clients used by a given user.
---

# filescom-user-sftp-client-uses

A UserSftpClientUse is a way to see the exact set of SFTP clients used by a given user.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli user-sftp-client-uses list`

List User SFTP Client Uses.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID. If provided, will return uses for this user. |

