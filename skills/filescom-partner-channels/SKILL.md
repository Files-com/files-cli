---
name: filescom-partner-channels
description: |
  A PartnerChannel defines a structured communication path within a Partner root folder, including directional folder names and partner-scoped routing configuration.
---

# filescom-partner-channels

A PartnerChannel defines a structured communication path within a Partner root folder, including directional folder names and partner-scoped routing configuration.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli partner-channels list`

List Partner Channels.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `path` or `partner_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `partner_id` and `workspace_id`. Valid field combinations are `[ workspace_id, partner_id ]`. |

### `files-cli partner-channels find`

Show Partner Channel.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Channel ID. **Required.** |

### `files-cli partner-channels create [path]`

Create Partner Channel.

| Flag | Type | Description |
| --- | --- | --- |
| `--from-partner-folder-name` | string | Optional Channel-level from-Partner folder name override. |
| `--from-partner-route-path` | string | Optional route path for files uploaded by the Partner. |
| `--to-partner-folder-name` | string | Optional Channel-level to-Partner folder name override. |
| `--to-partner-route-path` | string | Optional route path for files delivered to the Partner. |
| `--partner-id` | int64 | ID of the Partner this Channel belongs to. **Required.** |
| `--path` | string | Channel path relative to the Partner root folder. **Required.** |
| `--workspace-id` | int64 | ID of the Workspace associated with this Partner Channel. |

### `files-cli partner-channels update [path]`

Update Partner Channel.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Channel ID. **Required.** |
| `--from-partner-folder-name` | string | Optional Channel-level from-Partner folder name override. |
| `--from-partner-route-path` | string | Optional route path for files uploaded by the Partner. |
| `--to-partner-folder-name` | string | Optional Channel-level to-Partner folder name override. |
| `--to-partner-route-path` | string | Optional route path for files delivered to the Partner. |
| `--path` | string | Channel path relative to the Partner root folder. |

### `files-cli partner-channels delete`

Delete Partner Channel.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Channel ID. **Required.** |

