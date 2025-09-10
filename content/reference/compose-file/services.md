---
linkTitle: Services
title: Define services in Docker Compose
description: Explore all the attributes the services top-level element can have.
keywords: compose, compose specification, services, compose file reference
aliases:
 - /compose/compose-file/05-services/
weight: 20
---

{{% include "compose/services.md" %}}

A Compose file must declare a `services` top-level element as a map whose keys are string representations of service names,
and whose values are service definitions. A service definition contains the configuration that is applied to each
service container.

Each service may also include a `build` section, which defines how to create the Docker image for the service.
Compose supports building Docker images using this service definition. If not used, the `build` section is ignored and the Compose file is still considered valid. Build support is an optional aspect of the Compose Specification, and is
described in detail in the [Compose Build Specification](build.md) documentation.

Each service defines runtime constraints and requirements to run its containers. The `deploy` section groups
these constraints and lets the platform adjust the deployment strategy to best match containers' needs with
available resources. Deploy support is an optional aspect of the Compose Specification, and is
described in detail in the [Compose Deploy Specification](deploy.md) documentation.
If not implemented the `deploy` section is ignored and the Compose file is still considered valid.

## Examples

### Simple example

The following example demonstrates how to define two simple services, set their images, map ports, and configure basic environment variables using Docker Compose.

```yaml
services:
  web:
    image: nginx:latest
    ports:
      - "8080:80"

  db:
    image: postgres:13
    environment:
      POSTGRES_USER: example
      POSTGRES_DB: exampledb
```

### Advanced example

In the following example, the `proxy` service uses the Nginx image, mounts a local Nginx configuration file into the container, exposes port `80` and depends on the `backend` service.

The `backend` service builds an image from the Dockerfile located in the `backend` directory that is set to build at stage `builder`.

```yaml
services:
  proxy:
    image: nginx
    volumes:
      - type: bind
        source: ./proxy/nginx.conf
        target: /etc/nginx/conf.d/default.conf
        read_only: true
    ports:
      - 80:80
    depends_on:
      - backend

  backend:
    build:
      context: backend
      target: builder
```

