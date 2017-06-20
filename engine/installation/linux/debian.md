---
description: Instructions for installing Docker on Debian
keywords: Docker, Docker documentation, requirements, apt, installation, debian, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/debian/
- /engine/installation/linux/raspbian/
title: Get Docker for Debian
---

{% assign minor-version = "17.03" %}

To get started with Docker on Debian, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### Docker EE customers

Docker EE is not supported on Debian. For a list of supported operating systems
and distributions for different Docker editions, see
[Docker variants](/engine/installation/#docker-variants).

### OS requirements

To install Docker, you need the 64-bit version of one of these Debian or
Raspbian versions:

- Stretch (testing)
- Jessie 8.0 (LTS) / Raspbian Jessie
- Wheezy 7.7 (LTS)

Docker CE is supported on both `x86_64` and `armhf` architectures for Jessie and
Stretch.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them:

```bash
$ sudo apt-get remove docker docker-engine
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker CE package is now called `docker-ce`.

### Extra steps for Wheezy 7.7

- You need at least version 3.10 of the Linux kernel. Debian Wheezy ships with
  version 3.2, so you may need to
  [update the kernel](https://wiki.debian.org/HowToUpgradeKernel){: target="_blank" class="_" }.
  To check your kernel version:

  ```bash
  $ uname -r
  ```

- Enable the `backports` repository. See the
  [Debian documentation](https://backports.debian.org/Instructions/){: target="_blank" class"_"}.

## Install Docker CE

You can install Docker CE in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the DEB package and install it manually and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

### Install using the repository

Before you install Docker CE for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker
from the repository.

#### Set up the repository

{% assign download-url-base = "https://download.docker.com/linux/debian" %}

1.  Install packages to allow `apt` to use a repository over HTTPS:

    **Jessie or Stretch**:

    ```bash
    $ sudo apt-get install \
         apt-transport-https \
         ca-certificates \
         curl \
         gnupg2 \
         software-properties-common
    ```

    **Wheezy**:

    ```bash
    $ sudo apt-get install \
         apt-transport-https \
         ca-certificates \
         curl \
         python-software-properties
    ```

2.  Add Docker's official GPG key:

    ```bash
    $ curl -fsSL {{ download-url-base}}/gpg | sudo apt-key add -
    ```

    Verify that the key ID is `9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88`.

    ```bash
    $ sudo apt-key fingerprint 0EBFCD88

    pub   4096R/0EBFCD88 2017-02-22
          Key fingerprint = 9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
    uid                  Docker Release (CE deb) <docker@docker.com>
    sub   4096R/F273FCD8 2017-02-22
    ```

3.  Use the following command to set up the **stable** repository. You always
    need the **stable** repository, even if you want to install **edge** builds
    as well.

    > **Note**: The `lsb_release -cs` sub-command below returns the name of your
    > Debian distribution, such as `jessie`.

    To also add the **edge** repository, add `edge` after `stable` on the last
    line of the command.

    **amd64**:

    ```bash
    $ sudo add-apt-repository \
       "deb [arch=amd64] {{ download-url-base }} \
       $(lsb_release -cs) \
       stable"
    ```

    **armhf**:

    ```bash
    $ echo "deb [arch=armhf] {{ download-url-base }} \
         $(lsb_release -cs) stable" | \
        sudo tee /etc/apt/sources.list.d/docker.list
    ```

4.  **Wheezy only**: The version of `add-apt-repository` on Wheezy adds a `deb-src`
    repository that does not exist. You need to comment out this repository or
    running `apt-get update` will fail. Edit `/etc/apt/sources.list`. Find the
    line like the following, and comment it out or remove it:

    ```none
    deb-src [arch=amd64] https://download.docker.com/linux/debian wheezy stable
    ```

    Save and exit the file.

    [Learn about **stable** and **edge** channels](/engine/installation/).

#### Install Docker CE

> **NOTE**: On Debian for ARM you can continue following this step. For Raspbian,
  scroll down to follow its specific steps.

1.  Update the `apt` package index.

    ```bash
    $ sudo apt-get update
    ```

2.  Install the latest version of Docker, or go to the next step to install a
    specific version. Any existing installation of Docker is replaced.

    Use this command to install the latest version of Docker:

    ```bash
    $ sudo apt-get install docker-ce
    ```

    > **Warning**:
    > If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `apt-get install` or
    > `apt-get update` command will always install the highest possible version,
    > which may not be appropriate for your stability needs.
    {:.warning}

3.  On production systems, you should install a specific version of Docker
    instead of always using the latest. This output is truncated. List the
    available versions:

    ```bash
    $ apt-cache madison docker-ce

    docker-ce | {{ minor-version }}.0~ce-0~debian-jessie | {{ download-url-base}} jessie/stable amd64 Packages
    ```

    The contents of the list depend upon which repositories are enabled,
    and will be specific to your version of Debian (indicated by the `jessie`
    suffix on the version, in this example). Choose a specific version to
    install. The second column is the version string. The third column is the
    repository name, which indicates which repository the package is from and
    by extension its stability level. To install a specific version, append the
    version string to the package name and separate them by an equals sign (`=`):

    ```bash
    $ sudo apt-get install docker-ce=<VERSION_STRING>
    ```

    The Docker daemon starts automatically.

4.  Verify that Docker CE is installed correctly by running the `hello-world`
    image.

    **amd64**:

    ```bash
    $ sudo docker run hello-world
    ```

    **armhf**:

    ```bash
    $ sudo docker run armhf/hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker CE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Linux postinstall](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker CE

To upgrade Docker, first run `sudo apt-get update`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.


### Install on Raspbian (Raspberry Pi)
>**Warning**: This isn't necessary if you used the recommended
>`$ curl -sSL https://get.docker.com | sh` command!

Once you have added the Docker repo to `/etc/apt/sources.list.d/`, you should
see `docker.list` if you:

```bash
$ ls /etc/apt/sources.list.d/
```

And the contents of the `docker.list` should read:

`deb [arch=armhf] https://apt.dockerproject.org/repo raspbian-jessie main`

If you don't see that in `docker.list`, then either comment the line out, or
`rm` the `docker.list` file.

Once you have verified that you have the correct repository, you may continue
installing Docker.

1.  Update the `apt` package index.

    ```bash
    $ sudo apt-get update
    ```
2.  Install the latest version of Docker, or go to the next step to install a
    specific version. Any existing installation of Docker is replaced.

    Use this command to install the latest version of Docker:

    ```bash
    $ sudo apt-get install docker-engine
    ```
    > **NOTE**: By default, Docker on Raspbian is Docker Community Edition, so
    > there is no need to specify docker-ce.

    > **NOTE**: If `curl -sSL https://get.docker.com | sh` isn't used,
    > then docker won't have auto-completion! You'll have to add it manually.

3.  Verify that Docker is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hypriot/armhf-hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

#### (Optional) Install Docker Compose for Raspbian

This functionality is provided by [Hypriot](https://blog.hypriot.com/). Add the Hypriot repo:

```bash
curl -s https://packagecloud.io/install/repositories/Hypriot/Schatzkiste/script.deb.sh | sudo bash
```

Install `docker-compose`:

```bash
sudo apt-get install docker-compose
```

### Install from a package

If you cannot use Docker's repository to install Docker CE, you can download the
`.deb` file for your release and install it manually. You will need to download
a new file each time you want to upgrade Docker.

1.  Go to [{{ download-url-base }}/dists/]({{ download-url-base }}/dists/),
    choose your Debian version, browse to `stable/pool/stable/`, choose either
    `amd64` or `armhf`,and download the `.deb` file for the Docker version you
    want to install and for your version of Debian.

    > **Note**: To install an **edge**  package, change the word
    > `stable` in the  URL to `edge`.
    > [Learn about **stable** and **edge** channels](/engine/installation/).

2.  Install Docker CE, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo dpkg -i /path/to/package.deb
    ```

    The Docker daemon starts automatically.

3.  Verify that Docker CE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker CE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Post-installation steps for Linux](linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps.

#### Upgrade Docker

To upgrade Docker, download the newer package file and repeat the
[installation procedure](#install-from-a-package), pointing to the new file.

## Uninstall Docker

1.  Uninstall the Docker package:

    ```bash
    $ sudo apt-get purge docker-ce
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
