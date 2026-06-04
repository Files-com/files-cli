---
name: filescom-locks
description: |
  A Lock can be used by your custom-developed applications to implement file locking and concurrency features.
---

# filescom-locks

A Lock can be used by your custom-developed applications to implement file locking and concurrency features. These locks are advisory, meaning that while a lock can be created, it does not prevent other API requests from being processed concurrently.  You are responsible for checking locks prior to accessing a file.

The lock feature is designed to emulate the locking functionality provided by WebDAV. For a deeper understanding of how the lock mechanism works, refer to the WebDAV specification, which outlines how these endpoints function.

Files.com's WebDAV offering and desktop app leverage this locking API to manage concurrent file operations, ensuring consistency when multiple users or systems interact with the same files.  It is not used within the Files.com web interface.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

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
| `--allow-access-by-any-user` | bool | Can lock be modified by users other than its creator? |
| `--exclusive` | bool | Is lock exclusive? |
| `--recursive` | bool | Does lock apply to subfolders? |
| `--timeout` | int64 | Lock timeout in seconds |

### `files-cli locks delete [path]`

Delete Lock.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path **Required.** |
| `--token` | string | Lock token **Required.** |

