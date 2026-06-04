---
name: filescom-sftp-host-keys
description: |
  An SFTP Host Key is a cryptographic key used to verify the identity of the server during an SFTP connection.
---

# filescom-sftp-host-keys

An SFTP Host Key is a cryptographic key used to verify the identity of the server during an SFTP connection. This allows the client to be sure that it is connecting to the intended server, preventing man-in-the-middle attacks and ensuring secure communication between the client and Files.com.

Files.com allows you to provide custom SFTP Host Keys, which is particularly useful when migrating to Files.com from an existing SFTP server, allowing the Files.com platform to match your previously-installed host key for a seamless transition.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli sftp-host-keys list`

List SFTP Host Keys.

No flags beyond the global ones.

### `files-cli sftp-host-keys find`

Show SFTP Host Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sftp Host Key ID. **Required.** |

### `files-cli sftp-host-keys create`

Create SFTP Host Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | The friendly name of this SFTP Host Key. |
| `--private-key` | string | The private key data. |

### `files-cli sftp-host-keys update`

Update SFTP Host Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sftp Host Key ID. **Required.** |
| `--name` | string | The friendly name of this SFTP Host Key. |
| `--private-key` | string | The private key data. |

### `files-cli sftp-host-keys delete`

Delete SFTP Host Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Sftp Host Key ID. **Required.** |

