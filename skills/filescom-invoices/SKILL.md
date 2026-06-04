---
name: filescom-invoices
description: |
  An AccountLineItem is a single line item in the accounting ledger for a billing account.
---

# filescom-invoices

An AccountLineItem is a single line item in the accounting ledger for a billing account. These include payments and invoices.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli invoices list`

List Invoices.

No flags beyond the global ones.

### `files-cli invoices find`

Show Invoice.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Invoice ID. **Required.** |

