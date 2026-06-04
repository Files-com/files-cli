---
name: filescom-apps
description: |
  An App represents one of the various integrations provided by Files.com.
---

# filescom-apps

An App represents one of the various integrations provided by Files.com.  These are what are listed in the Integrations Catalog in the web interface.

Currently, all apps are developed internally by Files.com, though we do have the infrastructure to support vendor-developed apps.  If you are a third party vendor interested in developing for the Files.com ecosystem, please contact us.  We'd love to hear more.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli apps list`

List Apps.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name` and `app_type`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `name` and `app_type`. Valid field combinations are `[ name, app_type ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |

