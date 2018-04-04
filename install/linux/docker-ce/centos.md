---
description: Instructions for installing Docker CE on CentOS
keywords: requirements, apt, installation, centos, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/centos/
- /engine/installation/linux/docker-ce/centos/
title: Get Docker CE for CentOS
toc_max: 4
---

To get started with Docker CE on CentOS, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ce).

## Prerequisites

### Docker EE customers

To install Docker Enterprise Edition (Docker EE), go to
[Get Docker EE for CentOS](/install/linux/docker-ee/centos/)
**instead of this topic**.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### OS requirements

To install Docker CE, you need a maintained version of CentOS 7. Archived
versions aren't supported or tested.

The `centos-extras` repository must be enabled. This repository is enabled by
default, but if you have disabled it, you need to
[re-enable it](https://wiki.centos.org/AdditionalResources/Repositories){: target="_blank" class="_" }.

The `overlay2` storage driver is recommended.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```bash
$ sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine
```

It's OK if `yum` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker CE package is now called `docker-ce`.

## Install Docker CE

You can install Docker CE in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the RPM package and
  [install it manually](#install-from-a-package) and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

- In testing and development environments, some users choose to use automated
  [convenience scripts](#install-using-the-convenience-script) to install Docker.

### Install using the repository

Before you install Docker CE for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker
from the repository.

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/centos" %}

1.  Install required packages. `yum-utils` provides the `yum-config-manager`
    utility, and `device-mapper-persistent-data` and `lvm2` are required by the
    `devicemapper` storage driver.

    ```bash
    $ sudo yum install -y yum-utils \
      device-mapper-persistent-data \
      lvm2
    ```

2.  Use the following command to set up the **stable** repository. You always
    need the **stable** repository, even if you want to install builds from the
    **edge** or **test** repositories as well.

    ```bash
    $ sudo yum-config-manager \
        --add-repo \
        {{ download-url-base }}/docker-ce.repo
    ```

3.  **Optional**: Enable the **edge** and **test** repositories. These
    repositories are included in the `docker.repo` file above but are disabled
    by default. You can enable them alongside the stable repository.

    ```bash
    $ sudo yum-config-manager --enable docker-ce-edge
    ```

    ```bash
    $ sudo yum-config-manager --enable docker-ce-test
    ```

    You can disable the **edge** or **test** repository by running the
    `yum-config-manager` command with the `--disable` flag. To re-enable it, use
    the `--enable` flag. The following command disables the **edge** repository.

    ```bash
    $ sudo yum-config-manager --disable docker-ce-edge
    ```

    > **Note**: Starting with Docker 17.06, stable releases are also pushed to
    > the **edge** and **test** repositories.

    [Learn about **stable** and **edge** builds](/install/index.md).

#### Install Docker CE

1.  Install the _latest version_ of Docker CE, or go to the next step to install a specific version:

    ```bash
    $ sudo yum install docker-ce
    ```

    If prompted to accept the GPG key, verify that the fingerprint matches
    `060A 61C5 1B55 8A7F 742B 77AA C52F EB6B 621E 9F35`, and if so, accept it.

    > Got multiple Docker repositories?
    >
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `yum install` or
    > `yum update` command always installs the highest possible version,
    > which may not be appropriate for your stability needs.

    Docker is installed but not started. The `docker` group is created, but no users are added to the group.

2.  To install a _specific version_ of Docker CE, list the available versions
    in the repo, then select and install:

    a. List and sort the versions available in your repo. This example sorts
       results by version number, highest to lowest, and is truncated:

    ```bash
    $ yum list docker-ce --showduplicates | sort -r

    docker-ce.x86_64            {{ site.docker_ce_stable_version }}.0.ce-1.el7.centos             docker-ce-stable
    ```

    The list returned depends on which repositories are enabled, and is specific
    to your version of CentOS (indicated by the `.el7` suffix in this example).

    b. Install a specific version by its fully qualified package name, which is
       the package name (`docker-ce`) plus the version string (2nd column) up to
       the first hyphen, separated by a hyphen (`-`), for example,
       `docker-ce-18.03.0.ce`.

    ```bash
    $ sudo yum install docker-ce-<VERSION STRING>
    ```

    Docker is installed but not started. The `docker` group is created, but no users are added to the group.

3.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

4.  Verify that `docker` is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker CE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Linux postinstall](/install/linux/linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker CE

To upgrade Docker CE, follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker, you can download the
`.rpm` file for your release and install it manually. You need to download
a new file each time you want to upgrade Docker.

1.  Go to
    [{{ download-url-base }}/7/x86_64/stable/Packages/]({{ download-url-base }}/7/x86_64/stable/Packages/)
    and download the `.rpm` file for the Docker version you want to install.

    > **Note**: To install an **edge**  package, change the word
    > `stable` in the above URL to `edge`.
    > [Learn about **stable** and **edge** channels](/install/index.md).

2.  Install Docker CE, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo yum install /path/to/package.rpm
    ```

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

3.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

4.  Verify that `docker` is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker CE is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker CE

To upgrade Docker CE, download the newer package file and repeat the
[installation procedure](#install-from-a-package), using `yum -y upgrade`
instead of `yum -y install`, and pointing to the new file.

{% include install-script.md %}

## Uninstall Docker CE

1.  Uninstall the Docker package:

    ```bash
    $ sudo yum remove docker-ce
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
