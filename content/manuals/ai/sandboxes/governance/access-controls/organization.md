---
title: Organization policies
linkTitle: Org policies
weight: 20
description: Centrally manage sandbox network, filesystem, and MCP policies for your organization.
keywords: docker sandboxes, governance, organization policy, AI governance, Docker Home, network access, filesystem access, mcp policy
aliases:
  - /ai/sandboxes/security/governance/
  - /ai/sandboxes/governance/org/
---

The governance described here applies to local sandboxes. Organization
governance is not available for cloud sandboxes in this release. See
[Cloud network policy](../../cloud/network-policy.md) for cloud controls.

[Local policies](local.md) give individual developers control over what their
sandboxes can access. Organization policy moves that control to the admin level:
organization policies apply to local sandboxes across the organization, either to
every member or to specific teams. When organization governance is active, only
organization allow rules grant access: local `sbx policy` allow rules are no
longer evaluated and can't expand what the organization permits. Local network
deny rules remain active, so developers can restrict access further but never
loosen it.

Admins can manage organization policies through the Docker Home UI. For
programmatic management of network and filesystem policies, use the
[Governance API](/reference/api/ai-governance/).

By default, only organization
[owners](/manuals/security/roles-and-permissions/core-roles.md) can
view and manage AI Governance policies. To let someone other than an owner
manage policies, create a
[custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
with the **Governance** permissions and assign it to a user or team.

> [!NOTE]
> Sandbox organization governance is available on a separate paid
> subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Create a policy

Manage policies from the **AI Platform** section in the left-hand navigation
of [Docker Home](https://app.docker.com).

To create a policy:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select
   **Network access**, **Filesystem access**, or **MCP access**.
1. Select **Create policy**.
1. Enter a **Policy name**.
1. Set the **Scope** to **Organization** or **Teams**. If you select **Teams**,
   choose the teams the policy applies to. See
   [Scope policies to teams](#scope-policies-to-teams).
1. Define the policy rules. For network and filesystem policies, select
   **Add rule** for each rule. For MCP policies, enter Cedar statements in the
   policy editor. For syntax and examples, use the relevant access-control page
   in [Choose a policy type](#choose-a-policy-type).
1. For a network policy, set **Require approval before access** if developers
   should confirm each destination before a sandbox can reach it. See
   [Require approval for a network policy](#require-approval-for-a-network-policy).

Existing policies are listed with their name, scope, rule count, and last
update. Use the action menu (⋮) to edit or delete a policy.

### Require approval for a network policy

Turning on **Require approval before access** means the destinations a network
policy allows aren't reachable until the developer confirms each one. For how
approval behaves and what satisfies it, see
[Approval-required access](network.md#approval-required-access).

To set it on an existing policy:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select
   **Network access**.
1. Select the policy, then choose **Edit**.
1. Turn on **Require approval before access**.
1. Select **Save changes**.

The policy's detail page reports approval as **Required** or **Not required**.
Editing a policy replaces it in full, so turning the setting off removes the
requirement from every rule in that policy.

## Configure a support message

Admins can add an optional support message that appears after the policy denial
details when a sandbox action is blocked by organization governance. Use it to
point members to an internal support channel, ticket queue, or security contact.

To set the message:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select **Manage**.
1. In **Support message**, enter up to 500 characters.
1. Select **Save changes**.

Docker shows the message only for denials caused by organization governance
policy. If you leave it blank, Docker shows the policy denial without additional
contact text.

## Choose a policy type

Organization policies are managed by access surface. Use the access-control
pages for syntax, examples, and enforcement details:

- [Network access policies](network.md): control outbound network access from
  sandboxes.
- [Filesystem access policies](filesystem.md): control which host paths
  sandboxes can mount as workspaces.
- [MCP access policies](mcp.md): control MCP server registration, tool calls,
  resources, prompts, and approval gates with Cedar policy.

When organization governance is active, local and kit-defined allow rules are
not evaluated, while deny rules from those sources still apply. See
[Precedence](../concepts.md#precedence). To see which rules are active on a
developer machine, use
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
[teams](/manuals/accounts/organization/manage/manage-a-team.md), so a team must
exist before you can scope a policy to it. Create teams and manage their members
in one of two ways:

- Manually, in Docker Home.
- Automatically, by using
  [group mapping](/manuals/security/provisioning/scim/group-mapping.md)
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
  immediately to subsequent requests.

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
