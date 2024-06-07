---
description: Learn how to install Docker Engine on CentOS. These instructions cover
  the different installation methods, how to uninstall, and next steps.
keywords: requirements, yum, installation, centos, install, uninstall, docker engine, upgrade, update
title: Install Docker Engine on CentOS
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
[meet the prerequisites](#prerequisites), and then follow the
[installation steps](#installation-methods).

## Prerequisites

### OS requirements

To install Docker Engine, you need a maintained version of one of the following
CentOS versions:

- CentOS 7 (EOL: [June 30, 2024](https://blog.centos.org/2023/04/end-dates-are-coming-for-centos-stream-8-and-centos-linux-7/#centos-linux-7-end-of-life-june-30-2024))
- CentOS 9 (stream)

The `centos-extras` repository must be enabled. This repository is enabled by
default. If you have disabled it, you need to re-enable it.

### Uninstall old versions

Older versions of Docker went by `docker` or `docker-engine`.
Uninstall any such older versions before attempting to install a new version,
along with associated dependencies.

```console
$ sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
```

`yum` might report that you have none of these packages installed.

Images, containers, volumes, and networks stored in `/var/lib/docker/` aren't
automatically removed when you uninstall Docker.

## Installation methods

You can install Docker Engine in different ways, depending on your needs:

- You can
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- You can download the RPM package,
  [install it manually](#install-from-a-package), and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

- In testing and development environments, you can use automated
  [convenience scripts](#install-using-the-convenience-script) to install Docker.

### Install using the rpm repository {#install-using-the-repository}

Before you install Docker Engine for the first time on a new host machine, you
need to set up the Docker repository. Afterward, you can install and update
Docker from the repository.

#### Set up the repository

Install the `yum-utils` package (which provides the `yum-config-manager`
utility) and set up the repository.

```console
$ sudo yum install -y yum-utils
$ sudo yum-config-manager --add-repo {{% param "download-url-base" %}}/docker-ce.repo
```

#### Install Docker Engine

1. Install Docker Engine, containerd, and Docker Compose:

   {{< tabs >}}
   {{< tab name="Latest" >}}

   To install the latest version, run:

   ```console
   $ sudo yum install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
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
   $ yum list docker-ce --showduplicates | sort -r

   docker-ce.x86_64    3:26.1.3-1.el9    docker-ce-stable
   docker-ce.x86_64    3:26.1.2-1.el9    docker-ce-stable
   <...>
   ```

   The list returned depends on which repositories are enabled, and is specific
   to your version of CentOS (indicated by the `.el9` suffix in this example).

   Install a specific version by its fully qualified package name, which is
   the package name (`docker-ce`) plus the version string (2nd column),
   separated by a hyphen (`-`). For example, `docker-ce-3:26.1.2-1.el9`.

   Replace `<VERSION_STRING>` with the desired version and then run the following
   command to install:

   ```console
   $ sudo yum install docker-ce-<VERSION_STRING> docker-ce-cli-<VERSION_STRING> containerd.io docker-buildx-plugin docker-compose-plugin
   ```

   This command installs Docker, but it doesn't start Docker. It also creates a
   `docker` group, however, it doesn't add any users to the group by default.
   
   {{< /tab >}}
   {{< /tabs >}}


2. Start Docker.

   ```console
   $ sudo systemctl start docker
   ```

3. Verify that the Docker Engine installation is successful by running the
   `hello-world` image.

   ```console
   $ sudo docker run hello-world
   ```

   This command downloads a test image and runs it in a container. When the
   container runs, it prints a confirmation message and exits.

You have now successfully installed and started Docker Engine.

{{< include "root-errors.md" >}}

#### Upgrade Docker Engine

To upgrade Docker Engine, follow the [installation instructions](#install-using-the-repository),
choosing the new version you want to install.

### Install from a package

If you can't use Docker's `rpm` repository to install Docker Engine, you can
download the `.rpm` file for your release and install it manually. You need to
download a new file each time you want to upgrade Docker Engine.

<!-- markdownlint-disable-next-line -->
1. Go to [{{% param "download-url-base" %}}/]({{% param "download-url-base" %}}/)
   and choose your version of CentOS. Then browse to `x86_64/stable/Packages/`
   and download the `.rpm` file for the Docker version you want to install.

2. Install Docker Engine, changing the following path to the path where you downloaded
   the Docker package.

   ```console
   $ sudo yum install /path/to/package.rpm
   ```

   Docker is installed but not started. The `docker` group is created, but no
   users are added to the group.

3. Start Docker.

   ```console
   $ sudo systemctl start docker
   ```

4. Verify that the Docker Engine installation is successful by running the
   `hello-world` image.

   ```console
   $ sudo docker run hello-world
   ```

   This command downloads a test image and runs it in a container. When the
   container runs, it prints a confirmation message and exits.

You have now successfully installed and started Docker Engine.

{{< include "root-errors.md" >}}

#### Upgrade Docker Engine

To upgrade Docker Engine, download the newer package files and repeat the
[installation procedure](#install-from-a-package), using `yum -y upgrade`
instead of `yum -y install`, and point to the new files.

{{< include "install-script.md" >}}

## Uninstall Docker Engine

1. Uninstall the Docker Engine, CLI, containerd, and Docker Compose packages:

   ```console
   $ sudo yum remove docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras
   ```

2. Images, containers, volumes, or custom configuration files on your host
   aren't automatically removed. To delete all images, containers, and volumes:

   ```console
   $ sudo rm -rf /var/lib/docker
   $ sudo rm -rf /var/lib/containerd
   ```

You have to delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md).
