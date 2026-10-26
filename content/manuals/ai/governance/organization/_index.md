---
title: Organization governance
weight: 20
description: Centrally manage sandbox network, filesystem, and MCP policies and review audit logs for your organization.
keywords: docker sandboxes, governance, organization policy, AI governance, Docker Home, network access, filesystem access, mcp policy, sign-in enforcement, audit logs
aliases:
  - /ai/sandboxes/security/governance/
  - /ai/sandboxes/governance/org/
  - /ai/sandboxes/governance/access-controls/organization/
---

{{< summary-bar feature_name="AI Governance" >}}

Agents in sandboxes act on their own, on your developers' machines. With
[local policy](../local/_index.md), each developer decides what their
sandboxes can reach. Organization governance moves that decision to you, the
administrator. You set the rules once in Docker Home, and they apply to the
sandboxes on every governed developer's machine. Organization policy is default
deny, so those sandboxes can reach, mount, and use only what your policies
allow. You also get a record of what those sandboxes did.

## How organization governance works

Organization governance gives you central control, in three parts: you decide
who is governed, you set what their sandboxes can do, and you review what their
sandboxes did.

### Decide who is governed

Organization governance applies to the members of your organization who hold
an AI Governance license. You can apply each policy to the whole organization
or only to selected teams.

A developer could avoid your rules by signing in to Docker Sandboxes with a
personal account. [Sign-in enforcement](sign-in-enforcement.md) closes that
gap by allowing only members of your organization to use sandboxes on your
managed machines.

### Set what sandboxes can do

Organization policies control what sandboxes can do. You can create three types
of policy:

- [Organization network access policies](policies/network.md): the
  destinations sandboxes can reach, down to HTTP methods and paths.
- [Organization filesystem access policies](policies/filesystem.md): the host
  paths developers can mount into sandboxes.
- [Organization MCP access policies](policies/mcp.md): the MCP servers
  developers can register and what agents can do with them.

When an organization policy is enforced for a developer, it changes which of
their local rules apply. For more details, see
[How organization policy relates to local and kit rules](policies/_index.md#how-organization-policy-relates-to-local-and-kit-rules).

You manage policies in Docker Home, or programmatically through the
[AI Governance API](/reference/api/ai-governance/).

### Review what sandboxes did

[Audit logs](audit/_index.md) record governance activity for your
organization, such as the policy decision made for each sandbox request. You
can view and export them in Docker Home, or forward them to your SIEM.

## What developers see

Developers keep using the same `sbx` commands. When an organization policy is
enforced for them, `sbx policy ls` shows `Governance: Managed by <org>` in its
status line. A request that your policy blocks returns `Blocked by org policy`,
followed by your
[support message](policies/create-manage.md#configure-a-support-message) if
you configured one. For more details, see
[Inspect policy status and rules](../local/inspect.md).

## Who can manage governance

By default, only organization
[owners](/manuals/security/roles-and-permissions/core-roles.md) can view and
manage AI Governance policies in Docker Home. To let someone else manage
policies, create a
[custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
with the **Governance** permissions and assign it to a user or team.

## Get started

To set up organization governance in this order, from licenses and sign-in
enforcement to policies and audit logs, see
[Get started with organization governance](get-started.md). To request an AI
Governance subscription,
[contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales).
