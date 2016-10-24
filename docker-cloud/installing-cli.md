---
aliases:
- /docker-cloud/getting-started/intermediate/installing-cli/
- /docker-cloud/getting-started/installing-cli/
- /docker-cloud/tutorials/installing-cli/
description: Using the Docker Cloud CLI on Linux, Windows, and macOS, installing,
  updating, uninstall
keywords:
- cloud, command-line, CLI
menu:
  main:
    parent: docker-cloud
title: The Docker Cloud CLI
---

# The Docker Cloud CLI

Docker Cloud maintains a Command Line Interface (CLI) tool that you can use
to interact with the service. We highly recommend installing the CLI, as it will
allow you to script and automate actions in Docker Cloud without using the web
interface. If you will only ever be using the web interface, this is not
necessary.

## Install

Install the docker-cloud CLI either by running a Docker container, or by using the package manager for your system.

#### Install using a Docker container

If you have Docker Engine installed locally, you can simply run the following command regardless of which operating system you are using.

```
docker run dockercloud/cli -h
```

This runs a container that installs the docker-cloud CLI for you. Learn more about this container [here](https://github.com/docker/dockercloud-cli#docker-image).

#### Install for Linux or Windows

Open your terminal or command shell and execute the following command:

```bash
$ pip install docker-cloud
```
(On Linux machines, ensure that python-dev is installed.)

#### Install on macOS

We recommend installing Docker CLI for macOS using Homebrew. If you don't have `brew` installed, follow the instructions here: <a href="http://brew.sh" target="_blank">http://brew.sh</a>

Once Homebrew is installed, open Terminal and run the following command:

```bash
$ brew install docker-cloud
```

#### Validate the installation

Check that the CLI installed correctly:

```bash
$ docker-cloud -v
docker-cloud 1.0.0
```

## Getting Started

First, you should log in using the `docker` CLI and the `docker login` command.
Your Docker ID, which you also use to log in to Docker Hub, is also used for
logging in to Docker Cloud.

```bash
$ docker login
Username: user
Password:
Email: user@example.org
Login succeeded!
```

#### What's next?

See the [Developer documentation](/apidocs/docker-cloud.md) for more information on using the CLI and our APIs.

## Upgrade the docker-cloud CLI

Periodically, Docker will add new features and fix bugs in the existing CLI. To use these new features, you must upgrade the CLI.

#### Upgrade on the docker-cloud CLI on Linux or Windows

```
$ pip install -U docker-cloud
```

#### Upgrade the docker-cloud CLI on macOS

```
$ brew update && brew upgrade docker-cloud
```

## Uninstall the docker-cloud CLI

If you are having trouble using the docker-cloud CLI, or find that it conflicts
with other applications on your system, you may want to uninstall and reinstall.

#### Uninstall on Linux or Windows

Open your terminal or command shell and execute the following command:

```
$ pip uninstall docker-cloud
```

#### Uninstall on macOS

Open your Terminal application and execute the following command:

```
$ brew uninstall docker-cloud
```
