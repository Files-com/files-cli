---
name: filescom-ip-addresses
description: |
  An IPAddress is a record of IP addresses that you can use to automate keeping your firewall's configuration up to date.
---

# filescom-ip-addresses

An IPAddress is a record of IP addresses that you can use to automate keeping your firewall's configuration up to date.

Customers who maintain custom firewall configurations may require knowing the public IP addresses of Files.com's edge servers.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli ip-addresses list`

List IP Addresses associated with the current site.

No flags beyond the global ones.

### `files-cli ip-addresses get-smartfile-reserved`

List all possible public SmartFile IP addresses.

No flags beyond the global ones.

### `files-cli ip-addresses get-exavault-reserved`

List all possible public ExaVault IP addresses.

No flags beyond the global ones.

### `files-cli ip-addresses get-reserved`

List all possible public IP addresses.

No flags beyond the global ones.

