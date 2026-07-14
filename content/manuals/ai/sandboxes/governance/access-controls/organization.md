---
title: Organization policies
linkTitle: Org policies
weight: 20
description: Centrally manage sandbox network, filesystem, and MCP policies for your organization.
keywords: docker sandboxes, governance, organization policy, AI governance, admin console, network access, filesystem access, mcp policy
aliases:
  - /ai/sandboxes/security/governance/
  - /ai/sandboxes/governance/org/
---

[Local policies](local.md) give individual developers control over what their
sandboxes can access. Organization policy moves that control to the admin level:
rules defined in **Admin Console** apply to sandboxes across the organization,
either to every member or to specific teams. When organization governance is
active, it replaces local `sbx policy` rules entirely — local rules are no
longer evaluated and can't be used to supplement or override the organization
policy.

Admins can manage organization policies through the Admin Console UI. For
programmatic management of network and filesystem policies, use the
[Governance API](/reference/api/ai-governance/).

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

The steps in this section cover network and filesystem policies. MCP policies
use Cedar rather than rule rows. For MCP examples, see
[MCP tool access](mcp.md).

To create a policy:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Select **AI Platform**, then the governance section you want.
1. Select **Network access** or **Filesystem access**, then **Create policy**.
1. Enter a **Policy name**.
1. Set the **Scope** to **Organization** or **Teams**. If you select **Teams**,
   choose the teams the policy applies to. See
   [Scope policies to teams](#scope-policies-to-teams).
1. Select **Add rule** to add each rule. For rule syntax, use the relevant
   access-control page in [Choose a policy type](#choose-a-policy-type).

Existing policies are listed with their name, scope, rule count, and last
update. Use the action menu (⋮) to edit or delete a policy.

## Choose a policy type

Organization policies are managed by access surface. Use the access-control
pages for syntax, examples, and enforcement details:

- [Network access rules](network.md): control outbound network access from
  sandboxes.
- [Filesystem access rules](filesystem.md): control which host paths sandboxes
  can mount as workspaces.
- [MCP tool access](mcp.md): control MCP server registration, tool calls,
  resources, prompts, and approval gates with Cedar policy.

When organization governance is active, local and kit-defined rules are not
evaluated. To see which rules are active on a developer machine, use
[Monitoring policies](../monitor-and-enforce/monitoring.md).

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

### How scoped policies combine

A user is governed by all of their
[effective policies](../concepts.md#policy-scope): every org-wide policy, plus
the team-scoped policies for the teams they belong to. Use org-wide policies
for guardrails that must apply everywhere, and team-scoped policies for access
that only some teams need.

For precedence between local and organization policies, and for how allow and
deny rules combine, see [Policy concepts](../concepts.md).

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

#### Network versus filesystem enforcement timing

Network policy and filesystem policy differ in when a change takes effect:

- Network policy is evaluated on every outbound request. Once a policy
  change has synced to the developer's machine (up to 5 minutes), it applies
  immediately to subsequent requests.

- Filesystem policy is only checked when a workspace is mounted — that
  is, when a sandbox is created. Once a sandbox is running, changing the
  filesystem policy has no effect on that sandbox. The sandbox continues to
  access the previously allowed path until it is removed and a new one is
  created.

To apply a filesystem policy change immediately, remove the running sandbox
and create a new one.
