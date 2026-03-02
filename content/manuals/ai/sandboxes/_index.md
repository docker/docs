---
title: Docker Sandboxes
description: Run AI agents in isolated environments
weight: 20
params:
  sidebar:
    group: AI
    badge:
      color: violet
      text: Experimental
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

Docker Sandboxes lets you run AI coding agents in isolated environments on your
machine. Sandboxes provides a secure way to give agents autonomy without
compromising your system.

## Why use Docker Sandboxes

AI agents need to execute commands, install packages, and test code. Running
them directly on your host machine means they have full access to your files,
processes, and network. Docker Sandboxes isolates agents in microVMs, each with
its own Docker daemon. Agents can spin up test containers and modify their
environment without affecting your host.

You get:

- Agent autonomy without host system risk
- YOLO mode by default - agents work without asking permission
- Private Docker daemon for running test containers
- File sharing between host and sandbox
- Network access control

For a comparison between Docker Sandboxes and other approaches to isolating
coding agents, see [Comparison to alternatives](./architecture.md#comparison-to-alternatives).

> [!NOTE]
> MicroVM-based sandboxes require macOS or Windows (experimental). Linux users
> can use legacy container-based sandboxes with
> [Docker Desktop 4.57](/desktop/release-notes/#4570).

> [!NOTE]
> For Windows users: On Windows Docker Sandboxes only work inside of a `cmd` or `PowerShell` window.
> At this time there is no support for WSL 2


## How to use sandboxes

To create and run a sandbox:

```console
$ cd ~/my-project
$ docker sandbox run claude
```

Replace `claude` with your [preferred agent](./agents/_index.md). This command
creates a sandbox for your workspace (`~/my-project`) and starts the agent. The
agent can now work with your code, install tools, and run containers inside the
isolated sandbox.

## How it works

Sandboxes run in lightweight microVMs with private Docker daemons. Each sandbox
is completely isolated - the agent runs inside the VM and can't access your
host Docker daemon, containers, or files outside the workspace.

Your workspace directory syncs between host and sandbox at the same absolute
path, so file paths in error messages match between environments.

Sandboxes don't appear in `docker ps` on your host because they're VMs, not
containers. Use `docker sandbox ls` to see them.

For technical details on the architecture, isolation model, and networking, see
[Architecture](architecture.md).

### Multiple sandboxes

Create separate sandboxes for different projects:

```console
$ docker sandbox run claude ~/project-a
$ docker sandbox run claude ~/project-b
```

Each sandbox is completely isolated from the others. Sandboxes persist until
you remove them, so installed packages and configuration stay available for
that workspace.

## Supported agents

Docker Sandboxes works with multiple AI coding agents:

- **Claude Code** - Anthropic's coding agent (production-ready)
- **Codex** - OpenAI's Codex agent (in development)
- **Copilot** - GitHub Copilot agent (in development)
- **Gemini** - Google's Gemini agent (in development)
- **OpenCode** - Multi-provider agent with TUI interface (in development)
- **cagent** - Docker's multi-provider coding agent (in development)
- **Kiro** - Interactive agent with device flow auth (in development)
- **Shell** - Minimal sandbox for manual agent installation

For detailed configuration instructions, see [Supported agents](agents/).

## Get started

Head to the [Get started guide](get-started.md) to run your first sandboxed agent.

## Troubleshooting

See [Troubleshooting](./troubleshooting) for common configuration errors, or
report issues on the [Docker Desktop issue tracker](https://github.com/docker/desktop-feedback).
