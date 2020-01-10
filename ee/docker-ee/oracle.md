---
description: Instructions for installing Docker Engine - Enterprise on Oracle Linux
keywords: requirements, installation, oracle, ol, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/oracle/
- /engine/installation/linux/oracle/
- /engine/installation/linux/docker-ee/oracle/
- /install/linux/docker-ee/oracle/
title: Get Docker Engine - Enterprise for Oracle Linux
---

{% assign linux-dist = "oraclelinux" %}
{% assign linux-dist-cap = "OL" %}
{% assign linux-dist-url-slug = "oraclelinux" %}
{% assign linux-dist-long = "Oracle Linux" %}
{% assign package-format = "RPM" %}
{% assign gpg-fingerprint = "77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9" %}

>{% include enterprise_label_shortform.md %}

There are two ways to install and upgrade [Docker Enterprise](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }
on {{ linux-dist-long }}:

- [YUM repository](#repo-install-and-upgrade): Set up a Docker repository and install Docker Engine - Enterprise from it. This is the recommended approach because installation and upgrades are managed with YUM and easier to do.

- [RPM package](#package-install-and-upgrade): Download the {{ package-format }} package, install it manually, and manage upgrades manually. This is useful when installing Docker Engine - Enterprise on air-gapped systems with no access to the internet.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

Docker Engine - Community is _not_ supported on {{ linux-dist-long }}.

<!---
Shared between rhel.md, oracle.md
--->

## Prerequisites

This section lists what you need to consider before installing Docker Engine -
Enterprise. Items that require action are explained below.

- Use {{ linux-dist-cap }} 64-bit 7.3 or higher on RHCK 3.10.0-514 or higher.
- Use the `devicemapper` storage driver only (`direct-lvm` mode in production).
- Find the URL for your Docker Engine - Enterprise repo at [Docker Hub](https://hub.docker.com/my-content){: target="_blank" class="_" }.
- Uninstall old versions of Docker.
- Remove old Docker repos from `/etc/yum.repos.d/`.

### Architectures and storage drivers

Docker Engine - Enterprise supports {{ linux-dist-long }} 64-bit, versions 7.3
and higher, running the Red Hat Compatible kernel (RHCK) 3.10.0-514 or higher.
Older versions of {{ linux-dist-long }} are not supported.

On {{ linux-dist-long }}, Docker Engine - Enterprise only supports the
`devicemapper` storage driver. In production, you must use it in `direct-lvm`
mode, which requires one or more dedicated block devices. Fast storage such as
solid-state media (SSD) is recommended. Do not start Docker until properly
configured per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

### Find your Docker Engine - Enterprise repo URL

To install Docker Enterprise, you will need the URL of the Docker Enterprise repository associated with your trial or subscription:

1.  Go to [https://hub.docker.com/my-content](https://hub.docker.com/my-content){: target="_blank" class="_" }. All of your subscriptions and trials are listed.
2.  Click the **Setup** button for **Docker Enterprise Edition for {{ linux-dist-long }}**.
3.  Copy the URL from **Copy and paste this URL to download your Edition** and save it for later use.

You will use this URL in a later step to create a variable called, `DOCKERURL`.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

### Uninstall old Docker versions

The Docker Engine - Enterprise package is called `docker-ee`. Older versions
were called `docker` or `docker-engine`. Uninstall all older versions and
associated dependencies. The contents of `/var/lib/docker/` are preserved,
including images, containers, volumes, and networks.

```bash
$ sudo yum remove docker \
                  docker-engine \
                  docker-engine-selinux
```

## Repo install and upgrade

The advantage of using a repository from which to install Docker Engine - Enterprise (or any software) is that it provides a certain level of automation. RPM-based distributions such as {{ linux-dist-long }}, use a tool called YUM that work with your repositories to manage dependencies and provide automatic updates.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

### Set up the repository

You only need to set up the repository once, after which you can install Docker Engine - Enterprise _from_ the repo and repeatedly upgrade as necessary.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->


1.  Remove existing Docker repositories from `/etc/yum.repos.d/`:

    ```bash
    $ sudo rm /etc/yum.repos.d/docker*.repo
    ```

2.  Temporarily store the URL (that you [copied above](#find-your-docker-ee-repo-url)) in an environment variable. Replace `<DOCKER-EE-URL>` with your URL in the following command. This variable assignment does not persist when the session ends:

    ```bash
    $ export DOCKERURL="<DOCKER-EE-URL>"
    ```

3.  Store the value of the variable, `DOCKERURL` (from the previous step), in a `yum` variable in `/etc/yum/vars/`:

    ```bash
    $ sudo -E sh -c 'echo "$DOCKERURL/{{ linux-dist-url-slug }}" > /etc/yum/vars/dockerurl'
    ```

4.  Install required packages: `yum-utils` provides the _yum-config-manager_ utility, and `device-mapper-persistent-data` and `lvm2` are required by the _devicemapper_ storage driver:

    ```bash
    $ sudo yum install -y yum-utils \
      device-mapper-persistent-data \
      lvm2
    ```

<!---
Shared between centOS.md, oracle.md
--->

5.  Enable the `ol7_addons` Oracle repository. This ensures access to the `container-selinux` package required by `docker-ee`.

    ```bash
    $ sudo yum-config-manager --enable ol7_addons
    ```

<!---
Not shared
--->


6.  Add the Docker Engine - Enterprise **stable** repository:

    ```bash
    $ sudo -E yum-config-manager \
        --add-repo \
        "$DOCKERURL/{{ linux-dist-url-slug }}/docker-ee.repo"
    ```

<!---
Shared between centOS.md, oracle.md
--->

### Install from the repository


> **Note**: If you need to run Docker Engine - Enterprise 2.0, please see the following instructions:
> * [18.03](https://docs.docker.com/v18.03/ee/supported-platforms/) - Older Docker Engine - Enterprise Engine only release
> * [17.06](https://docs.docker.com/v17.06/engine/installation/) - Docker Enterprise Edition 2.0 (Docker Engine,
> UCP, and DTR).

1.  Install the latest patch release, or go to the next step to install a specific version:

    ```bash
    $ sudo yum -y install docker-ee docker-ee-cli containerd.io
    ```

    If prompted to accept the GPG key, verify that the fingerprint matches `{{ gpg-fingerprint }}`, and if so, accept it.


2.  To install a _specific version_ of Docker Engine - Enterprise (recommended in production), list versions and install:

    a.  List and sort the versions available in your repo. This example sorts results by version number, highest to lowest, and is truncated:

    ```bash
    $ sudo yum list docker-ee  --showduplicates | sort -r

    docker-ee.x86_64      {{ site.docker_ee_version }}.ee.2-1.el7.{{ linux-dist }}      docker-ee-stable-18.09
    ```

    The list returned depends on which repositories you enabled, and is specific to your version of {{ linux-dist-long }} (indicated by `.el7` in this example).

    b.  Install a specific version by its fully qualified package name, which is the package name (`docker-ee`) plus the version string (2nd column) starting at the first colon (`:`), up to the first hyphen, separated by a hyphen (`-`). For example, `docker-ee-18.09.1`.

    ```bash
    $ sudo yum -y install docker-ee-<VERSION_STRING> docker-ee-cli-<VERSION_STRING> containerd.io
    ```

    For example, if you want to install the 18.09 version run the following:

    ```bash
    sudo yum-config-manager --enable docker-ee-stable-18.09
    ```

    Docker is installed but not started. The `docker` group is created, but no users are added to the group.

3.  Start Docker:

    > If using `devicemapper`, ensure it is properly configured before starting Docker, per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

    ```bash
    $ sudo systemctl start docker
    ```

4.  Verify that Docker Engine - Enterprise is installed correctly by running the `hello-world`
    image. This command downloads a test image, runs it in a container, prints
    an informational message, and exits:

    ```bash
    $ sudo docker run hello-world
    ```

    Docker Engine - Enterprise is installed and running. Use `sudo` to run Docker commands. See
    [Linux postinstall](/install/linux/linux-postinstall.md){: target="_blank" class="_" } to allow
    non-privileged users to run Docker commands.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

### Upgrade from the repository

1.  [Add the new repository](#set-up-the-repository).

2.  Follow the [installation instructions](#install-from-the-repository) and install a new version.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

## Package install and upgrade

To manually install Docker Enterprise, download the `.{{ package-format | downcase }}` file for your release. You need to download a new file each time you want to upgrade Docker Enterprise.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

{{ selinux-warning }}

### Install with a package

1.  Go to the Docker Engine - Enterprise repository URL associated with your
    trial or subscription in your browser. Go to
    `{{ linux-dist-url-slug }}/`. Choose your {{ linux-dist-long }} version,
    architecture, and Docker version. Download the
    `.{{ package-format | downcase }}` file from the `Packages` directory.

<!---
Not shared
--->

2.  Install Docker Enterprise, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo yum install /path/to/package.rpm
    ```

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

3.  Start Docker:

    > If using `devicemapper`, ensure it is properly configured before starting Docker, per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

    ```bash
    $ sudo systemctl start docker
    ```

4.  Verify that Docker Engine - Enterprise is installed correctly by running the `hello-world`
    image. This command downloads a test image, runs it in a container, prints
    an informational message, and exits:

    ```bash
    $ sudo docker run hello-world
    ```

    Docker Engine - Enterprise is installed and running. Use `sudo` to run Docker commands. See
    [Linux postinstall](/install/linux/linux-postinstall.md){: target="_blank" class="_" } to allow
    non-privileged users to run Docker commands.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

### Upgrade with a package

1.  Download the newer package file.

2.  Repeat the [installation procedure](#install-with-a-package), using
    `yum -y upgrade` instead of `yum -y install`, and point to the new file.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

## Uninstall Docker Engine - Enterprise

1.  Uninstall the Docker Engine - Enterprise package:

    ```bash
    $ sudo yum -y remove docker-ee
    ```

2.  Delete all images, containers, and volumes (because these are not automatically removed from your host):

    ```bash
    $ sudo rm -rf /var/lib/docker
    ```

3.  Delete other Docker related resources:
    ```bash
    $ sudo rm -rf /run/docker
    $ sudo rm -rf /var/run/docker
    $ sudo rm -rf /etc/docker
    ```

4.  If desired, remove the `devicemapper` thin pool and reformat the block
    devices that were part of it.

You must delete any edited configuration files manually.

<!---
Shared between centOS.md, rhel.md, oracle.md
--->

## Next steps


- Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md){: target="_blank" class="_" }

- Continue with user guides on [Universal Control Plane (UCP)](/ee/ucp/){: target="_blank" class="_" } and [Docker Trusted Registry (DTR)](/ee/dtr/){: target="_blank" class="_" }
