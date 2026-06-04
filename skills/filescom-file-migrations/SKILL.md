---
name: filescom-file-migrations
description: |
  A FileMigration is a background operation on one or more files, such as a copy or a region migration.
---

# filescom-file-migrations

A FileMigration is a background operation on one or more files, such as a copy or a region migration.

If no `operation` or `dest_path` is present, then the record represents a region migration.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli file-migrations find`

Show File Migration.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | File Migration ID. **Required.** |

