---
name: filescom-child-site-management-policies
description: |
  A Child Site Management Policy is a centralized policy defined by a parent site to enforce consistent configurations across child sites.
---

# filescom-child-site-management-policies

A Child Site Management Policy is a centralized policy defined by a parent site to enforce consistent configurations across child sites. These policies allow parent sites to maintain control over specific aspects of their child sites' functionality and appearance.

Policies can be applied to all child sites, or specific sites can be exempted from policy management by adding their site ID to the `skip_child_site_ids` parameter.

The `value` field contains the policy configuration data, with the format varying based on the policy type. When a policy is active, its managed configurations are automatically enforced on applicable child sites, and attribute modifications are not permitted.

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
| `--skip-child-site-ids` | []int64 | IDs of child sites that this policy has been exempted from. If `skip_child_site_ids` is empty, the policy will be applied to all child sites. To apply a policy to a child site that has been exempted, remove it from `skip_child_site_ids` or set it to an empty array (`[]`). |
| `--policy-type` | enum | Type of policy.  Valid values: `settings`. One of: `settings`. **Required.** |
| `--name` | string | Name for this policy. |
| `--description` | string | Description for this policy. |

### `files-cli child-site-management-policies update`

Update Child Site Management Policy.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Child Site Management Policy ID. **Required.** |
| `--value` | object | Policy configuration data. Attributes differ by policy type. For more information, refer to the Value Hash section of the developer documentation. |
| `--skip-child-site-ids` | []int64 | IDs of child sites that this policy has been exempted from. If `skip_child_site_ids` is empty, the policy will be applied to all child sites. To apply a policy to a child site that has been exempted, remove it from `skip_child_site_ids` or set it to an empty array (`[]`). |
| `--policy-type` | enum | Type of policy.  Valid values: `settings`. One of: `settings`. |
| `--name` | string | Name for this policy. |
| `--description` | string | Description for this policy. |

### `files-cli child-site-management-policies delete`

Delete Child Site Management Policy.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Child Site Management Policy ID. **Required.** |

