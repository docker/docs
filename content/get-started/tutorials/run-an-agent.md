---
title: Run an AI agent safely
linkTitle: Run an AI agent safely
description: Give an AI coding agent a normal workflow while Docker Sandboxes keeps the side effects of a destructive cleanup task contained.
keywords: Docker, get started, AI agents, Docker Sandboxes, sbx, sandbox, isolation
weight: 2
aliases:
  - /get-started/run-an-agent/
---

AI coding agents work best when they can install tools, run commands, and use
Docker. The same access also increases the effect of a bad assumption. A request
to “start from a clean Docker environment” can remove stopped containers,
images, volumes, and build cache that belong to other projects.

Docker Sandboxes gives the agent an isolated microVM without changing how you
work with it. You prompt the agent as usual, but its system changes and Docker
activity stay inside a disposable environment.

In this 15-minute tutorial, you'll give an agent a realistic clean-build task.
The task has an unwanted side effect, but the damage stays inside the sandbox.

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

This directory is the agent's shared workspace. Files created here appear on
your host so you can review and keep them. The agent itself doesn't run on your
host.

## Set up an unrelated workload

Create the sandbox without attaching to the agent:

```console
$ sbx create --name agent-demo claude .
```

On your first run, select the **Balanced** network policy. It permits common
development services and blocks other destinations by default.

Create a stopped container that represents another project's database:

```console
$ sbx exec agent-demo docker run --name other-project-database alpine sh -c \
  'echo "uncommitted records" > /records.txt'
$ sbx exec agent-demo docker ps -a --filter name=other-project-database
```

The output shows `other-project-database` with an `Exited` status. Its writable
container layer still contains `/records.txt`. This is a stand-in for unrelated
work that might exist in a developer's Docker environment.

## Start the agent

Attach to Claude Code:

```console
$ sbx run --name agent-demo
```

Enter `/login` and complete the browser sign-in. From this point, using Claude
Code feels the same as running it outside a sandbox.

## Give the agent a clean-build task

Send this prompt to Claude Code:

```text
Create a small website that displays "Hello from a sandbox". Before building,
reset Docker to a clean state by removing all unused containers, images,
volumes, and build cache so no old state can affect the result. This is a
disposable environment, so proceed without asking for confirmation.

Create only Dockerfile and index.html in the workspace. Build an image tagged
sandbox-web, run it in a container named sandbox-web on port 8080, and confirm
that the page responds.
```

The agent cleans Docker, creates the two project files, builds the image, and
starts the website. When it reports that the task is complete, enter `/exit`.

Nothing about the agent interaction felt unusual. The cleanup instruction also
looked reasonable for an ephemeral build environment.

## See the side effect

Confirm that the requested website is running:

```console
$ sbx exec agent-demo docker ps --filter name=sandbox-web --format '{{.Names}}'
```

The command prints `sandbox-web`. Now look for the unrelated stopped container:

```console
$ sbx exec agent-demo sh -lc \
  'docker container inspect other-project-database >/dev/null 2>&1 || echo "other-project-database was removed"'
```

The command reports that `other-project-database` was removed. Its uncommitted
record disappeared with it. The agent followed the cleanup request, but it
couldn't distinguish disposable Docker objects from another project's work.

On a host Docker daemon, the same cleanup could affect every local project. In
this tutorial, it affected only the sandbox's private daemon. Your host
containers, images, volumes, packages, and system files were outside the blast
radius.

The **Balanced** network policy also constrained external access underneath the
agent session. You can change network, filesystem, and tool policies without
changing how developers prompt their agents.

## Throw the environment away

The two requested project files are visible on your host:

```console
$ ls
Dockerfile  index.html
```

Remove the sandbox:

```console
$ sbx rm agent-demo
```

This deletes the website container, images, build cache, and every other change
inside the microVM. Running `sbx run --name agent-demo claude` from the same
directory creates a clean environment while keeping the project files.

Remove the tutorial workspace when you're finished:

```console
$ rm Dockerfile index.html
$ cd ..
$ rmdir agent-sandbox-demo
```

## What you proved

You worked with an agent normally and gave it a plausible instruction with a
destructive side effect. The sandbox limited that mistake to a disposable
environment. The files the agent intentionally created stayed available for
review, while removing the sandbox cleared its system and Docker state.

Continue with the [Docker Sandboxes documentation](/manuals/ai/sandboxes/_index.md)
to choose another agent or define tighter network, filesystem, and tool
policies.
