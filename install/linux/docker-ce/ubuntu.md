---
description: Instructions for installing Docker CE on Ubuntu
keywords: requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
- /engine/installation/linux/ubuntulinux/
- /engine/installation/linux/docker-ce/ubuntu/
title: Get Docker CE for Ubuntu
toc_max: 4
---

To get started with Docker CE on Ubuntu, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ce).

## Prerequisites

### Docker EE customers

To install Docker Enterprise Edition (Docker EE), go to
[Get Docker EE for Ubuntu](/install/linux/docker-ee/ubuntu.md)
**instead of this topic**.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### OS requirements

To install Docker CE, you need the 64-bit version of one of these Ubuntu
versions:

- Artful 17.10 (Docker CE 17.11 Edge and higher only)
- Xenial 16.04 (LTS)
- Trusty 14.04 (LTS)

Docker CE is supported on Ubuntu on `x86_64`, `armhf`, `s390x` (IBM Z), and `ppc64le` (IBM Power) architectures.

> **`ppc64le` and `s390x` limitations**: Packages for IBM Z and Power architectures are only available on Ubuntu Xenial and above.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them:

```bash
$ sudo apt-get remove docker docker-engine docker.io
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker CE package is now called `docker-ce`.

### Supported storage drivers

Docker CE on Ubuntu supports `overlay2` and `aufs` storage drivers.

- For new installations on version 4 and higher of the Linux kernel, `overlay2`
  is supported and preferred over `aufs`.
- For version 3 of the Linux kernel, `aufs` is supported because `overlay` or
  `overlay2` drivers are not supported by that kernel version.

If you need to use `aufs`, you need to do additional preparation as
outlined below.

#### Extra steps for aufs

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#aufs_prep_xenial">Xenial 16.04 and newer</a></li>
  <li><a data-toggle="tab" data-target="#aufs_prep_trusty">Trusty 14.04</a></li>
</ul>
<div class="tab-content">
<div id="aufs_prep_xenial" class="tab-pane fade in active" markdown="1">

For Ubuntu 16.04 and higher, the Linux kernel includes support for OverlayFS,
and Docker CE uses the `overlay2` storage driver by default. If you need
to use `aufs` instead, you need to configure it manually.
See [aufs](/engine/userguide/storagedriver/aufs-driver.md)

</div>
<div id="aufs_prep_trusty" class="tab-pane fade" markdown="1">

Unless you have a strong reason not to, install the
`linux-image-extra-*` packages, which allow Docker to use the `aufs` storage
drivers.

```bash
$ sudo apt-get update

$ sudo apt-get install \
    linux-image-extra-$(uname -r) \
    linux-image-extra-virtual
```

</div>
</div> <!-- tab-content -->

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
    > like Linux Mint, you might need to change `$(lsb_release -cs)`
    > to your parent Ubuntu distribution. For example, if you are using
    >  `Linux Mint Rafaela`, you could use `trusty`.


    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-target="#x86_64_repo">x86_64 / amd64</a></li>
      <li><a data-toggle="tab" data-target="#armhf">armhf</a></li>
      <li><a data-toggle="tab" data-target="#ppc64le_repo">IBM Power (ppc64le)</a></li>
      <li><a data-toggle="tab" data-target="#s390x_repo">IBM Z (s390x)</a></li>
    </ul>
    <div class="tab-content">
    <div id="x86_64_repo" class="tab-pane fade in active" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=amd64] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    <div id="armhf" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=armhf] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    <div id="ppc64le_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=ppc64el] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    <div id="s390x_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=s390x] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    </div>
    </div> <!-- tab-content -->

    > **Note**: Starting with Docker 17.06, stable releases are also pushed to
    > the **edge** and **test** repositories.

    [Learn about **stable** and **edge** channels](/install/index.md).


#### Install Docker CE

1.  Update the `apt` package index.

    ```bash
    $ sudo apt-get update
    ```

2.  Install the _latest version_ of Docker CE, or go to the next step to install a specific version:

    ```bash
    $ sudo apt-get install docker-ce
    ```

    > Got multiple Docker repositories?
    >
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `apt-get install` or
    > `apt-get update` command always installs the highest possible version,
    > which may not be appropriate for your stability needs.

3.  To install a _specific version_ of Docker CE, list the available versions in the repo, then select and install:

    a. List the versions available in your repo:

    ```bash
    $ apt-cache madison docker-ce

    docker-ce | {{ site.docker_ce_stable_version }}.0~ce-0~ubuntu | {{ download-url-base }} xenial/stable amd64 Packages
    ```

    b. Install a specific version by its fully qualified package name, which is
       the package name (`docker-ce`) plus the version string (2nd column) up to
       the first hyphen, separated by a an equals sign (`=`), for example,
       `docker-ce=18.03.0.ce`.

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

Docker CE is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](/install/linux/linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker CE

To upgrade Docker CE, first run `sudo apt-get update`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker CE, you can download the
`.deb` file for your release and install it manually. You need to download
a new file each time you want to upgrade Docker CE.

1.  Go to [{{ download-url-base }}/dists/]({{ download-url-base }}/dists/),
    choose your Ubuntu version, browse to `pool/stable/` and choose `amd64`,
    `armhf`, `ppc64el`, or `s390x`. Download the `.deb` file for the Docker
    version you want to install.

    > **Note**: To install an **edge**  package, change the word
    > `stable` in the  URL to `edge`.
    > [Learn about **stable** and **edge** channels](/install/index.md).

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

Docker CE is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md) to allow
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

- Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md)

- Continue with the [User Guide](/engine/userguide/index.md).
