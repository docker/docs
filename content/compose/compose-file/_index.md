---
description: Find the latest recommended version of the Docker Compose file format
  for defining multi-container applications.
keywords: docker compose file, docker compose yml, docker compose reference, docker
  compose cmd, docker compose user, docker compose image, yaml spec, docker compose
  syntax, yaml specification, docker compose specification
title: Overview
toc_max: 4
toc_min: 1
grid:
- title: Version and name top-level element
  description: Understand version and name attributes for Compose.
  icon: feed
  link: /compose/compose-file/04-version-and-name/
- title: Services top-level element
  description: Explore all services attributes for Compose.
  icon: construction
  link: /compose/compose-file/05-services/
- title: Networks top-level element
  description: Find all networks attributes for Compose.
  icon: lan
  link: /compose/compose-file/06-networks/
- title: Volumes top-level element
  description: Explore all volumes attributes for Compose.
  icon: database
  link: /compose/compose-file/07-volumes/
- title: Configs top-level element
  description: Find out about configs in Compose.
  icon: settings_suggest
  link: /compose/compose-file/08-configs/
- title: Secrets top-level element
  description: Learn about secrets in Compose.
  icon: lock
  link: /compose/compose-file/09-secrets/
aliases:
- /compose/yaml/
- /compose/compose-file/compose-file-v1/
---

>**New to Compose?**
>
> Find more information about the [key features and use cases of Docker Compose](../intro/features-uses.md) or [try the get started guide](../gettingstarted.md).
{ .tip }

The Compose Specification is the latest and recommended version of the Compose file format. It helps you define a [Compose file](../compose-application-model.md) which is used to configure your Docker application’s services, networks, volumes, and more.

Legacy versions 2.x and 3.x of the Compose file format were merged into the Compose Specification. It is implemented in versions 1.27.0 and above (also known as Compose V2) of the Docker Compose CLI.

The Compose Specification on Docker Docs is the Docker Compose implementation. If you wish to implement your own version of the Compose Specification, see the [Compose Specication repository](https://github.com/compose-spec/compose-spec).

> **Note**
>
> Compose V1 no longer receives updates and is not available in new releases of Docker Desktop.
>
> Compose V2 is included with all currently supported versions of Docker Desktop.
> For more information, see [Migrate to Compose V2](/compose/migrate).

Use the following links to navigate key sections of the Compose Specification. 

{{< grid >}}