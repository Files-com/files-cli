---
name: filescom-site-subdomain-redirects
description: |
  A SiteSubdomainRedirect object represents an old Files.com subdomain that continues to work after the site's Files.com subdomain changes.
---

# filescom-site-subdomain-redirects

A SiteSubdomainRedirect object represents an old Files.com subdomain that continues to work after the site's Files.com subdomain changes.
HTTPS requests redirect to the current subdomain, and other protocols such as FTP and SFTP are routed through DNS.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli site-subdomain-redirects list`

List Site Subdomain Redirects.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `id`. |

### `files-cli site-subdomain-redirects find`

Show Site Subdomain Redirect.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Site Subdomain Redirect ID. **Required.** |

### `files-cli site-subdomain-redirects delete`

Delete Site Subdomain Redirect.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Site Subdomain Redirect ID. **Required.** |

