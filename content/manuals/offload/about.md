---
title: About Docker Offload
linktitle: About
weight: 15
description: Learn about Docker Offload, its features, and how it works.
keywords: cloud, vdi, offload
---

Docker Offload is a fully managed service designed for developers working in virtual
desktop infrastructure (VDI) environments or on systems that don't support nested
virtualization. It lets you use Docker Desktop to build and run containers in
the cloud, providing full Docker functionality without the need for local compute
resources or nested virtualization support.

## Key features

Docker Offload includes the following capabilities to support modern container workflows:

- VDI-friendly: [Run Docker Desktop for Windows in a VM or VDI environment](/desktop/setup/vm-vdi/) without nested
  virtualization.
- Ephemeral cloud runners: Automatically provision and tear down cloud environments for each container session.
- Secure communication: Use encrypted tunnels between Docker Desktop and the cloud environment.
- Port forwarding and bind mounts: Retain a local development experience even when running containers in the cloud.

## How Docker Offload works

Docker Offload replaces the need to build or run containers locally by connecting Docker Desktop to secure, dedicated
cloud resources.

When you use Docker Offload to build images or run containers, Docker Desktop creates a secure SSH tunnel to a Docker
daemon running in the cloud. Your images are built and containers are started and managed entirely in that remote
environment.

Here's what happens when you run containers:

1. Docker Desktop connects to the cloud and triggers container creation.
2. Docker Offload pulls or builds the required images and starts containers in the cloud.
3. The connection stays open while the containers run.
4. When the containers stop running, the environment shuts down and is cleaned up automatically.

This setup avoids the overhead of running containers locally and enables fast, reliable containers even on low-powered
machines, including machines that do not support nested virtualization. This makes Docker Offload ideal for developers
using environments such as virtual desktops, cloud-hosted development machines, or older hardware.

Despite running remotely, features like bind mounts and port forwarding continue to work seamlessly, providing a
local-like experience from within Docker Desktop and the CLI.

Docker Offload provisions an ephemeral cloud environment for each session. The environment remains active while you are
interacting with Docker Desktop or actively using containers. If no activity is detected for about 5 minutes, the
session shuts down automatically. This includes any containers, images, or volumes in that environment, which are
deleted when the session ends.

## What's next

[Contact sales](https://www.docker.com/pricing/contact-sales/) to subscribe to Docker Offload, and [get
started](/offload/quickstart) using it.