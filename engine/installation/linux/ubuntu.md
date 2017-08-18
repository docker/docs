---
description: Instructions for installing Docker on Ubuntu
keywords: Docker, Docker documentation, requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
- /engine/installation/linux/ubuntulinux/
title: Get Docker for Ubuntu
---

{% assign minor-version = "17.03" %}

To get started with Docker on Ubuntu, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker).

## Prerequisites

### Docker EE customers

To install Docker Enterprise Edition (Docker EE), you need to know the Docker EE
repository URL associated with your trial or subscription. To get this information:

- Go to [https://store.docker.com/?overlay=subscriptions](https://store.docker.com/?overlay=subscriptions).
- Choose **Get Details** / **Setup Instructions** within the
  **Docker Enterprise Edition for Ubuntu** section.
- Copy the URL from the field labeled
  **Copy and paste this URL to download your Edition**.

Where the installation instructions differ for Docker EE and Docker CE, use this
URL when you see the placeholder text `<DOCKER-EE-URL>`.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### OS requirements

To install Docker, you need the 64-bit version of one of these Ubuntu versions:

- Yakkety 16.10
- Xenial 16.04 (LTS)
- Trusty 14.04 (LTS)

Docker CE is supported on both `x86_64` and `armhf` architectures.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them:

```bash
$ sudo apt-get remove docker docker-engine
```

It's OK if `apt-get` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker CE package is now called `docker-ce`, and
the Docker EE package is now called `docker-ee`.

### Recommended extra packages for Trusty 14.04

Unless you have a strong reason not to, install the
`linux-image-extra-*` packages, which allow Docker to use the `aufs` storage
drivers.

```bash
$ sudo apt-get update

$ sudo apt-get install \
    linux-image-extra-$(uname -r) \
    linux-image-extra-virtual
```

## Install Docker

You can install Docker in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the DEB package and install it manually and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

### Install using the repository

Before you install Docker for the first time on a new host machine, you need to
set up the Docker repository. Afterward, you can install and update Docker from
the repository.

#### Set up the repository

The procedure for setting up the repository is different for Docker CE and
Docker EE.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-group="ce" data-target="#ce-repo-setup">Docker CE</a></li>
  <li><a data-toggle="tab" data-group="ee" data-target="#ee-repo-setup">Docker EE</a></li>
</ul>
<div class="tab-content">
  <div id="ce-repo-setup" class="tab-pane fade in active" markdown="1">
  {% assign download-url-base = "https://download.docker.com/linux/ubuntu" %}

  1.  Install packages to allow `apt` to use a repository over HTTPS:

      ```bash
      $ sudo apt-get install \
          apt-transport-https \
          ca-certificates \
          curl \
          software-properties-common
      ```

  2.  Add Docker's official GPG key:

      ```bash
      $ curl -fsSL {{ download-url-base}}/gpg | sudo apt-key add -
      ```

      Verify that the key fingerprint is `9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88`.

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

      > **Note**: The `lsb_release -cs` sub-command below returns the
      name of your Ubuntu distribution, such as `xenial`. Sometimes,
      in a distribution like Linux Mint, you might have to
      change `$(lsb_release -cs)` to your parent Ubuntu distribution.
      For example: If you are using `Linux Mint Rafaela`, you could use
      `trusty`.

      **amd64**:

      ```bash
      $ sudo add-apt-repository \
         "deb [arch=amd64] {{ download-url-base }} \
         $(lsb_release -cs) \
         stable"
      ```

      **armhf**:

      ```bash
      $ sudo add-apt-repository \
         "deb [arch=armhf] {{ download-url-base }} \
         $(lsb_release -cs) \
         stable"
      ```

      [Learn about **stable** and **edge** channels](/engine/installation/).
  </div>
  <div id="ee-repo-setup" class="tab-pane fade" markdown="1">
  1.  Install packages to allow `apt` to use a repository over HTTPS:

      ```bash
      $ sudo apt-get install \
          apt-transport-https \
          ca-certificates \
          curl \
          software-properties-common
      ```

  2.  Add Docker's official GPG key using your customer Docker EE repository URL:

      ```bash
      $ curl -fsSL <DOCKER-EE-URL>/ubuntu/gpg | sudo apt-key add -
      ```

      Verify that the key fingerprint is `DD91 1E99 5A64 A202 E859  07D6 BC14 F10B 6D08 5F96`.

      ```bash
      $ apt-key fingerprint 0EBFCD88

      pub   4096R/6D085F96 2017-02-22
          Key fingerprint = DD91 1E99 5A64 A202 E859  07D6 BC14 F10B 6D08 5F96
      uid       [ultimate] Docker Release (EE deb) <docker@docker.com>
      sub   4096R/91A29FA3 2017-02-22
      ```

  3.  Use the following command to set up the **stable** repository, replacing
      `<DOCKER-EE-URL>` with the URL you noted down in the
      [prerequisites](#prerequisites).

      > **Note**: The `lsb_release -cs` sub-command below returns the name of your
      > Ubuntu distribution, such as `xenial`.
      >

      ```bash
      $ sudo add-apt-repository \
         "deb [arch=amd64] <DOCKER-EE-URL>/ubuntu \
         $(lsb_release -cs) \
         stable-{{ minor-version }}"
      ```
  </div>
</div>

#### Install Docker

1.  Update the `apt` package index.

    ```bash
    $ sudo apt-get update
    ```

2.  Install the latest version of Docker, or go to the next step to install a
    specific version. Any existing installation of Docker is replaced.

    Use this command to install the latest version of Docker:

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-group="ce" data-target="#ce-install-docker">Docker CE</a></li>
      <li><a data-toggle="tab" data-group="ee" data-target="#ee-install-docker">Docker EE</a></li>
    </ul>
    <div class="tab-content">
    <div id="ce-install-docker" class="tab-pane fade in active" markdown="1">

    ```bash
    $ sudo apt-get install docker-ce
    ```

    </div>
    <div id="ee-install-docker" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo apt-get install docker-ee
    ```

    </div>
    </div>


    > **Warning**: If you have multiple Docker repositories enabled, installing
    > or updating without specifying a version in the `apt-get install` or
    > `apt-get update` command will always install the newest possible version,
    > which may not be appropriate for your stability needs.
    {: .warning-vanilla}

3.  On production systems, you should install a specific version of Docker
    instead of always using the latest. This output is truncated. List the
    available versions. For Docker EE customers, use `docker-ee` where you see
    `docker-ce`.

    ```bash
    $ apt-cache madison docker-ce

    docker-ce | {{ minor-version }}.0~ce-0~ubuntu-xenial | {{ download-url-base}} xenial/stable amd64 Packages
    ```

    The contents of the list depend upon which repositories are enabled,
    and will be specific to your version of Ubuntu (indicated by the `xenial`
    suffix on the version, in this example). Choose a specific version to
    install. The second column is the version string. The third column is the
    repository name, which indicates which repository the package is from and
    by extension its stability level. To install a specific version, append the
    version string to the package name and separate them by an equals sign (`=`):

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-group="ce" data-target="#ce-install-version-docker">Docker CE</a></li>
      <li><a data-toggle="tab" data-group="ee" data-target="#ee-install-version-docker">Docker EE</a></li>
    </ul>
    <div class="tab-content">
    <div id="ce-install-version-docker" class="tab-pane fade in active" markdown="1">

    ```bash
    $ sudo apt-get install docker-ce=<VERSION>
    ```

    </div>
    <div id="ee-install-version-docker" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo apt-get install docker-ee=<VERSION>
    ```

    </div>
    </div>

    The Docker daemon starts automatically.

4.  Verify that Docker CE or Docker EE is installed correctly by running the
    `hello-world` image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker is installed and running. You need to use `sudo` to run Docker commands.
Continue to [Linux postinstall](linux-postinstall.md) to allow
non-privileged users to run Docker commands and for other optional configuration
steps.

#### Upgrade Docker

To upgrade Docker, first run `sudo apt-get update`, then follow the
[installation instructions](#install-docker), choosing the new version you want
to install.

### Install from a package

If you cannot use Docker's repository to install Docker, you can download the
`.deb` file for your release and install it manually. You will need to download
a new file each time you want to upgrade Docker.

1.  This step is different for Docker CE and Docker EE.

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-group="ce" data-target="#ce-install-from-package-docker">Docker CE</a></li>
      <li><a data-toggle="tab" data-group="ee" data-target="#ee-install-from-package-docker">Docker EE</a></li>
    </ul>
    <div class="tab-content">
    <div id="ce-install-from-package-docker" class="tab-pane fade in active" markdown="1">

    Go to [{{ download-url-base }}/dists/]({{ download-url-base }}/dists/), choose your
    Ubuntu version, browse to `pool/stable/`, choose either `amd64` or
    `armhf`,and download the `.deb` file for the Docker version you want to
    install and for your version of Ubuntu.

      > **Note**: To install an **edge**  package, change the word
      > `stable` in the  URL to `edge`.
      > [Learn about **stable** and **edge** channels](/engine/installation/).

    </div>
    <div id="ee-install-from-package-docker" class="tab-pane fade" markdown="1">

    Go to the Docker EE repository URL associated with your
    trial or subscription in your browser. Go to
    `ubuntu/x86_64/stable-{{ minor-version }}` and download the `.deb` file for the
    Docker version you want to install.

    </div>
    </div>

2.  Install Docker, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo dpkg -i /path/to/package.deb
    ```

    The Docker daemon starts automatically.

3.  Verify that Docker CE or Docker EE is installed correctly by running the
    `hello-world` image.

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
[installation procedure](#install-from-a-package), pointing to the new file.

## Uninstall Docker

1.  Uninstall the Docker package:

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-group="ce" data-target="#ce-uninstall-version-docker">Docker CE</a></li>
      <li><a data-toggle="tab" data-group="ee" data-target="#ee-uninstall-version-docker">Docker EE</a></li>
    </ul>
    <div class="tab-content">
    <div id="ce-uninstall-version-docker" class="tab-pane fade in active" markdown="1">

    ```bash
    $ sudo apt-get purge docker-ce
    ```

    </div>
    <div id="ee-uninstall-version-docker" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo apt-get purge docker-ee
    ```

    </div>
    </div>

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
