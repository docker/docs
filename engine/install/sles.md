---
description: Instructions for installing Docker Engine on SLES
keywords: requirements, apt, installation, centos, rpm, sles, install, uninstall, upgrade, update, s390x, ibm-z
redirect_from:
- /ee/docker-ee/sles/
- /ee/docker-ee/suse/
- /engine/installation/linux/docker-ce/sles/
- /engine/installation/linux/docker-ee/sles/
- /engine/installation/linux/docker-ee/suse/
- /engine/installation/linux/sles/
- /engine/installation/linux/SUSE/
- /engine/installation/linux/suse/
- /engine/installation/sles/
- /engine/installation/SUSE/
- /install/linux/docker-ce/sles/
- /install/linux/docker-ee/sles/
- /install/linux/docker-ee/suse/
- /install/linux/sles/
- /installation/sles/
title: Install Docker Engine on SLES
toc_max: 4
---

To get started with Docker Engine on SLES, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#installation-methods).

## Prerequisites

> **Note**
>
> We currently only provide packages for SLES on s390x (IBM Z). Other architectures
> are not yet supported for SLES.

### OS requirements

To install Docker Engine, you need a maintained version of SLES 15-SP2 or SLES 15-SP3 on s390x (IBM Z).
Archived versions aren't supported or tested.

The [`SCC SUSE`](https://scc.suse.com/packages?name=SUSE%20Linux%20Enterprise%20Server&version=15.2&arch=s390x)
repositories must be enabled. 

The `SELinux` repository must be enabled. This repository is not added by
default, you need to
[add it](https://download.opensuse.org/repositories/security).

```console
$ sudo zypper install -y lsb-release
$ sles_version=$(lsb_release -r | rev | cut -c1)
$ opensuse_repo="https://download.opensuse.org/repositories/security:SELinux/SLE_15_SP$sles_version/security:SELinux.repo"
$ sudo zypper addrepo $opensuse_repo 
```

The `overlay2` storage driver is recommended.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```console
$ sudo zypper remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine \
                  runc
```

It's OK if `zypper` reports that none of these packages are installed.

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

{% assign download-url-base = "https://download.docker.com/linux/sles" %}

Set up the **stable** repository.

```console
$ sudo zypper \
    addrepo \
    {{ download-url-base }}/docker-ce.repo
```

> **Optional**: Enable the **nightly** or **test** repositories.
>
> These repositories are included in the `docker.repo` file above but are disabled
> by default. You can enable them alongside the stable repository.  The following
> command enables the **nightly** repository.
>
> ```console
> $ sudo zypper mr -e docker-ce-nightly
> ```
>
> To enable the **test** channel, run the following command:
>
> ```console
> $ sudo zypper mr -e docker-ce-test
> ```
>
> You can disable the **nightly** or **test** repository by running the
> ```console
> $ sudo zypper mr -d docker-ce-nightly
> $ sudo zypper mr -d docker-ce-test
> ```
>
> [Learn about **nightly** and **test** channels](index.md).

#### Install Docker Engine

1.  Install the _latest version_ of Docker Engine and containerd, or go to the next step to install a specific version:

    ```console
    $ sudo zypper install docker-ce docker-ce-cli containerd.io
    ```

    If prompted to accept the GPG key, verify that the fingerprint matches
    `060A 61C5 1B55 8A7F 742B 77AA C52F EB6B 621E 9F35`, and if so, accept it.

    > Got multiple Docker repositories?
    >
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `zypper install` or
    > `zypper update` command always installs the highest possible version,
    > which may not be appropriate for your stability needs.

    This command installs Docker, but it doesn't start Docker. It also creates a
    `docker` group, however, it doesn't add any users to the group by default.

2.  To install a _specific version_ of Docker Engine, list the available versions
    in the repo, then select and install:

    a. List and sort the versions available in your repo. This example sorts
       results by version number, highest to lowest, and is truncated:

    ```console
    $ sudo zypper search -s --match-exact docker-ce | sort -r
    
      v  | docker-ce | package | 3:20.10.7-3 | s390x | Docker CE Stable - s390x
      v  | docker-ce | package | 3:20.10.8-3 | s390x | Docker CE Stable - s390x
    ```

    The list returned depends on which repositories are enabled, and is specific
    to your version of SLES.

    b. Install a specific version using the version string from the fourth column, For example, `docker-ce-3:20.10.7`.

    ```console
    $ sudo zypper install docker-ce-<VERSION_STRING> docker-ce-cli-<VERSION_STRING> containerd.io
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
    and choose your version of SLES. Then browse to `15/s390x/stable/Packages/`
    and download the `.rpm` file for the Docker version you want to install.

    > **Note**
    >
    > To install a **nightly** or **test** (pre-release) package,
    > change the word `stable` in the above URL to `nightly` or `test`.
    > [Learn about **nightly** and **test** channels](index.md).

2.  Install Docker Engine, changing the path below to the path where you downloaded
    the Docker package.

    ```console
    $ sudo zypper install /path/to/package.rpm
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
[installation procedure](#install-from-a-package), using `zypper -y upgrade`
instead of `zypper -y install`, and point to the new file.

{% include install-script.md %}

## Uninstall Docker Engine

1.  Uninstall the Docker Engine, CLI, and Containerd packages:

    ```console
    $ sudo zypper remove docker-ce docker-ce-cli containerd.io
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
