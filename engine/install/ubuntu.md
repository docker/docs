---
description: Instructions for installing Docker Engine on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /ee/docker-ee/ubuntu/
- /engine/installation/linux/docker-ce/ubuntu/
- /engine/installation/linux/docker-ee/ubuntu/
- /engine/installation/linux/ubuntu/
- /engine/installation/linux/ubuntulinux/
- /engine/installation/ubuntulinux/
- /install/linux/docker-ce/ubuntu/
- /install/linux/docker-ee/ubuntu/
- /install/linux/ubuntu/
- /installation/ubuntulinux/
title: Install Docker Engine on Ubuntu
toc_max: 4
---

> **Docker Desktop for Linux**
>
> Docker Desktop helps you build, share, and run containers easily on Mac and Windows as you do on Linux. Docker handles the complex setup and allows you to focus on writing the code. Thanks to the positive support we received on the [subscription updates](https://www.docker.com/blog/updating-product-subscriptions/){: target="_blank" rel="noopener" class="_" id="dkr_docs_cta"}, we've started working on [Docker Desktop for Linux](https://www.docker.com/blog/accelerating-new-features-in-docker-desktop/){: target="_blank" rel="noopener" class="_" id="dkr_docs_cta"} which is the second-most popular feature request in our public roadmap. If you are interested in early access, sign up for our [Developer Preview program](https://www.docker.com/community/get-involved/developer-preview){: target="_blank" rel="noopener" class="_" id="dkr_docs_cta"}.
{: .important}

To get started with Docker Engine on Ubuntu, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#installation-methods).

## Prerequisites

### OS requirements

To install Docker Engine, you need the 64-bit version of one of these Ubuntu
versions:

- Ubuntu Impish 21.10
- Ubuntu Hirsute 21.04
- Ubuntu Focal 20.04 (LTS)
- Ubuntu Bionic 18.04 (LTS)

Docker Engine is supported on `x86_64` (or `amd64`), `armhf`, `arm64`, and `s390x` architectures.

> Ubuntu 16.04 LTS "Xenial Xerus" end-of-life
> 
> Ubuntu Linux 16.04 LTS reached the end of its five-year LTS window on April
> 30th 2021 and is no longer supported. Docker no longer releases packages for
> this distribution (including patch- and security releases). Users running
> Docker on Ubuntu 16.04 are recommended to update their system to a currently
> supported LTS version of Ubuntu.
{: .important }

### Uninstall old versions

Older versions of Docker were called `docker`, `docker.io`, or `docker-engine`.
If these are installed, uninstall them:

```console
$ sudo apt-get remove docker docker-engine docker.io containerd runc
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. If you do not need to save your existing data, and want to
start with a clean installation, refer to the [uninstall Docker Engine](#uninstall-docker-engine)
section at the bottom of this page.

### Supported storage drivers

Docker Engine on Ubuntu supports `overlay2`, `aufs` and `btrfs` storage drivers.

Docker Engine uses the `overlay2` storage driver by default. If you need to use
`aufs` instead, you need to configure it manually.
See [use the AUFS storage driver](../../storage/storagedriver/aufs-driver.md)

## Installation methods

You can install Docker Engine in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the DEB package and
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

{% assign download-url-base = "https://download.docker.com/linux/ubuntu" %}

1.  Update the `apt` package index and install packages to allow `apt` to use a
    repository over HTTPS:

    ```console
    $ sudo apt-get update

    $ sudo apt-get install \
        ca-certificates \
        curl \
        gnupg \
        lsb-release
    ```

2.  Add Docker's official GPG key:

    ```console
    $ curl -fsSL {{ download-url-base }}/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    ```

3.  Use the following command to set up the **stable** repository. To add the
    **nightly** or **test** repository, add the word `nightly` or `test` (or both)
    after the word `stable` in the commands below. [Learn about **nightly** and **test** channels](index.md).

    ```console
    $ echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] {{ download-url-base }} \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```

#### Install Docker Engine

1. Update the `apt` package index, and install the _latest version_ of Docker
   Engine and containerd, or go to the next step to install a specific version:

    ```console
    $ sudo apt-get update
    $ sudo apt-get install docker-ce docker-ce-cli containerd.io
    ```

    > Got multiple Docker repositories?
    >
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `apt-get install` or
    > `apt-get update` command always installs the highest possible version,
    > which may not be appropriate for your stability needs.

2.  To install a _specific version_ of Docker Engine, list the available versions
    in the repo, then select and install:

    a. List the versions available in your repo:

    ```console
    $ apt-cache madison docker-ce

      docker-ce | 5:18.09.1~3-0~ubuntu-xenial | {{ download-url-base }}  xenial/stable amd64 Packages
      docker-ce | 5:18.09.0~3-0~ubuntu-xenial | {{ download-url-base }}  xenial/stable amd64 Packages
      docker-ce | 18.06.1~ce~3-0~ubuntu       | {{ download-url-base }}  xenial/stable amd64 Packages
      docker-ce | 18.06.0~ce~3-0~ubuntu       | {{ download-url-base }}  xenial/stable amd64 Packages
    ```

    b. Install a specific version using the version string from the second column,
       for example, `5:18.09.1~3-0~ubuntu-xenial`.

    ```console
    $ sudo apt-get install docker-ce=<VERSION_STRING> docker-ce-cli=<VERSION_STRING> containerd.io
    ```

3.  Verify that Docker Engine is installed correctly by running the `hello-world`
    image.

    ```console
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints a message and exits.

Docker Engine is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](linux-postinstall.md) to allow non-privileged
users to run Docker commands and for other optional configuration steps.

#### Upgrade Docker Engine

To upgrade Docker Engine, first run `sudo apt-get update`, then follow the
[installation instructions](#install-using-the-repository), choosing the new
version you want to install.

### Install from a package

If you cannot use Docker's repository to install Docker Engine, you can download the
`.deb` file for your release and install it manually. You need to download
a new file each time you want to upgrade Docker.

1.  Go to [`{{ download-url-base }}/dists/`]({{ download-url-base }}/dists/){: target="_blank" rel="noopener" class="_" },
    choose your Ubuntu version, then browse to `pool/stable/`, choose `amd64`,
    `armhf`, `arm64`, or `s390x`, and download the `.deb` file for the Docker Engine
    version you want to install.

    > **Note**
    >
    > To install a **nightly** or **test** (pre-release) package,
    > change the word `stable` in the above URL to `nightly` or `test`.
    > [Learn about **nightly** and **test** channels](index.md).

2.  Install Docker Engine, changing the path below to the path where you downloaded
    the Docker package.

    ```console
    $ sudo dpkg -i /path/to/package.deb
    ```

    The Docker daemon starts automatically.

3.  Verify that Docker Engine is installed correctly by running the `hello-world`
    image.

    ```console
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints a message and exits.

Docker Engine is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker Engine

To upgrade Docker Engine, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

{% include install-script.md %}

## Uninstall Docker Engine

1.  Uninstall the Docker Engine, CLI, and Containerd packages:

    ```console
    $ sudo apt-get purge docker-ce docker-ce-cli containerd.io
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
