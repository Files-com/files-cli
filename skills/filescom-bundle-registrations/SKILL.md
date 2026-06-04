---
name: filescom-bundle-registrations
description: |
  A BundleRegistration is a registration record when a user fills out the form to access the bundle.
---

# filescom-bundle-registrations

A BundleRegistration is a registration record when a user fills out the form to access the bundle.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli bundle-registrations list`

List Share Link Registrations.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `bundle_id` or `created_at`. |
| `--bundle-id` | int64 | ID of the associated Bundle |

