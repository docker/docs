---
title: Governance
weight: 120
description: Control what sandboxes can access, from local developer rules to org-wide enforcement.
keywords: docker sandboxes, governance, policy, network access, filesystem access, mcp policy, organization policy
aliases:
  - /ai/sandboxes/governance/access-controls/
  - /ai/sandboxes/governance/monitor-and-enforce/
---

The governance described here applies to local sandboxes. Cloud sandboxes
use separate network policy configuration. See
[Cloud network policy](../cloud/network-policy.md) for cloud controls.

Sandbox governance covers the policy system that controls what sandboxes can
access over the network, on the filesystem, and through MCP. For MCP setup and
server registration, see [MCP gateway](../mcp-gateway.md). Governance operates
at two layers:

**Local policy** is configured per machine using the `sbx policy` CLI. It
lets individual developers customize which domains their sandboxes can reach.
See [Local policy](local/_index.md).

**Organization policy** is configured centrally in Docker Home. Network and
filesystem policies can also be managed via the
[Governance API](/reference/api/ai-governance/). Controls defined at the org
level apply uniformly across every local sandbox in the organization. Organization
governance can also include MCP policies for sandbox MCP activity. When
organization governance is active, only organization allow rules grant access:
local `sbx policy` allow rules are no longer evaluated, while local deny rules
still apply on top. See
[Organization policies](organization/_index.md).

Alongside this access-control policy, admins can require developers to sign in
as members of their organization before using sandboxes at all.
[Sign-in enforcement](organization/sign-in-enforcement.md) is deployed
through endpoint management and ensures developers can't bypass organization
policy by using a personal account.

> [!NOTE]
> Organization governance is available on a separate paid subscription.
> [Contact Docker Sales](https://www.docker.com/products/ai-governance/#contact-sales)
> to request access.

## Resource model

Docker sandbox governance is built around two resource types: **policies** and
**rules**.

A **policy** is a named collection of rules that controls sandbox access.
Policies exist at two levels:

- **Local**: configured per machine using the `sbx policy` CLI. Applies to
  sandboxes on that machine only.
- **Organization**: configured in Docker Home. Network and filesystem policies
  can also be managed via the
  [Governance API](/reference/api/ai-governance/). Applies to sandboxes across
  the organization. An organization can have several policies, each applying
  either org-wide or to specific teams. See [Policy scope](organization/policies/_index.md#policy-scope).

When organization governance is active, only organization allow rules can grant
access. Local and kit-defined deny rules still apply on top. See
[Precedence](organization/policies/_index.md#precedence).

A **rule** is the unit of access control within a policy. Each rule has:

- **Name**: a human-readable label
- **Actions**: the type of access the rule controls
- **Resources**: the targets the rule matches against
- **Decision**: `allow` or `deny`

Rules are grouped by domain. Network and filesystem rules in a policy must
share the same domain, either `network` or `filesystem`. MCP policies use Cedar
statements written in the `MCP` namespace instead of the network and filesystem
rule format.

An organization network policy can also require approval, which turns every
allow in that policy into a request the developer must confirm before access is
granted. Approval is set on the policy rather than on individual rules, so it
applies to all of the policy's allow rules at once. Without organization
governance, a request with no matching allow or deny rule also asks for
approval. See
[Approval-required access](organization/policies/network.md#approval-required-access).
