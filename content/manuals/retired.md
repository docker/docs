---
title: Deprecated and retired Docker products and features
linkTitle: Deprecated products and features
description: |
  Explore deprecated and retired Docker features, products, and open source
  projects, including details on transitioned tools and archived initiatives.
params:
  sidebar:
    group: Products
aliases:
  - /cloud/
  - /cloud/aci-compose-features/
  - /cloud/aci-container-features/
  - /cloud/aci-integration/
  - /cloud/ecs-architecture/
  - /cloud/ecs-compose-examples/
  - /cloud/ecs-compose-features/
  - /cloud/ecs-integration/
  - /engine/context/aci-integration/
  - /engine/context/ecs-integration/
  - /machine/
  - /machine/drivers/hyper-v/
  - /machine/get-started/
  - /machine/install-machine/
  - /machine/overview/
  - /registry/
  - /registry/compatibility/
  - /registry/configuration/
  - /registry/deploying/
  - /registry/deprecated/
  - /registry/garbage-collection/
  - /registry/help/
  - /registry/insecure/
  - /registry/introduction/
  - /registry/notifications/
  - /registry/recipes/
  - /registry/recipes/apache/
  - /registry/recipes/nginx/
  - /registry/recipes/osx-setup-guide/
  - /registry/spec/
  - /registry/spec/api/
  - /registry/spec/auth/
  - /registry/spec/auth/jwt/
  - /registry/spec/auth/oauth/
  - /registry/spec/auth/scope/
  - /registry/spec/auth/token/
  - /registry/spec/deprecated-schema-v1/
  - /registry/spec/implementations/
  - /registry/spec/json/
  - /registry/spec/manifest-v2-1/
  - /registry/spec/manifest-v2-2/
  - /registry/spec/menu/
  - /registry/storage-drivers/
  - /registry/storage-drivers/azure/
  - /registry/storage-drivers/filesystem/
  - /registry/storage-drivers/gcs/
  - /registry/storage-drivers/inmemory/
  - /registry/storage-drivers/oss/
  - /registry/storage-drivers/s3/
  - /registry/storage-drivers/swift/
  - /toolbox/
  - /toolbox/overview/
  - /toolbox/toolbox_install_mac/
  - /toolbox/toolbox_install_windows/
  - /desktop/features/dev-environments/
  - /desktop/features/dev-environments/create-dev-env/
  - /desktop/features/dev-environments/set-up/
  - /desktop/features/dev-environments/share/
  - /desktop/features/dev-environments/dev-cli/
  - /desktop/dev-environments/
---

This document provides an overview of Docker features, products, and
open-source projects that have been deprecated, retired, or transitioned.

> [!NOTE]
>
> This page does not cover deprecated and removed Docker Engine features.
> For a detailed list of deprecated Docker Engine features, refer to the
> [Docker Engine Deprecated Features documentation](/manuals/engine/deprecated.md).

## Products and features

Support for these deprecated or retired features is no longer provided by
Docker, Inc. The projects that have been transitioned to third parties continue
to receive updates from their new maintainers.

### Docker Machine

Docker Machine was a tool for provisioning and managing Docker hosts across
various platforms, including virtual machines and cloud providers. It is no
longer maintained, and users are encouraged to use [Docker Desktop](/manuals/desktop/_index.md)
or [Docker Engine](/manuals/engine/_index.md) directly on supported platforms.
Machine's approach to creating and configuring hosts has been superseded by
more modern workflows that integrate more closely with Docker Desktop.

### Docker Toolbox

Docker Toolbox was used on older systems where Docker Desktop could not run. It
bundled Docker Machine, Docker Engine, and Docker Compose into a single
installer. Toolbox is no longer maintained and is effectively replaced by
[Docker Desktop](/manuals/desktop/_index.md) on current systems. References to
Docker Toolbox occasionally appear in older documentation or community
tutorials, but it is not recommended for new installations.

### Docker Cloud integrations

