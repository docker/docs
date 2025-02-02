---
description: Learn how to install Docker Engine on SLES. These instructions cover
  the different installation methods, how to uninstall, and next steps.
keywords: requirements, apt, installation, install docker engine, centos, rpm, sles, install, uninstall,
  upgrade, update, s390x, ibm-z
title: Install Docker Engine on SLES (s390x)
linkTitle: SLES (s390x)
weight: 70
toc_max: 4
aliases:
- /ee/docker-ee/sles/
- /ee/docker-ee/suse/
- /engine/installation/linux/docker-ce/sles/
- /engine/installation/linux/docker-ee/sles/
- /engine/installation/linux/docker-ee/suse/
- /engine/installation/linux/sles/
- /engine/installation/linux/SUSE/
- /engine/installation/linux/suse/
- /engine/installation/sles/
- /engine/installation/SUSE/
- /install/linux/docker-ce/sles/
- /install/linux/docker-ee/sles/
- /install/linux/docker-ee/suse/
- /install/linux/sles/
- /installation/sles/
download-url-base: https://download.docker.com/linux/sles
---

> [!NOTE]
>
> The installation instructions on this page refer to packages for SLES on the
> **s390x** architecture (IBM Z). Other architectures, including x86_64, aren't
> supported for SLES.

To get started with Docker Engine on SLES, make sure you
[meet the prerequisites](#prerequisites), and then follow the
[installation steps](#installation-methods).

## Prerequisites

### OS requirements

To install Docker Engine, you need a maintained version of one of the following
SLES versions:

- SLES 15-SP4 on s390x (IBM Z)
- SLES 15-SP5 on s390x (IBM Z)

You must enable the [`SCC SUSE`](https://scc.suse.com/packages?name=SUSE%20Linux%20Enterprise%20Server&version=15.5&arch=s390x)
repositories.

You must add the [OpenSUSE `SELinux` repository](https://download.opensuse.org/repositories/security:/SELinux/). This repository is not added by default. Run the following commands to add it:

```console
$ opensuse_repo="https://download.opensuse.org/repositories/security:/SELinux/openSUSE_Factory/security:SELinux.repo"
$ sudo zypper addrepo $opensuse_repo
```

### Uninstall old versions

Before you can install Docker Engine, you need to uninstall any conflicting packages.

Your Linux distribution may provide unofficial Docker packages, which may conflict
with the official packages provided by Docker. You must uninstall these packages
before you install the official version of Docker Engine.

```console
$ sudo zypper remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine \
                  runc
```

`zypper` might report that you have none of these packages installed.

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

Set up the repository.

```console
$ sudo zypper addrepo {{% param "download-url-base" %}}/docker-ce.repo
```

#### Install Docker Engine

1. Install the Docker packages.

   {{< tabs >}}
   {{< tab name="Latest" >}}

   To install the latest version, run:

   ```console
   $ sudo zypper install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
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
   $ sudo zypper search -s --match-exact docker-ce | sort -r
 
     v  | docker-ce | package | 3:{{% param "docker_ce_version" %}}-1 | s390x | Docker CE Stable - s390x
     v  | docker-ce | package | 3:{{% param "docker_ce_version_prev" %}}-1 | s390x | Docker CE Stable - s390x
   ```

   The list returned depends on which repositories are enabled, and is specific
   to your version of SLES.

   Install a specific version by its fully qualified package name, which is
   the package name (`docker-ce`) plus the version string (2nd column),
   separated by a hyphen (`-`). For example, `docker-ce-3:{{% param "docker_ce_version" %}}`.

   Replace `<VERSION_STRING>` with the desired version and then run the following
   command to install:

   ```console
   $ sudo zypper install docker-ce-<VERSION_STRING> docker-ce-cli-<VERSION_STRING> containerd.io docker-buildx-plugin docker-compose-plugin
   ```

   This command installs Docker, but it doesn't start Docker. It also creates a
   `docker` group, however, it doesn't add any users to the group by default.
  
   {{< /tab >}}
   {{< /tabs >}}

2. Start Docker Engine.

   ```console
   $ sudo systemctl enable --now docker
   ```

   This configures the Docker systemd service to start automatically when you
   boot your system. If you don't want Docker to start automatically, use `sudo
   systemctl start docker` instead.

3. Verify that the installation is successful by running the `hello-world` image:

   ```console
   $ sudo docker run hello-world
   ```

   This command downloads a test image and runs it in a container. When the
   container runs, it prints a confirmation message and exits.

You have now successfully installed and started Docker Engine.

{{% include "root-errors.md" %}}

#### Upgrade Docker Engine

To upgrade Docker Engine, follow the [installation instructions](#install-using-the-repository),
choosing the new version you want to install.

### Install from a package

If you can't use Docker's `rpm` repository to install Docker Engine, you can
download the `.rpm` file for your release and install it manually. You need to
download a new file each time you want to upgrade Docker Engine.

<!-- markdownlint-disable-next-line -->
1. Go to [{{% param "download-url-base" %}}/]({{% param "download-url-base" %}}/)
   and choose your version of SLES. Then browse to `s390x/stable/Packages/`
   and download the `.rpm` file for the Docker version you want to install.

2. Install Docker Engine, changing the following path to the path where you downloaded
   the Docker package.

   ```console
   $ sudo zypper install /path/to/package.rpm
   ```

   Docker is installed but not started. The `docker` group is created, but no
   users are added to the group.

3. Start Docker Engine.

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
[installation procedure](#install-from-a-package), using `zypper -y upgrade`
instead of `zypper -y install`, and point to the new files.

{{% include "install-script.md" %}}

## Uninstall Docker Engine

1. Uninstall the Docker Engine, CLI, containerd, and Docker Compose packages:

   ```console
   $ sudo zypper remove docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras
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
