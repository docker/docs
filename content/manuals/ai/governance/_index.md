---
title: AI Governance
weight: 15
description: Control what sandboxes can access, from local developer rules to org-wide enforcement.
keywords: docker sandboxes, governance, policy, network access, filesystem access, mcp policy, organization policy
aliases:
  - /ai/sandboxes/governance/
params:
  sidebar:
    group: AI and agents
---

Manage organization policies, sign-in enforcement, and audit delivery for
local Docker Sandboxes on developer machines. Organization governance requires
a separate paid subscription. [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
to request access.

Cloud sandbox account policies are personal controls, separate from this
organization governance. Delivering local audit records to the cloud does not
extend governance to cloud sandboxes.

For developer tasks, use [local CLI policies](/manuals/ai/sandboxes/cli/network-local.md)
or [cloud CLI policies](/manuals/ai/sandboxes/cli/network-cloud.md). See
[Network access and policies](/manuals/ai/sandboxes/concepts/network.md) for the
relationship between personal controls and organization enforcement.

## Learn more

Start with [Policy concepts](/manuals/ai/governance/concepts.md) for organization and team scope,
rule syntax, and precedence.

### Access controls

[Set up organization policies](/manuals/ai/governance/access-controls/organization.md), then configure
[network](/manuals/ai/governance/access-controls/network.md), [filesystem](/manuals/ai/governance/access-controls/filesystem.md),
and [MCP](/manuals/ai/governance/access-controls/mcp.md) access. Use team assignments where policies
need different scopes within the organization.

### Monitor and enforce

[Enforce organization sign-in](/manuals/ai/governance/monitor-and-enforce/sign-in-enforcement.md)
through endpoint management. [Inspect enforcement](/manuals/ai/governance/monitor-and-enforce/monitoring.md)
to diagnose policy decisions on developer machines.

[Audit logs](/manuals/ai/governance/audit/_index.md) cover local collection, hosted event delivery,
retention, export, and SIEM integration.

### Reference

Use the [Governance API and MCP policy references](/manuals/ai/governance/reference/_index.md) for
programmatic administration and exact policy fields, and the
[audit record reference](/manuals/ai/governance/audit/record-reference.md) for event fields.
