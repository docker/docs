+++
title = "Manually Install the CS Docker Engine"
description = "Install instructions for the commercially supported Docker Engine"
keywords = ["docker, documentation, about, technology, enterprise, hub, commercially supported Docker Engine, CS engine, registry"]
[menu.main]
parent="smn_dhe_install"
+++


# Manually Install the CS Docker Engine

This document describes the process of installing the Commercially Supported
Docker Engine (CS Engine). Installing the CS Engine is a prerequisite for
installing the Docker Trusted Registry. Use these instructions if you
are installing the CS Engine on physical or cloud infrastructures.

Note that you first install the CS Engine before you install Docker Trusted
Registry. If you are upgrading, you reverse that order and upgrade the Trusted
Registry first. To upgrade, see the [upgrade documentation](upgrade.md). You will need to install the latest version of the CS Engine to run with the latest
version of the Trusted Registry. You will also want to install the CS Engine on
any clients, especially in your production environment.

If your cloud provider is AWS, you have the option of installing the CS Engine
using an Amazon Machine Image (AMI). For more information, read the [installation overview](index.md) to understand your options.

## Prerequisites

You need a login to Docker Hub. If you have not already done so, go to Docker Hub and [sign up for an account](https://hub.docker.com). You do not need a license for the CS Engine, only for the Docker Trusted Registry.

## CentOS 7.1 & RHEL 7.0/7.1 (YUM-based systems)

This section explains how to install on CentOS 7.1 & RHEL 7.0/7.1. Only these
versions are supported. CentOS 7.0 is not supported. On RHEL, depending on your
current level of updates, you may need to reboot your server to update its RHEL
kernel.

1. Log into the system as a user with root or sudo permissions.

2. Add Docker's public key for CS packages:

    `$ sudo rpm --import "https://pgp.mit.edu/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e"`

3. Install yum-utils if necessary:

    `$ sudo yum install -y yum-utils`

4. Add the repository:

    ```
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.9/yum/repo/main/centos/7
    ```

5. Install the CS Engine with the following command:

        $ sudo yum install docker-engine

6. Enable the Docker daemon as a service and then start it.

        $ sudo systemctl enable docker.service
        $ sudo systemctl start docker.service

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

8. Optionally, add non-sudo access to the Docker socket by adding your user to the `docker` group.

        $ sudo usermod -a -G docker $USER

9. Log out and log back in to have your new permissions take effect.


## Install on Ubuntu 14.04 LTS

1. Log into the system as a user with root or sudo permissions.

2. Add Docker's public key for CS packages:

    `$ curl -s 'https://pgp.mit.edu/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e' | sudo apt-key add --import`

3. Install the HTTPS helper for apt (your system may already have it):

    `$ sudo apt-get update && sudo apt-get install apt-transport-https`

4. Install additional virtual drivers not in the base image.

        $ sudo apt-get install -y linux-image-extra-virtual

      You may need to reboot your server after updating the LTS kernel.

5. Add the repository for the new version:

    `$ echo "deb https://packages.docker.com/1.9/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list`

        **Note**: modify the "ubuntu-trusty" string for your flavor of ubuntu or debian.
        * debian-jessie (Debian 8)
        * debian-stretch (future release)
        * debian-wheezy (Debian 7)
        * ubuntu-precise (Ubuntu 12.04)
        * ubuntu-trusty (Ubuntu 14.04)
        * ubuntu-utopic (Ubuntu 14.10)
        * ubuntu-vivid (Ubuntu 15.04)
        * ubuntu-wily (Ubuntu 15.10)

6. Run the following to install commercially supported Docker Engine and its dependencies:

    `$ sudo apt-get update && sudo apt-get install docker-engine`

7. Confirm the Docker daemon is running with `sudo service docker start`.

        $ sudo service docker start

8. Optionally, add non-sudo access to the Docker socket by adding your user to the `docker` group.

        $ sudo usermod -a -G docker $USER

    Log out and log back in to have your new permissions take effect.


## Next step
You are ready to install [Docker Trusted Registry](install-dtr.md).
