---
name: recipe-generating-reports
description: |
  Choose the right built-in Files.com report or export instead of assembling one by hand; use this when the user asks for a report, audit, export, or summary — storage/usage, bandwidth/transfer, activity/audit logs, webhook delivery, permissions, group membership, or share-link audits. Files.com already produces these as first-class resources (`usage-snapshots`, `bandwidth-snapshots`, `history-exports`, `action-notification-exports`, and `scheduled-exports`). Prefer them; only home-roll a report by listing and aggregating records when no built-in report covers the need.
---

# recipe-generating-reports

Files.com ships built-in reports and exports for the questions people ask most often. **Reach for an existing report before building one yourself.** Home-rolling a report — listing records and aggregating them in code — is slower, burns API rate limit, and usually reproduces data the platform already computes (often more accurately, since the built-in figures account for deleted-but-retained files, billing high-water marks, and the like).

## When to use this skill

- Any request phrased as a "report", "audit", "export", "summary", or "breakdown".
- Before writing a loop that lists a resource just to count, sum, or tabulate it — check whether a built-in report already answers the question.

## The built-in reports

| Question | Use | Notes |
| --- | --- | --- |
| Storage usage, file counts, user counts, API-call usage | `usage-snapshots`, `usage-daily-snapshots` | Read immediately. Per-top-level-folder size and count via `usage_by_top_level_dir`. See `recipe-folder-size-and-counts`. |
| Bandwidth / transfer over time | `bandwidth-snapshots` | Supports `--filter`, `--filter-gt/-gteq/-lt/-lteq`, and `--sort-by` on the time range. |
| Activity / audit log — who did what, when | `history-exports` → `history-export-results` | **Readable.** Asynchronous export of the action history, richly filterable; results download via `results_url`. See the flow below. |
| Webhook / Action Notification delivery log | `action-notification-exports` → `action-notification-export-results` | **Readable.** Asynchronous, same create-then-read pattern as history exports. |
| Folder size across every folder | `scheduled-exports` with `--export-type=folder_size_audit` | **Emailed, not readable.** Recurring, sent to a Site Admin. For readable folder sizes use `usage-snapshots`; see `recipe-folder-size-and-counts`. |
| Permissions audit — who can access what | `scheduled-exports` with `--export-type=permission_audit` | **Emailed, not readable.** Recurring. `--export-options` supports `group_by` of `user` or `path`. See `filescom-permissions`. |
| Group membership audit | `scheduled-exports` with `--export-type=group_membership_audit` | **Emailed, not readable.** Recurring. |
| Share Link (bundle) audit | `scheduled-exports` with `--export-type=share_link_audit` | **Emailed, not readable.** Recurring. |

**Readable vs. emailed is the key distinction.** Snapshots and the `history` / `action-notification` exports return their data to the API, so the agent can read and present it. The `scheduled-exports` audit reports are delivered only by email to a Site Admin — pick one of those only when the user wants a recurring report sent to a person, never as a way for the agent to obtain the data.

## Activity / audit reports — the history-export flow

Exporting the action history is a two-step asynchronous operation. Create the export with whatever filters you need, poll it until it is `ready`, then read the results.

1. Create, narrowing with `query_*` filters and a date range:

       files-cli history-exports create --start-at=2026-01-01T00:00:00Z --end-at=2026-01-31T23:59:59Z --query-action=login --format=json

   Filters include `--query-action`, `--query-interface`, `--query-user-id`, `--query-path`, `--query-folder`, `--query-ip`, and more — see the `filescom-history-exports` skill.

2. Poll until `status` is `ready` (it starts `building`; it can also become `failed`):

       files-cli history-exports find --id=EXPORT_ID --format=json

3. Read the rows. When ready, the export carries a `results_url` (a single CSV download of all rows), or page through the rows directly:

       files-cli history-export-results list --history-export-id=EXPORT_ID --format=json

`action-notification-exports` follows the same pattern, paired with `action-notification-export-results list --action-notification-export-id=EXPORT_ID`.

## Scheduled Exports are recurring and emailed

The `scheduled-exports` audit reports (`folder_size_audit`, `permission_audit`, `group_membership_audit`, `share_link_audit`) run on a schedule and are **emailed** to the Site Administrator named in `--user-id`. Unlike `history-exports` and `action-notification-exports`, a `ScheduledExport` has no `status`, no `results_url`, and no results endpoint — there is nothing for the agent to read back, no on-demand "run now" trigger, and no API way to download the generated file. Create one only to set up recurring delivery to a person; if the user needs an immediate, downloadable copy, that is done from the relevant tab in the Files.com web interface. See the `filescom-scheduled-exports` skill for the scheduling flags.

## When home-rolling is justified

Only assemble a report yourself when no built-in report covers the question. If you must:

- List the narrowest resource that holds the data and aggregate from there.
- Remember that list endpoints are paginated (`--cursor`, `--per-page`, `--max-pages`) and rate-limited — do not walk large trees.
- Note that usage figures and search results are best-effort and not real time (see the caveats in the resource skills).

## Related skills

- `filescom-usage-snapshots`, `filescom-usage-daily-snapshots` — storage, transfer, user, and API-call usage.
- `filescom-bandwidth-snapshots` — bandwidth/transfer series.
- `filescom-history-exports`, `filescom-history-export-results` — action/audit history exports.
- `filescom-action-notification-exports`, `filescom-action-notification-export-results` — webhook delivery exports.
- `filescom-scheduled-exports` — recurring audit reports (folder size, permissions, group membership, share links).
- `recipe-folder-size-and-counts` — storage size and file counts specifically.
