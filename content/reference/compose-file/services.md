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

A service definition supports all attributes defined in the [container specification](container-spec.md)
and the [workload specification](workload-spec.md), as well as the service-specific attributes on this page.

## Where attributes are defined

Services are long-running processes, but Compose also needs to describe tasks that run once and exit, such as database migrations, batch processing, and scheduled cleanups. The [`jobs`](jobs.md) top-level element declares those tasks. To avoid repeating the full container configuration for each element type, Compose defines the shared attributes once, in layers, and splits the reference to match:

- The [container specification](container-spec.md) covers the attributes of a single container. Services, jobs, and the init containers that [`pre_start`](#pre_start) declares all share it.
- The [workload specification](workload-spec.md) covers the attributes of an orchestrated workload: how Compose builds, connects, checks, and exposes it. Services and jobs share it, but `pre_start` init containers don't.
- This page covers the attributes that apply only to services.

A service accepts the attributes from all three pages. Use the following table to find the page that documents an attribute.

| Page | Used by | Attributes |
| ---- | ------- | ---------- |
| [Container specification](container-spec.md) | Services, jobs, `pre_start` init containers | [`annotations`](container-spec.md#annotations) [`blkio_config`](container-spec.md#blkio_config) [`cpu_count`](container-spec.md#cpu_count) [`cpu_percent`](container-spec.md#cpu_percent) [`cpu_shares`](container-spec.md#cpu_shares) [`cpu_period`](container-spec.md#cpu_period) [`cpu_quota`](container-spec.md#cpu_quota) [`cpu_rt_runtime`](container-spec.md#cpu_rt_runtime) [`cpu_rt_period`](container-spec.md#cpu_rt_period) [`cpus`](container-spec.md#cpus) [`cpuset`](container-spec.md#cpuset) [`cap_add`](container-spec.md#cap_add) [`cap_drop`](container-spec.md#cap_drop) [`cgroup`](container-spec.md#cgroup) [`cgroup_parent`](container-spec.md#cgroup_parent) [`command`](container-spec.md#command) [`configs`](container-spec.md#configs) [`credential_spec`](container-spec.md#credential_spec) [`device_cgroup_rules`](container-spec.md#device_cgroup_rules) [`devices`](container-spec.md#devices) [`dns`](container-spec.md#dns) [`dns_opt`](container-spec.md#dns_opt) [`dns_search`](container-spec.md#dns_search) [`domainname`](container-spec.md#domainname) [`driver_opts`](container-spec.md#driver_opts) [`entrypoint`](container-spec.md#entrypoint) [`env_file`](container-spec.md#env_file) [`environment`](container-spec.md#environment) [`extra_hosts`](container-spec.md#extra_hosts) [`gpus`](container-spec.md#gpus) [`group_add`](container-spec.md#group_add) [`hostname`](container-spec.md#hostname) [`image`](container-spec.md#image) [`init`](container-spec.md#init) [`ipc`](container-spec.md#ipc) [`isolation`](container-spec.md#isolation) [`labels`](container-spec.md#labels) [`label_file`](container-spec.md#label_file) [`logging`](container-spec.md#logging) [`mac_address`](container-spec.md#mac_address) [`mem_limit`](container-spec.md#mem_limit) [`mem_reservation`](container-spec.md#mem_reservation) [`mem_swappiness`](container-spec.md#mem_swappiness) [`memswap_limit`](container-spec.md#memswap_limit) [`models`](container-spec.md#models) [`network_mode`](container-spec.md#network_mode) [`networks`](container-spec.md#networks) [`oom_kill_disable`](container-spec.md#oom_kill_disable) [`oom_score_adj`](container-spec.md#oom_score_adj) [`pid`](container-spec.md#pid) [`pids_limit`](container-spec.md#pids_limit) [`platform`](container-spec.md#platform) [`privileged`](container-spec.md#privileged) [`pull_policy`](container-spec.md#pull_policy) [`read_only`](container-spec.md#read_only) [`runtime`](container-spec.md#runtime) [`secrets`](container-spec.md#secrets) [`security_opt`](container-spec.md#security_opt) [`shm_size`](container-spec.md#shm_size) [`stop_grace_period`](container-spec.md#stop_grace_period) [`stop_signal`](container-spec.md#stop_signal) [`storage_opt`](container-spec.md#storage_opt) [`sysctls`](container-spec.md#sysctls) [`tmpfs`](container-spec.md#tmpfs) [`ulimits`](container-spec.md#ulimits) [`use_api_socket`](container-spec.md#use_api_socket) [`user`](container-spec.md#user) [`userns_mode`](container-spec.md#userns_mode) [`uts`](container-spec.md#uts) [`volumes`](container-spec.md#volumes) [`volumes_from`](container-spec.md#volumes_from) [`working_dir`](container-spec.md#working_dir) |
| [Workload specification](workload-spec.md) | Services, jobs | [`build`](workload-spec.md#build) [`depends_on`](workload-spec.md#depends_on) [`expose`](workload-spec.md#expose) [`healthcheck`](workload-spec.md#healthcheck) [`ports`](workload-spec.md#ports) [`stdin_open`](workload-spec.md#stdin_open) [`tty`](workload-spec.md#tty) |
| Services (this page) | Services | [`attach`](#attach) [`container_name`](#container_name) [`deploy`](#deploy) [`develop`](#develop) [`extends`](#extends) [`external_links`](#external_links) [`links`](#links) [`post_start`](#post_start) [`pre_start`](#pre_start) [`pre_stop`](#pre_stop) [`profiles`](#profiles) [`provider`](#provider) [`restart`](#restart) [`scale`](#scale) |

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
    image: postgres:18
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

### `attach`

{{< summary-bar feature_name="Compose attach" >}}

When `attach` is defined and set to `false` Compose does not collect service logs,
until you explicitly request it to.

The default service configuration is `attach: true`.

### `container_name`

`container_name` is a string that specifies a custom container name, rather than a name generated by default.

```yml
container_name: my-web-container
```

Compose does not scale a service beyond one container if the Compose file specifies a
`container_name`. Attempting to do so results in an error.

`container_name` follows the regex format of `[a-zA-Z0-9][a-zA-Z0-9_.-]+`

### `deploy`

`deploy` specifies the configuration for the deployment and lifecycle of services, as defined [in the Compose Deploy Specification](deploy.md).

### `develop`

{{< summary-bar feature_name="Compose develop" >}}

`develop` specifies the development configuration for maintaining a container in sync with source, as defined in the [Development Section](develop.md).

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

`extends` is not supported when deploying with `docker stack deploy`.

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
[`depends_on`](workload-spec.md#depends_on), so they determine the order of service startup.

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

### pre_start

{{< summary-bar feature_name="Compose pre_start" >}}

`pre_start` defines a sequence of init containers to run before the service container is started. Each step runs to completion, in declared order, and the service container only starts once every step has exited `0`. A non-zero exit fails the bring-up of the service and its dependents.

Unlike `post_start` and `pre_stop`, which run a command inside the running service container, each `pre_start` step runs in its own ephemeral container, created after the service container is created but before it is started.

A `pre_start` step is a full [container specification](container-spec.md), so you can declare any container attribute on the step. The step inherits every attribute it doesn't declare from the service. A `pre_start` container is the service container with the image, the command, and anything else the step declares replaced. In particular:

- `image`: If omitted, the parent service's image is used.
- `command`: Optional when the chosen image's entrypoint already runs the intended command.
- `environment`: The step inherits the environment of the service's main command. Declared values are appended to the inherited ones or override them.
- `volumes`: If omitted, the step shares the service's declared volume mounts. Re-declaring a mount overrides how the step mounts it, for example read-write where the service mounts it read-only.

Re-declaring an inherited attribute follows the same rules as [merging Compose files](merge.md): mappings merge, scalars replace, sequences append, and the `!reset` and `!override` tags are available.

Attributes that only apply to an orchestrated workload aren't part of the container specification, so you can't declare them on a step. These include `ports`, `healthcheck`, `depends_on`, and `build`. The [workload specification](workload-spec.md) defines them.

In addition to container attributes, a step supports:

- `per_replica`: Whether the step runs once per service replica (`true`), or once for the service as a whole before any replica starts (`false`, the default). Set it to `true` for per-instance work, such as initializing a replica's anonymous volumes or registering a per-replica network identity.

`pre_start` steps only run once the service's `depends_on` conditions have been satisfied, so a step can rely on those dependencies the same way the main service command does. A`pre_start` container joins the same networks as the service, so it can reach services declared in `depends_on`, and shares the service's declared volume mounts so files it produces in a shared volume are visible to the service. 

With `per_replica: false` and a scaled service, only mounts that are shared across replicas (named volumes, bind mounts) are usable. Per-instance mounts (`tmpfs`, anonymous volumes) cannot be addressed by a single run. This is not an error. The steps run without access to per-instance mounts. Data a `per_replica: false` step must share with the service belongs in a named volume or bind mount.

A `pre_start` step that has already succeeded for its current definition is not re-run on a subsequent `up`, nor when the service container restarts under its `restart` policy. A step runs again when its definition changes, when the previous run did not succeed, or when the service is recreated. For example, after a change to the service configuration or an explicit forced recreation.

```yaml
services:
  app:
    image: myapp:latest
    depends_on:
      db:
        condition: service_healthy
    pre_start:
      - command: ["./manage.py", "migrate"]
      - image: busybox
        command: sh -c 'chown -R 1000:1000 /data'
    volumes:
      - data:/data

  db:
    image: postgres:16

volumes:
  data:
```

### `pre_stop`

{{< summary-bar feature_name="Compose pre stop" >}}

`pre_stop` defines a sequence of lifecycle hooks to run before the container is stopped. These hooks won't run if the container stops by itself or is terminated suddenly.

Configuration is equivalent to [`post_start`](#post_start).

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

The mechanism used by Compose to delegate the service lifecycle to an external binary is described in the [Compose extensibility documentation](https://github.com/docker/compose/tree/main/docs/extension.md).

For more information on using the `provider` attribute, see [Use provider services](/manuals/compose/how-tos/provider-services.md).

#### `type`

`type` attribute is required. It defines the external component used by Compose to manage setup and tear down lifecycle
events.

#### `options`

`options` are specific to the selected provider and not validated by the compose specification

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
[Restart Policies (--restart)](/reference/cli/docker/container/run/#restart)
section of the Docker run reference page.

### `scale`

`scale` specifies the default number of containers to deploy for this service.
When both are set, `scale` must be consistent with the `replicas` attribute in the [Deploy Specification](deploy.md#replicas).
