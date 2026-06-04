---
name: filescom-form-field-sets
description: |
  A Form Field Set is a custom form to be used for bundle and inbox registrations.
---

# filescom-form-field-sets

A Form Field Set is a custom form to be used for bundle and inbox registrations.

Each Form Field Set contains one or more Form Fields. A form and all of its form fields are submitted in a single create request. The order of form fields in the array is the order they will be displayed.

Once created, a form field set can then be associated with one or more bundle(s) and/or inbox(s). Once associated, you will be required to submit well-formatted form-data when creating a bundle-registration or inbox registration.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli form-field-sets list`

List Form Field Sets.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |

### `files-cli form-field-sets find`

Show Form Field Set.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Form Field Set ID. **Required.** |

### `files-cli form-field-sets create`

Create Form Field Set.

| Flag | Type | Description |
| --- | --- | --- |
| `--user-id` | int64 | User ID.  Provide a value of `0` to operate the current session's user. |
| `--title` | string | Title to be displayed |
| `--workspace-id` | int64 | Workspace ID |
| `--skip-email` | bool | Skip validating form email |
| `--skip-name` | bool | Skip validating form name |
| `--skip-company` | bool | Skip validating company |
| `--form-fields` | []object | (no description) |

### `files-cli form-field-sets update`

Update Form Field Set.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Form Field Set ID. **Required.** |
| `--title` | string | Title to be displayed |
| `--workspace-id` | int64 | Workspace ID |
| `--skip-email` | bool | Skip validating form email |
| `--skip-name` | bool | Skip validating form name |
| `--skip-company` | bool | Skip validating company |
| `--form-fields` | []object | (no description) |

### `files-cli form-field-sets delete`

Delete Form Field Set.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Form Field Set ID. **Required.** |

