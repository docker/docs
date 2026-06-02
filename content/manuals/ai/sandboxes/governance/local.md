---
title: Local policy
weight: 10
description: Configure network access rules for sandboxes on your local machine.
keywords: docker sandboxes, local policy, network access, allow rules, deny rules, sbx policy
aliases:
  - /ai/sandboxes/security/policy/
---

The `sbx policy` command manages network access rules on your local machine.
Rules apply to all sandboxes on the machine when you use the global scope, or
to a single sandbox when scoped by name.

Local rules apply only when your organization doesn't enforce governance:

- **No org governance**: local rules fully control what sandboxes can access.
- **Org governance active**: the organization policy replaces local policy.
  Local rules are inactive, and `sbx policy allow` and `sbx policy deny` have
  no effect. Local rules still appear in `sbx policy ls` with an `inactive`
  status.

See [Organization policy](org.md) for how organization governance works.

For domain patterns, wildcards, CIDR ranges, and filesystem path syntax, see
[Policy concepts](concepts.md#rule-syntax).

## Default preset

The only way traffic can leave a sandbox is through an HTTP/HTTPS proxy on
your host, which enforces access rules on every outbound request. Non-HTTP TCP
traffic, including SSH, can be allowed by adding a policy rule for the
destination IP and port (for example, `sbx policy allow network -g
"10.1.2.3:22"`). UDP and ICMP are blocked at the network layer and can't be
unblocked with policy rules.

On first start, and after running `sbx policy reset`, the daemon prompts you
to choose a network preset:

```plaintext
Choose a default network policy:

     1. Open         — All network traffic allowed, no restrictions.
     2. Balanced     — Default deny, with common dev sites allowed.
     3. Locked Down  — All network traffic blocked unless you allow it.

  Use ↑/↓ to navigate, Enter to select, or press 1–3.
```

| Preset | Description |
| ------ | ----------- |
| Open | All outbound traffic is allowed. Equivalent to adding a wildcard allow rule with `sbx policy allow network -g "**"`. |
| Balanced | Default deny, with a baseline allowlist covering AI provider APIs, package managers, code hosts, container registries, and common cloud services. |
| Locked Down | All outbound traffic is blocked, including model provider APIs (for example, `api.anthropic.com`). You must explicitly allow everything you need. |

The **Balanced** preset's baseline allowlist is a good starting point for most
workflows. Run `sbx policy ls` to see exactly which rules it includes.

> [!NOTE]
> If your organization manages sandbox policies centrally, organization rules
> take precedence over the preset you select here. See
> [Organization policy](org.md).

### Non-interactive environments

In non-interactive environments such as CI pipelines or headless servers, the
interactive prompt can't be displayed. Use `sbx policy set-default` to set the
preset before running any other `sbx` commands:

```console
$ sbx policy set-default balanced
```

Available values are `allow-all`, `balanced`, and `deny-all`.

## Managing rules

Use [`sbx policy allow`](/reference/cli/sbx/policy/allow/) and
[`sbx policy deny`](/reference/cli/sbx/policy/deny/) to add or restrict access
on top of the active preset. Changes take effect immediately. Pass `-g` to
apply a rule globally to all sandboxes:

```console
$ sbx policy allow network -g api.anthropic.com
$ sbx policy deny network -g ads.example.com
```

Pass a sandbox name to scope a rule to one sandbox:

```console
$ sbx policy allow network my-sandbox api.example.com
$ sbx policy deny network my-sandbox ads.example.com
```

Specify multiple hosts in one command with a comma-separated list:

```console
$ sbx policy allow network -g "api.anthropic.com,*.npmjs.org,*.pypi.org"
```

Remove a rule by resource or by rule ID:

```console
$ sbx policy rm network -g --resource ads.example.com
$ sbx policy rm network -g --id 2d3c1f0e-4a73-4e05-bc9d-f2f9a4b50d67
```

To remove a sandbox-scoped rule, include the sandbox name:

```console
$ sbx policy rm network my-sandbox --resource api.example.com
```

To inspect which rules are active and where they come from, use
`sbx policy ls`. See [Monitoring](monitoring.md).

### Resetting

To remove all custom rules and start fresh with a new preset, use
`sbx policy reset`:

```console
$ sbx policy reset
```

This deletes the local policy store and stops the daemon. When the daemon
restarts on the next command, you are prompted to choose a new preset. Running
sandboxes stop when the daemon shuts down. Pass `--force` to skip the
confirmation prompt:

```console
$ sbx policy reset --force
```

## Troubleshooting

### Local rules have no effect

If rules you add with `sbx policy allow` or `sbx policy deny` don't change
sandbox behavior, your organization likely has governance enabled. Run `sbx
policy ls` to check: if the output starts with a `Governance: managed by <org>`
header, org governance is active. When it's active, the organization policy
replaces local policy, so your rules appear with `inactive` status and have no
effect.

Organization policy can't be supplemented from your machine. To change what
your sandboxes can access, ask your admin to update the organization policy in
the Admin Console.

### A domain is still blocked after adding an allow rule

If a domain remains blocked after you add a local allow rule, your organization
likely enforces governance, which makes local rules inactive. Run `sbx policy
ls` to check whether org governance is active and whether your rule shows an
`inactive` status. If so, the block can only be lifted by updating the org
policy in the Admin Console or via the [API](/reference/api/ai-governance/).
