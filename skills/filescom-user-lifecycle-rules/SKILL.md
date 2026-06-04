---
name: filescom-user-lifecycle-rules
description: |
  A UserLifecycleRule represents a rule that applies to users based on their inactivity, state and authentication method.
---

# filescom-user-lifecycle-rules

A UserLifecycleRule represents a rule that applies to users based on their inactivity, state and authentication method.

The rule either disable or delete users who have been inactive or disabled for a specified number of days.

The authentication_method property specifies the authentication method for the rule, which can be set to "all", "all_non_sso", or a specific authentication method.

The rule can also include or exclude site and folder admins from the action.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli user-lifecycle-rules list`

List User Lifecycle Rules.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id` and `workspace_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |

### `files-cli user-lifecycle-rules find`

Show User Lifecycle Rule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Lifecycle Rule ID. **Required.** |

### `files-cli user-lifecycle-rules create`

Create User Lifecycle Rule.

| Flag | Type | Description |
| --- | --- | --- |
| `--apply-to-all-workspaces` | bool | If true, a default-workspace rule also applies to users in all workspaces. |
| `--authentication-method` | enum | User authentication method for which the rule will apply. Use `all_non_sso` to target every non-SSO authentication method with one rule. One of: `all`, `password`, `sso`, `none`, `email_signup`, `password_with_imported_hash`, `password_and_ssh_key`, `all_non_sso`. |
| `--group-ids` | []int64 | Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users. |
| `--inactivity-days` | int64 | Number of days of inactivity before the rule applies |
| `--include-site-admins` | bool | If true, the rule will apply to site admins. |
| `--include-folder-admins` | bool | If true, the rule will apply to folder admins. |
| `--name` | string | User Lifecycle Rule name |
| `--notify-users` | bool | If true, users will be emailed before the rule disables or deletes them. |
| `--partner-tag` | string | If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `--user-state` | enum | State of the users to apply the rule to (inactive or disabled). One of: `inactive`, `disabled`. |
| `--user-tag` | string | If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli user-lifecycle-rules update`

Update User Lifecycle Rule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Lifecycle Rule ID. **Required.** |
| `--apply-to-all-workspaces` | bool | If true, a default-workspace rule also applies to users in all workspaces. |
| `--authentication-method` | enum | User authentication method for which the rule will apply. Use `all_non_sso` to target every non-SSO authentication method with one rule. One of: `all`, `password`, `sso`, `none`, `email_signup`, `password_with_imported_hash`, `password_and_ssh_key`, `all_non_sso`. |
| `--group-ids` | []int64 | Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users. |
| `--inactivity-days` | int64 | Number of days of inactivity before the rule applies |
| `--include-site-admins` | bool | If true, the rule will apply to site admins. |
| `--include-folder-admins` | bool | If true, the rule will apply to folder admins. |
| `--name` | string | User Lifecycle Rule name |
| `--notify-users` | bool | If true, users will be emailed before the rule disables or deletes them. |
| `--partner-tag` | string | If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `--user-state` | enum | State of the users to apply the rule to (inactive or disabled). One of: `inactive`, `disabled`. |
| `--user-tag` | string | If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |

### `files-cli user-lifecycle-rules delete`

Delete User Lifecycle Rule.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User Lifecycle Rule ID. **Required.** |

