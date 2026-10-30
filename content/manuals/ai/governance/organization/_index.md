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

Organization governance lets you, the administrator, control what sandboxes can
do across your organization. You set the rules once in Docker Home, and they
apply to the sandboxes on every governed developer's machine, instead of each
developer deciding with local policy. Organization policy is default deny, so
those sandboxes can reach, mount, and use only what your policies allow.

## How organization governance works

Organization governance has three parts. Licenses and sign-in enforcement
decide who is governed, organization policies control what their sandboxes can
do, and audit logs record the decisions those policies make.

### Licenses and sign-in enforcement

Organization governance applies to the members of your organization who hold
an AI Governance license. You can apply each policy to the whole organization
or only to selected teams.

A developer could avoid your rules by signing in to Docker Sandboxes with a
personal account. [Sign-in enforcement](sign-in-enforcement.md) closes that
gap by allowing only members of your organization to use sandboxes on your
managed machines.

### Organization policies

Organization policies control which destinations sandboxes can reach, which
host paths developers can mount, and which MCP servers developers can register
and use. You create them in Docker Home, or programmatically through the
[AI Governance API](/reference/api/ai-governance/). For more details, see
[Organization policies](policies/_index.md).

### Audit logs

[Audit logs](audit/_index.md) record governance activity for your
organization, such as the policy decision made for each sandbox request. You
can view and export them in Docker Home, or forward them to your SIEM.

## Developer experience

Developers keep using the same `sbx` commands. When an organization policy is
enforced for them, `sbx policy ls` shows `Governance: Managed by <org>` in its
status line, and only your allow rules grant access, so their own local and kit
allow rules stop applying. For more details, see
[How organization policy relates to local and kit rules](policies/_index.md#how-organization-policy-relates-to-local-and-kit-rules).
A request that your policy blocks returns `Blocked by org policy`, followed by
your [support message](policies/create-manage.md#configure-a-support-message)
if you configured one.

## Get started

To set up organization governance step by step, see
[Get started with organization governance](get-started.md). To get an AI
Governance plan,
[contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales).
