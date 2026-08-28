---
title: Run an AI agent in a sandbox
linkTitle: Run an AI agent
description: Give an AI coding agent an isolated workspace, let it change and test code, and review its work without changing your host working tree.
keywords: Docker, get started, AI agents, Docker Sandboxes, sbx, sandbox
weight: 2
---

In this 10-minute tutorial, you'll give an AI coding agent a real task inside an
isolated microVM. The agent can edit code, install tools, and run containers,
while your host working tree stays unchanged.

This path uses Claude Code and clone mode to provide one direct route through
the experience.

## Before you start

- [Install Docker Sandboxes](/manuals/ai/sandboxes/install.md) and run `sbx login`
- Install [Git](https://git-scm.com/downloads)
- Have access to a Claude subscription

The `sbx` CLI runs without Docker Desktop or Docker Engine on the host.

## Get a project

Clone the sample application and open its directory:

```console
$ git clone https://github.com/docker/getting-started-todo-app
$ cd getting-started-todo-app
```

## Start the agent

Launch Claude Code in a named sandbox with a private clone of the project:

```console
$ sbx run --clone --name agent-demo claude
```

On the first run, select the **Balanced** network policy. This policy blocks
network destinations by default and includes access to common development
services.

Enter `/login` inside Claude Code and complete the browser sign-in.

## Give the agent a task

Send the following prompt to Claude Code:

```text
Create a branch named sandbox-demo. Change the greeting in
backend/src/routes/getGreeting.js to "Hello from a sandbox!". Run the backend
tests. Do not push the branch.
```

The agent works in a private clone inside the sandbox. Your source repository is
mounted in the sandbox with read-only access.

## Check the boundary

After the agent finishes, enter `/exit`. Then check the host working tree:

```console
$ git status --short
```

The command produces no output. The agent edited only the private clone.

Fetch the branch from the sandbox and review the change without applying it to
your working tree:

```console
$ git fetch sandbox-agent-demo
$ git diff main..sandbox-agent-demo/sandbox-demo -- backend/src/routes/getGreeting.js
```

The sandbox acts as a Git remote. You choose when to fetch a branch and whether
to keep it.

## Remove the sandbox

Remove the sandbox when you're finished:

```console
$ sbx rm agent-demo
```

This deletes the private clone, packages, images, containers, and other files
inside the microVM. The project on your host remains unchanged.

## What made this safer

- The agent ran inside a microVM with its own filesystem and Docker daemon
- Clone mode gave the agent a private working copy and read-only access to the
  host repository
- The network policy limited which destinations the sandbox could reach

Using a sandbox narrows what an agent can access. Review the network rules,
credentials, and files you share before assigning sensitive work.

Continue with the [Docker Sandboxes documentation](/manuals/ai/sandboxes/_index.md)
to use another agent or configure tighter controls.
