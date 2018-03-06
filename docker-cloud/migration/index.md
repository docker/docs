---
description: Migrating from Docker Cloud
keywords: cloud, migration
title: Migrate Docker Cloud runtime apps
---

## Introduction

<span class="badge badge-warning">Important</span>  **Cluster management services in Docker Cloud are shutting down 11 May. You must migrate your applications from Docker Cloud to another platform and deregister your Swarms. Docker SaaS features are _not_ affected.**

The Docker Cloud runtime, and _only_ the runtime, is being discontinued. This means that you must migrate your applications to another platform, and if applicable, deregister your Swarms from Docker Cloud. The documents in this section explain how.

- [Migrate Docker Cloud stacks to Docker CE swarm](cloud-to-swarm){: target="_blank" class="_"}
- [Migrate Docker Cloud stacks to Azure Container Service](cloud-to-kube-aks){: target="_blank" class="_"}
- [Migrate Docker Cloud stacks to Google Kubernetes Engine](cloud-to-kube-gke){: target="_blank" class="_"}
- [Deregister Swarms on Docker Cloud](deregister-swarms){: target="_blank" class="_"}
- [Kubernetes primer](kube-primer){: target="_blank" class="_"}

## What stays the same

**How users and external systems interact with your Docker applications**. Your Docker images, autobuilds, and overall application functionality, remains the same. For example, if your application uses a Docker image called `myorg/webfe:v3`, and publishes container port `80` to external port `80`, none of this changes.

Docker Cloud SaaS features stay! We are _not_ removing automated builds and registry storage services.

## What changes

**How you manage your Docker applications**. We are removing node cluster management (specifically, Docker Cloud runtime) from Docker Cloud. Depending on how you migrate, you may lose the ability to:

- Deploy and manage applications with the Docker Cloud web UI
- Authorize users in the Docker platform with their Docker ID
- Autoredeploy your applications.

> **Autoredeploy options**: Autoredeploy is a Docker Cloud feature that automatically updates running applications every time you push an image. It is not native to Docker CE, AKS or GKE, but you may be able to regain it with Docker Cloud auto-builds, using web-hooks from the Docker Cloud repository for your image back to the CI/CD pipeline in your dev/staging/production environment.

While you will lose some features as part of the migration, you may be able to regain them elsewhere and even add advanced features that were not available as part of Docker Cloud.
