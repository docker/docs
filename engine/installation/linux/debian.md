---
description: Instructions for installing Docker on Debian
keywords: Docker, Docker documentation, requirements, apt, installation, debian, raspbian, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/debian/
- /engine/installation/linux/raspbian/
title: Get Docker for Debian or Raspbian
---

To get started with Docker on Debian, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### OS requirements

To install Docker, you need the 64-bit version of one of these Debian or
Raspbian versions:

- Stretch (testing)
- Jessie 8.0 (LTS) / Raspbian Jessie
- Wheezy 7.7 (LTS)

#### Extra steps for Wheezy 7.7

- You need at least version 3.10 of the Linux kernel. Debian Wheezy ships with
  version 3.2, so you may need to
  [update the kernel](https://wiki.debian.org/HowToUpgradeKernel){: target="_blank" class="_" }.
  To check your kernel version:

  ```bash
  $ uname -r
  ```

- Enable the `backports` repository. See the
  [Debian documentation](https://backports.debian.org/Instructions/){: target="_blank" class"_"}.

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

    **Jessie or Stretch**:

    ```bash
    $ sudo apt-get install -y --no-install-recommends \
         apt-transport-https \
         ca-certificates \
         curl \
         software-properties-common
    ```

    **Wheezy**:

    ```bash
    $ sudo apt-get install -y --no-install-recommends \
         apt-transport-https \
         ca-certificates \
         curl \
         python-software-properties
    ```

2.  Add Docker's official GPG key:

    ```bash
    $ curl -fsSL https://apt.dockerproject.org/gpg | sudo apt-key add -
    ```

    Verify that the key ID is `58118E89F3A912897C070ADBF76221572C52609D`.

    ```bash
    $ sudo apt-key fingerprint 58118E89F3A912897C070ADBF76221572C52609D

      pub   4096R/2C52609D 2015-07-14
            Key fingerprint = 5811 8E89 F3A9 1289 7C07  0ADB F762 2157 2C52 609D
      uid                  Docker Release Tool (releasedocker) <docker@docker.com>
    ```

3.  Use the following command to set up the **stable** repository.

    > **Note**: The `lsb_release -cs` sub-command below returns the name of your
    > Debian distribution, such as `jessie`.

    - **Debian**:

      ```bash
      $ sudo add-apt-repository \
             "deb https://apt.dockerproject.org/repo/ \
             debian-$(lsb_release -cs) \
             main"
      ```

    - **Raspbian**:

      ```bash
      $ sudo add-apt-repository \
             "deb https://apt.dockerproject.org/repo/ \
             raspbian-$(lsb_release -cs) \
             main"
      ```

    To enable the `testing` repository, you can edit `/etc/apt/sources.list`
    and add the word `testing` after `main` on the appropriate line of the file.
    **Do not use unstable repositories on production systems or for non-testing workloads.**

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
    > updating to the latest version of Docker by not specifying a version in
    > the `apt-get install` or `apt-get update` command will always install the
    > highest possible version, which will almost certainly be an unstable one.

3.  On production systems, you should install a specific version of Docker
    instead of always using the latest. This output is truncated. List the
    available versions:

    ```bash
    $ apt-cache madison docker-engine
    docker-engine | 1.13.0-0~stretch | https://apt.dockerproject.org/repo debian-stretch/main amd64 Packages
    docker-engine | 1.12.3-0~stretch | https://apt.dockerproject.org/repo debian-stretch/main amd64 Packages
    docker-engine | 1.12.2-0~stretch | https://apt.dockerproject.org/repo debian-stretch/main amd64 Packages
    docker-engine | 1.12.1-0~stretch | https://apt.dockerproject.org/repo debian-stretch/main amd64 Packages
    ```

    The contents of the list depend upon which repositories are enabled,
    and will be specific to your version of Debian (indicated by the `stretch`
    suffix on the version, in this example). Choose a specific version to
    install. The second column is the version string. The third column is the
    repository name, which indicates which repository the package is from and
    by extension its stability level. To install a specific version, append the
    version string to the package name and separate them by an equals sign (`=`):

    ```bash
    $ sudo apt-get -y install docker-engine=<VERSION_STRING>
    ```

    On Debian, the Docker daemon starts automatically. On Raspbian, start Docker:

    ```bash
    $ sudo service docker start
    ```

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
    for your version of Debian or Raspbian.

    > **Note**: To install a testing version, change the word `main` in the
    > URL to `testing`. Do not use unstable versions of Docker in production
    > or for non-testing workloads.

2.  Install Docker, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo dpkg -i /path/to/package.deb
    ```

    On Debian, the Docker daemon starts automatically. On Raspbian, start Docker:

    ```bash
    $ sudo service docker start
    ```

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
