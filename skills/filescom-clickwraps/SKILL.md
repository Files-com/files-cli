---
name: filescom-clickwraps
description: |
  A Clickwrap is a legal agreement (such as an NDA or Terms of Use) that your Users and/or Bundle/Inbox participants will need to agree to via a "Clickwrap" UI before accessing the site, bundle, or inbox.
---

# filescom-clickwraps

A Clickwrap is a legal agreement (such as an NDA or Terms of Use) that your Users and/or Bundle/Inbox participants will need to agree to via a "Clickwrap" UI before accessing the site, bundle, or inbox.

The values for `use_with_users`, `use_with_bundles`, `use_with_inboxes` are explained as follows:

* `none` - This Clickwrap may not be used in this context.
* `available_to_all_users` - This Clickwrap may be assigned in this context by any user.
* `available` - This Clickwrap may be assigned in this context, but only by Site Admins.  We recognize that the name of this setting is somewhat ambiguous, but we maintain it for legacy reasons.
* `required` - This Clickwrap will always be used in this context, and may not be overridden.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli clickwraps list`

List Clickwraps.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |

### `files-cli clickwraps find`

Show Clickwrap.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Clickwrap ID. **Required.** |

### `files-cli clickwraps create`

Create Clickwrap.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.) |
| `--body` | string | Body text of Clickwrap (supports Markdown formatting). |
| `--use-with-bundles` | enum | Use this Clickwrap for Bundles?. One of: `none`, `available`, `require`, `available_to_all_users`. |
| `--use-with-inboxes` | enum | Use this Clickwrap for Inboxes?. One of: `none`, `available`, `require`, `available_to_all_users`. |
| `--use-with-users` | enum | Use this Clickwrap for Users?  Values: `none`, `require` (new user signup via email invitation only), `require_all_users_once` (show to all users at their next web login; once accepted, not shown again), `require_all_users_always` (show to all users on every web login). One of: `none`, `require`, `require_all_users_once`, `require_all_users_always`. |

### `files-cli clickwraps update`

Update Clickwrap.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Clickwrap ID. **Required.** |
| `--name` | string | Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.) |
| `--body` | string | Body text of Clickwrap (supports Markdown formatting). |
| `--use-with-bundles` | enum | Use this Clickwrap for Bundles?. One of: `none`, `available`, `require`, `available_to_all_users`. |
| `--use-with-inboxes` | enum | Use this Clickwrap for Inboxes?. One of: `none`, `available`, `require`, `available_to_all_users`. |
| `--use-with-users` | enum | Use this Clickwrap for Users?  Values: `none`, `require` (new user signup via email invitation only), `require_all_users_once` (show to all users at their next web login; once accepted, not shown again), `require_all_users_always` (show to all users on every web login). One of: `none`, `require`, `require_all_users_once`, `require_all_users_always`. |

### `files-cli clickwraps delete`

Delete Clickwrap.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Clickwrap ID. **Required.** |

