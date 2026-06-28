---
name: filescom-partners
description: |
  A Partner is a first-class entity that cleanly represents an external organization, enables delegated administration, and enforces strict boundaries.
---

# filescom-partners

A Partner is a first-class entity that cleanly represents an external organization, enables delegated administration, and enforces strict boundaries.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli partners list`

List Partners.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |

### `files-cli partners find`

Show Partner.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner ID. **Required.** |

### `files-cli partners create`

Create Partner.

| Flag | Type | Description |
| --- | --- | --- |
| `--ai-assistant-personality-id` | int64 | AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it. |
| `--allowed-ips` | string | A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists. |
| `--allow-bypassing-2fa-policies` | bool | Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users. |
| `--allow-credential-changes` | bool | Allow Partner Admins to change or reset credentials for users belonging to this Partner. |
| `--allow-providing-gpg-keys` | bool | Allow Partner Admins to provide GPG keys. |
| `--allow-user-creation` | bool | Allow Partner Admins to create users. |
| `--cc-emails-to-responsible-party` | bool | When `true`, emails sent to Partner users are copied to the responsible User or Group. |
| `--notes` | string | Notes about this Partner. |
| `--responsible-group-id` | int64 | ID of the Group responsible for this Partner. |
| `--responsible-user-id` | int64 | ID of the User responsible for this Partner. |
| `--tags` | string | Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens. |
| `--name` | string | The name of the Partner. **Required.** |
| `--root-folder` | string | The root folder path for this Partner. **Required.** |
| `--workspace-id` | int64 | ID of the Workspace associated with this Partner. |

### `files-cli partners update`

Update Partner.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner ID. **Required.** |
| `--ai-assistant-personality-id` | int64 | AI Assistant Personality ID assigned to this Partner, if any. Users in the Partner inherit it unless a direct per-user assignment overrides it. |
| `--allowed-ips` | string | A list of allowed IPs for this Partner. Newline delimited. Partner User IP access is allowed when the IP matches the Partner, User, or Site allowed IP lists. |
| `--allow-bypassing-2fa-policies` | bool | Allow Partner Admins to change Two-Factor Authentication requirements for Partner Users. |
| `--allow-credential-changes` | bool | Allow Partner Admins to change or reset credentials for users belonging to this Partner. |
| `--allow-providing-gpg-keys` | bool | Allow Partner Admins to provide GPG keys. |
| `--allow-user-creation` | bool | Allow Partner Admins to create users. |
| `--cc-emails-to-responsible-party` | bool | When `true`, emails sent to Partner users are copied to the responsible User or Group. |
| `--notes` | string | Notes about this Partner. |
| `--responsible-group-id` | int64 | ID of the Group responsible for this Partner. |
| `--responsible-user-id` | int64 | ID of the User responsible for this Partner. |
| `--tags` | string | Comma-separated list of Tags for this Partner. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens. |
| `--name` | string | The name of the Partner. |
| `--root-folder` | string | The root folder path for this Partner. |

### `files-cli partners delete`

Delete Partner.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Partner ID. **Required.** |

