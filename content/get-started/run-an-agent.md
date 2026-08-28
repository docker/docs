---
title: Run an AI agent in a sandbox
linkTitle: Run an AI agent
description: Give an AI coding agent an isolated workspace, let it change and test code, and review its work without changing your host working tree.
keywords: Docker, get started, AI agents, Docker Sandboxes, sbx, sandbox
weight: 2
---

An AI coding agent is most useful when it can edit files, install tools, and run
commands. Running those actions directly on your host also increases the impact
of a mistake.

In this 10-minute tutorial, you'll give an agent a real task in an isolated
microVM. The agent works in a private clone, while you decide when its changes
cross back to your Git repository.

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

The prepared project gives the agent something concrete to change and test. You
will review one source file rather than learn the application itself.

## Start the agent

Launch Claude Code in a named sandbox with a private clone of the project:

```console
$ sbx run --clone --name agent-demo claude
```

On the first run, select the **Balanced** network policy. This policy blocks
network destinations by default and includes access to common development
services.

The CLI then creates the microVM and starts Claude Code. When the Claude Code
prompt appears, the sandbox is ready.

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

When Claude Code reports that the task is complete, confirm that it created the
`sandbox-demo` branch and that the backend tests passed.

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

The diff should show `Hello world!` replaced by `Hello from a sandbox!`. Fetching
the branch made the change available for review without applying it to your
working tree.

## Remove the sandbox

Remove the sandbox when you're finished:

```console
$ sbx rm agent-demo
```

This deletes the private clone, packages, images, containers, and other files
inside the microVM. The project on your host remains unchanged.

## What you proved

The agent changed code and ran tests with its own filesystem, Docker daemon, and
private working copy. Your host repository remained read-only to the agent, and
the network policy limited which destinations the sandbox could reach.

Using a sandbox narrows what an agent can access. Review the network rules,
credentials, and files you share before assigning sensitive work.

Continue with the [Docker Sandboxes documentation](/manuals/ai/sandboxes/_index.md)
to use another agent or configure tighter controls.
