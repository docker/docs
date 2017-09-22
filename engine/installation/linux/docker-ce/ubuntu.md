---
description: Instructions for installing Docker CE on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
- /engine/installation/linux/ubuntulinux/
title: Get Docker CE for Ubuntu
---

{% assign minor-version = "17.06" %}

To get started with Docker CE on Ubuntu, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### Docker EE customers

To install Docker Enterprise Edition (Docker EE), go to
[Get Docker EE for Ubuntu](/engine/installation/linux/docker-ee/ubuntu/)
**instead of this topic**.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### OS requirements

To install Docker CE, you need the 64-bit version of one of these Ubuntu
versions:

- Zesty 17.04
- Xenial 16.04 (LTS)
- Trusty 14.04 (LTS)

Docker CE is supported on Ubuntu on `x86_64`, `armhf`, and `s390x` (IBM z
Systems) architectures.

> **`s390x` limitations**: System Z is only supported on Ubuntu Xenial and Zesty.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them:

```bash
$ sudo apt-get remove docker docker-engine docker.io
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker CE package is now called `docker-ce`.

### Recommended extra packages for Trusty 14.04

Unless you have a strong reason not to, install the
`linux-image-extra-*` packages, which allow Docker to use the `aufs` storage
drivers.

```bash
$ sudo apt-get update

$ sudo apt-get install \
    linux-image-extra-$(uname -r) \
    linux-image-extra-virtual
```

For Ubuntu 16.04 and higher, the Linux kernel includes support for OverlayFS,
and Docker CE will use the `overlay2` storage driver by default.

## Install Docker CE

You can install Docker CE in different ways, depending on your needs:

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

Before you install Docker CE for the first time on a new host machine, you need to
set up the Docker repository. Afterward, you can install and update Docker from
the repository.

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/ubuntu" %}

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

3.  Add Docker's official GPG key:

    ```bash
    $ curl -fsSL {{ download-url-base }}/gpg | sudo apt-key add -
    ```

    Verify that you now have the key with the fingerprint
    `9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88`, by searching for the
    last 8 characters of the fingerprint.

    ```bash
    $ sudo apt-key fingerprint 0EBFCD88

    pub   4096R/0EBFCD88 2017-02-22
          Key fingerprint = 9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
    uid                  Docker Release (CE deb) <docker@docker.com>
    sub   4096R/F273FCD8 2017-02-22
    ```

4.  Use the following command to set up the **stable** repository. You always
    need the **stable** repository, even if you want to install builds from the
    **edge** or **test** repositories as well. To add the **edge** or
    **test** repository, add the word `edge` or `test` (or both) after the
    word `stable` in the commands below.

    > **Note**: The `lsb_release -cs` sub-command below returns the name of your
    > Ubuntu distribution, such as `xenial`. Sometimes, in a distribution
    > like Linux Mint, you might have to change `$(lsb_release -cs)`
    > to your parent Ubuntu distribution. For example, if you are using
    >  `Linux Mint Rafaela`, you could use `trusty`.

    **amd64**:

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=amd64] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    **armhf**:

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=armhf] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    **s390x**:

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=s390x] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    > **Note**: Starting with Docker 17.06, stable releases are also pushed to
    > the **edge** and **test** repositories.

    [Learn about **stable** and **edge** channels](/engine/installation/).


#### Install Docker CE

1.  Update the `apt` package index.

    ```bash
    $ sudo apt-get update
    ```

2.  Install the latest version of Docker CE, or go to the next step to install a
    specific version. Any existing installation of Docker is replaced.

    ```bash
    $ sudo apt-get install docker-ce
    ```

    > Got multiple Docker repositories?
    >
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `apt-get install` or
    > `apt-get update` command will always install the highest possible version,
    > which may not be appropriate for your stability needs.
    {:.warning-vanilla}

3.  On production systems, you should install a specific version of Docker CE
    instead of always using the latest. This output is truncated. List the
    available versions.

    ```bash
    $ apt-cache madison docker-ce

    docker-ce | {{ minor-version }}.0~ce-0~ubuntu | {{ download-url-base }} xenial/stable amd64 Packages
    ```

    The contents of the list depend upon which repositories are enabled. Choose
    a specific version to install. The second column is the version string. The
    third column is the repository name, which indicates which repository the
    package is from and by extension its stability level. To install a specific
    version, append the version string to the package name and separate them by
    an equals sign (`=`):

    ```bash
    $ sudo apt-get install docker-ce=<VERSION>
    ```

    The Docker daemon starts automatically.

4.  Verify that Docker CE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker CE is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](../linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker CE

To upgrade Docker CE, first run `sudo apt-get update`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker CE, you can download the
`.deb` file for your release and install it manually. You will need to download
a new file each time you want to upgrade Docker CE.

1.  Go to [{{ download-url-base }}/dists/]({{ download-url-base }}/dists/),
    choose your Ubuntu version, browse to `pool/stable/` and choose `amd64`,
    `armhf`, or `s390x`. Download the `.deb` file for the Docker version you
    want to install.

    > **Note**: To install an **edge**  package, change the word
    > `stable` in the  URL to `edge`.
    > [Learn about **stable** and **edge** channels](/engine/installation/).

2.  Install Docker CE, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo dpkg -i /path/to/package.deb
    ```

    The Docker daemon starts automatically.

3.  Verify that Docker CE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker CE is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](/engine/installation/linux/linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker CE

To upgrade Docker CE, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

{% include install-script.md %}

## Uninstall Docker CE

1.  Uninstall the Docker CE package:

    ```bash
    $ sudo apt-get purge docker-ce
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
