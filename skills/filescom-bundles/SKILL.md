---
name: filescom-bundles
description: |
  A Bundle is the API/SDK term for the feature called Share Links in the web interface. Creating a Share Link does not send any email. To deliver the link by email, either call `bundles share` after creation, or pass `--share-after-create=true` when adding a recipient with `bundle-recipients create`.
---

# filescom-bundles

A Bundle is the API/SDK term for the feature called Share Links in the web interface.
The API provides the full set of actions related to Share Links, including sending them via E-Mail.

Please note that we very closely monitor the E-Mailing feature and any abuse will result in disabling of your site.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli bundles list`

List Share Links.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `expires_at`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `expires_at`, `code`, `user_id` or `bypasses_site_expiration_rules`. Valid field combinations are `[ user_id, expires_at ]`. |
| `--filter-gt` | object | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at` and `expires_at`. |
| `--filter-gteq` | object | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at` and `expires_at`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `code`. |
| `--filter-lt` | object | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at` and `expires_at`. |
| `--filter-lteq` | object | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at` and `expires_at`. |
| `--deleted` | bool | If true, only list deleted Share Links. |

### `files-cli bundles find`

Show Share Link.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Bundle ID. **Required.** |
| `--deleted` | bool | If true, show a deleted Share Link. |

### `files-cli bundles create`

Create Share Link.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--paths` | []string | A list of paths to include in this bundle. **Required.** |
| `--password` | string | Password for this bundle. |
| `--bypasses-site-expiration-rules` | bool | If true, this Share Link bypasses site-wide expiration rules. Only site admins may set this. |
| `--form-field-set-id` | int64 | Id of Form Field Set to use with this bundle |
| `--create-snapshot` | bool | If true, create a snapshot of this bundle's contents. |
| `--dont-separate-submissions-by-folder` | bool | Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required. |
| `--expires-at` | datetime | Bundle expiration date/time |
| `--finalize-snapshot` | bool | If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true. |
| `--max-uses` | int64 | Maximum number of times bundle can be accessed |
| `--group-id` | int64 | Owning group ID. If set, members of this group can view, edit, and share this Share Link. |
| `--description` | string | Public description |
| `--note` | string | Bundle internal note |
| `--code` | string | Bundle code.  This code forms the end part of the Public URL. |
| `--path-template` | string | Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data. |
| `--path-template-time-zone` | string | Timezone to use when rendering timestamps in path templates. |
| `--permissions` | enum | Permissions that apply to Folders in this Share Link. One of: `read`, `write`, `read_write`, `full`, `none`, `preview_only`. |
| `--require-registration` | bool | Show a registration page that captures the downloader's name and email address? |
| `--clickwrap-id` | int64 | ID of the clickwrap to use with this bundle. |
| `--inbox-id` | int64 | ID of the associated inbox, if available. |
| `--require-share-recipient` | bool | Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI? |
| `--send-one-time-password-to-recipient-at-registration` | bool | If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set. |
| `--send-email-receipt-to-uploader` | bool | Send delivery receipt to the uploader. Note: For writable share only |
| `--skip-email` | bool | BundleRegistrations can be saved without providing email? |
| `--skip-name` | bool | BundleRegistrations can be saved without providing name? |
| `--skip-company` | bool | BundleRegistrations can be saved without providing company? |
| `--start-access-on-date` | datetime | Date when share will start to be accessible. If `nil` access granted right after create. |
| `--snapshot-id` | int64 | ID of the snapshot containing this bundle's contents. |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |
| `--watermark-attachment-file` | file | Preview watermark image applied to all bundle items. |

### `files-cli bundles share`

Send email(s) with a link to bundle.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Bundle ID. **Required.** |
| `--to` | []string | A list of email addresses to share this bundle with. Required unless `recipients` is used. |
| `--note` | string | Note to include in email. |
| `--recipients` | []object | A list of recipients to share this bundle with. Required unless `to` is used. |

### `files-cli bundles update`

