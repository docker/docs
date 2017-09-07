---
description: Instructions for installing Docker EE on Oracle Linux
keywords: requirements, installation, oracle, ol, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/oracle/
- /engine/installation/linux/oracle/
title: Get Docker EE for Oracle Linux
---

{% assign minor-version = "17.06" %}

To get started with Docker EE on Oracle Linux, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ee).

## Prerequisites

Docker Community Edition (Docker CE) is not supported on Oracle Linux.

### Docker EE repository URL

To install Docker Enterprise Edition (Docker EE), you need to know the Docker EE
repository URL associated with your trial or subscription. These instructions
work for Docker EE for Oracle Linux and for Docker EE for Linux, which includes
access to Docker EE for all Linux distributions.To get this information:

- Go to [https://store.docker.com/my-content](https://store.docker.com/my-content).
- Each subscription or trial you have access to is listed. Click the **Setup**
  button for **Docker Enterprise Edition for Oracle Linux**.
- Copy the URL from the field labeled
  **Copy and paste this URL to download your Edition**.

Use this URL when you see the placeholder text `<DOCKER-EE-URL>`.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### OS requirements

To install Docker EE, you need the 64-bit version of Oracle Linux 7.3 running the
Red Hat Compatible kernel (RHCK) 3.10.0-514 or higher. Older versions of Oracle
Linux are not supported.

In addition, you must use the `devicemapper` storage driver if you use
Docker EE. On production systems, you must use `direct-lvm` mode, which
requires one or more dedicated block devices. Fast storage such as solid-state
media (SSD) is recommended.

> **Docker EE will not install on Oracle Linux with `selinux` enabled!**
>
> If you have `selinux` enabled and you attempt to install Docker EE 17.06.1,
> you will get an error that the `container-selinux` package cannot be found.
{:.warning }

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```bash
$ sudo yum remove docker \
                  docker-engine \
                  docker-engine-selinux
```

It's OK if `yum` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker EE package is now called `docker-ee`.

## Install Docker EE

You can install Docker EE in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the RPM package and install it manually and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

### Install using the repository

Before you install Docker EE for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker EE
from the repository.

#### Set up the repository

1.  Remove any existing Docker repositories from `/etc/yum.repos.d/`.

2.  Temporarily store the Docker EE repository URL you noted down in the
    [prerequisites](#prerequisites) in an environment variable.
    This will not persist when the current session ends.

    ```bash
    $ export DOCKERURL='<DOCKER-EE-URL>'
    ```

3.  Store your Docker EE repository URL in a `yum` variable in `/etc/yum/vars/`.
    This command relies on the variable you stored in the previous step.

    ```bash
    $ sudo sh -c 'echo "$DOCKERURL/oraclelinux" > /etc/yum/vars/dockerurl'
    ```

4.  Install required packages. `yum-utils` provides the `yum-config-manager`
    utility, and `device-mapper-persistent-data` and `lvm2` are required by the
    `devicemapper` storage driver.

    ```bash
    $ sudo yum install -y yum-utils \
      device-mapper-persistent-data \
      lvm2
    ```

5.  Use the following command to add the **stable** repository. This command
    relies on the variable you stored earlier.

    ```bash
    $ sudo yum-config-manager \
      --add-repo \
      "$DOCKERURL/oraclelinux/docker-ee.repo"
    ```

#### Install Docker EE

1.  Install the latest version of Docker EE, or go to the next step to install a
    specific version.

    ```bash
    $ sudo yum -y install docker-ee
    ```

    If this is the first time you have refreshed the package index since adding
    the Docker repositories, you will be prompted to accept the GPG key, and
    the key's fingerprint will be shown. Verify that the fingerprint matches
    `77fe da13 1a83 1d29 a418 d3e8 99e5 ff2e 7668 2bc9` and if so, accept the
    key.

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

2.  On production systems, you should install a specific version of Docker
    instead of always using the latest. List the available versions.
    This example uses the `sort -r` command to sort the results by version
    number, highest to lowest. The output is truncated.

    > **Note**: This `yum list` command only shows binary packages. To show
    > source packages as well, omit the `.x86_64` from the package name.

    ```bash
    $ yum list docker-ee.x86_64  --showduplicates | sort -r

    docker-ee.x86_64            {{ minor-version }}.ee.2-1.el7            docker-ee-stable-17.06
    ```

    The contents of the list depend upon which repositories you have enabled.
    Choose a specific version to install. The second column is the version
    string. You can use the entire version string, but **you need to include at
    least to the first hyphen**. The third column is the repository name, which
    indicates which repository the package is from and by extension its
    stability level. To install a specific version, append the version string to
    the package name and separate them by a hyphen (`-`):

    > **Note**: The version string is the package name plus the version up to
    > the first hyphen. In the example above, the fully qualified package name
    > is `docker-ee-17.06.1.ee.2`.

    ```bash
    $ sudo yum -y install <FULLY-QUALIFIED-PACKAGE-NAME>
    ```

3.  Edit `/etc/docker/daemon.json`. If it does not yet exist, create it. Assuming
    that the file was empty, add the following contents.

    ```json
    {
      "storage-driver": "devicemapper"
    }
    ```

4.  For production systems, you must use `direct-lvm` mode, which requires you
    to prepare the block devices. Follow the procedure in the
    [devicemapper storage driver guide](/engine/userguide/storagedriver/device-mapper-driver.md#configure-direct-lvm-mode-for-production){: target="_blank" class="_" }
    **before starting Docker**.

5.  Start the Docker daemon.

    ```bash
    $ sudo systemctl start docker
    ```

6.  Verify that Docker EE is installed correctly by running the `hello-world`
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

To upgrade Docker EE:

1.  If upgrading to a new major Docker EE version (such as when going from
    Docker 17.03.x to Docker 17.06.x),
    [add the new repository](#set-up-the-repository){: target="_blank" class="_" }.

2.  Run `sudo yum makecache fast`.

3.  Follow the
    [installation instructions](#install-docker), choosing the new version you want
    to install.

### Install from a package

If you cannot use the official Docker repository to install Docker EE, you can
download the `.rpm` file for your release and install it manually. You will
need to download a new file each time you want to upgrade Docker EE.

1.  Go to the Docker EE repository URL associated with your
    trial or subscription in your browser. Browse to
    `oraclelinux/7/x86_64/stable-{{ minor-version }}/Packages` and download the
    `.rpm` file for the Docker version you want to install.

2.  Install Docker EE, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo yum install /path/to/package.rpm
    ```

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

3.  Edit `/etc/docker/daemon.json`. If it does not yet exist, create it. Assuming
    that the file was empty, add the following contents.

    ```json
    {
      "storage-driver": "devicemapper"
    }
    ```

4.  For production systems, you must use `direct-lvm` mode, which requires you
    to prepare the block devices. Follow the procedure in the
    [devicemapper storage driver guide](/engine/userguide/storagedriver/device-mapper-driver.md#configure-direct-lvm-mode-for-production){: target="_blank" class="_" }
    **before starting Docker**.

5.  Start the Docker daemon.

    ```bash
    $ sudo systemctl start docker
    ```

6.  Verify that Docker EE is installed correctly by running the `hello-world`
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
[installation procedure](#install-from-a-package), using `yum -y upgrade`
instead of `yum -y install`, and pointing to the new file.

## Uninstall Docker EE

1.  Uninstall the Docker EE package:

    ```bash
    $ sudo yum remove docker-ee
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```bash
    $ sudo rm -rf /var/lib/docker
    ```

3.  If desired, remove the `devicemapper` thin pool and reformat the block
    devices that were part of it.

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](/engine/installation/linux/linux-postinstall.md)

- Continue with the [User Guide](/engine/userguide/index.md).
