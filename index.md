---
description: Home page for Docker's documentation
keywords: Docker, documentation, manual, guide, reference, api
layout: docs
title: Welcome to the Docs
---

# Welcome to the Docs

Docker provides a way to run applications securely isolated in a container, packaged with all its dependencies and libraries. Because your application can always be run with the environment it expects right in the build image, testing and deployment is simpler than ever, as your build will be fully portable and ready to run as designed in any environment. And because containers are lightweight and run without the extra load of a hypervisor, you can run many applications that all rely on different libraries and environments on a single kernel, each one never interfering with the other. This allows you to get more out of your hardware by shifting the "unit of scale" for your application from a virtual or physical machine, to a container instance.

## Typical Docker Platform Workflow

1. Get your code and its dependencies into Docker [containers](engine/getstarted/step_two.md):
   - [Write a Dockerfile](engine/getstarted/step_four.md) that specifies the execution
     environment and pulls in your code.
   - If your app depends on external applications (such as Redis, or
     MySQL), simply [find them on a registry such as Docker Hub](docker-hub/repos.md), and refer to them in
     [a Docker Compose file](compose/overview.md), along with a reference to your application, so they'll run
     simultaneously.
     - Software providers also distribute paid software via the [Docker Store](https://store.docker.com).
   - Build, then run your containers on a virtual host via [Docker Machine](machine/overview.md) as you develop.
2. Configure [networking](engine/tutorials/networkingcontainers.md) and
   [storage](engine/tutorials/dockervolumes.md) for your solution, if needed.
3. Upload builds to a registry ([ours](engine/tutorials/dockerrepos.md), [yours](docker-trusted-registry/index.md), or your cloud provider's), to collaborate with your team.
4. If you're gonna need to scale your solution across multiple hosts (VMs or physical machines), [plan
   for how you'll set up your Swarm cluster](engine/swarm/key-concepts.md) and [scale it to meet demand](engine/swarm/swarm-tutorial/index.md).
   - Note: Use [Universal Control Plane](ucp/overview.md) and you can manage your
     Swarm cluster using a friendly UI!
5. Finally, deploy to your preferred
   cloud provider (or, for redundancy, *multiple* cloud providers) with [Docker Cloud](docker-cloud/overview.md). Or, use [Docker Datacenter](https://www.docker.com/products/docker-datacenter), and deploy to your own on-premise hardware.


## Components

### [![docker-for-mac](images/docker-for-mac.png) Docker for Mac](docker-for-mac/)