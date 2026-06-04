---
name: filescom-as2-outgoing-messages
description: |
  An AS2OutgoingMessage is a single record created for each individual AS2 file transfer out to a Partner.
---

# filescom-as2-outgoing-messages

An AS2OutgoingMessage is a single record created for each individual AS2 file transfer out to a Partner. It contains the message details and tracks the lifecycle status with milestones of the transfer.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli as2-outgoing-messages list`

List AS2 Outgoing Messages.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `as2_partner_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `as2_station_id`, `workspace_id` or `as2_partner_id`. Valid field combinations are `[ as2_station_id, created_at ]`, `[ workspace_id, created_at ]`, `[ as2_partner_id, created_at ]`, `[ workspace_id, as2_station_id ]`, `[ workspace_id, as2_partner_id ]`, `[ workspace_id, as2_station_id, created_at ]` or `[ workspace_id, as2_partner_id, created_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |

