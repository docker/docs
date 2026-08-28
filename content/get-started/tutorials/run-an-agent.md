---
title: Run an AI agent safely
linkTitle: Run an AI agent safely
description: Run an AI coding agent in an isolated environment with its own Docker daemon and system tools.
keywords: Docker, get started, AI agents, Docker Sandboxes, sbx, sandbox, isolation
weight: 2
aliases:
  - /get-started/run-an-agent/
---

AI coding agents are more useful when they can run commands, install tools, and
use Docker. Giving an autonomous process those capabilities directly on your
host increases the impact of a mistake.

In this 10-minute tutorial, you'll run an agent in an isolated microVM, ask it
to build and test a project, inspect what it created, and remove the entire
environment.

## Before you start

- [Install Docker Sandboxes](/manuals/ai/sandboxes/install.md) and run `sbx login`
- Install [Git](https://git-scm.com/downloads)
- Have access to a Claude subscription

Docker Desktop and Docker Engine aren't required on the host.

## Get a project

Clone a prepared application and open its directory:

```console
$ git clone https://github.com/docker/getting-started-todo-app
$ cd getting-started-todo-app
```

The project gives the agent a real Dockerfile and test suite to work with.

## Start the agent

Launch Claude Code in a named sandbox:

```console
$ sbx run --name agent-demo claude
```

On your first run, select the **Balanced** network policy. It permits common
development services and blocks other destinations by default.

When Claude Code starts, enter `/login` and complete the browser sign-in. The
agent is then running inside the sandbox with the current project as its
workspace.

## Give the agent a task

Send this prompt to Claude Code:

```text
Build the Dockerfile's test stage as an image tagged sandbox-demo. Report
whether the tests pass. Do not change any project files.
```

The agent can use Docker without controlling a daemon on your host. Its images,
containers, installed packages, and system files stay inside the sandbox.

When the agent reports that the build and tests passed, enter `/exit`.

## Check the boundary

Confirm that the project files are unchanged:

```console
$ git status --short
```

The command produces no output. Inspect the image in the sandbox's private
Docker daemon:

```console
$ sbx exec agent-demo docker image inspect sandbox-demo --format '{{ join .RepoTags ", " }}'
```

The command prints the `sandbox-demo` tag. The agent created the image inside
the microVM, not on your host.

The selected project directory is shared read-write with the sandbox so an
agent can work on it. The agent runtime, Docker daemon, installed packages, and
system filesystem remain inside the microVM. Review project changes as you
would review any code change before keeping them.

## Remove the sandbox

Delete the sandbox and everything the agent created inside it:

```console
$ sbx rm agent-demo
```

The project directory remains on your host. The sandbox's image, Docker daemon,
packages, and filesystem are removed together.

## What you proved

You gave an AI agent the tools to build and test a real project without giving
it the rest of your host system. The sandbox contained the agent's runtime and
Docker activity, while the shared workspace kept its code changes visible for
review.

Continue with the [Docker Sandboxes documentation](/manuals/ai/sandboxes/_index.md)
to use another supported agent or configure organization policies.
