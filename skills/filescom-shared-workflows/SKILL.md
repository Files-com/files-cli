---
name: filescom-shared-workflows
description: |
  Files.com Shared Workflows via files-cli.
---

# filescom-shared-workflows

Files.com Shared Workflows management via files-cli.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli shared-workflows list`

List Shared Workflows.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |

### `files-cli shared-workflows find`

Show Shared Workflow.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Shared Workflow ID. **Required.** |

### `files-cli shared-workflows find-run`

Show a Shared Workflow run started by the recipient.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Shared Workflow ID. **Required.** |
| `--run-id` | int64 | ID of the run returned by the Shared Workflow submission. **Required.** |

### `files-cli shared-workflows run-workflow`

Run a Shared Workflow.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Shared Workflow ID. **Required.** |
| `--submission-revision` | string | Revision returned when opening the Shared Workflow. **Required.** |
| `--submission-id` | string | Client-generated UUID used to retry the same submission without starting another run. **Required.** |
| `--file-path` | string | Path of a file the recipient can read, including a file uploaded to a Snapshot created for this workflow. **Required.** |

