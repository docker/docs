---
description: Migrating from Docker Cloud
keywords: cloud, migration
title: Migration overview
---

## Introduction

<span class="badge badge-warning">Important</span>  **Cluster and application management services in Docker Cloud are shutting down on May 21. You must migrate your applications from Docker Cloud to another platform and deregister your Swarms.**

The Docker Cloud runtime is being discontinued. This means that you will no longer be able to manage your nodes, swarm clusters, and the applications that run on them in Docker Cloud. To protect your applications, you must must migrate them to another platform, and if applicable, deregister your Swarms from Docker Cloud. The documents in this section explain how.

- [Migrate Docker Cloud stacks to Docker CE swarm](cloud-to-swarm){: target="_blank" class="_"}
- [Migrate Docker Cloud stacks to Azure Container Service](cloud-to-kube-aks){: target="_blank" class="_"}
- [Migrate Docker Cloud stacks to Google Kubernetes Engine](cloud-to-kube-gke){: target="_blank" class="_"}
- [Deregister Swarms on Docker Cloud](deregister-swarms){: target="_blank" class="_"}
- [Kubernetes primer](kube-primer){: target="_blank" class="_"}

## What stays the same

**How users and external systems interact with your Docker applications**. Your Docker images, autobuilds, automated tests, and overall application functionality remain the same. For example, if your application uses a Docker image called `myorg/webfe:v3`, and publishes container port `80` to external port `80`, none of this changes.

Docker Cloud SaaS features stay! We are _not_ removing automated builds and registry storage services.

## What changes

**How you manage your applications**. We are removing cluster management and the ability to deploy and manage Docker Cloud stacks. As part of the migration, you will no longer be able to:

- Manage your nodes and clusters in Docker Cloud.
- Deploy and manage applications from the Docker Cloud web UI.
- Autoredeploy your applications.
- Integrate users with other parts the Docker platform with their Docker ID.

> **Autoredeploy options**: Autoredeploy is a Docker Cloud feature that automatically updates running applications every time you push an image. It is not native to Docker CE, AKS or GKE, but you may be able to regain it with Docker Cloud auto-builds, using web-hooks from the Docker Cloud repository for your image back to the CI/CD pipeline in your dev/staging/production environment.
