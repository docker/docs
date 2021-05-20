---
description: Compose V2 in Docker CLI
keywords: documentation, docs, docker, compose, containers
title: Compose V2 beta
---

## Compose V2 and new docker compose command

> Important
>
> Compose V2 and the `compose` command  in the Docker CLI is currently available as a beta version. We recommend that you do not use this in production environments.
>
> Your feedback is important to us. Let us know your feedback on the new 'compose' command by creating an issue in the [Compose-CLI](https://github.com/docker/compose-cli/issues){:target="_blank" rel="noopener" class="_"} GitHub repository.
{: .important}

The Docker CLI now supports the `compose` command, including most of the `docker-compose` features and flags, without the need for a separate tool.

You can replace the dash (`-`) with a space when you use `docker-compose` to switch over to `docker compose`.

Starting with Docker Desktop 3.4.0, you can also enable drop-in replacement if you don't want to modify your invocations, and still use `docker-compose` commands but use Compose V2. This can be turned on in Docker Desktop Experimental Settings, or using `docker-compose enable-v2`.

We will progressively turn this option on automatically for Docker Desktop users, so that users can seamlessly move to Docker Compose V2 without migration efforts. You can switch back to Compose V1 if needed in Docker Desktop Experimental Settings, or using `docker-compose disable-v2`.

## Context of Docker Compose evolution

Introduction of the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"} makes a clean distinction between the Compose YAML file model and the `docker-compose` implementation. Making this change has enabled a number of enhancements, including the launch of [Amazon ECS](/cloud/ecs-integration) and [Microsoft ACI](/cloud/aci-integration), being able to “up” a Compose application on cloud platforms simply by switching Docker context, and adding compose command directly into the Docker CLI.
As the Compose specification evolves, new features land faster in the Docker CLI. While `docker-compose` is still supported and maintained, Compose in the Docker CLI Go implementation relies directly on the compose-go bindings which are maintained as part of the specification. This allows us to include community proposals, experimental implementations by the Docker CLI and/or Engine, and deliver features faster to users. Compose in the Docker CLI already supports some of the newer additions to the Compose specification such as profiles and GPU devices.

For more information about the flags that are not yet supported in the new `compose` command, see the [docker-compose compatibility list](cli-command-compatibility.md).

## Installing Compose V2

### Install Compose V2 on Mac and Windows

**Docker Desktop for Mac and for Windows** version 3.2.1 and above includes the new Compose command along with the Docker CLI. Therefore, Windows and Mac users do not need to install the Compose V2 separately.

**Docker Desktop for Mac and for Windows** version 3.4.0 and above also includes docker-compose drop-in replacement, allowing users to choose to use Compose V2 when using the `docker-compose` command.

We will progressively turn this option on automatically for Docker Desktop users, so that users can seamlessly move to Docker Compose V2 without migration efforts. You can switch back to Compose V1 if needed in Docker Desktop Experimental Settings, or using `docker-compose disable-v2`.

For Docker Desktop installation instructions, see:

- [Install Docker Desktop on Mac](../docker-for-mac/install.md)
- [Install Docker Desktop on Windows](../docker-for-windows/install.md).

### Install Compose V2 on Linux

You can install the new Compose CLI, including Compose V2, using the following install script:

```console
$ curl -L https://raw.githubusercontent.com/docker/compose-cli/main/scripts/install/install_linux.sh | sh
```
