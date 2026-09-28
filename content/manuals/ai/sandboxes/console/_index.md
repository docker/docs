---
title: Docker Agentic Platform
description: Run agents and tools in isolated cloud sandboxes with Docker Agentic Platform.
keywords: docker agentic platform, agents, kits, sandboxes, mcp, secrets, network policies
weight: 50
params:
  sidebar:
    badge:
      color: violet
      text: Experimental
grid:
  - title: Sign up
    description: Activate cloud access and review billing.
    icon: credit-card
    link: /ai/sandboxes/cloud-access/
  - title: Get started
    description: Start your first sandbox.
    icon: rocket-launch
    link: /ai/sandboxes/console/get-started/
  - title: Kits
    description: Find a kit or run your own public kit.
    icon: cube
    link: /ai/sandboxes/console/kits/
  - title: Sandboxes
    description: Access, pause, resume, and delete your sandboxes.
    icon: command-line
    link: /ai/sandboxes/console/sandboxes/
  - title: MCP
    description: Connect predefined or custom MCP servers.
    icon: cpu-chip
    link: /ai/sandboxes/console/mcp/
  - title: Secrets
    description: Manage model provider and service credentials.
    icon: key
    link: /ai/sandboxes/console/secrets/
  - title: Policies
    description: Control outbound network access from sandboxes.
    icon: shield-check
    link: /ai/sandboxes/console/policies/
  - title: FAQ
    description: Find answers about access, billing, and supported features.
    icon: question-mark-circle
    link: /ai/sandboxes/console/faq/
aliases:
  - /agentic-platform/
linkTitle: Use the Console
---

> [!NOTE]
> Docker Agentic Platform is experimental. Features and behavior may change.

Docker Agentic Platform lets you run agents and tools in isolated cloud
sandboxes. Your agent keeps working when you close the Console, disconnect
your computer, or put it to sleep.

The `sbx` CLI supports both
[local sandboxes](/manuals/ai/sandboxes/cli/get-started-local.md) and
[cloud sandboxes](/manuals/ai/sandboxes/cli/get-started-cloud.md). These pages describe
the web Console experience. The Console and cloud CLI share a secret store. See
[Secrets and credential injection](/manuals/ai/sandboxes/concepts/secrets.md).

In the Console, choose a [kit](/manuals/ai/sandboxes/console/kits.md) and configure
the sandbox's credentials, network access, tools, and compute size. Once it
starts, use its terminal to work with the agent. You can return to running or
paused sandboxes from **Sandboxes**.

You can reuse these settings across sandboxes:

- [MCP](/manuals/ai/sandboxes/console/mcp.md) connects external tools.
- [Secrets](/manuals/ai/sandboxes/console/secrets.md) provide credentials
  without placing their values inside a sandbox.
- [Network policies](/manuals/ai/sandboxes/console/policies.md) control
  which hosts and services a sandbox can reach.

To begin, [activate your subscription](/manuals/ai/sandboxes/cloud-access/_index.md#activate-cloud-access), then
[start a sandbox](/manuals/ai/sandboxes/console/get-started.md). You pay for compute by the second while your
sandbox runs. See [Signup and billing](/manuals/ai/sandboxes/cloud-access/_index.md) for account and payment
information.

{{< grid >}}
