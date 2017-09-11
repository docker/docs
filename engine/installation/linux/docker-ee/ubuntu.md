---
description: Instructions for installing Docker EE on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
- /engine/installation/linux/ubuntulinux/
title: Get Docker EE for Ubuntu
---

{% assign minor-version = "17.06" %}

To get started with Docker EE on Ubuntu, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ee).

## Prerequisites

Docker CE users should go to
[Get Docker CE for Ubuntu](/engine/installation/linux/docker-ce/ubuntu.md)
**instead of this topic**.

To install Docker Enterprise Edition (Docker EE), you need to know the Docker EE
repository URL associated with your trial or subscription. These instructions
work for Docker EE for Ubuntu and for Docker EE for Linux, which includes access
to Docker EE for all Linux distributions. To get this information:

- Go to [https://store.docker.com/my-content](https://store.docker.com/my-content).
- Each subscription or trial you have access to is listed. Click the **Setup**
  button for **Docker Enterprise Edition for Ubuntu**.
- Copy the URL from the field labeled
  **Copy and paste this URL to download your Edition**.

Use this URL when you see the placeholder text `<DOCKER-EE-URL>`.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### OS requirements

To install Docker EE, you need the 64-bit version of one of these Ubuntu versions:

- Xenial 16.04 (LTS)
- Trusty 14.04 (LTS)

Docker EE is supported on `x86_64` (or `amd64`) and `s390x` (IBM Z) architectures.

The only supported storage driver for Docker EE on Ubuntu is `aufs`.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. In addition,
if you are upgrading from Docker CE to Docker EE, remove the Docker CE package.

```bash
$ sudo apt-get remove docker docker-engine docker-ce docker.io
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker EE package is now called `docker-ee`.

### Extra packages for Trusty 14.04

Docker EE users must use the `aufs` storage driver on production systems. Install
the `linux-image-extra-*` packages, which allow Docker EE to use the `aufs`
storage driver.

```bash
$ sudo apt-get update

$ sudo apt-get install \
    linux-image-extra-$(uname -r) \
    linux-image-extra-virtual
```

## Install Docker EE

You can install Docker EE in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the DEB package and install it manually and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

### Install using the repository

Before you install Docker EE for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker EE
from the repository.

#### Set up the repository

1.  Update the `apt` package index:

    ```bash
    $ sudo apt-get update
    ```

2.  Install packages to allow `apt` to use a repository over HTTPS:

    ```bash
    $ sudo apt-get install \
        apt-transport-https \
        ca-certificates \
        curl \
        software-properties-common
    ```

3.  Add Docker's official GPG key using your customer Docker EE repository URL:

    ```bash
    $ curl -fsSL <DOCKER-EE-URL>/ubuntu/gpg | sudo apt-key add -
    ```

    Verify that you now have the key with the fingerprint
    `DD91 1E99 5A64 A202 E859  07D6 BC14 F10B 6D08 5F96`, by searching for the
    last 8 characters of the fingerprint.

    ```bash
    $ sudo apt-key fingerprint 6D085F96

    pub   4096R/0EBFCD88 2017-02-22
          Key fingerprint = DD91 1E99 5A64 A202 E859  07D6 BC14 F10B 6D08 5F96
    uid                  Docker Release (CE deb) <docker@docker.com>
    sub   4096R/6D085F96 2017-02-22
    ```

4.  Use the following command to set up the **stable** repository, replacing
    `<DOCKER-EE-URL>` with the URL you noted down in the
    [prerequisites](#prerequisites).

    > **Note**: The `lsb_release -cs` sub-command below returns the name of your
    > Ubuntu distribution, such as `xenial`.
    >

    **x86_64**:

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=amd64] <DOCKER-EE-URL>/ubuntu \
       $(lsb_release -cs) \
       stable-{{ minor-version }}"
    ```

    **s390x**:

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=s390x] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable-{{ minor-version }}"
    ```

#### Install Docker EE

1.  Update the `apt` package index.

    ```bash
    $ sudo apt-get update
    ```

2.  Install the latest version of Docker EE, or go to the next step to install a
    specific version. Any existing installation of Docker EE is replaced.

    Use this command to install the latest version of Docker EE:

    ```bash
    $ sudo apt-get install docker-ee
    ```

    > **Warning**: If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `apt-get install` or
    > `apt-get update` command will always install the highest possible version,
    > which may not be appropriate for your stability needs.
    {:.warning}

3.  On production systems, you should install a specific version of Docker EE
    instead of always using the latest. This output is truncated. List the
    available versions.

    ```bash
    $ apt-cache madison docker-ee

    docker-ee | {{ minor-version }}.0~ee-0~ubuntu-xenial | <DOCKER-EE-URL>/ubuntu xenial/stable amd64 Packages
    ```

    The contents of the list depend upon which repositories are enabled,
    and will be specific to your version of Ubuntu (indicated by the `xenial`
    suffix on the version, in this example). Choose a specific version to
    install. The second column is the version string. The third column is the
    repository name, which indicates which repository the package is from and
    by extension its stability level. To install a specific version, append the
    version string to the package name and separate them by an equals sign (`=`):

    ```bash
    $ sudo apt-get install docker-ee=<VERSION>
    ```

    The Docker daemon starts automatically.

4.  Verify that Docker EE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker EE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Linux postinstall](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker EE

To upgrade Docker EE:

1.  If upgrading to a new major Docker EE version (such as when going from
    Docker 17.03.x to Docker 17.06.x),
    [add the new repository](#set-up-the-repository){: target="_blank" class="_" }.

2.  Run `sudo apt-get update`.

3.  Follow the
    [installation instructions](#install-docker), choosing the new version you want
    to install.

### Install from a package

If you cannot use Docker's repository to install Docker EE, you can download the
`.deb` file for your release and install it manually. You will need to download
a new file each time you want to upgrade Docker EE.

1.  Go to the Docker EE repository URL associated with your
    trial or subscription in your browser. Go to
    `ubuntu/x86_64/stable-{{ minor-version }}` and download the `.deb` file for the
    Docker EE version you want to install.

2.  Install Docker EE, changing the path below to the path where you downloaded
    the Docker EE package.

    ```bash
    $ sudo dpkg -i /path/to/package.deb
    ```

    The Docker daemon starts automatically.

3.  Verify that Docker EE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker EE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Post-installation steps for Linux](/engine/installation/linux/linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps.

#### Upgrade Docker EE

To upgrade Docker EE, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

## Uninstall Docker EE

1.  Uninstall the Docker EE package:

    ```bash
    $ sudo apt-get purge docker-ee
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```bash
    $ sudo rm -rf /var/lib/docker
    ```

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](/engine/installation/linux/linux-postinstall.md)

- Continue with the [User Guide](/engine/userguide/index.md).
