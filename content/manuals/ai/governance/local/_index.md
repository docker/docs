---
title: Local policy
weight: 10
description: Configure local network access rules for sandboxes on your machine.
keywords: docker sandboxes, local policy, network access, allow rules, deny rules, sbx policy
aliases:
  - /ai/sandboxes/security/policy/
  - /ai/sandboxes/governance/access-controls/local/
  - /ai/sandboxes/governance/local/
---

Local policy lets you, the developer, decide which network destinations the
sandboxes on your machine can reach. You manage it with the `sbx policy`
command: you start from a preset and add rules of your own, and the kits your
sandboxes use add the rules their agents need. Local policy covers network
access only, so unless an organization policy applies to you, sandboxes can
mount any host path and use any MCP server you register.

## How local policy works

Every outbound TCP request from a sandbox, including HTTP and SSH, passes
through a proxy on your host, which checks the request against your rules
before it leaves your machine. A rule names a destination and says whether to
allow or deny it. Other protocols are blocked by default. UDP needs the
experimental feature described in
[Allow outbound UDP](network-rules.md#allow-outbound-udp), and ICMP is always
blocked.

Most rules match a hostname, such as `api.github.com`, or an IP address or CIDR
range. For HTTP traffic, a rule can also match the request's method, path, and
query parameters, so you can allow part of an API without allowing all of it. A
rule for other TCP traffic, such as SSH, can name a port, for example
`sbx policy allow network "myhost:22"`. For the full pattern syntax, including
wildcards, see [Network and HTTP rule syntax](../reference/network-rules.md).

## Start from a preset

Local policy starts from a network preset, which sets how much network access
sandboxes have before you add rules of your own. If you haven't chosen a preset,
the CLI prompts you to choose one when you run a sandbox in an interactive
terminal. If a preset is already set, including one carried over from an
earlier installation, the CLI uses it without prompting.

| Preset      | CLI value   | Description                                                                                                                                                                                 |
| ----------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Open        | `allow-all` | All outbound TCP traffic is allowed. Equivalent to adding a wildcard allow rule with `sbx policy allow network "**"`.                                                                       |
| Balanced    | `balanced`  | Default deny, with a baseline allowlist covering AI provider APIs, package managers, OS package repositories, code hosts, container registries, cloud services, and certificate validation. |
| Locked Down | `deny-all`  | No baseline allow rules. Destinations need an allow rule from you or a kit.                                                                                                                 |

The **Balanced** preset's baseline allowlist is a good starting point for most
workflows. After you choose it, run `sbx policy ls` to see exactly which rules
it includes.

To set the preset without the prompt, or to choose again later, see
[Set or reset the preset](network-rules.md#set-or-reset-the-preset).

## Add your own rules

The preset sets your global rules, which apply to every sandbox on your
machine. On top of the preset, you can add your own allow and deny rules,
either globally or scoped to one sandbox by name. For more details, see
[Manage local network rules](network-rules.md).

## Use kits that add rules

Kits, including the built-in agent kits, can add network rules to the sandboxes
that use them, such as the API endpoints an agent needs. For more details, see
[Use kits](../../sandboxes/customize/use-kits.md).

## How rules combine

Without organization governance, your rules and your kits' rules combine: a
destination is allowed if any rule allows it, and a deny rule from any source
takes precedence over an allow rule. Kit rules apply under any preset,
including **Locked Down**, because a preset isn't an explicit deny rule.

How a request that no rule matches is handled depends on your preset. Under
**Open**, every destination is allowed. Under **Balanced** and **Locked Down**,
the request is blocked and asks for your approval instead of being denied
outright, so you can open access to each destination as a sandbox needs it.
For more details, see [Respond to approval requests](approvals.md).

When an organization policy is enforced for you, it changes which of your
local rules apply. For more details, see
[How organization policy relates to local and kit rules](../organization/policies/_index.md#how-organization-policy-relates-to-local-and-kit-rules).

## Try it out

[Try out local policy](try-it-out.md) walks you through an example that allows
and blocks destinations for one sandbox.
