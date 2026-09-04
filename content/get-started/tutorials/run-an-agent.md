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
agent gets a private environment with its own operating system and Docker
daemon. Your project remains available on your host, while tools the agent
installs and system changes it makes stay inside an environment you can
discard.

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

Imported skills become available to supported agents across your sandboxes.
This feature is experimental. See
[Share agent skills](/manuals/ai/sandboxes/workflows/agent-skills.md) for
supported agents and security considerations.

## Choose how to authenticate

Your agent needs access to a model provider. Subscription sign-in uses OAuth,
while API keys are stored on your host and supplied to the agent through the
sandbox proxy. Choose your agent and authentication method for the preparation
steps:

{{< sandbox-auth >}}

If you already exported a supported API key in your shell, you can import it
instead of setting it separately:

```console
$ sbx secret import --dry-run
$ sbx secret import
```

The dry run shows which exported keys `sbx` found. The import command prompts
you to confirm each key before storing it. If the dry run finds nothing, use
the `sbx secret set` command from the picker. See
[Credentials](/manuals/ai/sandboxes/configuration/credentials.md) for other
secret sources and [supported coding agents](/manuals/ai/sandboxes/agents/_index.md)
for agent-specific authentication.

## Run your agent

Open your project and start your preferred agent. This example uses Codex:

```console
$ cd ~/my-project
$ sbx run --name my-project codex
```

Replace `codex` with another supported agent identifier, such as `claude`,
`copilot`, `cursor`, or `gemini`.

The first time you run a sandbox, `sbx` asks you to choose a default network
policy. This policy controls which external services your sandboxes can reach.
Select **Balanced** to permit common development services and block other
destinations by default. You can change these rules later with
[`sbx policy`](/manuals/ai/sandboxes/governance/access-controls/local.md).

The built-in integrations start coding agents in their full-autonomy mode. For
example, Codex bypasses approvals, Claude Code skips permission prompts, and
Cursor, Copilot, and Gemini use YOLO mode. You don't need to add those flags.

Give the agent the same task you would give it on your host. Source changes
appear in your working tree, so you can inspect them with your usual tools:

```console
$ git diff
```

Your project directory is the exception to the sandbox boundary. It is shared
read-write, so the agent can modify or delete its files and you can see those
changes immediately. Keep your work under version control. Tools the agent
installs and changes to the sandbox's operating system stay inside the sandbox.

## Return or start over

Exit the agent when you're finished. Reconnect to the same environment later:

```console
$ sbx run --name my-project
```

When you want a clean environment, remove the sandbox:

```console
$ sbx rm my-project
```

Removing the sandbox deletes the environment and everything installed inside
it. It doesn't delete your project directory.

## What changed

You kept your project, agent, skills, credentials, and prompting workflow. The
only essential change was launching the agent with `sbx run`. That moved its
full-autonomy execution into an environment you control and can throw away.

## What's next

Continue with the Docker Sandboxes manuals:

- [Manage your sandboxes](/manuals/ai/sandboxes/usage.md) with day-to-day
  commands
- [Configure your coding agent](/manuals/ai/sandboxes/agents/_index.md) for its
  authentication and settings
- [Review the security defaults](/manuals/ai/sandboxes/security/defaults.md) for
  workspace, network, and credential access
