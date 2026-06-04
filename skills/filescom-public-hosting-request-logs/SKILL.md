---
name: filescom-public-hosting-request-logs
description: |
  A PublicHostingRequestLog is an audit log for monitoring requests we received to access your publicly served folders.
---

# filescom-public-hosting-request-logs

A PublicHostingRequestLog is an audit log for monitoring requests we received to access your publicly served folders.

Logs available through this endpoint are retained for 6 months, after which they are automatically discarded. For longer retention, use Files.com SIEM integrations to stream logs in real time to your preferred SIEM, or configure SIEM streaming to a file.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli public-hosting-request-logs list`

List Public Hosting Request Logs.

| Flag | Type | Description |
| --- | --- | --- |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `path`, `remote_ip`, `success` or `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `path`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ path ]`, `[ remote_ip ]`, `[ success ]`, `[ created_at ]`, `[ path, remote_ip ]`, `[ path, success ]`, `[ path, created_at ]`, `[ remote_ip, success ]`, `[ remote_ip, created_at ]`, `[ success, created_at ]`, `[ path, remote_ip, success ]`, `[ path, remote_ip, created_at ]`, `[ path, success, created_at ]`, `[ remote_ip, success, created_at ]` or `[ path, remote_ip, success, created_at ]`. |

