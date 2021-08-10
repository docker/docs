---
description: Compose V2 RC1 in the Docker CLI
keywords: compose, V2, release candidate RC 1
title: Compose V2 release candidate
---

## Compose V2 and the new `docker compose` command

> Important
>
> The new Compose V2, which supports the `compose` command as part of the Docker CLI, is available as a release candidate with the Docker Desktop 3.6 release.
>
> Compose V2 integrates compose functions into the Docker platform, continuing to support most of the previous `docker-compose` features and flags. You can test the Compose V2 by simply replacing the dash (`-`) with a space, and by running `docker compose`, instead of `docker-compose`.
>
> As Docker Compose V2 is a release candidate, we recommend that you extensively test before using it in production environments.
{: .important}

Starting with Docker Desktop 3.4.0, you can run Compose V2 commands without modifying your invocations, by enabling the drop-in replacement of the previous `docker-compose` with the new command.  See the section [Installing Compose v2](#installing-compose-v2) for detailed instructions how to enable the drop-in replacement.

We will gradually turn this option on automatically for Docker Desktop users, so that users can seamlessly move to Docker Compose V2 without the need to upgrade any of their scripts. If you run into any problems with Compose V2, you can easily switch back to Compose v1 by either by making changes in Docker Desktop **Experimental** Settings, or by running the command `docker-compose disable-v2`.

Your feedback is important to us. Let us know your feedback on the new 'compose' command by creating an issue in the [Compose-CLI](https://github.com/docker/compose-cli/issues){:target="_blank" rel="noopener" class="_"} GitHub repository.
{: .important}

## Context of Docker Compose evolution

Introduction of the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"} makes a clean distinction between the Compose YAML file model and the `docker-compose` implementation. Making this change has enabled a number of enhancements, including adding the `compose` command directly into the Docker CLI,  being able to “up” a Compose application on cloud platforms by simply switching the Docker context, and launching of [Amazon ECS](/cloud/ecs-integration) and [Microsoft ACI](/cloud/aci-integration). As the Compose specification evolves, new features land faster in the Docker CLI.

While docker-compose is still supported and maintained, Compose V2 implementation relies directly on the compose-go bindings which are maintained as part of the specification. This allows us to include community proposals, experimental implementations by the Docker CLI and/or Engine, and deliver features faster to users. Compose V2 also  supports some of the newer additions to the Compose specification, such as [profiles](profiles.md) and [GPU](gpu-support.md) devices.

Additionally, Compose V2 also supports [Apple silicon](../docker-for-mac/apple-silicon.md).

For more information about the flags that are supported in the new compose command, see the [docker-compose compatibility list](cli-command-compatibility.md).

## Installing Compose V2

This section contains instructions on how to install Compose V2.

### Install on Mac and Windows

**Docker Desktop for Mac and for Windows version 3.2.1** and above includes the new Compose command along with the Docker CLI. Therefore, Windows and Mac users do not need to install the Compose V2 separately.

**Docker Desktop for Mac and for Windows version 3.4.0** and above also includes `docker-compose` drop-in replacement, allowing users to choose to use Compose V2 when using the `docker-compose` command.

We will progressively turn Docker Compose V2 on automatically for Docker Desktop users, so that users can seamlessly move to Docker Compose V2 without the need to change any of their scripts.  If you run into any problems with Compose V2, you can simply switch back to Compose v1, either in Docker Desktop, or in the CLI.

For Docker Desktop installation instructions, see:

- [Install Docker Desktop on Mac](../docker-for-mac/install.md)
- [Install Docker Desktop on Windows](../docker-for-windows/install.md)

To disable Docker Compose V2 using Docker Desktop:

1. From the Docker menu, click **Preferences** (**Settings** on Windows) > **Experimental features**.
2. Clear the **Use Docker Compose V2** check box.

To disable Docker Compose V2 using the CLI, run:

```console
$ docker-compose disable-v2
```

### Install on Linux

You can install Compose V2 by downloading the adequate binary for your system 
from the [project release page](https://github.com/docker/compose-cli/releases){:target="_blank" rel="noopener" class="_"} and copying it into `$HOME/.docker/cli-plugins` as `docker-compose`.

```console
$ mkdir -p ~/.docker/cli-plugins/
$ curl -SL https://github.com/docker/compose-cli/releases/download/v2.0.0-rc.1/docker-compose-linux-amd64 -o ~/.docker/cli-plugins/docker-compose
$ chmod +x ~/.docker/cli-plugins/docker-compose
```
