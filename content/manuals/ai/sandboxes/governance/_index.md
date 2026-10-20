---
title: Governance
weight: 120
description: Control what sandboxes can access, from local developer rules to org-wide enforcement.
keywords: docker sandboxes, governance, policy, network access, filesystem access, mcp policy, organization policy
aliases:
  - /ai/sandboxes/governance/access-controls/
  - /ai/sandboxes/governance/monitor-and-enforce/
---

Sandbox governance controls what sandboxes can access over the network, on the
filesystem, and through MCP. For MCP setup and server registration, see
[MCP gateway](../mcp-gateway.md). Governance has two parts:

- **Local policy** is the policy you set on your machine with the `sbx policy`
  CLI. It lets individual developers customize which domains their sandboxes
  can reach. See [Local policy](local/_index.md).
- **Organization governance** is policy, sign-in enforcement, and audit logs
  that administrators set for organization members, either organization-wide
  or for selected teams. Organization policy is configured centrally in Docker
  Home. Network, filesystem, and MCP policies can also be managed through the
  [Governance API](/reference/api/ai-governance/). See
  [Organization governance](organization/_index.md).

Both apply to local sandboxes. Cloud sandboxes use separate network policies.
See [Cloud network policy](../cloud/network-policy.md) and Docker Agentic
Platform [network policies](/manuals/agentic-platform/policies.md).

Organization policies apply to members who hold an AI Governance license.
Developer machines must run sbx 0.35.0 or later to enforce HTTP method and
path rules.

When an organization policy is enforced for you, access is default deny: a
request that no organization allow rule matches is blocked, and a matching deny
rule takes precedence over any allow rule. Governed MCP activity is default
deny in the same way. Without organization governance, a request that no rule
matches is handled according to your local preset. See
[Local policy](local/_index.md).

Administrators can also require developers to sign in as members of their
organization before using sandboxes at all.
[Sign-in enforcement](organization/sign-in-enforcement.md) is deployed
through endpoint management and ensures developers can't bypass organization
policy by using a personal account.

> [!NOTE]
> Organization governance requires a Docker AI Governance subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Resource model

Docker sandbox governance is built around two resource types: **policies** and
**rules**.

A **policy** is a named collection of rules that controls sandbox access.
Policies exist at two levels:

- **Local**: configured per machine using the `sbx policy` CLI. Applies to
  sandboxes on that machine only.
- **Organization**: configured in Docker Home or through the
  [Governance API](/reference/api/ai-governance/). Applies to licensed members
  of the organization. An organization can have several policies, each applying
  either org-wide or to specific teams. See
  [Policy scope](organization/policies/_index.md#policy-scope).

When an organization policy is enforced for you, only organization allow rules
can grant access. Local and kit-defined deny rules still apply on top. See
[Precedence](organization/policies/_index.md#precedence).

A **rule** is the unit of access control within a policy. Each rule has:

- **Name**: a human-readable label
- **Actions**: the type of access the rule controls
- **Resources**: the targets the rule matches against
- **Decision**: `allow` or `deny`

Each policy covers one kind of access: network or filesystem. MCP policies use
Cedar statements written in the `MCP` namespace instead of the network and
filesystem rule format.
