---
title: Docker Sandboxes
description: Run AI coding agents in isolated environments
keywords: docker sandboxes, sbx, ai agents, sandboxed agents, microVM
weight: 10
params:
  sidebar:
    group: AI and agents
grid:
  - title: "Agents"
    description: "Choose an agent and configure it for local sandboxes."
    link: /ai/sandboxes/agents/
    icon: terminal
  - title: "Workflows"
    description: "Use Git, local development tools, authentication, and automation in sandboxes."
    link: /ai/sandboxes/workflows/
    icon: code
  - title: "Configuration"
    description: "Manage credentials, project environments, GPU passthrough, and upstream proxies."
    link: /ai/sandboxes/configuration/
    icon: settings
  - title: "Integrations"
    description: "Connect editors and applications to a sandbox over SSH."
    link: /ai/sandboxes/integrations/
    icon: code
  - title: "MCP gateway"
    description: "Register MCP servers and connect them to sandboxed agents."
    link: /ai/sandboxes/mcp-gateway/
    icon: box
  - title: "Kits"
    description: "Package tools and configuration into reusable sandbox environments."
    link: /ai/sandboxes/customize/
    icon: box
  - title: "Architecture"
    description: "Understand microVM isolation, workspace mounting, and networking."
    link: /ai/sandboxes/architecture/
    icon: book-open
  - title: "Security"
    description: "Review isolation, credential handling, and network policies."
    link: /ai/sandboxes/security/
    icon: settings
  - title: "CLI reference"
    description: "Find commands and options for the `sbx` CLI."
    link: /reference/cli/sbx/
    icon: terminal
  - title: "Troubleshooting"
    description: "Diagnose common issues and find fixes."
    link: /ai/sandboxes/troubleshooting/
    icon: search
  - title: "FAQ"
    description: "Find answers about login requirements, telemetry, and other sandbox behavior."
    link: /ai/sandboxes/faq/
    icon: book-open
---

Docker Sandboxes run AI coding agents in isolated environments on your machine
or on Docker-managed cloud infrastructure. Use the `sbx` CLI to create and
manage either kind of sandbox.

The `sbx` CLI and local sandbox compute are free to use, including for commercial
work. Cloud compute uses a
[pay-as-you-go subscription](/manuals/agentic-platform/signup.md#billing).
Model-provider charges are separate.

Organization admins can
[centrally manage sandbox network, filesystem, and MCP policies](../governance/organization/_index.md),
for local sandboxes across developer machines.
Available on a separate paid subscription.

## Get started

[Install the `sbx` CLI](install.md) and sign in, then choose where to run your
agent:

| Environment | Use it for | Start here |
| --- | --- | --- |
| Local sandboxes | Work with files and supported hardware on your machine | [Get started locally](get-started.md) |
| Cloud sandboxes | Run on Docker-managed compute without local virtualization | [Get started in the cloud](cloud/_index.md#get-started) |

The two environments have separate credentials, network policies, and lifecycle
controls. See [Compare local and cloud sandboxes](cloud/local-vs-cloud.md)
before adapting a workflow.

To create and manage cloud sandboxes from your application, see
[Sandboxes API and SDK](../sandboxes-api/_index.md).

## Explore the docs

The following guides describe local sandbox workflows. For cloud workflows,
see [Cloud sandboxes](cloud/).

{{< grid >}}

## Feedback

Your feedback shapes what gets built next. If you run into a bug, hit a
missing feature, or have a suggestion, open an issue at
[github.com/docker/sbx-releases/issues](https://github.com/docker/sbx-releases/issues).
