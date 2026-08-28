---
title: Run an AI agent safely
linkTitle: Run an AI agent safely
description: Work with an AI coding agent as usual while Docker Sandboxes keeps its tools, containers, and system changes in a disposable environment.
keywords: Docker, get started, AI agents, Docker Sandboxes, sbx, sandbox, isolation
weight: 2
aliases:
  - /get-started/run-an-agent/
---

AI coding agents work best when they can install tools, run commands, and use
Docker. On your host, those actions can change system packages, interfere with
containers, or leave behind an environment you have to repair.

Docker Sandboxes gives the agent an isolated microVM without changing how you
work with it. You prompt the agent as usual. Its system changes stay in a
disposable environment that you can remove and recreate.

In this 15-minute tutorial, you'll let an agent install a package, start a
container, and write a file. Then you'll discard its environment and start
again from a clean state.

## Before you start

- [Install Docker Sandboxes](/manuals/ai/sandboxes/install.md) and run `sbx login`
- Have access to a Claude subscription

Docker Desktop and Docker Engine aren't required on the host.

## Create a workspace

Create an empty directory for the agent and open it:

```console
$ mkdir agent-sandbox-demo
$ cd agent-sandbox-demo
```

This directory is the agent's shared workspace. Changes made here appear on
your host so you can review and keep them. The agent doesn't run directly on
your host.

## Start the agent

Launch Claude Code in a sandbox:

```console
$ sbx run --name agent-demo claude
```

On your first run, select the **Balanced** network policy. It permits common
development services and blocks other destinations by default.

When Claude Code starts, enter `/login` and complete the browser sign-in. From
this point, using the agent feels the same as running it outside a sandbox.

## Give the agent a task

Send this prompt to Claude Code:

```text
Install the tree command system-wide using the system package manager. Start an
nginx:alpine container named sandbox-web on port 8080 and confirm that it
responds. Create a file named sandbox-result.txt in the current workspace that
summarizes what you did. Do not create any other project files.
```

The agent installs the package, pulls the image, starts the container, and
writes the result file. When it reports that the task is complete, enter
`/exit`.

## See what stayed in the sandbox

Read the file from your host:

```console
$ cat sandbox-result.txt
```

The file is visible because the workspace is shared. Now inspect the agent's
environment:

```console
$ sbx exec agent-demo sh -lc 'command -v tree && docker ps --filter name=sandbox-web --format "{{.Names}}"'
```

The command prints the path to `tree` and the name `sandbox-web`. Both exist in
the sandbox. The agent didn't install the package on your host or use a Docker
daemon on your host.

You interacted with Claude Code in the usual way. Underneath that experience,
the agent ran with its own system filesystem and Docker daemon. The network
policy also controlled which external destinations it could reach.

## Throw the environment away

Remove the sandbox:

```console
$ sbx rm agent-demo
```

This deletes the package, container, image, and every other change inside the
microVM. The shared workspace remains, including `sandbox-result.txt`.

Start the same agent again:

```console
$ sbx run --name agent-demo claude
```

Send this prompt:

```text
Check whether the tree command is installed and whether a container named
sandbox-web exists. Do not install or start anything.
```

The agent reports that neither exists. The command looks the same, but Docker
created a clean environment for the new session. Enter `/exit`, then remove the
new sandbox:

```console
$ sbx rm agent-demo
```

Remove the tutorial workspace when you're finished:

```console
$ rm sandbox-result.txt
$ cd ..
$ rmdir agent-sandbox-demo
```

## What you proved

You worked with an agent normally while its package installation, Docker
activity, and system changes stayed in a disposable environment. Only the file
written to the selected workspace appeared on your host. You also started over
from a clean environment without repairing or uninstalling anything.

Continue with the [Docker Sandboxes documentation](/manuals/ai/sandboxes/_index.md)
to choose another agent or define tighter network, filesystem, and tool
policies.
