---
description: Learn how to install the commercially supported version of Docker Engine.
keywords: docker, engine, dtr, install
redirect_from:
- /docker-trusted-registry/install/engine-ami-launch/
- /docker-trusted-registry/install/install-csengine/
- /docker-trusted-registry/cs-engine/install/
title: Install Commercially Supported Docker Engine
---

Follow these instructions to install CS Docker Engine, the commercially
supported version of Docker Engine.

CS Docker Engine can be installed on the following operating systems:


* [CentOS 7.1/7.2 & RHEL 7.0/7.1/7.2 (YUM-based systems)](install.md#install-on-centos-7172--rhel-707172-yum-based-systems)
* [Ubuntu 14.04 LTS](install.md#install-on-ubuntu-1404-lts)
* [SUSE Linux Enterprise 12](install.md#install-on-suse-linux-enterprise-123)


## Install on CentOS 7.1/7.2 & RHEL 7.0/7.1/7.2 (YUM-based systems)

This section explains how to install on CentOS 7.1/7.2 & RHEL 7.0/7.1/7.2. Only
these versions are supported. CentOS 7.0 is **not** supported. On RHEL,
depending on your current level of updates, you may need to reboot your server
to update its RHEL kernel.

1. Log into the system as a user with root or sudo permissions.

2.  Add the Docker public key for CS packages:

    ```bash
    $ sudo rpm --import "https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e"
    ```

3.  Install yum-utils if necessary:

    ```bash
    $ sudo yum install -y yum-utils
    ```

4.  Add the Docker repository:

    ```bash
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.12/yum/repo/main/centos/7
    ```

    This adds the repository of the latest version of CS Docker Engine. You can
    customize the URL to install an older version.

5.  Install Docker CS Engine:

    ```bash
    $ sudo yum install docker-engine
    ```

6.  Configure devicemapper:

    By default, the `devicemapper` graph driver does not come pre-configured in a production ready state. Follow the documented step by step instructions to [configure devicemapper with direct-lvm for production](../../engine/userguide/storagedriver/device-mapper-driver/#/for-a-direct-lvm-mode-configuration) in order to achieve the best performance and reliability for your environment.

7.  Enable the Docker daemon as a service and start it.

    ```bash
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

8.  Confirm the Docker daemon is running:

    ```bash
    $ sudo docker info
    ```

9.  Optionally, add non-sudo access to the Docker socket by adding your user
to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

10. Log out and log back in to have your new permissions take effect.

## Install on Ubuntu 14.04 LTS

1. Log into the system as a user with root or sudo permissions.

2.  Add Docker's public key for CS packages:

    ```bash
    $ curl -s 'https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e' | sudo apt-key add --import
    ```

3.  Install the HTTPS helper for apt (your system may already have it):

    ```bash
    $ sudo apt-get update && sudo apt-get install apt-transport-https
    ```

4.  Install additional kernel modules to add AUFS support.

    ```bash
    $ sudo apt-get install -y linux-image-extra-$(uname -r) linux-image-extra-virtual
    ```

5.  Add the repository for the new version:

    ```bash
    $ echo "deb https://packages.docker.com/1.12/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
    ```

6.  Run the following to install commercially supported Docker Engine and its
dependencies:

    ```bash
    $ sudo apt-get update && sudo apt-get install docker-engine
    ```

7. Confirm the Docker daemon is running:

    ```bash
    $ sudo docker info
    ```

8.  Optionally, add non-sudo access to the Docker socket by adding your
user to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

    Log out and log back in to have your new permissions take effect.


## Install on SUSE Linux Enterprise 12.3

1. Log into the system as a user with root or sudo permissions.

2.  Refresh your repository so that curl commands and CA certificates
are available:

    ```bash
    $ sudo zypper ref
    ```

3.  Add the Docker repository and public key:

    ```bash
    $ sudo zypper ar -t YUM https://packages.docker.com/1.12/yum/repo/main/opensuse/12.3 docker-1.12
    $ sudo rpm --import 'https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e'
    ```

    This adds the repository of the latest version of CS Docker Engine. You can
    customize the URL to install an older version.

4.  Install the Docker daemon package:

    ```bash
    $ sudo zypper install docker-engine
    ```

5.  Enable the Docker daemon as a service and then start it:

    ```bash
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

6.  Confirm the Docker daemon is running:

    ```bash
    $ sudo docker info
    ```

7.  Optionally, add non-sudo access to the Docker socket by adding your user
to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

8. Log out and log back in to have your new permissions take effect.
