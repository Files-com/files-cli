# files-cli Skills — Index

Domain skills mirror the CLI's top-level command set, one per domain, grouped by docs category.

See [`README.md`](README.md) for an overview and [`../CONTEXT.md`](../CONTEXT.md) for CLI-wide conventions.

## Recipes

Cross-cutting workflows that span multiple domains.

| Skill | Description |
| --- | --- |
| [`recipe-folder-size-and-counts`](recipes/recipe-folder-size-and-counts/SKILL.md) | Determine storage usage and file/folder counts on Files.com from the data the agent can actually read back, without recursively walking the tree; use this when the user asks "how big is this folder", "how much storage are we using", or "how many files do we have". |
| [`recipe-generating-reports`](recipes/recipe-generating-reports/SKILL.md) | Choose the right built-in Files.com report or export instead of assembling one by hand; use this when the user asks for a report, audit, export, or summary — storage/usage, bandwidth/transfer, activity/audit logs, webhook delivery, permissions, group membership, or share-link audits. |
| [`recipe-searching-for-files`](recipes/recipe-searching-for-files/SKILL.md) | Find files and folders on Files.com effectively, and pick the right tool for the job; use this when the user wants to locate files by name, path, modified time, or custom metadata. |
| [`recipe-share-and-notify`](recipes/recipe-share-and-notify/SKILL.md) | Share a Files.com file or folder via a Share Link (bundle) and email it to recipients. |

## Domains

One skill per top-level CLI command.

### Automations

| Skill | Resource |
| --- | --- |
| [`filescom-ai-tasks`](filescom-ai-tasks/SKILL.md) | Ai Tasks |
| [`filescom-automation-runs`](filescom-automation-runs/SKILL.md) | Automation Runs |
| [`filescom-automations`](filescom-automations/SKILL.md) | Automations |
| [`filescom-expectation-evaluations`](filescom-expectation-evaluations/SKILL.md) | Expectation Evaluations |
| [`filescom-expectation-incidents`](filescom-expectation-incidents/SKILL.md) | Expectation Incidents |
| [`filescom-expectations`](filescom-expectations/SKILL.md) | Expectations |

### Billing

| Skill | Resource |
| --- | --- |
| [`filescom-invoices`](filescom-invoices/SKILL.md) | Invoices |
| [`filescom-payments`](filescom-payments/SKILL.md) | Payments |

### Developers

| Skill | Resource |
| --- | --- |
| [`filescom-api-keys`](filescom-api-keys/SKILL.md) | API Keys |
| [`filescom-apps`](filescom-apps/SKILL.md) | Apps |
| [`filescom-sessions`](filescom-sessions/SKILL.md) | Sessions |

### Encryption

| Skill | Resource |
| --- | --- |
| [`filescom-gpg-keys`](filescom-gpg-keys/SKILL.md) | GPG Keys |
| [`filescom-public-keys`](filescom-public-keys/SKILL.md) | Public Keys |
| [`filescom-sftp-host-keys`](filescom-sftp-host-keys/SKILL.md) | SFTP Host Keys |

### File System

| Skill | Resource |
| --- | --- |
| [`filescom-file-comment-reactions`](filescom-file-comment-reactions/SKILL.md) | File Comment Reactions |
| [`filescom-file-comments`](filescom-file-comments/SKILL.md) | File Comments |
| [`filescom-file-migrations`](filescom-file-migrations/SKILL.md) | File Migrations |
| [`filescom-files`](filescom-files/SKILL.md) | Files |
| [`filescom-folders`](filescom-folders/SKILL.md) | Folders |
| [`filescom-locks`](filescom-locks/SKILL.md) | Locks |
| [`filescom-message-comment-reactions`](filescom-message-comment-reactions/SKILL.md) | Message Comment Reactions |
| [`filescom-message-comments`](filescom-message-comments/SKILL.md) | Message Comments |
| [`filescom-message-reactions`](filescom-message-reactions/SKILL.md) | Message Reactions |
| [`filescom-messages`](filescom-messages/SKILL.md) | Messages |
| [`filescom-projects`](filescom-projects/SKILL.md) | Projects |
| [`filescom-requests`](filescom-requests/SKILL.md) | Requests |
| [`filescom-restores`](filescom-restores/SKILL.md) | Restores |
| [`filescom-snapshots`](filescom-snapshots/SKILL.md) | Snapshots |

### Integrations

| Skill | Resource |
| --- | --- |
| [`filescom-remote-mount-backends`](filescom-remote-mount-backends/SKILL.md) | Remote Mount Backends |
| [`filescom-remote-server-credentials`](filescom-remote-server-credentials/SKILL.md) | Remote Server Credentials |
| [`filescom-remote-servers`](filescom-remote-servers/SKILL.md) | Remote Servers |
| [`filescom-secrets`](filescom-secrets/SKILL.md) | Secrets |
| [`filescom-sync-runs`](filescom-sync-runs/SKILL.md) | Sync Runs |
| [`filescom-syncs`](filescom-syncs/SKILL.md) | Syncs |

