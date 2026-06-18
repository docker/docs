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
rules defined in **Admin Console** apply to sandboxes across the organization,
either to every member or to specific teams. When organization governance is active, it replaces local `sbx policy`
rules entirely — local rules are no longer evaluated and can't be used to
supplement or override the organization policy.

Admins can manage organization policies through the Admin Console UI or
programmatically using the [Governance API](/reference/api/ai-governance/).

By default, only organization
[owners](/manuals/enterprise/security/roles-and-permissions/core-roles.md) can
view and manage AI Governance policies. To let someone other than an owner
manage policies, create a
[custom role](/manuals/enterprise/security/roles-and-permissions/custom-roles.md)
with the **Governance** permissions and assign it to a user or team.

> [!NOTE]
> Sandbox organization governance is available on a separate paid
> subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Create a policy

Manage policies under **Admin Console**, a section in the left-hand navigation
of [Docker Home](https://app.docker.com). Network and filesystem policies are
managed separately, under **Network access** and **Filesystem access**.

To create a policy:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Select **Admin Console**, then **AI governance**.
1. Select **Network access** or **Filesystem access**, then **Create policy**.
1. Enter a **Policy name**.
1. Set the **Scope** to **Organization** or **Teams**. If you select **Teams**,
   choose the teams the policy applies to. See
   [Scope policies to teams](#scope-policies-to-teams).
1. Select **Add rule** to add each rule. For rule syntax, see
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
The organization policy is the only policy in effect. `sbx policy ls` hides
these inactive local rules by default. See
[Monitoring](monitoring.md#showing-inactive-rules) for how to list them and read
the rule view.

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
to the whole organization or to specific teams. Scoping lets you apply different
rules to different parts of the organization.

A policy's [**Scope**](#create-a-policy) controls who it applies to. Set it to
**Organization** to apply the policy to every member, or to **Teams** to apply
it only to members of the teams you select.

### Before you start

Team scoping targets your organization's existing
[teams](/manuals/admin/organization/manage/manage-a-team.md), so a team must
exist before you can scope a policy to it. Create teams and manage their members
in one of two ways:

- Manually, in the Admin Console.
- Automatically, by using
  [group mapping](/manuals/enterprise/security/provisioning/scim/group-mapping.md)
  to synchronize your identity provider's groups with the teams in your
  organization. Group mapping creates teams that don't already exist and keeps
  their membership in step with your IdP groups.

Because policies apply by team, a user's policies update automatically as their
team membership changes, including changes synced from your IdP.

### How org-wide and team-scoped policies combine

A user is governed by all of their
[effective policies](concepts.md#policy-scope) at once — every org-wide policy,
plus the team-scoped policies for the teams they belong to. The rules combine
into a single evaluation in which deny always wins, so a team-scoped policy can
grant access on top of the org-wide policies but can't loosen a restriction they
impose. For the full evaluation model, see
[Rule evaluation](concepts.md#rule-evaluation).

This makes org-wide policies the natural home for guardrails that must hold
everywhere. For example, an org-wide policy can deny a category of domains for
all members, while a team-scoped policy grants a research team access to extra
package mirrors. Research-team members get the extra access, but the org-wide
deny still applies.

## Precedence

When organization governance is active, local rules are not evaluated. Only
organization rules set in the Admin Console determine what is allowed or denied,
and they can't be supplemented or overridden from a developer's machine. The
same applies to filesystem policies: organization rules replace local behavior
entirely. For how a user's organization policies are evaluated together, see
[Policy concepts](concepts.md#rule-evaluation).

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
