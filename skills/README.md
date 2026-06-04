# files-cli Skills

Agent-loadable skill packages for the [`files-cli`](https://github.com/Files-com/files-cli) CLI.

## Structure

- **Domain skills** (`filescom-<domain>/`) — one per top-level CLI command (`bundles`, `users`, `folders`, …). Each skill body lists the subcommands and flags for that domain.

The full index is in `INDEX.md`.

## Consuming

These skills follow the Anthropic SKILL.md frontmatter contract: each `SKILL.md` declares a `name` and a `description`. An agent loads the relevant skill based on the `description` when its current task matches.

For Claude Code, Codex, or any agent that supports filesystem-based skills, point the skills directory at this folder.

For agents that don't support skill loading directly, the same content is readable as plain markdown — load the relevant `SKILL.md` into context manually.

## Companion files

- `../CONTEXT.md` — CLI-wide invocation contract, authentication, global flags. Every skill assumes you've read this.
- `../agents/tool-catalog.json` — machine-readable catalog of every command and parameter. Use for programmatic surface discovery.
- `../agents/error-catalog.json` — machine-readable catalog of every known error type with HTTP code.
