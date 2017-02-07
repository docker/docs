---
description: Instructions for installing Docker on Ubuntu
keywords: Docker, Docker documentation, requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
- /engine/installation/linux/ubuntulinux/
title: Get Docker for Ubuntu
---

To get started with Docker on Ubuntu, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### OS requirements

To install Docker, you need the 64-bit version of one of these Ubuntu versions:

- Yakkety 16.10
- Xenial 16.04 (LTS)
- Trusty 14.04 (LTS)

### Recommended extra packages

You need `curl` if you don't have it. Unless you have a strong reason not to,
install the `linux-image-extra-*` packages, which allow Docker to use the `aufs`
storage drivers. **This applies to all versions of Ubuntu**.

```bash
$ sudo apt-get update

$ sudo apt-get install curl \
    linux-image-extra-$(uname -r) \
    linux-image-extra-virtual
```

## Install Docker

You can install Docker in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the DEB package and install it manually and manage
  upgrades completely manually.

- Some users cannot use the official Docker repositories, and must rely on
  the version of Docker that comes with their operating system. This version of
  Docker may be out of date. Those users should consult their operating system
  documentation and not follow these procedures.

### Install using the repository

Before you install Docker for the first time on a new host machine, you need to
set up the Docker repository. Afterward, you can install, update, or downgrade
Docker from the repository.

#### Set up the repository

1.  Install packages to allow `apt` to use a repository over HTTPS:

    ```bash
    $ sudo apt-get install apt-transport-https \
                           ca-certificates
    ```

2.  Add Docker's official GPG key:

    ```bash
    $ curl -fsSL https://yum.dockerproject.org/gpg | sudo apt-key add -
    ```

    > **Note**: The URL is correct, even for Linux distributions that use `APT`.

    Verify that the key ID is `58118E89F3A912897C070ADBF76221572C52609D`.

    ```bash
    $ apt-key fingerprint 58118E89F3A912897C070ADBF76221572C52609D

      pub   4096R/2C52609D 2015-07-14
            Key fingerprint = 5811 8E89 F3A9 1289 7C07  0ADB F762 2157 2C52 609D
      uid                  Docker Release Tool (releasedocker) <docker@docker.com>
    ```

3.  Use the following command to set up the **stable** repository. To also
    enable the **testing** repository, add the words `testing` after `main` on
    the last line.
    **Do not use these unstable repositories on production systems or for non-testing workloads.**

    ```bash
    $ sudo apt-get install software-properties-common
    $ sudo add-apt-repository \
           "deb https://apt.dockerproject.org/repo/ \
           ubuntu-$(lsb_release -cs) \
           main"
    ```

    To disable the `testing` repository, you can edit `/etc/apt/sources.list`
    and remove the word `testing` from the appropriate line in the file.
    
    > **Note**: Sometimes, in a distribution like Linux Mint, you might have
    > to change **ubuntu-$(lsb_release -cs)** to your parent Ubuntu distribution.
    > example: If you are using **Linux Mint Rafaela**, you could type in **ubuntu-trusty**


#### Install Docker

1.  Update the `apt` package index.

    ```bash
    $ sudo apt-get update
    ```

2.  Install the latest version of Docker, or go to the next step to install a
    specific version. Any existing installation of Docker is replaced.

    Use this command to install the latest version of Docker:

    ```bash
    $ sudo apt-get -y install docker-engine
    ```

    > **Warning**: If you have both stable and unstable repositories enabled,
    > installing or updating without specifying a version in the
    > `apt-get install` or `apt-get update` command will always install the
    > highest possible version, which will almost certainly be an unstable one.

3.  On production systems, you should install a specific version of Docker
    instead of always using the latest. This output is truncated. List the
    available versions.

    ```bash
    $ apt-cache madison docker-engine

    docker-engine | 1.13.0-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.12.3-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.12.2-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.12.1-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    ```

    The contents of the list depend upon which repositories are enabled,
    and will be specific to your version of Ubuntu (indicated by the `xenial`
    suffix on the version, in this example). Choose a specific version to
    install. The second column is the version string. The third column is the
    repository name, which indicates which repository the package is from and
    by extension its stability level. To install a specific version, append the
    version string to the package name and separate them by an equals sign (`=`):

    ```bash
    $ sudo apt-get -y install docker-engine=<VERSION_STRING>
    ```

    The Docker daemon starts automatically.

4.  Verify that `docker` is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker

To upgrade Docker, first run `sudo apt-get update`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker, you can download the
`.deb` file for your release and install it manually. You will need to download
a new file each time you want to upgrade Docker.

1.  Go to [https://apt.dockerproject.org/repo/pool/main/d/docker-engine/](https://apt.dockerproject.org/repo/pool/main/d/docker-engine/)
    and download the `.deb` file for the Docker version you want to install and
    for your version of Ubuntu.

    > **Note**: To install a testing version, change the word `main` in the
    > URL to `testing`. Do not use unstable versions of Docker in production
    > or for non-testing workloads.

2.  Install Docker, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo dpkg -i /path/to/package.deb
    ```

    The Docker daemon starts automatically.

3.  Verify that `docker` is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker

To upgrade Docker, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

## Uninstall Docker

1.  Uninstall the Docker package:

    ```bash
    $ sudo apt-get purge docker-engine
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```bash
    $ sudo rm -rf /var/lib/docker
    ```

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md)

- Continue with the [User Guide](../../userguide/index.md).
