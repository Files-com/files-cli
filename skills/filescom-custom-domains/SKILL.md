---
name: filescom-custom-domains
description: |
  A CustomDomain object represents an additional customer-owned domain that routes to a Files.com site.
---

# filescom-custom-domains

A CustomDomain object represents an additional customer-owned domain that routes to a Files.com site.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli custom-domains list`

List Custom Domains.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `id`. |

### `files-cli custom-domains find`

Show Custom Domain.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Custom Domain ID. **Required.** |

### `files-cli custom-domains create`

Create Custom Domain.

| Flag | Type | Description |
| --- | --- | --- |
| `--destination` | enum | Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. One of: `site_alias`, `public_hosting`, `s3_endpoint`, `unassigned`. |
| `--folder-behavior-id` | int64 | Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior.  Preserved as historical context when `destination` becomes `unassigned`. |
| `--ssl-certificate-id` | int64 | Current SSL certificate ID. |
| `--domain` | string | Customer-owned domain name. **Required.** |

### `files-cli custom-domains update`

Update Custom Domain.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Custom Domain ID. **Required.** |
| `--destination` | enum | Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. One of: `site_alias`, `public_hosting`, `s3_endpoint`, `unassigned`. |
| `--folder-behavior-id` | int64 | Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior.  Preserved as historical context when `destination` becomes `unassigned`. |
| `--ssl-certificate-id` | int64 | Current SSL certificate ID. |
| `--domain` | string | Customer-owned domain name. |

### `files-cli custom-domains delete`

Delete Custom Domain.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Custom Domain ID. **Required.** |

