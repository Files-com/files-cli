---
name: filescom-locks
description: |
  A Lock can be used by your custom-developed applications to implement file locking and concurrency features.
---

# filescom-locks

A Lock can be used by your custom-developed applications to implement file locking and concurrency features. These locks are advisory, meaning that while a lock can be created, it does not prevent other API requests from being processed concurrently.  You are responsible for checking locks prior to accessing a file.

The lock feature is designed to emulate the locking functionality provided by WebDAV. For a deeper understanding of how the lock mechanism works, refer to the WebDAV specification, which outlines how these endpoints function.

Files.com's WebDAV offering and desktop app leverage this locking API to manage concurrent file operations, ensuring consistency when multiple users or systems interact with the same files.  It is not used within the Files.com web interface.

The optional owner parameter is a descriptive label, not the lock creator or a grant of permission. It can be set when creating a lock; refreshing a lock or replacing its token preserves it.

To refresh only an existing lock or replace its token, send expected_token, token, and timeout to the create endpoint. Set token to expected_token to refresh, or to a different value to replace. The expected token must identify an existing, unexpired lock on that path, and the caller must have permission to modify it. The token check and update happen together; invalid replacement values leave the stored lock unchanged.

A missing, expired, or mismatched expected token returns processing-failure/resource-locked with data.lock_token containing an active token on that path, or an empty string when none exists. Omitting expected_token retains the existing acquire-or-refresh behavior. Shared locks retain their existing semantics.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli locks list-for [path]`

List Locks by Path.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--include-children` | bool | Include locks from children objects? |

### `files-cli locks create [path]`

Create Lock.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path **Required.** |
| `--token` | string | Lock token. With expected_token, use the same value to refresh or a different value to replace the existing token. |
| `--expected-token` | string | Require this existing, unexpired token before refreshing or replacing a lock. Set token to the same value to refresh, or a different value to replace. |
| `--allow-access-by-any-user` | bool | Can lock be modified by users other than its creator? |
| `--exclusive` | bool | Is lock exclusive? |
| `--recursive` | bool | Does lock apply to subfolders? |
| `--owner` | string | Arbitrary descriptive label for the lock. Does not change the lock creator or permissions. |
| `--timeout` | int64 | Lock timeout in seconds |

### `files-cli locks delete [path]`

Delete Lock.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path **Required.** |
| `--token` | string | Lock token **Required.** |

