---
name: filescom-database-destinations
description: |
  A Database Destination lets an Automation append rows to an existing table.
---

# filescom-database-destinations

A Database Destination lets an Automation append rows to an existing table. Each destination references a credential in the Secret Store.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli database-destinations list`

List Database Destinations.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id` and `name`. Valid field combinations are `[ workspace_id, name ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |

### `files-cli database-destinations find`

Show Database Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Database Destination ID. **Required.** |

### `files-cli database-destinations create`

Create Database Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Name of the destination. **Required.** |
| `--provider` | enum | Database provider. One of: `snowflake`, `bigquery`, `supabase`, `postgresql`, `mysql`. **Required.** |
| `--config` | object | Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation. **Required.** |
| `--secret-id` | int64 | ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login. **Required.** |
| `--workspace-id` | int64 | Workspace ID. 0 means the default workspace. |

### `files-cli database-destinations update`

Update Database Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Database Destination ID. **Required.** |
| `--name` | string | Name of the destination. |
| `--secret-id` | int64 | ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login. |
| `--config` | object | Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation. |

### `files-cli database-destinations delete`

Delete Database Destination.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Database Destination ID. **Required.** |

