---
title: Docker Sandboxes
description: Run AI coding agents in isolated environments
keywords: docker sandboxes, sbx, ai agents, sandboxed agents, microVM
weight: 10
params:
  sidebar:
    group: AI and agents
---

Docker Sandboxes run AI coding agents in isolated microVM sandboxes. Each
sandbox gets its own Docker daemon, filesystem, and network — the agent can
build containers, install packages, and modify files without accessing host
resources beyond those you share.

> [!NOTE]
> The `sbx` CLI is free to use, including for commercial work. Only
> [organization governance](governance/) requires a separate paid subscription.

Organization admins can
[centrally manage sandbox network, filesystem, and MCP policies](governance/access-controls/organization.md),
so the same controls apply uniformly across every developer's machine.
Available on a separate paid subscription.

## Get started

Follow the [installation guide](install.md) to check the system requirements,
install the `sbx` CLI, and sign in.

Then launch an agent in a sandbox:

```console
$ cd ~/my-project
$ sbx run claude
```

See the [get started guide](get-started.md) for a first-session walkthrough, or
jump to the [usage guide](usage.md) for basic commands.

## Learn more

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
