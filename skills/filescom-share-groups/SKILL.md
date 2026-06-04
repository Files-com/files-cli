---
name: filescom-share-groups
description: |
  A ShareGroup is a way for you to store and name groups of email contacts to be used for sending share and inbox invitations.
---

# filescom-share-groups

A ShareGroup is a way for you to store and name groups of email contacts to be used for sending share and inbox invitations.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli share-groups list`

List Share Groups.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |

### `files-cli share-groups find`

Show Share Group.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Share Group ID. **Required.** |

### `files-cli share-groups create`

Create Share Group.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--notes` | string | Additional notes of the share group |
| `--name` | string | Name of the share group **Required.** |
| `--members` | []object | A list of share group members. **Required.** |

### `files-cli share-groups update`

Update Share Group.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Share Group ID. **Required.** |
| `--notes` | string | Additional notes of the share group |
| `--name` | string | Name of the share group |
| `--members` | []object | A list of share group members. |

### `files-cli share-groups delete`

Delete Share Group.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Share Group ID. **Required.** |

