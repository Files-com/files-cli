---
name: filescom-users
description: |
  A User represents a human or system/service user with the ability to connect to Files.com via any of the available connectivity methods (unless restricted to specific protocols). Pick `--authentication-method` carefully at creation. Each method behaves differently: `email_signup` mails a "set your password" link; `password` requires the admin to set one upfront; `sso` defers authentication to a configured SSO strategy; `none` allows only API key or SSH key login.
---

# filescom-users

A User represents a human or system/service user with the ability to connect to Files.com via any of the available connectivity methods (unless restricted to specific protocols).

Users are associated with API Keys, SSH (SFTP) Keys, Notifications, Permissions, and Group memberships.


## Authentication

The `authentication_method` property on a User determines exactly how that user can login and authenticate to their Files.com account. Files.com offers a variety of authentication methods to ensure flexibility, security, migration, and compliance.

These authentication methods can be configured during user creation and can be modified at any time by site administrators. The meanings of the available values are as follows:

* `password` - Allows authentication via a password.  If API Keys or SSH (SFTP) Keys are also configured, those can be used *instead* of the password.  If Two Factor Authentication (2FA) methods are also configured, a valid 2nd factor is required in addition to the password.
* `email_signup` - When set upon user creation, an email will be sent to the new user with a link for them to create their password. Once the user has created their password, their authentication type will change to `password`.
* `sso` - Allows authentication via a linked Single Sign On provider.  If API Keys or SSH (SFTP) Keys are also configured, those can be used *instead* of Single Sign On.  If Two Factor Authentication (2FA) methods are also configured, a valid 2nd factor is required in addition to Single Sign On.  When using this method, you must also provide a valid `sso_strategy_id` to associate the User to the appropriate SSO provider.
* `password_with_imported_hash` - Works like the `password` method but allows importing a hashed password in MD5, SHA-1, or SHA-256 format.  Provide the imported hash in the field `imported_password_hash`.  Upon first use, the password will be converted to Files.com's internal storage format and the authentication type will change to `password`. Typically only used when migrating to Files.com from another MFT solution.
* `none` - Does not allow authentication via username and password, but does allow authentication via API Key or SSH (SFTP) Key.  Typically only used for service users.
* `password_and_ssh_key` - Allows authentication only by providing a password and also a valid SSH (SFTP) Key in a single attempt.  If API Keys are also configured, those can be used *instead* of the password and key combination.  This method only works with (typically enterprise) SSH/SFTP clients capable of sending both authentication methods at once.  Typically only used for service users.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli users list`

List Users.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `site_id`, `workspace_id`, `company`, `name`, `disabled`, `authenticate_until`, `username`, `email`, `site_admin`, `last_desktop_login_at`, `last_login_at`, `password_validity_days` or `ssl_required`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `username`, `name`, `email`, `company`, `site_admin`, `password_validity_days`, `ssl_required`, `last_login_at`, `authenticate_until`, `not_site_admin`, `disabled`, `partner_id`, `primary_group_id` or `workspace_id`. Valid field combinations are `[ site_admin, username ]`, `[ not_site_admin, username ]`, `[ workspace_id, username ]`, `[ company, name ]`, `[ workspace_id, name ]`, `[ workspace_id, email ]`, `[ workspace_id, company ]`, `[ workspace_id, site_admin ]`, `[ workspace_id, not_site_admin ]`, `[ workspace_id, disabled ]`, `[ workspace_id, partner_id ]`, `[ workspace_id, site_admin, username ]`, `[ workspace_id, not_site_admin, username ]`, `[ workspace_id, disabled, username ]`, `[ workspace_id, partner_id, username ]` or `[ workspace_id, company, name ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `username`, `name`, `email` or `company`. Valid field combinations are `[ company, name ]`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `password_validity_days`, `last_login_at` or `authenticate_until`. |
| `--ids` | string | comma-separated list of User IDs |
| `--include-parent-site-users` | bool | Include users from the parent site. |
| `--search` | string | Searches for partial matches of name, username, or email. |

### `files-cli users find`

Show User.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User ID. **Required.** |

### `files-cli users create`

Create User.

| Flag | Type | Description |
| --- | --- | --- |
| `--avatar-file` | file | An image file for your user avatar. |
| `--avatar-delete` | bool | If true, the avatar will be deleted. |
| `--change-password` | string | Used for changing a password on an existing user. |
| `--change-password-confirmation` | string | Optional, but if provided, we will ensure that it matches the value sent in `change_password`. |
| `--email` | string | User's email. |
| `--grant-permission` | string | Permission to grant on the User Root upon user creation. Can be blank or `full`, `read`, `write`, `list`, `read+write`, or `list+write` |
| `--group-id` | int64 | Group ID to associate this user with. |
| `--group-ids` | string | A list of group ids to associate this user with.  Comma delimited. |
| `--imported-password-hash` | string | Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256. |
| `--password` | string | User password. |
| `--password-confirmation` | string | Optional, but if provided, we will ensure that it matches the value sent in `password`. |
| `--announcements-read` | bool | Signifies that the user has read all the announcements in the UI. |
| `--ai-assistant-personality-id` | int64 | AI Assistant Personality ID assigned directly to this user, if any. |
| `--allowed-ips` | string | A list of allowed IPs if applicable.  Newline delimited |
| `--attachments-permission` | bool | DEPRECATED: If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead. |
| `--authenticate-until` | datetime | Scheduled Date/Time at which user will be deactivated |
| `--authentication-method` | enum | How is this user authenticated?. One of: `password`, `sso`, `none`, `email_signup`, `password_with_imported_hash`, `password_and_ssh_key`. |
| `--billing-permission` | bool | Allow this user to perform operations on the account, payments, and invoices? |
| `--bypass-user-lifecycle-rules` | bool | Exempt this user from user lifecycle rules? |
| `--bypass-site-allowed-ips` | bool | Allow this user to skip site-wide IP blacklists? |
| `--dav-permission` | bool | Can the user connect with WebDAV? |
| `--desktop-configuration-profile-id` | int64 | Desktop Configuration Profile ID assigned directly to this user, if any. |
| `--default-workspace-id` | int64 | Workspace ID the user should land in by default when more than one Workspace is available. |
| `--disabled` | bool | Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date. |
| `--filesystem-layout` | enum | File system layout. One of: `site_root`, `user_root`, `partner_root`, `integration_centric`, `workspace_root`. |
| `--ftp-permission` | bool | Can the user access with FTP/FTPS? |
| `--header-text` | string | Text to display to the user in the header of the UI |
| `--language` | string | Preferred language |
| `--notification-daily-send-time` | int64 | Hour of the day at which daily notifications should be sent. Can be in range 0 to 23 |
| `--name` | string | User's full name |
| `--company` | string | User's company |
| `--notes` | string | Any internal notes on the user |
| `--office-integration-enabled` | bool | Enable integration with Office for the web? |
| `--partner-admin` | bool | Is this user a Partner administrator? |
| `--partner-id` | int64 | Partner ID if this user belongs to a Partner |
| `--password-validity-days` | int64 | Number of days to allow user to use the same password |
| `--primary-group-id` | int64 | Primary group ID for Group Admin scoping |
| `--readonly-site-admin` | bool | Is the user an allowed to view all (non-billing) site configuration for this site? |
| `--receive-admin-alerts` | bool | Deprecated. Use notify_on_all_site_warnings and granular failure notification preferences instead. |
| `--notify-on-all-site-warnings` | bool | Should the user receive site warnings via email? |
| `--notify-on-all-sso-failures` | bool | Should the user receive sso/scim/ldap configuration/sync failures via email? |
| `--notify-on-all-user-security-events` | bool | Should the user receive user security events via email? |
| `--notify-on-all-pending-work-failures` | bool | Should the user receive pending work failures via email? |
| `--notify-on-all-siem-http-destination-failures` | bool | Should the user receive siem failures via email? |
| `--notify-on-all-sync-failures` | bool | Should the user receive sync failures via email? |
| `--notify-on-all-automation-failures` | bool | Should the user receive automation failures via email? |
| `--notify-on-all-expectation-failures` | bool | Should the user receive expectation failures and misses via email? |
| `--require-login-by` | datetime | Require user to login by specified date otherwise it will be disabled. |
| `--require-password-change` | bool | Is a password change required upon next user login? |
| `--restapi-permission` | bool | Can this user access the Web app, Desktop app, SDKs, or REST API?  (All of these tools use the API internally, so this is one unified permission set.) |
| `--self-managed` | bool | Does this user manage it's own credentials or is it a shared/bot user? |
| `--sftp-permission` | bool | Can the user access with SFTP? |
| `--site-admin` | bool | Is the user an administrator for this site? |
| `--skip-welcome-screen` | bool | Skip Welcome page in the UI? |
| `--ssl-required` | enum | SSL required setting. One of: `use_system_setting`, `always_require`, `never_require`. |
| `--sso-strategy-id` | int64 | SSO (Single Sign On) strategy ID for the user, if applicable. |
| `--subscribe-to-newsletter` | bool | Is the user subscribed to the newsletter? |
| `--require-2fa` | enum | 2FA required setting. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed. One of: `use_system_setting`, `always_require`, `never_require`. |
| `--tags` | string | Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens. |
| `--time-zone` | string | User time zone |
| `--user-root` | string | Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set).  Note that this is not used for API, Desktop, or Web interface. |
| `--user-home` | string | Home folder for FTP/SFTP.  Note that this is not used for API, Desktop, or Web interface. |
| `--workspace-admin` | bool | Is the user a Workspace administrator?  Applicable only to the workspace ID related to this user, if one is set. |
| `--username` | string | User's username **Required.** |
| `--workspace-id` | int64 | Workspace ID |

