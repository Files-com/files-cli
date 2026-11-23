---
name: filescom-certificates
description: |
  A Certificate is a way to store your custom SSL certificate for your Files.com site.
---

# filescom-certificates

A Certificate is a way to store your custom SSL certificate for your Files.com site.
Files.com is also capable of registering SSL certificates itself, so this is only required if you have specific security requirements to use your own certificate.

All subcommands also accept the flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli certificates list`

List Certificates.

No flags beyond the global ones.

### `files-cli certificates find`

Show Certificate.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Certificate ID. **Required.** |

### `files-cli certificates create`

Create CSR or import existing SSL Certificate.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Internal name of the SSL certificate. **Required.** |
| `--certificate-domain` | string | Domain for certificate. |
| `--certificate-extra-domains` | []string | Additional domains for the certificate CSR. |
| `--certificate-country` | string | Country. |
| `--certificate-state-or-province` | string | State or province. |
| `--certificate-city-or-locale` | string | City or locale. |
| `--certificate-company-name` | string | Company name. |
| `--csr-ou1` | string | Department name or organization unit 1 |
| `--csr-ou2` | string | Department name or organization unit 2 |
| `--csr-ou3` | string | Department name or organization unit 3 |
| `--certificate-email-address` | string | Email address for the certificate owner. |
| `--key-type` | string | Any supported key type.  Defaults to `RSA-4096`. |
| `--certificate` | string | The certificate.  PEM or PKCS7 formats accepted. |
| `--private-key` | string | Certificate private key. |
| `--password` | string | Certificate password. |
| `--intermediates` | string | Intermediate certificates.  PEM or PKCS7 formats accepted. |

### `files-cli certificates activate`

Activate SSL Certificate.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Certificate ID. **Required.** |
| `--replace-cert` | string | Leave blank or set to `any` to replace any currently active certificate, if applicable.  Set to `new_ips` to activate this certificate as an additional certificate on your Site by allocating a new set of Dedicated IPs (may require a Plan upgrade).  Set to the ID of a currently active certificate to replace that certificate on its set of dedicated IPs.  A certificate belonging to a Custom Domain activates only for that domain: leave blank, use `any`, or specify that domain's current certificate ID. |

### `files-cli certificates deactivate`

Deactivate SSL Certificate.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Certificate ID. **Required.** |

### `files-cli certificates replace-with-files`

Replace this active custom SSL Certificate for the Primary Custom Domain with a Files.com SSL Certificate.  The Primary Custom Domain's CNAME record must point at the Site's subdomain.  The Files.com certificate is issued and activated in the background.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Certificate ID. **Required.** |

### `files-cli certificates update`

Update Certificate.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Certificate ID. **Required.** |
| `--name` | string | Internal certificate name. |
| `--intermediates` | string | Any intermediate certificates.  PEM or PKCS7 formats accepted. |
| `--certificate` | string | The certificate.  PEM or PKCS7 formats accepted. |

### `files-cli certificates delete`

Delete Certificate.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Certificate ID. **Required.** |

