---
name: filescom-webhook-tests
description: |
  A WebhookTest is an object that can be sent to your Webhook endpoint for connection and response tests.
---

# filescom-webhook-tests

A WebhookTest is an object that can be sent to your Webhook endpoint for connection and response tests.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli webhook-tests create`

Create Webhook Test.

| Flag | Type | Description |
| --- | --- | --- |
| `--url` | string | URL for testing the webhook. **Required.** |
| `--method` | string | HTTP method(GET or POST). |
| `--encoding` | string | HTTP encoding method.  Can be JSON, XML, or RAW (form data). |
| `--headers` | object | Additional request headers. |
| `--body` | object | Additional body parameters. |
| `--raw-body` | string | raw body text |
| `--file-as-body` | bool | Send the file data as the request body? |
| `--file-form-field` | string | Send the file data as a named parameter in the request POST body |
| `--use-dedicated-ips` | bool | Use dedicated IPs for sending the webhook? |

