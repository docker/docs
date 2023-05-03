---
description: Key features and use cases of Docker Compose
keywords: >-
  documentation, docs, docker, compose, orchestration, containers, uses,
  features
title: Evolution of Compose
redirect_from:
 - /compose/cli-command-compatibility/
---
{% include compose-eol.md %}

This page provides information on the history of Compose and explains the key
differences between Compose V1 and Compose V2.

## History

The first release of Compose, written in Python, happened at the end of 2014.
Two other noticeable versions of Compose happened between 2014 and 2017, which
introduced new file format versions:

- [Compose 1.6.0 with file format V2](../compose-file/compose-file-v2/)
- [Compose 1.10.0 with file format V3](../compose-file/compose-file-v3/)

These three key file format versions and releases prior to v1.29.2 are
collectively referred to as Compose V1.

In mid-2020 Compose V2 was released. It merged Compose file formats V2 and V3
and was written in Go. The file format is defined by the
[Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"}.
Compose V2 is the latest and recommended version of Compose and is compatible
with Docker Engine version 19.03.0 and later. It provides improved integration
with other Docker command-line features, and simplifies installation on macOS,
Windows, and Linux.

It makes a clean distinction between the Compose YAML file model and the
`docker-compose` implementation. This change has enabled a number of
enhancements, including adding the `compose` command directly into the Docker
CLI,  being able to "up" a Compose application on cloud platforms by simply
switching the Docker context, and launching of
[Amazon ECS](../../cloud/ecs-integration.md) and
[Microsoft ACI](../../cloud/aci-integration.md). As the Compose specification
evolves, new features land faster in the Docker CLI.

> **A note about version numbers**
>
>In addition to Compose file format versions described here, the Compose binary
itself is on a release schedule, as shown in
[Compose releases](https://github.com/docker/compose/releases/). File format
versions do not necessarily increment with each release. For example, Compose
file format V3 was first introduced in Compose release 1.10.0, and versioned
gradually in subsequent releases.
>
>The latest Compose file format, defined by the Compose Specification, was
implemented by Docker Compose 1.27.0+.

## Differences between Compose V1 and Compose V2

Compose V2 integrates compose functions into the Docker platform, continuing to
support most of the previous `docker-compose` commands and flags. You can run
Compose V2 by replacing the hyphen (`-`) with a space, using `docker compose`
instead of `docker-compose`. The `compose` command in the Docker CLI,
`docker compose`, is expected to be a drop-in replacement for `docker-compose`.

If you see any Compose functionality that is not available in the `compose`
command, create an issue in
[Compose](https://github.com/docker/compose/issues){:target="_blank" rel="noopener" class="_"}
GitHub repository.

Compose V2 relies directly on the compose-go bindings that are maintained as
part of the specification. As a result, features, community proposals and
experimental implementations by Docker CLI and Docker Engine land faster into
Compose V2.

### Differences in command line interface

This section documents commands and flags that are different in Compose V2
compared to Compose V1.

`compose build --memory`: BuildKit does not yet support this option. The flag
is supported, but is hidden to avoid breaking existing Compose usage. It does
not have any effect.

There are no plans to support the commands and flags in the list below. Some of
these are already deprecated in Compose V1, some are not relevant to Compose
V2.

* `compose ps --filter KEY-VALUE` Not relevant due to its complicated usage
with the `service` command and lack of proper documentation in Compose V1.
* `compose rm --all` Deprecated in Compose V1.
* `compose scale` Deprecated in Compose V1 (use `compose up --scale` instead)
* `--compatibility` has different meaning in Compose V2. It means that V2 will
behave as V1 used to do.

Compose V2 uses `-` as word separator in container names while V1 used `_`.
Providing `--compatibility` flag to V2 will make it use `_` as in V1. Make sure
to stick to one of them, otherwise Compose will not be able to recognize the
container as an instance of the service.

#### New commands and flags in Compose V2

##### Copy

The `cp` command copies files or folders between service containers and the
local filesystem. This command is bidirectional, it can copy from or to the
service containers.

Copy a file from service container to local filesystem:

```console
$ docker compose cp my-service:~/path/to/myfile ~/local/path/to/copied/file
```

Copy a file from local filesystem to service container:

```console
$ docker compose cp ~/local/path/to/source/file my-service:~/path/to/copied/file
```

##### List

The `ls` command lists Compose projects. With no flags, it lists only running
projects. `--all` and `--filter` can be provided to further customize projects
to list. The output of `ls` can be further customized with `--format`. For
example:

```console
$ docker compose ls --all --format json
[{"Name":"dockergithubio","Status":"exited(1)","ConfigFiles":"/path/to/docs/docker-compose.yml"}]
```

##### Specify project with `--project-name`

With Compose V1, compose commands had to run either from the project directory
or by specifying `--file` or `--project-directory` command line flags. With
Compose V2, you can use `--project-name` or its shorthand `-p` to run commands
against a loaded project from any directory. For example:

```console
$ docker compose -p my-loaded-project restart my-service
```

This option works with `start`, `stop`, `restart` and `down` commands.

##### `config`

The `config` command shows the configuration used by Docker Compose after
normalization and templating. The resulting output might contain superficial
differences in formatting and style. For example, some fields in the Compose
Specification support both short and a long format so the output structure
might not match the input structure but is semantically equivalent. Comments in
the source file are not preserved. For example, with the following
`docker-compose.yaml`:

```yaml
services:
  web:
    image: nginx:latest
    ports:
      - 80:80
```

`docker compose config` will show `ports` expanded.
```yaml
name: docs-example
services:
  web:
    image: nginx:latest
    networks:
      default: null
    ports:
    - mode: ingress
      target: 80
      published: "80"
      protocol: tcp
networks:
  default:
    name: basic_default
```

The result is the configuration Compose will use to run the project.