### `files-cli users unlock`

Unlock user who has been locked out due to failed logins.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User ID. **Required.** |

### `files-cli users resend-welcome-email`

Resend user welcome email.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User ID. **Required.** |

### `files-cli users user-2fa-reset`

Trigger 2FA Reset process for user who has lost access to their existing 2FA methods.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User ID. **Required.** |

### `files-cli users update`

Update User.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User ID. **Required.** |
| `--avatar-file` | file | An image file for your user avatar. |
| `--avatar-delete` | bool | If true, the avatar will be deleted. |
| `--change-password` | string | Used for changing a password on an existing user. |
| `--change-password-confirmation` | string | Optional, but if provided, we will ensure that it matches the value sent in `change_password`. |
| `--email` | string | User's email. |
| `--grant-permission` | string | Permission to grant on the User Root upon user creation. Can be blank or `full`, `read`, `write`, `list`, `read+write`, or `list+write` |
| `--group-id` | int64 | Group ID to associate this user with. |
| `--group-ids` | string | A list of group ids to associate this user with.  Comma delimited. |
| `--imported-password-hash` | string | Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256. |
| `--password` | string | User password. |
| `--password-confirmation` | string | Optional, but if provided, we will ensure that it matches the value sent in `password`. |
| `--announcements-read` | bool | Signifies that the user has read all the announcements in the UI. |
| `--ai-assistant-personality-id` | int64 | AI Assistant Personality ID assigned directly to this user, if any. |
| `--allowed-ips` | string | A list of allowed IPs if applicable.  Newline delimited |
| `--attachments-permission` | bool | DEPRECATED: If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead. |
| `--authenticate-until` | datetime | Scheduled Date/Time at which user will be deactivated |
| `--authentication-method` | enum | How is this user authenticated?. One of: `password`, `sso`, `none`, `email_signup`, `password_with_imported_hash`, `password_and_ssh_key`. |
| `--billing-permission` | bool | Allow this user to perform operations on the account, payments, and invoices? |
| `--bypass-user-lifecycle-rules` | bool | Exempt this user from user lifecycle rules? |
| `--bypass-site-allowed-ips` | bool | Allow this user to skip site-wide IP blacklists? |
| `--dav-permission` | bool | Can the user connect with WebDAV? |
| `--desktop-configuration-profile-id` | int64 | Desktop Configuration Profile ID assigned directly to this user, if any. |
| `--default-workspace-id` | int64 | Workspace ID the user should land in by default when more than one Workspace is available. |
| `--disabled` | bool | Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date. |
| `--filesystem-layout` | enum | File system layout. One of: `site_root`, `user_root`, `partner_root`, `integration_centric`, `workspace_root`. |
| `--ftp-permission` | bool | Can the user access with FTP/FTPS? |
| `--header-text` | string | Text to display to the user in the header of the UI |
| `--language` | string | Preferred language |
| `--notification-daily-send-time` | int64 | Hour of the day at which daily notifications should be sent. Can be in range 0 to 23 |
| `--name` | string | User's full name |
| `--company` | string | User's company |
| `--notes` | string | Any internal notes on the user |
| `--office-integration-enabled` | bool | Enable integration with Office for the web? |
| `--partner-admin` | bool | Is this user a Partner administrator? |
| `--partner-id` | int64 | Partner ID if this user belongs to a Partner |
| `--password-validity-days` | int64 | Number of days to allow user to use the same password |
| `--primary-group-id` | int64 | Primary group ID for Group Admin scoping |
| `--readonly-site-admin` | bool | Is the user an allowed to view all (non-billing) site configuration for this site? |
| `--receive-admin-alerts` | bool | Deprecated. Use notify_on_all_site_warnings and granular failure notification preferences instead. |
| `--notify-on-all-site-warnings` | bool | Should the user receive site warnings via email? |
| `--notify-on-all-sso-failures` | bool | Should the user receive sso/scim/ldap configuration/sync failures via email? |
| `--notify-on-all-user-security-events` | bool | Should the user receive user security events via email? |
| `--notify-on-all-pending-work-failures` | bool | Should the user receive pending work failures via email? |
| `--notify-on-all-siem-http-destination-failures` | bool | Should the user receive siem failures via email? |
| `--notify-on-all-sync-failures` | bool | Should the user receive sync failures via email? |
| `--notify-on-all-automation-failures` | bool | Should the user receive automation failures via email? |
| `--notify-on-all-expectation-failures` | bool | Should the user receive expectation failures and misses via email? |
| `--require-login-by` | datetime | Require user to login by specified date otherwise it will be disabled. |
| `--require-password-change` | bool | Is a password change required upon next user login? |
| `--restapi-permission` | bool | Can this user access the Web app, Desktop app, SDKs, or REST API?  (All of these tools use the API internally, so this is one unified permission set.) |
| `--self-managed` | bool | Does this user manage it's own credentials or is it a shared/bot user? |
| `--sftp-permission` | bool | Can the user access with SFTP? |
| `--site-admin` | bool | Is the user an administrator for this site? |
| `--skip-welcome-screen` | bool | Skip Welcome page in the UI? |
| `--ssl-required` | enum | SSL required setting. One of: `use_system_setting`, `always_require`, `never_require`. |
| `--sso-strategy-id` | int64 | SSO (Single Sign On) strategy ID for the user, if applicable. |
| `--subscribe-to-newsletter` | bool | Is the user subscribed to the newsletter? |
| `--require-2fa` | enum | 2FA required setting. `use_system_setting` uses the site-wide setting, including SSO exemptions. `always_require` and `never_require` override the site-wide setting when user-level overrides are allowed. One of: `use_system_setting`, `always_require`, `never_require`. |
| `--tags` | string | Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens. |
| `--time-zone` | string | User time zone |
| `--user-root` | string | Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set).  Note that this is not used for API, Desktop, or Web interface. |
| `--user-home` | string | Home folder for FTP/SFTP.  Note that this is not used for API, Desktop, or Web interface. |
| `--workspace-admin` | bool | Is the user a Workspace administrator?  Applicable only to the workspace ID related to this user, if one is set. |
| `--username` | string | User's username |
| `--workspace-id` | int64 | Workspace ID |
| `--clear-2fa` | bool | If true when changing authentication_method from `password` to `sso`, remove all two-factor methods. Ignored in all other cases. |
| `--convert-to-partner-user` | bool | If true, convert this user to a partner user by assigning the partner_id provided. |

