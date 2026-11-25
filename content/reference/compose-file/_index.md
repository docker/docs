---
description:
  Find the latest recommended version of the Docker Compose file format
  for defining multi-container applications.
keywords:
  docker compose file, docker compose yml, docker compose reference, docker
  compose cmd, docker compose user, docker compose image, yaml spec, docker compose
  syntax, yaml specification, docker compose specification
title: Compose file reference
toc_max: 4
toc_min: 1
grid:
  - title: Services top-level element
    description: Explore all services attributes for Compose.
    icon: wrench
    link: /reference/compose-file/services/
  - title: Container specification
    description: Find the attributes shared by services, jobs, and init containers.
    icon: container
    link: /reference/compose-file/container-spec/
  - title: Workload specification
    description: Find the attributes shared by services and jobs.
    icon: layers
    link: /reference/compose-file/workload-spec/
  - title: Networks top-level element
    description: Find all networks attributes for Compose.
    icon: server
    link: /reference/compose-file/networks/
  - title: Volumes top-level element
    description: Explore all volumes attributes for Compose.
    icon: database
    link: /reference/compose-file/volumes/
  - title: Include top-level element
    description: Reuse and modularize Compose configurations by including other Compose files.
    icon: file-text
    link: /reference/compose-file/include/
  - title: Configs top-level element
    description: Find out about configs in Compose.
    icon: settings
    link: /reference/compose-file/configs/
  - title: Secrets top-level element
    description: Learn about secrets in Compose.
    icon: lock
    link: /reference/compose-file/secrets/
  - title: Jobs top-level element
    description: Run containers to completion, on demand or on a schedule.
    icon: clock
    link: /reference/compose-file/jobs/
aliases:
  - /compose/compose-file/
  - /compose/reference/overview/
---

> **New to Docker Compose?**
>
> Find more information about the [key features and use cases of Docker Compose](/manuals/compose/intro/features-uses.md) or [try the quickstart guide](/manuals/compose/gettingstarted.md).

The Compose Specification is the latest and recommended version of the Compose file format. It helps you define a [Compose file](/manuals/compose/intro/compose-application-model.md) which is used to configure your Docker application’s services, networks, volumes, and more.

Legacy versions 2.x and 3.x of the Compose file format were merged into the Compose Specification. It is implemented in versions 1.27.0 and above (also known as Compose v2) of the Docker Compose CLI.

The Compose Specification on Docker Docs is the Docker Compose implementation. If you wish to implement your own version of the Compose Specification, see the [Compose Specification repository](https://github.com/compose-spec/compose-spec).

Use the following links to navigate key sections of the Compose Specification.

{{< grid level="2" >}}
