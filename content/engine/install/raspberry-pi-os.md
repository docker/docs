---
description: Instructions for installing Docker Engine on a 32-bit Raspberry Pi OS system
keywords: requirements, apt, installation, Raspberry Pi OS, install, uninstall, upgrade,
  update
title: Install Docker Engine on Raspberry Pi OS (32-bit)
toc_max: 4
aliases:
- /engine/installation/linux/raspbian/
- /engine/install/raspbian/
download-url-base: https://download.docker.com/linux/raspbian
---

To get started with Docker Engine on Raspberry Pi OS, make sure you
[meet the prerequisites](#prerequisites), and then follow the
[installation steps](#installation-methods).

> **Important**
>
> This installation instruction refers to the 32-bit (armhf) version of
> Raspberry Pi OS. If you're using the 64-bit (arm64) version, follow the
> instructions for [Debian](debian.md).
{ .important }

## Prerequisites

> **Note**
>
> If you use ufw or firewalld to manage firewall settings, be aware that
> when you expose container ports using Docker, these ports bypass your
> firewall rules. For more information, refer to
> [Docker and ufw](../../network/packet-filtering-firewalls.md#docker-and-ufw).

### OS requirements

The following OS versions are supported:

- 32-bit Raspberry Pi OS Bookworm 12 (stable)
- 32-bit Raspberry Pi OS Bullseye 11 (oldstable)

### Uninstall old versions

Before you can install Docker Engine, you must first make sure that any
conflicting packages are uninstalled.

Distro maintainers provide an unofficial distributions of Docker packages in
APT. You must uninstall these packages before you can install the official
version of Docker Engine.

The unofficial packages to uninstall are:

- `docker.io`
- `docker-compose`
- `docker-doc`
- `podman-docker`

Moreover, Docker Engine depends on `containerd` and `runc`. Docker Engine
bundles these dependencies as one bundle: `containerd.io`. If you have
installed the `containerd` or `runc` previously, uninstall them to avoid
conflicts with the versions bundled with Docker Engine.

Run the following command to uninstall all conflicting packages:

```console
$ for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt-get remove $pkg; done
```

`apt-get` might report that you have none of these packages installed.

Images, containers, volumes, and networks stored in `/var/lib/docker/` aren't
automatically removed when you uninstall Docker. If you want to start with a
clean installation, and prefer to clean up any existing data, read the
[uninstall Docker Engine](#uninstall-docker-engine) section.

## Installation methods

You can install Docker Engine in different ways, depending on your needs:

- Docker Engine comes bundled with
  [Docker Desktop for Linux](../../desktop/install/linux-install.md). This is
  the easiest and quickest way to get started.

- Set up and install Docker Engine from
  [Docker's `apt` repository](#install-using-the-repository).

- [Install it manually](#install-from-a-package) and manage upgrades manually.

- Use a [convenience scripts](#install-using-the-convenience-script). Only
  recommended for testing and development environments.

### Install using the apt repository {#install-using-the-repository}

Before you install Docker Engine for the first time on a new host machine, you
need to set up the Docker Apt repository. Afterward, you can install and update
Docker from the repository.

1. Set up Docker's Apt repository.

   ```bash
   # Add Docker's official GPG key:
   sudo apt-get update
   sudo apt-get install ca-certificates curl gnupg
   sudo install -m 0755 -d /etc/apt/keyrings
   curl -fsSL {{% param "download-url-base" %}}/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
   sudo chmod a+r /etc/apt/keyrings/docker.gpg

   # Set up Docker's Apt repository:
   echo \
     "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] {{% param "download-url-base" %}} \
     "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
     sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
   sudo apt-get update
   ```

2. Install the Docker packages.

   {{< tabs >}}
   {{< tab name="Latest" >}}

   To install the latest version, run:

   ```console
   $ sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
   ```
  
   {{< /tab >}}
   {{< tab name="Specific version" >}}
  
   To install a specific version of Docker Engine, start by listing the
   available versions in the repository:

   ```console
   # List the available versions:
   $ apt-cache madison docker-ce | awk '{ print $3 }'

   5:24.0.0-1~raspbian.11~bullseye
   5:23.0.6-1~raspbian.11~bullseye
   ...
   ```

   Select the desired version and install:

   ```console
   $ VERSION_STRING=5:24.0.0-1~raspbian.11~bullseye
   $ sudo apt-get install docker-ce=$VERSION_STRING docker-ce-cli=$VERSION_STRING containerd.io docker-buildx-plugin docker-compose-plugin
   ```

   {{< /tab >}}
   {{< /tabs >}}

3. Verify that the installation is successful by running the `hello-world`
   image:

   ```console
   $ sudo docker run hello-world
   ```

   This command downloads a test image and runs it in a container. When the
   container runs, it prints a confirmation message and exits.

You have now successfully installed and started Docker Engine.

{{< include "root-errors.md" >}}

#### Upgrade Docker Engine

To upgrade Docker Engine, follow step 2 of the
[installation instructions](#install-using-the-repository),
choosing the new version you want to install.

### Install from a package

If you can't use Docker's `apt` repository to install Docker Engine, you can
download the `deb` file for your release and install it manually. You need to
download a new file each time you want to upgrade Docker Engine.

1. Go to [`{{% param "download-url-base" %}}/dists/`]({{% param "download-url-base" %}}/dists/).

2. Select your Raspberry Pi OS version in the list.

3. Go to `pool/stable/` and select the applicable architecture (`amd64`,
   `armhf`, `arm64`, or `s390x`).

4. Download the following `deb` files for the Docker Engine, CLI, containerd,
   and Docker Compose packages:

   - `containerd.io_<version>_<arch>.deb`
   - `docker-ce_<version>_<arch>.deb`
   - `docker-ce-cli_<version>_<arch>.deb`
   - `docker-buildx-plugin_<version>_<arch>.deb`
   - `docker-compose-plugin_<version>_<arch>.deb`

5. Install the `.deb` packages. Update the paths in the following example to
   where you downloaded the Docker packages.

   ```console
   $ sudo dpkg -i ./containerd.io_<version>_<arch>.deb \
     ./docker-ce_<version>_<arch>.deb \
     ./docker-ce-cli_<version>_<arch>.deb \
     ./docker-buildx-plugin_<version>_<arch>.deb \
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

You have now successfully installed and started Docker Engine.

{{< include "root-errors.md" >}}

#### Upgrade Docker Engine

To upgrade Docker Engine, download the newer package files and repeat the
[installation procedure](#install-from-a-package), pointing to the new files.

{{< include "install-script.md" >}}

## Uninstall Docker Engine

1.  Uninstall the Docker Engine, CLI, containerd, and Docker Compose packages:

    ```console
    $ sudo apt-get purge docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras
    ```

2.  Images, containers, volumes, or custom configuration files on your host
    aren't automatically removed. To delete all images, containers, and volumes:

    ```console
    $ sudo rm -rf /var/lib/docker
    $ sudo rm -rf /var/lib/containerd
    ```

You have to delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md).
- Review the topics in [Develop with Docker](../../develop/index.md) to learn
  how to build new applications using Docker.
