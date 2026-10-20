---
title: Create and manage policies
linkTitle: Create and manage
weight: 10
description: Create, scope, edit, and delete organization network, filesystem, and MCP policies in Docker Home.
keywords: docker sandboxes, organization policy, create policy, Docker Home, team scope, policy limits
---

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
1. Define the policy rules.
   - Network and filesystem policies: select **Add rule** for each rule. For a
     network policy, see [Network rule fields](network.md#network-rule-fields).
     For a filesystem policy, see
     [Filesystem mount rule syntax](../../reference/filesystem-rules.md).
   - MCP policies: enter Cedar statements in the policy editor. See
     [MCP access policies](mcp.md).
1. For a network policy, set **Require approval before access** if developers
   should confirm each destination before a sandbox can reach it. See
   [Require approval for a network policy](network.md#require-approval-for-a-network-policy).

Existing policies are listed with their name, scope, rule count, and last
update. Use the action menu (⋮) to edit or delete a policy.

If you manage policies through the
[Governance API](/reference/api/ai-governance/), a network or filesystem
policy's `domain` field is `network` or `filesystem`. It's derived from the
policy's rule actions, so one policy can't mix network and filesystem rules.

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

## Limits

Organization policies have the following limits, which help ensure fair usage
and resource availability across organizations:

| Limit                     | Value                                               |
|---------------------------|-----------------------------------------------------|
| Policies per organization | 100                                                 |
| Rules per policy          | 250                                                 |
| Policy size               | 400 KB total, shared across all of a policy's rules |

Typical policies use only a small fraction of the policy size limit. Domain
and file path values have no separate length limit beyond valid format.

If these limits don't fit your organization's needs,
[contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
to discuss options.
