---
description: Compose file reference
keywords: fig, composition, compose, versions, upgrading, docker
title: Compose file versions and upgrading
---

The Compose file is a [YAML](http://yaml.org/) file defining services,
networks, and volumes for a Docker application.

The Compose file formats are now described in these references, specific to each version.

| **Reference file**       | **What changed in this version** |
|  -------------------     |    ------------------            |
|[Version 3](index.md) (most current, and recommended) |[Version 3 updates](#version-3) |
| [Version 2](compose-file-v2.md)  |[Version 2 updates](#version-2) |
| [Version 1](compose-file-v1.md) | [Version 1 updates](#version-1) |

The topics below explain the differences among the versions, Docker Engine
compatibility, and [how to upgrade](#upgrading).

## Compatibility matrix

{% include content/compose-matrix.md %}

## Versioning

There are currently three versions of the Compose file format:

- Version 1, the legacy format. This is specified by
omitting a `version` key at the root of the YAML.

- Version 2.x. This is specified with a `version: '2'` or `version: '2.1'` entry at the root of the YAML.

- Version 3.x, the latest and recommended version, designed to
be cross-compatible between Compose and the Docker Engine's
[swarm mode](/engine/swarm/index.md). This is specified with a `version: '3'` or `version: '3.1'`, etc., entry at the root of the YAML.


The [Compatibility Matrix](#compatibility-matrix) shows Compose file versions mapped to Docker Engine releases.

To move your project to a later version, see the [Upgrading](#upgrading)
section.

> **Note**: If you're using
> [multiple Compose files](extends.md#different-environments) or
> [extending services](extends.md#extending-services), each file must be of the
> same version - you cannot, for example, mix version 1 and 2 in a single
> project.

Several things differ depending on which version you use:

- The structure and permitted configuration keys
- The minimum Docker Engine version you must be running
- Compose's behaviour with regards to networking

These differences are explained below.

### Version 1

Compose files that do not declare a version are considered "version 1". In those
files, all the [services](index.md#service-configuration-reference) are
declared at the root of the document.

Version 1 is supported by **Compose up to 1.6.x**. It will be deprecated in a
future Compose release.

Version 1 files cannot declare named
[volumes](index.md#volume-configuration-reference), [networks](index.md#network-configuration-reference) or
[build arguments](index.md#args).

Compose does not take advantage of [networking](index.md#networking.md) when you use
version 1: every container is placed on the default `bridge` network and is
reachable from every other container at its IP address. You will need to use
[links](compose-file-v1.md#links) to enable discovery between containers.

Example:

    web:
      build: .
      ports:
       - "5000:5000"
      volumes:
       - .:/code
      links:
       - redis
    redis:
      image: redis

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
[Networking in Compose](compose-file-v2.md#networking.md).

Simple example:

    version: '2'
    services:
      web:
        build: .
        ports:
         - "5000:5000"
        volumes:
         - .:/code
      redis:
        image: redis

A more extended example, defining volumes and networks:

    version: '2'
    services:
      web:
        build: .
        ports:
         - "5000:5000"
        volumes:
         - .:/code
        networks:
          - front-tier
          - back-tier
      redis:
        image: redis
        volumes:
          - redis-data:/var/lib/redis
        networks:
          - back-tier
    volumes:
      redis-data:
        driver: local
    networks:
      front-tier:
        driver: bridge
      back-tier:
        driver: bridge

Several other options were added to support networking, such as:

* [`aliases`](compose-file-v2.md#aliases)

* The [`depends_on`](compose-file-v2.md#dependson) option can be used in place of links to indicate dependencies
between services and startup order.

      version: '2'
      services:
        web:
          build: .
          depends_on:
            - db
            - redis
        redis:
          image: redis
        db:
          image: postgres

* [`ipv4_address`, `ipv6_address`](compose-file-v2.md#ipv4address-ipv6address)

[Variable substitution](compose-file-v2.md#variable-substitution) also was added in Version 2.

### Version 2.1

An upgrade of [version 2](#version-2) that introduces new parameters only
available with Docker Engine version **1.12.0+**

Introduces the following additional parameters:

- [`link_local_ips`](compose-file-v2.md#linklocalips)
- [`isolation`](compose-file-v2.md#isolation)
- `labels` for [volumes](compose-file-v2.md#volume-configuration-reference) and
  [networks](compose-file-v2.md#network-configuration-reference)
- [`userns_mode`](compose-file-v2.md#userns_mode)
- [`healthcheck`](compose-file-v2.md#healthcheck)
- [`sysctls`](compose-file-v2.md#sysctls)

### Version 2.2

An upgrade of [version 2.1](#version-21) that introduces new parameters only
available with Docker Engine version **1.13.0+**. This version also allows
to specify default scale numbers inside the service's configuration.

Introduces the following additional parameters:

- [`init`](compose-file-v2.md#init)
- [`scale`](compose-file-v2.md#scale)


### Version 3

Designed to be cross-compatible between Compose and the Docker Engine's
[swarm mode](/engine/swarm/index.md), version 3 removes several options and adds
several more.

- Removed: `volume_driver`, `volumes_from`, `cpu_shares`, `cpu_quota`, `cpuset`,
  `mem_limit`, `memswap_limit`, `extends`, `group_add`. See the [upgrading](#upgrading)
  guide for how to migrate away from these.

- Added: [deploy](index.md#deploy)

## Upgrading

### Version 2.x to 3.x

Between versions 2.x and 3.x, the structure of the Compose file is the same, but
several options have been removed:

-   `volume_driver`: Instead of setting the volume driver on the service, define
    a volume using the
    [top-level `volumes` option](index.md#volume-configuration-reference)
    and specify the driver there.

        version: "3"
        services:
          db:
            image: postgres
            volumes:
              - data:/var/lib/postgresql/data
        volumes:
          data:
            driver: mydriver

-   `volumes_from`: To share a volume between services, define it using the
    [top-level `volumes` option](index.md#volume-configuration-reference)
    and reference it from each service that shares it using the
    [service-level `volumes` option](index.md#volumes-volumedriver).

-   `cpu_shares`, `cpu_quota`, `cpuset`, `mem_limit`, `memswap_limit`: These
    have been replaced by the [resources](index.md#resources) key under
    `deploy`. Note that `deploy` configuration only takes effect when using
    `docker stack deploy`, and is ignored by `docker-compose`.

-   `extends`: This option has been removed for `version: "3.x"` Compose files.
-   `group_add`: This option has been removed for `version: "3.x"` Compose files.

### Version 1 to 2.x

In the majority of cases, moving from version 1 to 2 is a very simple process:

1. Indent the whole file by one level and put a `services:` key at the top.
2. Add a `version: '2'` line at the top of the file.

It's more complicated if you're using particular configuration features:

-   `dockerfile`: This now lives under the `build` key:

        build:
          context: .
          dockerfile: Dockerfile-alternate

-   `log_driver`, `log_opt`: These now live under the `logging` key:

        logging:
          driver: syslog
          options:
            syslog-address: "tcp://192.168.0.42:123"

-   `links` with environment variables: As documented in the
    [environment variables reference](link-env-deprecated.md), environment variables
    created by
    links have been deprecated for some time. In the new Docker network system,
    they have been removed. You should either connect directly to the
    appropriate hostname or set the relevant environment variable yourself,
    using the link hostname:

        web:
          links:
            - db
          environment:
            - DB_PORT=tcp://db:5432

-   `external_links`: Compose uses Docker networks when running version 2
    projects, so links behave slightly differently. In particular, two
    containers must be connected to at least one network in common in order to
    communicate, even if explicitly linked together.

    Either connect the external container to your app's
    [default network](networking.md), or connect both the external container and
    your service's containers to an
    [external network](networking.md#using-a-pre-existing-network).

-   `net`: This is now replaced by [network_mode](compose-file-v1.md#network_mode):

        net: host    ->  network_mode: host
        net: bridge  ->  network_mode: bridge
        net: none    ->  network_mode: none

    If you're using `net: "container:[service name]"`, you must now use
    `network_mode: "service:[service name]"` instead.

        net: "container:web"  ->  network_mode: "service:web"

    If you're using `net: "container:[container name/id]"`, the value does not
    need to change.

        net: "container:cont-name"  ->  network_mode: "container:cont-name"
        net: "container:abc12345"   ->  network_mode: "container:abc12345"

-   `volumes` with named volumes: these must now be explicitly declared in a
    top-level `volumes` section of your Compose file. If a service mounts a
    named volume called `data`, you must declare a `data` volume in your
    top-level `volumes` section. The whole file might look like this:

        version: '2'
        services:
          db:
            image: postgres
            volumes:
              - data:/var/lib/postgresql/data
        volumes:
          data: {}

    By default, Compose creates a volume whose name is prefixed with your
    project name. If you want it to just be called `data`, declare it as
    external:

        volumes:
          data:
            external: true

## Compose file format references

- [Compose file version 3](index.md)
- [Compose file version 2](compose-file-v2.md)
- [Compose file version 1](compose-file-v1.md)
