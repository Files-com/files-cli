---
name: filescom-as2-stations
description: |
  An AS2Station is a remote AS2 server that can send data into Files.com and receive data from Files.com.
---

# filescom-as2-stations

An AS2Station is a remote AS2 server that can send data into Files.com and receive data from Files.com.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli as2-stations list`

List AS2 Stations.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |

### `files-cli as2-stations find`

Show AS2 Station.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | As2 Station ID. **Required.** |

### `files-cli as2-stations create`

Create AS2 Station.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | The station's formal AS2 name. **Required.** |
| `--workspace-id` | int64 | ID of the Workspace associated with this AS2 Station. |
| `--public-certificate` | string | Public certificate used for message security. |
| `--private-key` | string | PEM-encoded private key matching public_certificate. |
| `--private-key-password` | string | Password for the PEM-encoded private key. |
| `--pkcs12` | string | Base64-encoded PKCS#12 (.pfx or .p12) file containing the identity's certificate and private key. Provide this instead of public_certificate and private_key. |
| `--pkcs12-password` | string | Password for pkcs12. The file and password are used only for import; the extracted certificate and private key are stored as PEM. |

### `files-cli as2-stations update`

Update AS2 Station.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | As2 Station ID. **Required.** |
| `--name` | string | The station's formal AS2 name. |
| `--public-certificate` | string | Public certificate used for message security. |
| `--private-key` | string | PEM-encoded private key matching public_certificate. |
| `--private-key-password` | string | Password for the PEM-encoded private key. |
| `--pkcs12` | string | Base64-encoded PKCS#12 (.pfx or .p12) file containing the identity's replacement certificate and private key. Provide this instead of public_certificate and private_key. |
| `--pkcs12-password` | string | Password for pkcs12. The file and password are used only for import; the extracted certificate and private key are stored as PEM. |

### `files-cli as2-stations delete`

Delete AS2 Station.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | As2 Station ID. **Required.** |

