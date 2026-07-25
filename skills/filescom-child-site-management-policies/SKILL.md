---
name: filescom-child-site-management-policies
description: |
  A Child Site Management Policy is a centralized policy defined by a parent site to enforce consistent configurations across child sites.
---

# filescom-child-site-management-policies

A Child Site Management Policy is a centralized policy defined by a parent site to enforce consistent configurations across child sites. These policies allow parent sites to maintain control over specific aspects of their child sites' functionality and appearance.

Non-default policies apply only to the child sites listed in `child_site_ids`, and each child site can be explicitly assigned to only one policy.

One policy can be designated as the default policy. It applies to every child site not explicitly assigned to another policy or listed in its `skip_child_site_ids`, including newly created child sites. Only a default policy can exclude child sites. The `value` field contains the policy configuration data, with the format varying based on the policy type. When a policy is active, its managed configurations are automatically enforced on applicable child sites, and attribute modifications are not permitted.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli child-site-management-policies list`

List Child Site Management Policies.

No flags beyond the global ones.

### `files-cli child-site-management-policies find`

Show Child Site Management Policy.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Child Site Management Policy ID. **Required.** |

### `files-cli child-site-management-policies create`

Create Child Site Management Policy.

| Flag | Type | Description |
| --- | --- | --- |
| `--value` | object | Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation. |
| `--skip-child-site-ids` | []int64 | IDs of child sites excluded from this default policy. |
| `--child-site-ids` | []int64 | IDs of child sites explicitly assigned to this non-default policy. |
| `--default-policy` | bool | Whether this policy applies to child sites not explicitly assigned to another policy. |
| `--policy-type` | enum | Type of policy.  Valid values: `settings`. One of: `settings`. **Required.** |
| `--name` | string | Name for this policy. |
| `--description` | string | Description for this policy. |

### `files-cli child-site-management-policies update`

Update Child Site Management Policy.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Child Site Management Policy ID. **Required.** |
| `--value` | object | Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation. |
| `--skip-child-site-ids` | []int64 | IDs of child sites excluded from this default policy. |
| `--child-site-ids` | []int64 | IDs of child sites explicitly assigned to this non-default policy. |
| `--default-policy` | bool | Whether this policy applies to child sites not explicitly assigned to another policy. |
| `--policy-type` | enum | Type of policy.  Valid values: `settings`. One of: `settings`. |
| `--name` | string | Name for this policy. |
| `--description` | string | Description for this policy. |

### `files-cli child-site-management-policies delete`

Delete Child Site Management Policy.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Child Site Management Policy ID. **Required.** |

