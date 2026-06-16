---
name: recipe-folder-size-and-counts
description: |
  Determine storage usage and file/folder counts on Files.com from the data the agent can actually read back, without recursively walking the tree; use this when the user asks "how big is this folder", "how much storage are we using", or "how many files do we have". The agent-readable source is `usage-snapshots` / `usage-daily-snapshots` (size and file count per top-level folder) and `site/usage` (site totals). Note the important limits: usage data only breaks down by top-level folder, and the `folder_size_audit` scheduled export is emailed to a Site Admin rather than returned to the API.
---

# recipe-folder-size-and-counts

Files.com computes storage usage and file counts on a recurring schedule. **Read that precomputed data — do not recursively walk a folder tree to add up sizes or count files.** A manual walk is slow, burns API rate limit, and still will not match the platform's own accounting. The catch is that not every size figure is readable back through the API; this recipe is about getting an answer the agent can actually *see*.

## When to use this skill

- "How much storage is this folder / our site using?"
- "How many files do we have?" / "How many files are under this top-level folder?"
- A previous attempt tried to `list` folders recursively to sum sizes — stop and use the data below.

## What the agent can read directly

`usage-snapshots` and `usage-daily-snapshots` are plain `list` reads — the results come straight back, no email, no polling.

    files-cli usage-snapshots list --format=json

Each snapshot includes:

- `current_storage`, `high_water_storage`, `root_storage` — site-wide storage in GB (`high_water_storage` is what billing uses).
- `usage_by_top_level_dir` — an array of `{ "dir": <name>, "size": <usage>, "count": <file count> }`, one entry per **top-level** folder. This is your readable size-and-count breakdown.

For a day-by-day series with the same per-top-level-folder breakdown:

    files-cli usage-daily-snapshots list --format=json

Site-wide totals on their own are also available from `site/usage`.

## Limitation: only top-level folders are broken out

`usage_by_top_level_dir` covers the **root-level** folders only. There is no API/CLI read that returns a recursive total size for an arbitrary *nested* folder. Folder listings (`folders list-for`) return each entry's own `size`, not a subtree total — so a precise size for a deep folder would require walking its subfolders, which is exactly what to avoid. If the user needs sizes below the top level, say so plainly and offer the top-level breakdown, the single-folder option below, or the emailed audit (next section) for a human.

## The `folder_size_audit` scheduled export is email-only — you cannot read it back

Files.com can produce a size report across **every** folder via a Scheduled Export of type `folder_size_audit`. **But the agent cannot retrieve the result.** A `ScheduledExport` has no `status`, no `results_url`, and no results endpoint — the finished report is **emailed** to the Site Administrator named in `--user-id`, and runs only on a schedule (no on-demand "run now"). Unlike `history-exports` / `action-notification-exports`, which return a downloadable `results_url`, there is nothing here for an agent to read.

So only create one when the user explicitly wants a recurring report **delivered by email to a person**:

    files-cli scheduled-exports create --name="Folder size audit" --export-type=folder_size_audit --user-id=SITE_ADMIN_USER_ID --trigger=daily --format=json

Do not create a `folder_size_audit` export expecting to read the numbers from it — use the readable `usage-snapshots` data above. For an immediate, downloadable per-folder export, that is done from the **Usage** tab in the Files.com web interface.

## Caveats on the readable usage data

- **Not real time.** Usage is recalculated several times a day. After deletions, the updated figures and folder breakdown can take up to 24 hours to catch up.
- **Storage figures are the platform's accounting** and may include deleted-but-retained backup files — read the field descriptions in the `filescom-usage-snapshots` skill to pick the right number for the question.

## When a manual count is genuinely unavoidable

To confirm the live contents of a **single** folder (e.g. whether a small folder is empty), list just that one folder and filter by type — do not recurse:

    files-cli folders list-for "/path/to/folder" --type=file --format=json
    files-cli folders list-for "/path/to/folder" --type=folder --format=json

## Related skills

- `filescom-usage-snapshots` — readable site and top-level-folder storage usage and file counts.
- `filescom-usage-daily-snapshots` — daily readable usage series with per-top-level-folder breakdown.
- `filescom-scheduled-exports` — `folder_size_audit` and other emailed (non-readable) audit reports.
- `filescom-folders` — folder listing (`list-for`), for the single-folder case only.
- `recipe-generating-reports` — which reports are readable by the agent vs. emailed to a person.
