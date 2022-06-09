---
description: Instructions for installing Docker Engine on CentOS
keywords: requirements, apt, installation, centos, rpm, install, uninstall, upgrade, update
redirect_from:
- /ee/docker-ee/centos/
- /engine/installation/centos/
- /engine/installation/linux/centos/
- /engine/installation/linux/docker-ce/centos/
- /engine/installation/linux/docker-ee/centos/
- /install/linux/centos/
- /install/linux/docker-ce/centos/
- /install/linux/docker-ee/centos/
title: Install Docker Engine on CentOS
toc_max: 4
---

To get started with Docker Engine on CentOS, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#installation-methods).

## Prerequisites

### OS requirements

To install Docker Engine, you need a maintained version of CentOS 7, CentOS 8 (stream),
or CentOS 9 (stream). Archived versions aren't supported or tested.

The `centos-extras` repository must be enabled. This repository is enabled by
default, but if you have disabled it, you need to
[re-enable it](https://wiki.centos.org/AdditionalResources/Repositories){: target="_blank" rel="noopener" class="_" }.

The `overlay2` storage driver is recommended.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```console
$ sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
```

It's OK if `yum` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker Engine package is now called `docker-ce`.

## Installation methods

You can install Docker Engine in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the RPM package and
  [install it manually](#install-from-a-package) and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

- In testing and development environments, some users choose to use automated
  [convenience scripts](#install-using-the-convenience-script) to install Docker.

### Install using the repository

Before you install Docker Engine for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker
from the repository.

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/centos" %}

Install the `yum-utils` package (which provides the `yum-config-manager`
utility) and set up the repository.

```console
$ sudo yum install -y yum-utils

$ sudo yum-config-manager \
    --add-repo \
    {{ download-url-base }}/docker-ce.repo
```

#### Install Docker Engine

1.  Install the _latest version_ of Docker Engine, containerd, and Docker Compose
    or go to the next step to install a specific version:

    ```console
    $ sudo yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin
    ```

    If prompted to accept the GPG key, verify that the fingerprint matches
    `060A 61C5 1B55 8A7F 742B 77AA C52F EB6B 621E 9F35`, and if so, accept it.

    This command installs Docker, but it doesn't start Docker. It also creates a
    `docker` group, however, it doesn't add any users to the group by default.

2.  To install a _specific version_ of Docker Engine, list the available versions
    in the repo, then select and install:

    a. List and sort the versions available in your repo. This example sorts
       results by version number, highest to lowest, and is truncated:

    ```console
    $ yum list docker-ce --showduplicates | sort -r

    docker-ce.x86_64  3:18.09.1-3.el7                     docker-ce-stable
    docker-ce.x86_64  3:18.09.0-3.el7                     docker-ce-stable
    docker-ce.x86_64  18.06.1.ce-3.el7                    docker-ce-stable
    docker-ce.x86_64  18.06.0.ce-3.el7                    docker-ce-stable
    ```

    The list returned depends on which repositories are enabled, and is specific
    to your version of CentOS (indicated by the `.el7` suffix in this example).

    b. Install a specific version by its fully qualified package name, which is
       the package name (`docker-ce`) plus the version string (2nd column)
       starting at the first colon (`:`), up to the first hyphen, separated by
       a hyphen (`-`). For example, `docker-ce-18.09.1`.

    ```console
    $ sudo yum install docker-ce-<VERSION_STRING> docker-ce-cli-<VERSION_STRING> containerd.io docker-compose-plugin
    ```

    This command installs Docker, but it doesn't start Docker. It also creates a
    `docker` group, however, it doesn't add any users to the group by default.

3.  Start Docker.

    ```console
    $ sudo systemctl start docker
    ```

4.  Verify that Docker Engine is installed correctly by running the `hello-world`
    image.

    ```console
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints a message and exits.

This installs and runs Docker Engine. Use `sudo` to run Docker
commands. Continue to [Linux postinstall](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker Engine

To upgrade Docker Engine, follow the [installation instructions](#install-using-the-repository),
choosing the new version you want to install.

### Install from a package

If you cannot use Docker's repository to install Docker, you can download the
`.rpm` file for your release and install it manually. You need to download
a new file each time you want to upgrade Docker Engine.

1.  Go to [{{ download-url-base }}/]({{ download-url-base }}/){: target="_blank" rel="noopener" class="_" }
    and choose your version of CentOS. Then browse to `x86_64/stable/Packages/`
    and download the `.rpm` file for the Docker version you want to install.

2.  Install Docker Engine, changing the path below to the path where you downloaded
    the Docker package.

    ```console
    $ sudo yum install /path/to/package.rpm
    ```

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

3.  Start Docker.

    ```console
    $ sudo systemctl start docker
    ```

4.  Verify that Docker Engine is installed correctly by running the `hello-world`
    image.

    ```console
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints a message and exits.

This installs and runs Docker Engine. Use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker Engine

To upgrade Docker Engine, download the newer package file and repeat the
[installation procedure](#install-from-a-package), using `yum -y upgrade`
instead of `yum -y install`, and point to the new file.

{% include install-script.md %}

## Uninstall Docker Engine

1.  Uninstall the Docker Engine, CLI, Containerd, and Docker Compose packages:

    ```console
    $ sudo yum remove docker-ce docker-ce-cli containerd.io docker-compose-plugin
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```console
    $ sudo rm -rf /var/lib/docker
    $ sudo rm -rf /var/lib/containerd
    ```

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md).
- Review the topics in [Develop with Docker](../../develop/index.md) to learn how to build new applications using Docker.
