---
description: Instructions for installing Docker Engine on RHEL
keywords: requirements, apt, installation, rhel, rpm, install, uninstall, upgrade, update, s390x, ibm-z
redirect_from:
- /ee/docker-ee/rhel/
- /engine/installation/linux/docker-ce/rhel/
- /engine/installation/linux/docker-ee/rhel/
- /engine/installation/linux/rhel/
- /engine/installation/rhel/
- /engine/installation/rhel/
- /install/linux/docker-ee/rhel/
- /installation/rhel/
title: Install Docker Engine on RHEL
toc_max: 4
---

To get started with Docker Engine on RHEL, make sure you
[meet the prerequisites](#prerequisites), and then follow the
[installation steps](#installation-methods).

## Prerequisites

> **Note**
>
> We currently only provide packages for RHEL on s390x (IBM Z). Other architectures
> are not yet supported for RHEL, but you may be able to install the CentOS packages
> on RHEL. Refer to the [Install Docker Engine on CentOS](centos.md) page for details.

### OS requirements

To install Docker Engine, you need a maintained version of one of the following
RHEL versions:

- RHEL 7 on s390x (IBM Z)
- RHEL 8 on s390x (IBM Z)
- RHEL 9 on s390x (IBM Z)

### Uninstall old versions

Older versions of Docker went by the names of `docker` or `docker-engine`.
Uninstall any such older versions before attempting to install a new version,
along with associated dependencies. Also uninstall `Podman` and the associated
dependencies if installed already:

```console
$ sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine \
                  podman \
                  runc
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

- You can download the RPM package and
  [install it manually](#install-from-a-package) and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

- In testing and development environments, you can use automated
  [convenience scripts](#install-using-the-convenience-script) to install Docker.

### Install using the rpm repository {#install-using-the-repository}

Before you install Docker Engine for the first time on a new host machine, you
need to set up the Docker repository. Afterward, you can install and update
Docker from the repository.

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/rhel" %}

Install the `yum-utils` package (which provides the `yum-config-manager`
utility) and set up the repository.

```console
$ sudo yum install -y yum-utils
$ sudo yum-config-manager --add-repo {{ download-url-base }}/docker-ce.repo
```

#### Install Docker Engine

1. Install Docker Engine, containerd, and Docker Compose:

   <ul class="nav nav-tabs">
    <li class="active"><a data-toggle="tab" data-target="#tab-latest">Latest</a></li>
    <li><a data-toggle="tab" data-target="#tab-version">Specific version</a></li>
   </ul>
   <div class="tab-content">
   <br>
   <div id="tab-latest" class="tab-pane fade in active" markdown="1">

   To install the latest version, run:

   ```console
   $ sudo yum install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
   ```

   If prompted to accept the GPG key, verify that the fingerprint matches
   `060A 61C5 1B55 8A7F 742B 77AA C52F EB6B 621E 9F35`, and if so, accept it.

   This command installs Docker, but it doesn't start Docker. It also creates a
   `docker` group, however, it doesn't add any users to the group by default.

   </div>
   <div id="tab-version" class="tab-pane fade" markdown="1">

   To install a specific version, start by listing the available versions in
   the repository:

   ```console
   $ yum list docker-ce --showduplicates | sort -r

   docker-ce.s390x    3:24.0.0-1.el8    docker-ce-stable
   docker-ce.s390x    3:23.0.6-1.el8    docker-ce-stable
   <...>
   ```

   The list returned depends on which repositories are enabled, and is specific
   to your version of RHEL (indicated by the `.el8` suffix in this example).

   Install a specific version by its fully qualified package name, which is
   the package name (`docker-ce`) plus the version string (2nd column),
   separated by a hyphen (`-`). For example, `docker-ce-3:24.0.0-1.el8`.

   Replace `<VERSION_STRING>` with the desired version and then run the following
   command to install:

   ```console
   $ sudo yum install docker-ce-<VERSION_STRING> docker-ce-cli-<VERSION_STRING> containerd.io docker-buildx-plugin docker-compose-plugin
   ```

   This command installs Docker, but it doesn't start Docker. It also creates a
   `docker` group, however, it doesn't add any users to the group by default.

   </div>
   <hr>
   </div>

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

{% include root-errors.md %}

#### Upgrade Docker Engine

To upgrade Docker Engine, follow the [installation instructions](#install-using-the-repository),
choosing the new version you want to install.

### Install from a package

If you can't use Docker's `rpm` repository to install Docker Engine, you can
download the `.rpm` file for your release and install it manually. You need to
download a new file each time you want to upgrade Docker Engine.

1. Go to [{{ download-url-base }}/]({{ download-url-base }}/){: target="_blank" rel="noopener" class="_" }
   and choose your version of RHEL. Then browse to `s390x/stable/Packages/`
   and download the `.rpm` file for the Docker version you want to install.

2. Install Docker Engine, changing the path below to the path where you downloaded
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

{% include root-errors.md %}

#### Upgrade Docker Engine

To upgrade Docker Engine, download the newer package files and repeat the
[installation procedure](#install-from-a-package), using `yum -y upgrade`
instead of `yum -y install`, and point to the new files.

{% include install-script.md %}

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
- Review the topics in [Develop with Docker](../../develop/index.md) to learn
  how to build new applications using Docker.
