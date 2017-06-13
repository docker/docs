---
description: Instructions for installing Docker on Fedora
keywords: Docker, Docker documentation, requirements, apt, installation, fedora, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/fedora/
title: Get Docker for Fedora
---

{% assign minor-version = "17.03" %}

To get started with Docker on Fedora, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### Docker EE customers

Docker EE is not supported on Fedora. For a list of supported operating systems
and distributions for different Docker editions, see
[Docker variants](/engine/installation/#docker-variants).

### OS requirements

To install Docker, you need the 64-bit version of one of these Fedora versions:

- 24
- 25

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```bash
$ sudo dnf remove docker \
                  docker-common \
                  container-selinux \
                  docker-selinux \
                  docker-engine
```

It's OK if `dnf` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker CE package is now called `docker-ce`.

## Install Docker CE

You can install Docker CE in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the RPM package and install it manually and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

### Install using the repository

Before you install Docker for the first time on a new host machine, you need to
set up the Docker repository. Afterward, you can install and update Docker from
the repository.

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/fedora" %}

1.  Install the `dnf-plugins-core` package which provides the commands to manage
    your DNF repositories from the command line.

    ```bash
    $ sudo dnf -y install dnf-plugins-core
    ```

2.  Use the following command to set up the **stable** repository. You always
    need the **stable** repository, even if you want to install **edge** builds
    as well

    ```bash
    $ sudo dnf config-manager \
        --add-repo \
        {{ download-url-base }}/docker-ce.repo
    ```

3.  **Optional**: Enable the **edge** repository. This repository is included
    in the `docker.repo` file above but is disabled by default. You can enable
    it alongside the stable repository.

    ```bash
    $ sudo dnf config-manager --set-enabled docker-ce-edge
    ```

    You can disable the **edge** repository by running the `dnf config-manager`
    command with the `--disable` flag. To re-enable it, use the
    `--enable` flag. The following command disables the **edge** repository.

    ```bash
    $ sudo dnf config-manager --set-disabled docker-ce-edge
    ```

    [Learn about **stable** and **edge** channels](/engine/installation/).

#### Install Docker CE

1.  Update the `dnf` package index.

    ```bash
    $ sudo dnf makecache fast
    ```

    If this is the first time you have refreshed the package index since adding
    the Docker repositories, you will be prompted to accept the GPG key, and
    the key's fingerprint will be shown. Verify that the fingerprint matches
    `060A 61C5 1B55 8A7F 742B  77AA C52F EB6B 621E 9F35` and if so, accept the
    key.

2.  Install the latest version of Docker, or go to the next step to install a
    specific version.

    ```bash
    $ sudo dnf install docker-ce
    ```

    > **Warning**:
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `dnf install` or
    > `dnf update` command will always install the highest possible version,
    > which may not be appropriate for your stability needs.
    {:.warning}

3.  On production systems, you should install a specific version of Docker
    instead of always using the latest. List the available versions. This
    example uses the `sort -r` command to sort the results by version number,
    highest to lowest, and is truncated.

    > **Note**: This `dnf list` command only shows binary packages. To show
    > source packages as well, omit the `.x86_64` from the package name.

    ```bash
    $ dnf list docker-ce.x86_64  --showduplicates |sort -r

    docker-ce.x86_64  {{ minor-version }}.0.fc24                               docker-ce-stable  
    ```

    The contents of the list depend upon which repositories are enabled, and
    will be specific to your version of Fedora (indicated by the `.fc24` suffix
    on the version, in this example). Choose a specific version to install. The
    second column is the version string. The third column is the repository
    name, which indicates which repository the package is from and by extension
    its stability level. To install a specific version, append the version
    string to the package name and separate them by a hyphen (`-`):

    ```bash
    $ sudo dnf -y install docker-ce-<VERSION>
    ```

4.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

5.  Verify that Docker CE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](linux-postinstall.md) to allow non-privileged
users to run Docker commands and for other optional configuration steps.

#### Upgrade Docker CE

To upgrade Docker, first run `sudo dnf makecache fast`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker, you can download the
`.rpm` file for your release and install it manually. You will need to download
a new file each time you want to upgrade Docker.

1.  Go to [{{ download-url-base }}/]({{ download-url-base }}/) and choose your
    version of Fedora. Go to `x86_64/stable/Packages/`
    and download the `.rpm` file for the Docker version you want to install.

    > **Note**: To install an **edge**  package, change the word
    > `stable` in the URL to `edge`.

2.  Install Docker, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo dnf -y install /path/to/package.rpm
    ```

3.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

4.  Verify that Docker CE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker CE is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker

To upgrade Docker, download the newer package file and repeat the
[installation procedure](#install-from-a-package), using `dnf -y upgrade`
instead of `dnf -y install`, and pointing to the new file.

## Uninstall Docker

1.  Uninstall the Docker package:

    ```bash
    $ sudo dnf remove docker-ce
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```bash
    $ sudo rm -rf /var/lib/docker
    ```

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](linux-postinstall.md)

- Continue with the [User Guide](../../userguide/index.md).
