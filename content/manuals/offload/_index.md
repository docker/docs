---
title: Docker Offload
weight: 20
description: Find documentation on Docker Offload to help you offload containers to the cloud.
keywords: build, cloud, cloud build, remote builder, offload
params:
  sidebar:
    group: Application development
grid:
- title: Quickstart
  description: Get up and running with Docker Offload in just a few steps.
  icon: rocket-launch
  link: /offload/quickstart/
- title: About
  description: Learn about Docker Offload and how it works.
  icon: information-circle
  link: /offload/about/
- title: Manage cloud sandboxes
  description: Start, stop, inspect, and resume the cloud sandbox that runs your containers.
  icon: command-line
  link: /offload/sandboxes/
- title: Usage
  description: Learn how Docker Offload is billed and how to monitor your spend.
  icon: heart
  link: /offload/usage/
- title: Optimize
  description: Improve performance and cost efficiency in Docker Offload.
  icon: bolt
  link: /offload/optimize/
- title: Troubleshoot
  description: Learn how to troubleshoot issues with Docker Offload.
  icon: bug-ant
  link: /offload/troubleshoot/
- title: Feedback
  description: Provide feedback on Docker Offload.
  icon: chat-bubble-left
  link: /offload/feedback/
aliases:
- /harmonia/
---

{{< summary-bar feature_name="Docker Offload" >}}

Docker Offload is a fully managed service that lets you offload building and
running containers to the cloud using the Docker tools you already know. It
runs your containers in a cloud sandbox instead of on your local machine, so
you can work on hardware that can't run your workload, in virtual desktop
infrastructure (VDI) environments, or on systems that don't support nested
virtualization.

Docker Offload is for individual developers, and it's included with a [Docker
Agentic Platform](/manuals/subscription-billing/plans/docker-agentic-platform.md) subscription. You pay for the time your cloud sandbox runs.

> [!NOTE]
>
> If you start Docker Offload with `docker offload start`, or with the cloud toggle
> in the Docker Desktop Dashboard header, your organization subscribes to
> [Docker Offload v1](v1/_index.md). To use the latest version of Docker Offload,
> see [Move to the current version of Docker
> Offload](v1/configuration.md#move-to-the-current-version-of-docker-offload).

In the following topics, learn about Docker Offload, how to set it up, use it
for your workflows, and troubleshoot common issues.

{{< grid >}}
