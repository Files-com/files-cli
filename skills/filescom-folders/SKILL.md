---
name: filescom-folders
description: |
  Files.com Folders via files-cli.
---

# filescom-folders

Files.com Folders management via files-cli.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli folders list-for [path]`

List Folders by Path.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--preview-size` | string | Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`. |
| `--sort-by` | object | Search by field and direction. Valid fields are `path`, `size`, `modified_at_datetime`, `provided_modified_at`.  Valid directions are `asc` and `desc`.  Defaults to `{"path":"asc"}`. |
| `--search` | string | If specified, will search the folders/files list by name. Ignores text before last `/`. This is the same API used by the search bar in the web UI when running 'Search This Folder'.  Search results are a best effort, not real time, and not guaranteed to perfectly match the latest folder listing.  Results may be truncated if more than 1,000 possible matches exist.  This field should only be used for ad-hoc (human) searching, and not as part of an automated process. |
| `--search-custom-metadata-key` | string | If provided, the search string in `search` will search for files where this custom metadata key matches the value sent in `search`.  Set this to `*` to allow any metadata key to match the value sent in `search`. |
| `--search-all` | bool | Search entire site?  If true, we will search the entire site.  Do not provide a path when using this parameter.  This is the same API used by the search bar in the web UI when running 'Search All Files'.  Search results are a best effort, not real time, and not guaranteed to match every file.  This field should only be used for ad-hoc (human) searching, and not as part of an automated process. |
| `--with-previews` | bool | Include file previews? |
| `--with-priority-color` | bool | Include file priority color information? |
| `--type` | string | Type of objects to return.  Can be `folder` or `file`. |
| `--modified-at-datetime` | datetime | If provided, will only return files/folders modified after this time. Can be used only in combination with `type` filter. |

### `files-cli folders create [path]`

Create Folder.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--mkdir-parents` | bool | Create parent directories if they do not exist? |
| `--provided-mtime` | datetime | User provided modification time. |

