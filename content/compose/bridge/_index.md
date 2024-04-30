---
description: Understand what Compose Bridge is and how it can be useful
keywords: compose, orchestration, kubernetes, bridge
title: Overview of Compose Bridge
---

{{< include "compose-bridge-early-access.md" >}}

Compose Bridge lets you translate your Compose configuration file into configuration files for different platforms, primarily focusing on Kubernetes. It generates Kubernetes manifests and a Kustomize overlay when used together with Docker Desktop. 

It's a flexible tool that let you eaither take advantage of the [default transformation](usage.md) or [customize it](customize.md) to suit specific project needs and requirements.

Compose Bridge significantly simplifies the transition from Docker Compose to Kubernetes, making it easier for you to leverage the power of Kubernetes while maintaining the simplicity and efficiency of Docker Compose.

## How it works

Compose bridge uses transformations to let you convert a compose model into another form. 

A transformation is packaged as a Docker image that receive the fully-resolved Compose model as `/in/compose.yaml` and can produce any target format file under `/out`.

We provide our own transformation for Kubernetes using Go templates, so that it is easy to extend for customization by just replacing or appending your own templates.

For more detailed information on how these transformations work and how you can customize them for your projects, see [Customize](customize.md).

## Setup

To get started with Compose Bridge, you need to:

1. Download and install a version of Docker Desktop that supports Compose Bridge.
2. Navigate to the **Features in development** tab in **Settings**. 
3. From the **Experimental features** tab, select **Enable Compose Bridge**.

## Feedback

To give feedback, report bugs, or receive support, email `desktop-preview@docker.com`. There is also a dedicated Slack channel. To join, simply send an email to the provided address.

## What's next?

- [Use Compose Bridge](usage.md)
- [Explore how you can customize Compose Bridge](customize.md)
- [Explore the advanced integration](advanced-integration.md)
