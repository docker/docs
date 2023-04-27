---
description: Key features and use cases of Docker Compose
keywords: documentation, docs, docker, compose, orchestration, containers, uses, features
title: Evolution of Compose
redirect_from:
 - /compose/cli-command-compatibility/
---
{% include compose-eol.md %}

The first release of Compose, written in Python, happened at the end of 2014. 
Between 2014 and 2017 two other noticeable versions of Compose, which introduced new file format versions, were released:

- [Compose 1.6.0 with file format V2](../compose-file/compose-file-v2/)
- [Compose 1.10.0 with file format V3](../compose-file/compose-file-v3/)

These three key file format versions and releases prior to v1.29.2 are collectively referred to as Compose V1. 

In mid-2020 Compose V2 was released. It merged Compose file format V2 and V3 and was written in Go. The file format is defined by the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"}. Compose V2 is the latest and recommended version of Compose and is compatible with Docker Engine version 19.03.0 and later. It provides improved integration with other Docker command-line features, and simplified installation on macOS, Windows, and Linux.  

It makes a clean distinction between the Compose YAML file model and the `docker-compose`
implementation. Making this change has enabled a number of enhancements, including
adding the `compose` command directly into the Docker CLI,  being able to "up" a
Compose application on cloud platforms by simply switching the Docker context,
and launching of [Amazon ECS](../../cloud/ecs-integration.md) and [Microsoft ACI](../../cloud/aci-integration.md).

Compose V2 relies directly on the compose-go bindings which are maintained as part
of the specification. This allows us to include community proposals, experimental
implementations by the Docker CLI and/or Engine, and deliver features faster to
users. 

> **A note about version numbers**
>
>In addition to Compose file format versions described above, the Compose binary itself is on a release schedule, as shown in [Compose releases](https://github.com/docker/compose/releases/). File format versions do not necessarily increment with each release. For example, Compose file format V3 was first introduced in Compose release 1.10.0, and versioned gradually in subsequent releases.
>
>The latest Compose file format, defined by the Compose Specification, was implemented by Docker Compose 1.27.0+.

