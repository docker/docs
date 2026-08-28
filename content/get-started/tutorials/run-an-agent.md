---
title: Run an AI agent safely
linkTitle: Run an AI agent safely
description: Let an AI agent reset a real Compose project without risking the database in another checkout.
keywords: Docker, get started, AI agents, Docker Sandboxes, sbx, sandbox, isolation
weight: 2
aliases:
  - /get-started/run-an-agent/
---

AI coding agents work best when they can install tools, run commands, and use
Docker. The same access increases the effect of a mistaken assumption. For
example, an agent can reset the database for one checkout and accidentally
delete data from another checkout that Compose identifies as the same project.

Docker Sandboxes gives the agent its own Docker daemon and Linux environment.
You prompt the agent as usual, but its containers, packages, and system changes
stay inside a disposable microVM.

In this 15-minute tutorial, you'll run a real to-do application with saved data,
then ask an agent to reset a parallel checkout. The agent completes the task,
while the application on your host remains untouched.

## Before you start

- Install and start [Docker Desktop](../get-docker.md) or
  [Docker Engine](/manuals/engine/install/_index.md)
- [Install Docker Sandboxes](/manuals/ai/sandboxes/install.md) and run `sbx login`
- Install [Git](https://git-scm.com/downloads)
- Have access to a Claude subscription

Docker Sandboxes doesn't require a host Docker runtime for normal use. This
tutorial uses one to make the isolation visible with a real application on your
host.

## Create two checkouts

Agent tools and workspace managers often keep parallel checkouts under separate
workspace roots while preserving the repository name. Create that common layout
with two checkouts of Docker's sample to-do application:

```console
$ mkdir agent-sandbox-demo
$ cd agent-sandbox-demo
$ mkdir developer agent
$ git clone https://github.com/dockersamples/todo-list-app developer/todo-list-app
$ git clone https://github.com/dockersamples/todo-list-app agent/todo-list-app
```

The checkouts have different parent directories, but both end in
`todo-list-app`. That detail will matter later.

## Start your application

Open the developer checkout and start its application stack on your host:

```console
$ cd developer/todo-list-app
$ docker compose up -d
```

Wait for the API to become available:

```console
$ curl --retry 30 --retry-all-errors --retry-delay 2 -fsS http://localhost:3000/items
```

Open [http://localhost:3000](http://localhost:3000), then add a to-do named
`Prepare the release notes`.

The item is stored in the application's MySQL database. Confirm that the API
returns it:

```console
$ curl -fsS http://localhost:3000/items
```

Keep this application running.

## Create an environment for the agent

Move to the agent's checkout and create a sandbox without attaching to it:

```console
$ cd ../../agent/todo-list-app
$ sbx create --name onboarding-test claude .
```

On your first run, select the **Balanced** network policy. It permits common
development services and blocks other destinations by default.

The checkout is shared with the sandbox, so the agent can work with the project
files. The agent's operating system and Docker daemon are separate from your
host.

## Give the agent a normal task

Attach to Claude Code:

```console
$ sbx run --name onboarding-test
```

Enter `/login` and complete the browser sign-in. Then send this prompt:

```text
I'm testing the first-run experience in this checkout. Reset this checkout's
Compose environment, including its database, then start the application. Add a
to-do named "Review onboarding copy" through the application API and confirm
that it is the only item. Don't modify the source files. Handle the task end to
end.
```

Claude inspects the project, prepares a clean Compose environment, starts the
stack, and uses the API. These are the same commands and the same interaction
you would expect outside a sandbox. When Claude reports that the task is
complete, enter `/exit`.

## Preview what could have happened

You are still in the agent checkout. Ask Compose on your host to preview the
same kind of database reset:

```console
$ docker compose --dry-run down --volumes
```

The dry run identifies the running `todo-list-app` containers and its MySQL
volume for removal. Do not repeat the command without `--dry-run` yet.

Compose uses the [directory name as the default project name](/manuals/compose/how-tos/project-name.md).
Both checkouts are named `todo-list-app`, so on a shared Docker daemon they have
the same identity. An agent working directly on your host could interpret
“reset this checkout” correctly and still delete the database used by your
other checkout.

The sandbox has a private Docker daemon, so the agent's identically named
containers and volume never collided with the ones on your host.

## Verify both applications

Check the application on your host again:

```console
$ curl -fsS http://localhost:3000/items
```

The response still contains `Prepare the release notes`. Now query the
application inside the sandbox:

```console
$ sbx exec onboarding-test curl -fsS http://localhost:3000/items
```

This response contains only `Review onboarding copy`. The same source project
has two independent application environments, and the agent could reset only
its own.

The selected project directory is shared read-write, so source edits would
still appear on your host. Keep project files under version control and review
agent changes as usual. The sandbox isolates everything around that workspace:
the operating system, installed packages, Docker objects, and network policy.

## Throw the environment away

Remove the sandbox and everything the agent created inside it:

```console
$ sbx rm onboarding-test
```

The agent's application stack, database, images, packages, and other system
changes are deleted together. The source checkout and the application running
on your host remain.

When you're finished with the host application, remove its stack and volume:

```console
$ cd ../../developer/todo-list-app
$ docker compose down --volumes
```

You can also delete the `agent-sandbox-demo` directory and both checkouts.

## What you proved

You gave an agent a reasonable development task without changing how you
prompted or supervised it. A subtle scope collision could have deleted a real
database on a shared Docker daemon. The sandbox kept the agent's Docker and
system state in a disposable environment, while its intended project remained
available on your host.

Continue with the [Docker Sandboxes documentation](/manuals/ai/sandboxes/_index.md)
to choose another agent or define tighter network, filesystem, and tool
policies.
