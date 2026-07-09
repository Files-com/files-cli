---
name: filescom-partner-channel-templates
description: |
  A PartnerChannelTemplate defines reusable Partner Channel configuration that can be applied to Partners.
---

# filescom-partner-channel-templates

A PartnerChannelTemplate defines reusable Partner Channel configuration that can be applied to Partners.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli partner-channel-templates list`

List Partner Channel Templates.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |

### `files-cli partner-channel-templates find`

Show Partner Channel Template.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Channel Template ID. **Required.** |

### `files-cli partner-channel-templates create [path]`

Create Partner Channel Template.

| Flag | Type | Description |
| --- | --- | --- |
| `--from-partner-folder-name` | string | Optional Channel-level from-Partner folder name override. |
| `--from-partner-managed-folder-paths` | []string | Managed folder paths inside the from-Partner folder. |
| `--from-partner-route-path` | string | Optional route path for files uploaded by the Partner. |
| `--to-partner-folder-name` | string | Optional Channel-level to-Partner folder name override. |
| `--to-partner-managed-folder-paths` | []string | Managed folder paths inside the to-Partner folder. |
| `--to-partner-route-path` | string | Optional route path for files delivered to the Partner. |
| `--name` | string | The name of the Partner Channel Template. **Required.** |
| `--path` | string | Channel path relative to the Partner root folder. **Required.** |
| `--workspace-id` | int64 | ID of the Workspace associated with this Partner Channel Template. |

### `files-cli partner-channel-templates update [path]`

Update Partner Channel Template.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Channel Template ID. **Required.** |
| `--from-partner-folder-name` | string | Optional Channel-level from-Partner folder name override. |
| `--from-partner-managed-folder-paths` | []string | Managed folder paths inside the from-Partner folder. |
| `--from-partner-route-path` | string | Optional route path for files uploaded by the Partner. |
| `--to-partner-folder-name` | string | Optional Channel-level to-Partner folder name override. |
| `--to-partner-managed-folder-paths` | []string | Managed folder paths inside the to-Partner folder. |
| `--to-partner-route-path` | string | Optional route path for files delivered to the Partner. |
| `--name` | string | The name of the Partner Channel Template. |
| `--path` | string | Channel path relative to the Partner root folder. |

### `files-cli partner-channel-templates delete`

Delete Partner Channel Template.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner Channel Template ID. **Required.** |

