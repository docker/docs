---
title: Docker Agentic Platform
description: Run agents and tools in isolated cloud sandboxes with Docker Agentic Platform.
keywords: docker agentic platform, agents, kits, sandboxes, mcp, secrets, network policies
weight: 5
params:
  sidebar:
    group: AI and agents
    badge:
      color: violet
      text: Experimental
grid:
  - title: Sign up
    description: Activate cloud access and review billing.
    icon: credit-card
    link: /agentic-platform/signup/
  - title: Get started
    description: Start your first sandbox.
    icon: rocket-launch
    link: /agentic-platform/get-started/
  - title: Kits
    description: Find a kit or run your own public kit.
    icon: cube
    link: /agentic-platform/kits/
  - title: Sandboxes
    description: Access, pause, resume, and delete your sandboxes.
    icon: command-line
    link: /agentic-platform/sandboxes/
  - title: MCP
    description: Connect predefined or custom MCP servers.
    icon: cpu-chip
    link: /agentic-platform/mcp/
  - title: Secrets
    description: Manage model provider and service credentials.
    icon: key
    link: /agentic-platform/secrets/
  - title: Policies
    description: Control outbound network access from sandboxes.
    icon: shield-check
    link: /agentic-platform/policies/
  - title: FAQ
    description: Find answers about access, billing, and supported features.
    icon: question-mark-circle
    link: /agentic-platform/faq/
---

> [!NOTE]
> Docker Agentic Platform is experimental. Features and behavior may change.

Docker Agentic Platform lets you run agents and tools in isolated cloud
sandboxes. Your agent keeps working when you close the Console, disconnect
your computer, or put it to sleep.

For sandboxes that run on your development machine through the `sbx` CLI, see
[Docker Sandboxes](/manuals/ai/sandboxes/_index.md).

In the Console, choose a [kit](/manuals/agentic-platform/kits.md) and configure
the sandbox's credentials, network access, tools, and compute size. Once it
starts, use its terminal to work with the agent. You can return to running or
paused sandboxes from **Sandboxes**.

You can reuse these settings across sandboxes:

- [MCP](/manuals/agentic-platform/mcp.md) connects external tools.
- [Secrets](/manuals/agentic-platform/secrets.md) provide credentials
  without placing their values inside a sandbox.
- [Network policies](/manuals/agentic-platform/policies.md) control
  which hosts and services a sandbox can reach.

To begin, [activate your subscription](signup.md#activate-cloud-access), then
[start a sandbox](get-started.md). You pay for compute by the second while your
sandbox runs. See [Signup and billing](signup.md) for account and payment
information.

{{< grid >}}
