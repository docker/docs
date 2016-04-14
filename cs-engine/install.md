<!--[metadata]>
+++
aliases = [ "/docker-trusted-registry/install/engine-ami-launch/",
            "/docker-trusted-registry/install/install-csengine/"]
title = "Install CS Docker Engine"
description = "Learn how to install the commercially supported version of Docker Engine."
keywords = ["docker, engine, dtr, install"]
[menu.main]
parent="menu_csengine"
weight=0
+++
<![end-metadata]-->


# Install CS Docker Engine

This document describes the process of installing the commercially supported
Docker Engine (CS Engine). Installing the CS Engine is a prerequisite for
installing Docker Trusted Registry and/or the Universal Control Plane (UCP).
Follow these instructions if you are installing the CS Engine on physical or
cloud infrastructures.


You first install the CS Engine before you install Docker Trusted Registry.
However, if you are upgrading, you reverse that order and upgrade the Trusted
Registry first. To upgrade, see the [upgrade documentation](upgrade.md).
You will need to install the latest version of the CS Engine to run with the
latest version of the Trusted Registry. You will also want to install the
CS Engine on
any clients, especially in your production environment.

If your cloud provider is AWS, you have the option of installing the CS Engine
using an Amazon Machine Image (AMI). For more information, read
the [installation overview](index.md) to understand your options.

The CS Engine is supported on the following operating systems:


* [CentOS 7.1/7.2 & RHEL 7.0/7.1 (YUM-based systems)](#install-on-centos-7-1-rhel-7-0-7-1-yum-based-systems)
* [Ubuntu 14.04 LTS](#install-on-ubuntu-14-04-lts)
* [SUSE Linux Enterprise 12](#install-on-suse-linux-enterprise-12-3)


## Install CentOS 7.1/7.2 & RHEL 7.0/7.1 (YUM-based systems)

This section explains how to install on CentOS 7.1/7.2 & RHEL 7.0/7.1. Only
these versions are supported. CentOS 7.0 is **not** supported. On RHEL,
depending on your current level of updates, you may need to reboot your server
to update its RHEL kernel.

1. Log into the system as a user with root or sudo permissions.

2. Add Docker's public key for CS packages:

    ```bash
    $ sudo rpm --import "https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e"
    ```

3. Install yum-utils if necessary:

    ```bash
    $ sudo yum install -y yum-utils
    ```

4. Add the repository. Notice in the following code that it gets the latest
version of the CS Engine. Each time you either install or upgrade, ensure that
you are requesting the version and the OS that you want.


    ```bash
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.10/yum/repo/main/centos/7
    ```

5. Install the CS Engine with the following command:

    ```bash
    $ sudo yum install docker-engine
    ```

6. Enable the Docker daemon as a service and then start it.

    ```bash
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

7. Verify the installation was successful by running a simple container.

        $ sudo docker run hello-world
        Unable to find image 'hello-world:latest' locally
        latest: Pulling from hello-world
        b901d36b6f2fd75: Pull complete
        0a6ba66e537a53a: Pull complete
        hello-world:latest: The image you are pulling has been verified. Important: image verification is a tech preview feature and should not be relied on to provide security.
        Digest: sha256:517f03be3f8169d84711c9ffb2b3235a4d27c1eb4ad147f6248c8040adb93113
        Status: Downloaded newer image for hello-world:latest

        Hello from Docker.
        This message shows that your installation appears to be working correctly.

        To generate this message, Docker took the following steps:
         1. The Docker client contacted the Docker daemon.
         2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
         3. The Docker daemon created a new container from that image which runs the
            executable that produces the output you are currently reading.
         4. The Docker daemon streamed that output to the Docker client, which sent it
            to your terminal.

        To try something more ambitious, you can run an Ubuntu container with:
         $ docker run -it ubuntu bash


8. (Optional) Add non-sudo access to the Docker socket by adding your user
to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

9. Log out and log back in to have your new permissions take effect.


## Install on Ubuntu 14.04 LTS

1. Log into the system as a user with root or sudo permissions.

2. Add Docker's public key for CS packages:

    ```bash
    $ curl -s 'https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e' | sudo apt-key add --import
    ```

3. Install the HTTPS helper for apt (your system may already have it):

    ```bash
    $ sudo apt-get update && sudo apt-get install apt-transport-https
    ```

4. Install additional virtual drivers not in the base image.

    ```bash
    $ sudo apt-get install -y linux-image-extra-virtual
    ```

    You may need to reboot your server after updating the LTS kernel.

5. Add the repository for the new version:

    ```bash
    $ echo "deb https://packages.docker.com/1.10/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
    ```

    You must modify the "ubuntu-trusty" string for your flavor of ubuntu or
    debian as seen in the following options:

    * debian-jessie (Debian 8)
    * debian-stretch (future release)
    * debian-wheezy (Debian 7)
    * ubuntu-precise (Ubuntu 12.04)
    * ubuntu-trusty (Ubuntu 14.04)
    * ubuntu-utopic (Ubuntu 14.10)
    * ubuntu-vivid (Ubuntu 15.04)
    * ubuntu-wily (Ubuntu 15.10)

6. Run the following to install commercially supported Docker Engine and its
dependencies:

    ```bash
    $ sudo apt-get update && sudo apt-get install docker-engine
    ```

7. Confirm the Docker daemon is running:

    ```bash
    $ sudo docker info
    ```

8. Optionally, add non-sudo access to the Docker socket by adding your
user to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

    Log out and log back in to have your new permissions take effect.


## Install on SUSE Linux Enterprise 12.3

1. Log into the system as a user with root or sudo permissions.

2. Refresh your repository so that curl commands and CA certificates
are available.

    ```bash
    $ sudo zypper ref
    ```

3. Add the repository and the signing key. Notice in the following code
that it gets the latest version of the CS Engine. Each time you either
install or upgrade, ensure that the you are requesting the version and the
OS that you want.

    ```bash
    $ sudo zypper ar -t YUM https://packages.docker.com/1.10/yum/repo/main/opensuse/12.3 docker-1.10
    $ sudo rpm --import 'https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e'
    ```

4. Install the Docker daemon package:

    ```bash
    $ sudo zypper install docker-engine
    ```

5. Enable the Docker daemon as a service and then start it:

    ```bash
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

6. Confirm the Docker daemon is running:

    ```bash
    $ sudo docker info
    ```

7. Optionally, add non-sudo access to the Docker socket by adding your user
to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

8. Log out and log back in to have your new permissions take effect.
