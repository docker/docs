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

Use org-wide policies for guardrails that must apply everywhere, and
team-scoped policies for access that only some teams need.

## Rule evaluation

When an organization policy is enforced for you, the rules from all of your
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

A request that an approval-required policy allows produces a third outcome.
Rather than being allowed outright, it's held back until the developer confirms
the destination, and the confirmation governs later requests to it. This holds
even when another policy allows the same request without requiring approval. A
matching deny still wins, so a denied destination is blocked without asking.
See [Approval-required access](network.md#approval-required-access).

### Precedence

What applies depends on whether an organization policy is enforced for you:

- No organization policy: local rules and any kit-defined network rules
  determine what sandboxes can access.
- Organization policy enforced: organization policy determines what access can
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
[Showing inactive rules](../../local/inspect.md#showing-inactive-rules) for how
to list them.

A local deny takes precedence over an organization approval requirement as
well, so the request is blocked and no approval is requested. Rules that a
developer gains by approving a request are the one exception to local allow
rules being inactive, because they record an answer to the organization's own
approval requirement rather than granting new access. See
[Approval-required access](network.md#approval-required-access).

## Troubleshooting

### Policy changes not taking effect

After updating organization policies, changes take up to 5 minutes to
propagate to developer machines. To apply them right away on a developer
machine, run `sbx policy refresh`, which requires sbx 0.49.0 or later. It
fetches and applies the latest organization policy without restarting the
daemon or stopping running sandboxes.

The refreshed policy applies to new connections only. A connection that a
sandbox opened while the previous policy allowed it stays open even if the new
policy denies it. Restart the sandbox to close those connections.

When a synced change takes effect also depends on the policy type. See
[Network access policies](network.md#when-changes-take-effect),
[Filesystem access policies](filesystem.md#when-changes-take-effect), and
[MCP access policies](mcp.md#when-changes-take-effect).
