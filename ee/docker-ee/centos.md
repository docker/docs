---
description: Instructions for installing Docker Engine - Enterprise on CentOS
keywords: requirements, apt, installation, centos, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/centos/
- /engine/installation/linux/docker-ee/centos/
- /install/linux/docker-ee/centos/
title: Get Docker Engine - Enterprise for CentOS
---

{% assign linux-dist = "centos" %}
{% assign linux-dist-cap = "CentOS" %}
{% assign linux-dist-url-slug = "centos" %}
{% assign linux-dist-long = "Centos" %}
{% assign package-format = "RPM" %}
{% assign gpg-fingerprint = "77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9" %}

>{% include enterprise_label_shortform.md %}

> **Important** 
> 
> Docker Engine - Community users should go to
[Get Docker Engine - Community for Centos](/install/linux/docker-ce/centos.md)
**instead of this topic**. 
{: .important}

## Prerequisites

Confirm that all prerequisites are met before installing Docker Engine - Enterprise on CentOS. These prerequisites include:

- Confirming architecture and storage drivers:
    * {{ linux-dist-cap }} 64-bit 7.1 and higher on `x86_64`
    * Storage driver `overlay2` or `devicemapper` (`direct-lvm` mode in
  production)
- Locating the URL for the Docker Engine - Enterprise repo
- Uninstalling all old Docker versions
- Removing old Docker repos from `/etc/yum.repos.d/`

