---
name: filescom-groups
description: |
  A Group is a powerful tool for permissions and user management on Files.com.
---

# filescom-groups

A Group is a powerful tool for permissions and user management on Files.com.  Users can belong to multiple groups.

All permissions can be managed via Groups, and Groups can also be synced to your identity platform via LDAP or SCIM.

Files.com's Group Admin feature allows you to define Group Admins, who then have access to add and remove users within their groups.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli groups list`

List Groups.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id` or `name`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `name` and `workspace_id`. Valid field combinations are `[ workspace_id, name ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |
| `--ids` | string | Comma-separated list of group ids to include in results. |
| `--include-parent-site-groups` | bool | Include groups from the parent site. |

### `files-cli groups find`

Show Group.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Group ID. **Required.** |

### `files-cli groups create`

Create Group.

| Flag | Type | Description |
| --- | --- | --- |
| `--notes` | string | Group notes. |
| `--user-ids` | string | A list of user ids. If sent as a string, should be comma-delimited. |
| `--admin-ids` | string | A list of group admin user ids. If sent as a string, should be comma-delimited. |
| `--ai-assistant-personality-id` | int64 | AI Assistant Personality ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user or Partner assignment overrides it. |
| `--ftp-permission` | bool | If true, users in this group can use FTP to login.  This will override a false value of `ftp_permission` on the user level. |
| `--sftp-permission` | bool | If true, users in this group can use SFTP to login.  This will override a false value of `sftp_permission` on the user level. |
| `--dav-permission` | bool | If true, users in this group can use WebDAV to login.  This will override a false value of `dav_permission` on the user level. |
| `--restapi-permission` | bool | If true, users in this group can use the REST API to login.  This will override a false value of `restapi_permission` on the user level. |
| `--desktop-configuration-profile-id` | int64 | Desktop Configuration Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `--integration-centric-profile-id` | int64 | Integration Centric Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `--allowed-ips` | string | A list of allowed IPs if applicable.  Newline delimited |
| `--name` | string | Group name. **Required.** |
| `--workspace-id` | int64 | Workspace ID |

### `files-cli groups update`

Update Group.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Group ID. **Required.** |
| `--notes` | string | Group notes. |
| `--user-ids` | string | A list of user ids. If sent as a string, should be comma-delimited. |
| `--admin-ids` | string | A list of group admin user ids. If sent as a string, should be comma-delimited. |
| `--ai-assistant-personality-id` | int64 | AI Assistant Personality ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user or Partner assignment overrides it. |
| `--ftp-permission` | bool | If true, users in this group can use FTP to login.  This will override a false value of `ftp_permission` on the user level. |
| `--sftp-permission` | bool | If true, users in this group can use SFTP to login.  This will override a false value of `sftp_permission` on the user level. |
| `--dav-permission` | bool | If true, users in this group can use WebDAV to login.  This will override a false value of `dav_permission` on the user level. |
| `--restapi-permission` | bool | If true, users in this group can use the REST API to login.  This will override a false value of `restapi_permission` on the user level. |
| `--desktop-configuration-profile-id` | int64 | Desktop Configuration Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `--integration-centric-profile-id` | int64 | Integration Centric Profile ID assigned to this Group, if any. Users in the Group inherit it unless a direct per-user assignment overrides it. |
| `--allowed-ips` | string | A list of allowed IPs if applicable.  Newline delimited |
| `--name` | string | Group name. |

### `files-cli groups delete`

Delete Group.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Group ID. **Required.** |

