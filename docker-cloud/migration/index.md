---
description: Migrating from Docker Cloud
keywords: cloud, swarm, migration
title: Migrate Docker Cloud runtime apps
---

## Introduction

<span class="badge badge-warning">Important</span>  **Cluster management** (_but not_ private repos, team management, and autobuilds) **in Docker Cloud will be discontinued on May 11, 2018. You must migrate your applications from Docker Cloud to another platform**. The documents in this section explain how.**

The following documents explain how to migrate your Docker Cloud applications to Docker Swarm (on Docker CE) and Kubernetes (on [Microsoft AKS](https://azure.microsoft.com/en-us/services/container-service/) and [Google GKE]()).

At a high-level, migrating a Docker Cloud application requires that you:

- **Set up** a target environment (Docker Swarm or Kubernetes cluster)
- **Convert** your Docker Cloud YAML stack files
- **Test** the converted YAML stack files in the new environment
- **Migrate** your applications from Docker Cloud to the new environment
- **Re-route** application endpoints

### What does not change

**How users and external systems interact with your Docker applications** does _not_ change after migrating your applications from Docker Cloud. Your Docker images, autobuilds, and overall application functionality, remains the same. For example, if your application uses a Docker image called `myorg/webfe:v3`, and publishes container port `80` to external port `80`, none of this changes.

### What does change

**How you manage your Docker applications** changes. We are removing node cluster management from Docker Cloud. Depending on how you migrate, you may lose the ability to:

- Deploy and manage applications with the Docker Cloud web UI
- Authorize users in the Docker platform with their Docker ID
- Autoredeploy your applications.

> **Autoredeploy** is a Docker Cloud feature that automatically updates running applications every time an updated image is pushed or built. It is not native to Docker CE, AKS or GKE, but you may be able to regain it with Docker Cloud auto-builds, using web-hooks from the image in your Docker Cloud repository back to the CI/CD pipeline in your dev/staging/production environment.

While you will lose some features as part of the migration, you may be able to regain them elsewhere and even add advanced features that were not available as part of Docker Cloud.
