# files-cli Skills

Agent-loadable skill packages for the [`files-cli`](https://github.com/Files-com/files-cli) CLI.

## Structure

- **Domain skills** (`filescom-<domain>/`) — one per top-level CLI command (`bundles`, `users`, `folders`, …). Each skill body lists the subcommands and flags for that domain.

The full index is in `INDEX.md`.

## Consuming

These skills follow the Anthropic SKILL.md frontmatter contract: each `SKILL.md` declares a `name` and a `description`. An agent loads the relevant skill based on the `description` when its current task matches.

For Claude Code, Codex, or any agent that supports filesystem-based skills, point the skills directory at this folder.

For agents that don't support skill loading directly, the same content is readable as plain markdown — load the relevant `SKILL.md` into context manually.

## From the CLI

The installed binary serves the same recipes, matched to its version, without credentials or network access: `files-cli workflows` lists them and `files-cli workflows show <name>` prints one. `files-cli commands` describes the binary's exact commands and flags, including commands such as `upload` that the domain skills do not cover.

## Companion files

- `../CONTEXT.md` — CLI-wide invocation contract, authentication, global flags. Every skill assumes you've read this.
- `../agents/tool-catalog.json` — machine-readable catalog of the API resource commands and their parameters, generated from the API schema. `files-cli commands describe` covers every command of the installed binary.
- `../agents/error-catalog.json` — machine-readable catalog of every known error type with HTTP code.
