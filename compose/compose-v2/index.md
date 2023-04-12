---
description: Key features and use cases of Docker Compose
keywords: documentation, docs, docker, compose, orchestration, containers, uses, features
title: Evolution of Compose
redirect_from:
 - /compose/cli-command-compatibility/
---
{% include compose-eol.md %}

This page provides information on the history of Compose and explains the key differences between Compose V1 and Compose V2. 

## History

The first release of Compose, written in Python, happened at the end of 2014. 
Between 2014 and 2017 two other noticeable versions of Compose, which introduced new file format versions, were released:

- [Compose 1.6.0 with file format V2](../compose-file/compose-file-v2/)
- [Compose 1.10.0 with file format V3](../compose-file/compose-file-v3/)

These three key file format versions and releases prior to v1.29.2 are collectively referred to as Compose V1. 

In mid-2020 Compose V2 was released. It merged Compose file format V2 and V3 and was written in Go. The file format is defined by the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"}. Compose V2 is the latest and recommended version of Compose. It provides improved integration with other Docker command-line features, and simplified installation on macOS, Windows, and Linux.  

It makes a clean distinction between the Compose YAML file model and the `docker-compose`
implementation. Making this change has enabled a number of enhancements, including
adding the `compose` command directly into the Docker CLI,  being able to "up" a
Compose application on cloud platforms by simply switching the Docker context,
and launching of [Amazon ECS](../../cloud/ecs-integration.md) and [Microsoft ACI](../../cloud/aci-integration.md).
As the Compose specification evolves, new features land faster in the Docker CLI.

> **A note about version numbers**
>
>In addition to Compose file format versions described above, the Compose binary itself is on a release schedule, as shown in [Compose releases](https://github.com/docker/compose/releases/). File format versions do not necessarily increment with each release. For example, Compose file format V3 was first introduced in Compose release 1.10.0, and versioned gradually in subsequent releases.
>
>The latest Compose file format, defined by the Compose Specification, was implemented by Docker Compose 1.27.0+.

## Differences between Compose V1 and Compose V2

Compose V2 integrates compose functions into the Docker platform, continuing to support most of the previous `docker-compose` features and flags. You can run Compose V2 by replacing the hyphen (`-`) with a space, using `docker compose`, instead of `docker-compose`.

The `compose` command in the Docker CLI supports most of the `docker-compose` commands and flags. It is expected to be a drop-in replacement for `docker-compose`. 

If you see any Compose functionality that is not available in the `compose` command, create an issue in the [Compose](https://github.com/docker/compose/issues){:target="_blank" rel="noopener" class="_"} GitHub repository, so we can prioritize it.

Compose V2 relies directly on the compose-go bindings which are maintained as part
of the specification. This allows us to include community proposals, experimental
implementations by the Docker CLI and/or Engine, and deliver features faster to
users. 

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab1">Commands not yet implemented</a></li>
  <li><a data-toggle="tab" data-target="#tab2">Flags not be implemented</a></li>
  <li><a data-toggle="tab" data-target="#tab3">New commands in Compose v2</a></li>
</ul>
<div class="tab-content">
<div id="tab1" class="tab-pane fade in active" markdown="1">

The following commands have not been implemented yet, and may be implemented at a later time.
Let us know if these commands are a higher priority for your use cases.

`compose build --memory`: This option is not yet supported by BuildKit. The flag is currently supported, but is hidden to avoid breaking existing Compose usage. It does not have any effect.

<hr>
</div>
<div id="tab2" class="tab-pane fade" markdown="1">

The list below includes the flags that we are not planning to support in Compose in the Docker CLI,
either because they are already deprecated in `docker-compose`, or because they are not relevant for Compose in the Docker CLI.

* `compose ps --filter KEY-VALUE` Not relevant due to its complicated usage with the `service` command and also because it is not documented properly in `docker-compose`.
* `compose rm --all` Deprecated in docker-compose.
* `compose scale` Deprecated in docker-compose (use `compose up --scale` instead)

Global flags:

* `--compatibility` has been resignified Docker Compose V2. This now means that in the command running V2 will behave as V1 used to do.
  * One difference is in the word separator on container names. V1 used to use `_` as separator while V2 uses `-` to keep the names more hostname friendly. So when using `--compatibility` Docker 
    Compose should use `_` again. Just make sure to stick to one of them otherwise Docker Compose will not be able to recognize the container as an instance of the service.
<hr>
</div>
<div id="tab3" class="tab-pane fade" markdown="1">

#### Copy

The `cp` command is intended to copy files or folders between service containers and the local filesystem.  
This command is a bidirectional command, we can copy **from** or **to** the service containers.

Copy a file from a service container to the local filesystem:

```console
$ docker compose cp my-service:~/path/to/myfile ~/local/path/to/copied/file
```

We can also copy from the local filesystem to all the running containers of a service:

```console
$ docker compose cp --all ~/local/path/to/source/file my-service:~/path/to/copied/file
```


#### List

The ls command is intended to list the Compose projects. By default, the command only lists the running projects, 
we can use flags to display the stopped projects, to filter by conditions and change the output to `json` format for example.

```console
$ docker compose ls --all --format json
[{"Name":"dockergithubio","Status":"exited(1)","ConfigFiles":"/path/to/docs/docker-compose.yml"}]
```

### Use `--project-name` with Compose commands

With the GA version of Compose, you can run some commands:
- outside of directory containing the project compose file
- or without specifying the path of the Compose with the `--file` flag
- or without specifying the project directory with the `--project-directory` flag

When a compose project has been loaded once, we can just use the `-p` or `--project-name` to reference it:

```console
$ docker compose -p my-loaded-project restart my-service
```

This option works with the `start`, `stop`, `restart` and `down` commands.

### Config command

The config command is intended to show the configuration used by Docker Compose to run the actual project after normalization and templating. The resulting output might contain superficial differences in formattting and style.
For example, some fields in the Compose Specification support both short and a long format so the output structure might not match the input structure but is guaranteed to be semantically equivalent.

Similarly, comments in the source file are not preserved.

In the example below we can see the config command expanding the `ports` section:

docker-compose.yml:
```yaml
services:
  web:
    # default to latest but allow overriding the tag
    image: nginx:${TAG-latest}
    ports:
      - 80:80
```
With `$ docker compose config` the output turns into:
```yaml
name: docs-example
services:
  web:
    image: nginx:stable-alpine
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

The result above is a full size configuration of what will be used by Docker Compose to run the project.
<hr>
</div>
</div>
