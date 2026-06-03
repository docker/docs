---
title: Organization policy
linkTitle: Org policy
weight: 20
description: Centrally manage sandbox network and filesystem policies for your organization.
keywords: docker sandboxes, governance, organization policy, AI governance, admin console, network access, filesystem access
aliases:
  - /ai/sandboxes/security/governance/
---

[Local policies](local.md) give individual developers control over what their
sandboxes can access. Organization policy moves that control to the admin level:
rules defined in the [Docker Admin Console](https://app.docker.com/admin) apply
uniformly to every sandbox in the organization. When organization governance is
active, it replaces local `sbx policy` rules entirely — local rules are no
longer evaluated and can't be used to supplement or override the organization
policy.

Admins can manage organization policies through the Admin Console UI or
programmatically using the [Governance API](/reference/api/ai-governance/).

> [!NOTE]
> Sandbox organization governance is available on a separate paid
> subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Network policies

### Configuring org-level network rules

Define network allow and deny rules in the Admin Console under
**AI governance > Network access**. Each rule takes a network target and an
action (allow or deny). You can add multiple entries at once, one per line.

For the full syntax reference (exact hostnames, wildcard subdomains, port
suffixes, and CIDR ranges), see [Policy concepts](concepts.md#network-rules).

When organization governance is active, local network rules are not evaluated.
The organization policy is the only policy in effect. Local rules still appear
in `sbx policy ls` but with an `inactive` status. See [Monitoring](monitoring.md)
for how to read the rule view.

## Filesystem policies

Filesystem policies control which host paths a sandbox can mount as
workspaces. By default, sandboxes can mount any directory the user has
access to.

Admins can restrict which paths are mountable by defining filesystem allow
and deny rules in the Admin Console under **AI governance > Filesystem
access**. Each rule takes a path pattern and an action (allow or deny).

For path pattern syntax including the difference between `*` and `**`, see
[Policy concepts](concepts.md#filesystem-rules).

## Scope policies to teams

An organization can have more than one policy, and each policy applies either
to the whole organization or to specific teams. Use scoping to give different
parts of the organization different rules: for example, a permissive policy for
a security research team alongside a stricter default for everyone else.

When you create or edit a policy in the Admin Console, the **Scope** setting
controls who it applies to:

- Leave the team list empty to make the policy org-wide. It applies to every
  member of the organization.
- Add one or more teams to scope the policy to those teams. It applies only to
  members of the listed teams.

Teams are the same [teams](/manuals/admin/organization/manage/manage-a-team.md)
you manage for your organization, so assignment follows existing team
membership. This lets you manage policies for an organization with thousands of
users without configuring anything per user. When team membership changes in
your identity provider, the policies a user receives change with it.

### How org-wide and team-scoped policies combine

A user receives every org-wide policy plus every team-scoped policy for a team
they belong to. The rules from all of these policies are combined and evaluated
together:

- **Allows are additive.** A request is allowed if any of the user's effective
  policies allow it.
- **Denies are absolute.** A request is blocked if any of the user's effective
  policies deny it.

Because deny always wins, a deny rule in an org-wide policy acts as a guardrail
that team-scoped policies can't override. Use org-wide policies for rules that
must apply everywhere, and team-scoped policies to grant additional access to
specific teams.

## Precedence

Deny rules beat allow rules. If a resource matches both an allow and a deny, in
the same policy or across a user's effective policies, it's blocked regardless
of specificity. Outbound traffic is blocked unless a rule allows it.

When organization governance is active, local rules are not evaluated. Only
organization rules set in the Admin Console determine what is allowed or
denied, and they can't be supplemented or overridden from a developer's machine.
The same model applies to filesystem policies: organization rules replace local
behavior entirely.

To unblock a domain when organization governance is active, update the rule in
the Admin Console or via the [API](/reference/api/ai-governance/). Without
organization governance, remove the local rule with `sbx policy rm`.

## Troubleshooting

### Policy changes not taking effect

After updating organization policies in the Admin Console, changes take up
to 5 minutes to propagate to developer machines. To apply changes
immediately, users can run `sbx policy reset`, which stops the daemon and
forces it to pull the latest organization policies on the next `sbx`
command.

> [!WARNING]
> `sbx policy reset` deletes all locally configured policy rules. The command
> prompts for confirmation before proceeding.

### Sandbox cannot mount workspace

If a sandbox fails to mount with a `mount policy denied` error, verify that
the filesystem allow rule in the Admin Console uses `**` rather than `*`. A
single `*` doesn't match across directory separators.