For more example Compose files, explore the [Awesome Compose samples](https://github.com/docker/awesome-compose).

## Attributes

<!-- vale off(Docker.HeadingSentenceCase.yml) -->

### `annotations`

`annotations` defines annotations for the container. `annotations` can use either an array or a map.

```yml
annotations:
  com.example.foo: bar
```

```yml
annotations:
  - com.example.foo=bar
```

### `attach`

{{< summary-bar feature_name="Compose attach" >}}

When `attach` is defined and set to `false` Compose does not collect service logs,
until you explicitly request it to.

The default service configuration is `attach: true`.

### `build`

`build` specifies the build configuration for creating a container image from source, as defined in the [Compose Build Specification](build.md).

### `blkio_config`

`blkio_config` defines a set of configuration options to set block I/O limits for a service.

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

#### `device_read_bps`, `device_write_bps`

Set a limit in bytes per second for read / write operations on a given device.
Each item in the list must have two keys:

- `path`: Defines the symbolic path to the affected device.
- `rate`: Either as an integer value representing the number of bytes or as a string expressing a byte value.

#### `device_read_iops`, `device_write_iops`

Set a limit in operations per second for read / write operations on a given device.
Each item in the list must have two keys:

- `path`: Defines the symbolic path to the affected device.
- `rate`: As an integer value representing the permitted number of operations per second.

#### `weight`

Modify the proportion of bandwidth allocated to a service relative to other services.
Takes an integer value between 10 and 1000, with 500 being the default.

#### `weight_device`

Fine-tune bandwidth allocation by device. Each item in the list must have two keys:

- `path`: Defines the symbolic path to the affected device.
- `weight`: An integer value between 10 and 1000.

### `cpu_count`

`cpu_count` defines the number of usable CPUs for service container.

### `cpu_percent`

`cpu_percent` defines the usable percentage of the available CPUs.

### `cpu_shares`

`cpu_shares` defines, as integer value, a service container's relative CPU weight versus other containers.

### `cpu_period`

`cpu_period` configures CPU CFS (Completely Fair Scheduler) period when a platform is based
on Linux kernel.

### `cpu_quota`

`cpu_quota` configures CPU CFS (Completely Fair Scheduler) quota when a platform is based
on Linux kernel.

### `cpu_rt_runtime`

`cpu_rt_runtime` configures CPU allocation parameters for platforms with support for real-time scheduler. It can be either
an integer value using microseconds as unit or a [duration](extension.md#specifying-durations).

```yml
 cpu_rt_runtime: '400ms'
 cpu_rt_runtime: '95000'
```

### `cpu_rt_period`

`cpu_rt_period` configures CPU allocation parameters for platforms with support for real-time scheduler. It can be either
an integer value using microseconds as unit or a [duration](extension.md#specifying-durations).

```yml
 cpu_rt_period: '1400us'
 cpu_rt_period: '11000'
```

### `cpus`

`cpus` define the number of (potentially virtual) CPUs to allocate to service containers. This is a fractional number.
`0.000` means no limit.

When set, `cpus` must be consistent with the `cpus` attribute in the [Deploy Specification](deploy.md#cpus).

### `cpuset`

`cpuset` defines the explicit CPUs in which to permit execution. Can be a range `0-3` or a list `0,1`

### `cap_add`

`cap_add` specifies additional container [capabilities](https://man7.org/linux/man-pages/man7/capabilities.7.html)
as strings.

```yaml
cap_add:
  - ALL
```

### `cap_drop`

`cap_drop` specifies container [capabilities](https://man7.org/linux/man-pages/man7/capabilities.7.html) to drop
as strings.

```yaml
cap_drop:
  - NET_ADMIN
  - SYS_ADMIN
```

### `cgroup`

{{< summary-bar feature_name="Compose cgroup" >}}

`cgroup` specifies the cgroup namespace to join. When unset, it is the container runtime's decision to
select which cgroup namespace to use, if supported.

- `host`: Runs the container in the Container runtime cgroup namespace.
- `private`: Runs the container in its own private cgroup namespace.

### `cgroup_parent`

`cgroup_parent` specifies an optional parent [cgroup](https://man7.org/linux/man-pages/man7/cgroups.7.html) for the container.

```yaml
cgroup_parent: m-executor-abcd
```

### `command`

`command` overrides the default command declared by the container image, for example by Dockerfile's `CMD`.

```yaml
command: bundle exec thin -p 3000
```

If the value is `null`, the default command from the image is used.

If the value is `[]` (empty list) or `''` (empty string), the default command declared by the image is ignored, or in other words overridden to be empty.

> [!NOTE]
>
> Unlike the `CMD` instruction in a Dockerfile, the `command` field doesn't automatically run within the context of the [`SHELL`](/reference/dockerfile.md#shell-form) instruction defined in the image. If your `command` relies on shell-specific features, such as environment variable expansion, you need to explicitly run it within a shell. For example:
>
> ```yaml
> command: /bin/sh -c 'echo "hello $$HOSTNAME"'
> ```

The value can also be a list, similar to the [exec-form syntax](/reference/dockerfile.md#exec-form)
used by the [Dockerfile](/reference/dockerfile.md#exec-form).

### `configs`

`configs` let services adapt their behaviour without the need to rebuild a Docker image.
Services can only access configs when explicitly granted by the `configs` attribute. Two different syntax variants are supported.

Compose reports an error if `config` doesn't exist on the platform or isn't defined in the
[`configs` top-level element](configs.md) in the Compose file.

There are two syntaxes defined for configs: a short syntax and a long syntax.

You can grant a service access to multiple configs, and you can mix long and short syntax.

#### Short syntax

The short syntax variant only specifies the config name. This grants the
container access to the config and mounts it as files into a service’s container’s filesystem. The location of the mount point within the container defaults to `/<config_name>` in Linux containers, and `C:\<config-name>` in Windows containers.

The following example uses the short syntax to grant the `redis` service
access to the `my_config` and `my_other_config` configs. The value of
`my_config` is set to the contents of the file `./my_config.txt`, and
`my_other_config` is defined as an external resource, which means that it has
already been defined in the platform. If the external config does not exist,
the deployment fails.

```yml
services:
  redis:
    image: redis:latest
    configs:
      - my_config
      - my_other_config
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
- `uid` and `gid`: The numeric uid or gid that owns the mounted config file
  within the service's task containers. Default value when not specified is `USER`.
- `mode`: The [permissions](https://wintelguy.com/permissions-calc.pl) for the file that is mounted within the service's
  task containers, in octal notation. Default value is world-readable (`0444`).
  Writable bit must be ignored. The executable bit can be set.

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

### `container_name`

`container_name` is a string that specifies a custom container name, rather than a name generated by default.

```yml
container_name: my-web-container
```

Compose does not scale a service beyond one container if the Compose file specifies a
`container_name`. Attempting to do so results in an error.

`container_name` follows the regex format of `[a-zA-Z0-9][a-zA-Z0-9_.-]+`

### `credential_spec`

`credential_spec` configures the credential spec for a managed service account.

If you have services that use Windows containers, you can use `file:` and
`registry:` protocols for `credential_spec`. Compose also supports additional
protocols for custom use-cases.

The `credential_spec` must be in the format `file://<filename>` or `registry://<value-name>`.

```yml
credential_spec:
  file: my-credential-spec.json
```

When using `registry:`, the credential spec is read from the Windows registry on
the daemon's host. A registry value with the given name must be located in:

```bash
HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization\Containers\CredentialSpecs
```

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
    file: ./my-credential-spec.json
```

### `depends_on`

{{% include "compose/services-depends-on.md" %}}

#### Short syntax

The short syntax variant only specifies service names of the dependencies.
Service dependencies cause the following behaviors:

- Compose creates services in dependency order. In the following
  example, `db` and `redis` are created before `web`.

- Compose removes services in dependency order. In the following
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

Compose guarantees dependency services have been started before
starting a dependent service.
Compose waits for dependency services to be "ready" before
starting a dependent service.

#### Long syntax

The long form syntax enables the configuration of additional fields that can't be
expressed in the short form.

- `restart`: When set to `true` Compose restarts this service after it updates the dependency service.
  This applies to an explicit restart controlled by a Compose operation, and excludes automated restart by the container runtime
  after the container dies. Introduced in Docker Compose version [2.17.0](/manuals/compose/releases/release-notes.md#2170).

- `condition`: Sets the condition under which dependency is considered satisfied
  - `service_started`: An equivalent of the short syntax described previously
  - `service_healthy`: Specifies that a dependency is expected to be "healthy"
    (as indicated by [`healthcheck`](#healthcheck)) before starting a dependent
    service.
  - `service_completed_successfully`: Specifies that a dependency is expected to run
    to successful completion before starting a dependent service.
- `required`: When set to `false` Compose only warns you when the dependency service isn't started or available. If it's not defined
    the default value of `required` is `true`. Introduced in Docker Compose version [2.20.0](/manuals/compose/releases/release-notes.md#2200).

Service dependencies cause the following behaviors:

- Compose creates services in dependency order. In the following
  example, `db` and `redis` are created before `web`.

- Compose waits for healthchecks to pass on dependencies
  marked with `service_healthy`. In the following example, `db` is expected to
  be "healthy" before `web` is created.

- Compose removes services in dependency order. In the following
  example, `web` is removed before `db` and `redis`.

```yml
services:
  web:
    build: .
    depends_on:
      db:
        condition: service_healthy
        restart: true
      redis:
        condition: service_started
  redis:
    image: redis
  db:
    image: postgres
```

Compose guarantees dependency services are started before
starting a dependent service.
Compose guarantees dependency services marked with
`service_healthy` are "healthy" before starting a dependent service.

### `deploy`

`deploy` specifies the configuration for the deployment and lifecycle of services, as defined [in the Compose Deploy Specification](deploy.md).

### `develop`

{{< summary-bar feature_name="Compose develop" >}}

`develop` specifies the development configuration for maintaining a container in sync with source, as defined in the [Development Section](develop.md).

### `device_cgroup_rules`

`device_cgroup_rules` defines a list of device cgroup rules for this container.
The format is the same format the Linux kernel specifies in the [Control Groups
Device Whitelist Controller](https://www.kernel.org/doc/html/latest/admin-guide/cgroup-v1/devices.html).

```yml
device_cgroup_rules:
  - 'c 1:3 mr'
  - 'a 7:* rmw'
```

### `devices`

`devices` defines a list of device mappings for created containers in the form of
`HOST_PATH:CONTAINER_PATH[:CGROUP_PERMISSIONS]`.

```yml
devices:
  - "/dev/ttyUSB0:/dev/ttyUSB0"
  - "/dev/sda:/dev/xvda:rwm"
```

`devices` can also rely on the [CDI](https://github.com/cncf-tags/container-device-interface) syntax to let the container runtime select a device:

```yml
devices:
  - "vendor1.com/device=gpu"
```

### `dns`

`dns` defines custom DNS servers to set on the container network interface configuration. It can be a single value or a list.

```yml
dns: 8.8.8.8
```

```yml
dns:
  - 8.8.8.8
  - 9.9.9.9
```

### `dns_opt`

`dns_opt` list custom DNS options to be passed to the container’s DNS resolver (`/etc/resolv.conf` file on Linux).

```yml
dns_opt:
  - use-vc
  - no-tld-query
```

### `dns_search`

`dns_search` defines custom DNS search domains to set on container network interface configuration. It can be a single value or a list.

```yml
dns_search: example.com
```

```yml
dns_search:
  - dc1.example.com
  - dc2.example.com
```

### `domainname`

`domainname` declares a custom domain name to use for the service container. It must be a valid RFC 1123 hostname.

### `driver_opts`

{{< summary-bar feature_name="Compose driver opts" >}}

`driver_opts` specifies a list of options as key-value pairs to pass to the driver. These options are
driver-dependent.

```yml
services:
  app:
    networks:
      app_net:
        driver_opts:
          com.docker.network.bridge.host_binding_ipv4: "127.0.0.1"
```

Consult the [network drivers documentation](/manuals/engine/network/_index.md) for more information.

### `entrypoint`

`entrypoint` declares the default entrypoint for the service container.
This overrides the `ENTRYPOINT` instruction from the service's Dockerfile.

If `entrypoint` is non-null, Compose ignores any default command from the image, for example the `CMD`
instruction in the Dockerfile.

See also [`command`](#command) to set or override the default command to be executed by the entrypoint process.

In its short form, the value can be defined as a string:
```yml
entrypoint: /code/entrypoint.sh
```

Alternatively, the value can also be a list, in a manner similar to the
[Dockerfile](https://docs.docker.com/reference/dockerfile/#cmd):

```yml
entrypoint:
  - php
  - -d
  - zend_extension=/usr/local/lib/php/extensions/no-debug-non-zts-20100525/xdebug.so
  - -d
  - memory_limit=-1
  - vendor/bin/phpunit
```

If the value is `null`, the default entrypoint from the image is used.

If the value is `[]` (empty list) or `''` (empty string), the default entrypoint declared by the image is ignored, or in other words, overridden to be empty.

### `env_file`

{{% include "compose/services-env-file.md" %}}

```yml
env_file: .env
```

Relative paths are resolved from the Compose file's parent folder. As absolute paths prevent the Compose
file from being portable, Compose warns you when such a path is used to set `env_file`.

Environment variables declared in the [`environment`](#environment) section override these values. This holds true even if those values are
empty or undefined.

`env_file` can also be a list. The files in the list are processed from the top down. For the same variable
specified in two environment files, the value from the last file in the list stands.

```yml
env_file:
  - ./a.env
  - ./b.env
```

List elements can also be declared as a mapping, which then lets you set additional
attributes.

#### `required`

{{< summary-bar feature_name="Compose required" >}}

The `required` attribute defaults to `true`. When `required` is set to `false` and the `.env` file is missing, Compose silently ignores the entry.

```yml
env_file:
  - path: ./default.env
    required: true # default
  - path: ./override.env
    required: false
```

#### `format`

{{< summary-bar feature_name="Compose format" >}}

The `format` attribute lets you use an alternative file format for the `env_file`. When not set, `env_file` is parsed according to the Compose rules outlined in [`Env_file` format](#env_file-format).

`raw` format lets you use an `env_file` with key=value items, but without any attempt from Compose to parse the value for interpolation.
This let you pass values as-is, including quotes and `$` signs.

```yml
env_file:
  - path: ./default.env
    format: raw
```

#### `Env_file` format

Each line in an `.env` file must be in `VAR[=[VAL]]` format. The following syntax rules apply:

- Lines beginning with `#` are processed as comments and ignored.
- Blank lines are ignored.
- Unquoted and double-quoted (`"`) values have [Interpolation](interpolation.md) applied.
- Each line represents a key-value pair. Values can optionally be quoted.
  - `VAR=VAL` -> `VAL`
  - `VAR="VAL"` -> `VAL`
  - `VAR='VAL'` -> `VAL`
- Inline comments for unquoted values must be preceded with a space.
  - `VAR=VAL # comment` -> `VAL`
  - `VAR=VAL# not a comment` -> `VAL# not a comment`
- Inline comments for quoted values must follow the closing quote.
  - `VAR="VAL # not a comment"` -> `VAL # not a comment`
  - `VAR="VAL" # comment` -> `VAL`
- Single-quoted (`'`) values are used literally.
  - `VAR='$OTHER'` -> `$OTHER`
  - `VAR='${OTHER}'` -> `${OTHER}`
- Quotes can be escaped with `\`.
  - `VAR='Let\'s go!'` -> `Let's go!`
  - `VAR="{\"hello\": \"json\"}"` -> `{"hello": "json"}`
- Common shell escape sequences including `\n`, `\r`, `\t`, and `\\` are supported in double-quoted values.
  - `VAR="some\tvalue"` -> `some  value`
  - `VAR='some\tvalue'` -> `some\tvalue`
  - `VAR=some\tvalue` -> `some\tvalue`

`VAL` may be omitted, in such cases the variable value is an empty string.
`=VAL` may be omitted, in such cases the variable is unset.

```bash
# Set Rails/Rack environment
RACK_ENV=development
VAR="quoted"
```

### `environment`

{{% include "compose/services-environment.md" %}}

Environment variables can be declared by a single key (no value to equals sign). In this case Compose
relies on you to resolve the value. If the value is not resolved, the variable
is unset and is removed from the service container environment.

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

### `expose`

`expose` defines the (incoming) port or a range of ports that Compose exposes from the container. These ports must be
accessible to linked services and should not be published to the host machine. Only the internal container
ports can be specified.

Syntax is `<portnum>/[<proto>]` or `<startport-endport>/[<proto>]` for a port range.
When not explicitly set, `tcp` protocol is used.

```yml
expose:
  - "3000"
  - "8000"
  - "8080-8085/tcp"
```

> [!NOTE]
>
> If the Dockerfile for the image already exposes ports, it is visible to other containers on the network even if `expose` is not set in your Compose file.

### `extends`

`extends` lets you share common configurations among different files, or even different projects entirely. With `extends` you can define a common set of service options in one place and refer to it from anywhere. You can refer to another Compose file and select a service you want to also use in your own application, with the ability to override some attributes for your own needs.

You can use `extends` on any service together with other configuration keys. The `extends` value must be a mapping
defined with a required `service` and an optional `file` key.

```yaml
extends:
  file: common.yml
  service: webapp
```

- `service`: Defines the name of the service being referenced as a base, for example `web` or `database`.
- `file`: The location of a Compose configuration file defining that service.

#### Restrictions

When a service is referenced using `extends`, it can declare dependencies on other resources. These dependencies may be explicitly defined through attributes like `volumes`, `networks`, `configs`, `secrets`, `links`, `volumes_from`, or `depends_on`. Alternatively, dependencies can reference another service using the `service:{name}` syntax in namespace declarations such as `ipc`, `pid`, or `network_mode`.

Compose does not automatically import these referenced resources into the extended model. It is your responsibility to ensure all required resources are explicitly declared in the model that relies on extends.

Circular references with `extends` are not supported, Compose returns an error when one is detected.

#### Finding referenced service

`file` value can be:

- Not present.
  This indicates that another service within the same Compose file is being referenced.
- File path, which can be either:
  - Relative path. This path is considered as relative to the location of the main Compose
    file.
  - Absolute path.

A service denoted by `service` must be present in the identified referenced Compose file.
Compose returns an error if:

- The service denoted by `service` is not found.
- The Compose file denoted by `file` is not found.

#### Merging service definitions

Two service definitions, the main one in the current Compose file and the referenced one
specified by `extends`, are merged in the following way:

- Mappings: Keys in mappings of the main service definition override keys in mappings
  of the referenced service definition. Keys that aren't overridden are included as is.
- Sequences: Items are combined together into a new sequence. The order of elements is
  preserved with the referenced items coming first and main items after.
- Scalars: Keys in the main service definition take precedence over keys in the
  referenced one.

##### Mappings

The following keys should be treated as mappings: `annotations`, `build.args`, `build.labels`,
`build.extra_hosts`, `deploy.labels`, `deploy.update_config`, `deploy.rollback_config`,
`deploy.restart_policy`, `deploy.resources.limits`, `environment`, `healthcheck`,
`labels`, `logging.options`, `sysctls`, `storage_opt`, `extra_hosts`, `ulimits`.

One exception that applies to `healthcheck` is that the main mapping cannot specify
`disable: true` unless the referenced mapping also specifies `disable: true`. Compose returns an error in this case.
For example, the following input:

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

For example, the following input:

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

Produces the following configuration for the `cli` service. Note that the mounted path
now points to the new volume name and `ro` flag was applied.

```yaml
image: busybox
volumes:
- cli-volume:/var/lib/backup/data:ro
```

If the referenced service definition contains `extends` mapping, the items under it
are simply copied into the new merged definition. The merging process is then kicked
off again until no `extends` keys are remaining.

For example, the following input:

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

For example, the following input:

```yaml
services:
  common:
    image: busybox
    security_opt:
      - label=role:ROLE
  cli:
    extends:
      service: common
    security_opt:
      - label=user:USER
```

Produces the following configuration for the `cli` service.

```yaml
image: busybox
security_opt:
- label=role:ROLE
- label=user:USER
```

In case list syntax is used, the following keys should also be treated as sequences:
`dns`, `dns_search`, `env_file`, `tmpfs`. Unlike sequence fields mentioned previously,
duplicates resulting from the merge are not removed.

##### Scalars

Any other allowed keys in the service definition should be treated as scalars.

### `external_links`

`external_links` link service containers to services managed outside of your Compose application.
`external_links` define the name of an existing service to retrieve using the platform lookup mechanism.
An alias of the form `SERVICE:ALIAS` can be specified.

```yml
external_links:
  - redis
  - database:mysql
  - database:postgresql
```

### `extra_hosts`

`extra_hosts` adds hostname mappings to the container network interface configuration (`/etc/hosts` for Linux).

#### Short syntax

Short syntax uses plain strings in a list. Values must set hostname and IP address for additional hosts in the form of `HOSTNAME=IP`.

```yml
extra_hosts:
  - "somehost=162.242.195.82"
  - "otherhost=50.31.209.229"
  - "myhostv6=::1"
```

IPv6 addresses can be enclosed in square brackets, for example:

```yml
extra_hosts:
  - "myhostv6=[::1]"
```

The separator `=` is preferred, but `:` can also be used. Introduced in Docker Compose version [2.24.1](/manuals/compose/releases/release-notes.md#2241). For example:

```yml
extra_hosts:
  - "somehost:162.242.195.82"
  - "myhostv6:::1"
```

#### Long syntax

Alternatively, `extra_hosts` can be set as a mapping between hostname(s) and IP(s)

```yml
extra_hosts:
  somehost: "162.242.195.82"
  otherhost: "50.31.209.229"
  myhostv6: "::1"
```

Compose creates a matching entry with the IP address and hostname in the container's network
configuration, which means for Linux `/etc/hosts` get extra lines:

```console
162.242.195.82  somehost
50.31.209.229   otherhost
::1             myhostv6
```

### `gpus`

{{< summary-bar feature_name="Compose gpus" >}}

`gpus` specifies GPU devices to be allocated for container usage. This is equivalent to a [device request](deploy.md#devices) with
an implicit `gpu` capability.

```yaml
services:
  model:
    gpus:
      - driver: 3dfx
        count: 2
```

`gpus` also can be set as string `all` to allocate all available GPU devices to the container.

```yaml
services:
  model:
    gpus: all
```

### `group_add`

`group_add` specifies additional groups, by name or number, which the user inside the container must be a member of.

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

Running `id` inside the created container must show that the user belongs to the `mail` group, which would not have
been the case if `group_add` were not declared.

### `healthcheck`

{{% include "compose/services-healthcheck.md" %}}

For more information on `HEALTHCHECK`, see the [Dockerfile reference](/reference/dockerfile.md#healthcheck).

```yml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost"]
  interval: 1m30s
  timeout: 10s
  retries: 3
  start_period: 40s
  start_interval: 5s
```

`interval`, `timeout`, `start_period`, and `start_interval` are [specified as durations](extension.md#specifying-durations). Introduced in Docker Compose version [2.20.2](/manuals/compose/releases/release-notes.md#2202)

`test` defines the command Compose runs to check container health. It can be
either a string or a list. If it's a list, the first item must be either `NONE`, `CMD` or `CMD-SHELL`.
If it's a string, it's equivalent to specifying `CMD-SHELL` followed by that string.

```yml
# Hit the local web app
test: ["CMD", "curl", "-f", "http://localhost"]
```

Using `CMD-SHELL` runs the command configured as a string using the container's default shell
(`/bin/sh` for Linux). Both of the following forms are equivalent:

```yml
test: ["CMD-SHELL", "curl -f http://localhost || exit 1"]
```

```yml
test: curl -f https://localhost || exit 1
```

`NONE` disables the healthcheck, and is mostly useful to disable the Healthcheck Dockerfile instruction set by the service's Docker image. Alternatively,
the healthcheck set by the image can be disabled by setting `disable: true`:

```yml
healthcheck:
  disable: true
```

### `hostname`

`hostname` declares a custom host name to use for the service container. It must be a valid RFC 1123 hostname.

### `image`

`image` specifies the image to start the container from. `image` must follow the Open Container Specification
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

If the image does not exist on the platform, Compose attempts to pull it based on the `pull_policy`.
If you are also using the [Compose Build Specification](build.md), there are alternative options for controlling the precedence of
pull over building the image from source, however pulling the image is the default behavior.

`image` may be omitted from a Compose file as long as a `build` section is declared. If you are not using the Compose Build Specification, Compose won't work if `image` is missing from the Compose file.

### `init`

`init` runs an init process (PID 1) inside the container that forwards signals and reaps processes.
Set this option to `true` to enable this feature for the service.

```yml
services:
  web:
    image: alpine:latest
    init: true
```

The init binary that is used is platform specific.

### `ipc`

`ipc` configures the IPC isolation mode set by the service container.

- `shareable`: Gives the container its own private IPC namespace, with a
  possibility to share it with other containers.
- `service:{name}`: Makes the container join another container's
  (`shareable`) IPC namespace.

```yml
    ipc: "shareable"
    ipc: "service:[service name]"
```

### `isolation`

`isolation` specifies a container’s isolation technology. Supported values are platform specific.

### `labels`

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

Compose creates containers with canonical labels:

- `com.docker.compose.project` set on all resources created by Compose to the user project name
- `com.docker.compose.service` set on service containers with service name as defined in the Compose file

The `com.docker.compose` label prefix is reserved. Specifying labels with this prefix in the Compose file
results in a runtime error.

### `label_file`

{{< summary-bar feature_name="Compose label file" >}}

The `label_file` attribute lets you load labels for a service from an external file or a list of files. This provides a convenient way to manage multiple labels without cluttering the Compose file.

The file uses a key-value format, similar to `env_file`. You can specify multiple files as a list. When using multiple files, they are processed in the order they appear in the list. If the same label is defined in multiple files, the value from the last file in the list overrides earlier ones.

```yaml
services:
  one:
    label_file: ./app.labels

  two:
    label_file:
      - ./app.labels
      - ./additional.labels
```

If a label is defined in both the `label_file` and the `labels` attribute, the value in [labels](#labels) takes precedence.

### `links`

`links` defines a network link to containers in another service. Either specify both the service name and
a link alias (`SERVICE:ALIAS`), or just the service name.

```yml
web:
  links:
    - db
    - db:database
    - redis
```

Containers for the linked service are reachable at a hostname identical to the alias, or the service name
if no alias is specified.

Links are not required to enable services to communicate. When no specific network configuration is set,
any service is able to reach any other service at that service’s name on the `default` network.
If services specify the networks they are attached to, `links` does not override the network configuration. Services that are not connected to a shared network are not be able to communicate with each other. Compose doesn't warn you about a configuration mismatch.

Links also express implicit dependency between services in the same way as
[`depends_on`](#depends_on), so they determine the order of service startup.

### `logging`

`logging` defines the logging configuration for the service.

```yml
logging:
  driver: syslog
  options:
    syslog-address: "tcp://192.168.0.42:123"
```

The `driver` name specifies a logging driver for the service's containers. The default and available values
are platform specific. Driver specific options can be set with `options` as key-value pairs.

### `mac_address`

> Available with Docker Compose version 2.24.0 and later.

`mac_address` sets a Mac address for the service container.

> [!NOTE]
> Container runtimes might reject this value, for example Docker Engine >= v25.0. In that case, you should use [networks.mac_address](#mac_address) instead.

### `mem_limit`

`mem_limit` configures a limit on the amount of memory a container can allocate, set as a string expressing a [byte value](extension.md#specifying-byte-values).

When set, `mem_limit` must be consistent with the `limits.memory` attribute in the [Deploy Specification](deploy.md#memory).

### `mem_reservation`

`mem_reservation` configures a reservation on the amount of memory a container can allocate, set as a string expressing a [byte value](extension.md#specifying-byte-values).

When set, `mem_reservation` must be consistent with the `reservations.memory` attribute in the [Deploy Specification](deploy.md#memory).

### `mem_swappiness`

`mem_swappiness` defines as a percentage, a value between 0 and 100, for the host kernel to swap out
anonymous memory pages used by a container.

- `0`: Turns off anonymous page swapping.
- `100`: Sets all anonymous pages as swappable.

The default value is platform specific.

### `memswap_limit`

`memswap_limit` defines the amount of memory the container is allowed to swap to disk. This is a modifier
attribute that only has meaning if [`memory`](deploy.md#memory) is also set. Using swap lets the container write excess
memory requirements to disk when the container has exhausted all the memory that is available to it.
There is a performance penalty for applications that swap memory to disk often.

- If `memswap_limit` is set to a positive integer, then both `memory` and `memswap_limit` must be set. `memswap_limit` represents the total amount of memory and swap that can be used, and `memory` controls the amount used by non-swap memory. So if `memory`="300m" and `memswap_limit`="1g", the container can use 300m of memory and 700m (1g - 300m) swap.
- If `memswap_limit` is set to 0, the setting is ignored, and the value is treated as unset.
- If `memswap_limit` is set to the same value as `memory`, and `memory` is set to a positive integer, the container does not have access to swap.
- If `memswap_limit` is unset, and `memory` is set, the container can use as much swap as the `memory` setting, if the host container has swap memory configured. For instance, if `memory`="300m" and `memswap_limit` is not set, the container can use 600m in total of memory and swap.
- If `memswap_limit` is explicitly set to -1, the container is allowed to use unlimited swap, up to the amount available on the host system.

### `models`

{{< summary-bar feature_name="Compose models" >}}

`models` defines which AI models the service should use at runtime. Each referenced model must be defined under the [`models` top-level element](models.md).

```yaml
services:
  short_syntax:
    image: app
    models:
      - my_model
  long_syntax:
    image: app
    models:
      my_model:
        endpoint_var: MODEL_URL
        model_var: MODEL
```

When a service is linked to a model, Docker Compose injects environment variables to pass connection details and model identifiers to the container. This allows the application to locate and communicate with the model dynamically at runtime, without hard-coding values.

#### Long syntax

The long syntax gives you more control over the environment variable names.

- `endpoint_var` sets the name of the environment variable that holds the model runner’s URL.
- `model_var` sets the name of the environment variable that holds the model identifier.

If either is omitted, Compose automatically generates the environment variable names based on the model key using the following rules:

 - Convert the model key to uppercase
 - Replace any '-' characters with '_'
 - Append `_URL` for the endpoint variable

### `network_mode`

`network_mode` sets a service container's network mode.

- `none`: Turns off all container networking.
- `host`: Gives the container raw access to the host's network interface.
- `service:{name}`: Gives the container access to the specified container by referring to its service name.
- `container:{name}`: Gives the container access to the specified container by referring to its container ID.

For more information container networks, see the [Docker Engine documentation](/manuals/engine/network/_index.md#container-networks).

```yml
    network_mode: "host"
    network_mode: "none"
    network_mode: "service:[service name]"
```

When set, the [`networks`](#networks) attribute is not allowed and Compose rejects any
Compose file containing both attributes.

### `networks`

{{% include "compose/services-networks.md" %}}

```yml
services:
  some-service:
    networks:
      - some-network
      - other-network
```
For more information about the `networks` top-level element, see [Networks](networks.md).

#### Implicit default network

If `networks` is empty or absent from the Compose file, Compose considers an implicit definition for the service to be
connected to the `default` network:

```yml
services:
  some-service:
    image: foo
```
This example is actually equivalent to:

```yml
services:
  some-service:
    image: foo
    networks:
      default: {}
```

If you want the service to not be connected a network, you must set [`network_mode: none`](#network_mode).

#### `aliases`

`aliases` declares alternative hostnames for the service on the network. Other containers on the same
network can use either the service name or an alias to connect to one of the service's containers.

Since `aliases` are network-scoped, the same service can have different aliases on different networks.

> [!NOTE]
> A network-wide alias can be shared by multiple containers, and even by multiple services.
> If it is, then exactly which container the name resolves to is not guaranteed.

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

In the following example, service `frontend` is able to reach the `backend` service at
the hostname `backend` or `database` on the `back-tier` network. The service `monitoring`
is able to reach same `backend` service at `backend` or `mysql` on the `admin` network.

```yml
services:
  frontend:
    image: example/webapp
    networks:
      - front-tier
      - back-tier

  monitoring:
    image: example/monitoring
    networks:
      - admin

  backend:
    image: example/backend
    networks:
      back-tier:
        aliases:
          - database
      admin:
        aliases:
          - mysql

networks:
  front-tier: {}
  back-tier: {}
  admin: {}
```

#### `interface_name`

{{< summary-bar feature_name="Compose interface-name" >}}

`interface_name` lets you specify the name of the network interface used to connect a service to a given network. This ensures consistent and predictable interface naming across services and networks.

```yaml
services:
  backend:
    image: alpine
    command: ip link show
    networks:
      back-tier:
        interface_name: eth0
```

Running the example Compose application shows:

```console
backend-1  | 11: eth0@if64: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
```

#### `ipv4_address`, `ipv6_address`

Specify a static IP address for a service container when joining the network.

The corresponding network configuration in the [top-level networks section](networks.md) must have an
`ipam` attribute with subnet configurations covering each static address.

```yml
services:
  frontend:
    image: example/webapp
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

#### `link_local_ips`

`link_local_ips` specifies a list of link-local IPs. Link-local IPs are special IPs which belong to a well
known subnet and are purely managed by the operator, usually dependent on the architecture where they are
deployed.

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

#### `mac_address`

{{< summary-bar feature_name="Compose mac address" >}}

`mac_address` sets the Mac address used by the service container when connecting to this particular network.

#### `driver_opts`

`driver_opts` specifies a list of options as key-value pairs to pass to the driver. These options are
driver-dependent. Consult the driver's documentation for more information.

```yml
services:
  app:
    networks:
      app_net:
        driver_opts:
          foo: "bar"
          baz: 1
```

#### `gw_priority`

{{< summary-bar feature_name="Compose gw priority" >}}

The network with the highest `gw_priority` is selected as the default gateway for the service container.
If unspecified, the default value is 0.

In the following example, `app_net_2` will be selected as the default gateway.

```yaml
services:
  app:
    image: busybox
    command: top
    networks:
      app_net_1:
      app_net_2:
        gw_priority: 1
      app_net_3:
networks:
  app_net_1:
  app_net_2:
  app_net_3:
```

#### `priority`

`priority` indicates in which order Compose connects the service’s containers to its
networks. If unspecified, the default value is 0.

If the container runtime accepts a `mac_address` attribute at service level, it is
applied to the network with the highest `priority`. In other cases, use attribute
`networks.mac_address`.

`priority` does not affect which network is selected as the default gateway. Use the
[`gw_priority`](#gw_priority) attribute instead.

`priority` does not control the order in which networks connections are added to
the container, it cannot be used to determine the device name (`eth0` etc.) in the
container.

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

### `oom_kill_disable`

If `oom_kill_disable` is set, Compose configures the platform so it won't kill the container in case
of memory starvation.

### `oom_score_adj`

`oom_score_adj` tunes the preference for containers to be killed by platform in case of memory starvation. Value must
be within -1000,1000 range.

### `pid`

`pid` sets the PID mode for container created by Compose.
Supported values are platform specific.

### `pids_limit`

`pids_limit` tunes a container’s PIDs limit. Set to -1 for unlimited PIDs.

```yml
pids_limit: 10
```

When set, `pids_limit` must be consistent with the `pids` attribute in the [Deploy Specification](deploy.md#pids).

### `platform`

`platform` defines the target platform the containers for the service run on. It uses the `os[/arch[/variant]]` syntax.

The values of `os`, `arch`, and `variant` must conform to the convention used by the [OCI Image Spec](https://github.com/opencontainers/image-spec/blob/v1.0.2/image-index.md).

Compose uses this attribute to determine which version of the image is pulled
and/or on which platform the service’s build is performed.

```yml
platform: darwin
platform: windows/amd64
platform: linux/arm64/v8
```

### `ports`

{{% include "compose/services-ports.md" %}}

> [!NOTE]
>
> Port mapping must not be used with `network_mode: host`. Doing so causes a runtime error because `network_mode: host` already exposes container ports directly to the host network, so port mapping isn’t needed.

#### Short syntax

The short syntax is a colon-separated string to set the host IP, host port, and container port
in the form:

`[HOST:]CONTAINER[/PROTOCOL]` where:

- `HOST` is `[IP:](port | range)` (optional). If it is not set, it binds to all network interfaces (`0.0.0.0`).
- `CONTAINER` is `port | range`.
- `PROTOCOL` restricts ports to a specified protocol either `tcp` or `udp`(optional). Default is `tcp`.

Ports can be either a single value or a range. `HOST` and `CONTAINER` must use equivalent ranges.

You can either specify both ports (`HOST:CONTAINER`), or just the container port. In the latter case,
the container runtime automatically allocates any unassigned port of the host.

`HOST:CONTAINER` should always be specified as a (quoted) string, to avoid conflicts
with [YAML base-60 float](https://yaml.org/type/float.html).

IPv6 addresses can be enclosed in square brackets.

Examples:

```yml
ports:
  - "3000"
  - "3000-3005"
  - "8000:8000"
  - "9090-9091:8080-8081"
  - "49100:22"
  - "8000-9000:80"
  - "127.0.0.1:8001:8001"
  - "127.0.0.1:5000-5010:5000-5010"
  - "::1:6000:6000"
  - "[::1]:6001:6001"
  - "6060:6060/udp"
```

> [!NOTE]
>
> If host IP mapping is not supported by a container engine, Compose rejects
> the Compose file and ignores the specified host IP.

#### Long syntax

The long form syntax lets you configure additional fields that can't be
expressed in the short form.

- `target`: The container port.
- `published`: The publicly exposed port. It is defined as a string and can be set as a range using syntax `start-end`. It means the actual port is assigned a remaining available port, within the set range.
- `host_ip`: The host IP mapping. If it is not set, it binds to all network interfaces (`0.0.0.0`).
- `protocol`: The port protocol (`tcp` or `udp`). Defaults to `tcp`.
- `app_protocol`: The application protocol (TCP/IP level 4 / OSI level 7) this port is used for. This is optional and can be used as a hint for Compose to offer richer behavior for protocols that it understands. Introduced in Docker Compose version [2.26.0](/manuals/compose/releases/release-notes.md#2260).
- `mode`: Specifies how the port is published in a Swarm setup. If set to `host`, it publishes the port on every node in Swarm. If set to `ingress`, it allows load balancing across the nodes in Swarm. Defaults to `ingress`.
- `name`: A human-readable name for the port, used to document its usage within the service.

```yml
ports:
  - name: web
    target: 80
    host_ip: 127.0.0.1
    published: "8080"
    protocol: tcp
    app_protocol: http
    mode: host

  - name: web-secured
    target: 443
    host_ip: 127.0.0.1
    published: "8083-9000"
    protocol: tcp
    app_protocol: https
    mode: host
```

### `post_start`

{{< summary-bar feature_name="Compose post start" >}}

`post_start` defines a sequence of lifecycle hooks to run after a container has started. The exact timing of when the command is run is not guaranteed.

- `command`: Specifies the command to run once the container starts. This attribute is required, and you can choose to use either the shell form or the exec form.
- `user`: The user to run the command. If not set, the command is run with the same user as the main service command.
- `privileged`: Lets the `post_start` command run with privileged access.
- `working_dir`: The working directory in which to run the command. If not set, it is run in the same working directory as the main service command.
- `environment`: Sets environment variables specifically for the `post_start` command. While the command inherits the environment variables defined for the service’s main command, this section lets you add new variables or override existing ones.

```yaml
services:
  test:
    post_start:
      - command: ./do_something_on_startup.sh
        user: root
        privileged: true
        environment:
          - FOO=BAR
```

For more information, see [Use lifecycle hooks](/manuals/compose/how-tos/lifecycle.md).

### `pre_stop`

{{< summary-bar feature_name="Compose pre stop" >}}

`pre_stop` defines a sequence of lifecycle hooks to run before the container is stopped. These hooks won't run if the container stops by itself or is terminated suddenly.

Configuration is equivalent to [post_start](#post_start).

### `privileged`

`privileged` configures the service container to run with elevated privileges. Support and actual impacts are platform specific.

### `profiles`

`profiles` defines a list of named profiles for the service to be enabled under. If unassigned, the service is always started but if assigned, it is only started if the profile is activated.

If present, `profiles` follow the regex format of `[a-zA-Z0-9][a-zA-Z0-9_.-]+`.

```yaml
services:
  frontend:
    image: frontend
    profiles: ["frontend"]

  phpmyadmin:
    image: phpmyadmin
    depends_on:
      - db
    profiles:
      - debug
```

### `provider`

{{< summary-bar feature_name="Compose provider services" >}}

`provider` can be used to define a service that Compose won't manage directly. Compose delegated the service lifecycle to a dedicated or third-party component.

```yaml
  database:
    provider:
      type: awesomecloud
      options:
        type: mysql
        foo: bar
  app:
    image: myapp
    depends_on:
       - database
```

As Compose runs the application, the `awesomecloud` binary is used to manage the `database` service setup.
Dependent service `app` receives additional environment variables prefixed by the service name so it can access the resource.

For illustration, assuming `awesomecloud` execution produced variables `URL` and `API_KEY`, the `app` service
runs with environment variables `DATABASE_URL` and `DATABASE_API_KEY`.

As Compose stops the application, the `awesomecloud` binary is used to manage the `database` service tear down.

The mechanism used by Compose to delegate the service lifecycle to an external binary is described [here](https://github.com/docker/compose/tree/main/docs/extension.md).

For more information on using the `provider` attribute, see [Use provider services](/manuals/compose/how-tos/provider-services.md).

#### `type`

`type` attribute is required. It defines the external component used by Compose to manage setup and tear down lifecycle
events.

#### `options`

`options` are specific to the selected provider and not validated by the compose specification

### `pull_policy`

`pull_policy` defines the decisions Compose makes when it starts to pull images. Possible values are:

- `always`: Compose always pulls the image from the registry.
- `never`: Compose doesn't pull the image from a registry and relies on the platform cached image.
   If there is no cached image, a failure is reported.
- `missing`: Compose pulls the image only if it's not available in the platform cache.
   This is the default option if you are not also using the [Compose Build Specification](build.md).
  `if_not_present` is considered an alias for this value for backward compatibility. The `latest` tag is always pulled even when the `missing` pull policy is used.
- `build`: Compose builds the image. Compose rebuilds the image if it's already present.
- `daily`: Compose checks the registry for image updates if the last pull took place more than 24 hours ago.
- `weekly`: Compose checks the registry for image updates if the last pull took place more than 7 days ago.
- `every_<duration>`: Compose checks the registry for image updates if the last pull took place before `<duration>`. Duration can be expressed in weeks (`w`), days (`d`), hours (`h`), minutes (`m`), seconds (`s`) or a combination of these.

```yaml
services:
  test:
    image: nginx
    pull_policy: every_12h
```

### `read_only`

`read_only` configures the service container to be created with a read-only filesystem.

### `restart`

`restart` defines the policy that the platform applies on container termination.

- `no`: The default restart policy. It does not restart the container under any circumstances.
- `always`: The policy always restarts the container until its removal.
- `on-failure[:max-retries]`: The policy restarts the container if the exit code indicates an error.
Optionally, limit the number of restart retries the Docker daemon attempts.
- `unless-stopped`: The policy restarts the container irrespective of the exit code but stops
  restarting when the service is stopped or removed.

```yml
    restart: "no"
    restart: always
    restart: on-failure
    restart: on-failure:3
    restart: unless-stopped
```

You can find more detailed information on restart policies in the
[Restart Policies (--restart)](/reference/cli/docker/container/run.md#restart)
section of the Docker run reference page.

### `runtime`

`runtime` specifies which runtime to use for the service’s containers.

For example, `runtime` can be the name of [an implementation of OCI Runtime Spec](https://github.com/opencontainers/runtime-spec/blob/master/implementations.md), such as "runc".

```yml
web:
  image: busybox:latest
  command: true
  runtime: runc
```

The default is `runc`. To use a different runtime, see [Alternative runtimes](/manuals/engine/daemon/alternative-runtimes.md).

### `scale`

`scale` specifies the default number of containers to deploy for this service.
When both are set, `scale` must be consistent with the `replicas` attribute in the [Deploy Specification](deploy.md#replicas).

### `secrets`

{{% include "compose/services-secrets.md" %}}

Two different syntax variants are supported; the short syntax and the long syntax. Long and short syntax for secrets may be used in the same Compose file.

Compose reports an error if the secret doesn't exist on the platform or isn't defined in the
[`secrets` top-level section](secrets.md) of the Compose file.

Defining a secret in the top-level `secrets` must not imply granting any service access to it.
Such grant must be explicit within service specification as [secrets](secrets.md) service element.

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
    image: example/webapp
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
  service's task container, or absolute path of the file if an alternate location is required. Defaults to `source` if not specified.
- `uid` and `gid`: The numeric uid or gid that owns the file within
  `/run/secrets/` in the service's task containers. Default value is `USER`.
- `mode`: The [permissions](https://wintelguy.com/permissions-calc.pl) for the file to be mounted in `/run/secrets/`
  in the service's task containers, in octal notation.
  The default value is world-readable permissions (mode `0444`).
  The writable bit must be ignored if set. The executable bit may be set.

Note that support for `uid`, `gid`, and `mode` attributes are not implemented in Docker Compose when the source of the secret is a [`file`](secrets.md). This is because bind-mounts used under the hood don't allow uid remapping.

The following example sets the name of the `server-certificate` secret file to `server.cert`
within the container, sets the mode to `0440` (group-readable), and sets the user and group
to `103`. The value of `server-certificate` is set
to the contents of the file `./server.cert`.

```yml
services:
  frontend:
    image: example/webapp
    secrets:
      - source: server-certificate
        target: server.cert
        uid: "103"
        gid: "103"
        mode: 0o440
secrets:
  server-certificate:
    file: ./server.cert
```

### `security_opt`

`security_opt` overrides the default labeling scheme for each container.

```yml
security_opt:
  - label=user:USER
  - label=role:ROLE
```

For further default labeling schemes you can override, see [Security configuration](/reference/cli/docker/container/run.md#security-opt).

### `shm_size`

`shm_size` configures the size of the shared memory (`/dev/shm` partition on Linux) allowed by the service container.
It's specified as a [byte value](extension.md#specifying-byte-values).

### `stdin_open`

`stdin_open` configures a service's container to run with an allocated stdin. This is the same as running a container with the
`-i` flag. For more information, see [Keep stdin open](/reference/cli/docker/container/run.md#interactive).

Supported values are `true` or `false`.

### `stop_grace_period`

`stop_grace_period` specifies how long Compose must wait when attempting to stop a container if it doesn't
handle SIGTERM (or whichever stop signal has been specified with
[`stop_signal`](#stop_signal)), before sending SIGKILL. It's specified
as a [duration](extension.md#specifying-durations).

```yml
    stop_grace_period: 1s
    stop_grace_period: 1m30s
```

Default value is 10 seconds for the container to exit before sending SIGKILL.

### `stop_signal`

`stop_signal` defines the signal that Compose uses to stop the service containers.
If unset containers are stopped by Compose by sending `SIGTERM`.

```yml
stop_signal: SIGUSR1
```

### `storage_opt`

`storage_opt` defines storage driver options for a service.

```yml
storage_opt:
  size: '1G'
```

### `sysctls`

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
parameters (sysctls) at runtime](/reference/cli/docker/container/run.md#sysctl).

### `tmpfs`

`tmpfs` mounts a temporary file system inside the container. It can be a single value or a list.

```yml
tmpfs:
 - <path>
 - <path>:<options>
```

- `path`: The path inside the container where the tmpfs will be mounted.
- `options`: Comma-separated list of options for the tmpfs mount.

Available options:

- `mode`: Sets the file system permissions.
- `uid`: Sets the user ID that owns the mounted tmpfs.
- `gid`: Sets the group ID that owns the mounted tmpfs.

```yml
services:
  app:
    tmpfs:
      - /data:mode=755,uid=1009,gid=1009
      - /run
```

### `tty`

`tty` configures a service's container to run with a TTY. This is the same as running a container with the
`-t` or `--tty` flag. For more information, see [Allocate a pseudo-TTY](/reference/cli/docker/container/run.md#tty).

Supported values are `true` or `false`.

### `ulimits`

`ulimits` overrides the default `ulimits` for a container. It's specified either as an integer for a single limit
or as mapping for soft/hard limits.

```yml
ulimits:
  nproc: 65535
  nofile:
    soft: 20000
    hard: 40000
```

### `use_api_socket`

When `use_api_socket` is set, the container is able to interact with the underlying container engine through the API socket.
Your credentials are mounted inside the container so the container acts as a pure delegate for your commands relating to the container engine.
Typically, commands ran by container can `pull` and `push` to your registry.

### `user`

`user` overrides the user used to run the container process. The default is set by the image, for example Dockerfile `USER`. If it's not set, then `root`.

### `userns_mode`

`userns_mode` sets the user namespace for the service. Supported values are platform specific and may depend
on platform configuration.

```yml
userns_mode: "host"
```

### `uts`

{{< summary-bar feature_name="Compose uts" >}}

`uts` configures the UTS namespace mode set for the service container. When unspecified
it is the runtime's decision to assign a UTS namespace, if supported. Available values are:

- `'host'`: Results in the container using the same UTS namespace as the host.

```yml
    uts: "host"
```

### `volumes`

{{% include "compose/services-volumes.md" %}}

The following example shows a named volume (`db-data`) being used by the `backend` service,
and a bind mount defined for a single service.

```yml
services:
  backend:
    image: example/backend
    volumes:
      - type: volume
        source: db-data
        target: /data
        volume:
          nocopy: true
          subpath: sub
      - type: bind
        source: /var/run/postgres/postgres.sock
        target: /var/run/postgres/postgres.sock

volumes:
  db-data:
```

For more information about the `volumes` top-level element, see [Volumes](volumes.md).

#### Short syntax

The short syntax uses a single string with colon-separated values to specify a volume mount
(`VOLUME:CONTAINER_PATH`), or an access mode (`VOLUME:CONTAINER_PATH:ACCESS_MODE`).

- `VOLUME`: Can be either a host path on the platform hosting containers (bind mount) or a volume name.
- `CONTAINER_PATH`: The path in the container where the volume is mounted.
- `ACCESS_MODE`: A comma-separated `,` list of options:
  - `rw`: Read and write access. This is the default if none is specified.
  - `ro`: Read-only access.
  - `z`: SELinux option indicating that the bind mount host content is shared among multiple containers.
  - `Z`: SELinux option indicating that the bind mount host content is private and unshared for other containers.

> [!NOTE]
>
> The SELinux re-labeling bind mount option is ignored on platforms without SELinux.

> [!NOTE]
> Relative host paths are only supported by Compose that deploy to a
> local container runtime. This is because the relative path is resolved from the Compose file’s parent
> directory which is only applicable in the local case. When Compose deploys to a non-local
> platform it rejects Compose files which use relative host paths with an error. To avoid ambiguities
> with named volumes, relative paths should always begin with `.` or `..`.

> [!NOTE]
>
> For bind mounts, the short syntax creates a directory at the source path on the host if it doesn't exist. This is for backward compatibility with `docker-compose` legacy.
> It can be prevented by using long syntax and setting `create_host_path` to `false`.

#### Long syntax

The long form syntax lets you configure additional fields that can't be
expressed in the short form.

- `type`: The mount type. Either `volume`, `bind`, `tmpfs`, `image`, `npipe`, or `cluster`
- `source`: The source of the mount, a path on the host for a bind mount, a Docker image reference for an image mount, or the
  name of a volume defined in the
  [top-level `volumes` key](volumes.md). Not applicable for a tmpfs mount.
- `target`: The path in the container where the volume is mounted.
- `read_only`: Flag to set the volume as read-only.
- `bind`: Used to configure additional bind options:
  - `propagation`: The propagation mode used for the bind.
  - `create_host_path`: Creates a directory at the source path on host if there is nothing present. Defaults to `true`.
  - `selinux`: The SELinux re-labeling option `z` (shared) or `Z` (private)
- `volume`: Configures additional volume options:
  - `nocopy`: Flag to disable copying of data from a container when a volume is created.
  - `subpath`: Path inside a volume to mount instead of the volume root.
- `tmpfs`: Configures additional tmpfs options:
  - `size`: The size for the tmpfs mount in bytes (either numeric or as bytes unit).
  - `mode`: The file mode for the tmpfs mount as Unix permission bits as an octal number. Introduced in Docker Compose version [2.14.0](/manuals/compose/releases/release-notes.md#2260).
- `image`: Configures additional image options:
  - `subpath`: Path inside the source image to mount instead of the image root. Available in [Docker Compose version 2.35.0](/manuals/compose/releases/release-notes.md#2350)
- `consistency`: The consistency requirements of the mount. Available values are platform specific.

> [!TIP]
>
> Working with large repositories or monorepos, or with virtual file systems that are no longer scaling with your codebase?
> Compose now takes advantage of [Synchronized file shares](/manuals/desktop/features/synchronized-file-sharing.md) and automatically creates file shares for bind mounts.
> Ensure you're signed in to Docker with a paid subscription and have enabled both **Access experimental features** and **Manage Synchronized file shares with Compose** in Docker Desktop's settings.

### `volumes_from`

`volumes_from` mounts all of the volumes from another service or container. You can optionally specify
read-only access `ro` or read-write `rw`. If no access level is specified, then read-write access is used.

You can also mount volumes from a container that is not managed by Compose by using the `container:` prefix.

```yaml
volumes_from:
  - service_name
  - service_name:ro
  - container:container_name
  - container:container_name:rw
```

### `working_dir`

`working_dir` overrides the container's working directory which is specified by the image, for example Dockerfile's `WORKDIR`.