### `files-cli users delete`

Delete User.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | User ID. **Required.** |
| `--new-owner-id` | int64 | Provide a User ID here to transfer ownership of certain resources such as Automations and Share Links (Bundles) to that new user. |

## Limitations and considerations

**Authentication method determines how the user signs in.** The `--authentication-method` flag accepts these values:

- `password` — site administrator specifies the password; the user can change it later. Pair with `--password` on create.
- `email_signup` — Files.com sends a Welcome Email with a link for the user to set their own password. The method automatically changes to `password` after they do.
- `sso` — the user authenticates through a configured SSO provider. Pair with `--sso-strategy-id`. Users can still authenticate with an API key or SSH/SFTP key.
- `none` — no password login at all. API key and SSH/SFTP key only. Typically used for service accounts.
- `password_and_ssh_key` — SFTP only. Both a password and an SSH key are required to authenticate over SFTP. API keys still work as an alternative.
- `password_with_imported_hash` — for migration from another system. Pair with `--imported-password-hash`. After the first successful login the system re-stores the password in its internal format and the method changes to `password`.

**Do not apply User Lifecycle Rules to SSO or SCIM users.** The identity provider controls the lifecycle state of these users. Use `--bypass-user-lifecycle-rules` to exempt a specific user.

**Disable vs delete.** Setting `--disabled=true` blocks login but preserves the user record. Disabled users cannot log in and do not count for billing purposes. Deletion is a separate, permanent operation.

