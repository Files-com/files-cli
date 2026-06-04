---
name: filescom-public-keys
description: |
  A PublicKey is used to authenticate to Files.com via SFTP (SSH File Transfer Protocol).
---

# filescom-public-keys

A PublicKey is used to authenticate to Files.com via SFTP (SSH File Transfer Protocol). This method of authentication allows users to use their private key (which is never shared with Files.com) to authenticate themselves against the PublicKey stored on Files.com.

When a user configures their PublicKey, it allows them to bypass traditional password-based authentication, leveraging the security of key-based authentication instead.

Note that Files.com's SSH support is limited to file operations only. While users can securely transfer files and manage their data via SFTP, they do not have access to a full shell environment for executing arbitrary commands.

When generating new SSH keys, here are the available options: Files.com supports multiple SSH key algorithms: RSA (default 4096 bits, range 1024-4096 in 8-bit increments), DSA (1024 bits only), Ed25519 (256 bits), and ECDSA (256, 384, or 521 bits). When generating keys, the system uses these default lengths unless a specific length is specified.

Files.com also supports importing additional key types that cannot be generated: security key types (sk-ecdsa-sha2-nistp256, sk-ssh-ed25519). RSA keys up to 8192 bits are also supported for import.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli public-keys list`

List Public Keys.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `user_id`, `title` or `created_at`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `workspace_id`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |

### `files-cli public-keys find`

Show Public Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Public Key ID. **Required.** |

### `files-cli public-keys create`

Create Public Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--title` | string | Internal reference for key. **Required.** |
| `--public-key` | string | Actual contents of SSH key. |
| `--generate-keypair` | bool | If true, generate a new SSH key pair. Can not be used with `public_key` |
| `--generate-private-key-password` | string | Password for the private key. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `--generate-algorithm` | string | Type of key to generate.  One of rsa, dsa, ecdsa, ed25519. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `--generate-length` | int64 | Length of key to generate. If algorithm is ecdsa, this is the signature size. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |

### `files-cli public-keys update`

Update Public Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Public Key ID. **Required.** |
| `--title` | string | Internal reference for key. **Required.** |

### `files-cli public-keys delete`

Delete Public Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Public Key ID. **Required.** |

