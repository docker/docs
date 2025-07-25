---
description: Learn how to install Docker Engine on CentOS. These instructions cover
  the different installation methods, how to uninstall, and next steps.
keywords: requirements, dnf, yum, installation, centos, install, uninstall, docker engine, upgrade, update
title: Install Docker Engine on CentOS
linkTitle: CentOS
weight: 60
toc_max: 4
aliases:
- /ee/docker-ee/centos/
- /engine/installation/centos/
- /engine/installation/linux/centos/
- /engine/installation/linux/docker-ce/centos/
- /engine/installation/linux/docker-ee/centos/
- /install/linux/centos/
- /install/linux/docker-ce/centos/
- /install/linux/docker-ee/centos/
download-url-base: https://download.docker.com/linux/centos
---

To get started with Docker Engine on CentOS, make sure you
[meet the prerequisites](#prerequisites), and follow the
[installation steps](#installation-methods).

## Prerequisites

### OS requirements

To install Docker Engine, you need a maintained version of one of the following
CentOS versions:

- CentOS 9 (stream)

> [!NOTE]
> The `centos-extras` repository must be enabled. This repository is enabled by default.
> If you have disabled it, re-enable it before proceeding.

### Uninstall old versions

Uninstall unofficial or conflicting Docker packages provided
by your Linux distribution before installing Docker Engine.
Remove the following packages if present:

```console
$ sudo dnf remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
```

`dnf` might report that you have none of these packages installed.

> [!NOTE]
> Images, containers, volumes, and networks stored in `/var/lib/docker/` are
> not automatically removed when you uninstall Docker.

## Installation methods

You can install Docker Engine in different ways:

- Recommended: [Set up Docker's repositories](#install-using-the-repository) and install from them for easy installation and upgrades.
- Manual: [Download the RPM package](#install-from-a-package) and install it manually. Useful for air-gapped systems.
- Development/testing: Use [convenience scripts](#install-using-the-convenience-script) for quick setup.

### Install using the rpm repository {#install-using-the-repository}

Before installing Docker Engine for the first time on a new host, set up the Docker repository. Afterward, you can install and update Docker from the repository.

#### Set up the repository

Install the `dnf-plugins-core` package (provides commands to manage DNF repositories) and set up the repository.

```console
$ sudo dnf -y install dnf-plugins-core
$ sudo dnf config-manager --add-repo {{% param "download-url-base" %}}/docker-ce.repo
```

#### Install Docker Engine

1. Install the Docker packages.

   {{< tabs >}}
   {{< tab name="Latest" >}}

   To install the latest version, run:

   ```console
   $ sudo dnf install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
   ```

   If prompted to accept the GPG key, verify that the fingerprint matches
   `060A 61C5 1B55 8A7F 742B 77AA C52F EB6B 621E 9F35`, and if so, accept it.

   This command installs Docker, but it doesn't start Docker. It also creates a
   `docker` group, however, it doesn't add any users to the group by default.

   {{< /tab >}}
   {{< tab name="Specific version" >}}

   To install a specific version, start by listing the available versions in
   the repository:

   ```console
   $ dnf list docker-ce --showduplicates | sort -r

   docker-ce.x86_64    3:{{% param "docker_ce_version" %}}-1.el9    docker-ce-stable
   docker-ce.x86_64    3:{{% param "docker_ce_version_prev" %}}-1.el9    docker-ce-stable
   <...>
   ```

   The list returned depends on which repositories are enabled, and is specific
   to your version of CentOS (indicated by the `.el9` suffix in this example).

   Install a specific version by its fully qualified package name, which is
   the package name (`docker-ce`) plus the version string (2nd column),
   separated by a hyphen (`-`). For example, `docker-ce-3:{{% param "docker_ce_version" %}}-1.el9`.

   Replace `<VERSION_STRING>` with the desired version and then run the following
   command to install:

   ```console
   $ sudo dnf install docker-ce-<VERSION_STRING> docker-ce-cli-<VERSION_STRING> containerd.io docker-buildx-plugin docker-compose-plugin
   ```

   This command installs Docker, but it doesn't start Docker. It also creates a
   `docker` group, however, it doesn't add any users to the group by default.

   {{< /tab >}}
   {{< /tabs >}}

1. Start Docker Engine:

   ```console
   $ sudo systemctl enable --now docker
   ```

   This configures the Docker systemd service to start automatically when you
   boot your system. If you don't want Docker to start automatically, use `sudo
   systemctl start docker` instead.

1. Verify the installation by running the `hello-world` image:

   ```console
   $ sudo docker run hello-world
   ```

   This downloads a test image and runs it in a container. When the container runs, it prints a confirmation message and exits.

You have now installed and started Docker Engine.

{{% include "root-errors.md" %}}

#### Upgrade Docker Engine

To upgrade Docker Engine, follow the [installation instructions](#install-using-the-repository),
choosing the new version you want to install.

### Install from a package

If you cannot use Docker's `rpm` repository, download the `.rpm` file for your release and install it manually. Download a new file each time you want to upgrade Docker Engine.

<!-- markdownlint-disable-next-line -->
1. Go to [{{% param "download-url-base" %}}/]({{% param "download-url-base" %}}/)
   and choose your version of CentOS. Then browse to `x86_64/stable/Packages/`
   and download the `.rpm` file for the Docker version you want to install.

2. Install Docker Engine, changing the following path to the path where you downloaded
   the Docker package.

   ```console
   $ sudo dnf install /path/to/package.rpm
   ```

   Docker is installed but not started. The `docker` group
   is created, but no users are added to the group.

3. Start Docker Engine:

   ```console
   $ sudo systemctl enable --now docker
   ```

   This configures the Docker systemd service to start automatically when you
   boot your system. If you don't want Docker to start automatically, use `sudo
   systemctl start docker` instead.

4. Verify that the installation is successful by running the `hello-world` image:

   ```console
   $ sudo docker run hello-world
   ```

   This command downloads a test image and runs it in a container. When the
   container runs, it prints a confirmation message and exits.

You have now successfully installed and started Docker Engine.

{{% include "root-errors.md" %}}

#### Upgrade Docker Engine

To upgrade Docker Engine, download the newer package files and repeat the
[installation procedure](#install-from-a-package), using `dnf upgrade`
instead of `dnf install`, and point to the new files.

{{% include "install-script.md" %}}

## Uninstall Docker Engine

1. Uninstall the Docker Engine, CLI, containerd, and Docker Compose packages:

   ```console
   $ sudo dnf remove docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras
   ```

1. Images, containers, volumes, or custom configuration files on your host
   aren't automatically removed. To delete all images, containers, and volumes:

   ```console
   $ sudo rm -rf /var/lib/docker
   $ sudo rm -rf /var/lib/containerd
   ```

You have to delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md).
