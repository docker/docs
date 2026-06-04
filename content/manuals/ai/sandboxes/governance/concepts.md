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
Policies exist at two levels:

- **Local**: configured per machine using the `sbx policy` CLI. Applies to
  sandboxes on that machine only.
- **Organization**: configured in the Docker Admin Console or via the
  [Governance API](/reference/api/ai-governance/). Applies to sandboxes across
  the organization. An organization can have several policies, each applying
  either org-wide or to specific teams. See [Policy scope](#policy-scope).

When organization governance is active, organization policies replace local
policies entirely. See [Precedence](#precedence).

A **rule** is the unit of access control within a policy. Each rule has:

- **Name**: a human-readable label
- **Actions**: the type of access the rule controls
- **Resources**: the targets the rule matches against
- **Decision**: `allow` or `deny`

Rules are grouped by domain: all rules in a policy must share the same domain,
either `network` or `filesystem`.

## Policy scope

Each organization policy applies either across the whole organization or only
to specific teams:

- Org-wide: with no teams assigned, the policy applies to every member of the
  organization.
- Team-scoped: with one or more teams assigned, the policy applies only to
  members of those teams.

Teams are the same [teams](/manuals/admin/organization/manage/manage-a-team.md)
you manage for your organization; Docker matches a policy's teams against each
user's team membership. Because an organization can mix org-wide and team-scoped
policies, a single user is often subject to several at once. The policies that
apply to a given user are their _effective policies_: every org-wide policy,
plus every team-scoped policy for a team they belong to. See
[Rule evaluation](#rule-evaluation) for how a user's effective policies combine.

## Rule syntax

### Network rules

Network rules use the actions `connect:tcp` and `connect:udp`. Resources are
hostnames, CIDR ranges, or ports.

**Hostname patterns**

| Pattern               | Example           | Matches                                            |
| --------------------- | ----------------- | -------------------------------------------------- |
| Exact hostname        | `example.com`     | `example.com` only, not subdomains                 |
| Single-level wildcard | `*.example.com`   | One subdomain level: `api.example.com`             |
| Multi-level wildcard  | `**.example.com`  | Any depth: `api.example.com`, `v2.api.example.com` |
| Hostname with port    | `example.com:443` | `example.com` on port 443 only                     |

`example.com` and `*.example.com` don't cover each other. Specify both if you
need to match the root domain and its subdomains.

**CIDR ranges**

Both IPv4 and IPv6 notation are supported: `10.0.0.0/8`, `192.168.1.0/24`,
`2001:db8::/32`.

### Filesystem rules

Filesystem rules use the actions `read` and `write`. Resources are host paths
that sandboxes can mount as workspaces.

| Pattern            | Example    | Matches                                                    |
| ------------------ | ---------- | ---------------------------------------------------------- |
| Exact path         | `/data`    | `/data` only                                               |
| Segment wildcard   | `/data/*`  | `/data/project`, one path segment only, not subdirectories |
| Recursive wildcard | `/data/**` | `/data/project`, `/data/project/src`, any depth            |

Use `**` when you intend to match a directory tree recursively. A single `*`
only matches within one path segment and won't cross directory boundaries.
For example, `~/**` matches all paths under the home directory, while `~/*`
matches only direct children of `~`.

## Rule evaluation

When organization governance is active, the rules from all of a user's
[effective policies](#policy-scope) are combined and evaluated together against
each request, following two principles:

- Deny wins: if any rule matches with `decision: deny`, the request is denied,
  regardless of any matching allow rules.
- Default deny: anything an allow rule doesn't match is blocked. Outbound
  network traffic is blocked unless a network rule allows the destination, and a
  host path can't be mounted unless a filesystem rule allows it.

Because every effective policy feeds the same evaluation, allows are additive (a
request is allowed if any effective policy allows it) and denies are absolute (a
request is blocked if any effective policy denies it). A deny rule in an
org-wide policy therefore applies to everyone and can't be overridden by a
team-scoped policy, which makes org-wide deny rules useful as guardrails.

Local rules take no part in this evaluation; see [Precedence](#precedence).

## Precedence

Local and organization policies don't combine. Which one applies depends on
whether your organization has governance enabled:

- No organization governance: local rules and any
  [kit-defined network rules](../customize/kits.md#control-network-access)
  determine what sandboxes can access.
- Organization governance active: organization rules apply across all developer
  machines, and local and kit-defined rules are not evaluated. They still appear
  in `sbx policy ls`, but with an `inactive` status.

When organization governance is active, a user's organization policies are
evaluated together, as described in [Rule evaluation](#rule-evaluation).
