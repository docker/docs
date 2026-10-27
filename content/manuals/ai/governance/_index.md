---
title: AI governance
weight: 20
description: Control what sandboxes can access, from local developer rules to org-wide enforcement.
keywords: docker sandboxes, governance, policy, network access, filesystem access, mcp policy, organization policy
params:
  sidebar:
    group: AI and agents
aliases:
  - /ai/sandboxes/governance/
  - /ai/sandboxes/governance/access-controls/
  - /ai/sandboxes/governance/monitor-and-enforce/
  - /ai/sandboxes/governance/reference/
---

Agents in [Docker Sandboxes](../sandboxes/_index.md) work on their own. They
install packages, call APIs, and use tools without asking you first.
Governance lets an agent work freely within rules that limit what it can
access. The rules are enforced outside the sandbox, so the agent can't change or
bypass them. They come from three sources: your local policy, the kits your
sandboxes use, and your organization.

> [!NOTE]
> This section covers sandboxes that run on your machine. Cloud sandboxes use
> the same kinds of allow and deny rules, but they don't use local policy or
> organization policy. They use
> [cloud network policy](../sandboxes/cloud/network-policy.md), which is stored in your
> Docker account and managed with `sbx --cloud policy` or in the Docker Agentic
> Platform [Console](/manuals/agentic-platform/policies.md).

## Where the rules come from

Rules come either from your own machine or from your organization.

### Your machine, with local policy and kits

[Local policy](local/_index.md) is the policy you set on your machine with the
`sbx policy` CLI. You start from a preset, then allow the network destinations
your sandboxes need and block the ones they shouldn't reach.

The [kits](../sandboxes/customize/use-kits.md) your sandboxes use, including
the built-in agent kits, can also add network rules, such as the API endpoints
an agent needs. You choose the kit, and the kit brings its rules. They work
like the rules you add yourself.

### Your organization, with organization governance

With an
[AI Governance plan](/manuals/subscription-billing/plans/ai-governance.md),
administrators use [organization governance](organization/_index.md) to set
policies for the licensed members of their organization and review the
decisions those policies make in audit logs. Organization policies control network
access, the host paths sandboxes can mount, and MCP activity.

## Get started

If you're a developer, start with
[Try out local policy](local/try-it-out.md), an example that allows and blocks
destinations for one sandbox.

If you're an administrator, start with
[Get started with organization governance](organization/get-started.md) to set
up policies, audit logs, licenses, sign-in enforcement, and Docker Sandboxes
for your organization.