### Logging

| Skill | Resource |
| --- | --- |
| [`filescom-action-logs`](filescom-action-logs/SKILL.md) | Action Logs |
| [`filescom-action-notification-export-results`](filescom-action-notification-export-results/SKILL.md) | Action Notification Export Results |
| [`filescom-action-notification-exports`](filescom-action-notification-exports/SKILL.md) | Action Notification Exports |
| [`filescom-api-request-logs`](filescom-api-request-logs/SKILL.md) | API Request Logs |
| [`filescom-automation-logs`](filescom-automation-logs/SKILL.md) | Automation Logs |
| [`filescom-chat-sessions`](filescom-chat-sessions/SKILL.md) | Chat Sessions |
| [`filescom-email-logs`](filescom-email-logs/SKILL.md) | Email Logs |
| [`filescom-exavault-api-request-logs`](filescom-exavault-api-request-logs/SKILL.md) | Exavault API Request Logs |
| [`filescom-external-events`](filescom-external-events/SKILL.md) | External Events |
| [`filescom-file-migration-logs`](filescom-file-migration-logs/SKILL.md) | File Migration Logs |
| [`filescom-ftp-action-logs`](filescom-ftp-action-logs/SKILL.md) | FTP Action Logs |
| [`filescom-histories`](filescom-histories/SKILL.md) | Histories |
| [`filescom-history-export-results`](filescom-history-export-results/SKILL.md) | History Export Results |
| [`filescom-history-exports`](filescom-history-exports/SKILL.md) | History Exports |
| [`filescom-inbound-s3-logs`](filescom-inbound-s3-logs/SKILL.md) | Inbound S3 Logs |
| [`filescom-outbound-connection-logs`](filescom-outbound-connection-logs/SKILL.md) | Outbound Connection Logs |
| [`filescom-pending-work-events`](filescom-pending-work-events/SKILL.md) | Pending Work Events |
| [`filescom-public-hosting-request-logs`](filescom-public-hosting-request-logs/SKILL.md) | Public Hosting Request Logs |
| [`filescom-scheduled-exports`](filescom-scheduled-exports/SKILL.md) | Scheduled Exports |
| [`filescom-scim-logs`](filescom-scim-logs/SKILL.md) | Scim Logs |
| [`filescom-settings-changes`](filescom-settings-changes/SKILL.md) | Settings Changes |
| [`filescom-sftp-action-logs`](filescom-sftp-action-logs/SKILL.md) | SFTP Action Logs |
| [`filescom-siem-http-destination-events`](filescom-siem-http-destination-events/SKILL.md) | SIEM HTTP Destination Events |
| [`filescom-siem-http-destinations`](filescom-siem-http-destinations/SKILL.md) | SIEM HTTP Destinations |
| [`filescom-sso-events`](filescom-sso-events/SKILL.md) | SSO Events |
| [`filescom-sync-logs`](filescom-sync-logs/SKILL.md) | Sync Logs |
| [`filescom-user-security-events`](filescom-user-security-events/SKILL.md) | User Security Events |
| [`filescom-web-dav-action-logs`](filescom-web-dav-action-logs/SKILL.md) | WebDAV Action Logs |

### Notifications

| Skill | Resource |
| --- | --- |
| [`filescom-event-channels`](filescom-event-channels/SKILL.md) | Event Channels |
| [`filescom-event-delivery-attempts`](filescom-event-delivery-attempts/SKILL.md) | Event Delivery Attempts |
| [`filescom-event-records`](filescom-event-records/SKILL.md) | Event Records |
| [`filescom-event-subscriptions`](filescom-event-subscriptions/SKILL.md) | Event Subscriptions |
| [`filescom-event-targets`](filescom-event-targets/SKILL.md) | Event Targets |
| [`filescom-notifications`](filescom-notifications/SKILL.md) | Notifications |
| [`filescom-webhook-tests`](filescom-webhook-tests/SKILL.md) | Webhook Tests |

### Security

| Skill | Resource |
| --- | --- |
| [`filescom-dns-records`](filescom-dns-records/SKILL.md) | DNS Records |
| [`filescom-ip-addresses`](filescom-ip-addresses/SKILL.md) | IP Addresses |
| [`filescom-key-lifecycle-rules`](filescom-key-lifecycle-rules/SKILL.md) | Key Lifecycle Rules |

### Services

| Skill | Resource |
| --- | --- |
| [`filescom-as2-incoming-messages`](filescom-as2-incoming-messages/SKILL.md) | AS2 Incoming Messages |
| [`filescom-as2-outgoing-messages`](filescom-as2-outgoing-messages/SKILL.md) | AS2 Outgoing Messages |
| [`filescom-as2-partners`](filescom-as2-partners/SKILL.md) | AS2 Partners |
| [`filescom-as2-stations`](filescom-as2-stations/SKILL.md) | AS2 Stations |

### Settings

