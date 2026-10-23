---
title: Organization policies
linkTitle: Policies
weight: 30
description: How organization policy scope, rule evaluation, and precedence work for Docker Sandboxes.
keywords: docker sandboxes, organization policy, policy scope, rule evaluation, precedence, governance
aliases:
  - /ai/sandboxes/governance/concepts/
---

{{< summary-bar feature_name="AI Governance" >}}

Organization policies control what sandboxes can access for organization
members. When a developer's sandbox makes a request, three questions decide
the outcome: which of your policies apply to that developer, what their rules
say about the request, and how your policy relates to the developer's own local
and kit rules. This page answers each in turn, then explains when your changes
reach developers. To create a policy, see
[Create and manage organization policies](create-manage.md).

## What a policy contains

An organization policy is a named collection of rules. Each rule has a name,
the actions it controls, the resources it matches, and a decision of `allow` or
`deny`. Each policy covers one kind of access:

- [Organization network access policies](network.md) control the destinations
  sandboxes can reach.
- [Organization filesystem access policies](filesystem.md) control the host
  paths sandboxes can mount.
- [Organization MCP access policies](mcp.md) control MCP activity. They use
  Cedar statements in the `MCP` namespace instead of rules.

## Which policies apply to a developer

Organization policies apply only to members who hold an AI Governance license.
Each policy applies to those members either across the whole organization or
only on specific teams:

- Org-wide: with no teams assigned, the policy applies to every licensed member
  of the organization.
- Team-scoped: with one or more teams assigned, the policy applies only to
  licensed members of those teams.

Teams are the same [teams](/manuals/accounts/organization/manage/manage-a-team.md)
you manage for your organization. Docker matches a policy's teams against each
user's team membership. Because an organization can mix org-wide and
team-scoped policies, a single user is often subject to several at once. The
policies that apply to a given user are their effective policies: every
org-wide policy, plus every team-scoped policy for a team they belong to.

Use org-wide policies for guardrails that must apply everywhere, and
team-scoped policies for access that only some teams need. The next section
explains why org-wide deny rules work well as guardrails.

## How rules decide a request

When an organization policy is enforced for you, the rules from all of your
effective policies are combined and evaluated together against each request,
following two principles:

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

A request that an approval-required policy allows produces a third outcome: it's
blocked, and the developer is asked to approve the destination. After they
approve it, later requests to that destination succeed. For more details, see
[Approval-required access](network.md#approval-required-access).

### How organization policy relates to local and kit rules

Your organization's rules aren't the only ones on a developer's machine.
Developers also have their own local rules, and kits can add rules to the
sandboxes that use them. What applies depends on whether an organization
policy is enforced for the developer:

- No organization policy: local rules and any kit-defined network rules
  determine what sandboxes can access.
- Organization policy enforced: organization policy determines what access can
  be granted. Only organization allow rules grant access, so local and
  kit-defined allow rules are inactive and can't expand what the organization
  permits. Deny rules apply from every source, so a local or kit-defined deny
  can still restrict access further.

The following table summarizes which rules are evaluated under organization
governance:

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
inactive rules by default. For more details on how to list them, see
[Show inactive rules](../../local/inspect.md#show-inactive-rules).

A local deny takes precedence over an organization approval requirement as
well, so the request is blocked and no approval is requested. Rules that a
developer gains by approving a request are the one exception to local allow
rules being inactive, because they record an answer to the organization's own
approval requirement rather than granting new access. For more details, see
[Approval-required access](network.md#approval-required-access).

## When policy changes reach developers

Policy changes reach developer machines within 5 minutes. When a synced change
takes effect also depends on the policy type. For more details, see
[Organization network access policies](network.md#when-changes-take-effect),
[Organization filesystem access policies](filesystem.md#when-changes-take-effect), and
[Organization MCP access policies](mcp.md#when-changes-take-effect).

To apply changes right away on a developer machine, or if a change doesn't seem
to apply, see
[Organization policy changes aren't taking effect](../../../troubleshooting.md#organization-policy-changes-arent-taking-effect).