Docker previously offered integrations for Amazon's Elastic Container Service
(ECS) and Azure Container Instances (ACI) to streamline container workflows.
These integrations have been deprecated, and users should now rely on native
cloud tools or third-party solutions to manage their workloads. The move toward
platform-specific or universal orchestration tools reduced the need for
specialized Docker Cloud integrations.

You can still view the relevant documentation for these integrations in the
[Compose CLI repository](https://github.com/docker-archive/compose-cli/tree/main/docs).

### Docker Enterprise Edition

Docker Enterprise Edition (EE) was Docker's commercial platform for deploying
and managing large-scale container environments. It was acquired by Mirantis in
2019, and users looking for enterprise-level functionality can now explore
Mirantis Kubernetes Engine or other products offered by Mirantis. Much of the
technology and features found in Docker EE have been absorbed into the Mirantis
product line.

> [!NOTE]  
> For information about enterprise-level features offered by Docker today,
> see the [Docker Business subscription](https://www.docker.com/pricing/).

### Docker Data Center and Docker Trusted Registry

Docker Data Center (DDC) was an umbrella term that encompassed Docker Universal
Control Plane (UCP) and Docker Trusted Registry (DTR). These components
provided a full-stack solution for managing containers, security, and registry
services in enterprise environments. They are now under the Mirantis portfolio
following the Docker Enterprise acquisition. Users still encountering
references to DDC, UCP, or DTR should refer to Mirantis's documentation for
guidance on modern equivalents.

### Dev Environments

Dev Environments was a feature introduced in Docker Desktop that allowed
developers to spin up development environments quickly. It was deprecated and removed from Docker Desktop version 4.42 and later. Similar workflows can be achieved through
Docker Compose or by creating custom configurations tailored to specific
project requirements.

### GitHub Copilot extension

The Docker for GitHub Copilot extension integrated Docker capabilities with
GitHub Copilot Chat, helping developers containerize applications, generate
Docker assets, and analyze vulnerabilities through conversational prompts. The
extension was available in early access on the GitHub Marketplace. GitHub
[deprecated Copilot Extensions](https://github.blog/changelog/2025-09-24-deprecate-github-copilot-extensions-github-apps/)
which led to the retirement of the Docker for GitHub Copilot extension. If
you're looking for AI-assisted Docker workflows, explore the Docker MCP Toolkit
and MCP Catalog, or use Ask Gordon in Docker Desktop and the Docker CLI.

## Open source projects

Several open-source projects originally maintained by Docker have been
archived, discontinued, or transitioned to other maintainers or organizations.

### Registry (now CNCF Distribution)

The Docker Registry served as the open-source implementation of a container
image registry. It was donated to the Cloud Native Computing Foundation (CNCF)
in 2019 and is maintained under the name "Distribution." It remains a
cornerstone for managing and distributing container images.

[CNCF Distribution](https://github.com/distribution/distribution)

### Docker Compose v1 (replaced by Compose v2)

Docker Compose v1 (`docker-compose`), a Python-based tool for defining
multi-container applications, has been superseded by Compose v2 (`docker
compose`), which is written in Go and integrates with the Docker CLI. Compose
v1 is no longer maintained, and users should migrate to Compose v2.

[Compose v2 Documentation](/manuals/compose/_index.md)

### InfraKit

InfraKit was an open-source toolkit designed to manage declarative
infrastructure and automate container deployments. It has been archived, and
users are encouraged to explore tools such as Terraform for infrastructure
provisioning and orchestration.

[InfraKit GitHub Repository](https://github.com/docker/infrakit)

### Docker Notary (now CNCF Notary)

Docker Notary was a system for signing and verifying the authenticity of
container content. It was donated to the CNCF in 2017 and continues to be
developed as "Notary." Users seeking secure content verification should consult
the CNCF Notary project.

[CNCF Notary](https://github.com/notaryproject/notary)

### SwarmKit

SwarmKit powers Docker Swarm mode by providing orchestration for container
deployments. While Swarm mode remains functional, development has slowed in
favor of Kubernetes-based solutions. Individuals evaluating container
orchestration options should investigate whether SwarmKit meets modern workload
requirements.

[SwarmKit GitHub Repository](https://github.com/docker/swarmkit)
