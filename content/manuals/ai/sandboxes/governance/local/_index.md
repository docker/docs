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

Local policy interacts with organization governance as follows:

- **No org governance**: the local policy controls what sandboxes can access.
- **Org governance active**: only organization allow rules grant access, so
  local allow rules are inactive and can't expand what the organization permits.
  Local deny rules are still evaluated, so you can restrict access further than
  the organization policy does. To list inactive rules, run
  `sbx policy ls --include-inactive`. See
  [Monitoring](inspect.md#showing-inactive-rules).

See [Organization policies](../organization/_index.md) for how organization governance
works.

For domain patterns, wildcards, CIDR ranges, and filesystem path syntax, see
[Policy concepts](../reference/network-rules.md).

Outbound TCP traffic passes through a proxy on your host, which enforces access
rules on every connection. Non-HTTP TCP traffic, including SSH, can be allowed
with a hostname rule (for example, `sbx policy allow network "myhost:22"`) or an
address-based rule. UDP requires the experimental feature and policy rules
described in [Allow outbound UDP](network-rules.md#allow-outbound-udp). ICMP is blocked.

Under **Balanced** and **Locked Down**, a sandbox request that no rule matches
is blocked and asks for your approval instead of being denied outright, so you
can open access to each destination as a sandbox needs it. See
[Approval-required access](approvals.md).

Without organization governance, a request with no matching allow or deny rule
also asks for approval rather than being denied outright, so access opens up as
the developer approves each destination.
