---
name: filescom-styles
description: |
  A Style is a custom set of branding that can be applied on a per-folder basis.
---

# filescom-styles

A Style is a custom set of branding that can be applied on a per-folder basis.
Currently these support a logo per folder and an optional click-through URL for public visitors.
In the future we may extend these to also support colors.
If you want to see that, please let us know so we can add your vote to the list.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli styles find [path]`

Show Style.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Style path. **Required.** |

### `files-cli styles update [path]`

Update Style.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Style path. **Required.** |
| `--file` | file | Logo for custom branding. Required when creating a new style. |
| `--logo-click-href` | string | URL to open when a public visitor clicks the logo. |

### `files-cli styles delete [path]`

Delete Style.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Style path. **Required.** |

