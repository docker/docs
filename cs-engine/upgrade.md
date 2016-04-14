<!--[metadata]>
+++
title = "Upgrade"
description = "Learn how to install the commercially supported version of Docker Engine."
keywords = ["docker, engine, dtr, upgrade"]
[menu.main]
parent="menu_csengine"
identifier="cse_upgrade"
weight=10
+++
<![end-metadata]-->

# Upgrade CS Docker Engine

This article explains how to upgrade your CS Docker Engine installation.
Before starting the upgrade, make sure all containers running on that CS
Docker Engine are stopped.

The way you upgrade your CS Docker Engine, depends on the version that is
currently installed and the version that you want to upgrade to:

* [Upgrade from the same minor version](#upgrade-from-the-same-minor-version), if
you're upgrading from 1.10 to 1.10.x,
* [Upgrade from the same major version](#upgrade-from-the-same-major-version), if
you're upgrading from 1.9 to 1.10,
* [Upgrade from a legacy version](#upgrade-from-a-legacy-version), if you're
upgrading from a version prior to 1.9.


## Upgrade from the same minor version

Use these instructions if you're upgrading your CS Docker Engine within the
same minor version. As an example, from 1.10.0 to 1.10.1.

### CentOS 7.1 & RHEL 7.0/7.1
Use these instructions to upgrade YUM-based systems.

1. Update your docker-engine package:

    ```bash
    $ sudo yum upgrade docker-engine
    ```

2. Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```

### Ubuntu 14.04 LTS
Use these instructions to upgrade APT-based systems.

1. Update your docker-engine package:

    ```bash
    $ sudo apt-get update && sudo apt-get upgrade docker-engine
    ```

2. Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```

### SUSE Enterprise 12.3

1. Update your docker-engine package:

    ```bash
    $ sudo zypper upgrade docker-engine
    ```

2. Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```


## Upgrade from the same major version

Use these instructions if you're upgrading your CS Docker Engine within the
same major version. As an example, from 1.9 to 1.10.


### CentOS 7.1 & RHEL 7.0/7.1
Use these instructions to upgrade YUM-based systems.

1. Add the docker engine repository.

    ```bash
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.10/yum/repo/main/centos/7
    ```

    In this example we are adding the docker engine 1.10 repository. But sure
    to add the repository of the version you want to install.

2. Install the new package:

    ```bash
    $ sudo yum update docker-engine
    ```

3. Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```

### Ubuntu 14.04 LTS
Use these instructions to update APT-based systems.


1. Add the docker engine repository.

    ```bash
    $ echo "deb https://packages.docker.com/1.10/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
    ```

    In this example we are adding the docker engine ubuntu-trusty 1.10
    repository. Be sure to use the Docker Engine version you want to upgrade to,
    and your linux distribution. The following distributions are available:

    * debian-jessie (Debian 8)
    * debian-stretch (future release)
    * debian-wheezy (Debian 7)
    * ubuntu-precise (Ubuntu 12.04)
    * ubuntu-trusty (Ubuntu 14.04)
    * ubuntu-utopic (Ubuntu 14.10)
    * ubuntu-vivid (Ubuntu 15.04)
    * ubuntu-wily (Ubuntu 15.10)

2. Update your docker-engine package.

    ```bash
    $ sudo apt-get update && sudo apt-get upgrade docker-engine
    ```

3. Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```

#### SUSE Enterprise 12.3

1. Add the docker engine repository.

      ```bash
      $ sudo zypper ar -t YUM https://packages.docker.com/1.10/yum/repo/main/opensuse/12.3 docker-1.10
      ```

      In this example we are adding the CS Docker Engine 1.0 repository.
      Be sure to use the Docker Engine version you want to upgrade to.

2. Install the new package:

    ```bash
    $ sudo zypper update docker-engine
    ```

3. Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```


## Upgrade from a legacy version

Use these instructions if you're upgrading your CS Docker Engine from a version
prior to 1.9. In this case you'll have to first uninstall CS Docker Engine, and
then install the latest version.

### CentOS 7.1 & RHEL 7.0/7.1
Use these instructions to upgrade YUM-based systems.

1. Remove the current CS Engine:

    ```bash
    $ sudo yum remove docker-engine-cs
    ```

2. Add Docker's public key for CS packages:

    ```bash
    $ sudo rpm --import "https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e"
    ```

3. Install yum-utils if necessary:

    ```bash
    $ sudo yum install -y yum-utils
    ```

4. Add the repository for the new version and disable the old one.

    ```bash
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.10/yum/repo/main/centos/7
    $ sudo yum-config-manager --disable 'Docker_cs*'
    ```

    In this example we are adding the CS Docker Engine 1.0 repository.
    Be sure to use the Docker Engine version you want to upgrade to.

5. Install the new package:

    ```bash
    $ sudo yum install docker-engine
    ```

6. Enable the Docker daemon as a service and start it.

    ```
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

### Ubuntu 14.04 LTS
Use these instructions to update APT-based systems.


1. Remove the current Engine:

    ```bash
    $ sudo apt-get remove docker-engine-cs
    ```

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

    In this example we are adding the docker engine ubuntu-trusty 1.10
    repository. Be sure to use the Docker Engine version you want to upgrade to,
    and your linux distribution. The following distributions are available:

    * debian-jessie (Debian 8)
    * debian-stretch (future release)
    * debian-wheezy (Debian 7)
    * ubuntu-precise (Ubuntu 12.04)
    * ubuntu-trusty (Ubuntu 14.04)
    * ubuntu-utopic (Ubuntu 14.10)
    * ubuntu-vivid (Ubuntu 15.04)
    * ubuntu-wily (Ubuntu 15.10)


6. Install the upgraded package:

    ```bash
    $ sudo apt-get upgrade docker-engine
    ```
