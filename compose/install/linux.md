---
description: How to install Docker Compose on Linux
keywords: compose, orchestration, install, installation, docker, documentation
toc_max: 3
title: Install the Compose plugin
redirect_from:
- /compose/compose-plugin/
- /compose/compose-linux/
---

On this page you can find instructions on how to install the Compose plugin on Linux from the command line.

To install the Compose plugin on Linux, you can either:
- [Set up Docker's repository on your Linux system](#install-using-the-repository).
- [Install Compose manually](#install-the-plugin-manually).

> **Note**
>
> These instructions assume you already have Docker Engine and Docker CLI installed and now want to install the Compose plugin.  
For Compose standalone, see [Install Compose Standalone](other.md).

### Install using the repository

1. Set up the repository. Find distro-specific instructions in:

    [Ubuntu](../../engine/install/ubuntu.md/#set-up-the-repository) |
    [CentOS](../../engine/install/centos.md/#set-up-the-repository) |
    [Debian](../../engine/install/debian.md/#set-up-the-repository) |
    [Fedora](../../engine/install/fedora.md/#set-up-the-repository) |
    [RHEL](../../engine/install/rhel.md/#set-up-the-repository) |
    [SLES](../../engine/install/sles.md/#set-up-the-repository).

2. Update the package index, and install the latest version of Docker Compose:

    * For Ubuntu and Debian, run:

        ```console
        $ sudo apt-get update
        $ sudo apt-get install docker-compose-plugin
        ```
    * For RPM-based distros, run:

        ```console
        $ sudo yum update
        $ sudo yum install docker-compose-plugin
        ```

3.  Verify that Docker Compose is installed correctly by checking the version.

    ```console
    $ docker compose version
    Docker Compose version vN.N.N
    ```

Where `vN.N.N` is placeholder text standing in for the latest version.

#### Update Compose

To update the Compose plugin, run the following commands:

* For Ubuntu and Debian, run:

    ```console
    $ sudo apt-get update
    $ sudo apt-get install docker-compose-plugin
    ```
* For RPM-based distros, run:

    ```console
    $ sudo yum update
    $ sudo yum install docker-compose-plugin
    ```

### Install the plugin manually

> **Note**
>
> This option requires you to manage upgrades manually. We recommend setting up Docker's repository for easier maintenance.

1.  To download and install the Compose CLI plugin, run:

    ```console
    $ DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
    $ mkdir -p $DOCKER_CONFIG/cli-plugins
    $ curl -SL https://github.com/docker/compose/releases/download/{{site.compose_version}}/docker-compose-linux-x86_64 -o $DOCKER_CONFIG/cli-plugins/docker-compose
    ```

    This command downloads the latest release of Docker Compose (from the Compose releases repository) and installs Compose for the active user under `$HOME` directory.

    To install:
    * Docker Compose for _all users_ on your system, replace `~/.docker/cli-plugins` with `/usr/local/lib/docker/cli-plugins`.
    * A different version of Compose, substitute `{{site.compose_version}}` with the version of Compose you want to use.  


2. Apply executable permissions to the binary:

    ```console
    $ chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose
    ```
    or, if you chose to install Compose for all users:

    ```console
    $ sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
    ```

3. Test the installation.

    ```console
    $ docker compose version
    Docker Compose version {{site.compose_version}}
    ```
