---
title: Run an AI agent safely
linkTitle: Run an AI agent safely
description: Give an AI agent a real debugging task inside a disposable Docker Sandbox.
keywords: Docker, get started, AI agents, Docker Sandboxes, sbx, sandbox, isolation
weight: 2
aliases:
  - /get-started/run-an-agent/
---

AI coding agents work best when they can install tools, run commands, and use
Docker. The same access makes ordinary requests risky. Asking an agent to
reproduce a problem from an empty database might lead it to delete persistent
Docker data on your machine.

Docker Sandboxes gives the agent its own Docker daemon and Linux environment.
You work with the agent as usual, but its containers, packages, and system
changes stay inside a disposable microVM.

In this tutorial, you'll give an agent a real Compose application and ask it to
reproduce a clean first run. The workflow feels normal, but its side effects
remain inside the sandbox.

## Before you start

- [Install Docker Sandboxes](/manuals/ai/sandboxes/install.md) and run `sbx login`
- Install [Git](https://git-scm.com/downloads)
- Have access to a Claude subscription

You don't need Docker Desktop or Docker Engine on your host.

## Get the application

Clone Docker's sample to-do application:

```console
$ git clone https://github.com/dockersamples/todo-list-app
$ cd todo-list-app
```

The project contains a Node.js application and a Compose file for the
application and its MySQL database. You don't need to install Node.js or MySQL.

## Start the agent

Run Claude Code in a sandbox with the project as its workspace:

```console
$ sbx run --name todo-debug claude
```

On your first run, select the **Balanced** network policy. It permits common
development services and blocks other destinations by default.

Enter `/login` and complete the browser sign-in. From this point, using Claude
Code feels the same as running it directly on your machine.

## Give the agent a debugging task

Send this prompt:

```text
Reproduce this application's first-run experience from a clean local state.
Reset its Compose environment, including persistent database data, then start
the application. Through its API, add to-dos named "Draft release notes" and
"Review onboarding copy". Verify that the API returns exactly those two items.
Don't modify the source files. Report what you did and what you found.
```

Claude inspects the project, prepares a clean Compose environment, starts the
application, and uses its API. It can pull images, start containers, install
tools, and troubleshoot without asking you to prepare its machine.

When Claude reports that the task is complete, enter `/exit`.

## See the isolated environment

List the containers the agent started:

```console
$ sbx exec todo-debug docker ps --format 'table {{.Names}}\t{{.Status}}'
```

The output shows the application and MySQL containers. Query the application
inside the sandbox:

```console
$ sbx exec todo-debug curl -fsS http://localhost:3000/items
```

The response contains the two requested to-dos.

The agent used Docker as it would on a developer machine, including resetting
persistent data. The difference is scope: it could see only the sandbox's
private Docker daemon. Your host's containers, images, volumes, packages, and
system files were outside the environment.

The project directory is shared read-write, so source edits would still appear
on your host. Keep project files under version control and review agent changes
as usual.

## Throw the environment away

Remove the sandbox:

```console
$ sbx rm todo-debug
```

This deletes the application stack, MySQL database, images, installed packages,
and every other change inside the microVM. The project files remain on your
host. Running the same `sbx run` command creates a clean environment for the
next task.

## What you proved

You gave an agent the freedom to complete a normal development task without
running its tools and Docker workloads on your host. The agent could reset
Docker state and install what it needed, while the sandbox constrained the
effects to an environment you could discard with one command.

Continue with the [Docker Sandboxes documentation](/manuals/ai/sandboxes/_index.md)
to choose another agent or define tighter network, filesystem, and tool
policies.