Update Share Link.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Bundle ID. **Required.** |
| `--paths` | []string | A list of paths to include in this bundle. |
| `--password` | string | Password for this bundle. |
| `--bypasses-site-expiration-rules` | bool | If true, this Share Link bypasses site-wide expiration rules. Only site admins may set this. |
| `--form-field-set-id` | int64 | Id of Form Field Set to use with this bundle |
| `--clickwrap-id` | int64 | ID of the clickwrap to use with this bundle. |
| `--code` | string | Bundle code.  This code forms the end part of the Public URL. |
| `--create-snapshot` | bool | If true, create a snapshot of this bundle's contents. |
| `--description` | string | Public description |
| `--dont-separate-submissions-by-folder` | bool | Do not create subfolders for files uploaded to this share. Note: there are subtle security pitfalls with allowing anonymous uploads from multiple users to live in the same folder. We strongly discourage use of this option unless absolutely required. |
| `--expires-at` | datetime | Bundle expiration date/time |
| `--finalize-snapshot` | bool | If true, finalize the snapshot of this bundle's contents. Note that `create_snapshot` must also be true. |
| `--inbox-id` | int64 | ID of the associated inbox, if available. |
| `--max-uses` | int64 | Maximum number of times bundle can be accessed |
| `--group-id` | int64 | Owning group ID. If set, members of this group can view, edit, and share this Share Link. |
| `--note` | string | Bundle internal note |
| `--path-template` | string | Template for creating submission subfolders. Can use the uploader's name, email address, ip, company, `strftime` directives, and any custom form data. |
| `--path-template-time-zone` | string | Timezone to use when rendering timestamps in path templates. |
| `--permissions` | enum | Permissions that apply to Folders in this Share Link. One of: `read`, `write`, `read_write`, `full`, `none`, `preview_only`. |
| `--require-registration` | bool | Show a registration page that captures the downloader's name and email address? |
| `--require-share-recipient` | bool | Only allow access to recipients who have explicitly received the share via an email sent through the Files.com UI? |
| `--send-one-time-password-to-recipient-at-registration` | bool | If true, require_share_recipient bundles will send a one-time password to the recipient when they register. Cannot be enabled if the bundle has a password set. |
| `--send-email-receipt-to-uploader` | bool | Send delivery receipt to the uploader. Note: For writable share only |
| `--skip-company` | bool | BundleRegistrations can be saved without providing company? |
| `--start-access-on-date` | datetime | Date when share will start to be accessible. If `nil` access granted right after create. |
| `--skip-email` | bool | BundleRegistrations can be saved without providing email? |
| `--skip-name` | bool | BundleRegistrations can be saved without providing name? |
| `--workspace-id` | int64 | Workspace ID. `0` means the default workspace. |
| `--user-id` | int64 | The owning user id. Only site admins can set this. |
| `--watermark-attachment-delete` | bool | If true, will delete the file stored in watermark_attachment |
| `--watermark-attachment-file` | file | Preview watermark image applied to all bundle items. |

### `files-cli bundles delete`

Delete Share Link.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Bundle ID. **Required.** |

## Limitations and considerations

**Sending the link is a separate step from creating it.** `bundles create` returns the URL and code but does not email anyone. Two delivery paths exist: `bundles share` sends to one or more addresses; `bundle-recipients create --share-after-create=true` adds a recipient and emails them at creation. Do not delete and recreate the bundle to trigger delivery — neither is necessary, and the per-recipient rate limit below makes recreate-and-resend ineffective anyway.

**Per-recipient invitation emails are rate-limited to once every 24 hours.** A `bundles share` call to the same address within that window will not actually re-send.

**Snapshot Share Links are immutable.** A Share Link created with `--create-snapshot --finalize-snapshot` captures the file set at creation time and cannot be edited afterward. Adding files to the source folder later does not change what a Snapshot link delivers. Live Share Links include source changes; Snapshot Share Links do not.

**Password protection and one-time-password are mutually exclusive.** Only one authentication type is permitted per Share Link. If both `--password` and `--send-one-time-password-to-recipient-at-registration` are set, the create will fail.

**Adding access control after the fact does not revoke active sessions.** Enabling access control on an existing Share Link blocks new visitors but does not remove access from visitors already in a session.

**Site-wide expiration changes are retroactive.** Lowering the site-wide Share Link expiration setting immediately expires any links that have been active longer than the new threshold.

**`--dont-separate-submissions-by-folder` is a security risk.** It allows anonymous uploads from multiple users to live in the same folder, which has subtle security pitfalls. Do not enable without a specific reason.

## Common patterns

### Create a Share Link and email it after creation

1. Create:

       files-cli bundles create --paths=/demo/q3-report.pdf --description="Q3 report" --format=json

   The response includes an `id`. Use it in step 2.

2. Email:

       files-cli bundles share --id=<BUNDLE_ID> --to=bob@example.com --note="Q3 report attached" --format=json

   `bundles share` can deliver to up to 50 addresses in a single call.

### Create a Share Link and email it at the same time (recipient route)

    files-cli bundles create --paths=/demo/q3-report.pdf --description="Q3 report" --format=json
    files-cli bundle-recipients create --bundle-id=<BUNDLE_ID> --recipient=bob@example.com --share-after-create=true --format=json

This adds the recipient row and triggers the invitation email in one call.

### Restrict access to named recipients only

A Share Link with `--require-share-recipient` is only accessible to recipients who received an emailed invitation through Files.com. Create then deliver via `bundles share`:

    files-cli bundles create --paths=/legal/contract.pdf --require-share-recipient --format=json
    files-cli bundles share --id=<BUNDLE_ID> --to=counsel@example.com --to=ceo@example.com --format=json

### Audit who downloaded a Share Link

    files-cli bundle-downloads list --bundle-id=<BUNDLE_ID> --format=json
