---
description: How to install Docker Compose on Linux
keywords: compose, orchestration, install, installation, docker, documentation
toc_max: 3

title: Install on Linux
redirect_from:
- /compose/compose-plugin/
---

On this page you can find instructions on how to install the Compose on Linux from the command line.

## Install Compose

To install Compose:
* Option 1: [Set up Docker's repository on your Linux system](#install-using-the-repository).
* Option 2: [Install Compose manually](#install-the-plugin-manually).

### Install using the repository

> **Note**
>
> These instructions assume you already have Docker Engine and Docker CLI installed and now want to install the Compose plugin.  
For Compose standalone, see [Install Compose Standalone](./compose-other.md#on-linux).

If you have already set up the Docker repository, jump to step 2.

1. Set up the repository in:
[Ubuntu](../../engine/install/ubuntu.md/#set-up-the-repository) |
[CentOS](../../engine/install/centos.md/#set-up-the-repository) |
[Debian](../../engine/install/fedora.md/#set-up-the-repository) |
[RHEL](../../engine/install/fedora.md/#set-up-the-repository) |
[SLES](../../engine/install/sles.md/#set-up-the-repository).

2. Update the `apt` package index, and install the _latest version_ of Docker Compose:

    > Or, if using a different distro, use the equivalent package manager instructions.

    ```console
    $ sudo apt-get update
    $ sudo apt-get install docker-compose-plugin
    ```

    Alternatively, to install a specific version of the Compose CLI plugin:

    a. List the versions available in your repo:

      ```console
      $ apt-cache madison docker-compose-plugin

      docker-compose-plugin | 2.6.0~ubuntu-bionic | https://download.docker.com/linux/ubuntu bionic/stable amd64 Packages
      docker-compose-plugin | 2.5.0~ubuntu-bionic | https://download.docker.com/linux/ubuntu bionic/stable amd64 Packages
      docker-compose-plugin | 2.3.3~ubuntu-bionic | https://download.docker.com/linux/ubuntu bionic/stable amd64 Packages
      ```

    b. From the list obtained use the version string you can in the second column to specify the version you wish to install.

    c. Install the selected version:


      ```console
      $ sudo apt-get install docker-compose-plugin=<VERSION_STRING>
      ```
    where `<VERSION_STRING>` is, for example,`2.3.3~ubuntu-focal`.

3.  Verify that Docker Compose is installed correctly by checking the version.

    ```console
    $ docker compose version
    Docker Compose version v2.3.3
    ```

> **Note**
>
> To run Compose as a non-root user, see [Manage Docker as a non-root user](../../engine/install/linux-postinstall.md){:target="_blank" rel="noopener" class="_"}.


### Install the plugin manually

> **Note**
>
> This option requires you to manage upgrades manually. We recommend setting up Docker's repository for an easier maintenance.

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

> **Note**
>
> To run Compose as a non-root user, see [Manage Docker as a non-root user](../../engine/install/linux-postinstall.md){:target="_blank" rel="noopener" class="_"}.


