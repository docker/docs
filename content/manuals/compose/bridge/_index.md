---
description: Learn how Compose Bridge transforms Docker Compose files into Kubernetes manifests for seamless platform transitions
keywords: docker compose bridge, compose to kubernetes, docker compose kubernetes integration, docker compose kustomize, compose bridge docker desktop
title: Overview of Compose Bridge
linkTitle: Compose Bridge
weight: 50
---

{{< summary-bar feature_name="Compose bridge" >}}

Compose Bridge converts your Docker Compose configuration into platform-specific formats—primarily Kubernetes manifests. The default transformation generates Kubernetes manifests and a Kustomize overlay which are designed for deployment on Docker Desktop with Kubernetes enabled.  

It's a flexible tool that lets you either take advantage of the [default transformation](usage.md) or [create a custom transformation](customize.md) to suit specific project needs and requirements.  

Compose Bridge significantly simplifies the transition from Docker Compose to Kubernetes, making it easier for you to leverage the power of Kubernetes while maintaining the simplicity and efficiency of Docker Compose.

## How it works

Compose Bridge uses transformations to let you convert a Compose model into another form. 

A transformation is packaged as a Docker image that receives the fully resolved Compose model as `/in/compose.yaml` and can produce any target format file under `/out`.

Compose Bridge provides its own transformation for Kubernetes using Go templates, so that it is easy to extend for customization by replacing or appending your own templates.

For more detailed information on how these transformations work and how you can customize them for your projects, see [Customize](customize.md).

## What's next?

- [Use Compose Bridge](usage.md)
- [Explore how you can customize Compose Bridge](customize.md)