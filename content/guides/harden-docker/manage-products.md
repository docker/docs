---
title: Manage Docker products
description: Learn how to manage organization access to Docker products included in your subscription.
weight: 5
---

In hardened security environments, you may or may not want to use all of the
Docker products available in your Business subscription.

This guide covers how to enable or disable the products and features
included in your subscription to best suit your organization's security needs.

## Docker products and features

By default, a Docker Business subscription comes with the following
products and features:

- [Docker Desktop](): The industry-leading
container-first development solution that includes, Docker Engine, Docker CLI,
Docker Compose, Docker Build/BuildKit, and Kubernetes.
- [Docker Hub](): The world's largest cloud-based container registry.
- [Docker Build Cloud](): Powerful cloud-based builders that accelerate build times
by up to 39x.
- [Docker Scout](): Tooling for software supply chain security that lets you quickly
assess image health and accelerate security improvements.
- [Testcontainers Cloud](): Container-based testing automation that provides faster
tests, a unified developer experience, and more.

## Manage access to Docker products

| Docker product | Default access | Can be disabled? | How to manage access |
|----------------|----------------|------------------|----------------------|
| Docker Desktop | Enabled | No | You can't technically disable Docker Desktop, but you can manage user access by enforcing sign-in so only organization members can use Docker Desktop. You can also manage Docker Desktop settings using Settings Management to control what settings and features your users have access to. |
| Docker Hub | Enabled | Yes | Use the Docker Admin Console to configure Registry Access Management or Image Access Management. |
| Docker Build Cloud | Enabled | Yes | Lock Docker Build Cloud for your organization. |
| Docker Scout | Enabled | Yes | |
| Testcontainers Cloud | Enabled | Yes | Lock Testcontainers Cloud for your organization. |
| Docker Hardened Images (DHI) | Disabled | Yes | Only available for users who have signed up for DHI. |
