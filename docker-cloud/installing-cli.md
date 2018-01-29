---
description: Using the Docker Cloud CLI on Linux, Windows, and macOS, installing, updating, uninstall
keywords: cloud, command-line, CLI
redirect_from:
- /docker-cloud/getting-started/intermediate/installing-cli/
- /docker-cloud/getting-started/installing-cli/
- /docker-cloud/tutorials/installing-cli/
title: The Docker Cloud CLI
---

Docker Cloud maintains a Command Line Interface (CLI) tool that you can use
to interact with the service. We highly recommend installing the CLI, as it
allows you to script and automate actions in Docker Cloud without using the web
interface. If you only ever use the web interface, this is not necessary.

## Install

Install the docker-cloud CLI either by running a Docker container, or by using the package manager for your system.

#### Run the CLI in a Docker container

If you have Docker Engine installed locally, you can run the following `docker`
command in your shell regardless of which operating system you are using.

```none
$ docker run dockercloud/cli -h
```

This command runs the `docker-cloud` CLI image in a container for you. Learn
more about how to use this container
[here](https://github.com/docker/dockercloud-cli#docker-image).

#### Install for Linux or Windows

You can install the CLI locally using the [pip](https://pip.pypa.io/en/stable/)
package manager, which is a package manager for
[Python](https://www.python.org/) applications.

* If you already have Python 2.x or 3.x installed, you probably have `pip` and
`setuptools`, but need to upgrade per the instructions
[here](https://packaging.python.org/installing/).

> The Docker Cloud CLI does not currently support Python 3.x.
>
> we recommend using Python 2.x. To learn more,
see the Python and CLI issues described in
[Known issues in Docker Cloud](/docker-cloud/docker-errors-faq.md).

* If you do not have Python or `pip` installed, you can either [install
Python](https://wiki.python.org/moin/BeginnersGuide/Download) or use this
[standalone pip
installer](https://pip.pypa.io/en/latest/installing/#installing-with-get-pip-py). You do not need Python for our purposes, just `pip`.

Now that you have `pip`, open a shell or terminal
window and run the following command to install the docker-cloud CLI:

```bash
$ pip install docker-cloud
```

If you encounter errors on Linux machines, make sure that `python-dev` is
installed. For example, on Ubuntu, run the following command:

```
$ apt-get install python-dev
```

#### Install on macOS

We recommend installing Docker CLI for macOS using Homebrew. If you don't have
`brew` installed, follow the instructions here: [http://brew.sh](http://brew.sh){: target="_blank" class="_"}

Once Homebrew is installed, open Terminal and run the following command:

```bash
$ brew install docker-cloud
```

> **Note**: You can also use [pip](https://pip.pypa.io/en/stable/) to install on macOS, but we suggest Homebrew since it is a package manager designed for the
Mac.

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

```none
$ docker login
Username: user
Password:
Email: user@example.org
Login succeeded!
```

#### What's next?

See the [Developer documentation](/apidocs/docker-cloud.md) for more information on using the CLI and our APIs.


## Use the docker-cloud CLI with an organization

When you use the docker-cloud CLI, it authenticates against the Docker Cloud
service with the user credentials saved by the `docker login` command. To use
the CLI to interact with objects belonging to an [Organization](orgs.md), you
must override the `DOCKERCLOUD_NAMESPACE` environment variable that sets this
user.

For example:

```none
$ export DOCKERCLOUD_NAMESPACE=myorganization
```

You can also set the `DOCKERCLOUD_NAMESPACE` variable before each CLI command.
For example:

```none
$ DOCKERCLOUD_NAMESPACE=myteam docker container ps
```

To learn more, see the [Docker Cloud CLI README](https://github.com/docker/dockercloud-cli#namespace).


## Upgrade the docker-cloud CLI

Periodically, Docker adds new features and fixes bugs in the existing CLI. To use these new features, you must upgrade the CLI.

#### Upgrade the docker-cloud CLI on Linux or Windows

```none
$ pip install -U docker-cloud
```

#### Upgrade the docker-cloud CLI on macOS

```none
$ brew update && brew upgrade docker-cloud
```

## Uninstall the docker-cloud CLI

If you are having trouble using the docker-cloud CLI, or find that it conflicts
with other applications on your system, you may want to uninstall and reinstall.

#### Uninstall on Linux or Windows

Open your terminal or command shell and execute the following command:

```none
$ pip uninstall docker-cloud
```

#### Uninstall on macOS

Open your Terminal application and execute the following command:

```none
$ brew uninstall docker-cloud
```
