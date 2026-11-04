---
title: Core concepts
weight: 15
description: Learn about Docker Offload v1, its features, and how it works.
keywords: cloud, offload, vdi
aliases:
- /offload/v1/about/
---

{{< summary-bar feature_name="Docker Offload v1" >}}

Docker Offload v1 is a fully managed service for building and running containers in
the cloud using the Docker tools you already know, including Docker Desktop, the
Docker CLI, and Docker Compose. It extends your local development workflow into a
scalable, cloud-powered environment, enabling developers to work efficiently even
in virtual desktop infrastructure (VDI) environments or systems that don't support
nested virtualization.

Docker Offload v1 replaces the need to build or run containers locally by connecting
Docker Desktop to secure, dedicated cloud resources.

## Running containers with Docker Offload v1

When you use Docker Offload v1 to build or run containers, Docker Desktop creates a secure
SSH tunnel to a Docker daemon running in the cloud. Your containers are started
and managed entirely in that remote environment.

Here's what happens:

1. Docker Desktop connects to the cloud and triggers container creation.
2. Docker Offload v1 builds or pulls the required images and starts containers in the cloud.
3. The connection stays open while the containers run and you remain active.
4. When the containers stop running, the environment shuts down and is cleaned
   up automatically.

This setup avoids the overhead of running containers locally and enables fast,
reliable containers even on low-powered machines, including machines that do not
support nested virtualization. This makes Docker Offload v1 ideal for developers
using environments such as virtual desktops, cloud-hosted development machines,
or older hardware.

Despite running remotely, features like bind mounts and port forwarding continue
to work, providing a local-like experience from within Docker Desktop
and the CLI.

## Cloud resources

Docker Offload v1 uses cloud hosts with 4 vCPUs and 8 GiB of memory. If you have
different requirements, <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_offload_about" class="link" rel="noopener">contact Docker</a> to explore options.

## Session management and idle state

Docker Offload v1 uses session management and idle state policies to ensure [fair use](/manuals/offload/v1/_index.md#fair-use) of cloud resources across all users.

Each user can run one Docker Offload v1 session at a time. When Docker Desktop is in an **Offload idle** state, it waits for activity on the Docker API and only connects to a cloud environment when needed. Once connected, the session moves to an **Offload running** state and stays connected as long as Docker detects activity. Activity includes any Docker API call, a running container, or an active build.

### When you'll see a prompt

While Docker Offload v1 is running, Docker Desktop shows prompts in the Dashboard to check if you're still active. Prompts appear in two cases:

1. No activity is detected for more than 3 minutes.
2. The session has been running for a long time.

When a prompt appears, you can:
   - Select **Ask me again later** to confirm you're still active and continue your session.
   - Select **Idle now** to return to an idle state immediately.
   - Do nothing, and the session returns to an idle state automatically.

### What happens when your session goes idle

After your session returns to an idle state, there is a 5-minute grace period. You can resume the session during this time by running any Docker command.

> [!IMPORTANT]
> If the idle period exceeds 5 minutes without activity, the session is terminated. Docker Offload v1 environments are ephemeral, so the remote environment and any containers, images, or volumes in it are deleted. To keep work between sessions, push images to a registry such as [Docker Hub](/manuals/docker-hub/_index.md) before your session ends.

### Long session prompts

Long session prompts appear every 3 hours during a session. After 8 hours of cumulative usage in a day, prompts appear every hour. The 8-hour counter resets at the start of each day.

## What's next

Get hands-on with Docker Offload v1 by following the [Docker Offload v1
quickstart](/manuals/offload/v1/quickstart.md).