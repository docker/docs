---
description: Learn how Compose Bridge transforms Docker Compose files into Kubernetes manifests for seamless platform transitions
keywords: docker compose bridge, compose to kubernetes, docker compose kubernetes integration, docker compose kustomize, compose bridge docker desktop
title: Overview of Compose Bridge
linkTitle: Compose Bridge
weight: 50
---

{{< summary-bar feature_name="Compose bridge" >}}

Compose Bridge converts your Docker Compose configuration into platform-specific deployment formats such as Kubernetes manifests. By default, it generates:

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

Compose Bridge also supports applications that use LLMs via Docker Model Runner.

For more details, see [Use Model Runner](use-model-runner.md).

## Apply organizational standards at scale 

Docker Compose Bridge can help companies to apply organizational policies through custom transformation templates that are applied when converting `compose.yaml` files to Kubernetes manifests. Developers continue to write standard Compose files, while required security contexts, resource limits, labels, and network policies are applied automatically and consistently during conversion.

When requirements change, platform teams update the templates once and all teams inherit the changes on their next deployment. This separation lets developers focus on application configuration while platform teams or organization administrators manage governance.

## What's next?

- [Use Compose Bridge](usage.md)
- [Explore how you can customize Compose Bridge](customize.md)
