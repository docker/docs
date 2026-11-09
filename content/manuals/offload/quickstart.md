---
title: Docker Offload quickstart
linktitle: Quickstart
weight: 10
description: Learn how to use Docker Offload to build and run your container images in the cloud.
keywords: cloud, quickstart, Docker Desktop, offload
---

{{< summary-bar feature_name="Docker Offload" >}}

This quickstart covers the steps developers need to get started with Docker
Offload.

It describes Docker Offload usage with `docker offload create` and
`docker --cloud`, which requires a Docker Agentic Platform subscription. If you
start Offload with `docker offload start` or the toggle in the Docker Desktop
Dashboard header, see [Docker Offload v1](/manuals/offload/v1/_index.md).

## Prerequisites

- [Docker Desktop](/desktop/) 4.97 or later. Docker recommends using the latest
  version to access the newest features and improvements in Docker Offload.
- A Docker Agentic Platform subscription. To subscribe, see [Signup and
  billing](/manuals/agentic-platform/signup.md).

## Step 1: Create a cloud sandbox

Run the following command to create a sandbox and connect Docker to it:

```console
$ docker offload create
```

Docker Offload creates a cloud sandbox and switches your Docker context to
`docker-cloud`. From this point, `docker` commands run in the cloud sandbox instead of
on your local machine.

To confirm which context is selected:

```console
$ docker context ls
```

## Step 2: Run a container

Run a container to verify that Docker Offload is working:

```console
$ docker run --rm hello-world
```

The container runs in the cloud sandbox, and its output is streamed to your
terminal. If Docker Offload is working, you see `Hello from Docker!` in the
output.

You can also send a single command to the sandbox without running `docker offload create` first.
The `--cloud` flag runs that command in the cloud and leaves your context
unchanged:

```console
$ docker --cloud build .
$ docker --cloud run --rm hello-world
```

If you don't already have a sandbox, `docker --cloud` creates one.

> [!IMPORTANT]
>
> `--cloud` is a global flag, so it goes before the subcommand. Passing it
> after, as in `docker run --cloud hello-world`, fails with
> `unknown flag: --cloud`.

## Step 3: Check your sandbox

You can view your sandbox, including its hourly price, in the [Docker
Agentic Platform Console](https://agentic-platform.docker.com/), where its name
starts with `offload-`.

You use one Docker Offload sandbox per account at a time. If you run `docker
offload create` while a sandbox exists, Docker Offload connects you to that
sandbox instead of creating a second one. This is how you work from more than
one machine.

To check the state of your sandbox and your connection to it:

```console
$ docker offload status

  Account: <your-account>

  Sandbox
    ID: sbx_001m49455jm42x4v6epbg8kd5fy
    State: Running
    Created: 4 minutes ago

  This machine
    State: Started (Active)
    Session started: 4 minutes ago
    Idle detection: on (5m)
    Kept awake by: checking activity
```

The output separates the sandbox, which every machine on your account shares,
from this machine's connection to it. For what each field means, see [Check the
status of your cloud
sandbox](sandboxes.md#check-the-status-of-your-cloud-sandbox).

## Step 4: Reach your local containers

Creating a sandbox doesn't stop your local containers, but it does change which
containers Docker commands apply to. While your context is `docker-cloud`:

- `docker ps -a` lists containers in the sandbox, not your local containers.
- Ports published by your local containers aren't reachable through the
  `docker-cloud` context.

To work with your local containers without leaving Docker Offload, pass the
local context explicitly:

```console
$ docker --context desktop-linux ps -a
```

Ports published by containers in the sandbox are forwarded to your machine, so
you reach them the same way you reach a local container's ports. For more
information, see [Ports and files](sandboxes.md#ports-and-files).


## Step 5: Stop Docker Offload

When you finish, stop Docker Offload:

```console
$ docker offload stop
```

Your Docker context switches back to your previous context immediately, and the
sandbox pauses in the background. Your local containers and their published
ports are available again.

Stopping doesn't discard your work. The sandbox keeps
its containers, images, and volumes, so you can resume where you left off. For
what's preserved and how billing stops, see [Manage
sandboxes](sandboxes.md).

If you forget to stop it, your sandbox pauses on its own after 5 minutes
without activity. To change that period or turn it off, see [Control when your
sandbox pauses](sandboxes.md#control-when-your-cloud-sandbox-pauses).

To pick up where you left off, run a command that needs the sandbox, such as
`docker start` or `docker run`. That resumes it. 

## Step 6: Delete your cloud sandbox

To delete the sandbox and
everything on it, run:

```console
$ docker offload rm
```

The sandbox is removed from the Docker Agentic Platform console as well.

## What's next

- [Explore the Docker Offload CLI reference](/reference/cli/docker/offload/).
- [Learn how Docker Offload works](core-concepts.md)
- [Manage your sandbox](sandboxes.md), including resuming and deleting it
- [Optimize your Docker Offload usage](optimize.md)
