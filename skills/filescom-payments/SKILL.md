---
name: filescom-payments
description: |
  An AccountLineItem is a single line item in the accounting ledger for a billing account.
---

# filescom-payments

An AccountLineItem is a single line item in the accounting ledger for a billing account. These include payments and invoices.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli payments list`

List Payments.

No flags beyond the global ones.

### `files-cli payments find`

Show Payment.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Payment ID. **Required.** |

