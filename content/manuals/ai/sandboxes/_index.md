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
machine. If you're building with agents like Claude Code, Sandboxes provides a
secure way to give agents autonomy without compromising your system.

## Why use Docker Sandboxes

AI agents need to execute commands, install packages, and test code. Running
them directly on your host machine means they have full access to your files,
processes, and network. Docker Sandboxes isolates agents in microVMs, each with
its own Docker daemon. Agents can spin up test containers and modify their
environment without affecting your host.

You get:

- Agent autonomy without host system risk
- Private Docker daemon for running test containers
- File sharing between host and sandbox
- Network access control

For a comparison between Docker Sandboxes and other approaches to isolating
coding agents, see [Comparison to alternatives](./architecture.md#comparison-to-alternatives).

> [!NOTE]
> MicroVM-based sandboxes require macOS or Windows (experimental). Linux users
> can use legacy container-based sandboxes with
> [Docker Desktop 4.57](/desktop/release-notes/#4570).

## How to use sandboxes

To create and run a sandbox:

```console
$ docker sandbox run claude ~/my-project
```

This command creates a sandbox for your workspace (`~/my-project`) and starts
the Claude Code agent inside it. The agent can now work with your code, install
tools, and run containers inside the isolated sandbox.

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

- **Claude Code** - Anthropic's coding agent
- **Codex** - OpenAI's Codex agent (partial support; in development)
- **Gemini** - Google's Gemini agent (partial support; in development)
- **cagent** - Docker's [cagent](/ai/cagent/) (partial support; in development)
- **Kiro** - by AWS (partial support; in development)

## Get started

Head to the [Get started guide](get-started.md) to run your first sandboxed agent.

## Troubleshooting

See [Troubleshooting](./troubleshooting) for common configuration errors, or
report issues on the [Docker Desktop issue tracker](https://github.com/docker/desktop-feedback).
