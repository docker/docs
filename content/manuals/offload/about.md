---
title: About Docker Offload
linktitle: About
weight: 15
description: Learn about Docker Offload, its features, and how it works.
keywords: cloud, offload, vdi
---

Docker Offload is a fully managed service for building and running containers in
the cloud using the Docker tools you already know, including Docker Desktop, the
Docker CLI, and Docker Compose. It extends your local development workflow into a
scalable, cloud-powered environment, enabling developers to work efficiently even
in virtual desktop infrastructure (VDI) environments or systems that don't support
nested virtualization.



## Key features

Docker Offload includes the following capabilities to support modern container
workflows:

- Ephemeral cloud runners: Automatically provision and tear down cloud
  environments for each container session.
- Hybrid workflows: Seamlessly transition between local and remote execution
  using Docker Desktop or CLI.
- Secure communication: Use encrypted tunnels between Docker Desktop and cloud
  environments with support for secure secrets and image pulling.
- Port forwarding and bind mounts: Retain a local development experience even
  when running containers in the cloud.
- VDI-friendly: [Use Docker Desktop](../desktop/setup/vm-vdi.md) in virtual desktop environments or systems that
  don't support nested virtualization.

## Why use Docker Offload?

Docker Offload is designed to support modern development teams working across
local and cloud environments. It helps you:

- Offload heavy builds and runs to fast, scalable infrastructure
- Run containers that require more resources than your local setup can provide
- Use Docker Compose to manage complex, multi-service apps that need cloud
  resources
- Maintain consistent environments without managing custom infrastructure
- Develop efficiently in restricted or low-powered environments like VDIs

Docker Offload is ideal for high-velocity development workflows
that need the flexibility of the cloud without sacrificing the simplicity of
local tools.

## How Docker Offload works

Docker Offload replaces the need to build or run containers locally by connecting
Docker Desktop to secure, dedicated cloud resources.

### Running containers with Docker Offload

When you use Docker Offload to build or run containers, a Docker Desktop creates a secure
SSH tunnel to a Docker daemon running in the cloud. Your containers are started
and managed entirely in that remote environment.

Here's what happens:

1. Docker Desktop connects to the cloud and triggers container creation.
2. Docker Offload builds or pulls the required images and starts containers in the cloud.
3. The connection stays open while the containers run.
4. When the containers stop running, the environment shuts down and is cleaned
   up automatically.

This setup avoids the overhead of running containers locally and enables fast,
reliable containers even on low-powered machines, including machines that do not
support nested virtualization. This makes Docker Offload ideal for developers
using environments such as virtual desktops, cloud-hosted development machines,
or older hardware.

Despite running remotely, features like bind mounts and port forwarding continue
to work seamlessly, providing a local-like experience from within Docker Desktop
and the CLI.

Docker Offload automatically transitions between active and idle states based on
usage. You're only charged when actively building or running containers. When
idle for more than 5 minutes, the session ends and resources are cleaned up. For
details about how this works and how to configure idle timeout, see [Active and
idle states](configuration.md#understand-active-and-idle-states).

## What's next

Get hands-on with Docker Offload by following the [Docker Offload quickstart](/offload/quickstart/).