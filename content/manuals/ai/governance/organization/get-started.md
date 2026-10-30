---
title: Get started with organization governance
linkTitle: Get started
weight: 10
description: Set up organization governance for Docker Sandboxes, from planning policies and configuring audit logs to assigning licenses, enforcing sign-in, deploying Docker Sandboxes, and verification.
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
  hosts.
- Host paths that sandboxes mount as workspaces. By default, `sbx run` mounts
  the directory it's run from, so allow the directories where developers keep
  their projects. A sandbox can't start if its workspace isn't allowed.
- MCP servers that developers add with `sbx mcp add`. A server that an agent
  connects to directly, without `sbx`, is a network destination, so plan for it
  with your network rules instead.

Before you assign licenses, you can run `sbx policy log` on a developer's
machine to see the network requests its sandboxes make. Once developers are
governed, your audit logs record their network requests, mounts, and MCP
activity in more detail. If you're not sure what sandboxes need, start the
pilot team with the rules you're confident about, then review your audit logs
for requests your policies deny and add allow rules where developers need them.
Audit records are available after you configure audit delivery and assign
licenses in steps 3 and 4.

Then decide how to scope your policies. Use org-wide policies for guardrails
that must apply everywhere, such as deny rules that block a destination even
when a team-scoped policy allows it, and team-scoped policies for access that
only some teams need. For an example of how the two combine, see
[How rules decide a request](policies/_index.md#how-rules-decide-a-request).
Plan to start with one pilot team before you roll out to everyone.

## Step 2: Create your policies

In [Docker Home](https://app.docker.com), expand **AI Platform** and create
your policies under **Network access**, **Filesystem access**, and **MCP
access**. Create at least a network policy and a filesystem policy before you
assign any licenses, so governed developers can keep working. Scope your first
policies to the pilot team.

For the steps, see
[Create and manage organization policies](policies/create-manage.md).

Consider also
[configuring a support message](policies/create-manage.md#configure-a-support-message).
Developers see it when your policy blocks a request, so you can tell them how
to ask for access.

## Step 3: Configure audit delivery

Audit logs record the policy decisions made for your developers' sandboxes.
Configure delivery before you assign licenses, so records are collected from
the first governed request. In Docker Home, go to **AI Platform** >
**Audit logs** to choose where records are delivered. For more details, see
[Configure audit delivery](audit/configure.md).

## Step 4: Assign licenses

Organization policies apply only to members who hold an AI Governance license.
Assign licenses to the members of your pilot team first. For more details, see
[Manage licenses](/manuals/accounts/organization/manage/manage-licenses.md).

## Step 5: Enforce sign-in

Without sign-in enforcement, a developer can sign in to Docker Sandboxes with a
personal account and avoid your policies. Deploy an enforcement configuration
to your developers' machines with your endpoint management tooling, listing
the organizations that developers must belong to. For more details, see
[Enforce sign-in](sign-in-enforcement.md).

## Step 6: Deploy Docker Sandboxes

If your developers don't have Docker Sandboxes yet, deploy it after you assign
licenses and enforce sign-in, so developers are governed from their first
sign-in. For installation options, see
[Install Docker Sandboxes](/manuals/ai/sandboxes/install.md).

## Step 7: Verify on a developer machine

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

## Step 8: Roll out to more teams

When the pilot team works as expected, assign licenses to more members, deploy
Docker Sandboxes to their machines, and widen the scope of your policies, or
add team-scoped policies for teams that need different access. As more
developers are governed, keep reviewing your audit logs for requests your
policies block, and add allow rules where developers need them. Each change you make reaches developer machines within 5 minutes.
For more details, see
[When policy changes reach developers](policies/_index.md#when-policy-changes-reach-developers).
