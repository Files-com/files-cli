---
name: filescom-partner-site-requests
description: |
  A PartnerSiteRequest represents a request to link a partner's Files.com site with another Files.com site.
---

# filescom-partner-site-requests

A PartnerSiteRequest represents a request to link a partner's Files.com site with another Files.com site.

The Site with the Partner can initiate a request, which generates a pairing key. The target site admin must then approve the request using the pairing key.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli partner-site-requests list`

List Partner Site Requests.

No flags beyond the global ones.

### `files-cli partner-site-requests find-by-pairing-key`

Find partner site request by pairing key.

| Flag | Type | Description |
| --- | --- | --- |
| `--pairing-key` | string | Pairing key for the partner site request **Required.** |

### `files-cli partner-site-requests create`

Create Partner Site Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--partner-id` | int64 | Partner ID to link with **Required.** |
| `--site-url` | string | Site URL to link to **Required.** |

### `files-cli partner-site-requests reject`

Reject partner site request.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Site Request ID. **Required.** |

### `files-cli partner-site-requests approve`

Approve partner site request.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Site Request ID. **Required.** |

### `files-cli partner-site-requests delete`

Delete Partner Site Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Site Request ID. **Required.** |

