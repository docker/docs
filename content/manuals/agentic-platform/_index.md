---
title: Docker Agentic Platform
description: Run agents and tools in isolated, hosted sandboxes with Docker Agentic Platform.
keywords: docker agentic platform, agents, sandboxes, mcp, secrets, network policies
weight: 5
sitemap: false
params:
  sidebar:
    group: AI and agents
    badge:
      color: violet
      text: Experimental
grid:
  - title: Get started
    description: Choose an agent environment and start a sandbox.
    icon: rocket-launch
    link: /agentic-platform/get-started/
  - title: Sandboxes
    description: Work with hosted agent environments.
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
    description: Review launch scope and product boundaries.
    icon: question-mark-circle
    link: /agentic-platform/faq/
---

> [!NOTE]
> Docker Agentic Platform is experimental. Features and behavior may change.

Docker Agentic Platform runs agents and agent-powered tools in isolated
sandboxes on Docker-managed cloud infrastructure. An active workload is not
tied to your computer remaining awake or connected. You can leave the Console
and return to the sandbox while the agent continues working.

For sandboxes that run on your development machine through the `sbx` CLI, see
[Docker Sandboxes](/manuals/ai/sandboxes/_index.md).

From the web Console, choose the type of sandbox to run and configure its model
credential, network access, tools, and compute. Docker creates the sandbox and
opens a live terminal for interacting with the agent. The **Sandboxes** page
provides one place to return to and manage your running and paused workloads.

Account-level configuration can be reused across sandboxes:

- [MCP](/manuals/agentic-platform/mcp.md) connects external tools.
- [Secrets](/manuals/agentic-platform/secrets.md) provide credentials
  without placing their values inside a sandbox.
- [Network policies](/manuals/agentic-platform/policies.md) control
  outbound destinations.

To begin, open [Docker Agentic Platform](https://agentic-platform.docker.com/)
and sign in with your Docker account. Docker meters sandbox compute per second.
For account and payment information, see [Docker Billing](/billing/).

{{< grid >}}
