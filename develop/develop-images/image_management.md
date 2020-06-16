---
redirect_from:
- /reference/api/hub_registry_spec/
- /userguide/image_management/
- /engine/userguide/eng-image/image_management/
description: Documentation for docker Registry and Registry API
keywords: docker, registry, api,  hub
title: Manage images
---

The easiest way to make your images available for use by others inside or
outside your organization is to use a Docker registry, such as [Docker Hub](#docker-hub),
or by running your own [private registry](#docker-registry).


## Docker Hub

[Docker Hub](../../docker-hub/index.md) is a public registry managed by Docker,
Inc. It centralizes information about organizations, user accounts, and images.
It includes a web UI, authentication and authorization using organizations, CLI
and API access using commands such as `docker login`, `docker pull`, and `docker
push`, comments, stars, search, and more.

## Docker Registry

The Docker Registry is a component of Docker's ecosystem. A registry is a
storage and content delivery system, holding named Docker images, available in
different tagged versions. For example, the image `distribution/registry`, with
tags `2.0` and `latest`. Users interact with a registry by using docker push and
pull commands such as `docker pull myregistry.com/stevvooe/batman:voice`.

Docker Hub is an instance of a Docker Registry.

## Content Trust

When transferring data among networked systems, *trust* is a central concern. In
particular, when communicating over an untrusted medium such as the internet, it
is critical to ensure the integrity and publisher of all of the data a system
operates on. You use Docker to push and pull images (data) to a registry.
Content trust gives you the ability to both verify the integrity and the
publisher of all the data received from a registry over any channel.

See [Content trust](../../engine/security/trust/index.md) for information about
configuring and using this feature on Docker clients.
