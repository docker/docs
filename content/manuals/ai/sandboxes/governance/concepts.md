---
title: Policy concepts
weight: 5
description: The resource model, rule syntax, and evaluation logic behind Docker sandbox governance.
keywords: docker sandboxes, policy concepts, rule syntax, network rules, filesystem rules, precedence, rule evaluation
---

## Resource model

Docker sandbox governance is built around two resource types: **policies** and
**rules**.

A **policy** is a named collection of rules that controls sandbox access.
Policies exist at different levels:

- **Local**: configured per machine using the `sbx policy` CLI. Applies to
  sandboxes on that machine only.
- **Organization**: configured in the Docker Admin Console or via the
  [Governance API](/reference/api/ai-governance/). Applies uniformly across every sandbox in the
  organization.
- **Team**: applies to sandboxes used by a specific team within an
  organization. Coming soon.

When multiple levels are active, organization policies take precedence over
local policies. See [Precedence](#precedence).

A **rule** is the unit of access control within a policy. Each rule has:

- **Name**: a human-readable label
- **Actions**: the type of access the rule controls
- **Resources**: the targets the rule matches against
- **Decision**: `allow` or `deny`

Rules are grouped by domain: all rules in a policy must share the same domain,
either `network` or `filesystem`.

## Rule syntax

### Network rules

Network rules use the actions `connect:tcp` and `connect:udp`. Resources are
hostnames, CIDR ranges, or ports.

**Hostname patterns**

| Pattern | Example | Matches |
| ------- | ------- | ------- |
| Exact hostname | `example.com` | `example.com` only, not subdomains |
| Single-level wildcard | `*.example.com` | One subdomain level: `api.example.com` |
| Multi-level wildcard | `**.example.com` | Any depth: `api.example.com`, `v2.api.example.com` |
| Hostname with port | `example.com:443` | `example.com` on port 443 only |

`example.com` and `*.example.com` don't cover each other. Specify both if you
need to match the root domain and its subdomains.

**CIDR ranges**

Both IPv4 and IPv6 notation are supported: `10.0.0.0/8`, `192.168.1.0/24`,
`2001:db8::/32`.

### Filesystem rules

Filesystem rules use the actions `read` and `write`. Resources are host paths
that sandboxes can mount as workspaces.

| Pattern | Example | Matches |
| ------- | ------- | ------- |
| Exact path | `/data` | `/data` only |
| Segment wildcard | `/data/*` | `/data/project`, one path segment only, not subdirectories |
| Recursive wildcard | `/data/**` | `/data/project`, `/data/project/src`, any depth |

Use `**` when you intend to match a directory tree recursively. A single `*`
only matches within one path segment and won't cross directory boundaries.
For example, `~/**` matches all paths under the home directory, while `~/*`
matches only direct children of `~`.

## Rule evaluation

All rules in a policy are evaluated against every request. The outcome follows
two principles:

**Deny wins.** If any rule matches with `decision: deny`, the request is
denied regardless of any matching allow rules.

**Default deny.** Outbound traffic is blocked unless an explicit allow rule
matches.

These principles apply within whichever policy is active. When organization
governance is active, only organization rules are evaluated; local rules have
no effect.

## Precedence

Local and organization policies don't combine. Which one applies depends on
whether your organization has governance enabled:

- **No organization governance**: local rules determine what sandboxes can
  access.
- **Organization governance active**: organization rules apply across all
  developer machines, and local rules are not evaluated. Local rules still
  appear in `sbx policy ls`, but with an `inactive` status.

Within the active policy, deny rules beat allow rules regardless of specificity
or order.
