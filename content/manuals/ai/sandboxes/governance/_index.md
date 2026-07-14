---
title: Governance
weight: 55
description: Control what sandboxes can access, from local developer rules to org-wide enforcement.
keywords: docker sandboxes, governance, policy, network access, filesystem access, mcp policy, organization policy
---

Sandbox governance covers the policy system that controls what sandboxes can
access over the network, on the filesystem, and through MCP. It operates at two
layers, and only one applies at a time:

**Local policy** is configured per machine using the `sbx policy` CLI. It
lets individual developers customize which domains their sandboxes can reach.
See [Local policy](access-controls/local.md).

**Organization policy** is configured centrally in the Docker Admin Console.
Network and filesystem policies can also be managed via the
[Governance API](/reference/api/ai-governance/). Rules defined at the org level
apply uniformly across every sandbox in the organization. Organization
governance can also include MCP policies for sandbox MCP activity. When
organization governance is active, it replaces local policy entirely: local
`sbx policy` rules are no longer evaluated. See
[Organization policies](access-controls/organization.md).

Alongside this access-control policy, admins can require developers to sign in
as members of their organization before using sandboxes at all.
[Sign-in enforcement](monitor-and-enforce/sign-in-enforcement.md) is deployed
through endpoint management and ensures developers can't bypass organization
policy by using a personal account.

> [!NOTE]
> Organization governance is available on a separate paid subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Learn more

Start with [Policy concepts](concepts.md) for the resource model, rule syntax,
MCP policy basics, evaluation, and precedence.

### Access controls

- [Local policy](access-controls/local.md): configure network rules on your
  machine with the `sbx policy` CLI.
- [Organization policies](access-controls/organization.md): centrally manage
  sandbox policies across your organization from the Admin Console.
- [Network access rules](access-controls/network.md): control outbound network
  access from sandboxes.
- [Filesystem access rules](access-controls/filesystem.md): control which host
  paths sandboxes can mount as workspaces.
- [MCP access policies](access-controls/mcp.md): control MCP server registration,
  tool calls, resources, prompts, and approval gates.

### Monitor and enforce

- [Monitoring policies](monitor-and-enforce/monitoring.md): inspect active
  rules and monitor sandbox network traffic with `sbx policy ls` and
  `sbx policy log`.
- [Audit logs](monitor-and-enforce/audit.md): capture a durable, structured
  record of policy decisions for SIEM ingestion and compliance.
- [Sign-in enforcement](monitor-and-enforce/sign-in-enforcement.md): require
  developers to sign in as organization members, enforced through endpoint
  management.

### Reference

- [AI Governance API](/reference/api/ai-governance/): manage network and
  filesystem org policies programmatically.
- [MCP policy reference](reference/mcp-policy.md): look up Docker MCP policy
  actions, resources, attributes, context fields, and approval behavior.
