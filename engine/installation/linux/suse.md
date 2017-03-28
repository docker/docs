---
description: Instructions for installing Docker on OpenSUSE and SLES
keywords: Docker, Docker documentation, requirements, apt, installation, suse, opensuse, sles, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/SUSE/
title: Get Docker for and SLES
---

{% assign minor-version = "17.03" %}

To get started with Docker on SUSE Linux Enterprise Server (SLES), make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### Docker EE URL

To install Docker Enterprise Edition (Docker EE), you need to know the Docker EE
repository URL associated with your trial or subscription. To get this information:

- Go to [https://store.docker.com/?overlay=subscriptions](https://store.docker.com/?overlay=subscriptions).
- Choose **Get Details** / **Setup Instructions** within the
  **Docker Enterprise Edition for SUSE Linux Enterprise Server** section.
- Copy the URL from the field labeled
  **Copy and paste this URL to download your Edition**.

Use this URL when you see the placeholder text `<DOCKER-EE-URL>`.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

Docker Community Edition (Docker CE) is not supported on SLES.

### OS requirements

To install Docker, you need the 64-bit version of SLES 12.x.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```bash
$ sudo zypper rm docker docker-engine
```

If removal of the `docker-engine` package fails, use the following command
instead:

```bash
$ sudo rpm -e docker-engine
```

It's OK if `zypper` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker EE package is now called `docker-ee`.

## Install Docker

You can install Docker in different ways, depending on your needs:

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

Use the following command to set up the **stable** repository, using the
Docker EE repository URL you located in the [prerequisites](#prerequisites).

```bash
$ sudo zypper addrepo \
    <DOCKER-EE-URL>/12.3/x86_64/stable-{{ minor-version }} \
    docker-ee-stable
```

#### Install Docker EE

1.  Update the `zypper` package index.

    ```bash
    $ sudo zypper refresh
    ```

    If this is the first time you have refreshed the package index since adding
    the Docker repositories, you will be prompted to accept the GPG key, and
    the key's fingerprint will be shown. Verify that the fingerprint matches
    `77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9` and if so, accept the
    key.

2.  Install the latest version of Docker EE, or go to the next step to install a
    specific version.

    ```bash
    $ sudo zypper install docker-ee
    ```

    Start Docker:

    ```bash
    $ sudo service docker start
    ```

3.  On production systems, you should install a specific version of Docker
    instead of always using the latest. List the available versions. The
    following example only lists binary packages and is truncated. To also list
    source packages, omit the `-t package` flag from the command.

    ```bash
    $ zypper search -s --match-exact -t package docker-ee

      Loading repository data...
      Reading installed packages...

      S | Name          | Type    | Version                               | Arch   | Repository    
      --+---------------+---------+---------------------------------------+--------+---------------
        | docker-ee     | package | {{ minor-version }}-1                 | x86_64 | docker-ee-stable
    ```

    The contents of the list depend upon which repositories you have enabled.
    Choose a specific version to install. The third column is the version
    string. The fifth column is the repository name, which indicates which
    repository the package is from and by extension its stability level. To
    install a specific version, append the version string to the package name
    and separate them by a hyphen (`-`):

    ```bash
    $ sudo zypper install docker-ee-<VERSION_STRING>
    ```

    Start Docker:

    ```bash
    $ sudo service docker start
    ```

4.  Verify that Docker EE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker EE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Linux postinstall](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker EE

To upgrade Docker EE, first run `sudo zypper refresh`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use the official Docker repository to install Docker, you can
download the `.rpm` file for your release and install it manually. You will
need to download a new file each time you want to upgrade Docker.

1.  Go to the Docker EE repository URL associated with your
    trial or subscription in your browser. Go to
    `12.3/x86_64/stable-{{ minor-version }}` and download the `.rpm` file for
    the Docker version you want to install.

2.  Import Docker's official GPG key:

    ```bash
    $ sudo rpm --import <DOCKER-EE-URL>/gpg
    ```

3.  Install Docker EE, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo zypper install /path/to/package.rpm
    ```

    Start Docker:

    ```bash
    $ sudo service docker start
    ```

4.  Verify that Docker EE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker EE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Post-installation steps for Linux](linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps.

#### Upgrade Docker EE

To upgrade Docker EE, download the newer package file and repeat the
[installation procedure](#install-from-a-package), using `zypper update`
instead of `zypper install`, and pointing to the new file.

## Uninstall Docker

1.  Uninstall the Docker EE package using the following command.

    ```bash
    $ sudo zypper rm docker-ee
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
