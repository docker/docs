---
description: Learn how to install the commercially supported version of Docker Engine.
keywords: docker, engine, dtr, install
title: Install CS Docker Engine
redirect_from:
- /cs-engine/1.12/install/
---

Follow these instructions to install CS Docker Engine, the commercially
supported version of Docker Engine.

CS Docker Engine can be installed on the following operating systems:

* [CentOS 7.1/7.2 & RHEL 7.0/7.1/7.2 (YUM-based systems)](#install-on-centos-7172--rhel-707172-yum-based-systems)
* [Ubuntu 14.04 LTS](#install-on-ubuntu-1404-lts)
* [SUSE Linux Enterprise 12](#install-on-suse-linux-enterprise-123)

You can install CS Docker Engine using a repository or using packages.

- If you [use a repository](#install-using-a-repository), your operating system
  will notify you when updates are available and you can upgrade or downgrade
  easily, but you need an internet connection. This approach is recommended.

- If you [use packages](#install-using-packages), you can install CS Docker
  Engine on air-gapped systems that have no internet connection. However, you
  are responsible for manually checking for updates and managing upgrades.

## Prerequisites

To install CS Docker Engine, you need root or sudo privileges and you need
access to a command line on the system.

## Install using a repository

### Install on CentOS 7.1/7.2 & RHEL 7.0/7.1/7.2/7.3 (YUM-based systems)

This section explains how to install on CentOS 7.1/7.2 & RHEL 7.0/7.1/7.2/7.3. Only
these versions are supported. CentOS 7.0 is **not** supported. On RHEL,
depending on your current level of updates, you may need to reboot your server
to update its RHEL kernel.

1.  Add the Docker public key for CS Docker Engine packages:

    ```bash
    $ sudo rpm --import "https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e"
    ```

    > **Note**: If the key server above does not respond, you can try one of these:
    >
    >  - pgp.mit.edu
    >  - keyserver.ubuntu.com

2.  Install yum-utils if necessary:

    ```bash
    $ sudo yum install -y yum-utils
    ```

3.  Add the Docker repository:

    ```bash
    $ sudo yum-config-manager --add-repo https://packages.docker.com/1.12/yum/repo/main/centos/7
    ```

    This adds the repository of the latest version of CS Docker Engine. You can
    customize the URL to install an older version.

4.  Install Docker CS Engine:

    - **Latest version**:

      ```bash
      $ sudo yum makecache fast

      $ sudo yum install docker-engine
      ```

    - **Specific version**:

      On production systems, you should install a specific version rather than
      relying on the latest.

      1.  List the available versions:

          ```bash
          $ yum list docker-engine.x86_64  --showduplicates |sort -r
          ```

          The second column represents the version.

      2.  Install a specific version by adding the version after `docker-engine`,
          separated by a hyphen (`-`):

          ```bash
          $ sudo yum install docker-engine-<version>
          ```

5.  Configure `devicemapper`:

    By default, the `devicemapper` graph driver does not come pre-configured in
    a production-ready state. Follow the documented step by step instructions to
    [configure devicemapper with direct-lvm for production](../../engine/userguide/storagedriver/device-mapper-driver/#configure-direct-lvm-mode-for-production)
    to achieve the best performance and reliability for your environment.

6.  Configure the Docker daemon to start automatically when the system starts,
    and start it now.

    ```bash
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

7.  Confirm the Docker daemon is running:

    ```bash
    $ sudo docker info
    ```

8.  Only users with `sudo` access can run `docker` commands.
    Optionally, add non-sudo access to the Docker socket by adding your user
    to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

9. Log out and log back in to have your new permissions take effect.


### Install on Ubuntu 14.04 LTS or 16.04 LTS

1.  Install packages to allow `apt` to use a repository over HTTPS:

    ```bash
    $ sudo apt-get update

    $ sudo apt-get install --no-install-recommends \
        apt-transport-https \
        curl \
        software-properties-common
    ```

    Optionally, install additional kernel modules to add AUFS support.

    ```bash
    $ sudo apt-get install -y --no-install-recommends \
        linux-image-extra-$(uname -r) \
        linux-image-extra-virtual
    ```

2.  Download and import Docker's public key for CS packages:

    ```bash
    $ curl -fsSL 'https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e' | sudo apt-key add -
    ```

    >**Note**: If the key server above does not respond, you can try one of these:
    >
    >   - pgp.mit.edu
    >   - keyserver.ubuntu.com

3.  Add the repository. In the  command below, the `lsb_release -cs` sub-command
    returns the name of your Ubuntu version, like `xenial` or `trusty`.

    ```bash
    $ sudo add-apt-repository \
       "deb https://packages.docker.com/1.12/apt/repo/ \
       ubuntu-$(lsb_release -cs) \
       main"
    ```

4.  Install CS Docker Engine:

    - **Latest version**:

      ```bash
      $ sudo apt-get update

      $ sudo apt-get -y install docker-engine
      ```

    - **Specific version**:

      On production systems, you should install a specific version rather than
      relying on the latest.

      1.  List the available versions:

          ```bash
          $ sudo apt-get update

          $ apt-cache madison docker-engine
          ```

          The second column represents the version.

      2.  Install a specific version by adding the version after `docker-engine`,
          separated by an equals sign (`=`):

          ```bash
          $ sudo apt-get install docker-engine=<version>
          ```

5.  Confirm the Docker daemon is running:

    ```bash
    $ sudo docker info
    ```

6.  Only users with `sudo` access can run `docker` commands.
    Optionally, add non-sudo access to the Docker socket by adding your user
    to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

    Log out and log back in to have your new permissions take effect.

### Install on SUSE Linux Enterprise 12.3

1.  Refresh your repository:

    ```bash
    $ sudo zypper update
    ```

2.  Add the Docker repository and public key:

    ```bash
    $ sudo zypper ar -t YUM https://packages.docker.com/1.12/yum/repo/main/opensuse/12.3 docker-1.13
    $ sudo rpm --import 'https://sks-keyservers.net/pks/lookup?op=get&search=0xee6d536cf7dc86e2d7d56f59a178ac6c6238f52e'
    ```

    This adds the repository of the latest version of CS Docker Engine. You can
    customize the URL to install an older version.

   **Note**: If the key server above does not respond, you can try one of these:
    >
    >   - pgp.mit.edu
    >   - keyserver.ubuntu.com

3.  Install CS Docker Engine.

    - **Latest version**:

      ```bash
      $ sudo zypper refresh

      $ sudo zypper install docker-engine
      ```

    - **Specific version**:

      On production systems, you should install a specific version rather than
      relying on the latest.

      1.  List the available versions:

          ```bash
          $ sudo zypper refresh

          $ zypper search -s --match-exact -t package docker-engine
          ```

          The third column is the version string.

      2.  Install a specific version by adding the version after `docker-engine`,
          separated by a hyphen (`-`):

          ```bash
          $ sudo zypper install docker-engine-<version>
          ```

4.  Configure the Docker daemon to start automatically when the system starts,
    and start it now.

    ```bash
    $ sudo systemctl enable docker.service
    $ sudo systemctl start docker.service
    ```

5.  Confirm the Docker daemon is running:

    ```bash
    $ sudo docker info
    ```

6.  Only users with `sudo` access can run `docker` commands.
    Optionally, add non-sudo access to the Docker socket by adding your user
    to the `docker` group.

    ```bash
    $ sudo usermod -a -G docker $USER
    ```

    Log out and log back in to have your new permissions take effect.

7.  [Configure Btrfs for graph storage](/engine/userguide/storagedriver/btrfs-driver.md).
    This is the only graph storage driver supported on SLES.

## Install using packages

If you need to install Docker on an air-gapped system with no access to the
internet, use the [package download link table](#package-download-links) to
download the Docker package for your operating system, then install it using the
[appropriate command](#general-commands). You are responsible for manually
upgrading Docker when a new version is available, and also for satisfying
Docker's dependencies.

### General commands

To install Docker from packages, use the following commands:

| Operating system      | Command |
|-----------------------|---------|
| RHEL / CentOS / SLES  | `$ sudo yum install /path/to/package.rpm` |
| Ubuntu                | `$ sudo dpkg -i /path/to/package.deb`     |

### Package download links

{% assign rpm-prefix = "https://packages.docker.com/1.12/yum/repo/main" %}
{% assign deb-prefix = "https://packages.docker.com/1.12/apt/repo/pool/main/d/docker-engine" %}

#### CS Docker Engine 1.12.6

{% comment %} Check on the S3 bucket for packages.docker.com for the versions. {% endcomment %}
{% assign rpm-version = "1.12.6.cs8-1" %}
{% assign rpm-rhel-version = "1.12.6.cs8-1" %}
{% assign deb-version = "1.12.6~cs8-0" %}

| Operating system      | Package links                                                                                                                                                                                                                                                                                                                                               |
|-----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RHEL 7.x and CentOS 7 | [docker-engine]({{ rpm-prefix }}/centos/7/Packages/docker-engine-{{ rpm-version}}.el7.centos.x86_64.rpm), [docker-engine-debuginfo]({{ rpm-prefix }}/centos/7/Packages/docker-engine-debuginfo-{{ rpm-version }}.el7.centos.x86_64.rpm), [docker-engine-selinux]({{ rpm-prefix }}/centos/7/Packages/docker-engine-selinux-{{ rpm-version}}1.el7.centos.noarch.rpm) |
| RHEL 7.2 (only use if you have problems with `selinux` with the packages above) | [docker-engine]({{ rpm-prefix }}/rhel/7.2/Packages/docker-engine-{{ rpm-rhel-version }}.el7.centos.x86_64.rpm), [docker-engine-debuginfo]({{ rpm-prefix }}/rhel/7.2/Packages/docker-engine-debuginfo-{{ rpm-rhel-version }}.el7.centos.x86_64.rpm), [docker-engine-selinux]({{ rpm-prefix }}/rhel/7.2/Packages/docker-engine-selinux-{{ rpm-rhel-version }}.el7.centos.noarch.rpm) |
| SLES 12               | [docker-engine]({{ rpm-prefix }}/opensuse/12.3/Packages/docker-engine-{{ rpm-version }}.x86_64.rpm)                                                                                                                                                                                                                                                                |
| Ubuntu Xenial         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-xenial_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Wily           | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-wily_amd64.deb)                                                                                                                                                                                                                                                                            |
| Ubuntu Trusty         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-trusty_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Precise        | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-precisel_amd64.deb)                                                                                                                                                                                                                                                                        |

#### CS Docker Engine 1.12.5

{% comment %} Check on the S3 bucket for packages.docker.com for the versions. {% endcomment %}
{% assign rpm-version = "1.12.5.cs5-1" %}
{% assign deb-version = "1.12.5~cs5-0" %}

| Operating system      | Package links                                                                                                                                                                                                                                                                                                                                               |
|-----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RHEL 7.x and CentOS 7 | [docker-engine]({{ rpm-prefix }}/centos/7/Packages/docker-engine-{{ rpm-version}}.el7.centos.x86_64.rpm), [docker-engine-debuginfo]({{ rpm-prefix }}/centos/7/Packages/docker-engine-debuginfo-{{ rpm-version }}.el7.centos.x86_64.rpm), [docker-engine-selinux]({{ rpm-prefix }}/centos/7/Packages/docker-engine-selinux-{{ rpm-version}}1.el7.centos.noarch.rpm) |
| SLES 12               | [docker-engine]({{ rpm-prefix }}/opensuse/12.3/Packages/docker-engine-{{ rpm-version }}.x86_64.rpm)                                                                                                                                                                                                                                                                |
| Ubuntu Xenial         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-xenial_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Wily           | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-wily_amd64.deb)                                                                                                                                                                                                                                                                            |
| Ubuntu Trusty         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-trusty_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Precise        | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-precisel_amd64.deb)                                                                                                                                                                                                                                                                        |

#### CS Docker Engine 1.12.3

{% comment %} Check on the S3 bucket for packages.docker.com for the versions. {% endcomment %}
{% assign rpm-version = "1.12.3.cs4-1" %}
{% assign deb-version = "1.12.3~cs4-0" %}

| Operating system      | Package links                                                                                                                                                                                                                                                                                                                                               |
|-----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RHEL 7.x and CentOS 7 | [docker-engine]({{ rpm-prefix }}/centos/7/Packages/docker-engine-{{ rpm-version}}.el7.centos.x86_64.rpm), [docker-engine-debuginfo]({{ rpm-prefix }}/centos/7/Packages/docker-engine-debuginfo-{{ rpm-version }}.el7.centos.x86_64.rpm), [docker-engine-selinux]({{ rpm-prefix }}/centos/7/Packages/docker-engine-selinux-{{ rpm-version}}1.el7.centos.noarch.rpm) |
| SLES 12               | [docker-engine]({{ rpm-prefix }}/opensuse/12.3/Packages/docker-engine-{{ rpm-version }}.x86_64.rpm)                                                                                                                                                                                                                                                                |
| Ubuntu Xenial         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-xenial_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Wily           | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-wily_amd64.deb)                                                                                                                                                                                                                                                                            |
| Ubuntu Trusty         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-trusty_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Precise        | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-precisel_amd64.deb)                                                                                                                                                                                                                                                                        |

#### CS Docker Engine 1.12.2

{% comment %} Check on the S3 bucket for packages.docker.com for the versions. {% endcomment %}
{% assign rpm-version = "1.12.2.cs2-1" %}
{% assign deb-version = "1.12.2~cs2-0" %}

| Operating system      | Package links                                                                                                                                                                                                                                                                                                                                               |
|-----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RHEL 7.x and CentOS 7 | [docker-engine]({{ rpm-prefix }}/centos/7/Packages/docker-engine-{{ rpm-version}}.el7.centos.x86_64.rpm), [docker-engine-debuginfo]({{ rpm-prefix }}/centos/7/Packages/docker-engine-debuginfo-{{ rpm-version }}.el7.centos.x86_64.rpm), [docker-engine-selinux]({{ rpm-prefix }}/centos/7/Packages/docker-engine-selinux-{{ rpm-version}}1.el7.centos.noarch.rpm) |
| SLES 12               | [docker-engine]({{ rpm-prefix }}/opensuse/12.3/Packages/docker-engine-{{ rpm-version }}.x86_64.rpm)                                                                                                                                                                                                                                                                |
| Ubuntu Xenial         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-xenial_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Wily           | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-wily_amd64.deb)                                                                                                                                                                                                                                                                            |
| Ubuntu Trusty         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-trusty_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Precise        | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-precisel_amd64.deb)                                                                                                                                                                                                                                                                        |

#### CS Docker Engine 1.12.1

{% comment %} Check on the S3 bucket for packages.docker.com for the versions. {% endcomment %}
{% assign rpm-version = "1.12.1.cs1-1" %}
{% assign deb-version = "1.12.1~cs1-0" %}

| Operating system      | Package links                                                                                                                                                                                                                                                                                                                                               |
|-----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RHEL 7.x and CentOS 7 | [docker-engine]({{ rpm-prefix }}/centos/7/Packages/docker-engine-{{ rpm-version}}.el7.centos.x86_64.rpm), [docker-engine-debuginfo]({{ rpm-prefix }}/centos/7/Packages/docker-engine-debuginfo-{{ rpm-version }}.el7.centos.x86_64.rpm), [docker-engine-selinux]({{ rpm-prefix }}/centos/7/Packages/docker-engine-selinux-{{ rpm-version}}1.el7.centos.noarch.rpm) |
| SLES 12               | [docker-engine]({{ rpm-prefix }}/opensuse/12.3/Packages/docker-engine-{{ rpm-version }}.x86_64.rpm)                                                                                                                                                                                                                                                                |
| Ubuntu Xenial         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-xenial_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Wily           | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-wily_amd64.deb)                                                                                                                                                                                                                                                                            |
| Ubuntu Trusty         | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-trusty_amd64.deb)                                                                                                                                                                                                                                                                          |
| Ubuntu Precise        | [docker-engine]({{ deb-prefix }}/docker-engine_{{ deb-version }}~ubuntu-precisel_amd64.deb)                                                                                                                                                                                                                                                                        |
