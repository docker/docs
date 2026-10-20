---
title: Local policy
weight: 10
description: Configure local network access rules for sandboxes on your machine.
keywords: docker sandboxes, local policy, network access, allow rules, deny rules, sbx policy
aliases:
  - /ai/sandboxes/security/policy/
  - /ai/sandboxes/governance/access-controls/local/
---

The `sbx policy` command manages the local policy on your machine. The local
policy contains network access rules. Rules apply to all sandboxes on the
machine when you use the global scope, or to a single sandbox when scoped by
name.

When an organization policy is enforced for you, your local allow rules stop
applying and your local deny rules still apply. See
[Precedence](../organization/policies/_index.md#precedence). To list the rules
that stopped applying, run `sbx policy ls --include-inactive`. See
[Showing inactive rules](inspect.md#showing-inactive-rules).

For domain patterns, wildcards, and CIDR ranges, see
[Network and HTTP rule syntax](../reference/network-rules.md).

Outbound TCP traffic passes through a proxy on your host, which enforces access
rules on every connection. Non-HTTP TCP traffic, including SSH, can be allowed
with a hostname rule (for example, `sbx policy allow network "myhost:22"`) or an
address-based rule. UDP requires the experimental feature and policy rules
described in [Allow outbound UDP](network-rules.md#allow-outbound-udp). ICMP is blocked.

Without organization governance, how a request that no rule matches is handled
depends on your preset. Under **Open**, every destination is allowed. Under **Balanced** and **Locked Down**,
the request is blocked and asks for your approval instead of being denied
outright, so you can open access to each destination as a sandbox needs it. See
[Respond to approval requests](approvals.md).
