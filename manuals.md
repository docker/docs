---
title: Product and tool manuals
notoc: true
---

The Docker platform is comprised of a family of tools and products. After
learning the general principles of the Docker workflow under [Guides](/), you
can find the documentation for these tools and products here.

## Supported platforms

Docker CE and EE are available on multiple platforms, on cloud and on-premises.
Use the following tables to choose the best installation path for you.

### Desktop

{% include docker_desktop_matrix.md %}

### Docker Certified Infrastructure

{% include docker_cloud_matrix.md %}

### Server

{% include docker_platform_matrix.md %}

## Tools

Free downloadables that help your device use Docker containers.

| Tool                                      | Description                                                                                            |
|:------------------------------------------|:-------------------------------------------------------------------------------------------------------|
| [Docker Compose](/compose/overview/)      | Enables you to define, build, and run multi-container applications                                     |
| [Docker Machine](/machine/overview/)      | Enables you to provision and manage Dockerized hosts                                                   |
| [Docker Notary](/notary/getting_started/) | Allows the signing of container images to enable Docker Content Trust                                  |
| [Docker Registry](/registry/)             | The software that powers Docker Hub and Docker Store, Registry stores and distributes container images |

## Products

Commercial Docker products that turn your container-based solution into a
production-ready application.

| Product                                                      | Description                                                                                                       |
|:-------------------------------------------------------------|:------------------------------------------------------------------------------------------------------------------|
| [Docker Cloud](/docker-cloud/)                               | Manages multi-container applications and host resources running on a cloud provider (such as Amazon Web Services) |
| [Universal Control Plane (UCP)](/datacenter/ucp/2.2/guides/) | Manages your Docker swarm on-premise, or on the cloud                                                             |
| [Docker Trusted Registry (DTR)](/datacenter/dtr/2.3/guides/) | Securely stores and scans your Docker images                                                                      |
| [Docker Store](/docker-store/)                               | Public, Docker-hosted registry that distributes free and paid images from various publishers                      |

## Superseded products and tools

* [Docker Hub](/docker-hub/) - Superseded by Docker Store and Docker Cloud
* [Docker Swarm](/swarm/overview/) - Functionality folded directly into native Docker, no longer a standalone tool
* [Docker Toolbox](/toolbox/overview/) - Superseded by Docker for Mac and Windows
