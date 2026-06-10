---
name: filescom-sites
description: |
  A Site is the place you'll come to update site settings, as well as manage site-wide API keys.
---

# filescom-sites

A Site is the place you'll come to update site settings, as well as manage site-wide API keys.

Most site settings can be set via the API.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli sites get`

Show Site Settings.

No flags beyond the global ones.

### `files-cli sites get-usage`

Get the most recent usage snapshot (usage data for billing purposes) for a Site.

No flags beyond the global ones.

### `files-cli sites update`

Update Site Settings.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Site name |
| `--subdomain` | string | Site subdomain |
| `--domain` | string | Custom domain |
| `--domain-hsts-header` | bool | Send HSTS (HTTP Strict Transport Security) header when visitors access the site via a custom domain? |
| `--domain-letsencrypt-chain` | string | Letsencrypt chain to use when registering SSL Certificate for domain. |
| `--email` | string | Main email for this site |
| `--reply-to-email` | string | Reply-to email for this site |
| `--allow-bundle-names` | bool | Are manual Bundle names allowed? |
| `--bundle-expiration` | int64 | Site-wide Bundle expiration in days |
| `--welcome-email-enabled` | bool | Will the welcome email be sent to new users? |
| `--ask-about-overwrites` | bool | If false, rename conflicting files instead of asking for overwrite confirmation.  Only applies to web interface. |
| `--show-request-access-link` | bool | Show request access link for users without access?  Currently unused. |
| `--always-mkdir-parents` | bool | Create parent directories if they do not exist during uploads?  This is primarily used to work around broken upload clients that assume servers will perform this step. |
| `--welcome-email-cc` | string | Include this email in welcome emails if enabled |
| `--welcome-email-subject` | string | Include this email subject in welcome emails if enabled |
| `--welcome-custom-text` | string | Custom text send in user welcome email |
| `--language` | string | Site default language |
| `--windows-mode-ftp` | bool | Does FTP user Windows emulation mode? |
| `--default-time-zone` | string | Site default time zone |
| `--desktop-app` | bool | Is the desktop app enabled? |
| `--desktop-app-session-ip-pinning` | bool | Is desktop app session IP pinning enabled? |
| `--desktop-app-session-lifetime` | int64 | Desktop app session lifetime (in hours) |
| `--mobile-app` | bool | Is the mobile app enabled? |
| `--mobile-app-session-ip-pinning` | bool | Is mobile app session IP pinning enabled? |
| `--mobile-app-session-lifetime` | int64 | Mobile app session lifetime (in hours) |
| `--folder-permissions-groups-only` | bool | If true, permissions for this site must be bound to a group (not a user). |
| `--welcome-screen` | string | Does the welcome screen appear? |
| `--office-integration-available` | bool | If true, allows users to use a document editing integration. |
| `--office-integration-type` | string | Which document editing integration to support. Files.com Editor or Microsoft Office for the Web. |
| `--pin-all-remote-servers-to-site-region` | bool | If true, we will ensure that all internal communications with any remote server are made through the primary region of the site. This setting overrides individual remote server settings. |
| `--motd-text` | string | A message to show users when they connect via FTP or SFTP. |
| `--motd-use-for-ftp` | bool | Show message to users connecting via FTP |
| `--motd-use-for-sftp` | bool | Show message to users connecting via SFTP |
| `--left-navigation-visibility` | object | Visibility settings for account navigation |
| `--disable-all-ai-features` | bool | If true, all AI features are disabled for this site. |
| `--ai-feature-availability` | object | Availability settings for AI features by user class |
| `--additional-text-file-types` | []string | Additional extensions that are considered text files |
| `--bundle-require-note` | bool | Do Bundles require internal notes? |
| `--bundle-send-shared-receipts` | bool | Do Bundle creators receive receipts of invitations? |
| `--bundles-default-owned-by-primary-group` | bool | If true, new Share Links created by a user with a primary group will default to that group as owner. |
| `--calculate-file-checksums-crc32` | bool | Calculate CRC32 checksums for files? |
| `--calculate-file-checksums-md5` | bool | Calculate MD5 checksums for files? |
| `--calculate-file-checksums-sha1` | bool | Calculate SHA1 checksums for files? |
| `--calculate-file-checksums-sha256` | bool | Calculate SHA256 checksums for files? |
| `--legacy-checksums-mode` | bool | Use legacy checksums mode? |
| `--migrate-remote-server-sync-to-sync` | bool | If true, we will migrate all remote server syncs to the new Sync model. |
| `--as2-message-retention-days` | int64 | Number of days to retain AS2 messages (incoming and outgoing). |
| `--username-display` | string | How usernames are displayed in the web UI. Can be `username_only`, `full_name_only`, `full_name_username`, `full_name_company`, or `full_name_username_company`. |
| `--session-expiry-minutes` | int64 | Session expiry in minutes |
| `--ssl-required` | bool | Is SSL required?  Disabling this is insecure. |
| `--sftp-insecure-ciphers` | bool | If true, we will allow weak and known insecure ciphers to be used for SFTP connections.  Enabling this setting severely weakens the security of your site and it is not recommend, except as a last resort for compatibility. |
| `--sftp-insecure-diffie-hellman` | bool | If true, we will allow weak Diffie Hellman parameters to be used within ciphers for SFTP that are otherwise on our secure list.  This has the effect of making the cipher weaker than our normal threshold for security, but is required to support certain legacy or broken SSH and MFT clients.  Enabling this weakens security, but not nearly as much as enabling the full `sftp_insecure_ciphers` option. |
| `--disable-files-certificate-generation` | bool | If set, Files.com will not set the CAA records required to generate future SSL certificates for this domain. |
| `--user-lockout` | bool | Will users be locked out after incorrect login attempts? |
| `--user-lockout-tries` | int64 | Number of login tries within `user_lockout_within` hours before users are locked out |
| `--user-lockout-within` | int64 | Number of hours for user lockout window |
| `--user-lockout-lock-period` | int64 | How many hours to lock user out for failed password? |
| `--include-password-in-welcome-email` | bool | Include password in emails to new users? |
| `--allowed-countries` | string | Comma separated list of allowed Country codes |
| `--allowed-ips` | string | List of allowed IP addresses |
| `--allow-user-level-2fa-override` | bool | Allow the site-wide two-factor authentication requirement to be overriden on a per-user-basis? |
| `--allow-user-level-allowed-ip-override` | bool | Allow the site-wide allowed IP restriction to be overriden on a per-user-basis? |
| `--allow-user-level-ssl-override` | bool | Allow the site-wide FTP SSL requirement to be overriden on a per-user-basis? |
| `--disallowed-countries` | string | Comma separated list of disallowed Country codes |
| `--days-to-retain-backups` | int64 | Number of days to keep deleted files |
| `--max-prior-passwords` | int64 | Number of prior passwords to disallow |
| `--password-validity-days` | int64 | Number of days password is valid |
| `--password-min-length` | int64 | Shortest password length for users |
| `--password-require-letter` | bool | Require a letter in passwords? |
| `--password-require-mixed` | bool | Require lower and upper case letters in passwords? |
| `--password-require-special` | bool | Require special characters in password? |
| `--password-require-number` | bool | Require a number in passwords? |
| `--password-require-unbreached` | bool | Require passwords that have not been previously breached? (see https://haveibeenpwned.com/) |
| `--require-logout-from-bundles-and-inboxes` | bool | If true, we will hide the 'Remember Me' box on Inbox and Bundle registration pages, requiring that the user logout and log back in every time they visit the page. |
| `--dav-user-root-enabled` | bool | Use user FTP roots also for WebDAV? |
| `--sftp-user-root-enabled` | bool | Use user FTP roots also for SFTP? |
| `--disable-password-reset` | bool | Is password reset disabled? |
| `--immutable-files` | bool | Are files protected from modification? |
| `--bundle-not-found-message` | string | Custom error message to show when bundle is not found. |
| `--bundle-password-required` | bool | Do Bundles require password protection? |
| `--bundle-require-registration` | bool | Do Bundles require registration? |
| `--bundle-require-share-recipient` | bool | Do Bundles require recipients for sharing? |
| `--bundle-registration-notifications` | string | Do Bundle owners receive registration notification? |
| `--bundle-activity-notifications` | string | Do Bundle owners receive activity notifications? |
| `--bundle-upload-receipt-notifications` | string | Do Bundle uploaders receive upload confirmation notifications? |
| `--document-edits-in-bundle-allowed` | bool | If true, allow public viewers of Bundles with full permissions to use document editing integrations. |
| `--password-requirements-apply-to-bundles` | bool | Require bundles' passwords, and passwords for other items (inboxes, public shares, etc.) to conform to the same requirements as users' passwords? |
| `--prevent-root-permissions-for-non-site-admins` | bool | If true, we will prevent non-administrators from receiving any permissions directly on the root folder.  This is commonly used to prevent the accidental application of permissions. |
| `--opt-out-global` | bool | Use servers in the USA only? |
| `--use-provided-modified-at` | bool | Allow uploaders to set `provided_modified_at` for uploaded files? |
| `--custom-namespace` | bool | Is this site using a custom namespace for users? |
| `--non-sso-groups-allowed` | bool | If true, groups can be manually created / modified / deleted by Site Admins. Otherwise, groups can only be managed via your SSO provider. |
| `--non-sso-users-allowed` | bool | If true, users can be manually created / modified / deleted by Site Admins. Otherwise, users can only be managed via your SSO provider. |
| `--sharing-enabled` | bool | Allow bundle creation |
| `--snapshot-sharing-enabled` | bool | Allow snapshot share links creation |
| `--user-requests-enabled` | bool | Enable User Requests feature |
| `--user-requests-notify-admins` | bool | Send email to site admins when a user request is received? |
| `--dav-enabled` | bool | Is WebDAV enabled? |
| `--ftp-enabled` | bool | Is FTP enabled? |
| `--sftp-enabled` | bool | Is SFTP enabled? |
| `--sftp-finalize-partial-uploads` | bool | Finalize partial SFTP uploads from interrupted connections? Default: true. |
| `--users-can-create-api-keys` | bool | Allow users to create their own API keys? |
| `--users-can-create-ssh-keys` | bool | Allow users to create their own SSH keys? |
| `--show-user-notifications-log-in-link` | bool | Show log in link in user notifications? |
| `--sftp-host-key-type` | string | Sftp Host Key Type |
| `--active-sftp-host-key-id` | int64 | Id of the currently selected custom SFTP Host Key |
| `--protocol-access-groups-only` | bool | If true, protocol access permissions on users will be ignored, and only protocol access permissions set on Groups will be honored.  Make sure that your current user is a member of a group with API permission when changing this value to avoid locking yourself out of your site. |
| `--revoke-bundle-access-on-disable-or-delete` | bool | Auto-removes bundles for disabled/deleted users and enforces bundle expiry within user access period. |
| `--bundle-watermark-value` | object | Preview watermark settings applied to all bundle items. Uses the same keys as Behavior.value |
| `--group-admins-can-add-users` | bool | Allow group admins to create users in their groups |
| `--group-admins-can-delete-users` | bool | Allow group admins to delete users in their groups |
| `--group-admins-can-enable-disable-users` | bool | Allow group admins to enable or disable users in their groups |
| `--group-admins-can-modify-users` | bool | Allow group admins to modify users in their groups |
| `--group-admins-can-bypass-user-lifecycle-rules` | bool | Allow group admins to exempt users in their groups from lifecycle rules |
| `--group-admins-can-reset-passwords` | bool | Allow group admins to reset passwords for users in their groups |
| `--group-admins-can-set-user-password` | bool | Allow group admins to set password authentication method |
| `--bundle-recipient-blacklist-free-email-domains` | bool | Disallow free email domains for Bundle/Inbox recipients? |
| `--bundle-recipient-blacklist-domains` | []string | List of email domains to disallow when entering a Bundle/Inbox recipients |
| `--admins-bypass-locked-subfolders` | bool | Allow admins to bypass the locked subfolders setting. |
| `--allowed-2fa-method-sms` | bool | Is SMS two factor authentication allowed? |
| `--allowed-2fa-method-totp` | bool | Is TOTP two factor authentication allowed? |
| `--allowed-2fa-method-webauthn` | bool | Is WebAuthn two factor authentication allowed? |
| `--allowed-2fa-method-yubi` | bool | Is yubikey two factor authentication allowed? |
| `--allowed-2fa-method-email` | bool | Is OTP via email two factor authentication allowed? |
| `--allowed-2fa-method-static` | bool | Is OTP via static codes for two factor authentication allowed? |
| `--allowed-2fa-method-bypass-for-ftp-sftp-dav` | bool | Are users allowed to configure their two factor authentication to be bypassed for FTP/SFTP/WebDAV? |
| `--require-2fa` | bool | Require two-factor authentication for all users? |
| `--require-2fa-exempt-all-sso-users` | bool | If true, SSO users using the default user-level two-factor authentication setting are exempt from the site-wide two-factor authentication requirement. |
| `--require-2fa-user-type` | string | What type of user is required to use two-factor authentication (when require_2fa is set to `true` for this site)? |
| `--color2-top` | string | Top bar background color |
| `--color2-left` | string | Page link and button color |
| `--color2-link` | string | Top bar link color |
| `--color2-text` | string | Page link and button color |
| `--color2-top-text` | string | Top bar text color |
| `--site-header` | string | Custom site header text for authenticated pages |
| `--site-footer` | string | Custom site footer text for authenticated pages |
| `--site-public-header` | string | Custom site header text for public pages |
| `--site-public-footer` | string | Custom site footer text for public pages |
| `--login-help-text` | string | Login help text |
| `--use-dedicated-ips-for-smtp` | bool | If using custom SMTP, should we use dedicated IPs to deliver emails? |
| `--email-footer-custom-text` | string | Custom footer text for system-generated emails. Supports standard strftime date/time patterns like %Y (4-digit year), %m (month), %d (day). |
| `--smtp-address` | string | SMTP server hostname or IP |
| `--smtp-authentication` | string | SMTP server authentication type |
| `--smtp-from` | string | From address to use when mailing through custom SMTP |
| `--smtp-username` | string | SMTP server username |
| `--smtp-port` | int64 | SMTP server port |
| `--ldap-enabled` | bool | Main LDAP setting: is LDAP enabled? |
| `--ldap-type` | string | LDAP type |
| `--ldap-host` | string | LDAP host |
| `--ldap-host-2` | string | LDAP backup host |
| `--ldap-host-3` | string | LDAP backup host |
| `--ldap-port` | int64 | LDAP port |
| `--ldap-secure` | bool | Use secure LDAP? |
| `--ldap-username` | string | Username for signing in to LDAP server. |
| `--ldap-username-field` | string | LDAP username field |
| `--ldap-domain` | string | Domain name that will be appended to usernames |
| `--ldap-user-action` | string | Should we sync users from LDAP server? |
| `--ldap-group-action` | string | Should we sync groups from LDAP server? |
| `--ldap-user-include-groups` | string | Comma or newline separated list of group names (with optional wildcards) - if provided, only users in these groups will be added or synced. |
| `--ldap-group-exclusion` | string | Comma or newline separated list of group names (with optional wildcards) to exclude when syncing. |
| `--ldap-group-inclusion` | string | Comma or newline separated list of group names (with optional wildcards) to include when syncing. |
| `--ldap-base-dn` | string | Base DN for looking up users in LDAP server |
| `--uploads-via-email-authentication` | bool | Do incoming emails in the Inboxes require checking for SPF/DKIM/DMARC? |
| `--icon16-file` | file | (no description) |
| `--icon16-delete` | bool | If true, will delete the file stored in icon16 |
| `--icon32-file` | file | (no description) |
| `--icon32-delete` | bool | If true, will delete the file stored in icon32 |
| `--icon48-file` | file | (no description) |
| `--icon48-delete` | bool | If true, will delete the file stored in icon48 |
| `--icon128-file` | file | (no description) |
| `--icon128-delete` | bool | If true, will delete the file stored in icon128 |
| `--logo-file` | file | (no description) |
| `--logo-delete` | bool | If true, will delete the file stored in logo |
| `--bundle-watermark-attachment-file` | file | (no description) |
| `--bundle-watermark-attachment-delete` | bool | If true, will delete the file stored in bundle_watermark_attachment |
| `--login-page-background-image-file` | file | (no description) |
| `--login-page-background-image-delete` | bool | If true, will delete the file stored in login_page_background_image |
| `--disable-2fa-with-delay` | bool | If set to true, we will begin the process of disabling 2FA on this site. |
| `--ldap-password-change` | string | New LDAP password. |
| `--ldap-password-change-confirmation` | string | Confirm new LDAP password. |
| `--redirect-old-subdomain` | bool | If true, and if changing the site subdomain, then create a redirect from the previous Files.com subdomain to the new Files.com subdomain. |
| `--smtp-password` | string | Password for SMTP server. |

