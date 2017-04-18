---
description: Instructions for installing Docker on CentOS
keywords: Docker, Docker documentation, requirements, apt, installation, centos, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/centos/
title: Get Docker for CentOS
---

{% assign minor-version = "17.03" %}

To get started with Docker on CentOS, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### Docker EE customers

To install Docker Enterprise Edition (Docker EE), you need to know the Docker EE
repository URL associated with your trial or subscription. To get this information:

- Go to [https://store.docker.com/?overlay=subscriptions](https://store.docker.com/?overlay=subscriptions).
- Choose **Get Details** / **Setup Instructions** within the
  **Docker Enterprise Edition for CentOS** section.
- Copy the URL from the field labeled
  **Copy and paste this URL to download your Edition**.

Where the installation instructions differ for Docker EE and Docker CE, use this
URL when you see the placeholder text `<DOCKER-EE-URL>`.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### OS requirements

To install Docker, you need the 64-bit version of CentOS 7.

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
networks, are preserved. The Docker CE package is now called `docker-ce`, and
the Docker EE package is now called `docker-ee`.

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

Repository set-up instructions are different for [Docker CE](#docker-ce) and
[Docker EE](#docker-ee).

##### Docker CE

{% assign download-url-base = "https://download.docker.com/linux/centos" %}

1.  Install `yum-utils`, which provides the `yum-config-manager` utility:

    ```bash
    $ sudo yum install -y yum-utils
    ```

2.  Use the following command to set up the **stable** repository. You always
    need the **stable** repository, even if you want to install **edge** builds
    as well.

    ```bash
    $ sudo yum-config-manager \
        --add-repo \
        {{ download-url-base }}/docker-ce.repo
    ```

3.  **Optional**: Enable the **edge** repository. This repository is included
    in the `docker.repo` file above but is disabled by default. You can enable
    it alongside the stable repository.

    ```bash
    $ sudo yum-config-manager --enable docker-ce-edge
    ```

    You can disable the **edge** repository by running the `yum-config-manager`
    command with the `--disable` flag. To re-enable it, use the
    `--enable` flag. The following command disables the **edge** repository.

    ```bash
    $ sudo yum-config-manager --disable docker-ce-edge
    ```

    [Learn about **stable** and **edge** builds](/engine/installation/).

##### Docker EE

1.  Remove any existing Docker repositories from `/etc/yum.repos.d/`.

2.  Store your Docker EE repository URL in a `yum` variable in `/etc/yum/vars/`.
    Replace `<DOCKER-EE-URL>` with the URL you noted down in the
    [prerequisites](#prerequisites).

    ```bash
    $ sudo sh -c 'echo "<DOCKER-EE-URL>" > /etc/yum/vars/dockerurl'
    ```

3.  Install `yum-utils`, which provides the `yum-config-manager` utility:

    ```bash
    $ sudo yum install -y yum-utils
    ```

4.  Use the following command to add the **stable** repository:

    ```bash
    $ sudo yum-config-manager \
        --add-repo \
        <DOCKER-EE-URL>/docker-ee.repo
    ```

#### Install Docker

1.  Update the `yum` package index.

    ```bash
    $ sudo yum makecache fast
    ```

    If this is the first time you have refreshed the package index since adding
    the Docker repositories, you will be prompted to accept the GPG key, and
    the key's fingerprint will be shown. Verify that the fingerprint is
    correct, and if so, accept the key.

    | Docker Edition | Fingerprint                                          |
    |----------------|------------------------------------------------------|
    | Docker CE      | `060A 61C5 1B55 8A7F 742B  77AA C52F EB6B 621E 9F35` |
    | Docker EE      | `DD91 1E99 5A64 A202 E859  07D6 BC14 F10B 6D08 5F96` |

2.  Install the latest version of Docker, or go to the next step to install a
    specific version.

    | Docker Edition | Command                             |
    |----------------|-------------------------------------|
    | Docker CE      | `sudo yum install docker-ce`        |
    | Docker EE      | `sudo yum install docker-ee`        |

    > **Warning**: If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `yum install` or
    > `yum update` command will always install the highest possible version,
    > which may not be appropriate for your stability needs.

3.  On production systems, you should install a specific version of Docker
    instead of always using the latest. List the available versions. This
    example uses the `sort -r` command to sort the results by version number,
    highest to lowest, and is truncated.

    > **Note**: This `yum list` command only shows binary packages. To show
    > source packages as well, omit the `.x86_64` from the package name.

    ```bash
    $ yum list docker-ce.x86_64  --showduplicates |sort -r

    docker-ce.x86_64  {{ minor-version }}.0.el7                               docker-ce-stable  
    ```

    The contents of the list depend upon which repositories are enabled, and
    will be specific to your version of CentOS (indicated by the `.el7` suffix
    on the version, in this example). Choose a specific version to install. The
    second column is the version string. The third column is the repository
    name, which indicates which repository the package is from and by extension
    its stability level. To install a specific version, append the version
    string to the package name and separate them by a hyphen (`-`):

    | Docker Edition | Command                                       |
    |----------------|-----------------------------------------------|
    | Docker CE      | `sudo yum install docker-ce-<VERSION>`        |
    | Docker EE      | `sudo yum install docker-ee-<VERSION>`        |

4.  Start Docker.

    ```bash
    $ sudo systemctl start docker
    ```

5.  Verify that `docker` is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

6.  If you don’t want to use sudo when you use the docker command, you can add 
    users to the `docker` group.

    ```bash
    $ sudo usermod -aG docker <username>
    ```
    
    Now apply group settings by running the following command.

    ```bash
    $ newgrp docker
    ```

    Docker is installed and running.


#### Upgrade Docker

To upgrade Docker, first run `sudo yum makecache fast`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker, you can download the
`.rpm` file for your release and install it manually. You will need to download
a new file each time you want to upgrade Docker.

1.  This step is different for Docker CE and Docker EE.

    - **Docker CE**: Go to
      [{{ download-url-base }}/7/x86_64/stable/Packages/]({{ download-url-base }}/7/x86_64/stable/Packages/)
      and download the `.rpm` file for the Docker version you want to install.

      > **Note**: To install an **edge**  package, change the word
      > `stable` in the > URL to `edge`.
      > [Learn about **stable** and **edge** channels](/engine/installation/).

    - **Docker EE**: Go to the Docker EE repository URL associated with your
      trial or subscription in your browser. Go to
      `7/x86_64/stable-{{ minor-version }}/Packages/` and download the `.rpm`
      file for the Docker version you want to install.

2.  Install Docker, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo yum install /path/to/package.rpm
    ```

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

Docker is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Post-installation steps for Linux](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker

To upgrade Docker, download the newer package file and repeat the
[installation procedure](#install-from-a-package), using `yum -y upgrade`
instead of `yum -y install`, and pointing to the new file.

## Uninstall Docker

1.  Uninstall the Docker package:

    | Docker Edition | Command                        |
    |----------------|--------------------------------|
    | Docker CE      | `sudo yum remove docker-ce`    |
    | Docker EE      | `sudo yum remove docker-ee`    |

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
