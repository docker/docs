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
Governance lets an agent work freely within limits you set. It has two parts:
local policy, which applies to one machine, and organization governance, which
applies across an organization and also records what sandboxes did in audit
logs.

> [!NOTE]
> This section covers sandboxes that run on your machine. Cloud sandboxes use
> the same kinds of allow and deny rules, but they don't use local policy or
> organization policy. They use
> [cloud network policy](../sandboxes/cloud/network-policy.md), which is stored in your
> Docker account and managed with `sbx --cloud policy` or in the Docker Agentic
> Platform [Console](/manuals/agentic-platform/policies.md).

## Who sets the rules

A rule allows or denies one kind of access, such as a network destination, and
a policy is a set of rules. The two parts of governance differ in who sets the
rules and who the rules apply to.

### You, with local policy and kits

[Local policy](local/_index.md) is the policy you set on your machine with the
`sbx policy` CLI. You start from a preset, then allow or block the network
destinations your sandboxes need. Local policy is available to everyone on their
machine.

The kits your sandboxes use, including the built-in agent kits, can add network
rules too, such as the API endpoints an agent needs. Kit rules work like the
rules you add yourself.

### Your organization, with organization governance

With [organization governance](organization/_index.md), administrators set
policies for the members of their organization who hold an AI Governance
license. Organization policies control network access, the host paths
sandboxes can mount, and MCP activity. Organization governance requires an
[AI Governance plan](/manuals/subscription-billing/plans/ai-governance.md).

## How the rules are enforced

When a sandbox tries to reach a network destination, the request is checked
against your rules. With organization policy, the same check also applies when
a sandbox mounts a host path or uses an MCP server.

The rules are enforced outside the sandbox, so an agent running inside it
can't change or bypass them. For where network rules are enforced as a request
leaves the sandbox, see [Networking](../sandboxes/architecture.md#networking).

Local policy and organization policy both follow this model. For how each one
works, see [How local policy works](local/_index.md#how-local-policy-works)
and
[How organization governance works](organization/_index.md#how-organization-governance-works).

Local policy, kit rules, and organization policy combine differently depending
on whether an organization policy is enforced for you. For more details, see
[How organization policy relates to local and kit rules](organization/policies/_index.md#how-organization-policy-relates-to-local-and-kit-rules).

## Get started

If you're a developer, start with
[Try out local policy](local/try-it-out.md), an example that allows and blocks
destinations for one sandbox.

If you're an administrator, start with
[Get started with organization governance](organization/get-started.md) to set
up licenses, sign-in enforcement, policies, and audit logs for your
organization.
