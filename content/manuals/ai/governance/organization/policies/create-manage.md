---
title: Create and manage organization policies
linkTitle: Create and manage
weight: 40
description: Create, scope, edit, and delete organization network, filesystem, and MCP policies in Docker Home.
keywords: docker sandboxes, organization policy, create policy, Docker Home, team scope, policy limits
---

{{< summary-bar feature_name="AI Governance" >}}

Create and manage organization policies in Docker Home. The steps are the same
for network, filesystem, and MCP policies. Only the rules you add differ. For
how scope, rule evaluation, and precedence work, see
[Organization policies](_index.md).

To manage policies programmatically instead, use the
[AI Governance API](/reference/api/ai-governance/).

## Before you begin

You need the organization
[owner](/manuals/security/roles-and-permissions/core-roles.md) role, or a
[custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
with the **Governance** permissions.

## Create a policy

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select
   **Network access**, **Filesystem access**, or **MCP access**.
1. Select **Create policy**.
1. Enter a **Policy name**.
1. Set the **Scope** to **Organization** or **Teams**. If you select **Teams**,
   choose the teams the policy applies to. For more details, see
   [Scope policies to teams](#scope-policies-to-teams).
1. Define the policy rules.
   - Network and filesystem policies: select **Add rule** for each rule. For a
     network policy, see
     [What a network rule can match](network.md#what-a-network-rule-can-match).
     For a filesystem policy, see
     [Filesystem mount rule syntax](../../reference/filesystem-rules.md).
   - MCP policies: enter Cedar statements in the policy editor. For more
     details, see [Organization MCP access policies](mcp.md).
1. For a network policy, set **Require approval before access** if developers
   should confirm each destination before a sandbox can reach it. The policy's
   detail page reports approval as **Required** or **Not required**. For more
   details, see [Approval-required access](network.md#approval-required-access).

For the number of policies and rules you can create, see
[Organization policy limits](../../reference/policy-limits.md).

## Scope policies to teams

A policy's **Scope** controls who it applies to. Set it to **Organization** to
apply the policy to every licensed member, or to **Teams** to apply it only to
licensed members of the teams you select. For more details on how org-wide and
team-scoped policies combine, see
[Policy scope](_index.md#policy-scope).

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

## Edit or delete a policy

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select
   **Network access**, **Filesystem access**, or **MCP access**.
1. Find the policy in the list. Each policy is listed with its name, scope,
   rule count, and last update.
1. Select the action menu (⋮) for the policy, then select **Edit**, or delete
   the policy.
1. If you're editing, make your changes and select **Save**.

Editing a policy replaces it in full. For example, turning off **Require
approval before access** on a network policy removes the requirement from every
rule in that policy. Changes reach developer machines within 5 minutes. For more
details, see
[When policy changes reach developers](_index.md#when-policy-changes-reach-developers).

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
policy and for requests an
[approval-required policy](network.md#approval-required-access) blocks. If you
leave it blank, Docker shows the policy denial without additional contact text.
