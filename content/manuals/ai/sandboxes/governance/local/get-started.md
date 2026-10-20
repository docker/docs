---
title: Get started with local policy
linkTitle: Get started
weight: 10
description: Choose and reset the default network preset for local sandbox policy.
keywords: docker sandboxes, local policy, network preset, sbx policy init, sbx policy reset
---

## Default preset

If you haven't chosen a default preset, the CLI prompts you before it runs a
sandbox. Running `sbx policy reset` clears the preset and prompts you to choose
again:

```plaintext
Initialize the global network policy for your sandboxes:

  Applies to all sandboxes, current and future — change it later with
  "sbx policy allow/deny/rm". Kits, including built-in agent kits, may
  also add per-sandbox rules.

     1. Open         — All network traffic allowed, no restrictions.
  ❯  2. Balanced     — Default deny, with common dev sites allowed.
     3. Locked Down  — All network traffic blocked unless you allow it.

  Use ↑/↓ or 1–3 to navigate, Enter to confirm, Esc to cancel.
```

| Preset      | Description                                                                                                                                       |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| Open        | All outbound TCP traffic is allowed. Equivalent to adding a wildcard allow rule with `sbx policy allow network "**"`. |
| Balanced    | Default deny, with a baseline allowlist covering AI provider APIs, package managers, code hosts, container registries, and common cloud services. |
| Locked Down | No baseline allow rules. Destinations need an allow rule from you or a kit. |

Presets initialize the global policy. Built-in agent kits and other kits can
add per-sandbox allow rules, including under **Locked Down** (`deny-all`). The
preset isn't an explicit deny rule that overrides those allowances.

To inspect the rules a kit adds to a sandbox, run:

```console
$ sbx policy ls my-sandbox --source kit --type network --wide
```

To block a destination allowed by a kit, add an explicit deny rule:

```console
$ sbx policy deny network --sandbox my-sandbox openrouter.ai
```

Deny rules take precedence over allow rules. See
[Policy precedence](../organization/policies/_index.md#precedence).

The **Balanced** preset's baseline allowlist is a good starting point for most
workflows. Run `sbx policy ls` to see exactly which rules it includes.

> [!NOTE]
> If your organization manages sandbox policies centrally, organization rules
> take precedence over the preset you select here. See
> [Organization policies](../organization/_index.md).

### Non-interactive environments

In non-interactive environments such as CI pipelines or headless servers, the
interactive prompt can't be displayed. Use `sbx policy init` to set the
preset before running any other `sbx` commands:

```console
$ sbx policy init balanced
```

Available values are `allow-all`, `balanced`, and `deny-all`.

## Resetting

To remove all custom rules and start fresh with a new preset, use
`sbx policy reset`:

```console
$ sbx policy reset
```

This deletes the local policy store, restarts the daemon, and prompts you to
choose a new preset. Running sandboxes stop when the daemon shuts down. Pass
`--force` to skip the confirmation prompt:

```console
$ sbx policy reset --force
```
