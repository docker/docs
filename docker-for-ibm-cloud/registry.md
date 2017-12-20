---
description: Set up a registry to use images in your clusters
keywords: ibm, ibm cloud, registry, dtr, iaas, tutorial
title: Registry overview
---
Learn about image registries, installing the registry CLI and setting up a namespace, logging in with private registry credentials, and creating a swarm service with a registry image.

## Image registries

Images are typically stored in a registry that can either be accessible by the public (public registry) or set up with limited access for a small group of users (private registry). Public registries, such as Docker Hub, can be used to get started with Docker and to create your first containerized app in a swarm.

When it comes to enterprise applications, use a private registry, like Docker Trusted Registry or IBM Cloud Container Registry, to protect your images from being used and changed by unauthorized users. Private registries must be set up by the registry admin to ensure that the credentials to access the private registry are available to the swarm users.

To deploy a container image in Docker swarm mode, you create a service that uses the image.

## Docker Trusted Registry

[Docker Trusted Registry (DTR)](/datacenter/dtr/2.4/guides/) is a private registry that runs on a Docker EE cluster. Once deployed, you can use the DTR GUI or Docker CLI to manage your Docker images. Set up DTR to save images on external storage. Use DTR when you need a secure image registry that's integrated with Docker EE UCP.

Before you can use DTR, [configure DTR to use IBM Cloud Object Storage](dtr-ibm-cos.md) to securely store your images externally.

## IBM Cloud Container Registry

[IBM Cloud Container Registry](https://console.bluemix.net/docs/services/Registry/registry_overview.html#registry_overview) is a global, multi-tenant private image registry that you can use to safely store and share Docker images with other users in your IBM Cloud account.

Images in the registry are automatically scanned by Vulnerability Adviser so that you build and scale your containers using secure images. Take control of your image usage and billing by setting quota limits to manage storage and pull traffic. Use IBM Cloud Container Registry when you need to share images across clusters and with users who do not have access to your swarm.

Learn how to [use IBM Cloud Container Registry with Docker EE for IBM Cloud](ibm-registry.md).
