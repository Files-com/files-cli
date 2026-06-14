---
name: filescom-partner-site-requests
description: |
  A PartnerSiteRequest represents a request for a Guest Partner to add their Files.com Site to their Partnership with the Host Partner.
---

# filescom-partner-site-requests

A PartnerSiteRequest represents a request for a Guest Partner to add their Files.com Site to their Partnership with the Host Partner. The Guest Partner's Files.com Site is referred to as the Guest Site in this relationship.

The Partner Admin user representing the Guest on the Host Partner can initiate a request, which generates a pairing key. The Guest Site admin must then approve the request. This ensures that the Partner Admin user representing the Guest on the Host Partner and the Site Admins of the Site are in agreement that the linking should occur.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli partner-site-requests list`

List Partner Site Requests.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `host_partner_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `host_partner_id`. |

### `files-cli partner-site-requests find-by-pairing-key`

Find partner site request by pairing key.

| Flag | Type | Description |
| --- | --- | --- |
| `--pairing-key` | string | Pairing key for the partner site request **Required.** |

### `files-cli partner-site-requests create`

Create Partner Site Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--host-partner-id` | int64 | Host Partner ID to link with **Required.** |
| `--guest-site-url` | string | Guest Site URL to link to **Required.** |

### `files-cli partner-site-requests reject`

Reject partner site request.

| Flag | Type | Description |
| --- | --- | --- |
| `--pairing-key` | string | Pairing key for the partner site request **Required.** |

### `files-cli partner-site-requests approve`

Approve partner site request.

| Flag | Type | Description |
| --- | --- | --- |
| `--pairing-key` | string | Pairing key for the partner site request **Required.** |

### `files-cli partner-site-requests delete`

Delete Partner Site Request.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Site Request ID. **Required.** |

