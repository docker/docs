---
title: Organization governance
weight: 20
description: Centrally manage sandbox network, filesystem, and MCP policies for your organization.
keywords: docker sandboxes, governance, organization policy, AI governance, Docker Home, network access, filesystem access, mcp policy
aliases:
  - /ai/sandboxes/security/governance/
  - /ai/sandboxes/governance/org/
  - /ai/sandboxes/governance/access-controls/organization/
---

Organization governance is policy, sign-in enforcement, and audit logs that
administrators set for organization members, either organization-wide or for
selected teams. Organization policies apply to the local sandboxes of members
who hold an AI Governance license. When an organization policy is enforced for
a member, it takes precedence over their [local policy](../local/_index.md).
See [Precedence](policies/_index.md#precedence).

Admins can manage organization policies through the Docker Home UI. For
programmatic management of network, filesystem, and MCP policies, use the
[Governance API](/reference/api/ai-governance/).

By default, only organization
[owners](/manuals/security/roles-and-permissions/core-roles.md) can
view and manage AI Governance policies. To let someone other than an owner
manage policies, create a
[custom role](/manuals/security/roles-and-permissions/custom-roles/_index.md)
with the **Governance** permissions and assign it to a user or team.

> [!NOTE]
> Organization governance requires a Docker AI Governance subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.
