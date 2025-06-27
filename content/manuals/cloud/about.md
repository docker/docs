---
title: About Docker Cloud
linktitle: About
weight: 15
description: Learn about Docker Cloud, its features, and how it works.
keywords: cloud, build, remote builder
---

Docker Cloud is a fully managed service for building and running containers in
the cloud using the Docker tools you already know, including Docker Desktop, the
Docker CLI, and Docker Compose. It extends your local development workflow into a
scalable, cloud-powered environment, so you can offload compute-heavy tasks,
accelerate builds, and securely manage container workloads across the software
lifecycle.

Docker Cloud also supports GPU-accelerated instances, allowing you to
containerize and run compute-intensive workloads such as Docker Model Runner and
other machine learning or data processing tasks that benefit from GPU.

You can use Docker Cloud in following ways:

- In Cloud mode, where you use Docker Desktop with cloud-based resources. This
  is ideal for virtual desktop environments (VDIs) where nested virtualization
  isn't supported, or when you need more CPU, memory, or GPU than your local
  machine can provide. In this mode, both builds and container runs happen in
  the cloud, but Docker Desktop maintains a local-like experience. To get
  started, see [Docker Cloud quickstart](/cloud/quickstart/).

- For only builds, without running containers in the cloud. This lets you offload image
  builds to Docker Cloud while continuing to run containers locally. It's useful
  when you want faster, consistent builds but don’t need to run containers in
  the cloud. To get started, see [Build with Docker Cloud](/cloud/build/).

- In CI environments where builds are performed entirely in the cloud. This lets
  you have fast, consistent, and scalable builds without the need to manage your
  own runners or infrastructure. To get started, see [Use Docker Cloud in
  CI](/cloud/ci-build/).

## Key features

Docker Cloud includes the following capabilities to support modern container
workflows:

- Cloud-based builds: Execute builds on remote, fully managed BuildKit instances
  with native support for multi-platform output.
- GPU acceleration: Use NVIDIA L4 GPU-backed environments for machine learning,
  media processing, and other compute-intensive workloads.
- Ephemeral cloud runners: Automatically provision and tear down cloud
  environments for each container session.
- Shared build cache: Speed up build times across machines and teammates with a
  smart, shared cache layer.
- Hybrid workflows: Seamlessly transition between local and remote execution
  using Docker Desktop, CLI, or CI tools.
- Secure communication: Use encrypted tunnels between Docker Desktop and cloud
  environments with support for secure secrets and image pulling.
- CI/CD integration: Trigger builds in CI pipelines using Buildx, GitHub
  Actions, or prebuilt integrations.
- Port forwarding and bind mounts: Retain a local development experience even
  when running containers in the cloud.
- VDI-friendly: Use Docker Cloud in virtual desktop environments or systems that
  don't support nested virtualization.

## Why use Docker Cloud?

Docker Cloud is designed to support modern development teams working across
local and cloud environments. It helps you:

- Use the same Docker workflows locally and in the cloud
- Offload builds to fast, high-performance infrastructure
- Run containers that require more CPU, memory, or GPU than your local setup can
  provide
- Run Docker Desktop in environments that don't support nested virtualization
- Speed up feedback loops by reducing wait times for builds and testing
  environments
- Avoid managing custom infrastructure while retaining full control over how
  containers are built and executed
- Ensure consistent, clean environments for every build or test
- Integrate easily with any CI system using simple scripts or prebuilt actions

Docker Cloud is ideal for hybrid teams that want to iterate quickly, test
reliably, and scale efficiently without compromising on developer experience.

## How Docker Cloud works

Docker Cloud replaces the need to build or run containers locally by connecting
Docker Desktop and your CI pipelines to secure, dedicated cloud resources.

### Building with Docker Cloud

When you use Docker Cloud for builds, the `docker buildx build` command sends
the build request to a remote BuildKit instance in the cloud, instead of
executing it locally. Your workflow stays the same, only the execution
environment changes.

The build runs on infrastructure provisioned and managed by Docker:

- Each cloud builder is an isolated Amazon EC2 instance with its own EBS volume
- Remote builders use a shared cache to speed up builds across machines and
  teammates
- Builds support native multi-platform output (for example, `linux/amd64`,
  `linux/arm64`)
- Build results are encrypted in transit and sent to your specified destination
  (such as a registry or local image store)

Docker Cloud manages the lifecycle of builders automatically. There's no need to
provision or maintain infrastructure.

> [!NOTE]
>
> Docker Cloud builders are currently hosted in the US East region. Users in
> other regions may experience increased latency.

### Running containers with Docker Cloud

When you use Docker Cloud to run containers, a Docker Desktop creates a secure
SSH tunnel to a Docker daemon running in the cloud. Your containers are started
and managed entirely in that remote environment.

Here's what happens:

1. Docker Desktop connects to the cloud and triggers container creation.
2. Docker Cloud pulls the required images and starts containers in the cloud.
3. The connection stays open while the containers run.
4. When the containers stop running, the environment shuts down and is cleaned
   up automatically.

This setup avoids the overhead of running containers locally and enables fast,
reliable containers even on low-powered machines, including machines that do not
support nested virtualization. This makes Docker Cloud ideal for developers
using environments such as virtual desktops, cloud-hosted development machines,
or older hardware.

Docker Cloud also supports GPU-accelerated workloads. Containers that require
GPU access can run on cloud instances provisioned with NVIDIA L4 GPUs for
efficient AI inferencing, media processing, and general-purpose GPU
acceleration. This enables compute-heavy workflows such as model evaluation,
image processing, and hardware-accelerated CI tests to run seamlessly in the
cloud.

Despite running remotely, features like bind mounts and port forwarding continue
to work seamlessly, providing a local-like experience from within Docker Desktop
and the CLI.

Docker Cloud provisions an ephemeral cloud environment for each session. The
environment remains active while you are interacting with Docker Desktop or
actively using containers. If no activity is detected for about 30 minutes, the
session shuts down automatically. This includes any containers, images, or
volumes in that environment, which are deleted when the session ends.

## What's next

Get hands-on with Docker Cloud by following the [Docker Cloud quickstart](/cloud/quickstart/).