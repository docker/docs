---
title: Local policy
weight: 10
description: Configure local network access rules for sandboxes on your machine.
keywords: docker sandboxes, local policy, network access, allow rules, deny rules, sbx policy
aliases:
  - /ai/sandboxes/security/policy/
  - /ai/sandboxes/governance/access-controls/local/
---

Local policy decides which network destinations the sandboxes on your machine
can reach. You manage it with the `sbx policy` command, and it holds the network
access rules for your machine. You start from a preset, add rules of your own,
and kits add the rules their agents need.

## How local policy works

Local policy covers network access. Rules usually match a hostname, and HTTP
traffic can also be matched by method, path, and query parameters. Non-HTTP
TCP traffic, including SSH, can be allowed with a hostname rule, for example
`sbx policy allow network "myhost:22"`, or with an address-based rule. UDP requires the experimental
feature and policy rules described in [Allow outbound
UDP](network-rules.md#allow-outbound-udp). ICMP is blocked.

For domain patterns, wildcards, and CIDR ranges, see
[Network and HTTP rule syntax](../reference/network-rules.md).

## Start from a preset

Local policy starts from a network preset, which sets how much network access
sandboxes have before you add rules of your own. If you haven't chosen a preset,
the CLI prompts you before it runs a sandbox:

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

The **Balanced** preset's baseline allowlist is a good starting point for most
workflows. Run `sbx policy ls` to see exactly which rules it includes.

To set the preset without the prompt, or to choose again later, see
[Set or reset the preset](network-rules.md#set-or-reset-the-preset).

## Add your own rules

The preset sets your global rules, which apply to every sandbox on your
machine. On top of the preset, you can add your own allow and deny rules,
either globally or scoped to one sandbox by name. For more details, see
[Manage local network rules](network-rules.md).

## Add rules from kits

Kits, including the built-in agent kits, can add network rules to the
sandboxes that use them, such as the API endpoints an agent needs. Kit rules
work like the rules you add yourself, and they apply under any preset,
including **Locked Down** (`deny-all`). A preset isn't an explicit deny rule,
so it doesn't override the allow rules a kit adds. For more details on kits,
see [Use kits](../../customize/use-kits.md).

To list the rules a kit added to a sandbox, see
[Filter the listing](inspect.md#filter-the-listing). To block a destination
that a kit allows, see
[Allow or block a destination for one sandbox](network-rules.md#allow-or-block-a-destination-for-one-sandbox).

## How rules combine

Without organization governance, your rules and your kits' rules combine: a
destination is allowed if any rule allows it, and a deny rule from any source
takes precedence over an allow rule.

How a request that no rule matches is handled depends on your preset. Under
**Open**, every destination is allowed. Under **Balanced** and **Locked Down**,
the request is blocked and asks for your approval instead of being denied
outright, so you can open access to each destination as a sandbox needs it.
For more details, see [Respond to approval requests](approvals.md).

When an organization policy is enforced for you, it changes which of your
local rules apply. For more details, see
[How organization policy relates to local and kit rules](../organization/policies/_index.md#how-organization-policy-relates-to-local-and-kit-rules).
To list the rules that stopped applying, run `sbx policy ls --wide
--include-inactive`. For more details, see
[Show inactive rules](inspect.md#show-inactive-rules).

[Approvals](approvals.md), [inspection](inspect.md), and
[testing](test-monitor.md) work the same way whether a rule comes from you, a
kit, or your organization.

## Try it out

[Try out local policy](try-it-out.md) walks you through an example that allows
and blocks destinations for one sandbox.
