---
description: Instructions for installing Docker CE on Debian
keywords: requirements, apt, installation, debian, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/debian/
- /engine/installation/linux/raspbian/
- /engine/installation/linux/debian/
- /engine/installation/linux/docker-ce/debian/
title: Get Docker CE for Debian
toc_max: 4
---

To get started with Docker CE on Debian, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ce).

## Prerequisites

### Docker EE customers

Docker EE is not supported on Debian. For a list of supported operating systems
and distributions for different Docker editions, see
[Docker variants](/install/index.md#docker-variants).

### OS requirements

To install Docker CE, you need the 64-bit version of one of these Debian or
Raspbian versions:

- Buster 10 (Docker CE 17.11 Edge only)
- Stretch 9 (stable) / Raspbian Stretch
- Jessie 8 (LTS) / Raspbian Jessie
- Wheezy 7.7 (LTS)

Docker CE is supported on both `x86_64` (or `amd64`)  and `armhf` architectures for Jessie and
Stretch.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them:

```bash
$ sudo apt-get remove docker docker-engine docker.io
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker CE package is now called `docker-ce`.

### Extra steps for Wheezy 7.7

- You need at least version 3.10 of the Linux kernel. Debian Wheezy ships with
  version 3.2, so you may need to
  [update the kernel](https://wiki.debian.org/HowToUpgradeKernel){: target="_blank" class="_" }.
  To check your kernel version:

  ```bash
  $ uname -r
  ```

- Enable the `backports` repository. See the
  [Debian documentation](https://backports.debian.org/Instructions/){: target="_blank" class"_"}.

## Install Docker CE

You can install Docker CE in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach, except for Raspbian.

- Some users download the DEB package and
  [install it manually](#install-from-a-package) and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

- In testing and development environments, some users choose to use automated
  [convenience scripts](#install-using-the-convenience-script) to install Docker.
  This is currently the only approach for Raspbian.

### Install using the repository

Before you install Docker CE for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker
from the repository.

> **Raspbian users cannot use this method!**
>
> For Raspbian, installing using the repository is not yet supported. You must
> instead use the [convenience script](#install-using-the-convenience-script).

#### Set up the repository

{% assign download-url-base = 'https://download.docker.com/linux/debian' %}

1.  Update the `apt` package index:

    ```bash
    $ sudo apt-get update
    ```

2.  Install packages to allow `apt` to use a repository over HTTPS:

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-target="#jessie">Jessie or newer</a></li>
      <li><a data-toggle="tab" data-target="#wheezy">Wheezy or older</a></li>
    </ul>
    <div class="tab-content">
    <div id="jessie" class="tab-pane fade in active" markdown="1">

    ```bash
    $ sudo apt-get install \
         apt-transport-https \
         ca-certificates \
         curl \
         gnupg2 \
         software-properties-common
    ```

    </div>
    <div id="wheezy" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo apt-get install \
         apt-transport-https \
         ca-certificates \
         curl \
         python-software-properties
    ```

    </div>
    </div> <!-- tab-content -->

3.  Add Docker's official GPG key:

    ```bash
    $ curl -fsSL {{ download-url-base}}/gpg | sudo apt-key add -
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
    > Debian distribution, such as `jessie`.

    To also add the **edge** repository, add `edge` after `stable` on the last
    line of the command.

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-target="#x86_64_repo">x86_64 / amd64</a></li>
      <li><a data-toggle="tab" data-target="#armhf_repo">armhf</a></li>
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
    <div id="armhf_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ echo "deb [arch=armhf] {{ download-url-base }} \
         $(lsb_release -cs) stable" | \
        sudo tee /etc/apt/sources.list.d/docker.list
    ```

    </div>
    </div> <!-- tab-content -->

5.  **Wheezy only**: The version of `add-apt-repository` on Wheezy adds a `deb-src`
    repository that does not exist. You need to comment out this repository or
    running `apt-get update` fails. Edit `/etc/apt/sources.list`. Find the
    line like the following, and comment it out or remove it:

    ```none
    deb-src [arch=amd64] https://download.docker.com/linux/debian wheezy stable
    ```

    Save and exit the file.

    > **Note**: Starting with Docker 17.06, stable releases are also pushed to
    > the **edge** and **test** repositories.

    [Learn about **stable** and **edge** channels](/install/index.md).

#### Install Docker CE

> **Note**: This procedure works for Debian on `x86_64` / `amd64`, Debian ARM,
> or Raspbian.

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

    docker-ce | {{ site.docker_ce_stable_version }}.0~ce-0~debian | https://download.docker.com/linux/debian jessie/stable amd64 Packages
    ```

    b. Install a specific version by its fully qualified package name, which is
       the package name (`docker-ce`) plus the version string (2nd column) up to
       the first hyphen, separated by a an equals sign (`=`), for example,
       `docker-ce=18.03.0.ce`.

    ```bash
    $ sudo apt-get install docker-ce=<VERSION_STRING>
    ```

    The Docker daemon starts automatically.

4.  Verify that Docker CE is installed correctly by running the `hello-world`
    image.

    **x86_64**:

    ```bash
    $ sudo docker run hello-world
    ```

    **armhf**:

    ```bash
    $ sudo docker run armhf/hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker CE is installed and running. The `docker` group is created but no users
are added to it. You need to use `sudo` to run Docker
commands. Continue to [Linux postinstall](/install/linux/linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps. For Raspbian, you can optionally
[install Docker Compose for Raspbian](#install-docker-compose-for-raspbian).

#### Upgrade Docker CE

To upgrade Docker CE, first run `sudo apt-get update`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker CE, you can download the
`.deb` file for your release and install it manually. You need to download
a new file each time you want to upgrade Docker.

1.  Go to `{{ download-url-base }}/dists/`,
    choose your Debian version, browse to `pool/stable/`, choose either
    `amd64` or `armhf`, and download the `.deb` file for the Docker CE version you
    want to install.

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
are added to it. You need to use `sudo` to run Docker
commands. Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps. For Raspbian, you can optionally
[install Docker Compose for Raspbian](#install-docker-compose-for-raspbian).

#### Upgrade Docker CE

To upgrade Docker, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

{% include install-script.md %}

## Install Docker Compose for Raspbian

You can install Docker Compose using `pip`:

```bash
$ sudo pip install docker-compose
```

[Hypriot](https://hypriot.com/){: target="_blank" class="_" } provides a static
binary of `docker-compose` for Raspbian. It may not always be up to date, but if
space is at a premium, you may find it useful. To use it, first follow Hypriot's
[instructions for setting up the repository](https://blog.hypriot.com/post/your-number-one-source-for-docker-on-arm/){: target="_blank" class="_" },
then run the following command:

```bash
sudo apt-get install docker-compose
```

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
