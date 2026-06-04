---
name: filescom-gpg-keys
description: |
  A GPGKey object on Files.com is used to securely store both the private and public keys associated with a GPG (GNU Privacy Guard) encryption key pair.
---

# filescom-gpg-keys

A GPGKey object on Files.com is used to securely store both the private and public keys associated with a GPG (GNU Privacy Guard) encryption key pair. This object enables the encryption and decryption of data using GPG, allowing you to protect sensitive information.

The private key is kept confidential and is used for decrypting data or signing messages to prove authenticity, while the public key is used to encrypt messages that only the owner of the private key can decrypt.

By storing both keys together in a GPGKey object, Files.com makes it easier to understand encryption operations, ensuring secure and efficient handling of encrypted data within the platform.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli gpg-keys list`

List GPG Keys.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `expires_at`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `partner_id` or `expires_at`. Valid field combinations are `[ workspace_id, expires_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `expires_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `expires_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `expires_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `expires_at`. |

### `files-cli gpg-keys find`

Show GPG Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Gpg Key ID. **Required.** |

### `files-cli gpg-keys create`

Create GPG Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--partner-id` | int64 | Partner ID who owns this GPG Key, if applicable. |
| `--public-key` | string | The GPG public key |
| `--private-key` | string | The GPG private key |
| `--private-key-password` | string | The GPG private key password |
| `--name` | string | GPG key name. **Required.** |
| `--workspace-id` | int64 | Workspace ID (0 for default workspace). |
| `--generate-expires-at` | datetime | Expiration date of the key. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `--generate-keypair` | bool | If true, generate a new GPG key pair. Can not be used with `public_key`/`private_key` |
| `--generate-full-name` | string | Full name of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |
| `--generate-email` | string | Email address of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false. |

### `files-cli gpg-keys update`

Update GPG Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Gpg Key ID. **Required.** |
| `--partner-id` | int64 | Partner ID who owns this GPG Key, if applicable. |
| `--public-key` | string | The GPG public key |
| `--private-key` | string | The GPG private key |
| `--private-key-password` | string | The GPG private key password |
| `--name` | string | GPG key name. |

### `files-cli gpg-keys delete`

Delete GPG Key.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Gpg Key ID. **Required.** |

