---
title: Create and manage settings policies
linkTitle: Create and manage
weight: 10
description: Create, assign, edit, delete, and set the precedence of settings policies for Docker Sandboxes in Docker Home.
keywords: docker sandboxes, settings policy, create policy, policy precedence, Docker Home, sbx settings, AI governance, locked settings
---

{{< summary-bar feature_name="AI Governance" >}}

Create and manage settings policies in Docker Home. For how setting states,
scope, and precedence work, see
[Sandboxes settings management](_index.md).

## Before you begin

- Your organization needs an
  [AI Governance plan](/manuals/subscription-billing/plans/ai-governance.md).
- You need the organization
  [owner](/manuals/security/roles-and-permissions/core-roles.md) role, or a
  [custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
  with the **Governance** permissions.
- Settings policies apply to developers who are signed in to Docker Sandboxes
  as members of your organization. To stop developers from avoiding your
  policies by signing in with a personal account, set up
  [sign-in enforcement](../sign-in-enforcement.md).

## Create a settings policy

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select
   **Sandboxes Settings**.
1. Select **Create policy**.
1. Enter a **Settings policy name** of at least four characters, and an
   optional **Description**.
1. In **Assign this policy**, select **All organization members**, **Teams**,
   or **Members**. For teams or members, select the ones the policy applies
   to.
1. In **Configure the settings for this policy**, choose a state for each
   setting you want to manage. Leave the others as **Not set**. For the options
   each setting type offers, see
   [Setting states by type](../../reference/settings-policy.md#setting-states-by-type).
1. Select **Create policy**.

For the number of policies you can create, see
[Organization policy limits](../../reference/policy-limits.md).

## Edit or delete a settings policy

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select
   **Sandboxes Settings**.
1. Select the action menu (⋮) for the policy, then select **Edit**, or delete
   the policy.
1. If you're editing, make your changes and select **Save**.

## Set the precedence of team policies

When a developer is on more than one team with a settings policy, and those
policies set the same setting, the order of your team policies decides which
value applies for that setting. Organization and member policies don't have
an order.

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. In the left-hand navigation, expand **AI Platform** and select
   **Sandboxes Settings**.
1. Select **Set precedence**.
1. Drag the team policies into order. The policy at the top takes precedence.
1. Select **Save**.
