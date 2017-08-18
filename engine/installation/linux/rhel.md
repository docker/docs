---
description: Instructions for installing Docker on RHEL
keywords: Docker, Docker documentation, requirements, installation, rhel, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/rhel/
- /installation/rhel/
title: Get Docker for Red Hat Enterprise Linux
---

To get started with Docker on Red Hat Enterprise Linux (RHEL), make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### Docker EE repository URL

To install Docker Enterprise Edition (Docker EE), you need to know the Docker EE
repository URL associated with your trial or subscription. To get this information:

- Go to [https://store.docker.com/?overlay=subscriptions](https://store.docker.com/?overlay=subscriptions).
- Choose **Get Details** / **Setup Instructions** within the
  **Docker Enterprise Edition for Red Hat Enterprise Linux** section.
- Copy the URL from the field labeled
  **Copy and paste this URL to download your Edition**.

Use this URL when you see the placeholder text `<DOCKER-EE-URL>`.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

Docker Community Edition (Docker CE) is not supported on Red Hat Enterprise
Linux.

### OS requirements

To install Docker, you need the 64-bit version of RHEL 7, running on an x86
hardware platform.

In addition, you must use the `devicemapper` storage driver. On production
systems, you must use `direct-lvm` mode, which requires one or more dedicated
block devices. Fast storage such as solid-state media (SSD) is recommended.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```bash
$ sudo yum remove docker \
                  docker-common \
                  container-selinux \
                  docker-selinux \
                  docker-engine
```

It's OK if `yum` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker EE package is now called `docker-ee`.

## Install Docker EE

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

1.  Remove any existing Docker repositories from `/etc/yum.repos.d/`.

2.  Store two `yum` variables in `/etc/yum/vars/`.

    - Store your EE repository URL in `/etc/yum/vars/dockerurl`. Replace
      `<DOCKER-EE-URL>` with the URL you noted down in the
      [prerequisites](#prerequisites).

      ```bash
      $ sudo sh -c 'echo "<DOCKER-EE-URL>/rhel" > /etc/yum/vars/dockerurl'
      ```

    - Store your RHEL version string in `/etc/yum/vars/dockerosversion`.
      Use the appropriate value from the following table. Most users should use
      `7`.

      | Version string | Description |
      |----------------|-------------|
      | `7`       | Unless you have specific requirements, you should use this version. Dependencies are not locked to specific versions but use the latest available version. |
      | `7.3`     | Dependencies are locked to specific packages for RHEL 7.3.                                                                                                 |
      | `7.2`     | Dependencies are locked to specific packages for RHEL 7.2.                                                                                                 |

      ```bash
      $ sudo sh -c 'echo "<VERSION-STRING>" > /etc/yum/vars/dockerosversion'
      ```

3.  Install required packages. `yum-utils` provides the `yum-config-manager`
    utility, and `device-mapper-persistent-data` and `lvm2` are required by the
    `devicemapper` storage driver.

    ```bash
    $ sudo yum install -y yum-utils device-mapper-persistent-data lvm2
    ```

4.  Enable the `extras` RHEL repository. This ensures access to the
    `container-selinux` package which is required by `docker-ee`.

    ```bash
    $ sudo yum-config-manager --enable rhel-7-server-extras-rpms
    ```

    Depending on cloud provider, you may also need to enable another repository.

    For AWS:

    ```bash
    $ sudo yum-config-manager --enable rhui-REGION-rhel-server-extras
    ```

    > **Note**: `REGION` here is literal, and does *not* represent the region
    > your machine is running in.

    For Azure:

    ```bash
    $ sudo yum-config-manager --enable rhui-rhel-7-server-rhui-extras-rpms
    ```

5.  Use the following command to add the **stable** repository:

    ```bash
    $ sudo yum-config-manager \
        --add-repo \
        <DOCKER-EE-URL>/rhel/docker-ee.repo
    ```

#### Install Docker

1.  Update the `yum` package index.

    ```bash
    $ sudo yum makecache fast
    ```

    If this is the first time you have refreshed the package index since adding
    the Docker repositories, you will be prompted to accept the GPG key, and
    the key's fingerprint will be shown. Verify that the fingerprint matches
    `DD91 1E99 5A64 A202 E859  07D6 BC14 F10B 6D08 5F96` and if so, accept the
    key.

2.  Install the latest version of Docker EE, or go to the next step to install a
    specific version.

    ```bash
    $ sudo yum -y install docker-ee
    ```

3.  On production systems, you should install a specific version of Docker
    instead of always using the latest. List the available versions.
    This example uses the `sort -r` command to sort the results by version
    number, highest to lowest, and is truncated.

    > **Note**: This `yum list` command only shows binary packages. To show
    > source packages as well, omit the `.x86_64` from the package name.

    {% assign minor-version = "17.03" %}

    ```bash
    $ yum list docker-ee.x86_64  --showduplicates |sort -r

    docker-ee.x86_64  {{ minor-version }}.0.el7                               docker-ee-stable   
    ```

    The contents of the list depend upon which repositories you have enabled,
    and will be specific to your version of RHEL (indicated by the `.el7` suffix
    on the version, in this example). Choose a specific version to install. The
    second column is the version string. The third column is the repository
    name, which indicates which repository the package is from and by extension
    extension its stability level. To install a specific version, append the
    version string to the package name and separate them by a hyphen (`-`):

    ```bash
    $ sudo yum -y install docker-ee-<VERSION_STRING>
    ```

4.  Edit `/etc/docker/daemon.json`. If it does not yet exist, create it. Assuming
    that the file was empty, add the following contents.

    ```json
    {
      "storage-driver": "devicemapper"
    }
    ```

5.  For production systems, you must use `direct-lvm` mode, which requires you
    to prepare the block devices. Follow the procedure in the
    [devicemapper storage driver guide](/engine/userguide/storagedriver/device-mapper-driver.md#configure-direct-lvm-mode-for-production){: target="_blank" class="_" }
    **before starting Docker**.

6.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

7.  Verify that Docker EE is installed correctly by running the `hello-world`
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

To upgrade Docker EE, first run `sudo yum makecache fast`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use the official Docker repository to install Docker, you can
download the `.rpm` file for your release and install it manually. You will
need to download a new file each time you want to upgrade Docker.

1.  Go to the Docker EE repository URL associated with your
    trial or subscription in your browser. Go to
    `rhel/7/x86_64/stable-{{ minor-version }}/Packages` and download the `.rpm` file
    for the Docker version you want to install.

    > **Note**: If you have trouble with `selinux` using the packages under the
    > `7` directory, try choosing the version-specific directory instead, such
    > as `7.3`.

2.  Install Docker EE, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo yum install /path/to/package.rpm
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

5.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

6.  Verify that `docker` is installed correctly by running the `hello-world`
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
    $ sudo yum -y remove docker-ee
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

- Continue to [Post-installation steps for Linux](linux-postinstall.md)

- Continue with the [User Guide](../../userguide/index.md).
