---
description: Learn how Compose Bridge transforms Docker Compose files into Kubernetes manifests for seamless platform transitions
keywords: docker compose bridge, compose to kubernetes, docker compose kubernetes integration, docker compose kustomize, compose bridge docker desktop
title: Overview of Compose Bridge
linkTitle: Compose Bridge
weight: 50
---

{{< summary-bar feature_name="Compose bridge" >}}

Compose Bridge converts your Docker Compose configuration into platform-specific deployment formats such as Kubernetes manifests. By default, it geneterates:

- Kubernetes manifests 
- A Kustomize overlay

These outputs are ready for deployment on Docker Desktop with [Kubernetes enabled](/manuals/desktop/settings-and-maintenance/settings.md#kubernetes).  

Compose Bridge helps you bridge the gap between Compose and Kubernetes, making it easier to adopt Kubernetes while keeping the simplicity and efficiency of Compose.

It's a flexible tool that lets you either take advantage of the [default transformation](usage.md) or [create a custom transformation](customize.md) to suit specific project needs and requirements. 

## How it works

Compose Bridge uses transformations to convert a Compose model into another form. 

A transformation is packaged as a Docker image that receives the fully resolved Compose model as `/in/compose.yaml` and can produce any target format file under `/out`.

Compose Bridge provides its own transformation for Kubernetes using Go templates, so that it is easy to extend for customization by replacing or appending your own templates.

For more detailed information on how these transformations work and how you can customize them for your projects, see [Customize](customize.md).

With Composer version X, Compose Bridge now also supports applications that use LLMs via a model runner such as Docker Model Runner.

When a service includes the `models` top-level element, Compose Bridge automatically generates environment variables that point to model endpoints and model names.  
Depending on your deployment needs:
- On Docker Desktop, containers can access Model Runner running on the host during local development.
- On Kubernetes, Compose Bridge can deploy an in-cluster Model Runner as part of the generated manifests.

For more details, see [Use Model Runner](use-model-runner.md).

## What's next?

- [Use Compose Bridge](usage.md)
- [Explore how you can customize Compose Bridge](customize.md)