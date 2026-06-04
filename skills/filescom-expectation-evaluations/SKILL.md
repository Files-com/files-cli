---
name: filescom-expectation-evaluations
description: |
  An ExpectationEvaluation records one open or closed window for an Expectation.
---

# filescom-expectation-evaluations

An ExpectationEvaluation records one open or closed window for an Expectation.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli expectation-evaluations list`

List Expectation Evaluations.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `expectation_id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `expectation_id` and `workspace_id`. Valid field combinations are `[ workspace_id, expectation_id ]`. |

### `files-cli expectation-evaluations find`

Show Expectation Evaluation.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Expectation Evaluation ID. **Required.** |

