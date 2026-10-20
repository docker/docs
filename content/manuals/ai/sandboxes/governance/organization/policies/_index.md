---
title: Policies
weight: 30
description: How organization policy scope, rule evaluation, and precedence work for Docker Sandboxes.
keywords: docker sandboxes, organization policy, policy scope, rule evaluation, precedence, governance
aliases:
  - /ai/sandboxes/governance/concepts/
---

## Policy scope

Each organization policy applies either across the whole organization or only
to specific teams:

- Org-wide: with no teams assigned, the policy applies to every member of the
  organization.
- Team-scoped: with one or more teams assigned, the policy applies only to
  members of those teams.

Teams are the same [teams](/manuals/accounts/organization/manage/manage-a-team.md)
you manage for your organization; Docker matches a policy's teams against each
user's team membership. Because an organization can mix org-wide and team-scoped
policies, a single user is often subject to several at once. The policies that
apply to a given user are their _effective policies_: every org-wide policy,
plus every team-scoped policy for a team they belong to. See
[Rule evaluation](#rule-evaluation) for how a user's effective policies combine.

## How scoped policies combine

A user is governed by all of their
[effective policies](#policy-scope): every org-wide policy, plus
the team-scoped policies for the teams they belong to. Use org-wide policies
for guardrails that must apply everywhere, and team-scoped policies for access
that only some teams need.

For precedence between local and organization policies, and for how allow and
deny rules combine, see [Precedence](#precedence).

## Rule evaluation

When organization governance is active, the rules from all of a user's
[effective policies](#policy-scope) are combined and evaluated together against
each request, following two principles:

- Deny wins: if any rule matches with `decision: deny`, the request is denied,
  regardless of any matching allow rules.
- Default deny: anything an allow rule doesn't match is blocked. Outbound
  network traffic is blocked unless a network rule allows the destination, and a
  host path can't be mounted unless a filesystem rule allows it. MCP activity is
  blocked unless an MCP `permit` allows it.

Because every effective policy feeds the same evaluation, allows are additive (a
request is allowed if any effective policy allows it) and denies are absolute (a
request is blocked if any effective policy denies it). A deny rule in an
org-wide policy therefore applies to everyone and can't be overridden by a
team-scoped policy, which makes org-wide deny rules useful as guardrails.

Local and kit-defined allow rules take no part in this evaluation. Deny rules
from those sources do still apply. See [Precedence](#precedence).

A request that an approval-required policy allows produces a third outcome.
Rather than being allowed outright, it's held back until the developer confirms
the destination, and the confirmation governs later requests to it. This holds
even when another policy allows the same request without requiring approval. A
matching deny still wins, so a denied destination is blocked without asking.
See [Approval-required access](network.md#approval-required-access).

## Precedence

What applies depends on whether your organization has governance enabled:

- No organization governance: local rules and any kit-defined network rules
  determine what sandboxes can access.
- Organization governance active: organization policy determines what access can
  be granted. Only organization allow rules grant access, so local and
  kit-defined allow rules are inactive and can't expand what the organization
  permits. Deny rules apply from every source, so a local or kit-defined deny
  can still restrict access further.

For kit-defined rules, see
[Network policies](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/network-policy@1.md)
in the kit specification.

Precedence is decided by a rule's decision rather than its source:

| Rule                | Evaluated under organization governance |
| ------------------- | --------------------------------------- |
| Organization allow  | Yes                                     |
| Organization deny   | Yes                                     |
| Local allow         | No                                      |
| Local deny          | Yes                                     |
| Kit-defined allow   | No                                      |
| Kit-defined deny    | Yes                                     |

Local and kit-defined rules cover network access only, so a deny that layers on
top of organization policy is always a network deny. `sbx policy ls` hides
inactive rules by default. See
[Monitoring](../../local/inspect.md#showing-inactive-rules) for how
to list them.

A local deny takes precedence over an organization approval requirement as
well, so the request is blocked and no approval is requested. Rules that a
developer gains by approving a request are the one exception to local allow
rules being inactive, because they record an answer to the organization's own
approval requirement rather than granting new access. See
[Approval-required access](network.md#approval-required-access).

When organization governance is active, a user's organization policies are
evaluated together, as described in [Rule evaluation](#rule-evaluation).

## Troubleshooting

### Policy changes not taking effect

After updating organization policies, changes take up to 5 minutes to
propagate to developer machines. To apply changes immediately, users can run
`sbx policy reset`, which stops the daemon and forces it to pull the latest
organization policies on the next `sbx` command.

> [!WARNING]
> `sbx policy reset` deletes all locally configured policy rules, including any
> destinations the developer has approved under an
> [approval-required policy](network.md#approval-required-access). Those
> destinations are requested again the next time a sandbox reaches them. The
> command prompts for confirmation before proceeding.

#### Enforcement timing by policy type

Policy types differ in when a change takes effect after it reaches the
developer machine:

- Network policy is evaluated on every outbound request. Once a policy
  change has synced to the developer's machine (up to 5 minutes), it applies
  immediately to subsequent requests. HTTP rules are evaluated per request in
  the same way.

- An approval requirement applies from the point the policy change syncs.
  Destinations a developer already approved stay reachable, because the
  approval is recorded on the developer's machine. To withdraw one, add a deny
  rule. A deny takes precedence over a recorded approval.

- Filesystem policy is only checked when a workspace is mounted — that
  is, when a sandbox is created. Once a sandbox is running, changing the
  filesystem policy has no effect on that sandbox. The sandbox continues to
  access the previously allowed path until it is removed and a new one is
  created.

- MCP registration policy is evaluated when a server is registered with
  `sbx mcp add`. Changing registration rules doesn't remove existing
  registrations or stop an already-loaded server by itself.

- MCP use-time policy is evaluated by the MCP gateway when a sandbox makes a
  governed MCP request, such as a tool call, resource read, prompt retrieval,
  or built-in gateway tool call. Once a policy change has synced, use-time
  rules apply to subsequent governed MCP requests through the gateway.

To apply a filesystem policy change immediately, remove the running sandbox
and create a new one. To prevent use of an MCP server that is already registered
or loaded, add use-time rules for the registered server name. For examples, see
[Withdraw server access](mcp.md#withdraw-server-access).
