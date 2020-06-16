---
description: Compose file reference
keywords: fig, composition, compose version 1, docker
redirect_from:
- /compose/yml
title: Compose file version 1 reference
toc_max: 4
toc_min: 1
hide_from_sitemap: true
---

## Reference and guidelines

These topics describe version 1 of the Compose file format. This is the oldest
version.

## Compose and Docker compatibility matrix

There are several versions of the Compose file format – 1, 2, 2.x, and 3.x The
table below is a quick look. For full details on what each version includes and
how to upgrade, see **[About versions and upgrading](compose-versioning.md)**.

{% include content/compose-matrix.md %}

## Service configuration reference

The Version 1 Compose file is a [YAML](http://yaml.org/) file that defines [services](#service-configuration-reference).

The default path for a Compose file is `./docker-compose.yml`.

> **Tip**: You can use either a `.yml` or `.yaml` extension for this file.
> They both work.

A service definition contains configuration which is applied to each
container started for that service, much like passing command-line parameters to
`docker run`.

As with `docker run`, options specified in the Dockerfile, such as `CMD`,
`EXPOSE`, `VOLUME`, `ENV`, are respected by default - you don't need to
specify them again in `docker-compose.yml`.

This section contains a list of all configuration options supported by a service
definition in version 1.

### build

Configuration options that are applied at build time.

`build` can specified as a string containing a path to the build
context.

```yaml
build: ./dir
```

> **Note**
>
> In [version 1 file format](compose-versioning.md#version-1), `build` is
> different in two ways:
>
> * Only the string form (`build: .`) is allowed - not the object
>   form that is allowed in Version 2 and up.
> * Using `build` together with [`image`](#image) is not allowed.
>   Attempting to do so results in an error.

#### dockerfile

Alternate Dockerfile.

Compose uses an alternate file to build with. A build path must also be
specified.

```yaml
build: .
dockerfile: Dockerfile-alternate
```

> **Note**
>
> In the [version 1 file format](compose-versioning.md#version-1), `dockerfile`
> is different from newer versions in two ways:
>
> * It appears alongside `build`, not as a sub-option:
> * Using `dockerfile` together with [`image`](#image) is not allowed.
>   Attempting to do so results in an error.

### cap_add, cap_drop

Add or drop container capabilities.
See `man 7 capabilities` for a full list.

```yaml
cap_add:
  - ALL

cap_drop:
  - NET_ADMIN
  - SYS_ADMIN
```

> **Note**: These options are ignored when
> [deploying a stack in swarm mode](../../engine/reference/commandline/stack_deploy.md)
> with a (version 3) Compose file.

### command

Override the default command.

```yaml
command: bundle exec thin -p 3000
```

The command can also be a list, in a manner similar to
[dockerfile](../../engine/reference/builder.md#cmd):

```yaml
command: ["bundle", "exec", "thin", "-p", "3000"]
```

### cgroup_parent

Specify an optional parent cgroup for the container.

```yaml
cgroup_parent: m-executor-abcd
```

### container_name

Specify a custom container name, rather than a generated default name.

```yaml
container_name: my-web-container
```

Because Docker container names must be unique, you cannot scale a service
beyond 1 container if you have specified a custom name. Attempting to do so
results in an error.

### devices

List of device mappings.  Uses the same format as the `--device` docker
client create option.

```yaml
devices:
  - "/dev/ttyUSB0:/dev/ttyUSB0"
```

### dns

Custom DNS servers. Can be a single value or a list.

```yaml
dns: 8.8.8.8
```

```yaml
dns:
  - 8.8.8.8
  - 9.9.9.9
```

### dns_search

Custom DNS search domains. Can be a single value or a list.

```yaml
dns_search: example.com
```

```yaml
dns_search:
  - dc1.example.com
  - dc2.example.com
```

### entrypoint

Override the default entrypoint.

```yaml
entrypoint: /code/entrypoint.sh
```

The entrypoint can also be a list, in a manner similar to
[dockerfile](../../engine/reference/builder.md#entrypoint):

```yaml
entrypoint: ["php", "-d", "memory_limit=-1", "vendor/bin/phpunit"]
```

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

```yaml
env_file: .env
```

```yaml
env_file:
  - ./common.env
  - ./apps/web.env
  - /opt/runtime_opts.env
```

Compose expects each line in an env file to be in `VAR=VAL` format. Lines
beginning with `#` are processed as comments and are ignored. Blank lines are
also ignored.

```console
# Set Rails/Rack environment
RACK_ENV=development
```

> **Note**: If your service specifies a [build](#build) option, variables
> defined in environment files are _not_ automatically visible during the
> build.

The value of `VAL` is used as is and not modified at all. For example if the
value is surrounded by quotes (as is often the case of shell variables), the
quotes are included in the value passed to Compose.

Keep in mind that _the order of files in the list is significant in determining
the value assigned to a variable that shows up more than once_. The files in the
list are processed from the top down. For the same variable specified in file
`a.env` and assigned a different value in file `b.env`, if `b.env` is
listed below (after), then the value from `b.env` stands. For example, given the
following declaration in `docker-compose.yml`:

```yaml
services:
  some-service:
    env_file:
      - a.env
      - b.env
```

And the following files:

```console
# a.env
VAR=1
```

and

```console
# b.env
VAR=hello
```

`$VAR` is `hello`.

### environment

Add environment variables. You can use either an array or a dictionary. Any
boolean values (true, false, yes, no) need to be enclosed in quotes to ensure
they are not converted to True or False by the YML parser.

Environment variables with only a key are resolved to their values on the
machine Compose is running on, which can be helpful for secret or host-specific values.

```yaml
environment:
  RACK_ENV: development
  SHOW: 'true'
  SESSION_SECRET:
```

```yaml
environment:
  - RACK_ENV=development
  - SHOW=true
  - SESSION_SECRET
```

> **Note**: If your service specifies a [build](#build) option, variables
> defined in `environment` are _not_ automatically visible during the
> build.

### expose

Expose ports without publishing them to the host machine - they'll only be
accessible to linked services. Only the internal port can be specified.

```yaml
expose:
  - "3000"
  - "8000"
```

### extends

Extend another service, in the current file or another, optionally overriding
configuration.

You can use `extends` on any service together with other configuration keys.
The `extends` value must be a dictionary defined with a required `service`
and an optional `file` key.

```yaml
extends:
  file: common.yml
  service: webapp
```

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
[the extends documentation](../extends.md#extending-services).

### external_links

Link to containers started outside this `docker-compose.yml` or even outside of
Compose, especially for containers that provide shared or common services.
`external_links` follow semantics similar to `links` when
specifying both the container name and the link alias (`CONTAINER:ALIAS`).

```yaml
external_links:
  - redis_1
  - project_db_1:mysql
  - project_db_1:postgresql
```

### extra_hosts

Add hostname mappings. Use the same values as the docker client `--add-host` parameter.

```yaml
extra_hosts:
  - "somehost:162.242.195.82"
  - "otherhost:50.31.209.229"
```

An entry with the ip address and hostname is created in `/etc/hosts` inside containers for this service, e.g:

```console
162.242.195.82  somehost
50.31.209.229   otherhost
```

### image

Specify the image to start the container from. Can either be a repository/tag or
a partial image ID.

```yaml
image: redis
```
```yaml
image: ubuntu:18.04
```
```yaml
image: tutum/influxdb
```
```yaml
image: example-registry.com:4000/postgresql
```
```yaml
image: a4bc65fd
```

If the image does not exist, Compose attempts to pull it, unless you have also
specified [build](#build), in which case it builds it using the specified
options and tags it with the specified tag.

> **Note**: In the [version 1 file format](compose-versioning.md#version-1),
> using [`build`](#build) together with `image` is not allowed. Attempting to do
> so results in an error.

### labels

Add metadata to containers using [Docker labels](../../config/labels-custom-metadata.md). You can use either an array or a dictionary.

It's recommended that you use reverse-DNS notation to prevent your labels from conflicting with those used by other software.

```yaml
labels:
  com.example.description: "Accounting webapp"
  com.example.department: "Finance"
  com.example.label-with-empty-value: ""
```

```yaml
labels:
  - "com.example.description=Accounting webapp"
  - "com.example.department=Finance"
  - "com.example.label-with-empty-value"
```

### links

Link to containers in another service. Either specify both the service name and
a link alias (`"SERVICE:ALIAS"`), or just the service name.

> Links are a legacy option. We recommend using
> [networks](../networking.md) instead.

```yaml
web:
  links:
    - "db"
    - "db:database"
    - "redis"
```

Containers for the linked service are reachable at a hostname identical to
the alias, or the service name if no alias was specified.

Links also express dependency between services in the same way as
[depends_on](compose-file-v2.md#depends_on), so they determine the order of service startup.

> **Note**
>
> If you define both links and [networks](index.md#networks), services with
> links between them must share at least one network in common in order to
> communicate.

### log_driver

> [Version 1 file format](compose-versioning.md#version-1) only. In version 2 and up, use
> [logging](index.md#logging).

Specify a log driver. The default is `json-file`.

```yaml
log_driver: syslog
```

### log_opt

> [Version 1 file format](compose-versioning.md#version-1) only. In version 2 and up, use
> [logging](index.md#logging).

Specify logging options as key-value pairs. An example of `syslog` options:

```yaml
log_opt:
  syslog-address: "tcp://192.168.0.42:123"
```

### net

> [Version 1 file format](compose-versioning.md#version-1) only. In version 2 and up, use
> [network_mode](index.md#network_mode) and [networks](index.md#networks).

Network mode. Use the same values as the docker client `--net` parameter.
The `container:...` form can take a service name instead of a container name or
id.

```yaml
net: "bridge"
```
```yaml
net: "host"
```
```yaml
net: "none"
```
```yaml
net: "service:[service name]"
```
```yaml
net: "container:[container name/id]"
```

### pid

```yaml
pid: "host"
```

Sets the PID mode to the host PID mode. This turns on sharing between
container and the host operating system the PID address space. Containers
launched with this flag can access and manipulate other
containers in the bare-metal machine's namespace and vice versa.

### ports

Expose ports. Either specify both ports (`HOST:CONTAINER`), or just the container
port (an ephemeral host port is chosen).

> **Note**: When mapping ports in the `HOST:CONTAINER` format, you may experience
> erroneous results when using a container port lower than 60, because YAML
> parses numbers in the format `xx:yy` as a base-60 value. For this reason,
> we recommend always explicitly specifying your port mappings as strings.

```yaml
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
```

### security_opt

Override the default labeling scheme for each container.

```yaml
security_opt:
  - label:user:USER
  - label:role:ROLE
```

### stop_signal

Sets an alternative signal to stop the container. By default `stop` uses
SIGTERM. Setting an alternative signal using `stop_signal` causes
`stop` to send that signal instead.

```yaml
stop_signal: SIGUSR1
```

### ulimits

Override the default ulimits for a container. You can either specify a single
limit as an integer or soft/hard limits as a mapping.

```yaml
ulimits:
  nproc: 65535
  nofile:
    soft: 20000
    hard: 40000
```

### volumes, volume\_driver

Mount paths or named volumes, optionally specifying a path on the host machine
(`HOST:CONTAINER`), or an access mode (`HOST:CONTAINER:ro`).
For [version 2 files](compose-versioning.md#version-2), named volumes need to be specified with the
[top-level `volumes` key](compose-file-v2.md#volume-configuration-reference).
When using [version 1](compose-versioning.md#version-1), the Docker Engine creates the named
volume automatically if it doesn't exist.

You can mount a relative path on the host, which expands relative to
the directory of the Compose configuration file being used. Relative paths
should always begin with `.` or `..`.

```yaml
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
```

If you do not use a host path, you may specify a `volume_driver`.

```yaml
volume_driver: mydriver
```

There are several things to note, depending on which
[Compose file version](compose-versioning.md#versioning) you're using:

- For [version 1 files](compose-versioning.md#version-1), both named volumes and
  container volumes use the specified driver.
- No path expansion is done if you have also specified a `volume_driver`.
  For example, if you specify a mapping of `./foo:/data`, the `./foo` part
  is passed straight to the volume driver without being expanded.

See [Docker Volumes](../../storage/volumes.md) and
[Volume Plugins](/engine/extend/plugins_volume/) for more information.

### volumes_from

Mount all of the volumes from another service or container, optionally
specifying read-only access (``ro``) or read-write (``rw``). If no access level
is specified, then read-write is used.

```yaml
volumes_from:
  - service_name
  - service_name:ro
```

### cpu\_shares, cpu\_quota, cpuset, domainname, hostname, ipc, mac\_address, mem\_limit, memswap\_limit, mem\_swappiness, privileged, read\_only, restart, shm\_size, stdin\_open, tty, user, working\_dir

Each of these is a single value, analogous to its
[docker run](../../engine/reference/run.md) counterpart.

```yaml
cpu_shares: 73
cpu_quota: 50000
cpuset: 0,1

user: postgresql
working_dir: /code

domainname: foo.com
hostname: foo
ipc: host
mac_address: 02:42:ac:11:65:43

mem_limit: 1000000000
memswap_limit: 2000000000
privileged: true

restart: always

read_only: true
shm_size: 64M
stdin_open: true
tty: true
```

## Compose documentation

- [User guide](../index.md)
- [Installing Compose](../install.md)
- [Compose file versions and upgrading](compose-versioning.md)
- [Samples](../../samples/index.md)
- [Command line reference](../reference/index.md)
