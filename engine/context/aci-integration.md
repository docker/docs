---
title: Deploying Docker containers on Azure
description: Deploying Docker containers on Azure
keywords: Docker, Azure, Integration, ACI, context, Compose, cli, deploy, containers, cloud
toc_min: 1
toc_max: 2
---

## Overview

The Docker Azure Integration enables developers to use native Docker commands to run applications in Azure Container Instances (ACI) when building cloud-native applications. The new experience provides a tight integration between Docker Desktop and Microsoft Azure allowing developers to quickly run applications using the Docker CLI or VS Code extension, to switch seamlessly from local development to cloud deployment.

In addition, the integration between Docker and Microsoft developer technologies allow developers to use the Docker CLI to:

- Easily log into Azure
- Set up an ACI context in one Docker command allowing you to switch from a local context to a cloud context and run applications quickly and easily
- Simplify single container and multi-container application development using the Compose specification, allowing a developer to invoke fully Docker-compatible commands seamlessly for the first time natively within a cloud container service

>**Note**
>
> Docker Azure Integration is currently a beta release. The commands and flags are subject to change in subsequent releases.
{:.important}

## Prerequisites

To deploy Docker containers on Azure, you must meet the following requirements:

1. Download and install Docker Desktop Edge version 2.3.2.0 or later.

    - [Download for Mac](https://desktop.docker.com/mac/edge/Docker.dmg){: target="_blank" class="_"}
    - [Download for Windows](https://desktop.docker.com/win/edge/Docker%20Desktop%20Installer.exe){: target="_blank" class="_"}

    Alternatively, install the [Docker ACI Integration for Linux](#install-the-docker-aci-integration-cli-on-linux).

2. Ensure you have an Azure subscription. You can get started with an [Azure free account](https://aka.ms/AA8r2pj){: target="_blank" class="_"}.

## Run Docker containers on ACI

Docker not only runs containers locally, but also enables developers to seamlessly deploy Docker containers on ACI using `docker run` or deploy multi-container applications defined in a Compose file using the `docker compose up` command.

The following sections contain instructions on how to deploy your Docker containers on ACI.

### Log into Azure

Run the following commands to log into Azure:

```console
docker login azure
```

This opens your web browser and prompts you to enter your Azure login credentials.

### Create an ACI context

After you have logged in, you need to create a Docker context associated with ACI to deploy containers in ACI. For example, let us create a new context called `myacicontext`:

```console
docker context create aci myacicontext
```

This command automatically uses your Azure login credentials to identify your subscription IDs and resource groups. You can then interactively select the subscription and group that you would like to use. If you prefer, you can specify these options in the CLI using the following flags: `--subscription-id`,
`--resource-group`, and `--location`.

If you don't have any existing resource groups in your Azure account, the `docker context create aci myacicontext` command creates one for you. You don’t have to specify any additional options to do this.

After you have created an ACI context, you can list your Docker contexts by running the `docker context ls` command:

```console
NAME                TYPE                DESCRIPTION                               DOCKER ENDPOINT                KUBERNETES ENDPOINT   ORCHESTRATOR
myacicontext        aci                 myResourceGroupGTA@eastus
default *           moby              Current DOCKER_HOST based configuration   unix:///var/run/docker.sock                          swarm
```

> **Note**
>
> If you need to change the subscription and create a new context, you must 
execute the `docker login azure` command again.

### Run a container

Now that you've logged in and created an ACI context, you can start using Docker commands to deploy containers on ACI.

There are two ways to use your new ACI context. You can use the `--context` flag with the Docker command to specify that you would like to run the command using your newly created ACI context.

```console
docker --context myacicontext run -p 80:80 nginx
```

Or, you can change context using `docker context use` to select the ACI context to be your focus for running Docker commands. For example, we can use the `docker context use` command to deploy an ngnix container:

```console
docker context use myacicontext
docker run -p 80:80 nginx
```

After you've switched to the `myacicontext` context, you can use docker ps to list your containers running on ACI.

To view logs from your container, run:

```console
docker logs <CONTAINER_ID>
```

To execute a command in a running container, run:

```console
docker exec -t <CONTAINER_ID> COMMAND
```

To stop and remove a container from ACI, run:

```console
docker rm <CONTAINER_ID>
```

## Running Compose applications

You can also deploy and manage multi-container applications defined in Compose files to ACI using the `docker compose` command. To do this:

1. Ensure you are using your ACI context. You can do this either by specifying the `--context myacicontext` flag or by setting the default context using the command  `docker context use myacicontext`.

2. Run `docker compose up` and `docker compose down` to start and then stop a full Compose application.

  By default, `docker compose up` uses the `docker-compose.yaml` file in the current folder. You can specify the working directory using the  --workdir  flag or specify the Compose file directly using the `--file` flag.

  You can also specify a name for the Compose application using the `--project-name` flag during deployment. If no name is specified, a name will be derived from the working directory.

  You can view logs from containers that are part of the Compose application using the command `docker logs <CONTAINER_ID>`. To know the container ID, run `docker ps`.

> **Note**
>
> The current Docker Azure integration does not allow fetching a combined log stream from all the containers that make up the Compose application.

## Using ACI resource groups as namespaces

You can create several Docker contexts associated with ACI. Each context must be associated with a unique Azure resource group. This allows you to use Docker contexts as namespaces. You can switch between namespaces using `docker context use <CONTEXT>`.

When you run the `docker ps` command, it only lists containers in your current Docker context. There won’t be any contention in container names or Compose application names between two Docker contexts.

## Install the Docker ACI Integration CLI on Linux

The Docker ACI Integration CLI adds support for running and managing containers on Azure Container Instances (ACI).

>**Note**
>
> **Docker Azure Integration is a beta release**. The installation process, commands, and flags will change in future releases.
{:.important}

### Prerequisites

* [Docker 19.03 or later](https://docs.docker.com/get-docker/)

### Install script

You can install the new CLI using the install script:

```console
curl -L https://raw.githubusercontent.com/docker/aci-integration-beta/main/scripts/install_linux.sh | sh
```

### Manual install

You can download the Docker ACI Integration CLI from the 
[latest release](https://github.com/docker/aci-integration-beta/releases/latest){: target="_blank" class="_"} page.

You will then need to make it executable:

```console
chmod +x docker-aci
```

To enable using the local Docker Engine and to use existing Docker contexts, you
must have the existing Docker CLI as `com.docker.cli` somewhere in your
`PATH`. You can do this by creating a symbolic link from the existing Docker
CLI:

```console
ln -s /path/to/existing/docker /directory/in/PATH/com.docker.cli
```

> **Note**
>
> The `PATH` environment variable is a colon-separated list of
> directories with priority from left to right. You can view it using
> `echo $PATH`. You can find the path to the existing Docker CLI using
> `which docker`. You may need root permissions to make this link.

On a fresh install of Ubuntu 20.04 with Docker Engine
[already installed](https://docs.docker.com/engine/install/ubuntu/):

```console
$ echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin
$ which docker
/usr/bin/docker
$ sudo ln -s /usr/bin/docker /usr/local/bin/com.docker.cli
```

You can verify that this is working by checking that the new CLI works with the
default context:

```console
$ ./docker-aci --context default ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
$ echo $?
0
```

To make this CLI with ACI integration your default Docker CLI, you must move it
to a directory in your `PATH` with higher priority than the existing Docker CLI.

Again, on a fresh Ubuntu 20.04:

```console
$ which docker
/usr/bin/docker
$ echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin
$ sudo mv docker-aci /usr/local/bin/docker
$ which docker
/usr/local/bin/docker
$ docker version
...
 Azure integration  0.1.4
...
```

### Supported commands

After you have installed the Docker ACI Integration CLI, run `--help` to see the current list of commands.

> **Note**
>
> Docker Azure Integration is a beta release. The commands and flags will change in future releases.
{:.important}

### Uninstall

To remove the Docker Azure Integration CLI, you need to remove the binary you downloaded and `com.docker.cli` from your `PATH`. If you installed using the script, this can be done as follows:

```console
sudo rm /usr/local/bin/docker /usr/local/bin/com.docker.cli
```

## Feedback

Thank you for trying out the Docker Azure Integration beta release. Your feedback is very important to us. Let us know your feedback by creating an issue in the [aci-integration-beta](https://github.com/docker/aci-integration-beta){: target="_blank" class="_"} GitHub repository.
