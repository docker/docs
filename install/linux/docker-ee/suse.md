---
description: Instructions for installing Docker EE on SLES
keywords: requirements, apt, installation, suse, opensuse, sles, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/SUSE/
- /engine/installation/linux/suse/
- /engine/installation/linux/docker-ee/suse/
title: Get Docker EE for SLES
toc_max: 4
---

To get started with Docker on SUSE Linux Enterprise Server (SLES), make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ee).

## Prerequisites

### Docker EE URL

To install Docker Enterprise Edition (Docker EE), you need to know the Docker EE
repository URL associated with your trial or subscription. These instructions
work for Docker EE for SLES and for Docker EE for Linux, which includes access
to Docker EE for all Linux distributions. To get this information:

- Go to [https://store.docker.com/my-content](https://store.docker.com/my-content).
- Each subscription or trial you have access to is listed. Click the **Setup**
  button for **Docker Enterprise Edition for SUSE Linux Enterprise Server**.
- Copy the URL from the field labeled
  **Copy and paste this URL to download your Edition**.

Use this URL when you see the placeholder text `<DOCKER-EE-URL>`.

To learn more about Docker EE, see
[Docker Enterprise Edition](https://www.docker.com/enterprise-edition/){: target="_blank" class="_" }.

Docker Community Edition (Docker CE) is not supported on SLES.

### OS requirements

To install Docker EE, you need the 64-bit version of SLES 12.x, running on
`x86_64`, `s390x` (IBM Z), or `ppc64le` (IBM Power) architectures. Docker EE is
not supported on OpenSUSE.

The only supported storage driver for Docker EE on SLES is Btrfs, which is
used by default if the underlying filesystem hosting `/var/lib/docker/` is a
BTRFS filesystem.

#### Firewall configuration

Docker creates a `DOCKER` iptables chain when it starts. The SUSE firewall may
block access to this chain, which can prevent you from running
containers with published ports. You may see errors such as the following:

```none
WARNING: IPv4 forwarding is disabled. Networking will not work.
docker: Error response from daemon: driver failed programming external
        connectivity on endpoint adoring_ptolemy
        (0bb5fa80bc476f8a0d343973929bb3b7c039fc6d7cd30817e837bc2a511fce97):
        (iptables failed: iptables --wait -t nat -A DOCKER -p tcp -d 0/0 --dport 80 -j DNAT --to-destination 172.17.0.2:80 ! -i docker0: iptables: No chain/target/match by that name.
 (exit status 1)).
```

If you see errors like this, adjust the start-up script order so that the
firewall is started before Docker, and Docker stops before the firewall stops.
See the
[SLES documentation on init script order](https://www.suse.com/documentation/sled11/book_sle_admin/data/sec_boot_init.html).

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If you use OS
images from a cloud provider, you may need to remove the `runc` package, which
conflicts with Docker EE. If these are installed, uninstall them, along with
associated dependencies.

```bash
$ sudo zypper rm docker docker-engine runc
```

If removal of the `docker-engine` package fails, use the following command
instead:

```bash
$ sudo rpm -e docker-engine
```

It's OK if `zypper` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker EE package is now called `docker-ee`.

## Configure the Btrfs filesystem

By default, SLES formats the `/` filesystem using Btrfs, so **most people do not
not need to do the steps in this section**. If you use OS images from a cloud
provider, you may need to do this step. If the filesystem that
hosts `/var/lib/docker/` is **not** a BTRFS filesystem, you must configure a
BTRFS filesystem and mount it on `/var/lib/docker/`.

1.  Check whether `/` (or `/var/` or `/var/lib/` or `/var/lib/docker/` if they
    are separate mount points) are formatted using Btrfs. If you do not have
    separate mount points for any of these, a duplicate result for `/` is
    returned.

    ```bash
    $ df -T / /var /var/lib /var/lib/docker
    ```

    You need to complete the rest of these steps **only if one of the following
    is true**:

    - You have a separate `/var/` filesystem that is not formatted with Btrfs
    - You do not have a separate `/var/` or `/var/lib/` or `/var/lib/docker/`
      filesystem and `/` is not formatted with Btrfs

    If `/var/lib/docker` is already a separate mount point and is not formatted
    with Btrfs, back up its contents so that you can restore them after step
    3.

2.  Format your dedicated block device or devices as a Btrfs filesystem. This
    example assumes that you are using two block devices called `/dev/xvdf` and
    `/dev/xvdg`. **Make sure you are using the right device names.**

    > Double-check the block device names because this is a
    destructive operation.
    {:.warning}

    ```bash
    $ sudo mkfs.btrfs -f /dev/xvdf /dev/xvdg
    ```

    There are many more options for Btrfs, including striping and RAID. See the
    [Btrfs documentation](https://btrfs.wiki.kernel.org/index.php/Using_Btrfs_with_Multiple_Devices).

3.  Mount the new Btrfs filesystem on the `/var/lib/docker/` mount point. You
    can specify any of the block devices used to create the Btrfs filesystem.

    ```bash
    $ sudo mount -t btrfs /dev/xvdf /var/lib/docker
    ```

    Don't forget to make the change permanent across reboots by adding an
    entry to `/etc/fstab`.

4.  If `/var/lib/docker` previously existed and you backed up its contents
    during step 1, restore them onto `/var/lib/docker`.


## Install Docker EE

You can install Docker EE in different ways, depending on your needs:

- Most users
  [set up Docker's repositories](#install-using-the-repository) and install
  from them, for ease of installation and upgrade tasks. This is the
  recommended approach.

- Some users download the RPM package and install it manually and manage
  upgrades completely manually. This is useful in situations such as installing
  Docker on air-gapped systems with no access to the internet.

### Install using the repository

Before you install Docker EE for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker EE
from the repository.

#### Set up the repository

1.  Temporarily add a `$DOCKER_EE_URL` variable into your environment. This
    only persists until you log out of the session. Replace `<DOCKER-EE-URL>`
    with the URL you noted down in the [prerequisites](#prerequisites).

    ```bash
    $ DOCKER_EE_URL="<DOCKER-EE-URL>"
    ```

2.  Use the following command to set up the **stable** repository. Use the
    command as-is. It works because of the variable you set in the previous
    step.

    <ul class="nav nav-tabs">
      <li class="active"><a data-toggle="tab" data-target="#x86_64_repo">x86_64 / amd64</a></li>
      <li><a data-toggle="tab" data-target="#s390x_repo">IBM Z (s390x)</a></li>
      <li><a data-toggle="tab" data-target="#ppc64le_repo">IBM Power (ppc64le)</a></li>
    </ul>
    <div class="tab-content">
    <div id="x86_64_repo" class="tab-pane fade in active" markdown="1">

    ```bash
    $ sudo zypper addrepo \
        "${DOCKER_EE_URL}/sles/12.3/x86_64/stable-{{ site.docker_ee_version }}" \
        docker-ee-stable
    ```

    </div>
    <div id="s390x_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo zypper addrepo \
        "${DOCKER_EE_URL}/sles/12.3/s390x/stable-{{ site.docker_ee_version }}" \
        docker-ee-stable
    ```

    </div>

    <div id="ppc64le_repo" class="tab-pane fade" markdown="1">

    ```bash
    $ sudo zypper addrepo \
        "${DOCKER_EE_URL}/sles/12.3/ppc64le/stable-{{ site.docker_ee_version }}" \
        docker-ee-stable
    ```

    </div>
    </div><!--tab-content-->


3.  Import the GPG key from the repository. Use the command as-is. It works
    because of the variable you set earlier.

    ```bash
    $ sudo rpm --import "${DOCKER_EE_URL}/sles/gpg"
    ```

#### Install Docker EE

1.  Update the `zypper` package index.

    ```bash
    $ sudo zypper refresh
    ```

    If this is the first time you have refreshed the package index since adding
    the Docker repositories, you are prompted to accept the GPG key, and
    the key's fingerprint is shown. Verify that the fingerprint matches
    `77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9` and if so, accept the
    key.

2.  Install the latest version of Docker EE, or go to the next step to install a
    specific version.

    ```bash
    $ sudo zypper install docker-ee
    ```

    Start Docker:

    ```bash
    $ sudo service docker start
    ```

3.  On production systems, you should install a specific version of Docker EE
    instead of always using the latest. List the available versions. The
    following example only lists binary packages and is truncated. To also list
    source packages, omit the `-t package` flag from the command.

    ```bash
    $ zypper search -s --match-exact -t package docker-ee

      Loading repository data...
      Reading installed packages...

      S | Name          | Type    | Version                               | Arch   | Repository    
      --+---------------+---------+---------------------------------------+--------+---------------
        | docker-ee     | package | {{ site.docker_ee_version }}-1                 | x86_64 | docker-ee-stable
    ```

    The contents of the list depend upon which repositories you have enabled.
    Choose a specific version to install. The third column is the version
    string. The fifth column is the repository name, which indicates which
    repository the package is from and by extension its stability level. To
    install a specific version, append the version string to the package name
    and separate them by a hyphen (`-`):

    ```bash
    $ sudo zypper install docker-ee-<VERSION_STRING>
    ```

    Docker is installed but not started. The `docker` group is created, but no
    users are added to the group.

4.  Configure Docker EE to use the Btrfs filesystem. **This is only required if
    the `/` filesystem is not using BTRFS.** However, explicitly specifying the
    `storage-driver` has no harmful side effects.

    Edit the file `/etc/docker/daemon.json` (create it if it does not exist) and
    add the following contents:

    ```json
    {
      "storage-driver": "btrfs"
    }
    ```

    Save and close the file.

5.  Start Docker:

    ```bash
    $ sudo service docker start
    ```

6.  Verify that Docker EE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker EE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Linux postinstall](/install/linux/linux-postinstall.md) to configure the
graph storage driver, allow non-privileged users to run Docker commands, and for
other optional configuration steps.

> **Important**: Be sure Docker is configured to start after the system
> firewall. See [Firewall configuration](#firewall-configuration).

#### Upgrade Docker EE

To upgrade Docker EE:

1.  If upgrading to a new major Docker EE version (such as when going from
    Docker 17.03.x to Docker 17.06.x),
    [add the new repository](#set-up-the-repository){: target="_blank" class="_" }.

2.  Run `sudo zypper refresh`.

3.  Follow the
    [installation instructions](#install-docker), choosing the new version you want
    to install.

### Install from a package

If you cannot use the official Docker repository to install Docker EE, you can
download the `.rpm` file for your release and install it manually. You
need to download a new file each time you want to upgrade Docker EE.

1.  Go to the Docker EE repository URL associated with your
    trial or subscription in your browser. Go to `sles/12.3/` choose the
    directory corresponding to your architecture and desired Docker EE version.
    Download the `.rpm` file from the `Packages` directory.

2.  Import Docker's official GPG key:

    ```bash
    $ sudo rpm --import <DOCKER-EE-URL>/sles/gpg
    ```

3.  Install Docker EE, changing the path below to the path where you downloaded
    the Docker package.

    ```bash
    $ sudo zypper install /path/to/package.rpm
    ```

    Docker EE is installed but not started. The `docker` group is created, but no
    users are added to the group.

4.  Configure Docker EE to use the Btrfs filesystem. **This is only required if
    the `/` filesystem is not using Btrfs.** However, explicitly specifying the
    `storage-driver` has no harmful side effects.

    Edit the file `/etc/docker/daemon.json` (create it if it does not exist) and
    add the following contents:

    ```json
    {
      "storage-driver": "btrfs"
    }
    ```

    Save and close the file.

5.  Start Docker:

    ```bash
    $ sudo service docker start
    ```

6.  Verify that Docker EE is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

Docker EE is installed and running. You need to use `sudo` to run Docker
commands. Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md)
to allow non-privileged users to run Docker commands and for other optional
configuration steps.

> **Important**: Be sure Docker is configured to start after the system
> firewall. See [Firewall configuration](#firewall-configuration).

#### Upgrade Docker EE

To upgrade Docker EE, download the newer package file and repeat the
[installation procedure](#install-from-a-package), using `zypper update`
instead of `zypper install`, and pointing to the new file.

## Uninstall Docker EE

1.  Uninstall the Docker EE package using the following command.

    ```bash
    $ sudo zypper rm docker-ee
    ```

2.  Images, containers, volumes, or customized configuration files on your host
    are not automatically removed. To delete all images, containers, and
    volumes:

    ```bash
    $ sudo rm -rf /var/lib/docker/*
    ```

    If you used a separate BTRFS filesystem to host the contents of
    `/var/lib/docker/`, you can unmount and format the Btrfs filesystem.

You must delete any edited configuration files manually.

## Next steps

- Continue to [Post-installation steps for Linux](/install/linux/linux-postinstall.md)

- Continue with the [User Guide](/engine/userguide/index.md).
