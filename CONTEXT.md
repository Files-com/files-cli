# files-cli — Agent Context

This file documents the agent-relevant invocation contract for `files-cli`. For deeper coverage of configuration, file operations, sort and filter, and pagination, see the README and <https://www.files.com/docs/client-apps/command-line-interface-cli-app>.

## Agent invocation

```bash
files-cli <domain> <subcommand> --format json --non-interactive [flags...]
```

- Pass `--format json` so output is structured. The default output format is a human-readable table.
- Pass `--non-interactive` so the CLI never blocks on a prompt.

## Authentication

Two supported approaches: API key (recommended for agents and automation) or session.

### API key

Pass `--api-key` per command:

```bash
files-cli --api-key=YOUR_API_KEY <domain> <subcommand> ...
```

API key authentication does not trigger a Two-Factor Authentication challenge, even on accounts that require 2FA for web login.

The CLI also reads the `FILES_API_KEY` environment variable. After an API key is provided once, it is written to the `files-cli` configuration file and reused on subsequent commands.

### Session

Configure once, then log in:

```bash
files-cli config set --subdomain=MYSITENAME --username=MYUSERNAME
files-cli login
```

Login sessions expire automatically after 6 hours, or sooner if the site's authentication settings dictate a shorter timeout. For custom domains, pass `--endpoint=fully.qualified.host` instead of `--subdomain`.

## Global flags (persistent on every subcommand)

Sourced from the CLI's root-command flag definitions.

| Flag | Purpose |
| --- | --- |
| `--api-key=KEY` | Set API key for single use. |
| `--session-id=ID` | Set session ID for single use. |
| `--profile=NAME` | Use a named connection profile. |
| `--workspace-id=ID` | Scope this command to a specific workspace. |
| `--endpoint=HOST` | Override the API endpoint (custom domains). |
| `--format=FORMAT` | Output format (`json` for agents). |
| `--fields=LIST` | Comma-separated field names to include. |
| `--output-path=PATH` | Write output to a file instead of stdout. |
| `--debug[=PATH]` | Enable verbose logging. `--debug=STDOUT` prints to the screen; `--debug=<file>` writes to a log file. |
| `--non-interactive` | Do not prompt for user input. |
| `--reauthentication` | Re-supply session user's password for security-sensitive operations. |
| `--ignore-version-check` | Skip the CLI version check on startup. |

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

When a command fails, the response includes a `type` field — a stable, hierarchical string like `bad-request/missing-field` or `not-authorized/reauthentication-needed-action`. Route on `type`.

The full machine-readable catalog of every known error type, with HTTP codes, is at `agents/error-catalog.json`.

## Tool catalog

The full machine-readable catalog of every command and parameter is at `agents/tool-catalog.json`. The per-domain skills under `skills/` cover the same surface in a per-command narrative form.