| Skill | Resource |
| --- | --- |
| [`filescom-ai-assistant-personalities`](filescom-ai-assistant-personalities/SKILL.md) | Ai Assistant Personalities |
| [`filescom-behaviors`](filescom-behaviors/SKILL.md) | Behaviors |
| [`filescom-child-site-management-policies`](filescom-child-site-management-policies/SKILL.md) | Child Site Management Policies |
| [`filescom-custom-domains`](filescom-custom-domains/SKILL.md) | Custom Domains |
| [`filescom-desktop-configuration-profiles`](filescom-desktop-configuration-profiles/SKILL.md) | Desktop Configuration Profiles |
| [`filescom-holiday-regions`](filescom-holiday-regions/SKILL.md) | Holiday Regions |
| [`filescom-integration-centric-profiles`](filescom-integration-centric-profiles/SKILL.md) | Integration Centric Profiles |
| [`filescom-metadata-categories`](filescom-metadata-categories/SKILL.md) | Metadata Categories |
| [`filescom-site-subdomain-redirects`](filescom-site-subdomain-redirects/SKILL.md) | Site Subdomain Redirects |
| [`filescom-sites`](filescom-sites/SKILL.md) | Sites |
| [`filescom-styles`](filescom-styles/SKILL.md) | Styles |
| [`filescom-workspaces`](filescom-workspaces/SKILL.md) | Workspaces |

### Sharing

| Skill | Resource |
| --- | --- |
| [`filescom-bundle-actions`](filescom-bundle-actions/SKILL.md) | Bundle Actions |
| [`filescom-bundle-downloads`](filescom-bundle-downloads/SKILL.md) | Bundle Downloads |
| [`filescom-bundle-notifications`](filescom-bundle-notifications/SKILL.md) | Bundle Notifications |
| [`filescom-bundle-recipients`](filescom-bundle-recipients/SKILL.md) | Bundle Recipients |
| [`filescom-bundle-registrations`](filescom-bundle-registrations/SKILL.md) | Bundle Registrations |
| [`filescom-bundles`](filescom-bundles/SKILL.md) | Bundles |
| [`filescom-clickwraps`](filescom-clickwraps/SKILL.md) | Clickwraps |
| [`filescom-email-incoming-messages`](filescom-email-incoming-messages/SKILL.md) | Email Incoming Messages |
| [`filescom-form-field-sets`](filescom-form-field-sets/SKILL.md) | Form Field Sets |
| [`filescom-inbox-recipients`](filescom-inbox-recipients/SKILL.md) | Inbox Recipients |
| [`filescom-inbox-registrations`](filescom-inbox-registrations/SKILL.md) | Inbox Registrations |
| [`filescom-inbox-uploads`](filescom-inbox-uploads/SKILL.md) | Inbox Uploads |
| [`filescom-share-groups`](filescom-share-groups/SKILL.md) | Share Groups |

### Usage

| Skill | Resource |
| --- | --- |
| [`filescom-bandwidth-snapshots`](filescom-bandwidth-snapshots/SKILL.md) | Bandwidth Snapshots |
| [`filescom-remote-bandwidth-snapshots`](filescom-remote-bandwidth-snapshots/SKILL.md) | Remote Bandwidth Snapshots |
| [`filescom-usage-daily-snapshots`](filescom-usage-daily-snapshots/SKILL.md) | Usage Daily Snapshots |
| [`filescom-usage-snapshots`](filescom-usage-snapshots/SKILL.md) | Usage Snapshots |

### User Accounts

| Skill | Resource |
| --- | --- |
| [`filescom-group-users`](filescom-group-users/SKILL.md) | Group Users |
| [`filescom-groups`](filescom-groups/SKILL.md) | Groups |
| [`filescom-partner-channel-templates`](filescom-partner-channel-templates/SKILL.md) | Partner Channel Templates |
| [`filescom-partner-channels`](filescom-partner-channels/SKILL.md) | Partner Channels |
| [`filescom-partner-site-requests`](filescom-partner-site-requests/SKILL.md) | Partner Site Requests |
| [`filescom-partner-sites`](filescom-partner-sites/SKILL.md) | Partner Sites |
| [`filescom-partners`](filescom-partners/SKILL.md) | Partners |
| [`filescom-permissions`](filescom-permissions/SKILL.md) | Permissions |
| [`filescom-sso-strategies`](filescom-sso-strategies/SKILL.md) | SSO Strategies |
| [`filescom-user-cipher-uses`](filescom-user-cipher-uses/SKILL.md) | User Cipher Uses |
| [`filescom-user-lifecycle-rules`](filescom-user-lifecycle-rules/SKILL.md) | User Lifecycle Rules |
| [`filescom-user-requests`](filescom-user-requests/SKILL.md) | User Requests |
| [`filescom-user-sftp-client-uses`](filescom-user-sftp-client-uses/SKILL.md) | User SFTP Client Uses |
| [`filescom-users`](filescom-users/SKILL.md) | Users |

