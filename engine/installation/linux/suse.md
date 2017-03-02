---
description: Instructions for installing Docker on OpenSUSE and SLES
keywords: Docker, Docker documentation, requirements, apt, installation, suse, opensuse, sles, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/SUSE/
- /engine/installation/linux/SUSE/
title: Get Docker for OpenSUSE and SLES
---

To get started with Docker on OpenSUSE or SLES, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### OS requirements

To install Docker, you need the 64-bit version one of the following:

- OpenSUSE Leap 42.x
- SLES 12.x

### Remove unofficial Docker packages

OpenSUSE's operating system repositories contain an older version of Docker,
with the package name `docker` instead of `docker-engine`. If you installed this
version of Docker on OpenSUSE or on SLES by using the OpenSUSE repositories,
remove it using the following command:

```bash
$ sudo zypper rm docker
```

The contents of `/var/lib/docker` are not removed, so any images, containers,
or volumes you created using the older version of Docker are preserved.

## Install Docker

You can install Docker in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the RPM package and install it manually and manage
  upgrades completely manually.

- Some users cannot use third-party repositories, and must rely on the version
  of Docker in the OpenSUSE or SLES repositories. This version of Docker may be
  out of date. Those users should consult the OpenSuSE or SLES documentation and
  not follow these procedures.

### Install using the repository

Before you install Docker for the first time on a new host machine, you need to
set up the Docker repository. Afterward, you can install, update, or downgrade
Docker from the repository.

#### Set up the repository

1.  Use the following command to set up the **stable** repository:

    ```bash
    $ sudo zypper addrepo \
        https://yum.dockerproject.org/repo/main/opensuse/13.2/ \
        docker-main
    ```

2.  **Optional**: Enable the **testing** repository. You can enable it alongside
    the stable repository. **Do not use unstable repositories on on production
    systems or for non-testing workloads.**

    > **Warning**: If you have both stable and unstable repositories enabled,
    > updating without specifying a version in the `zypper install` or
    > `zypper update` command will always install the highest possible version,
    > which will almost certainly be an unstable one.


    ```bash
    $ sudo zypper addrepo \
        https://yum.dockerproject.org/repo/testing/opensuse/13.2/ \
        docker-testing
    ```

    You can disable a repository at any time by running the `zypper rmrepo`
    command. The following command disables the `testing` repository.

    ```bash
    $ sudo zypper removerepo docker-testing
    ```

#### Install Docker

1.  Update the `zypper` package index.

    ```bash
    $ sudo zypper refresh
    ```

    If this is the first time you have refreshed the package index since adding
    the Docker repositories, you will be prompted to accept the GPG key, and
    the key's fingerprint will be shown. Verify that the fingerprint matches
    `58118E89F3A912897C070ADBF76221572C52609D` and if so, accept the key.

2.  Install the latest version of Docker, or go to the next step to install a
    specific version.

    ```bash
    $ sudo zypper install docker-engine
    ```

    > **Warning**: If you have both stable and unstable repositories enabled,
    > installing or updating Docker without specifying a version in the
    > `zypper install` or `zypper update` command will always install the highest
    > available version, which will almost certainly be an unstable one.

    The RPM will install, but you will receive the following error during the
    post-installation procedure, because Docker cannot start the service
    automatically:

    ```none
    Additional rpm output:
    /var/tmp/rpm-tmp.YGySzA: line 1: fg: no job control
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
    $ zypper search -s --match-exact -t package docker-engine

      Loading repository data...
      Reading installed packages...

      S | Name          | Type    | Version                               | Arch   | Repository    
      --+---------------+---------+---------------------------------------+--------+---------------
        | docker-engine | package | 1.13.0-1                              | x86_64 | docker-main
        | docker-engine | package | 1.12.6-1                              | x86_64 | docker-main   
        | docker-engine | package | 1.12.5-1                              | x86_64 | docker-main   
    ```

    The contents of the list depend upon which repositories you have enabled.
    Choose a specific version to
    install. The third column is the version string. The fifth column is the
    repository name, which indicates which repository the package is from and by
    extension its stability level. To install a specific version, append the
    version string to the package name and separate them by a hyphen (`-`):

    ```bash
    $ sudo zypper install docker-engine-<VERSION_STRING>
    ```

    The RPM will install, but you will receive the following error during the
    post-installation procedure, because Docker cannot start the service
    automatically:

    ```none
    Additional rpm output:
    /var/tmp/rpm-tmp.YGySzA: line 1: fg: no job control
    ```

    Start Docker:

    ```bash
    $ sudo service docker start
    ```

4.  Verify that `docker` is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](linux-postinstall.md) to allow non-privileged
users to run Docker commands and for other optional configuration steps.

#### Upgrade Docker

To upgrade Docker, first run `sudo zypper refresh`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker, you can download the
`.rpm` file for your release and install it manually. You will need to download
a new file each time you want to upgrade Docker.

1.  Go to [https://yum.dockerproject.org/repo/main/opensuse/13.2/Packages/](https://yum.dockerproject.org/repo/main/opensuse/13.2/Packages/)
    and download the `.rpm` file for the Docker version you want to install.

    > **Note**: To install a testing version, change the word `main` in the
    > URL to `testing`. Do not use unstable versions of Docker in production
    > or for non-testing workloads.

2.  Install Docker, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo yum -y install /path/to/package.rpm
    ```

    The RPM will install, but you will receive the following error during the
    post-installation procedure, because Docker cannot start the service
    automatically:

    ```none
    Additional rpm output:
    /var/tmp/rpm-tmp.YGySzA: line 1: fg: no job control
    ```

    Start Docker:

    ```bash
    $ sudo service docker start
    ```

3.  Verify that `docker` is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker

To upgrade Docker, download the newer package file and repeat the
[installation procedure](#install-from-a-package), using `zypper update`
instead of `zypper install`, and pointing to the new file.

## Uninstall Docker

1.  Uninstallation using `zypper rm` fails. Uninstall the Docker package using
    the following command:

    ```bash
    $ sudo rpm -e --noscripts docker-engine
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
