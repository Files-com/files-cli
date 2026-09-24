# files-cli — Agent Context

This file documents the agent-relevant invocation contract for `files-cli`. For deeper coverage of configuration, file operations, sort and filter, and pagination, see the README and <https://www.files.com/docs/client-apps/command-line-interface-cli-app>.

## Agent invocation

```bash
files-cli <domain> <subcommand> --format json --non-interactive [flags...]
```

- Pass `--format json` so output is structured. The default output format is a human-readable table. `--format` is defined on each resource command rather than globally.
- Pass `--non-interactive` so the CLI never blocks on a prompt.

## Discover commands offline

The installed binary describes its own commands and ships these guides. Discovery needs no credentials or network access and never reads or writes the config file.

```bash
files-cli commands                                      # top-level commands and groups
files-cli commands list folders                         # the commands in one group
files-cli commands search share link                    # keyword search, 20 results by default (--limit)
files-cli commands describe folders list-for --format json
files-cli commands describe users list --flag cursor    # one flag, full description
files-cli workflows                                     # task guides
files-cli workflows show recipe-searching-for-files     # one guide as Markdown
```

`commands describe` reports usage, positional arguments, and local and inherited flags with their type, static default, and enum values. `required` marks flags the CLI rejects a command without; `api_required` marks parameters the Files.com API requires, which the CLI sends without checking. Flag descriptions longer than 200 characters are cut and marked `truncated`; add `--full` or `--flag NAME` for the complete text.

## Bounded listing with continuation

List commands keep their existing output: `--format json` prints one JSON array and, unless `--max-pages` is set, fetches every page. To read a list incrementally, add `--json-envelope`:

```bash
files-cli users list --json-envelope --per-page=100 --format json,raw --non-interactive
```

```json
{"has_more":true,"next_cursor":"CURSOR","data":[{"id":1,"username":"ann"}]}
```

Pass `next_cursor` back with `--cursor` until `has_more` is `false`, when `next_cursor` is `null`:

```bash
files-cli users list --json-envelope --per-page=100 --cursor=CURSOR --format json,raw --non-interactive
```

- `--json-envelope` fetches one page unless `--max-pages` is given; `--max-pages=0` fetches every page.
- `--fields`, `--filter-by`, and the API `--filter*` and `--sort-by` flags work as usual. Client-side filters such as `--filter-by` apply after a page is fetched, so `data` can be empty while `has_more` is `true`.
- The envelope is always JSON: `--format json,raw` makes it compact, and a non-JSON `--format` is rejected before the list is requested.
- The envelope is written once every selected page has been fetched, so those records are held in memory until then. Keep `--max-pages` bounded for large lists; `--max-pages=0` holds the entire list.
- If a page fails, nothing is written to stdout and the command exits non-zero; retry with the same `--cursor`.
- The flag is available on every cursor-based list command, including `folders list-for --recursive`.

## Authentication

Two supported approaches: API key (recommended for agents and automation) or session.

### API key

Pass `--api-key` per command:

```bash
files-cli --api-key=YOUR_API_KEY <domain> <subcommand> ...
```

API key authentication does not trigger a Two-Factor Authentication challenge, even on accounts that require 2FA for web login.

The CLI also reads the `FILES_API_KEY` environment variable. Neither `--api-key` nor `FILES_API_KEY` is saved; to store a key in the `files-cli` configuration file for later commands, run `files-cli config set --api-key=YOUR_API_KEY`.

### Session

Configure once, then log in:

```bash
files-cli config set --subdomain=MYSITENAME --username=MYUSERNAME
files-cli login
```

Login sessions expire automatically after 6 hours, or sooner if the site's authentication settings dictate a shorter timeout. For custom domains, pass `--endpoint=https://files.example.com` to `config set` instead of `--subdomain`.

## Global flags (persistent on every subcommand)

Sourced from the CLI's root-command flag definitions.

| Flag | Purpose |
| --- | --- |
| `--api-key=KEY` | Set API key for single use. |
| `--session-id=ID` | Set session ID for single use. |
| `--profile=NAME` | Use a named connection profile. |
| `--workspace-id=ID` | Scope this command to a specific workspace. |
| `--output=PATH`, `-o PATH` | Write output to a file instead of stdout. |
| `--debug[=PATH]` | Enable verbose logging. `--debug=STDOUT` prints to the screen; `--debug=<file>` writes to a log file. |
| `--non-interactive` | Do not prompt for user input. |
| `--reauthentication` | Re-supply session user's password for security-sensitive operations. |
| `--ignore-version-check` | Skip the CLI version check on startup. |

Resource commands also define `--format` (`json` for agents), `--fields` (comma-separated field names to include), and `--use-pager`; list commands add `--cursor`, `--per-page`, `--max-pages`, and `--json-envelope`. `files-cli commands describe <command>` shows the exact flags of any command.

## Workspaces

A Workspace is a way to organize related resources inside a single Files.com site. Every site has an implicit Default workspace with ID `0`; resources not explicitly assigned to a named workspace belong to the Default workspace.

Scope a single command with `--workspace-id`:

```bash
files-cli --workspace-id=YOUR_WORKSPACE_ID folders list-for ''
```

Or store the workspace ID so every subsequent command is scoped:

```bash
files-cli config set --workspace-id=YOUR_WORKSPACE_ID
```

To clear the stored workspace ID:

```bash
files-cli config reset --workspace-id
```

### Workspace constraints

Workspaces isolate resources. A user assigned to a named Workspace (a non-zero `workspace_id`) can only see and act on resources within that Workspace; a Site Administrator is not confined to any Workspace.

**Resources cannot be moved between Workspaces.** No API or CLI operation reassigns a folder, automation, remote server, or any other resource from one Workspace to another, so do not attempt it. The only supported Workspace reassignment is a Site Administrator returning a *user* to the Default workspace by changing that user's `workspace_id` from a non-zero value to `0`:

```bash
files-cli users update --id=USER_ID --workspace-id=0
```

**Give a Default-workspace user access to other Workspaces with Permissions, not by reassigning a Workspace.** A user on `workspace_id=0` is not confined to the Default Workspace. Grant them access to resources that belong to other Workspaces by adding a Permission record for that user on the relevant path, where the path is prefixed like `_/Workspaces/$WORKSPACE_ID/$FOLDER_PATH`:

```bash
files-cli permissions create --path=_/Workspaces/WORKSPACE_ID/FOLDER_PATH --user-id=USER_ID --permission=LEVEL
```

If `$FOLDER_PATH` is empty, the permission applies to the Workspace's root folder; granting `admin` on a Workspace's root folder grants Workspace Admin access to the entire Workspace. See the `filescom-permissions` skill for permission levels and the `filescom-workspaces` skill for managing Workspaces.

## Errors

Check the process exit code: `0` indicates success, and a non-zero code indicates failure. Diagnostics are written to stderr. `--format json` selects result formatting; it does not provide a uniform JSON error envelope on stdout.

`agents/error-catalog.json` describes Files.com API error types and HTTP codes. CLI usage errors, local filesystem errors, and other failures do not necessarily carry an API error type.

## Tool catalog

A machine-readable catalog of the API resource commands and their API parameters is at `agents/tool-catalog.json`. The per-domain skills under `skills/` cover the same surface in a per-command narrative form. Both are generated from the API schema, so they omit commands such as `upload`, `download`, and `sync`; `files-cli commands` describes the exact commands and flags of the installed binary.
