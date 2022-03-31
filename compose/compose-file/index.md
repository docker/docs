---
description: Compose file reference
keywords: fig, composition, compose, docker
redirect_from:
- /compose/yaml
- /compose/compose-file/compose-file-v1/
title: Compose specification
toc_max: 4
toc_min: 1
---

The Compose file is a [YAML](https://yaml.org){: target="_blank" rel="noopener" class="_"} file defining services,
networks, and volumes for a Docker application. The latest and recommended
version of the Compose file format is defined by the [Compose
Specification](https://github.com/compose-spec/compose-spec/blob/master/spec.md){:
target="_blank" rel="noopener" class="_"}. The Compose spec merges the legacy
2.x and 3.x versions, aggregating properties across these formats and is implemented by **Compose 1.27.0+**.

## Status of this document

This document specifies the Compose file format used to define multi-containers applications. Distribution of this document is unlimited.

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119){: target="_blank" rel="noopener" class="_"}.

### Requirements and optional attributes

The Compose specification includes properties designed to target a local [OCI](https://opencontainers.org/){: target="_blank" rel="noopener" class="_"} container runtime,
exposing Linux kernel specific configuration options, but also some Windows container specific properties, as well as cloud platform features related to resource placement on a cluster, replicated application distribution and scalability.

We acknowledge that no Compose implementation is expected to support **all** attributes, and that support for some properties
is Platform dependent and can only be confirmed at runtime. The definition of a versioned schema to control the supported
properties in a Compose file, established by the [docker-compose](https://github.com/docker/compose){: target="_blank" rel="noopener" class="_"} tool where the Compose
file format was designed, doesn't offer any guarantee to the end-user attributes will be actually implemented.

The specification defines the expected configuration syntax and behavior, but - until noted - supporting any of those is OPTIONAL.

A Compose implementation to parse a Compose file using unsupported attributes SHOULD warn user. We recommend implementors
to support those running modes:

* default: warn user about unsupported attributes, but ignore them
* strict: warn user about unsupported attributes and reject the compose file
* loose: ignore unsupported attributes AND unknown attributes (that were not defined by the spec by the time implementation was created)

## The Compose application model

The Compose specification allows one to define a platform-agnostic container based application. Such an application is designed as a set of containers which have to both run together with adequate shared resources and communication channels.

Computing components of an application are defined as [Services](#services-top-level-element). A Service is an abstract concept implemented on platforms by running the same container image (and configuration) one or more times.

Services communicate with each other through [Networks](#networks-top-level-element). In this specification, a Network is a platform capability abstraction to establish an IP route between containers within services connected together. Low-level, platform-specific networking options are grouped into the Network definition and MAY be partially implemented on some platforms.

Services store and share persistent data into [Volumes](#volumes-top-level-element). The specification describes such a persistent data as a high-level filesystem mount with global options. Actual platform-specific implementation details are grouped into the Volumes definition and MAY be partially implemented on some platforms.

Some services require configuration data that is dependent on the runtime or platform. For this, the specification defines a dedicated concept: [Configs](#configs-top-level-element). From a Service container point of view, Configs are comparable to Volumes, in that they are files mounted into the container. But the actual definition involves distinct platform resources and services, which are abstracted by this type.

A [Secret](#secrets-top-level-element) is a specific flavor of configuration data for sensitive data that SHOULD NOT be exposed without security considerations. Secrets are made available to services as files mounted into their containers, but the platform-specific resources to provide sensitive data are specific enough to deserve a distinct concept and definition within the Compose specification.

Distinction within Volumes, Configs and Secret allows implementations to offer a comparable abstraction at service level, but cover the specific configuration of adequate platform resources for well identified data usages.

A **Project** is an individual deployment of an application specification on a platform. A project's name is used to group
resources together and isolate them from other applications or other installation of the same Compose specified application with distinct parameters. A Compose implementation creating resources on a platform MUST prefix resource names by project and
set the label `com.docker.compose.project`.

Project name can be set explicitly by top-level `name` attribute. Compose implementation MUST offer a way for user to set a custom project name and override this name, so that the same `compose.yaml` file can be deployed twice on the same infrastructure, without changes, by just passing a distinct name.

### Illustrative example

The following example illustrates Compose specification concepts with a concrete example application. The example is non-normative.

Consider an application split into a frontend web application and a backend service.

The frontend is configured at runtime with an HTTP configuration file managed by infrastructure, providing an external domain name, and an HTTPS server certificate injected by the platform's secured secret store.

The backend stores data in a persistent volume.

Both services communicate with each other on an isolated back-tier network, while frontend is also connected to a front-tier network and exposes port 443 for external usage.

```
(External user) --> 443 [frontend network]
                            |
                  +--------------------+
                  |  frontend service  |...ro...<HTTP configuration>
                  |      "webapp"      |...ro...<server certificate> #secured
                  +--------------------+
                            |
                        [backend network]
                            |
                  +--------------------+
                  |  backend service   |  r+w   ___________________
                  |     "database"     |=======( persistent volume )
                  +--------------------+        \_________________/
```

The example application is composed of the following parts:

- 2 services, backed by Docker images: `webapp` and `database`
- 1 secret (HTTPS certificate), injected into the frontend
- 1 configuration (HTTP), injected into the frontend
- 1 persistent volume, attached to the backend
- 2 networks

```yml
services:
  frontend:
    image: awesome/webapp
    ports:
      - "443:8043"
    networks:
      - front-tier
      - back-tier
    configs:
      - httpd-config
    secrets:
      - server-certificate

  backend:
    image: awesome/database
    volumes:
      - db-data:/etc/data
    networks:
      - back-tier

volumes:
  db-data:
    driver: flocker
    driver_opts:
      size: "10GiB"

configs:
  httpd-config:
    external: true

secrets:
  server-certificate:
    external: true

networks:
  # The presence of these objects is sufficient to define them
  front-tier: {}
  back-tier: {}
```

This example illustrates the distinction between volumes, configs and secrets. While all of them are all exposed
to service containers as mounted files or directories, only a volume can be configured for read+write access.
Secrets and configs are read-only. The volume configuration allows you to select a volume driver and pass driver options
to tweak volume management according to the actual infrastructure. Configs and Secrets rely on platform services,
and are declared `external` as they are not managed as part of the application lifecycle: the Compose implementation
will use a platform-specific lookup mechanism to retrieve runtime values.

## Compose file

The Compose file is a [YAML](http://yaml.org/) file defining
[version](#version-top-level-element) (DEPRECATED),
[services](#services-top-level-element) (REQUIRED),
[networks](#networks-top-level-element),
[volumes](#volumes-top-level-element),
[configs](#configs-top-level-element) and
[secrets](#secrets-top-level-element).
The default path for a Compose file is `compose.yaml` (preferred) or `compose.yml` in working directory.
Compose implementations SHOULD also support `docker-compose.yaml` and `docker-compose.yml` for backward compatibility.
If both files exist, Compose implementations MUST prefer canonical `compose.yaml` one.

Multiple Compose files can be combined together to define the application model. The combination of YAML files
MUST be implemented by appending/overriding YAML elements based on Compose file order set by the user. Simple
attributes and maps get overridden by the highest order Compose file, lists get merged by appending. Relative
paths MUST be resolved based on the **first** Compose file's parent folder, whenever complimentary files being
merged are hosted in other folders.

As some Compose file elements can both be expressed as single strings or complex objects, merges MUST apply to
the expanded form.

### Profiles

Profiles allow to adjust the Compose application model for various usages and environments. A Compose
implementation SHOULD allow the user to define a set of active profiles. The exact mechanism is implementation
specific and MAY include command line flags, environment variables, etc.

The Services top-level element supports a `profiles` attribute to define a list of named profiles. Services without
a `profiles` attribute set MUST always be enabled. A service MUST be ignored by the Compose
implementation when none of the listed `profiles` match the active ones, unless the service is
explicitly targeted by a command. In that case its `profiles` MUST be added to the set of active profiles.
All other top-level elements are not affected by `profiles` and are always active.

References to other services (by `links`, `extends` or shared resource syntax `service:xxx`) MUST not
automatically enable a component that would otherwise have been ignored by active profiles. Instead the
Compose implementation MUST return an error.

#### Illustrative example

```yaml
services:
  foo:
    image: foo
  bar:
    image: bar
    profiles:
      - test
  baz:
    image: baz
    depends_on:
      - bar
    profiles:
      - test
  zot:
    image: zot
    depends_on:
      - bar
    profiles:
      - debug
```

- Compose application model parsed with no profile enabled only contains the `foo` service.
- If profile `test` is enabled, model contains the services `bar` and `baz` which are enabled by the
  `test` profile and service `foo` which is always enabled.
- If profile `debug` is enabled, model contains both `foo` and `zot` services, but not `bar` and `baz`
  and as such the model is invalid regarding the `depends_on` constraint of `zot`.
- If profiles `debug` and `test` are enabled, model contains all services: `foo`, `bar`, `baz` and `zot`.
- If Compose implementation is executed with `bar` as explicit service to run, it and the `test` profile
  will be active even if `test` profile is not enabled _by the user_.
- If Compose implementation is executed with `baz` as explicit service to run, the service `baz` and the
  profile `test` will be active and `bar` will be pulled in by the `depends_on` constraint.
- If Compose implementation is executed with `zot` as explicit service to run, again the model will be
  invalid regarding the `depends_on` constraint of `zot` since `zot` and `bar` have no common `profiles`
  listed.
- If Compose implementation is executed with `zot` as explicit service to run and profile `test` enabled,
  profile `debug` is automatically enabled and service `bar` is pulled in as a dependency starting both
  services `zot` and `bar`.

## Version top-level element

Top-level `version` property is defined by the specification for backward compatibility but is only informative.

A Compose implementation SHOULD NOT use this version to select an exact schema to validate the Compose file, but
prefer the most recent schema at the time it has been designed.

Compose implementations SHOULD validate whether they can fully parse the Compose file. If some fields are unknown, typically
because the Compose file was written with fields defined by a newer version of the specification, Compose implementations
SHOULD warn the user. Compose implementations MAY offer options to ignore unknown fields (as defined by ["loose"](#requirements-and-optional-attributes) mode).

## Name top-level element

Top-level `name` property is defined by the specification as project name to be used if user doesn't set one explicitly. 
Compose implementations MUST offer a way for user to override this name, and SHOULD define a mechanism to compute a
default project name, to be used if the top-level `name` element is not set.

Whenever project name is defined by top-level `name` or by some custom mechanism, it MUST be exposed for 
[interpolation](#interpolation) and environment variable resolution as `COMPOSE_PROJECT_NAME`

```yml
services:
  foo:
    image: busybox
    environment:
      - COMPOSE_PROJECT_NAME
    command: echo "I'm running ${COMPOSE_PROJECT_NAME}"
```

## Services top-level element

A Service is an abstract definition of a computing resource within an application which can be scaled/replaced
independently from other components. Services are backed by a set of containers, run by the platform
according to replication requirements and placement constraints. Being backed by containers, Services are defined
by a Docker image and set of runtime arguments. All containers within a service are identically created with these
arguments.

A Compose file MUST declare a `services` root element as a map whose keys are string representations of service names,
and whose values are service definitions. A service  definition contains the configuration that is applied to each
container started for that service.

Each service MAY also include a Build section, which defines how to create the Docker image for the service.
Compose implementations MAY support building docker images using this service definition. If not implemented
the Build section SHOULD be ignored and the Compose file MUST still be considered valid.

Build support is an OPTIONAL aspect of the Compose specification, and is
described in detail in the [Build support](build.md) documentation.

Each Service defines runtime constraints and requirements to run its containers. The `deploy` section groups
these constraints and allows the platform to adjust the deployment strategy to best match containers' needs with
available resources.

Deploy support is an OPTIONAL aspect of the Compose specification, and is
described in detail in the [Deployment support](deploy.md) documentation.
not implemented the Deploy section SHOULD be ignored and the Compose file MUST still be considered valid.

### build

`build` specifies the build configuration for creating container image from source, as defined in the [Build support](build.md) documentation.


### blkio_config

`blkio_config` defines a set of configuration options to set block IO limits for this service.

```yml
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
```

#### device_read_bps, device_write_bps

Set a limit in bytes per second for read / write operations on a given device.
Each item in the list MUST have two keys:

- `path`: defining the symbolic path to the affected device.
- `rate`: either as an integer value representing the number of bytes or as a string expressing a byte value.

#### device_read_iops, device_write_iops

Set a limit in operations per second for read / write operations on a given device.
Each item in the list MUST have two keys:

- `path`: defining the symbolic path to the affected device.
- `rate`: as an integer value representing the permitted number of operations per second.

#### weight

Modify the proportion of bandwidth allocated to this service relative to other services.
Takes an integer value between 10 and 1000, with 500 being the default.

#### weight_device

Fine-tune bandwidth allocation by device. Each item in the list must have two keys:

- `path`: defining the symbolic path to the affected device.
- `weight`: an integer value between 10 and 1000.

### cpu_count

`cpu_count` defines the number of usable CPUs for service container.

### cpu_percent

`cpu_percent` defines the usable percentage of the available CPUs.

### cpu_shares

`cpu_shares` defines (as integer value) service container relative CPU weight versus other containers.

### cpu_period

`cpu_period` allow Compose implementations to configure CPU CFS (Completely Fair Scheduler) period when platform is based
on Linux kernel.

### cpu_quota

`cpu_quota` allow Compose implementations to configure CPU CFS (Completely Fair Scheduler) quota when platform is based
on Linux kernel.

### cpu_rt_runtime

`cpu_rt_runtime` configures CPU allocation parameters for platform with support for realtime scheduler. Can be either
an integer value using microseconds as unit or a [duration](#specifying-durations).

```yml
 cpu_rt_runtime: '400ms'
 cpu_rt_runtime: 95000`
```

### cpu_rt_period

`cpu_rt_period` configures CPU allocation parameters for platform with support for realtime scheduler. Can be either
an integer value using microseconds as unit or a [duration](#specifying-durations).

```yml
 cpu_rt_period: '1400us'
 cpu_rt_period: 11000`
```

### cpus

_DEPRECATED: use [deploy.reservations.cpus](deploy.md#cpus)_

`cpus` define the number of (potentially virtual) CPUs to allocate to service containers. This is a fractional number.
`0.000` means no limit.

### cpuset

`cpuset` defines the explicit CPUs in which to allow execution. Can be a range `0-3` or a list `0,1`


### cap_add

`cap_add` specifies additional container [capabilities](http://man7.org/linux/man-pages/man7/capabilities.7.html)
as strings.

```
cap_add:
  - ALL
```

### cap_drop

`cap_drop` specifies container [capabilities](http://man7.org/linux/man-pages/man7/capabilities.7.html) to drop
as strings.

```
cap_drop:
  - NET_ADMIN
  - SYS_ADMIN
```

### cgroup_parent

`cgroup_parent` specifies an OPTIONAL parent [cgroup](http://man7.org/linux/man-pages/man7/cgroups.7.html) for the container.

```
cgroup_parent: m-executor-abcd
```

### command

`command` overrides the the default command declared by the container image (i.e. by Dockerfile's `CMD`).

```
command: bundle exec thin -p 3000
```

The command can also be a list, in a manner similar to [Dockerfile](https://docs.docker.com/engine/reference/builder/#cmd):

```
command: [ "bundle", "exec", "thin", "-p", "3000" ]
```

### configs

`configs` grant access to configs on a per-service basis using the per-service `configs`
configuration. Two different syntax variants are supported.

Compose implementations MUST report an error if config doesn't exist on platform or isn't defined in the
[`configs`](#configs-top-level-element) section of this Compose file.

There are two syntaxes defined for configs. To remain compliant to this specification, an implementation
MUST support both syntaxes. Implementations MUST allow use of both short and long syntaxes within the same document.

#### Short syntax

The short syntax variant only specifies the config name. This grants the
container access to the config and mounts it at `/<config_name>`
within the container. The source name and destination mount point are both set
to the config name.

The following example uses the short syntax to grant the `redis` service
access to the `my_config` and `my_other_config` configs. The value of
`my_config` is set to the contents of the file `./my_config.txt`, and
`my_other_config` is defined as an external resource, which means that it has
already been defined in the platform. If the external config does not exist,
the deployment MUST fail.

```yml
services:
  redis:
    image: redis:latest
    configs:
      - my_config
configs:
  my_config:
    file: ./my_config.txt
  my_other_config:
    external: true
```

#### Long syntax

The long syntax provides more granularity in how the config is created within the service's task containers.

- `source`: The name of the config as it exists in the platform.
- `target`: The path and name of the file to be mounted in the service's
  task containers. Defaults to `/<source>` if not specified.
- `uid` and `gid`: The numeric UID or GID that owns the mounted config file
  within the service's task containers. Default value when not specified is USER running container.
- `mode`: The [permissions](http://permissions-calculator.org/) for the file that is mounted within the service's
  task containers, in octal notation. Default value is world-readable (`0444`).
  Writable bit MUST be ignored. The executable bit can be set.

The following example sets the name of `my_config` to `redis_config` within the
container, sets the mode to `0440` (group-readable) and sets the user and group
to `103`. The `redis` service does not have access to the `my_other_config`
config.

```yml
services:
  redis:
    image: redis:latest
    configs:
      - source: my_config
        target: /redis_config
        uid: "103"
        gid: "103"
        mode: 0440
configs:
  my_config:
    external: true
  my_other_config:
    external: true
```

You can grant a service access to multiple configs, and you can mix long and short syntax.

### container_name

`container_name` is a string that specifies a custom container name, rather than a generated default name.

```yml
container_name: my-web-container
```

Compose implementation MUST NOT scale a service beyond one container if the Compose file specifies a
`container_name`. Attempting to do so MUST result in an error.

If present, `container_name` SHOULD follow the regex format of `[a-zA-Z0-9][a-zA-Z0-9_.-]+`

### credential_spec

`credential_spec` configures the credential spec for a managed service account.

Compose implementations that support services using Windows containers MUST support `file:` and
`registry:` protocols for credential_spec. Compose implementations MAY also support additional
protocols for custom use-cases.

The `credential_spec` must be in the format `file://<filename>` or `registry://<value-name>`.

```yml
credential_spec:
  file: my-credential-spec.json
```

When using `registry:`, the credential spec is read from the Windows registry on
the daemon's host. A registry value with the given name must be located in:

    HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization\Containers\CredentialSpecs

The following example loads the credential spec from a value named `my-credential-spec`
in the registry:

```yml
credential_spec:
  registry: my-credential-spec
```

#### Example gMSA configuration

When configuring a gMSA credential spec for a service, you only need
to specify a credential spec with `config`, as shown in the following example:

```yml
services:
  myservice:
    image: myimage:latest
    credential_spec:
      config: my_credential_spec

configs:
  my_credentials_spec:
    file: ./my-credential-spec.json|
```

### depends_on

`depends_on` expresses startup and shutdown dependencies between services.

#### Short syntax

The short syntax variant only specifies service names of the dependencies.
Service dependencies cause the following behaviors:

- Compose implementations MUST create services in dependency order. In the following
  example, `db` and `redis` are created before `web`.

- Compose implementations MUST remove services in dependency order. In the following
  example, `web` is removed before `db` and `redis`.

Simple example:

```yml
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
```

Compose implementations MUST guarantee dependency services have been started before
starting a dependent service.
Compose implementations MAY wait for dependency services to be "ready" before
starting a dependent service.

#### Long syntax

The long form syntax enables the configuration of additional fields that can't be
expressed in the short form.

- `condition`: condition under which dependency is considered satisfied
  - `service_started`: is an equivalent of the short syntax described above
  - `service_healthy`: specifies that a dependency is expected to be "healthy"
    (as indicated by [healthcheck](#healthcheck)) before starting a dependent
    service.
  - `service_completed_successfully`: specifies that a dependency is expected to run
    to successful completion before starting a dependent service.

Service dependencies cause the following behaviors:

- Compose implementations MUST create services in dependency order. In the following
  example, `db` and `redis` are created before `web`.

- Compose implementations MUST wait for healthchecks to pass on dependencies
  marked with `service_healthy`. In the following example, `db` is expected to
  be "healthy" before `web` is created.

- Compose implementations MUST remove services in dependency order. In the following
  example, `web` is removed before `db` and `redis`.

Simple example:

```yml
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
    image: postgres
```

Compose implementations MUST guarantee dependency services have been started before
starting a dependent service.
Compose implementations MUST guarantee dependency services marked with
`service_healthy` are "healthy" before starting a dependent service.


### deploy

`deploy` specifies the configuration for the deployment and lifecycle of services, as defined [here](deploy.md).


### device_cgroup_rules

`device_cgroup_rules` defines a list of device cgroup rules for this container.
The format is the same format the Linux kernel specifies in the [Control Groups
Device Whitelist Controller](https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v1/devices.html){: target="_blank" rel="noopener" class="_"}.

```yml
device_cgroup_rules:
  - 'c 1:3 mr'
  - 'a 7:* rmw'
```

### devices

`devices` defines a list of device mappings for created containers in the form of
`HOST_PATH:CONTAINER_PATH[:CGROUP_PERMISSIONS]`.

```yml
devices:
  - "/dev/ttyUSB0:/dev/ttyUSB0"
  - "/dev/sda:/dev/xvda:rwm"
```

### dns

`dns` defines custom DNS servers to set on the container network interface configuration. Can be a single value or a list.

```yml
dns: 8.8.8.8
```

```yml
dns:
  - 8.8.8.8
  - 9.9.9.9
```

### dns_opt

`dns_opt` list custom DNS options to be passed to the container’s DNS resolver (`/etc/resolv.conf` file on Linux).

```yml
dns_opt:
  - use-vc
  - no-tld-query
```

### dns_search

`dns` defines custom DNS search domains to set on container network interface configuration. Can be a single value or a list.

```yml
dns_search: example.com
```

```yml
dns_search:
  - dc1.example.com
  - dc2.example.com
```

### domainname

`domainname` declares a custom domain name to use for the service container. MUST be a valid RFC 1123 hostname.

### entrypoint

`entrypoint` overrides the default entrypoint for the Docker image (i.e. `ENTRYPOINT` set by Dockerfile).
Compose implementations MUST clear out any default command on the Docker image - both `ENTRYPOINT` and `CMD` instruction
in the Dockerfile - when `entrypoint` is configured by a Compose file. If [`command`](#command) is also set,
it is used as parameter to `entrypoint` as a replacement for Docker image's `CMD`

```yml
entrypoint: /code/entrypoint.sh
```

The entrypoint can also be a list, in a manner similar to
[Dockerfile](https://docs.docker.com/engine/reference/builder/#cmd):

```yml
entrypoint:
  - php
  - -d
  - zend_extension=/usr/local/lib/php/extensions/no-debug-non-zts-20100525/xdebug.so
  - -d
  - memory_limit=-1
  - vendor/bin/phpunit
```

### env_file

`env_file` adds environment variables to the container based on file content.

```yml
env_file: .env
```

`env_file` can also be a list. The files in the list MUST be processed from the top down. For the same variable
specified in two env files, the value from the last file in the list MUST stand.

```yml
env_file:
  - ./a.env
  - ./b.env
```

Relative path MUST be resolved from the Compose file's parent folder. As absolute paths prevent the Compose
file from being portable, Compose implementations SHOULD warn users when such a path is used to set `env_file`.

Environment variables declared in the [environment](#environment) section
MUST override these values – this holds true even if those values are
empty or undefined.

#### Env_file format

Each line in an env file MUST be in `VAR[=[VAL]]` format. Lines beginning with `#` MUST be ignored.
Blank lines MUST also be ignored.

The value of `VAL` is used as a raw string and not modified at all. If the value is surrounded by quotes
(as is often the case for shell variables), the quotes MUST be **included** in the value passed to containers
created by the Compose implementation.

`VAL` MAY be omitted, in such cases the variable value is empty string.
`=VAL` MAY be omitted, in such cases the variable is **unset**.

```bash
# Set Rails/Rack environment
RACK_ENV=development
VAR="quoted"
```

### environment

`environment` defines environment variables set in the container. `environment` can use either an array or a
map. Any boolean values; true, false, yes, no, SHOULD be enclosed in quotes to ensure
they are not converted to True or False by the YAML parser.

Environment variables MAY be declared by a single key (no value to equals sign). In such a case Compose
implementations SHOULD rely on some user interaction to resolve the value. If they do not, the variable
is unset and will be removed from the service container environment.

Map syntax:

```yml
environment:
  RACK_ENV: development
  SHOW: "true"
  USER_INPUT:
```

Array syntax:

```yml
environment:
  - RACK_ENV=development
  - SHOW=true
  - USER_INPUT
```

When both `env_file` and `environment` are set for a service, values set by `environment` have precedence.

### expose

`expose` defines the ports that Compose implementations MUST expose from container. These ports MUST be
accessible to linked services and SHOULD NOT be published to the host machine. Only the internal container
ports can be specified.

```yml
expose:
  - "3000"
  - "8000"
```

### extends

Extend another service, in the current file or another, optionally overriding configuration. You can use
`extends` on any service together with other configuration keys. The `extends` value MUST be a mapping
defined with a required `service` and an optional `file` key.

```yaml
extends:
  file: common.yml
  service: webapp
```

If supported Compose implementations MUST process `extends` in the following way:

- `service` defines the name of the service being referenced as a base, for example `web` or `database`.
- `file` is the location of a Compose configuration file defining that service.

#### Restrictions

The following restrictions apply to the service being referenced:

- Services that have dependencies on other services cannot be used as a base. Therefore, any key
  that introduces a dependency on another service is incompatible with `extends`. The
  non-exhaustive list of such keys is: `links`, `volumes_from`, `container` mode (in `ipc`, `pid`,
  `network_mode` and `net`), `service` mode (in `ipc`, `pid` and `network_mode`), `depends_on`.
- Services cannot have circular references with `extends`

Compose implementations MUST return an error in all of these cases.

#### Finding referenced service

`file` value can be:

- Not present.
  This indicates that another service within the same Compose file is being referenced.
- File path, which can be either:
  - Relative path. This path is considered as relative to the location of the main Compose
    file.
  - Absolute path.

Service denoted by `service` MUST be present in the identified referenced Compose file.
Compose implementations MUST return an error if:

- Service denoted by `service` was not found
- Compose file denoted by `file` was not found

#### Merging service definitions

Two service definitions (_main_ one in the current Compose file and _referenced_ one
specified by `extends`) MUST be merged in the following way:

- Mappings: keys in mappings of _main_ service definition override keys in mappings
  of _referenced_ service definition. Keys that aren't overridden are included as is.
- Sequences: items are combined together into an new sequence. Order of elements is
  preserved with the _referenced_ items coming first and _main_ items after.
- Scalars: keys in _main_ service definition take precedence over keys in the
  _referenced_ one.

##### Mappings

The following keys should be treated as mappings: `build.args`, `build.labels`,
`build.extra_hosts`, `deploy.labels`, `deploy.update_config`, `deploy.rollback_config`,
`deploy.restart_policy`, `deploy.resources.limits`, `environment`, `healthcheck`,
`labels`, `logging.options`, `sysctls`, `storage_opt`, `extra_hosts`, `ulimits`.

One exception that applies to `healthcheck` is that _main_ mapping cannot specify
`disable: true` unless _referenced_ mapping also specifies `disable: true`. Compose
implementations MUST return an error in this case.

For example, the input below:

```yaml
services:
  common:
    image: busybox
    environment:
      TZ: utc
      PORT: 80
  cli:
    extends:
      service: common
    environment:
      PORT: 8080
```

Produces the following configuration for the `cli` service. The same output is
produced if array syntax is used.

```yaml
environment:
  PORT: 8080
  TZ: utc
image: busybox
```

Items under `blkio_config.device_read_bps`, `blkio_config.device_read_iops`,
`blkio_config.device_write_bps`, `blkio_config.device_write_iops`, `devices` and
`volumes` are also treated as mappings where key is the target path inside the
container.

For example, the input below:

```yaml
services:
  common:
    image: busybox
    volumes:
      - common-volume:/var/lib/backup/data:rw
  cli:
    extends:
      service: common
    volumes:
      - cli-volume:/var/lib/backup/data:ro
```

Produces the following configuration for the `cli` service. Note that mounted path
now points to the new volume name and `ro` flag was applied.

```yaml
image: busybox
volumes:
- cli-volume:/var/lib/backup/data:ro
```

If _referenced_ service definition contains `extends` mapping, the items under it
are simply copied into the new _merged_ definition. Merging process is then kicked
off again until no `extends` keys are remaining.

For example, the input below:

```yaml
services:
  base:
    image: busybox
    user: root
  common:
    image: busybox
    extends:
      service: base
  cli:
    extends:
      service: common
```

Produces the following configuration for the `cli` service. Here, `cli` services
gets `user` key from `common` service, which in turn gets this key from `base`
service.

```yaml
image: busybox
user: root
```

##### Sequences

The following keys should be treated as sequences: `cap_add`, `cap_drop`, `configs`,
`deploy.placement.constraints`, `deploy.placement.preferences`,
`deploy.reservations.generic_resources`, `device_cgroup_rules`, `expose`,
`external_links`, `ports`, `secrets`, `security_opt`.
Any duplicates resulting from the merge are removed so that the sequence only
contains unique elements.

For example, the input below:

```yaml
services:
  common:
    image: busybox
    security_opt:
      - label:role:ROLE
  cli:
    extends:
      service: common
    security_opt:
      - label:user:USER
```

Produces the following configuration for the `cli` service.

```yaml
image: busybox
security_opt:
- label:role:ROLE
- label:user:USER
```

In case list syntax is used, the following keys should also be treated as sequences:
`dns`, `dns_search`, `env_file`, `tmpfs`. Unlike sequence fields mentioned above,
duplicates resulting from the merge are not removed.

##### Scalars

Any other allowed keys in the service definition should be treated as scalars.

### external_links

`external_links` link service containers to services managed outside this Compose application.
`external_links` define the name of an existing service to retrieve using the platform lookup mechanism.
An alias of the form `SERVICE:ALIAS` can be specified.

```yml
external_links:
  - redis
  - database:mysql
  - database:postgresql
```

### extra_hosts

`extra_hosts` adds hostname mappings to the container network interface configuration (`/etc/hosts` for Linux).
Values MUST set hostname and IP address for additional hosts in the form of `HOSTNAME:IP`.

```yml
extra_hosts:
  - "somehost:162.242.195.82"
  - "otherhost:50.31.209.229"
```

Compose implementations MUST create matching entry with the IP address and hostname in the container's network
configuration, which means for Linux `/etc/hosts` will get extra lines:

```
162.242.195.82  somehost
50.31.209.229   otherhost
```

### group_add

`group_add` specifies additional groups (by name or number) which the user inside the container MUST be a member of.

An example of where this is useful is when multiple containers (running as different users) need to all read or write
the same file on a shared volume. That file can be owned by a group shared by all the containers, and specified in
`group_add`.

```yml
services:
  myservice:
    image: alpine
    group_add:
      - mail
```

Running `id` inside the created container MUST show that the user belongs to the `mail` group, which would not have
been the case if `group_add` were not declared.

### healthcheck

`healthcheck` declares a check that's run to determine whether or not containers for this
service are "healthy". This overrides
[HEALTHCHECK Dockerfile instruction](https://docs.docker.com/engine/reference/builder/#healthcheck)
set by the service's Docker image.

```yml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost"]
  interval: 1m30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

`interval`, `timeout` and `start_period` are [specified as durations](#specifying-durations).

`test` defines the command the Compose implementation will run to check container health. It can be
either a string or a list. If it's a list, the first item must be either `NONE`, `CMD` or `CMD-SHELL`.
If it's a string, it's equivalent to specifying `CMD-SHELL` followed by that string.

```yml
# Hit the local web app
test: ["CMD", "curl", "-f", "http://localhost"]
```

Using `CMD-SHELL` will run the command configured as a string using the container's default shell
(`/bin/sh` for Linux). Both forms below are equivalent:

```yml
test: ["CMD-SHELL", "curl -f http://localhost || exit 1"]
```

```yml
test: curl -f https://localhost || exit 1
```

`NONE` disable the healthcheck, and is mostly useful to disable Healthcheck set by image. Alternatively
the healthcheck set by the image can be disabled by setting `disable: true`:

```yml
healthcheck:
  disable: true
```

### hostname

`hostname` declares a custom host name to use for the service container. MUST be a valid RFC 1123 hostname.

### image

`image` specifies the image to start the container from. Image MUST follow the Open Container Specification
[addressable image format](https://github.com/opencontainers/org/blob/master/docs/docs/introduction/digests.md),
as `[<registry>/][<project>/]<image>[:<tag>|@<digest>]`.

```yml
    image: redis
    image: redis:5
    image: redis@sha256:0ed5d5928d4737458944eb604cc8509e245c3e19d02ad83935398bc4b991aac7
    image: library/redis
    image: docker.io/library/redis
    image: my_private.registry:5000/redis
```

If the image does not exist on the platform, Compose implementations MUST attempt to pull it based on the `pull_policy`.
Compose implementations with build support MAY offer alternative options for the end user to control precedence of
pull over building the image from source, however pulling the image MUST be the default behavior.

`image` MAY be omitted from a Compose file as long as a `build` section is declared. Compose implementations
without build support MUST fail when `image` is missing from the Compose file.

### init

`init` run an init process (PID 1) inside the container that forwards signals and reaps processes.
Set this option to `true` to enable this feature for the service.

```yml
services:
  web:
    image: alpine:latest
    init: true
```

The init binary that is used is platform specific.

### ipc

`ipc` configures the IPC isolation mode set by service container. Available
values are platform specific, but Compose specification defines specific values
which MUST be implemented as described if supported:

- `shareable` which gives the container own private IPC namespace, with a
  possibility to share it with other containers.
- `service:{name}` which makes the container join another (`shareable`)
   container's IPC namespace.

```yml
    ipc: "shareable"
    ipc: "service:[service name]"
```

### isolation

`isolation` specifies a container’s isolation technology. Supported values are platform-specific.

### labels

`labels` add metadata to containers. You can use either an array or a map.

It's recommended that you use reverse-DNS notation to prevent your labels from conflicting with
those used by other software.

```yml
labels:
  com.example.description: "Accounting webapp"
  com.example.department: "Finance"
  com.example.label-with-empty-value: ""
```

```yml
labels:
  - "com.example.description=Accounting webapp"
  - "com.example.department=Finance"
  - "com.example.label-with-empty-value"
```

Compose implementations MUST create containers with canonical labels:

- `com.docker.compose.project` set on all resources created by Compose implementation to the user project name
- `com.docker.compose.service` set on service containers with service name as defined in the Compose file

The `com.docker.compose` label prefix is reserved. Specifying labels with this prefix in the Compose file MUST
result in a runtime error.

### links

`links` defines a network link to containers in another service. Either specify both the service name and
a link alias (`SERVICE:ALIAS`), or just the service name.

```yml
web:
  links:
    - db
    - db:database
    - redis
```

Containers for the linked service MUST be reachable at a hostname identical to the alias, or the service name
if no alias was specified.

Links are not required to enable services to communicate - when no specific network configuration is set,
any service MUST be able to reach any other service at that service’s name on the `default` network. If services
do declare networks they are attached to, `links` SHOULD NOT override the network configuration and services not
attached to a shared network SHOULD NOT be able to communicate. Compose implementations MAY NOT warn the user
about this configuration mismatch.

Links also express implicit dependency between services in the same way as
[depends_on](#depends_on), so they determine the order of service startup.

### logging

`logging` defines the logging configuration for the service.

```yml
logging:
  driver: syslog
  options:
    syslog-address: "tcp://192.168.0.42:123"
```

The `driver` name specifies a logging driver for the service's containers. The default and available values
are platform specific. Driver specific options can be set with `options` as key-value pairs.

### network_mode

`network_mode` set service containers network mode. Available values are platform specific, but Compose
specification define specific values which MUST be implemented as described if supported:

- `none` which disable all container networking
- `host` which gives the container raw access to host's network interface
- `service:{name}` which gives the containers access to the specified service only

```yml
    network_mode: "host"
    network_mode: "none"
    network_mode: "service:[service name]"
```

### networks

`networks` defines the networks that service containers are attached to, referencing entries under the
[top-level `networks` key](#networks-top-level-element).

```yml
services:
  some-service:
    networks:
      - some-network
      - other-network
```

#### aliases

`aliases` declares alternative hostnames for this service on the network. Other containers on the same
network can use either the service name or this alias to connect to one of the service's containers.

Since `aliases` are network-scoped, the same service can have different aliases on different networks.

> **Note**: A network-wide alias can be shared by multiple containers, and even by multiple services.
> If it is, then exactly which container the name resolves to is not guaranteed.

The general format is shown here:

```yml
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
```

In the example below, service `frontend` will be able to reach the `backend` service at
the hostname `backend` or `database` on the `back-tier` network, and service `monitoring`
will be able to reach same `backend` service at `db` or `mysql` on the `admin` network.

```yml
services:
  frontend:
    image: awesome/webapp
    networks:
      - front-tier
      - back-tier

  monitoring:
    image: awesome/monitoring
    networks:
      - admin

  backend:
    image: awesome/backend
    networks:
      back-tier:
        aliases:
          - database
      admin:
        aliases:
          - mysql

networks:
  front-tier:
  back-tier:
  admin:
```

#### ipv4_address, ipv6_address

Specify a static IP address for containers for this service when joining the network.

The corresponding network configuration in the [top-level networks section](#networks) MUST have an
`ipam` block with subnet configurations covering each static address.

```yml
services:
  frontend:
    image: awesome/webapp
    networks:
      front-tier:
        ipv4_address: 172.16.238.10
        ipv6_address: 2001:3984:3989::10

networks:
  front-tier:
    ipam:
      driver: default
      config:
        - subnet: "172.16.238.0/24"
        - subnet: "2001:3984:3989::/64"
```

#### link_local_ips

`link_local_ips` specifies a list of link-local IPs. Link-local IPs are special IPs which belong to a well
known subnet and are purely managed by the operator, usually dependent on the architecture where they are
deployed. Implementation is Platform specific.

Example:

```yaml
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
```

#### priority

`priority` indicates in which order Compose implementation SHOULD connect the service’s containers to its
networks. If unspecified, the default value is 0.

In the following example, the app service connects to app_net_1 first as it has the highest priority. It then connects to app_net_3, then app_net_2, which uses the default priority value of 0.

```yaml
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
```

### mac_address

`mac_address` sets a MAC address for service container.

### mem_limit

_DEPRECATED: use [deploy.limits.memory](deploy.md#memory)_

### mem_reservation

_DEPRECATED: use [deploy.reservations.memory](deploy.md#memory)_

### mem_swappiness

`mem_swappiness` defines as a percentage (a value between 0 and 100) for the host kernel to swap out
anonymous memory pages used by a container.

- a value of 0 turns off anonymous page swapping.
- a value of 100 sets all anonymous pages as swappable.

Default value is platform specific.

### memswap_limit

`memswap_limit` defines the amount of memory container is allowed to swap to disk. This is a modifier
attribute that only has meaning if `memory` is also set. Using swap allows the container to write excess
memory requirements to disk when the container has exhausted all the memory that is available to it.
There is a performance penalty for applications that swap memory to disk often.

- If `memswap_limit` is set to a positive integer, then both `memory` and `memswap_limit` MUST be set. `memswap_limit` represents the total amount of memory and swap that can be used, and `memory` controls the amount used by non-swap memory. So if `memory`="300m" and `memswap_limit`="1g", the container can use 300m of memory and 700m (1g - 300m) swap.
- If `memswap_limit` is set to 0, the setting MUST be ignored, and the value is treated as unset.
- If `memswap_limit` is set to the same value as `memory`, and `memory` is set to a positive integer, the container does not have access to swap. See Prevent a container from using swap.
- If `memswap_limit` is unset, and `memory` is set, the container can use as much swap as the `memory` setting, if the host container has swap memory configured. For instance, if `memory`="300m" and `memswap_limit` is not set, the container can use 600m in total of memory and swap.
- If `memswap_limit` is explicitly set to -1, the container is allowed to use unlimited swap, up to the amount available on the host system.

### oom_kill_disable

If `oom_kill_disable` is set Compose implementation MUST configure the platform so it won't kill the container in case
of memory starvation.

### oom_score_adj

`oom_score_adj` tunes the preference for containers to be killed by platform in case of memory starvation. Value MUST
be within [-1000,1000] range.

### pid

`pid` sets the PID mode for container created by the Compose implementation.
Supported values are platform specific.

### pids_limit

_DEPRECATED: use [deploy.reservations.pids](deploy.md#pids)_

`pids_limit` tunes a container’s PIDs limit. Set to -1 for unlimited PIDs.

```yml
pids_limit: 10
```

### platform

`platform` defines the target platform containers for this service will run on, using the `os[/arch[/variant]]` syntax.
Compose implementation MUST use this attribute when declared to determine which version of the image will be pulled
and/or on which platform the service’s build will be performed.

```yml
platform: osx
platform: windows/amd64
platform: linux/arm64/v8
```

### ports

Exposes container ports.
Port mapping MUST NOT be used with `network_mode: host` and doing so MUST result in a runtime error.

#### Short syntax

The short syntax is a colon-separated string to set host IP, host port and container port
in the form:

`[HOST:]CONTAINER[/PROTOCOL]` where:

- `HOST` is `[IP:](port | range)`
- `CONTAINER` is `port | range`
- `PROTOCOL` to restrict port to specified protocol. `tcp` and `udp` values are defined by the specification,
  Compose implementations MAY offer support for platform-specific protocol names.

Host IP, if not set, MUST bind to all network interfaces. Port can be either a single
value or a range. Host and container MUST use equivalent ranges.

Either specify both ports (`HOST:CONTAINER`), or just the container port. In the latter case, the
Compose implementation SHOULD automatically allocate any unassigned host port.

`HOST:CONTAINER` SHOULD always be specified as a (quoted) string, to avoid conflicts
with [yaml base-60 float](https://yaml.org/type/float.html){: target="_blank" rel="noopener" class="_"}.

Samples:

```yml
ports:
  - "3000"
  - "3000-3005"
  - "8000:8000"
  - "9090-9091:8080-8081"
  - "49100:22"
  - "127.0.0.1:8001:8001"
  - "127.0.0.1:5000-5010:5000-5010"
  - "6060:6060/udp"
```

> **Note**: Host IP mapping MAY not be supported on the platform, in such case Compose implementations SHOULD reject
> the Compose file and MUST inform the user they will ignore the specified host IP.

#### Long syntax

The long form syntax allows the configuration of additional fields that can't be
expressed in the short form.

- `target`: the container port
- `published`: the publicly exposed port. Can be set as a range using syntax `start-end`, then actual port SHOULD be assigned within this range based on available ports.
- `host_ip`: the Host IP mapping, unspecified means all network interfaces (`0.0.0.0`) 
- `protocol`: the port protocol (`tcp` or `udp`), unspecified means any protocol
- `mode`: `host` for publishing a host port on each node, or `ingress` for a port to be load balanced.

```yml
ports:
  - target: 80
    host_ip: 127.0.0.1
    published: 8080
    protocol: tcp
    mode: host

  - target: 80
    host_ip: 127.0.0.1
    published: 8000-9000
    protocol: tcp
    mode: host
```

### privileged

`privileged` configures the service container to run with elevated privileges. Support and actual impacts are platform-specific.

### profiles

`profiles` defines a list of named profiles for the service to be enabled under. When not set, service is always enabled.

If present, `profiles` SHOULD follow the regex format of `[a-zA-Z0-9][a-zA-Z0-9_.-]+`.

### pull_policy

`pull_policy` defines the decisions Compose implementations will make when it starts to pull images. Possible values are:

* `always`: Compose implementations SHOULD always pull the image from the registry.
* `never`: Compose implementations SHOULD NOT pull the image from a registry and SHOULD rely on the platform cached image.
   If there is no cached image, a failure MUST be reported.
* `missing`: Compose implementations SHOULD pull the image only if it's not available in the platform cache.
   This SHOULD be the default option for Compose implementations without build support.
  `if_not_present` SHOULD be considered an alias for this value for backward compatibility
* `build`: Compose implementations SHOULD build the image. Compose implementations SHOULD rebuild the image if already present.

If `pull_policy` and `build` both presents, Compose implementations SHOULD build the image by default. Compose implementations MAY override this behavior in the toolchain.

### read_only

`read_only` configures service container to be created with a read-only filesystem.

### restart

`restart` defines the policy that the platform will apply on container termination.

- `no`: The default restart policy. Does not restart a container under any circumstances.
- `always`: The policy always restarts the container until its removal.
- `on-failure`: The policy restarts a container if the exit code indicates an error.
- `unless-stopped`: The policy restarts a container irrespective of the exit code but will stop
  restarting when the service is stopped or removed.

```yml
    restart: "no"
    restart: always
    restart: on-failure
    restart: unless-stopped
```

### runtime

`runtime` specifies which runtime to use for the service’s containers.

The value of `runtime` is specific to implementation.
For example, `runtime` can be the name of [an implementation of OCI Runtime Spec](https://github.com/opencontainers/runtime-spec/blob/master/implementations.md){: target="_blank" rel="noopener" class="_"}, such as "runc".

```yml
web:
  image: busybox:latest
  command: true
  runtime: runc
```

### scale

-DEPRECATED: use [deploy/replicas](deploy.md#replicas)_

`scale` specifies the default number of containers to deploy for this service.

### secrets

`secrets` grants access to sensitive data defined by [secrets](#secrets) on a per-service basis. Two
different syntax variants are supported: the short syntax and the long syntax.

Compose implementations MUST report an error if the secret doesn't exist on the platform or isn't defined in the
[`secrets`](#secrets-top-level-element) section of this Compose file.

#### Short syntax

The short syntax variant only specifies the secret name. This grants the
container access to the secret and mounts it as read-only to `/run/secrets/<secret_name>`
within the container. The source name and destination mountpoint are both set
to the secret name.

The following example uses the short syntax to grant the `frontend` service
access to the `server-certificate` secret. The value of `server-certificate` is set
to the contents of the file `./server.cert`.

```yml
services:
  frontend:
    image: awesome/webapp
    secrets:
      - server-certificate
secrets:
  server-certificate:
    file: ./server.cert
```

#### Long syntax

The long syntax provides more granularity in how the secret is created within
the service's containers.

- `source`: The name of the secret as it exists on the platform.
- `target`: The name of the file to be mounted in `/run/secrets/` in the
  service's task containers. Defaults to `source` if not specified.
- `uid` and `gid`: The numeric UID or GID that owns the file within
  `/run/secrets/` in the service's task containers. Default value is USER running container.
- `mode`: The [permissions](http://permissions-calculator.org/) for the file to be mounted in `/run/secrets/`
  in the service's task containers, in octal notation.
  Default value is world-readable permissions (mode `0444`).
  The writable bit MUST be ignored if set. The executable bit MAY be set.

The following example sets the name of the `server-certificate` secret file to `server.crt`
within the container, sets the mode to `0440` (group-readable) and sets the user and group
to `103`. The value of `server-certificate` secret is provided by the platform through a lookup and
the secret lifecycle not directly managed by the Compose implementation.

```yml
services:
  frontend:
    image: awesome/webapp
    secrets:
      - source: server-certificate
        target: server.cert
        uid: "103"
        gid: "103"
        mode: 0440
secrets:
  server-certificate:
    external: true
```

Services MAY be granted access to multiple secrets. Long and short syntax for secrets MAY be used in the
same Compose file. Defining a secret in the top-level `secrets` MUST NOT imply granting any service access to it.
Such grant must be explicit within service specification as [secrets](#secrets) service element.

### security_opt

`security_opt` overrides the default labeling scheme for each container.

```yml
security_opt:
  - label:user:USER
  - label:role:ROLE
```

### shm_size

`shm_size` configures the size of the shared memory (`/dev/shm` partition on Linux) allowed by the service container.
Specified as a [byte value](#specifying-byte-values).

### stdin_open

`stdin_open` configures service containers to run with an allocated stdin.

### stop_grace_period

`stop_grace_period` specifies how long the Compose implementation MUST wait when attempting to stop a container if it doesn't
handle SIGTERM (or whichever stop signal has been specified with
[`stop_signal`](#stopsignal)), before sending SIGKILL. Specified
as a [duration](#specifying-durations).

```yml
    stop_grace_period: 1s
    stop_grace_period: 1m30s
```

Default value is 10 seconds for the container to exit before sending SIGKILL.

### stop_signal

`stop_signal` defines the signal that the Compose implementation MUST use to stop the service containers.
If unset containers are stopped by the Compose Implementation by sending `SIGTERM`.

```yml
stop_signal: SIGUSR1
```

### storage_opt

`storage_opt` defines storage driver options for a service.

```yml
storage_opt:
  size: '1G'
```

### sysctls

`sysctls` defines kernel parameters to set in the container. `sysctls` can use either an array or a map.

```yml
sysctls:
  net.core.somaxconn: 1024
  net.ipv4.tcp_syncookies: 0
```

```yml
sysctls:
  - net.core.somaxconn=1024
  - net.ipv4.tcp_syncookies=0
```

You can only use sysctls that are namespaced in the kernel. Docker does not
support changing sysctls inside a container that also modify the host system.
For an overview of supported sysctls, refer to [configure namespaced kernel
parameters (sysctls) at runtime](https://docs.docker.com/engine/reference/commandline/run/#configure-namespaced-kernel-parameters-sysctls-at-runtime).

### tmpfs

`tmpfs` mounts a temporary file system inside the container. Can be a single value or a list.

```yml
tmpfs: /run
```

```yml
tmpfs:
  - /run
  - /tmp
```

### tty

`tty` configure service container to run with a TTY.

### ulimits

`ulimits` overrides the default ulimits for a container. Either specifies as a single limit as an integer or
soft/hard limits as a mapping.

```yml
ulimits:
  nproc: 65535
  nofile:
    soft: 20000
    hard: 40000
```

### user

`user` overrides the user used to run the container process. Default is that set by image (i.e. Dockerfile `USER`),
if not set, `root`.

### userns_mode

`userns_mode` sets the user namespace for the service. Supported values are platform specific and MAY depend
on platform configuration

```yml
userns_mode: "host"
```

### volumes

`volumes` defines mount host paths or named volumes that MUST be accessible by service containers.

If the mount is a host path and only used by a single service, it MAY be declared as part of the service
definition instead of the top-level `volumes` key.

To reuse a volume across multiple services, a named
volume MUST be declared in the [top-level `volumes` key](#volumes-top-level-element).

This example shows a named volume (`db-data`) being used by the `backend` service,
and a bind mount defined for a single service

```yml
services:
  backend:
    image: awesome/backend
    volumes:
      - type: volume
        source: db-data
        target: /data
        volume:
          nocopy: true
      - type: bind
        source: /var/run/postgres/postgres.sock
        target: /var/run/postgres/postgres.sock

volumes:
  db-data:
```

#### Short syntax

The short syntax uses a single string with colon-separated values to specify a volume mount
(`VOLUME:CONTAINER_PATH`), or an access mode (`VOLUME:CONTAINER_PATH:ACCESS_MODE`).

- `VOLUME`: MAY be either a host path on the platform hosting containers (bind mount) or a volume name
- `CONTAINER_PATH`: the path in the container where the volume is mounted
- `ACCESS_MODE`: is a comma-separated `,` list of options and MAY be set to:
  - `rw`: read and write access (default)
  - `ro`: read-only access
  - `z`: SELinux option indicates that the bind mount host content is shared among multiple containers
  - `Z`: SELinux option indicates that the bind mount host content is private and unshared for other containers

> **Note**: The SELinux re-labeling bind mount option is ignored on platforms without SELinux.

> **Note**: Relative host paths MUST only be supported by Compose implementations that deploy to a
> local container runtime. This is because the relative path is resolved from the Compose file’s parent
> directory which is only applicable in the local case. Compose Implementations deploying to a non-local
> platform MUST reject Compose files which use relative host paths with an error. To avoid ambiguities
> with named volumes, relative paths SHOULD always begin with `.` or `..`.

#### Long syntax

The long form syntax allows the configuration of additional fields that can't be
expressed in the short form.

- `type`: the mount type `volume`, `bind`, `tmpfs` or `npipe`
- `source`: the source of the mount, a path on the host for a bind mount, or the
  name of a volume defined in the
  [top-level `volumes` key](#volumes-top-level-element). Not applicable for a tmpfs mount.
- `target`: the path in the container where the volume is mounted
- `read_only`: flag to set the volume as read-only
- `bind`: configure additional bind options
  - `propagation`: the propagation mode used for the bind
  - `create_host_path`: create a directory at the source path on host if there is nothing present. 
    Do nothing if there is something present at the path. This is automatically implied by short syntax
    for backward compatibility with docker-compose legacy.
  - `selinux`: the SELinux re-labeling option `z` (shared) or `Z` (private)
- `volume`: configure additional volume options
  - `nocopy`: flag to disable copying of data from a container when a volume is created
- `tmpfs`: configure additional tmpfs options
  - `size`: the size for the tmpfs mount in bytes (either numeric or as bytes unit)
- `consistency`: the consistency requirements of the mount. Available values are platform specific

### volumes_from

`volumes_from` mounts all of the volumes from another service or container, optionally specifying
read-only access (ro) or read-write (rw). If no access level is specified, then read-write MUST be used.

String value defines another service in the Compose application model to mount volumes from. The
`container:` prefix, if supported, allows to mount volumes from a container that is not managed by the
Compose implementation.

```yaml
volumes_from:
  - service_name
  - service_name:ro
  - container:container_name
  - container:container_name:rw
```

### working_dir

`working_dir` overrides the container's working directory from that specified by image (i.e. Dockerfile `WORKDIR`).

## Networks top-level element

Networks are the layer that allow services to communicate with each other. The networking model exposed to a service
is limited to a simple IP connection with target services and external resources, while the Network definition allows
fine-tuning the actual implementation provided by the platform.

Networks can be created by specifying the network name under a top-level `networks` section.
Services can connect to networks by specifying the network name under the service [`networks`](#networks) subsection

In the following example, at runtime, networks `front-tier` and `back-tier` will be created and the `frontend` service
connected to the `front-tier` network and the `back-tier` network.

```yml
services:
  frontend:
    image: awesome/webapp
    networks:
      - front-tier
      - back-tier

networks:
  front-tier:
  back-tier:
```

### driver

`driver` specifies which driver should be used for this network. Compose implementations MUST return an error if the
driver is not available on the platform.

```yml
driver: overlay
```

Default and available values are platform specific. Compose specification MUST support the following specific drivers:
`none` and `host`

- `host` use the host's networking stack
- `none` disable networking

#### host or none

The syntax for using built-in networks such as `host` and `none` is different, as such networks implicitly exists outside
the scope of the Compose implementation. To use them one MUST define an external network with the name `host` or `none` and
an alias that the Compose implementation can use (`hostnet` or `nonet` in the following examples), then grant the service
access to that network using its alias.

```yml
services:
  web:
    networks:
      hostnet: {}

networks:
  hostnet:
    external: true
    name: host
```

```yml
services:
  web:
    ...
    networks:
      nonet: {}

networks:
  nonet:
    external: true
    name: none
```

### driver_opts

`driver_opts` specifies a list of options as key-value pairs to pass to the driver for this network. These options are
driver-dependent - consult the driver's documentation for more information. Optional.

```yml
driver_opts:
  foo: "bar"
  baz: 1
```

### attachable

If `attachable` is set to `true`, then standalone containers SHOULD be able attach to this network, in addition to services.
If a standalone container attaches to the network, it can communicate with services and other standalone containers
that are also attached to the network.

```yml
networks:
  mynet1:
    driver: overlay
    attachable: true
```

### enable_ipv6

`enable_ipv6` enable IPv6 networking on this network.

### ipam

`ipam` specifies custom a IPAM configuration. This is an object with several properties, each of which is optional:

- `driver`: Custom IPAM driver, instead of the default.
- `config`: A list with zero or more configuration elements, each containing:
  - `subnet`: Subnet in CIDR format that represents a network segment
  - `ip_range`: Range of IPs from which to allocate container IPs
  - `gateway`: IPv4 or IPv6 gateway for the master subnet
  - `aux_addresses`: Auxiliary IPv4 or IPv6 addresses used by Network driver, as a mapping from hostname to IP
- `options`: Driver-specific options as a key-value mapping.

A full example:

```yml
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
```

### internal

By default, Compose implementations MUST provides external connectivity to networks. `internal` when set to `true` allow to
create an externally isolated network.

### labels

Add metadata to containers using Labels. Can use either an array or a dictionary.

Users SHOULD use reverse-DNS notation to prevent labels from conflicting with those used by other software.

```yml
labels:
  com.example.description: "Financial transaction network"
  com.example.department: "Finance"
  com.example.label-with-empty-value: ""
```

```yml
labels:
  - "com.example.description=Financial transaction network"
  - "com.example.department=Finance"
  - "com.example.label-with-empty-value"
```

Compose implementations MUST set `com.docker.compose.project` and `com.docker.compose.network` labels.

### external

If set to `true`, `external` specifies that this network’s lifecycle is maintained outside of that of the application.
Compose Implementations SHOULD NOT attempt to create these networks, and raises an error if one doesn't exist.

In the example below, `proxy` is the gateway to the outside world. Instead of attempting to create a network, Compose
implementations SHOULD interrogate the platform for an existing network simply called `outside` and connect the
`proxy` service's containers to it.

```yml

services:
  proxy:
    image: awesome/proxy
    networks:
      - outside
      - default
  app:
    image: awesome/app
    networks:
      - default

networks:
  outside:
    external: true
```

### name

`name` sets a custom name for this network. The name field can be used to reference networks which contain special characters.
The name is used as is and will **not** be scoped with the project name.

```yml
networks:
  network1:
    name: my-app-net
```

It can also be used in conjunction with the `external` property to define the platform network that the Compose implementation
should retrieve, typically by using a parameter so the Compose file doesn't need to hard-code runtime specific values:

```yml
networks:
  network1:
    external: true
    name: "${NETWORK_ID}"
```

## Volumes top-level element

Volumes are persistent data stores implemented by the platform. The Compose specification offers a neutral abstraction
for services to mount volumes, and configuration parameters to allocate them on infrastructure.

The `volumes` section allows the configuration of named volumes that can be reused across multiple services. Here's
an example of a two-service setup where a database's data directory is shared with another service as a volume named
`db-data` so that it can be periodically backed up:

```yml
services:
  backend:
    image: awesome/database
    volumes:
      - db-data:/etc/data

  backup:
    image: backup-service
    volumes:
      - db-data:/var/lib/backup/data

volumes:
  db-data:
```

An entry under the top-level `volumes` key can be empty, in which case it uses the platform's default configuration for
creating a volume. Optionally, you can configure it with the following keys:

### driver

Specify which volume driver should be used for this volume. Default and available values are platform specific. If the driver is not available, the Compose implementation MUST return an error and stop application deployment.

```yml
driver: foobar
```

### driver_opts

`driver_opts` specifies a list of options as key-value pairs to pass to the driver for this volume. Those options are driver-dependent.

```yml
volumes:
  example:
    driver_opts:
      type: "nfs"
      o: "addr=10.40.0.199,nolock,soft,rw"
      device: ":/docker/example"
```

### external

If set to `true`, `external` specifies that this volume already exist on the platform and its lifecycle is managed outside
of that of the application. Compose implementations MUST NOT attempt to create these volumes, and MUST return an error if they
do not exist.

In the example below, instead of attempting to create a volume called
`{project_name}_db-data`, Compose looks for an existing volume simply
called `db-data` and mounts it into the `backend` service's containers.

```yml
services:
  backend:
    image: awesome/database
    volumes:
      - db-data:/etc/data

volumes:
  db-data:
    external: true
```

### labels

`labels` are used to add metadata to volumes. You can use either an array or a dictionary.

It's recommended that you use reverse-DNS notation to prevent your labels from
conflicting with those used by other software.

```yml
labels:
  com.example.description: "Database volume"
  com.example.department: "IT/Ops"
  com.example.label-with-empty-value: ""
```

```yml
labels:
  - "com.example.description=Database volume"
  - "com.example.department=IT/Ops"
  - "com.example.label-with-empty-value"
```

Compose implementation MUST set `com.docker.compose.project` and `com.docker.compose.volume` labels.

### name

`name` set a custom name for this volume. The name field can be used to reference volumes that contain special
characters. The name is used as is and will **not** be scoped with the stack name.

```yml
volumes:
  data:
    name: "my-app-data"
```

It can also be used in conjunction with the `external` property. Doing so the name of the volume used to lookup for
actual volume on platform is set separately from the name used to refer to it within the Compose file:

```yml
volumes:
  db-data:
    external:
      name: actual-name-of-volume
```

This make it possible to make this lookup name a parameter of a Compose file, so that the model ID for volume is
hard-coded but the actual volume ID on platform is set at runtime during deployment:

```yml
volumes:
  db-data:
    external:
      name: ${DATABASE_VOLUME}
```

## Configs top-level element

Configs allow services to adapt their behaviour without the need to rebuild a Docker image. Configs are comparable to Volumes from a service point of view as they are mounted into service's containers filesystem. The actual implementation detail to get configuration provided by the platform can be set from the Configuration definition.

When granted access to a config, the config content is mounted as a file in the container. The location of the mount point within the container defaults to `/<config-name>` in Linux containers and `C:\<config-name>` in Windows containers.

By default, the config MUST be owned by the user running the container command but can be overridden by service configuration.
By default, the config MUST have world-readable permissions (mode 0444), unless service is configured to override this.

Services can only access configs when explicitly granted by a [`configs`](#configs) subsection.

The top-level `configs` declaration defines or references
configuration data that can be granted to the services in this
application. The source of the config is either `file` or `external`.

- `file`: The config is created with the contents of the file at the specified path.
- `external`: If set to true, specifies that this config has already been created. Compose implementation does not
  attempt to create it, and if it does not exist, an error occurs.
- `name`: The name of config object on Platform to lookup. This field can be used to
  reference configs that contain special characters. The name is used as is
  and will **not** be scoped with the project name.

In this example, `http_config` is created (as `<project_name>_http_config`) when the application is deployed,
and `my_second_config` MUST already exist on Platform and value will be obtained by lookup.

In this example, `server-http_config` is created as `<project_name>_http_config` when the application is deployed,
by registering content of the `httpd.conf` as configuration data.

```yml
configs:
  http_config:
    file: ./httpd.conf
```

Alternatively, `http_config` can be declared as external, doing so Compose implementation will lookup `http_config` to expose configuration data to relevant services.

```yml
configs:
  http_config:
    external: true
```

External configs lookup can also use a distinct key by specifying a `name`. The following
example modifies the previous one to lookup for config using a parameter `HTTP_CONFIG_KEY`. Doing
so the actual lookup key will be set at deployment time by [interpolation](#interpolation) of
variables, but exposed to containers as hard-coded ID `http_config`.

```yml
configs:
  http_config:
    external: true
    name: "${HTTP_CONFIG_KEY}"
```

Compose file need to explicitly grant access to the configs to relevant services in the application.

## Secrets top-level element

Secrets are a flavour of Configs focussing on sensitive data, with specific constraint for this usage. As the platform implementation may significantly differ from Configs, dedicated Secrets section allows to configure the related resources.

The top-level `secrets` declaration defines or references sensitive data that can be granted to the services in this
application. The source of the secret is either `file` or `external`.

- `file`: The secret is created with the contents of the file at the specified path.
- `external`: If set to true, specifies that this secret has already been created. Compose implementation does
  not attempt to create it, and if it does not exist, an error occurs.
- `name`: The name of the secret object in Docker. This field can be used to
  reference secrets that contain special characters. The name is used as is
  and will **not** be scoped with the project name.

In this example, `server-certificate` is created as `<project_name>_server-certificate` when the application is deployed,
by registering content of the `server.cert` as a platform secret.

```yml
secrets:
  server-certificate:
    file: ./server.cert
```

Alternatively, `server-certificate` can be declared as external, doing so Compose implementation will lookup `server-certificate` to expose secret to relevant services.

```yml
secrets:
  server-certificate:
    external: true
```

External secrets lookup can also use a distinct key by specifying a `name`. The following
example modifies the previous one to look up for secret using a parameter `CERTIFICATE_KEY`. Doing
so the actual lookup key will be set at deployment time by [interpolation](#interpolation) of
variables, but exposed to containers as hard-coded ID `server-certificate`.

```yml
secrets:
  server-certificate:
    external: true
    name: "${CERTIFICATE_KEY}"
```

Compose file need to explicitly grant access to the secrets to relevant services in the application.

## Fragments

It is possible to re-use configuration fragments using [YAML anchors](http://www.yaml.org/spec/1.2/spec.html#id2765878).

```yml
volumes:
  db-data: &default-volume
    driver: default
  metrics: *default-volume
```

In previous sample, an _anchor_ is created as `default-volume` based on `db-data` volume specification. It is later reused by _alias_ `*default-volume` to define `metrics` volume. Same logic can apply to any element in a Compose file. Anchor resolution MUST take place
before [variables interpolation](#interpolation), so variables can't be used to set anchors or aliases.

It is also possible to partially override values set by anchor reference using the
[YAML merge type](http://yaml.org/type/merge.html). In following example, `metrics` volume specification uses alias
to avoid repetition but override `name` attribute:

```yml

services:
  backend:
    image: awesome/database
    volumes:
      - db-data
      - metrics
volumes:
  db-data: &default-volume
    driver: default
    name: "data"
  metrics:
    <<: *default-volume
    name: "metrics"
```

## Extension

Special extension fields can be of any format as long as their name starts with the `x-` character sequence. They can be used
within any structure in a Compose file. This is the sole exception for Compose implementations to silently ignore unrecognized field.

```yml
x-custom:
  foo:
    - bar
    - zot

services:
  webapp:
    image: awesome/webapp
    x-foo: bar
```

The contents of such fields are unspecified by Compose specification, and can be used to enable custom features. Compose implementation to encounter an unknown extension field MUST NOT fail, but COULD warn about unknown field.

For platform extensions, it is highly recommended to prefix extension by platform/vendor name, the same way browsers add
support for [custom CSS features](https://www.w3.org/TR/2011/REC-CSS2-20110607/syndata.html#vendor-keywords){: target="_blank" rel="noopener" class="_"}.

```yml
service:
  backend:
    deploy:
      placement:
        x-aws-role: "arn:aws:iam::XXXXXXXXXXXX:role/foo"
        x-aws-region: "eu-west-3"
        x-azure-region: "france-central"
```

### Informative Historical Notes

This section is informative. At the time of writing, the following prefixes are known to exist:

| prefix     | vendor/organization |
| ---------- | ------------------- |
| docker     | Docker              |
| kubernetes | Kubernetes          |

### Using extensions as fragments

With the support for extension fields, Compose file can be written as follows to improve readability of reused fragments:

```yml
x-logging: &default-logging
  options:
    max-size: "12m"
    max-file: "5"
  driver: json-file

services:
  frontend:
    image: awesome/webapp
    logging: *default-logging
  backend:
    image: awesome/database
    logging: *default-logging
```

### specifying byte values

Value express a byte value as a string in `{amount}{byte unit}` format:
The supported units are `b` (bytes), `k` or `kb` (kilo bytes), `m` or `mb` (mega bytes) and `g` or `gb` (giga bytes).

```
    2b
    1024kb
    2048k
    300m
    1gb
```

### specifying durations

Value express a duration as a string in the in the form of `{value}{unit}`.
The supported units are `us` (microseconds), `ms` (milliseconds), `s` (seconds), `m` (minutes) and `h` (hours).
Value can can combine multiple values and using without separator.

```
  10ms
  40s
  1m30s
  1h5m30s20ms
```

## Interpolation

Values in a Compose file can be set by variables, and interpolated at runtime. Compose files use a Bash-like
syntax `${VARIABLE}`

Both `$VARIABLE` and `${VARIABLE}` syntax are supported. Default values can be defined inline using typical shell syntax:
latest

- `${VARIABLE:-default}` evaluates to `default` if `VARIABLE` is unset or
  empty in the environment.
- `${VARIABLE-default}` evaluates to `default` only if `VARIABLE` is unset
  in the environment.

Similarly, the following syntax allows you to specify mandatory variables:

- `${VARIABLE:?err}` exits with an error message containing `err` if
  `VARIABLE` is unset or empty in the environment.
- `${VARIABLE?err}` exits with an error message containing `err` if
  `VARIABLE` is unset in the environment.

Interpolation can also be nested:

- `${VARIABLE:-${FOO}}`
- `${VARIABLE?$FOO}`
- `${VARIABLE:-${FOO:-default}}`

Other extended shell-style features, such as `${VARIABLE/foo/bar}`, are not
supported by the Compose specification.

You can use a `$$` (double-dollar sign) when your configuration needs a literal
dollar sign. This also prevents Compose from interpolating a value, so a `$$`
allows you to refer to environment variables that you don't want processed by
Compose.

```yml
web:
  build: .
  command: "$$VAR_NOT_INTERPOLATED_BY_COMPOSE"
```

If the Compose implementation can't resolve a substituted variable and no default value is defined, it MUST warn
the user and substitute the variable with an empty string.

As any values in a Compose file can be interpolated with variable substitution, including compact string notation
for complex elements, interpolation MUST be applied _before_ merge on a per-file-basis.

## Compose documentation

- [User guide](../index.md)
- [Installing Compose](../install.md)
- [Compose file versions and upgrading](compose-versioning.md)
- [Sample apps with Compose](../samples-for-compose.md)
- [Enabling GPU access with Compose](../gpu-support.md)
- [Command line reference](../reference/index.md)
