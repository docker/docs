---
title: Docker Sandboxes
description: Run AI coding agents in isolated environments
keywords: docker sandboxes, sbx, ai agents, sandboxed agents, microVM
weight: 10
params:
  sidebar:
    group: AI and agents
---

Docker Sandboxes run AI coding agents in isolated environments on your machine
or on Docker-managed cloud infrastructure. Use the `sbx` CLI to create and
manage either kind of sandbox.

The `sbx` CLI and local sandbox compute are free to use, including for commercial
work. Cloud compute is metered through a
[Docker Agentic Platform plan](/manuals/subscription-billing/plans/docker-agentic-platform.md).
Model-provider charges are separate.

Organization admins can
[centrally manage sandbox network, filesystem, and MCP policies](governance/access-controls/organization.md),
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

## Learn more

The following guides describe local sandbox workflows. For cloud workflows,
see [Cloud sandboxes](cloud/).

- [Agents](agents/) — supported agents and per-agent configuration
- [Workflows](workflows/) — patterns for Git, local development,
  authentication, agent skills, and automation
- [Configuration](configuration/) — manage credentials, declare project
  environments, turn on GPU passthrough, and configure an upstream proxy
- [Integrations](integrations/) — connect editors and apps like VS Code and
  Cursor to a sandbox over SSH
- [MCP gateway](mcp-gateway.md) — register MCP servers and connect them to
  sandboxed agents
- [Customize](customize/) — reusable templates and declarative kits for
  extending or tailoring sandboxes
- [Architecture](architecture.md) — microVM isolation, workspace mounting,
  networking
- [Security](security/) — isolation model, credential handling, and
  network policies
- [CLI reference](/reference/cli/sbx/) — full list of `sbx` commands and options
- [Troubleshooting](troubleshooting.md) — common issues and fixes
- [FAQ](faq.md) — login requirements, telemetry, etc

## Feedback

Your feedback shapes what gets built next. If you run into a bug, hit a
missing feature, or have a suggestion, open an issue at
[github.com/docker/sbx-releases/issues](https://github.com/docker/sbx-releases/issues).
