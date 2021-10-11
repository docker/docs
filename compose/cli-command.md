---
description: Docker Compose
keywords: compose, V2
title: Compose V2
---

## Compose V2 and the new `docker compose` command

> Important
>
> The new Compose V2, which supports the `compose` command as part of the Docker CLI, is now available.
>
> Compose V2 integrates compose functions into the Docker platform, continuing to support most of the previous `docker-compose` features and flags. You can test the Compose V2 by simply replacing the dash (`-`) with a space, and by running `docker compose`, instead of `docker-compose`.
{: .important}

Starting with Docker Desktop 3.4.0, you can run Compose V2 commands without modifying your invocations, by enabling the drop-in replacement of the previous `docker-compose` with the new command.  See the section [Installing Compose V2](#installing-compose-v2) for detailed instructions how to enable the drop-in replacement.

## Context of Docker Compose evolution

Introduction of the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"} makes a clean distinction between the Compose YAML file model and the `docker-compose` implementation. Making this change has enabled a number of enhancements, including adding the `compose` command directly into the Docker CLI,  being able to “up” a Compose application on cloud platforms by simply switching the Docker context, and launching of [Amazon ECS](/cloud/ecs-integration) and [Microsoft ACI](/cloud/aci-integration). As the Compose specification evolves, new features land faster in the Docker CLI.

While `docker-compose` is still supported and maintained, Compose V2 implementation relies directly on the compose-go bindings which are maintained as part of the specification. This allows us to include community proposals, experimental implementations by the Docker CLI and/or Engine, and deliver features faster to users. Compose V2 also  supports some of the newer additions to the Compose specification, such as [profiles](profiles.md) and [GPU](gpu-support.md) devices.

Additionally, Compose V2 also supports [Apple silicon](../desktop/mac/apple-silicon.md).

For more information about the flags that are supported in the new compose command, see the [docker-compose compatibility list](cli-command-compatibility.md).

## Transitioning to GA for Compose V2

We are currently working towards providing a standard way to install Compose V2 on Linux. When this is available, Compose V2 will be marked as Generally Available (GA).

**Compose V2 GA** means:

- New features and bug fixes will only be considered in the Compose V2 code base.
- Docker Compose V2 will be the default setting in Docker Desktop for Mac and Windows. You can still opt out through the Docker Desktop UI and the CLI. This means, when you run `docker-compose`, you will actually be running`docker compose`.
- Compose V2 will be included with the latest version of the Docker CLI. You can use [Compose Switch](#compose-switch) to redirect `docker-compose` to `docker compose`.
- [Compose V2 branch](https://github.com/docker/compose/tree/v2) will become the default branch.
- Docker Compose V1 will be maintained to address any security issues.

> **Important**
>
> We would like to make the Compose V2 transition to be as smooth as possible for all users. We currently don't have a concrete timeline to deprecate Compose V1. We will review the feedback from the community on the GA and the adoption on Linux, and come up with a plan to deprecate Compose V1. We are not planning to remove the aliasing of `docker-compose` to `docker compose`. We would like to make it easier for users to switch to V2 without breaking any existing scripts. We will follow up with a blog post with more information on the exact timeline on V1 deprecation and the end of support policies for security issues.
>
> Your feedback is important to us. Reach out to us and let us know your feedback on our [Public Roadmap](https://github.com/docker/roadmap/issues/257){:target="_blank" rel="noopener" class="_"}.
{: .important}

## Installing Compose V2

This section contains instructions on how to install Compose V2.

### Install on Mac and Windows

Docker Desktop for Mac and for Windows version 3.2.1 and above includes the new Compose command along with the Docker CLI. Therefore, Windows and Mac users do not need to install Compose V2 separately.

We will progressively turn Docker Compose V2 on automatically for Docker Desktop users, so that users can seamlessly move to Docker Compose V2 without the need to change any of their scripts. If you run into any problems with Compose V2, you can simply switch back to Compose v1, either in Docker Desktop, or in the CLI.

For Docker Desktop installation instructions, see:

- [Install Docker Desktop on Mac](../desktop/mac/install.md)
- [Install Docker Desktop on Windows](../desktop/windows/install.md)

To disable Docker Compose V2 using Docker Desktop:

1. From the Docker menu, click **Preferences** (**Settings** on Windows) > **General**.
2. Clear the **Use Docker Compose V2** check box.

To disable Docker Compose V2 using the CLI, run:

```console
$ docker-compose disable-v2
```

### Install on Linux

You can install Compose V2 by downloading the appropriate binary for your system
from the [project release page](https://github.com/docker/compose/releases){:target="_blank" rel="noopener" class="_"} and copying it into `$HOME/.docker/cli-plugins` as `docker-compose`.

1. Run the following command to download the current stable release of Docker Compose:

    ```console
    $ mkdir -p ~/.docker/cli-plugins/
    $ curl -SL https://github.com/docker/compose/releases/download/v2.0.1/docker-compose-linux-x86_64 -o ~/.docker/cli-plugins/docker-compose
    ```

    This command installs Compose V2 for the active user under `$HOME` directory. To install Docker Compose for all users on your system, replace `~/.docker/cli-plugins` with `/usr/local/lib/docker/cli-plugins`.

2. Apply executable permissions to the binary:

    ```console
    $ chmod +x ~/.docker/cli-plugins/docker-compose
    ```

3. Test your installation

    ```console
    $ docker compose version
    Docker Compose version 2.0.1
    ```

### Compose Switch

[Compose Switch](https://github.com/docker/compose-switch/){:target="_blank" rel="noopener" class="_"} is a replacement to the Compose V1 `docker-compose` (python) executable. Compose switch translates the command line into Compose V2 `docker compose` and then runs the latter.

To install Compose Switch automatically, run:

```console
$ curl -fL https://raw.githubusercontent.com/docker/compose-cli/main/scripts/install/install_linux.sh | sh
```

To install Compose Switch manually:

1. Download the `compose-switch` binary for your architecture

    ```console
    $ curl -fL https://github.com/docker/compose-switch/releases/download/v1.0.1/docker-compose-linux-amd64 -o /usr/local/bin/compose-switch
    ```

2. Run the following command to make it an executable:

    ```console
    $ chmod +x /usr/local/bin/compose-switch
    ```

3. Rename the `docker-compose` binary if you've already installed it as `/usr/local/bin/docker-compose`

    ```console
    $ mv /usr/local/bin/docker-compose /usr/local/bin/docker-compose-v1
    ```

4. Define an **alternatives** group for the `docker-compose` command:

    ```console
    $ update-alternatives --install /usr/local/bin/docker-compose docker-compose <PATH_TO_DOCKER_COMPOSE_V1> 1
    $ update-alternatives --install /usr/local/bin/docker-compose docker-compose /usr/local/bin/compose-switch 99
    ```

5. Verify your installation:

    ```console
    $ update-alternatives --display docker-compose
    docker-compose - auto mode
        link best version is /usr/local/bin/compose-switch
        link currently points to /usr/local/bin/compose-switch
        link docker-compose is /usr/local/bin/docker-compose
    /usr/bin/docker-compose - priority 1
    /usr/local/bin/compose-switch - priority 99
    ```

#### Uninstall Docker Compose

If you installed Docker Compose using curl, run the following command to uninstall:

```console
$ sudo rm ~/.docker/cli-plugins/docker-compose
```
