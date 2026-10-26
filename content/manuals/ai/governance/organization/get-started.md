---
title: Get started with organization governance
linkTitle: Get started
weight: 10
description: Set up organization governance for Docker Sandboxes, from planning policies and assigning licenses to sign-in enforcement, audit logs, and verification.
keywords: docker sandboxes, AI Governance, organization governance, get started, organization policy, sign-in enforcement, audit logs, rollout
---

{{< summary-bar feature_name="AI Governance" >}}

This page walks you through setting up organization governance, from planning
your first policies to confirming they work on a developer machine. Each step
explains what to do and what to decide, and links to the page with the full
instructions.

## Before you begin

You need:

- A Docker
  [AI Governance plan](/manuals/subscription-billing/plans/ai-governance.md)
  for your organization.
- Your developers added as
  [members of your organization](/manuals/accounts/organization/manage/members.md).
  To scope policies to groups of developers, such as a pilot team, also
  [create teams](/manuals/accounts/organization/manage/manage-a-team.md) for
  them.
- The organization
  [owner](/manuals/security/roles-and-permissions/core-roles.md) role, or a
  [custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
  with the **Governance** permissions.
- Access to the endpoint management tooling for your developers' machines, such
  as an MDM solution or Group Policy, to enforce sign-in.

## Step 1: Plan what to allow

Organization policy is default deny. Once an organization policy is enforced
for a developer, their sandboxes can reach only the network destinations your
policies allow, mount only the host paths your policies allow, and use MCP
servers only as your policies permit. For more details, see
[Organization policies](policies/_index.md).

Before you govern anyone, work out what their sandboxes need:

- Network destinations, such as AI provider APIs, package registries, and code
  hosts. On a developer machine, `sbx policy log` shows the hosts that sandboxes
  contact today.
- Host paths that sandboxes mount as workspaces, such as the directory where
  developers keep their projects. A sandbox can't start without a workspace it's
  allowed to mount.
- MCP servers, if your developers use them through Docker's MCP gateway.

Then decide how to scope your policies. Use org-wide policies for guardrails
that must apply everywhere, such as deny rules for destinations no one should
reach, and team-scoped policies for access that only some teams need. Plan to
start with one pilot team before you roll out to everyone.

## Step 2: Create your policies

In [Docker Home](https://app.docker.com), expand **AI Platform** and create
your policies under **Network access**, **Filesystem access**, and **MCP
access**. Create at least a network policy and a filesystem policy before you
assign any licenses, so governed developers can keep working. Scope your first
policies to the pilot team.

For the steps, see
[Create and manage organization policies](policies/create-manage.md). For the
rules each type of policy supports, see
[Organization network access policies](policies/network.md),
[Organization filesystem access policies](policies/filesystem.md), and
[Organization MCP access policies](policies/mcp.md).

Consider also
[configuring a support message](policies/create-manage.md#configure-a-support-message).
Developers see it when your policy blocks a request, so you can tell them how
to ask for access.

## Step 3: Assign licenses

Organization policies apply only to members who hold an AI Governance license.
Assign licenses to the members of your pilot team first. For more details, see
[Manage licenses](/manuals/accounts/organization/manage/manage-licenses.md).

## Step 4: Enforce sign-in

Without sign-in enforcement, a developer can sign in to Docker Sandboxes with a
personal account and avoid your policies. Deploy an enforcement configuration
to your developers' machines with your endpoint management tooling, listing
the organizations that developers must belong to. For more details, see
[Enforce sign-in](sign-in-enforcement.md).

## Step 5: Configure audit delivery

Audit logs record the policy decisions made for your developers' sandboxes. In
Docker Home, go to **AI Platform** > **Audit logs** to choose where records are
delivered. For more details, see
[Configure audit delivery](audit/configure.md).

## Step 6: Verify on a developer machine

On the machine of a developer in your pilot team, pick up the latest
organization policy right away instead of waiting for the next sync:

```console
$ sbx policy refresh
```

Confirm that organization policy is enforced. The status line should show your
organization:

```console
$ sbx policy ls
```

Check a destination your policy allows and one it doesn't:

```console
$ sbx policy check network api.anthropic.com
$ sbx policy check network blocked.example.com
```

The first check should report `Allowed` and the second `Denied`. Then run a
sandbox and use `sbx policy log` to see which requests it made and which
rules matched. For more details, see
[Inspect policy status and rules](../local/inspect.md) and
[Test and monitor network policy](../local/test-monitor.md).

## Step 7: Roll out to more teams

When the pilot team works as expected, assign licenses to more members and widen
the scope of your policies, or add team-scoped policies for teams that need
different access. Use your audit logs to find requests your policies block, and
add allow rules where developers need them. Policy changes reach developer
machines within 5 minutes. For more details, see
[When policy changes reach developers](policies/_index.md#when-policy-changes-reach-developers).
