---
title: Upgrade Commercially Supported Docker Engine
description: Learn how to upgrade the commercially supported version of Docker Engine.
keywords: docker, engine, upgrade
---

This article explains how to upgrade your CS Docker Engine.

The upgrade process depends on the version that is currently installed and the
version that you want to upgrade to:

* [Upgrade from the same minor version](upgrade.md#upgrade-from-the-same-minor-version)
* [Upgrade from the same major version](upgrade.md#upgrade-from-the-same-major-version)
* [Upgrade from a legacy version](upgrade.md#upgrade-from-a-legacy-version)

Before starting the upgrade, make sure you stop all containers running on the
host. This ensures your containers have time for cleaning up before exiting,
thus avoiding data loss or corruption.

## Upgrade from the same minor version

Use these instructions if you're upgrading your CS Docker Engine within the
same minor version. As an example, from 1.1.0 to 1.1.1.

### CentOS 7.1 & RHEL 7.0/7.1/7.2
Use these instructions to upgrade YUM-based systems.

1.  Update your docker-engine package:

    ```bash
    $ sudo yum upgrade docker-engine
    ```

2.  Check that the CS Docker Engine is running:

    ```bash
    $ sudo docker info
    ```

### Ubuntu 14.04 LTS or 16.04 LTS
Use these instructions to upgrade APT-based systems.

1.  Update your docker-engine package:

    ```bash
    $ sudo apt-get update && sudo apt-get upgrade docker-engine
    ```

2.  Check that the CS Docker Engine is running:

    ```bash
    $ sudo docker info
    ```

### SUSE Enterprise 12.3

1.  Update your docker-engine package:

    ```bash
    $ sudo zypper upgrade docker-engine
    ```

2.  Check that the CS Docker Engine is running:

    ```bash
    $ sudo docker info
    ```


## Upgrade from the same major version

Use these instructions if you're upgrading your CS Docker Engine within the
same major version. As an example, from 1.1.x to 1.2.x.


### CentOS 7.1 & RHEL 7.0/7.1/7.2
Use these instructions to upgrade YUM-based systems.

1.  Add the Docker Engine repository.

    ```bash
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.13/yum/repo/main/centos/7
    ```

    This adds the repository of the latest version of CS Docker Engine. You can
    customize the URL to install other versions.

2.  Install the new package:

    ```bash
    $ sudo yum update docker-engine
    ```

3.  Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```

### Ubuntu 14.04 LTS or 16.04 LTS
Use these instructions to update APT-based systems.


1.  Add the docker engine repository.

    ```bash
    $ echo "deb https://packages.docker.com/1.13/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
    ```

    This adds the repository of the latest version of CS Docker Engine for the
    Ubuntu Trusty distribution. Change the "ubuntu-trusty" string to the
    distribution you're using:

    * debian-jessie (Debian 8)
    * debian-stretch (future release)
    * debian-wheezy (Debian 7)
    * ubuntu-precise (Ubuntu 12.04)
    * ubuntu-trusty (Ubuntu 14.04)
    * ubuntu-xenial (Ubuntu 16.04)

2.  Update your docker-engine package.

    ```bash
    $ sudo apt-get update && sudo apt-get upgrade docker-engine
    ```

3.  Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```

#### SUSE Enterprise 12.3

1.  Add the docker engine repository.

    ```bash
    $ sudo zypper ar -t YUM https://packages.docker.com/1.13/yum/repo/main/opensuse/12.3 docker-1.13
    ```

    This adds the repository of the latest version of CS Docker Engine. You
    can customize the URL to install other versions.

2.  Install the new package:

    ```bash
    $ sudo zypper update docker-engine
    ```

3.  Check that the CS Engine is running:

    ```bash
    $ sudo docker info
    ```


## Upgrade from a legacy version

Use these instructions if you're upgrading your CS Docker Engine from a version
prior to 1.9. In this case, first uninstall CS Docker Engine, and
then install the latest version.

### CentOS 7.1 & RHEL 7.0/7.1
Use these instructions to upgrade YUM-based systems.

1.  Remove the current CS Engine:

    ```bash
    $ sudo yum remove docker-engine-cs
    ```

2.  Add the Docker public key for CS packages:

    ```bash
    $ sudo rpm --import "https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e"
    ```

    Note: if the key server above does not respond, you can try one of these:
       - pgp.mit.edu
       - keyserver.ubuntu.com

3.  Install yum-utils if necessary:

    ```bash
    $ sudo yum install -y yum-utils
    ```

4.  Add the repository for the new version and disable the old one.

    ```bash
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.13/yum/repo/main/centos/7
    $ sudo yum-config-manager --disable 'Docker_cs*'
    ```

    This adds the repository of the latest version of CS Docker Engine. You
    can customize the URL to install other versions.

5.  Install the new package:

    ```bash
    $ sudo yum install docker-engine
    ```

6.  Enable the Docker daemon as a service and start it.

    ```bash
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

### Ubuntu 14.04 LTS
Use these instructions to update APT-based systems.


1.  Remove the current Engine:

    ```bash
    $ sudo apt-get remove docker-engine-cs
    ```

2.  Add the Docker public key for CS packages:

    ```bash
    $ curl -s 'https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e' | sudo apt-key add --import
    ```

    Note: if the key server above does not respond, you can try one of these:
       - pgp.mit.edu
       - keyserver.ubuntu.com

3.  Install the HTTPS helper for apt (your system may already have it):

    ```bash
    $ sudo apt-get update && sudo apt-get install apt-transport-https
    ```

4.  Install additional virtual drivers not in the parent image.

    ```bash
    $ sudo apt-get install -y linux-image-extra-virtual
    ```

    You may need to reboot your server after updating the LTS kernel.

5.  Add the repository for the new version:

    ```bash
    $ echo "deb https://packages.docker.com/1.13/apt/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
    ```

    This adds the repository of the latest version of CS Docker Engine for the
    Ubuntu Trusty distribution. Change the "ubuntu-trusty" string to the
    distribution you're using:

    * debian-jessie (Debian 8)
    * debian-stretch (future release)
    * debian-wheezy (Debian 7)
    * ubuntu-precise (Ubuntu 12.04)
    * ubuntu-trusty (Ubuntu 14.04)
    * ubuntu-xenial (Ubuntu 16.04)



6.  Install the upgraded package:

    ```bash
    $ sudo apt-get upgrade docker-engine
    ```
