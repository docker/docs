+++
title = "Manually Install the CS Docker Engine"
description = "Install instructions for the commercially supported Docker Engine"
keywords = ["docker, documentation, about, technology, enterprise, hub, commercially supported Docker Engine, CS engine, registry"]
[menu.main]
parent="smn_dhe_install"
+++


# Manually Install the CS Docker Engine

This document describes the process of obtaining and installing the Commercially
Supported Docker Engine (CS Engine). Installing CS Engine is a prerequisite for
installing the Docker Trusted Registry. You can use these instructions if you are installing CS Engine on physical or cloud infrastructure.

If your cloud provider is AWS, you have the option of installing CS Engine using an Amazon Machine Image (AMI) instead. For more information, read the [installation overview](index.md) to understand your options.

## Prerequisite

Installing CS Engine requires that you have a login to Docker Hub. If you have
not already done so, go to Docker Hub and [sign up for an
account](https://hub.docker.com).

Also, you must have a license for Docker Trusted Registry. This license allows
you to run both Docker Trusted Registry and CS Engine. Before installing,
[purchase a license or sign up for a free, 30 day trial license]((https://hub.docker.com/enterprise/)).


## Install on CentOS 7.1 & RHEL 7.0/7.1

This section explains how to install on CentOS 7.1 & RHEL 7.0/7.1. Only these versions are supported. CentOS 7.0 is not supported. On RHEL, depending on your current level of updates, you may need to reboot your server to update its RHEL kernel.

1. Log into the system as a user with root or sudo permissions.

2. Update your `yum` repositories.

        $ sudo yum update && sudo yum upgrade

3. In a browser, log in to the [Docker Hub](https://hub.docker.com) with the account you used to obtain your license.

4. Once you're logged in, go to your account's [Licenses](https://hub.docker.com/account/licenses/) page.

5. In the "Download and Install CS Engine" locate the script appropriate to your system.

6. Copy the script, paste it into your terminal, and press Return.

        $ curl -s
        https://packagecloud.io/install/repositories/Docker/cs-public/script.rpm.sh |
        sudo bash sudo yum install docker-engine-cs

7. After the command completes, install the CS Engine with the following command:

        $ sudo yum install docker-engine-cs

8. Enable the Docker daemon as a service and then start it.

        $ sudo systemctl enable docker.service
        $ sudo systemctl start docker.service

9. Verify the installation was successful by running a simple container.

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

        Share images, automate workflows, and more with a free Docker Hub account:
         https://hub.docker.com

        For more examples and ideas, visit:
         https://docs.docker.com/userguide/

10. Optionally, add non-sudo access to the Docker socket by adding your user to the `docker` group.

        $ sudo usermod -a -G docker $USER

  Log out and log back in to have your new permissions take effect.



## Install on Ubuntu 14.04 LTS

1. Log into the system as a user with root or sudo permissions.

2. Update your `yum` repositories.

        $ sudo apt-get update && sudo apt-get upgrade

3. Install additional virtual drivers not in the base image.

        $ sudo apt-get install -y linux-image-extra-virtual

    You may need to reboot your server to after updating the LTS kernel.

4. In a browser, log in to the [Docker Hub](https://hub.docker.com) with the account you used to obtain your license.

5. Once you're logged in, go to your account's [Licenses](https://hub.docker.com/account/licenses/) page.

6. In the "Download and Install CS Engine" locate the script appropriate to your system.

7. Copy the script, paste it into your terminal, and press Return.

        $ curl -s
        https://packagecloud.io/install/repositories/Docker/cs-public/script.deb.sh | sudo bash sudo apt-get install docker-engine-cs

8. Run the following to install commercially supported Docker Engine and its dependencies:
​
        $ sudo apt-get install docker-engine-cs

9. Confirm the Docker daemon is running with `sudo service docker start`.

        $ sudo service docker start

10. Optionally, add non-sudo access to the Docker socket by adding your user to the `docker` group.

        $ sudo usermod -a -G docker $USER

    Log out and log back in to have your new permissions take effect.


## Next step
You are ready to install [Docker Trusted Registry](install-dtry.md).
