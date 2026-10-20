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

[Local policies](../local/_index.md) give individual developers control over what their
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