> **Note**
>
> Learn more about Docker Engine - Enterprise at [Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

### Confirming Architectures and Storage Drivers

Docker Engine - Enterprise supports {{ linux-dist-long }} 64-bit, latest
version, running on  `x86_64`.

On {{ linux-dist-long }}, Docker Engine - Enterprise supports storage drivers,
`overlay2` and `devicemapper`. In Docker Engine - Enterprise 17.06.2-ee-5 and
higher, `overlay2` is the recommended storage driver. The following limitations
apply:

- [OverlayFS](/storage/storagedriver/overlayfs-driver){: target="_blank" class="_" }:
  If `selinux` is enabled, the `overlay2` storage driver is supported on
  {{ linux-dist-cap }} 7.4 or higher. If `selinux` is disabled, `overlay2` is
  supported on {{ linux-dist-cap }} 7.2 or higher with kernel version 3.10.0-693
  and higher.

- [Device Mapper](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }:
  On production systems using `devicemapper`, you must use `direct-lvm` mode,
  which requires one or more dedicated block devices. Fast storage such as
  solid-state media (SSD) is recommended. Do not start Docker until properly
  configured per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

### Locating the URL for the Docker Engine - Enterprise Repo

The product URL associated with your trial or subscription is required to install Docker Engine - Enterprise, which can be attained using the following procedure (note that these instructions apply to Docker on CentOS and for Docker on Linux, which includes access to Docker Engine - Enterprise for all Linux distributions):

1. Go to [https://hub.docker.com/my-content](https://hub.docker.com/my-content).
2. Each subscription or trial you have access to is listed. Click the **Setup**
  button for **Docker Enterprise Edition for CentOS**.
3. Copy the URL from  **Copy and paste this URL to download your Edition** and save it for later use.

Note that the URL will be put to use in a later step to create a variable called `DOCKERURL`.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

### Uninstalling All Old Docker Versions

Use the `yum remove` command to uninstall older versions and associated dependencies of Docker Engine - Enterprise (called `docker` or `docker-engine`. Note that the contents of `/var/lib/docker/` are preserved, including images, containers, volumes, and networks. In addition, if you are upgrading from Docker Engine - Community to Docker Engine - Enterprise, remove the Docker Engine - Community package.

```bash
$ sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine
```

## Install Docker Engine - Enterprise

Docker Engine - Enterprise can be installed either via YUM repositories, or by downloading and installing the RPM package and thereafter manually managing all upgrades. The Docker repository method is recommended, for the ease it lends in terms of both installation and upgrade tasks. The more manual RPM package approach, however, is useful in certain situations, such as installing Docker on air-gapped system that have no access to the Internet.

### Installing from a YUM Repository

The advantage of using a repository from which to install Docker Engine - Enterprise (or any software) is that it provides a certain level of automation. RPM-based distributions such as {{ linux-dist-long }}, use a tool called YUM that work with your repositories to manage dependencies and provide automatic updates.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

#### Set up the Repository

Naturally, to install Docker Engine - Enterprise on a new host machine using the Docker repository you must first set the repository up on the machine. It is only necessary to set up the repository once, after which Docker Engine - Enterprise can be installed _from_ the repo and repeatedly upgraded as necessary.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

1.  Remove existing Docker repositories from `/etc/yum.repos.d/`.

    ```bash
    $ sudo rm /etc/yum.repos.d/docker*.repo
    ```

2. Temporarily add a `$DOCKER_EE_URL` variable into your environment (it persists only up until you log out of the session). Replace `<DOCKER-EE-URL>` with the URL you noted down in the [prerequisites](#prerequisites).

    ```bash
    $ export DOCKERURL="<DOCKER-EE-URL>"
    ```

3.  Store the value of the variable, `DOCKERURL` (from the previous step), in a `yum` variable in `/etc/yum/vars/`.

    ```bash
    $ sudo -E sh -c 'echo "$DOCKERURL/{{ linux-dist-url-slug }}" > /etc/yum/vars/dockerurl'
    ```

4.  Install required packages.

    `yum-utils` provides the _yum-config-manager_ utility, and `device-mapper-persistent-data` and `lvm2` are required by the _devicemapper_ storage driver.

    ```bash
    $ sudo yum install -y yum-utils \
      device-mapper-persistent-data \
      lvm2
    ```

5.  Add the Docker Engine - Enterprise **stable** repository.

    ```bash
    $ sudo -E yum-config-manager \
        --add-repo \
        "$DOCKERURL/{{ linux-dist-url-slug }}/docker-ee.repo"
    ```

<!---
Shared between centOS.md, oracle.md
--->

#### Install from the Repository

> **Note**
>
> To run Docker Engine - Enterprise 2.0, refer to:
> * [18.03](https://docs.docker.com/v18.03/ee/supported-platforms/) - Older Docker Engine - Enterprise Engine only release
> * [17.06](https://docs.docker.com/v17.06/engine/installation/) - Docker Enterprise Edition 2.0 (Docker Engine, UCP, and DTR).

1. Install either the latest patch release or a _specific version_ of Docker Engine - Enterprise.

    * To install the latest patch release:

        ```bash
        $ sudo yum -y install docker-ee docker-ee-cli containerd.io
        ```

        If prompted to accept the GPG key, verify that the fingerprint matches `{{ gpg-fingerprint }}`, and if so, accept it.

    **— or —**

    * To install a _specific version_ of Docker Engine - Enterprise (recommended in production), list versions and install:

        a. List and sort the versions available in your repo. This example sorts   results by version number, highest to lowest, and is truncated:

            ```bash
            $ sudo yum list docker-ee  --showduplicates | sort -r

            docker-ee.x86_64      {{ site.docker_ee_version }}.ee.2-1.el7.{{ linux-dist }}      docker-ee-stable-18.09
            ```

        The list returned depends on which repositories you enabled, and is specific to your version of {{ linux-dist-long }} (indicated by `.el7` in this example).

        b.  Install a specific version by its fully qualified package name, which is the package name (`docker-ee`) plus the version string (2nd column) starting at the first colon (`:`), up to the first hyphen, separated by a hyphen (`-`). For example, `docker-ee-18.09.1`.

        ```bash
        $ sudo yum -y install docker-ee-<VERSION_STRING> docker-ee-cli-<VERSION_STRING> containerd.io
        ```

        For example, run the following to install the 18.09 version:

        ```bash
        sudo yum-config-manager --enable docker-ee-stable-18.09
        ```

    Docker is installed but not started. The `docker` group is created, but no users are added to the group.

2.  Start Docker.

    > **Note**
    >
    > If using `devicemapper`, ensure it is properly configured before starting Docker, per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

    ```bash
    $ sudo systemctl start docker
    ```

3.  Verify that Docker Engine - Enterprise is installed correctly by running the `hello-world` image.

    ```bash
    $ sudo docker run hello-world
    ```
    This command downloads a test image, runs it in a container, prints an informational message, and exits.

Docker Engine - Enterprise is installed and running. Use `sudo` to run Docker commands. Continue to [Linux postinstall](/install/linux/linux-postinstall.md){: target="_blank" class="_" } to allow non-privileged users to run Docker commands and for other optional configuration steps.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

#### Upgrade from the Repository

1.  [Add the new repository](#set-up-the-repository).

2.  Follow the [installation instructions](#install-from-the-repository) and install a new version.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

### Installing and Upgrading from an RPM Package

To manually install Docker Enterprise, download the `.{{ package-format | downcase }}` file for your release. Note that it will be necessry to download a new file each time you want to upgrade Docker Enterprise.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

#### Install with a Package

1.  Download the file for Docker Enterprise - Engine.

    a. Use a browser to go to the Docker Engine - Enterprise repository URL associated with your trial or subscription.

    b. Go to `{{ linux-dist-url-slug }}/7/x86_64/stable-<VERSION>/Packages` and download the `.{{ package-format | downcase }}` file for the desired version.

    <!---
    Not shared
    --->

2.  Install Docker Enterprise.

    ```bash
    $ sudo yum install </path/to/package>/<file>.rpm
    ```

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

3.  Start Docker.

    > **Note**
    >
    >If using `devicemapper`, ensure it is properly configured before starting Docker, per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

    ```bash
    $ sudo systemctl start docker
    ```

4.  Verify that Docker Engine - Enterprise is installed correctly by running the `hello-world` image.

    ```bash
    $ sudo docker run hello-world
    ```
    This command downloads a test image, runs it in a container, prints an informational message, and exits.

Docker Engine - Enterprise is installed and running. Use `sudo` to run Docker commands. Continue to [Linux postinstall](/install/linux/linux-postinstall.md){: target="_blank" class="_" } to allow non-privileged users to run Docker commands and for other optional configuration steps.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

#### Upgrade with a Package

1.  Download the newer package file.

2.  Repeat the [installation procedure](#install-with-a-package), using
    `yum -y upgrade` instead of `yum -y install`, and point to the new file.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

## Uninstall Docker Engine - Enterprise

1.  Uninstall the Docker Engine - Enterprise package.

    ```bash
    $ sudo yum -y remove docker-ee
    ```

2.  Delete all images, containers, and volumes (as these are not automatically removed from your host).

    ```bash
    $ sudo rm -rf /var/lib/docker
    ```

3.  Delete other Docker related resources.
    ```bash
    $ sudo rm -rf /run/docker
    $ sudo rm -rf /var/run/docker
    $ sudo rm -rf /etc/docker
    ```

4.  If desired, remove the `devicemapper` thin pool and reformat the block
    devices that were part of it.

> **Note**
>
> Any edited configuration files must be manually deleted.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

## Next steps


- Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md){: target="_blank" class="_" }

- Continue with user guides on [Universal Control Plane (UCP)](/ee/ucp/){: target="_blank" class="_" } and [Docker Trusted Registry (DTR)](/ee/dtr/){: target="_blank" class="_" }
