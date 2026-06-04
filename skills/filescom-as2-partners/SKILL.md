---
name: filescom-as2-partners
description: |
  An AS2Partner is a counterparty of the Files.com site's AS2 connectivity.
---

# filescom-as2-partners

An AS2Partner is a counterparty of the Files.com site's AS2 connectivity.  Generally you will have one AS2 Partner created for each counterparty with whom you send and/or receive files via AS2.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli as2-partners list`

List AS2 Partners.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `as2_station_id` or `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `as2_station_id` and `workspace_id`. Valid field combinations are `[ workspace_id, as2_station_id ]`. |

### `files-cli as2-partners find`

Show AS2 Partner.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | As2 Partner ID. **Required.** |

### `files-cli as2-partners create`

Create AS2 Partner.

| Flag | Type | Description |
| --- | --- | --- |
| `--enable-dedicated-ips` | bool | If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 Partner. |
| `--http-auth-username` | string | Username to send to server for HTTP Authentication. |
| `--http-auth-password` | string | Password to send to server for HTTP Authentication. |
| `--mdn-validation-level` | enum | How should Files.com evaluate message transfer success based on a partner's MDN response?  This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates. `auto`: Automatically set the correct value for this setting based on next mdn received. One of: `none`, `weak`, `normal`, `strict`, `auto`. |
| `--signature-validation-level` | enum | Should Files.com require signatures on incoming AS2 messages?  `normal`: require that incoming messages are signed with a valid matching signature. `none`: Unsigned incoming messages are allowed. `auto`: Automatically set the correct value for this setting based on next message received. One of: `normal`, `none`, `auto`. |
| `--server-certificate` | enum | Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS? (This only applies to Outgoing AS2 message from Files.com to a Partner.). One of: `require_match`, `allow_any`. |
| `--default-mime-type` | string | Default mime type of the file attached to the encrypted message |
| `--additional-http-headers` | object | Additional HTTP Headers for outgoing message sent to this partner. |
| `--as2-station-id` | int64 | ID of the AS2 Station associated with this partner. **Required.** |
| `--name` | string | The partner's formal AS2 name. **Required.** |
| `--uri` | string | Public URI where we will send the AS2 messages (via HTTP/HTTPS). **Required.** |
| `--public-certificate` | string | Public certificate for AS2 Partner.  Note: This is the certificate for AS2 message security, not a certificate used for HTTPS authentication. **Required.** |

### `files-cli as2-partners update`

Update AS2 Partner.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | As2 Partner ID. **Required.** |
| `--enable-dedicated-ips` | bool | If `true`, we will use your site's dedicated IPs for all outbound connections to this AS2 Partner. |
| `--http-auth-username` | string | Username to send to server for HTTP Authentication. |
| `--http-auth-password` | string | Password to send to server for HTTP Authentication. |
| `--mdn-validation-level` | enum | How should Files.com evaluate message transfer success based on a partner's MDN response?  This setting does not affect MDN storage; all MDNs received from a partner are always stored. `none`: MDN is stored for informational purposes only, a successful HTTPS transfer is a successful AS2 transfer. `weak`: Inspect the MDN for MIC and Disposition only. `normal`: `weak` plus validate MDN signature matches body, `strict`: `normal` but do not allow signatures from self-signed or incorrectly purposed certificates. `auto`: Automatically set the correct value for this setting based on next mdn received. One of: `none`, `weak`, `normal`, `strict`, `auto`. |
| `--signature-validation-level` | enum | Should Files.com require signatures on incoming AS2 messages?  `normal`: require that incoming messages are signed with a valid matching signature. `none`: Unsigned incoming messages are allowed. `auto`: Automatically set the correct value for this setting based on next message received. One of: `normal`, `none`, `auto`. |
| `--server-certificate` | enum | Should we require that the remote HTTP server have a valid SSL Certificate for HTTPS? (This only applies to Outgoing AS2 message from Files.com to a Partner.). One of: `require_match`, `allow_any`. |
| `--default-mime-type` | string | Default mime type of the file attached to the encrypted message |
| `--additional-http-headers` | object | Additional HTTP Headers for outgoing message sent to this partner. |
| `--name` | string | The partner's formal AS2 name. |
| `--uri` | string | Public URI where we will send the AS2 messages (via HTTP/HTTPS). |
| `--public-certificate` | string | Public certificate for AS2 Partner.  Note: This is the certificate for AS2 message security, not a certificate used for HTTPS authentication. |

### `files-cli as2-partners delete`

Delete AS2 Partner.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | As2 Partner ID. **Required.** |

