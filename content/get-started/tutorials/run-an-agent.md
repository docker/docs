---
title: Run your coding agent in a sandbox
linkTitle: Sandbox a coding agent
description: Move your existing AI coding agent workflow into a disposable Docker Sandbox.
keywords: Docker, get started, AI agents, coding agents, Docker Sandboxes, sbx, sandbox, YOLO mode
weight: 2
aliases:
  - /get-started/run-an-agent/
---

Coding agents can work faster in full-autonomy modes that let them run commands,
install tools, and use Docker without stopping for approval. Giving an agent
that access directly on your machine also gives mistakes a wider reach.

Docker Sandboxes change where the agent runs, not how you work with it. The
agent gets a private microVM with its own operating system and Docker daemon.
Your project remains available on your host, while packages, containers, and
other system changes stay inside an environment you can discard.

In this tutorial, you'll move an existing coding-agent workflow into a sandbox.

## Before you start

- Have a project directory for the agent to work with
- Have access to a [supported coding agent](/manuals/ai/sandboxes/agents/_index.md)

You don't need Docker Desktop or Docker Engine on your host.

## Install and sign in

[Install Docker Sandboxes](/manuals/ai/sandboxes/install.md) for your operating
system, then sign in to Docker:

```console
$ sbx login
```

## Bring your agent setup

This step is optional. Skip it if you want to start with a clean agent setup.

If you use agent skills on your host, preview and import them:

```console
$ sbx skills import --dry-run
$ sbx skills import
```

Imported skills are stored separately from their host originals and shared
with supported sandboxes. This feature is experimental. See
[Share agent skills](/manuals/ai/sandboxes/workflows/agent-skills.md) for the
supported agents and trust boundary.

If your model-provider API keys are already exported in your shell, import them
into the `sbx` credential store:

```console
$ sbx secret import --dry-run
$ sbx secret import
```

You can skip this command when you use a subscription or OAuth. Your agent
prompts you to authenticate when it starts. See
[Credentials](/manuals/ai/sandboxes/configuration/credentials.md) for other
secret sources and provider-specific commands.

## Run your agent

Open your project and start your preferred agent. This example uses Codex:

```console
$ cd ~/my-project
$ sbx run --name my-project codex
```

Replace `codex` with another supported agent identifier, such as `claude`,
`copilot`, `cursor`, or `gemini`.

On your first run, select the **Balanced** network policy. It permits common
development services and blocks other destinations by default.

The built-in integrations start coding agents in their full-autonomy mode. For
example, Codex bypasses approvals, Claude Code skips permission prompts, and
Cursor, Copilot, and Gemini use YOLO mode. You don't need to add those flags.

Give the agent the same task you would give it on your host. Source changes
appear in your working tree, so you can inspect them with your usual tools:

```console
$ git diff
```

The project directory is shared read-write. The agent can modify or delete
files in that directory, so keep your work under version control. Everything
around the project—the agent process, installed packages, Docker objects, and
system files—stays inside the sandbox.

## Return or start over

Exit the agent when you're finished. Reconnect to the same environment later:

```console
$ sbx run --name my-project
```

When you want a clean environment, remove the sandbox:

```console
$ sbx rm my-project
```

Removing the sandbox deletes its packages, containers, images, and system
changes. Your project directory and imported skills remain on your host.

## What changed

You kept your project, agent, skills, credentials, and prompting workflow. The
only essential change was launching the agent with `sbx run`. That moved its
full-autonomy execution into an environment you control and can throw away.

Continue with the [Docker Sandboxes documentation](/manuals/ai/sandboxes/_index.md)
to configure tighter network, filesystem, and tool policies.
