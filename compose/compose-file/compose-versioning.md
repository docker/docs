---
description: Compose file reference
keywords: fig, composition, compose, versions, upgrading, docker
title: Compose file versions and upgrading
---

The Compose file is a [YAML](https://yaml.org) file defining services,
networks, and volumes for a Docker application.

## Versioning

> Earlier releases of Docker Compose used a `version` field in the Compose file to define supported attributes. The latest and recommended version of the Compose file format is defined by the [Compose Specification](https://github.com/compose-spec/compose-spec/blob/master/spec.md). This format merges the [legacy 2.x](compose-file-v2.md) and [3.x](compose-file-v3.md) file formats and is implemented by **Compose 1.27.0+** and **Compose V2**.
: .important}

For modern usage, you should **not** declare a `version` attribute in your Compose file. File format is
described by the [Compose Specification](https://github.com/compose-spec/compose-spec/blob/master/spec.md).

## Legacy 

There are three legacy versions of the Compose file format:

- Version 1. (Obsolete)

- Version 2.x. This is specified with a `version: '2'` or `version: '2.1'`, etc., entry at the root of the YAML.

- Version 3.x, designed to be cross-compatible between Compose and the Docker Engine's
[swarm mode](../../engine/swarm/index.md). This is specified with a `version: '3'` or `version: '3.1'`, etc., entry at the root of the YAML.


These differences are explained below.


### Version 2

Compose files using the version 2 syntax must indicate the version number at
the root of the document. All [services](compose-file-v2.md#service-configuration-reference)
must be declared under the `services` key.

Version 2 files are supported by **Compose 1.6.0+** and require a Docker Engine
of version **1.10.0+**.

Named [volumes](compose-file-v2.md#volume-configuration-reference) can be declared under the
`volumes` key, and [networks](compose-file-v2.md#network-configuration-reference) can be declared
under the `networks` key.

By default, every container joins an application-wide default network, and is
discoverable at a hostname that's the same as the service name. This means
[links](compose-file-v2.md#links) are largely unnecessary. For more details, see
[Networking in Compose](../networking.md).


### Version 3

Designed to be cross-compatible between Compose and the Docker Engine's
[swarm mode](/engine/swarm/), version 3 removes several options and adds
several more.

- Removed: `volume_driver`, `volumes_from`, `cpu_shares`, `cpu_quota`,
`cpuset`, `mem_limit`, `memswap_limit`, `extends`, `group_add`. See
the [upgrading](#upgrading) guide for how to migrate away from these.
(For more information on `extends`, see [Extending services](../extends.md#extending-services).)

- Added: [deploy](compose-file-v3.md#deploy)


## Upgrading

The [Compose Specification](https://github.com/compose-spec/compose-spec/blob/master/spec.md) defines all the attributes from both versions 2.x and 3.x, so basically you don't need anything changed for your Compose file to be compliant with the Specification. Just the `version` attribute will be ignored, you can remove it.
