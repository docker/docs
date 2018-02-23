---
description: Migrating from Docker Cloud
keywords: cloud, swarm, migration
title: Migrate Docker Cloud runtime apps
---

## Introduction

<span class="badge badge-warning">Important</span>  **The Docker Cloud runtime (and only the runtime) will be discontinued on April 30, 2018. To ensure that you do not lose your applications, you must migrate them from Docker Cloud to another platform. The documents in this section explain how.**

The following documents explain how to migrate your Docker Cloud applications to Docker Swarm (on Docker CE) and Kubernetes (on AKS and GKE).

At a high-level, migrating a Docker Cloud application requires that you:

- Set up a target environment (Docker Swarm or Kubernetes cluster)
- Convert your Docker Cloud YAML stack files
- Test the converted YAML stack files in the new environment
- Migrate your application from Docker Cloud to the new environment

### What does not change

How users and external systems interact with your Docker applications does _not_ change after migrating your applications from Docker Cloud. Your Docker images, and overall application functionality, remains the same. For example, if your application uses a Docker image called `myorg/webfe:v3`, and publishes container port `80` to external port `80`, none of this changes.

### What does change

What changes is the way in which you manage your Docker applications. For example, depending on how you migrate, you may lose the ability to:

- Deploy and manage applications with the Docker Cloud web UI
- Authorize users within the Docker platform
- **Autoredeploy** of your applications.

> Autoredeploy is a Docker Cloud feature that automatically updates your running application every time an updated image is pushed or built. It is not native to Docker CE, AKS or GKE, but you may be able to regain it using Docker Cloud auto-builds, and web-hooks on your Docker Cloud repository for the image, back to your CI/CD pipeline.

However, if you migrate to Docker Swarm on Docker CE, you should be able to keep all of your Docker Cloud _application features_ because Cloud and Swarm stacks have similar declarative syntax and lifecycle operations. Also, while you will lose some features as part of the migration, you may be able to add some advanced features that were not available as part of Docker Cloud.
