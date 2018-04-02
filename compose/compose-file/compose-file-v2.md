---
description: Compose file reference
keywords: fig, composition, compose version 3, docker
redirect_from:
- /compose/yml
title: Compose file version 2 reference
toc_max: 4
toc_min: 1
---

## Reference and guidelines

These topics describe version 2 of the Compose file format.

## Compose and Docker compatibility matrix

There are several versions of the Compose file format – 1, 2, 2.x, and 3.x The
table below is a quick look. For full details on what each version includes and
how to upgrade, see **[About versions and upgrading](compose-versioning.md)**.

{% include content/compose-matrix.md %}

## Service configuration reference

The Compose file is a [YAML](http://yaml.org/) file defining
[services](#service-configuration-reference),
[networks](#network-configuration-reference) and
[volumes](#volume-configuration-reference).
The default path for a Compose file is `./docker-compose.yml`.

>**Tip**: You can use either a `.yml` or `.yaml` extension for this file. They both work.

A [container](/engine/reference/glossary.md#container) definition contains configuration which are applied to each
container started for that service, much like passing command-line parameters to
`docker run`. Likewise, network and volume definitions are analogous to
`docker network create` and `docker volume create`.

As with `docker run`, options specified in the Dockerfile, such as `CMD`,
`EXPOSE`, `VOLUME`, `ENV`, are respected by default - you don't need to
specify them again in `docker-compose.yml`.

You can use environment variables in configuration values with a Bash-like
`${VARIABLE}` syntax - see [variable substitution](#variable-substitution) for
full details.

This section contains a list of all configuration options supported by a service
definition in version 2.

### blkio_config

A set of configuration options to set block IO limits for this service.

    version: '2.2'
    services:
      foo:
        image: busybox
        blkio_config:
          weight: 300
          weight_device:
            - path: /dev/sda
              weight: 400
          device_read_bps:
            - path: /dev/sdb
              rate: '12mb'
          device_read_iops:
            - path: /dev/sdb
              rate: 120
          device_write_bps:
            - path: /dev/sdb
              rate: '1024k'
          device_write_iops:
            - path: /dev/sdb
              rate: 30

#### device_read_bps, device_write_bps

Set a limit in bytes per second for read / write operations on a given device.
Each item in the list must have two keys:

* `path`, defining the symbolic path to the affected device
* `rate`, either as an integer value representing the number of bytes or as
  a string expressing a [byte value](#specifying-byte-values).

#### device_read_iops, device_write_iops

Set a limit in operations per second for read / write operations on a given
device. Each item in the list must have two keys:

* `path`, defining the symbolic path to the affected device
* `rate`, as an integer value representing the permitted number of operations
  per second.

#### weight

Modify the proportion of bandwidth allocated to this service relative to other
services. Takes an integer value between 10 and 1000, with 500 being the
default.

#### weight_device

Fine-tune bandwidth allocation by device. Each item in the list must have
two keys:

* `path`, defining the symbolic path to the affected device
* `weight`, an integer value between 10 and 1000

### build

Configuration options that are applied at build time.

`build` can be specified either as a string containing a path to the build
context, or an object with the path specified under [context](#context) and
optionally [dockerfile](#dockerfile) and [args](#args).

    build: ./dir

    build:
      context: ./dir
      dockerfile: Dockerfile-alternate
      args:
        buildno: 1

If you specify `image` as well as `build`, then Compose names the built image
with the `webapp` and optional `tag` specified in `image`:

    build: ./dir
    image: webapp:tag

This results in an image named `webapp` and tagged `tag`, built from `./dir`.

#### cache_from

> Added in [version 2.2](compose-versioning.md#version-22) file format

A list of images that the engine uses for cache resolution.

    build:
      context: .
      cache_from:
        - alpine:latest
        - corp/web_app:3.14

#### context

> [Version 2 file format](compose-versioning.md#version-2) and up. In version 1, just use
> [build](#build).

Either a path to a directory containing a Dockerfile, or a url to a git repository.

When the value supplied is a relative path, it is interpreted as relative to the
location of the Compose file. This directory is also the build context that is
sent to the Docker daemon.

Compose builds and tags it with a generated name, and use that image thereafter.

    build:
      context: ./dir

#### dockerfile

Alternate Dockerfile.

Compose uses an alternate file to build with. A build path must also be
specified.

    build:
      context: .
      dockerfile: Dockerfile-alternate

#### args

> [Version 2 file format](compose-versioning.md#version-2) and up.

Add build arguments, which are environment variables accessible only during the
build process.

First, specify the arguments in your Dockerfile:

    ARG buildno
    ARG gitcommithash

    RUN echo "Build number: $buildno"
    RUN echo "Based on commit: $gitcommithash"

Then specify the arguments under the `build` key. You can pass a mapping
or a list:

    build:
      context: .
      args:
        buildno: 1
        gitcommithash: cdc3b19

    build:
      context: .
      args:
        - buildno=1
        - gitcommithash=cdc3b19

You can omit the value when specifying a build argument, in which case its value
at build time is the value in the environment where Compose is running.

    args:
      - buildno
      - gitcommithash

> **Note**: YAML boolean values (`true`, `false`, `yes`, `no`, `on`, `off`) must
> be enclosed in quotes, so that the parser interprets them as strings.

#### extra_hosts

Add hostname mappings at build-time. Use the same values as the docker client `--add-host` parameter.

    extra_hosts:
     - "somehost:162.242.195.82"
     - "otherhost:50.31.209.229"

An entry with the ip address and hostname is created in `/etc/hosts` inside containers for this build, e.g:

    162.242.195.82  somehost
    50.31.209.229   otherhost

#### isolation

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Specify a build’s container isolation technology. On Linux, the only supported value
is `default`. On Windows, acceptable values are `default`, `process` and
`hyperv`. Refer to the
[Docker Engine docs](/engine/reference/commandline/run.md#specify-isolation-technology-for-container---isolation)
for details.

If unspecified, Compose will use the `isolation` value found in the service's definition
to determine the value to use for builds.

#### labels

> Added in [version 2.1](compose-versioning.md#version-21) file format

Add metadata to the resulting image using [Docker labels](/engine/userguide/labels-custom-metadata.md).
You can use either an array or a dictionary.

It's recommended that you use reverse-DNS notation to prevent your labels from conflicting with
those used by other software.

    build:
      context: .
      labels:
        com.example.description: "Accounting webapp"
        com.example.department: "Finance"
        com.example.label-with-empty-value: ""


    build:
      context: .
      labels:
        - "com.example.description=Accounting webapp"
        - "com.example.department=Finance"
        - "com.example.label-with-empty-value"

#### network

> Added in [version 2.2](compose-versioning.md#version-22) file format

Set the network containers connect to for the `RUN` instructions during
build.

    build:
      context: .
      network: host


    build:
      context: .
      network: custom_network_1

#### shm_size

> Added in [version 2.3](compose-versioning.md#version-23) file format

Set the size of the `/dev/shm` partition for this build's containers. Specify
as an integer value representing the number of bytes or as a string expressing
a [byte value](#specifying-byte-values).

    build:
      context: .
      shm_size: '2gb'


    build:
      context: .
      shm_size: 10000000

#### target

> Added in [version 2.3](compose-versioning.md#version-23) file format

Build the specified stage as defined inside the `Dockerfile`. See the
[multi-stage build docs](/engine/userguide/eng-image/multistage-build.md) for
details.

      build:
        context: .
        target: prod

### cap_add, cap_drop

Add or drop container capabilities.
See `man 7 capabilities` for a full list.

    cap_add:
      - ALL

    cap_drop:
      - NET_ADMIN
      - SYS_ADMIN

### command

Override the default command.

    command: bundle exec thin -p 3000

The command can also be a list, in a manner similar to
[dockerfile](/engine/reference/builder.md#cmd):

    command: ["bundle", "exec", "thin", "-p", "3000"]

### cgroup_parent

Specify an optional parent cgroup for the container.

    cgroup_parent: m-executor-abcd

### container_name

Specify a custom container name, rather than a generated default name.

    container_name: my-web-container

Because Docker container names must be unique, you cannot scale a service
beyond 1 container if you have specified a custom name. Attempting to do so
results in an error.

### cpu_rt_runtime, cpu_rt_period

> Added in [version 2.2](compose-versioning.md#version-22) file format

Configure CPU allocation parameters using the Docker daemon realtime scheduler.

    cpu_rt_runtime: '400ms'
    cpu_rt_period: '1400us'

    # Integer values will use microseconds as units
    cpu_rt_runtime: 95000
    cpu_rt_period: 11000


### device_cgroup_rules

> [Added in version 2.3 file format](compose-versioning.md#version-23).

Add rules to the cgroup allowed devices list.

    device_cgroup_rules:
      - 'c 1:3 mr'
      - 'a 7:* rmw'

### devices

List of device mappings.  Uses the same format as the `--device` docker
client create option.

    devices:
      - "/dev/ttyUSB0:/dev/ttyUSB0"

### depends_on

> [Version 2 file format](compose-versioning.md#version-2) and up.

Express dependency between services, which has two effects:

- `docker-compose up` starts services in dependency order. In the following
  example, `db` and `redis` are started before `web`.

- `docker-compose up SERVICE` automatically include `SERVICE`'s
  dependencies. In the following example, `docker-compose up web` also
  create and start `db` and `redis`.

Simple example:

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

> **Note**: `depends_on` does not wait for `db` and `redis` to be "ready" before
> starting `web` - only until they have been started. If you need to wait
> for a service to be ready, see [Controlling startup order](/compose/startup-order.md)
> for more on this problem and strategies for solving it.

> [Added in version 2.1 file format](compose-versioning.md#version-21).

A healthcheck indicates that you want a dependency to wait
for another container to be "healthy" (as indicated by a successful state from
the healthcheck) before starting.

Example:

    version: '2.1'
    services:
      web:
        build: .
        depends_on:
          db:
            condition: service_healthy
          redis:
            condition: service_started
      redis:
        image: redis
      db:
        image: redis
        healthcheck:
          test: "exit 0"

In the above example, Compose waits for the `redis` service to be started
(legacy behavior) and the `db` service to be healthy before starting `web`.

See the [healthcheck section](#healthcheck) for complementary
information.

### dns

Custom DNS servers. Can be a single value or a list.

    dns: 8.8.8.8
    dns:
      - 8.8.8.8
      - 9.9.9.9

### dns_opt

List of custom DNS options to be added to the container's `resolv.conf` file.

    dns_opt:
      - use-vc
      - no-tld-query

### dns_search

Custom DNS search domains. Can be a single value or a list.

    dns_search: example.com
    dns_search:
      - dc1.example.com
      - dc2.example.com

### tmpfs

Mount a temporary file system inside the container. Can be a single value or a list.

    tmpfs: /run
    tmpfs:
      - /run
      - /tmp

### entrypoint

Override the default entrypoint.

    entrypoint: /code/entrypoint.sh

The entrypoint can also be a list, in a manner similar to
[dockerfile](/engine/reference/builder.md#entrypoint):

    entrypoint:
        - php
        - -d
        - zend_extension=/usr/local/lib/php/extensions/no-debug-non-zts-20100525/xdebug.so
        - -d
        - memory_limit=-1
        - vendor/bin/phpunit

> **Note**: Setting `entrypoint` both overrides any default entrypoint set
> on the service's image with the `ENTRYPOINT` Dockerfile instruction, *and*
> clears out any default command on the image - meaning that if there's a `CMD`
> instruction in the Dockerfile, it is ignored.

### env_file

Add environment variables from a file. Can be a single value or a list.

If you have specified a Compose file with `docker-compose -f FILE`, paths in
`env_file` are relative to the directory that file is in.

Environment variables declared in the [environment](#environment) section
_override_ these values &ndash; this holds true even if those values are
empty or undefined.

    env_file: .env

    env_file:
      - ./common.env
      - ./apps/web.env
      - /opt/secrets.env

Compose expects each line in an env file to be in `VAR=VAL` format. Lines
beginning with `#` are processed as comments and are ignored. Blank lines are
also ignored.

    # Set Rails/Rack environment
    RACK_ENV=development

> **Note**: If your service specifies a [build](#build) option, variables
> defined in environment files are _not_ automatically visible during the
> build. Use the [args](#args) sub-option of `build` to define build-time
> environment variables.

The value of `VAL` is used as is and not modified at all. For example if the
value is surrounded by quotes (as is often the case of shell variables), the
quotes are included in the value passed to Compose.

Keep in mind that _the order of files in the list is significant in determining
the value assigned to a variable that shows up more than once_. The files in the
list are processed from the top down. For the same variable specified in file
`a.env` and assigned a different value in file `b.env`, if `b.env` is
listed below (after), then the value from `b.env` stands. For example, given the
following declaration in `docker_compose.yml`:

```yaml
services:
  some-service:
    env_file:
      - a.env
      - b.env
```

And the following files:

```none
# a.env
VAR=1
```

and

```none
# b.env
VAR=hello
```

$VAR is `hello`.

### environment

Add environment variables. You can use either an array or a dictionary. Any
boolean values; true, false, yes no, need to be enclosed in quotes to ensure
they are not converted to True or False by the YML parser.

Environment variables with only a key are resolved to their values on the
machine Compose is running on, which can be helpful for secret or host-specific values.

    environment:
      RACK_ENV: development
      SHOW: 'true'
      SESSION_SECRET:

    environment:
      - RACK_ENV=development
      - SHOW=true
      - SESSION_SECRET

> **Note**: If your service specifies a [build](#build) option, variables
> defined in `environment` are _not_ automatically visible during the
> build. Use the [args](#args) sub-option of `build` to define build-time
> environment variables.

### expose

Expose ports without publishing them to the host machine - they'll only be
accessible to linked services. Only the internal port can be specified.

    expose:
     - "3000"
     - "8000"

### extends

Extend another service, in the current file or another, optionally overriding
configuration.

You can use `extends` on any service together with other configuration keys.
The `extends` value must be a dictionary defined with a required `service`
and an optional `file` key.

    extends:
      file: common.yml
      service: webapp

The `service` the name of the service being extended, for example
`web` or `database`. The `file` is the location of a Compose configuration
file defining that service.

If you omit the `file` Compose looks for the service configuration in the
current file. The `file` value can be an absolute or relative path. If you
specify a relative path, Compose treats it as relative to the location of the
current file.

You can extend a service that itself extends another. You can extend
indefinitely. Compose does not support circular references and `docker-compose`
returns an error if it encounters one.

For more on `extends`, see the
[the extends documentation](/compose/extends.md#extending-services).

### external_links

Link to containers started outside this `docker-compose.yml` or even outside
of Compose, especially for containers that provide shared or common services.
`external_links` follow semantics similar to `links` when specifying both the
container name and the link alias (`CONTAINER:ALIAS`).

    external_links:
     - redis_1
     - project_db_1:mysql
     - project_db_1:postgresql

> **Note**: For version 2 file format, the
> externally-created containers must be connected to at least one of the same
> networks as the service which is linking to them.

### extra_hosts

Add hostname mappings. Use the same values as the docker client `--add-host` parameter.

    extra_hosts:
     - "somehost:162.242.195.82"
     - "otherhost:50.31.209.229"

An entry with the ip address and hostname is created in `/etc/hosts` inside containers for this service, e.g:

    162.242.195.82  somehost
    50.31.209.229   otherhost

### group_add

Specify additional groups (by name or number) which the user inside the
container should be a member of. Groups must exist in both the container and the
host system to be added. An example of where this is useful is when multiple
containers (running as different users) need to all read or write the same
file on the host system. That file can be owned by a group shared by all the
containers, and specified in `group_add`. See the
[Docker documentation](/engine/reference/run.md#additional-groups) for more
details.

A full example:

```
version: '2'
services:
  myservice:
    image: alpine
    group_add:
      - mail
```

Running `id` inside the created container shows that the user belongs to
the `mail` group, which would not have been the case if `group_add` were not
used.

### healthcheck

> [Version 2.1 file format](compose-versioning.md#version-21) and up.

Configure a check that's run to determine whether or not containers for this
service are "healthy". See the docs for the
[HEALTHCHECK Dockerfile instruction](/engine/reference/builder.md#healthcheck)
for details on how healthchecks work.

    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost"]
      interval: 1m30s
      timeout: 10s
      retries: 3
      start_period: 40s

`interval`, `timeout` and `start_period` are specified as
[durations](#specifying-durations).

`test` must be either a string or a list. If it's a list, the first item must be
either `NONE`, `CMD` or `CMD-SHELL`. If it's a string, it's equivalent to
specifying `CMD-SHELL` followed by that string.

    # Hit the local web app
    test: ["CMD", "curl", "-f", "http://localhost"]

    # As above, but wrapped in /bin/sh. Both forms below are equivalent.
    test: ["CMD-SHELL", "curl -f http://localhost && echo 'cool, it works'"]
    test: curl -f https://localhost && echo 'cool, it works'

To disable any default healthcheck set by the image, you can use `disable:
true`. This is equivalent to specifying `test: ["NONE"]`.

    healthcheck:
      disable: true

> **Note**: The `start_period` option is a more recent feature and is only
> available with the [2.3 file format](compose-versioning.md#version-23).

### image

Specify the image to start the container from. Can either be a repository/tag or
a partial image ID.

    image: redis
    image: ubuntu:14.04
    image: tutum/influxdb
    image: example-registry.com:4000/postgresql
    image: a4bc65fd

If the image does not exist, Compose attempts to pull it, unless you have also
specified [build](#build), in which case it builds it using the specified
options and tags it with the specified tag.

### init

> [Added in version 2.2 file format](compose-versioning.md#version-22).

Run an init inside the container that forwards signals and reaps processes.
Either set a boolean value to use the default `init`, or specify a path to
a custom one.

    version: '2.2'
    services:
      web:
        image: alpine:latest
        init: true


    version: '2.2'
    services:
      web:
        image: alpine:latest
        init: /usr/libexec/docker-init

### isolation

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Specify a container’s isolation technology. On Linux, the only supported value
is `default`. On Windows, acceptable values are `default`, `process` and
`hyperv`. Refer to the
[Docker Engine docs](/engine/reference/commandline/run.md#specify-isolation-technology-for-container---isolation)
for details.

### labels

Add metadata to containers using [Docker labels](/engine/userguide/labels-custom-metadata.md). You can use either an array or a dictionary.

It's recommended that you use reverse-DNS notation to prevent your labels from conflicting with those used by other software.

    labels:
      com.example.description: "Accounting webapp"
      com.example.department: "Finance"
      com.example.label-with-empty-value: ""

    labels:
      - "com.example.description=Accounting webapp"
      - "com.example.department=Finance"
      - "com.example.label-with-empty-value"

### links

Link to containers in another service. Either specify both the service name and
a link alias (`"SERVICE:ALIAS"`), or just the service name.

> Links are a legacy option. We recommend using
> [networks](#networks) instead.

    web:
      links:
       - "db"
       - "db:database"
       - "redis"

Containers for the linked service are reachable at a hostname identical to
the alias, or the service name if no alias was specified.

Links also express dependency between services in the same way as
[depends_on](#depends_on), so they determine the order of service startup.

> **Note**: If you define both links and [networks](#networks), services with
> links between them must share at least one network in common in order to
> communicate. We recommend using networks instead.

### logging


Logging configuration for the service.

    logging:
      driver: syslog
      options:
        syslog-address: "tcp://192.168.0.42:123"

The `driver`  name specifies a logging driver for the service's
containers, as with the ``--log-driver`` option for docker run
([documented here](/engine/admin/logging/overview.md)).

The default value is json-file.

    driver: "json-file"
    driver: "syslog"
    driver: "none"

> **Note**: Only the `json-file` and `journald` drivers make the logs available directly from
> `docker-compose up` and `docker-compose logs`. Using any other driver does not
> print any logs.

Specify logging options for the logging driver with the ``options`` key, as with the ``--log-opt`` option for `docker run`.

Logging options are key-value pairs. An example of `syslog` options:

    driver: "syslog"
    options:
      syslog-address: "tcp://192.168.0.42:123"

### network_mode

> [Version 2 file format](compose-versioning.md#version-2) and up. Replaces the version 1 [net](compose-file-v1.md#net) option.

Network mode. Use the same values as the docker client `--net` parameter, plus
the special form `service:[service name]`.

    network_mode: "bridge"
    network_mode: "host"
    network_mode: "none"
    network_mode: "service:[service name]"
    network_mode: "container:[container name/id]"

### networks

> [Version 2 file format](compose-versioning.md#version-2) and up. Replaces the version 1 [net](compose-file-v1.md#net) option.

Networks to join, referencing entries under the
[top-level `networks` key](#network-configuration-reference).

    services:
      some-service:
        networks:
         - some-network
         - other-network

#### aliases

Aliases (alternative hostnames) for this service on the network. Other containers on the same network can use either the service name or this alias to connect to one of the service's containers.

Since `aliases` is network-scoped, the same service can have different aliases on different networks.

> **Note**: A network-wide alias can be shared by multiple containers, and even by multiple services. If it is, then exactly which container the name resolves to is not guaranteed.

The general format is shown here.

    services:
      some-service:
        networks:
          some-network:
            aliases:
             - alias1
             - alias3
          other-network:
            aliases:
             - alias2

In the example below, three services are provided (`web`, `worker`, and `db`), along with two networks (`new` and `legacy`). The `db` service is reachable at the hostname `db` or `database` on the `new` network, and at `db` or `mysql` on the `legacy` network.

    version: '2'

    services:
      web:
        build: ./web
        networks:
          - new

      worker:
        build: ./worker
        networks:
          - legacy

      db:
        image: mysql
        networks:
          new:
            aliases:
              - database
          legacy:
            aliases:
              - mysql

    networks:
      new:
      legacy:

#### ipv4_address, ipv6_address

Specify a static IP address for containers for this service when joining the network.

The corresponding network configuration in the [top-level networks section](#network-configuration-reference) must have an `ipam` block with subnet and gateway configurations covering each static address. If IPv6 addressing is desired, the [`enable_ipv6`](#enableipv6) option must be set.

An example:

    version: '2.1'

    services:
      app:
        image: busybox
        command: ifconfig
        networks:
          app_net:
            ipv4_address: 172.16.238.10
            ipv6_address: 2001:3984:3989::10

    networks:
      app_net:
        driver: bridge
        enable_ipv6: true
        ipam:
          driver: default
          config:
          - subnet: 172.16.238.0/24
            gateway: 172.16.238.1
          - subnet: 2001:3984:3989::/64
            gateway: 2001:3984:3989::1

#### link_local_ips

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Specify a list of link-local IPs. Link-local IPs are special IPs which belong
to a well known subnet and are purely managed by the operator, usually
dependent on the architecture where they are deployed. Therefore they are not
managed by docker (IPAM driver).

Example usage:

    version: '2.1'
    services:
      app:
        image: busybox
        command: top
        networks:
          app_net:
            link_local_ips:
              - 57.123.22.11
              - 57.123.22.13
    networks:
      app_net:
        driver: bridge

#### priority

Specify a priority to indicate in which order Compose should connect the
service's containers to its networks. If unspecified, the default value is `0`.

In the following example, the `app` service connects to `app_net_1` first
as it has the highest priority. It then connects to `app_net_3`, then
`app_net_2`, which uses the default priority value of `0`.

    version: '2.3'
    services:
      app:
        image: busybox
        command: top
        networks:
          app_net_1:
            priority: 1000
          app_net_2:

          app_net_3:
            priority: 100
    networks:
      app_net_1:
      app_net_2:
      app_net_3:

> **Note:** If multiple networks have the same priority, the connection order
> is undefined.

### pid

    pid: "host"
    pid: "container:custom_container_1"
    pid: "service:foobar"

If set to one of the following forms: `container:<container_name>`,
`service:<service_name>`, the service shares the PID address space of the
designated container or service.

If set to "host", the service's PID mode is the host PID mode.  This turns
on sharing between container and the host operating system the PID address
space. Containers launched with this flag can access and manipulate
other containers in the bare-metal machine's namespace and vise-versa.

> **Note**: the `service:` and `container:` forms require
> [version 2.1](compose-versioning.md#version-21) or above

### pids_limit

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Tunes a container's PIDs limit. Set to `-1` for unlimited PIDs.

    pids_limit: 10


### platform

> [Added in version 2.4 file format](compose-versioning.md#version-24).

Target platform containers for this service will run on, using the
`os[/arch[/variant]]` syntax, e.g.

    platform: osx
    platform: windows/amd64
    platform: linux/arm64/v8

This parameter determines which version of the image will be pulled and/or
on which platform the service's build will be performed.

### ports

Expose ports. Either specify both ports (`HOST:CONTAINER`), or just the container
port (an ephemeral host port is chosen).

> **Note**: When mapping ports in the `HOST:CONTAINER` format, you may experience
> erroneous results when using a container port lower than 60, because YAML
> parses numbers in the format `xx:yy` as a base-60 value. For this reason,
> we recommend always explicitly specifying your port mappings as strings.

    ports:
     - "3000"
     - "3000-3005"
     - "8000:8000"
     - "9090-9091:8080-8081"
     - "49100:22"
     - "127.0.0.1:8001:8001"
     - "127.0.0.1:5000-5010:5000-5010"
     - "6060:6060/udp"
     - "12400-12500:1240"

### runtime

> [Added in version 2.3 file format](compose-versioning.md#version-23)

Specify which runtime to use for the service's containers. Default runtime
and available runtimes are listed in the output of `docker info`.

    web:
      image: busybox:latest
      command: true
      runtime: runc


### scale

> [Added in version 2.2 file format](compose-versioning.md#version-22)

Specify the default number of containers to deploy for this service. Whenever
you run `docker-compose up`, Compose creates or removes containers to match
the specified number. This value can be overridden using the
[`--scale`](/compose/reference/up.md) flag.

    web:
      image: busybox:latest
      command: echo 'scaled'
      scale: 3

### security_opt

Override the default labeling scheme for each container.

    security_opt:
      - label:user:USER
      - label:role:ROLE

### stop_grace_period

Specify how long to wait when attempting to stop a container if it doesn't
handle SIGTERM (or whatever stop signal has been specified with
[`stop_signal`](#stopsignal)), before sending SIGKILL. Specified
as a [duration](#specifying-durations).

    stop_grace_period: 1s
    stop_grace_period: 1m30s

By default, `stop` waits 10 seconds for the container to exit before sending
SIGKILL.

### stop_signal

Sets an alternative signal to stop the container. By default `stop` uses
SIGTERM. Setting an alternative signal using `stop_signal` causes
`stop` to send that signal instead.

    stop_signal: SIGUSR1

### storage_opt

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Set storage driver options for this service.

    storage_opt:
      size: '1G'

### sysctls

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Kernel parameters to set in the container. You can use either an array or a
dictionary.

    sysctls:
      net.core.somaxconn: 1024
      net.ipv4.tcp_syncookies: 0

    sysctls:
      - net.core.somaxconn=1024
      - net.ipv4.tcp_syncookies=0

### ulimits

Override the default ulimits for a container. You can either specify a single
limit as an integer or soft/hard limits as a mapping.


    ulimits:
      nproc: 65535
      nofile:
        soft: 20000
        hard: 40000

### userns_mode

> [Added in version 2.1 file format](compose-versioning.md#version-21).

    userns_mode: "host"

Disables the user namespace for this service, if Docker daemon is configured with user namespaces.
See [dockerd](/engine/reference/commandline/dockerd.md#disable-user-namespace-for-a-container) for
more information.

### volumes

Mount host folders or named volumes. Named volumes need to be specified with the
[top-level `volumes` key](#volume-configuration-reference).

You can mount a relative path on the host, which expands relative to
the directory of the Compose configuration file being used. Relative paths
should always begin with `.` or `..`.

#### Short syntax

The short syntax uses the generic `[SOURCE:]TARGET[:MODE]` format, where
`SOURCE` can be either a host path or volume name. `TARGET` is the container
path where the volume is mounted. Standard modes are `ro` for read-only
and `rw` for read-write (default).

    volumes:
      # Just specify a path and let the Engine create a volume
      - /var/lib/mysql

      # Specify an absolute path mapping
      - /opt/data:/var/lib/mysql

      # Path on the host, relative to the Compose file
      - ./cache:/tmp/cache

      # User-relative path
      - ~/configs:/etc/configs/:ro

      # Named volume
      - datavolume:/var/lib/mysql

#### Long syntax

> [Added in version 2.3 file format](compose-versioning.md#version-23).

The long form syntax allows the configuration of additional fields that can't be
expressed in the short form.

- `type`: the mount type `volume`, `bind`, `tmpfs` or `npipe`
- `source`: the source of the mount, a path on the host for a bind mount, or the
  name of a volume defined in the
  [top-level `volumes` key](#volume-configuration-reference). Not applicable for a tmpfs mount.
- `target`: the path in the container where the volume is mounted
- `read_only`: flag to set the volume as read-only
- `bind`: configure additional bind options
  - `propagation`: the propagation mode used for the bind
- `volume`: configure additional volume options
  - `nocopy`: flag to disable copying of data from a container when a volume is
    created
- `tmpfs`: configure additional tmpfs options
  - `size`: the size for the tmpfs mount in bytes


```none
version: "2.3"
services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - type: volume
        source: mydata
        target: /data
        volume:
          nocopy: true
      - type: bind
        source: ./static
        target: /opt/app/static

networks:
  webnet:

volumes:
  mydata:
```

> **Note:** When creating bind mounts, using the long syntax requires the
> referenced folder to be created beforehand. Using the short syntax
> creates the folder on the fly if it doesn't exist.
> See the [bind mounts documentation](/engine/admin/volumes/bind-mounts.md/#differences-between--v-and---mount-behavior)
> for more information.

### volume\_driver

Specify a default volume driver to be used for all declared volumes on this
service.

    volume_driver: mydriver

> **Note:** In [version 2 files](compose-versioning.md#version-2), this
> option only applies to anonymous volumes (those specified in the image,
> or specified under `volumes` without an explicit named volume or host path).
> To configure the driver for a named volume, use the `driver` key under the
> entry in the [top-level `volumes` option](#volume-configuration-reference).


See [Docker Volumes](/engine/userguide/dockervolumes.md) and
[Volume Plugins](/engine/extend/plugins_volume.md) for more information.

### volumes_from

Mount all of the volumes from another service or container, optionally
specifying read-only access (``ro``) or read-write (``rw``). If no access level is specified,
then read-write is used.

    volumes_from:
     - service_name
     - service_name:ro
     - container:container_name
     - container:container_name:rw

> **Notes**
>
>* The `container:...` formats are only supported in the
> [version 2 file format](compose-versioning.md#version-2).
>
>* In [version 1](compose-versioning.md#version-1), you can use
> container names without marking them as such:
>
>     - `service_name`
>     - `service_name:ro`
>     - `container_name`
>     - `container_name:rw`

### restart

`no` is the default restart policy, and it doesn't restart a container under any circumstance. When `always` is specified, the container always restarts. The `on-failure` policy restarts a container if the exit code indicates an on-failure error.

      - restart: no
      - restart: always
      - restart: on-failure

{: id="cpu-and-other-resources"}

### cpu_count, cpu_percent, cpu\_shares, cpu\_period, cpu\_quota, cpus, cpuset, domainname, hostname, ipc, mac\_address, mem\_limit, memswap\_limit, mem\_swappiness, mem\_reservation, oom_kill_disable, oom_score_adj, privileged, read\_only, shm\_size, stdin\_open, tty, user, working\_dir

Each of these is a single value, analogous to its
[docker run](/engine/reference/run.md) counterpart.

> **Note:** The following options were added in [version 2.2](compose-versioning.md#version-22):
> `cpu_count`, `cpu_percent`, `cpus`.
> The following options were added in [version 2.1](compose-versioning.md#version-21):
> `oom_kill_disable`, `cpu_period`

    cpu_count: 2
    cpu_percent: 50
    cpus: 0.5
    cpu_shares: 73
    cpu_quota: 50000
    cpu_period: 20ms
    cpuset: 0,1

    user: postgresql
    working_dir: /code

    domainname: foo.com
    hostname: foo
    ipc: host
    mac_address: 02:42:ac:11:65:43

    mem_limit: 1000000000
    memswap_limit: 2000000000
    mem_reservation: 512m
    privileged: true

    oom_score_adj: 500
    oom_kill_disable: true

    read_only: true
    shm_size: 64M
    stdin_open: true
    tty: true

{: id="orig-resources" }

## Specifying durations

Some configuration options, such as the `interval` and `timeout` sub-options for
[`healthcheck`](#healthcheck), accept a duration as a string in a
format that looks like this:

    2.5s
    10s
    1m30s
    2h32m
    5h34m56s

The supported units are `us`, `ms`, `s`, `m` and `h`.

## Specifying byte values

Some configuration options, such as the `device_read_bps` sub-option for
[`blkio_config`](#blkioconfig), accept a byte value as a string in a format
that looks like this:

    2b
    1024kb
    2048k
    300m
    1gb

The supported units are `b`, `k`, `m` and `g`, and their alternative notation `kb`,
`mb` and `gb`. Decimal values are not supported at this time.

## Volume configuration reference

While it is possible to declare volumes on the fly as part of the service
declaration, this section allows you to create named volumes that can be
reused across multiple services (without relying on `volumes_from`), and are
easily retrieved and inspected using the docker command line or API.
See the [docker volume](/engine/reference/commandline/volume_create.md)
subcommand documentation for more information.

Here's an example of a two-service setup where a database's data directory is
shared with another service as a volume so that it can be periodically backed
up:

    version: "2.2"

    services:
      db:
        image: db
        volumes:
          - data-volume:/var/lib/db
      backup:
        image: backup-service
        volumes:
          - data-volume:/var/lib/backup/data

    volumes:
      data-volume:

An entry under the top-level `volumes` key can be empty, in which case it
uses the default driver configured by the Engine (in most cases, this is the
`local` driver). Optionally, you can configure it with the following keys:

### driver

Specify which volume driver should be used for this volume. Defaults to whatever
driver the Docker Engine has been configured to use, which in most cases is
`local`. If the driver is not available, the Engine returns an error when
`docker-compose up` tries to create the volume.

     driver: foobar

### driver_opts

Specify a list of options as key-value pairs to pass to the driver for this
volume. Those options are driver-dependent - consult the driver's
documentation for more information. Optional.

     driver_opts:
       foo: "bar"
       baz: 1

### external

If set to `true`, specifies that this volume has been created outside of
Compose. `docker-compose up` does not attempt to create it, and raises
an error if it doesn't exist.

`external` cannot be used in conjunction with other volume configuration keys
(`driver`, `driver_opts`).

In the example below, instead of attempting to create a volume called
`[projectname]_data`, Compose looks for an existing volume simply
called `data` and mount it into the `db` service's containers.

    version: '2'

    services:
      db:
        image: postgres
        volumes:
          - data:/var/lib/postgresql/data

    volumes:
      data:
        external: true

You can also specify the name of the volume separately from the name used to
refer to it within the Compose file:

    volumes:
      data:
        external:
          name: actual-name-of-volume

> **Note**: In newer versions of Compose, the `external.name` property is
> deprecated in favor of simply using the `name` property.

### labels

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Add metadata to containers using
[Docker labels](/engine/userguide/labels-custom-metadata.md). You can use either
an array or a dictionary.

It's recommended that you use reverse-DNS notation to prevent your labels from
conflicting with those used by other software.

    labels:
      com.example.description: "Database volume"
      com.example.department: "IT/Ops"
      com.example.label-with-empty-value: ""

    labels:
      - "com.example.description=Database volume"
      - "com.example.department=IT/Ops"
      - "com.example.label-with-empty-value"


### name

> [Added in version 2.1 file format](compose-versioning.md#version-21)

Set a custom name for this volume.

    version: '2.1'
    volumes:
      data:
        name: my-app-data

It can also be used in conjuction with the `external` property:

    version: '2.1'
    volumes:
      data:
        external: true
        name: my-app-data

## Network configuration reference

The top-level `networks` key lets you specify networks to be created. For a full
explanation of Compose's use of Docker networking features, see the
[Networking guide](/compose/networking.md).

### driver

Specify which driver should be used for this network.

The default driver depends on how the Docker Engine you're using is configured,
but in most instances it is `bridge` on a single host and `overlay` on a
Swarm.

The Docker Engine returns an error if the driver is not available.

    driver: overlay

Starting in Compose file format 2.1, overlay networks are always created as
`attachable`, and this is not configurable. This means that standalone
containers can connect to overlay networks.

### driver_opts

Specify a list of options as key-value pairs to pass to the driver for this
network. Those options are driver-dependent - consult the driver's
documentation for more information. Optional.

      driver_opts:
        foo: "bar"
        baz: 1

### enable_ipv6

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Enable IPv6 networking on this network.

### ipam

Specify custom IPAM config. This is an object with several properties, each of
which is optional:

-   `driver`: Custom IPAM driver, instead of the default.
-   `config`: A list with zero or more config blocks, each containing any of
    the following keys:
    - `subnet`: Subnet in CIDR format that represents a network segment
    - `ip_range`: Range of IPs from which to allocate container IPs
    - `gateway`: IPv4 or IPv6 gateway for the master subnet
    - `aux_addresses`: Auxiliary IPv4 or IPv6 addresses used by Network driver,
      as a mapping from hostname to IP
-   `options`: Driver-specific options as a key-value mapping.

A full example:

    ipam:
      driver: default
      config:
        - subnet: 172.28.0.0/16
          ip_range: 172.28.5.0/24
          gateway: 172.28.5.254
          aux_addresses:
            host1: 172.28.1.5
            host2: 172.28.1.6
            host3: 172.28.1.7
      options:
        foo: bar
        baz: "0"

### internal

By default, Docker also connects a bridge network to it to provide external
connectivity. If you want to create an externally isolated overlay network,
you can set this option to `true`.

### labels

> [Added in version 2.1 file format](compose-versioning.md#version-21).

Add metadata to containers using
[Docker labels](/engine/userguide/labels-custom-metadata.md). You can use either
an array or a dictionary.

It's recommended that you use reverse-DNS notation to prevent your labels from
conflicting with those used by other software.

    labels:
      com.example.description: "Financial transaction network"
      com.example.department: "Finance"
      com.example.label-with-empty-value: ""

    labels:
      - "com.example.description=Financial transaction network"
      - "com.example.department=Finance"
      - "com.example.label-with-empty-value"

### external

If set to `true`, specifies that this network has been created outside of
Compose. `docker-compose up` does not attempt to create it, and raises
an error if it doesn't exist.

`external` cannot be used in conjunction with other network configuration keys
(`driver`, `driver_opts`, `group_add`, `ipam`, `internal`).

In the example below, `proxy` is the gateway to the outside world. Instead of
attempting to create a network called `[projectname]_outside`, Compose
looks for an existing network simply called `outside` and connect the `proxy`
service's containers to it.

    version: '2'

    services:
      proxy:
        build: ./proxy
        networks:
          - outside
          - default
      app:
        build: ./app
        networks:
          - default

    networks:
      outside:
        external: true

You can also specify the name of the network separately from the name used to
refer to it within the Compose file:

    networks:
      outside:
        external:
          name: actual-name-of-network


Not supported for version 2 `docker-compose` files. Use
[network_mode](#network_mode) instead.

### name

> [Added in version 2.1 file format](compose-versioning.md#version-21)

Set a custom name for this network.

    version: '2.1'
    networks:
      network1:
        name: my-app-net

It can also be used in conjuction with the `external` property:

    version: '2.1'
    networks:
      network1:
        external: true
        name: my-app-net

## Variable substitution

{% include content/compose-var-sub.md %}

## Extension fields

> [Added in version 2.1 file format](compose-versioning.md#version-21).

{% include content/compose-extfields-sub.md %}

## Compose documentation

- [User guide](/compose/index.md)
- [Installing Compose](/compose/install.md)
- [Compose file versions and upgrading](compose-versioning.md)
- [Get started with Django](/compose/django.md)
- [Get started with Rails](/compose/rails.md)
- [Get started with WordPress](/compose/wordpress.md)
- [Command line reference](/compose/reference/)
