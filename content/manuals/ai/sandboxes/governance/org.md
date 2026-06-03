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
to sandboxes across the organization, either to every member or to specific
teams. When organization governance is active, it replaces local `sbx policy`
rules entirely — local rules are no longer evaluated and can't be used to
supplement or override the organization policy.

Admins can manage organization policies through the Admin Console UI or
programmatically using the [Governance API](/reference/api/ai-governance/).

> [!NOTE]
> Sandbox organization governance is available on a separate paid
> subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Create a policy

Manage policies in the [Docker Admin Console](https://app.docker.com/admin)
under **AI governance**. Network and filesystem policies are managed
separately, under **Network access** and **Filesystem access**.

To create a policy:

1. Under **AI governance**, select **Network access** or **Filesystem access**.
2. Select **Create policy** and enter a **Policy name**.
3. Set the **Scope** to **Organization** or **Teams**. If you select **Teams**,
   choose the teams the policy applies to. See
   [Scope policies to teams](#scope-policies-to-teams).
4. Select **Add rule** to add each rule. For rule syntax, see
   [Policy concepts](concepts.md#rule-syntax).

Existing policies are listed with their name, scope, rule count, and last
update. Use the action menu (⋮) to edit or delete a policy.

## Network policies

### Configuring org-level network rules

A network rule takes a network target and an action (allow or deny). You can
add multiple entries at once, one per line.

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

Admins can restrict which paths are mountable with filesystem allow and deny
rules. Each rule takes a path pattern and an action (allow or deny).

For path pattern syntax including the difference between `*` and `**`, see
[Policy concepts](concepts.md#filesystem-rules).

## Scope policies to teams

An organization can have more than one policy, and each policy applies either
to the whole organization or to specific teams. Use scoping to give different
parts of the organization different rules: for example, a permissive policy for
a security research team alongside a stricter default for everyone else.

A policy's [**Scope**](#create-a-policy) controls who it applies to. Set it to
**Organization** to apply the policy to every member, or to **Teams** to apply
it only to members of the teams you select.

### Before you start

Team scoping targets your organization's existing
[teams](/manuals/admin/organization/manage/manage-a-team.md), so a team must
exist before you can scope a policy to it. Create teams and manage their members
in one of two ways:

- Manually, in the Admin Console.
- Automatically, by syncing them from your identity provider with
  [SSO group mapping](/manuals/enterprise/security/single-sign-on/manage.md), so
  that team membership follows your IdP groups.

Because policy assignment follows team membership, you can govern an
organization with thousands of users without per-user configuration. When a
user's team membership changes — whether you edit it directly or it syncs from
your IdP — the policies they receive change with it.

### How org-wide and team-scoped policies combine

A user receives every org-wide policy plus every team-scoped policy for a team
they belong to. The rules from all of these policies are combined and evaluated
together:

- **Allows are additive.** A request is allowed if any of the user's effective
  policies allow it.
- **Denies are absolute.** A request is blocked if any of the user's effective
  policies deny it.

Because deny always wins, a deny rule in an **Organization**-scoped policy acts
as a guardrail that **Teams**-scoped policies can't override. Keep rules that
must apply everywhere in an organization-scoped policy, and use team-scoped
policies to grant extra access to specific teams.

For example, an organization-scoped policy can deny a category of domains for
everyone, while a team-scoped policy grants a research team access to additional
package mirrors. Research-team members receive both policies: they get the extra
access, but the organization-wide deny still applies and can't be undone by the
team policy.

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
