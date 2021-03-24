---
description: Compose CLI Tech Preview
keywords: documentation, docs, docker, compose, containers
title: Compose CLI Tech Preview
---

## New docker compose command

> Important
>
> The `compose` command  in the Docker CLI is currently available as a Tech Preview. We recommend that you do not use this in production environments.
>
> Your feedback is important to us. Let us know your feedback on the new 'compose' command by creating an issue in the [Compose-CLI](https://github.com/docker/compose-cli/issues){:target="_blank" rel="noopener" class="_"} GitHub repository.
{: .important}

The Docker CLI now supports the `compose` command, including most of the `docker-compose` features and flags, without the need for a separate tool.

You can replace the dash (`-`) with a space when you use `docker-compose` to switch over to `docker compose`. You can also use them interchangeably, so that you are not locked-in with the new `compose` command and, if needed, you can still use `docker-compose`.

Introduction of the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"} makes a clean distinction between the Compose YAML file model and the `docker-compose` implementation. Making this change has enabled a number of enhancements, including the launch of [Amazon ECS](/cloud/ecs-integration) and [Microsoft ACI](/cloud/aci-integration), being able to “up” a Compose application on cloud platforms simply by switching Docker context, and adding compose command directly into the Docker CLI.
As the Compose specification evolves, new features land faster in the Docker CLI. While `docker-compose` is still supported and maintained, Compose in the Docker CLI Go implementation relies directly on the compose-go bindings which are maintained as part of the specification. This allows us to include community proposals, experimental implementations by the Docker CLI and/or Engine, and deliver features faster to users. Compose in the Docker CLI already supports some of the newer additions to the Compose specification such as profiles and GPU devices.

For more information about the flags that are not yet supported in the new `compose` command, see the [docker-compose compatibility list](cli-command-compatibility.md).

## Installing the Compose CLI Tech Preview

### Install Compose CLI Tech Preview on Mac and Windows

**Docker Desktop for Mac and for Windows** version 3.2.1 and above includes the new Compose command along with the Docker CLI. Therefore, Windows and Mac users do not need to install the Compose CLI Tech Preview separately.

For Docker Desktop installation instructions, see:

- [Install Docker Desktop on Mac](../docker-for-mac/install.md)
- [Install Docker Desktop on Windows](../docker-for-windows/install.md).

### Install Compose CLI Tech Preview on Linux

You can install the new Compose CLI, including this Tech Preview, using the following install script:

```console
$ curl -L https://raw.githubusercontent.com/docker/compose-cli/main/scripts/install/install_linux.sh | sh
```
