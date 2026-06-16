---
name: recipe-searching-for-files
description: |
  Find files and folders on Files.com effectively, and pick the right tool for the job; use this when the user wants to locate files by name, path, modified time, or custom metadata. The key distinction: the `--search` / `--search-all` flags on `folders list-for` are the web search bar — best-effort, not real time, and meant only for ad-hoc human lookups — while automated or precise lookups should list a folder directly and filter deterministically. Covers both, and when to use each.
---

# recipe-searching-for-files

There are two ways to find files on Files.com, and choosing the wrong one is the usual mistake. Name search (`--search` / `--search-all`) is the web UI's search bar exposed through the API: convenient for a person, but a best-effort index that lags reality. Direct folder listing (`folders list-for`) is exact and current. Match the tool to whether a human or an automated process is consuming the result.

## When to use this skill

- "Find the file called …", "where is …", "search for …".
- Locating files by name, modification time, or a custom metadata value.
- Before writing an automated job that relies on search results — read the warning below first.

## Ad-hoc human search — `--search` and `--search-all`

These are the same engine as the web UI's "Search This Folder" and "Search All Files" boxes.

Search within a folder (and its subfolders — name search is recursive):

    files-cli folders list-for "/projects" --search="invoice" --format=json

Search the entire site (do **not** pass a path; use an empty path with `--search-all`):

    files-cli folders list-for "" --search="invoice" --search-all=true --format=json

What to know about this search:

- **Not real time, and best effort.** The index reflects uploads, moves, and deletes only after a delay, and is not guaranteed to match the latest folder listing. Results may be **truncated past ~1,000 matches**.
- **For ad-hoc human use only.** The API explicitly states this search "should only be used for ad-hoc (human) searching, and not as part of an automated process." Do not build automation on it.
- **Names, not contents.** It matches file and folder *names*. In-file/content search is not available.
- **`--search` ignores text before the last `/`.** Pass just the name fragment, not a full path, as the query.
- **Remote-mounted folders are not indexed** — for those, search only checks the current folder.

### Search by custom metadata

Pair `--search-custom-metadata-key` with `--search` to match a metadata value (this is non-recursive — current folder only). Set the key to `*` to match the value against any metadata key:

    files-cli folders list-for "/projects" --search="active" --search-custom-metadata-key="status" --format=json

## Precise or automated lookup — `folders list-for`

For anything an automated process depends on, or when you need an exact, current answer, list the folder directly and filter. This reflects the true current state of the folder.

Return only files, or only folders, with `--type`:

    files-cli folders list-for "/projects" --type=file --format=json

Sort by name, size, or modified time (`--sort-by` is a JSON object; valid fields are `path`, `size`, `modified_at_datetime`, `provided_modified_at`):

    files-cli folders list-for "/projects" --type=file --sort-by='{"modified_at_datetime":"desc"}' --format=json

Return only items modified after a given time (`--modified-at-datetime` must be combined with `--type`):

    files-cli folders list-for "/projects" --type=file --modified-at-datetime=2026-01-01T00:00:00Z --format=json

Listings are paginated — page through with the global `--cursor` / `--per-page` / `--max-pages` flags. To search a subtree deterministically, list each folder you care about; descend only into the folders you actually need rather than walking the entire tree.

## Anti-patterns

- **Don't use `--search` / `--search-all` in automation.** It is non-real-time and capped; an automated job will silently miss recent or beyond-1,000 results. List the folder and filter instead.
- **Don't expect content matches.** Search is name-only.
- **Don't pass a `--path` together with `--search-all`.** Site-wide search takes an empty path.
- **Don't recursively walk a huge tree just to find files.** Narrow with `--type`, `--modified-at-datetime`, and targeted folder paths; descend only where needed.

## Related skills

- `filescom-folders` — full reference for `folders list-for` and its flags.
- `filescom-files` — operating on a file once you have its path (download, move, metadata).
- `recipe-folder-size-and-counts` — sizes and counts without walking the tree.
