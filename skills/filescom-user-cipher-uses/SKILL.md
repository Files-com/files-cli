---
name: filescom-user-cipher-uses
description: |
  A UserCipherUse is a way to see the exact set of encryption ciphers and protocols used by a given user.
---

# filescom-user-cipher-uses

A UserCipherUse is a way to see the exact set of encryption ciphers and protocols used by a given user.
This is most often used to support migrations from one TLS version to the next.  You can query each user and determine who is still using legacy ciphers.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli user-cipher-uses list`

List User Cipher Uses.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID. If provided, will return uses for this user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `updated_at`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `insecure` and `updated_at`. Valid field combinations are `[ insecure, updated_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `updated_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `updated_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `updated_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `updated_at`. |

