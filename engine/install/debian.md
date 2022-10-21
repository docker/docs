---
description: Instructions for installing Docker Engine on Debian
keywords:
  requirements, apt, installation, debian, install, uninstall, upgrade, update
redirect_from:
  - /engine/installation/debian/
  - /engine/installation/linux/raspbian/
  - /engine/installation/linux/debian/
  - /engine/installation/linux/docker-ce/debian/
  - /install/linux/docker-ce/debian/
title: Install Docker Engine on Debian
toc_max: 4
---

To get started with Docker Engine on Debian, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#installation-methods).

## Prerequisites

### OS requirements

To install Docker Engine, you need the 64-bit version of one of these Debian or
Raspbian versions:

- Debian Bullseye 11 (stable)
- Debian Buster 10 (oldstable)
- Raspbian Bullseye 11 (stable)
- Raspbian Buster 10 (oldstable)

Docker Engine is compatible with `x86_64` (or `amd64`), `armhf`, and `arm64`
architectures.

### Uninstall old versions

Older versions of Docker went by the names of `docker`, `docker.io`, or
`docker-engine`. Uninstall any such older versions before attempting to install
a new version:

```console
$ sudo apt-get remove docker docker-engine docker.io containerd runc
```

It's OK if `apt-get` reports that none of these packages are installed.

Images, containers, volumes, and networks stored in `/var/lib/docker/` aren't
automatically removed when you uninstall Docker. If you want to start with a
clean installation, and prefer to clean up any existing data, refer to the
[uninstall Docker Engine](#uninstall-docker-engine) section.

## Installation methods

You can install Docker Engine in different ways, depending on your needs:

- Docker Engine comes bundled with
  [Docker Desktop for Linux](../../desktop/install/linux-install.md). This is
  the easiest and quickest way to get started.

- You can also set up and install Docker Engine from
  [Docker's `apt` repository](#install-using-the-repository).

- [Install it manually](#install-from-a-package) and manage upgrades manually.

- Using a [convenience scripts](#install-using-the-convenience-script). Only
  recommended for testing and development environments. This is the only
  approach available for Raspbian.

### Install using the repository

Before you install Docker Engine for the first time on a new host machine, you
need to set up the Docker repository. Afterward, you can install and update
Docker from the repository.

> **Raspbian users can't use this method.**
>
> For Raspbian, installing using the repository isn't yet supported. You must
> instead use the [convenience script](#install-using-the-convenience-script).

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/debian" %}

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
    $ sudo mkdir -p /etc/apt/keyrings
    $ curl -fsSL {{ download-url-base }}/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    ```

3.  Use the following command to set up the repository:

    ```console
    $ echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] {{ download-url-base }} \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    ```

#### Install Docker Engine

This procedure works for Debian on `x86_64` / `amd64`, `armhf`, `arm64`, and
Raspbian.

1. Update the `apt` package index:

   ```console
   $ sudo apt-get update
   ```

   > Receiving a GPG error when running `apt-get update`?
   >
   > Your default [umask](https://en.wikipedia.org/wiki/Umask){: target="blank"
   > rel="noopener" } may be incorrectly configured, preventing detection of the
   > repository public key file. Try granting read permission for the Docker
   > public key file before updating the package index:
   >
   > ```console
   > $ sudo chmod a+r /etc/apt/keyrings/docker.gpg
   > $ sudo apt-get update
   > ```

2. Install Docker Engine, containerd, and Docker Compose.

   <ul class="nav nav-tabs">
    <li class="active"><a data-toggle="tab" data-target="#tab-latest">Latest</a></li>
    <li><a data-toggle="tab" data-target="#tab-version">Specific version</a></li>
   </ul>
   <div class="tab-content">
   <br>
   <div id="tab-latest" class="tab-pane fade in active" markdown="1">

   To install the latest version, run:

   ```console
    $ sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
   ```

   </div>
   <div id="tab-version" class="tab-pane fade" markdown="1">

   To install a specific version of Docker Engine, start by list the available
   versions in the repository:

   ```console
   # List the available versions:
   $ apt-cache madison docker-ce | awk '{ print $3 }'

   5:18.09.1~3-0~debian-stretch
   5:18.09.0~3-0~debian-stretch
   18.06.1~ce~3-0~debian
   18.06.0~ce~3-0~debian
   ```

   Select the desired version and install:

   ```console
   $ VERSION_STRING=5:18.09.0~3-0~debian-stretch
   $ sudo apt-get install docker-ce=$VERSION_STRING docker-ce-cli=$VERSION_STRING containerd.io docker-compose-plugin
   ```

   </div>
   <hr>
   </div>

3. Verify that the Docker Engine installation is successful by running the
   `hello-world` image:

   ```console
   $ sudo docker run hello-world
   ```

   This command downloads a test image and runs it in a container. When the
   container runs, it prints a confirmation message and exits.

You have now successfully installed and started Docker Engine. The `docker` user
group exists but contains no users, which is why you're required to use `sudo`
to run Docker commands. Continue to [Linux post-install](linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps.

#### Upgrade Docker Engine

To upgrade Docker Engine, follow the
[installation instructions](#install-docker-engine), choosing the new version
you want to install.

### Install from a package

If you can't use Docker's `apt` repository to install Docker Engine, you can
download the `deb` file for your release and install it manually. You need to
download a new file each time you want to upgrade Docker Engine.

1. Go to [`{{ download-url-base }}/dists/`]({{ download-url-base }}/dists/){:
   target="_blank" rel="noopener" class="_" }.

2. Select your Debian version in the list.

3. Go to `pool/stable/` and select the applicable architecture (`amd64`,
   `armhf`, `arm64`, or `s390x`).

4. Download the following `deb` files for the Docker Engine, CLI, containerd,
   and Docker Compose packages:

   - `containerd.io_<version>_<arch>.deb`
   - `docker-ce_<version>_<arch>.deb`
   - `docker-ce-cli_<version>_<arch>.deb`
   - `docker-compose-plugin_<version>_<arch>.deb`

5. Install the `.deb` packages. Update the paths in the following example to
   where you downloaded the Docker packages.

   ```console
   $ sudo dpkg -i ./containerd.io_<version>_<arch>.deb \
     ./docker-ce_<version>_<arch>.deb \
     ./docker-ce-cli_<version>_<arch>.deb \
     ./docker-compose-plugin_<version>_<arch>.deb
   ```

   The Docker daemon starts automatically.

6. Verify that the Docker Engine installation is successful by running the
   `hello-world` image:

   ```console
   $ sudo service docker start
   $ sudo docker run hello-world
   ```

   This command downloads a test image and runs it in a container. When the
   container runs, it prints a confirmation message and exits.

You have now successfully installed and started Docker Engine. The `docker` user
group exists but contains no users, which is why you're required to use `sudo`
to run Docker commands. Continue to [Linux post-install](linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps.

#### Upgrade Docker Engine

To upgrade Docker Engine, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

{% include install-script.md %}

## Uninstall Docker Engine

1.  Uninstall the Docker Engine, CLI, containerd, and Docker Compose packages:

    ```console
    $ sudo apt-get purge docker-ce docker-ce-cli containerd.io docker-compose-plugin
    ```

2.  Images, containers, volumes, or custom configuration files on your host
    aren't automatically removed. To delete all images, containers, and volumes:

    ```console
    $ sudo rm -rf /var/lib/docker
    $ sudo rm -rf /var/lib/containerd
    ```

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md).
- Review the topics in [Develop with Docker](../../develop/index.md) to learn
  how to build new applications using Docker.
