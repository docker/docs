---
description: Set up the application
keywords: Python, application, setup
redirect_from:
- /docker-cloud/getting-started/python/2_set_up/
- /docker-cloud/getting-started/golang/2_set_up/
title: Set up your environment
---

In this step you install the Docker Cloud CLI to interact with the service using your command shell. This tutorial uses CLI commands to complete actions.

## Install the Docker Cloud CLI

Install the docker-cloud CLI using the package manager for your system.

#### Run the CLI in a Docker container

If you have Docker Engine installed locally, you can run the following `docker`
command in your shell regardless of which operating system you are using.

```none
docker run dockercloud/cli -h
```

This command runs the `docker-cloud` CLI image in a container for you. Learn
more about how to use this container
[here](https://github.com/docker/dockercloud-cli#docker-image).

#### Install for Linux or Windows

You can install the CLI locally using the [pip](https://pip.pypa.io/en/stable/)
package manager, which is a package manager for
[Python](https://www.python.org/) applications.

* If you already have 2.x or Python 3.x installed, you probably have `pip` and
`setuptools`, but need to upgrade per the instructions
[here](https://packaging.python.org/installing/).

  > The Docker Cloud CLI does not currently support Python 3.x.
  >
  > We recommend using Python 2.x with Docker Cloud. To learn more,
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

#### Install on macOS

We recommend installing Docker CLI for macOS using Homebrew. If you don't have `brew` installed, follow the instructions at [http://brew.sh](http://brew.sh){:target="_blank" class="_"}.

Once Homebrew is installed, open Terminal and run the following command:

```bash
$ brew install docker-cloud
```

> **Note**: You can also use [pip](https://pip.pypa.io/en/stable/) to install on macOS, but we suggest Homebrew since it is a package manager designed for the Mac.

## Validate the CLI installation
Check that the CLI installed correctly, using the `docker-cloud -v` command. (This command is the same for every platform.)

```bash
$ docker-cloud -v
docker-cloud 1.0.0
```

You can now use the `docker-cloud` CLI commands from your shell.

The documentation for the Docker Cloud CLI tool and API [here](/apidocs/docker-cloud.md).


## Log in

Use the `login` CLI command to log in to Docker Cloud. Use the username and password you used when creating your Docker ID. If you use Docker Hub, you can use the same username and password you use to log in to Docker Hub.

```
$ docker login
Username: my-username
Password:
Login succeeded!
```

You must log in to continue this tutorial.

## Set your username as an environment variable

For simplicity in this tutorial, we use an environment variable for your Docker Cloud username. If you plan to copy and paste the tutorial commands, set the environment variable using the command below. (Change `my-username` to your username.)

If you don't want to do this, make sure you substitute your username for $DOCKER_ID_USER whenever you see it in the example commands.

```none
$ export DOCKER_ID_USER=my-username
```

**If you are running the tutorial with an organization's resources:**

By default, the `docker-cloud` CLI uses your default user namespace, meaning the
repositories, nodes, and services associated with your individual Docker ID
account name. To use the CLI to interact with objects that belong to an
[organization](../../orgs.md), prefix these commands with
`DOCKERCLOUD_NAMESPACE=my-organization`, or set this variable as in the example below.

```none
$ export DOCKERCLOUD_NAMESPACE=my-organization
```


 See the [CLI documentation](../../installing-cli.md#use-the-docker-cloud-cli-with-an-organization) for more information.


Next up, we [prepare the app](3_prepare_the_app.md).
