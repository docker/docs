---
description: Instructions for installing Docker on Ubuntu
keywords: Docker, Docker documentation, requirements, apt, installation,  ubuntu
redirect_from:
- /engine/installation/ubuntulinux/
- /installation/ubuntulinux/
title: Install Docker on Ubuntu
---

Docker is supported on these Ubuntu operating systems:

- Ubuntu Xenial 16.04 (LTS)
- Ubuntu Wily 15.10
- Ubuntu Trusty 14.04 (LTS)
- Ubuntu Precise 12.04 (LTS)

This page instructs you to install Docker on Ubuntu, using packages provided by
Docker. Using these packages ensures you get the latest official
release of Docker. If you are required to install using Ubuntu-managed packages,
consult the Ubuntu documentation. Some files and commands may be different if
you use Ubuntu-managed packages.

>**Note**: Ubuntu Utopic 14.10 and 15.04 exist in Docker's `APT` repository but
are no longer officially supported.

## Prerequisites

Docker has two important installation requirements:

- Docker only works on a 64-bit Linux installation.
- Docker requires version 3.10 or higher of the Linux kernel. Kernels older than
  3.10 lack some of the features required to run Docker containers and contain
  known bugs which cause data loss and frequently panic under certain conditions.

  To check your current kernel version, open a terminal and use `uname -r` to
  display your kernel version:

  ```bash
  $ uname -r
  3.11.0-15-generic
  ```

### Update your apt sources

To set `APT` to use packages from the Docker repository:

1.  Log into your machine as a user with `sudo` or `root` privileges.

2.  Open a terminal window.

3.  Update package information, ensure that APT works with the `https` method,
    and that CA certificates are installed.

    ```bash
    $ sudo apt-get update
    $ sudo apt-get install apt-transport-https ca-certificates
    ```
4.  Add the new `GPG` key. This commands downloads the key with the ID
    `58118E89F3A912897C070ADBF76221572C52609D` from the keyserver
    `hkp://ha.pool.sks-keyservers.net:80` and adds it to the `adv` keychain.
    For more info, see the output of `man apt-key`.

    ```bash
    $ sudo apt-key adv \
                   --keyserver hkp://ha.pool.sks-keyservers.net:80 \
                   --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
    ```
    
    If the above keyserver is not available, try `hkp://pgp.mit.edu:80` or 
    `hkp://keyserver.ubuntu.com:80`.

5.  Find the entry in the table below which corresponds to your Ubuntu version.
    This determines where APT will search for Docker packages. When possible,
    run a long-term support (LTS) edition of Ubuntu.

    | Ubuntu version      | Repository                                                  |
    | ------------------- | ----------------------------------------------------------- |
    | Precise 12.04 (LTS) | `deb https://apt.dockerproject.org/repo ubuntu-precise main`|
    | Trusty 14.04 (LTS)  | `deb https://apt.dockerproject.org/repo ubuntu-trusty main` |
    | Wily 15.10          | `deb https://apt.dockerproject.org/repo ubuntu-wily main`   |
    | Xenial 16.04 (LTS)  | `deb https://apt.dockerproject.org/repo ubuntu-xenial main` |


    >**Note**: Docker does not provide packages for all architectures. Binary artifacts
    are built nightly, and you can download them from
    https://master.dockerproject.org. To install docker on a multi-architecture
    system, add an `[arch=...]` clause to the entry. Refer to
    [Debian Multiarch wiki](https://wiki.debian.org/Multiarch/HOWTO#Setting_up_apt_sources)
    for details.

6.  Run the following command, substituting the entry for your operating system
    for the placeholder `<REPO>`.

    ```bash
    $ echo "<REPO>" | sudo tee /etc/apt/sources.list.d/docker.list
    ```

7.  Update the `APT` package index.

    ```bash
    $ sudo apt-get update
    ```

8.  Verify that `APT` is pulling from the right repository.

    When you run the following command, an entry is returned for each version of
    Docker that is available for you to install. Each entry should have the URL
    `https://apt.dockerproject.org/repo/`. The version currently installed is
    marked with `***`.The output below is truncated.

    ```bash
    $ apt-cache policy docker-engine

      docker-engine:
        Installed: 1.12.2-0~trusty
        Candidate: 1.12.2-0~trusty
        Version table:
       *** 1.12.2-0~trusty 0
              500 https://apt.dockerproject.org/repo/ ubuntu-trusty/main amd64 Packages
              100 /var/lib/dpkg/status
           1.12.1-0~trusty 0
              500 https://apt.dockerproject.org/repo/ ubuntu-trusty/main amd64 Packages
           1.12.0-0~trusty 0
              500 https://apt.dockerproject.org/repo/ ubuntu-trusty/main amd64 Packages
    ```
From now on when you run `apt-get upgrade`, `APT` pulls from the new repository.

### Prerequisites by Ubuntu Version

#### Ubuntu Xenial 16.04 (LTS), Wily 15.10, Trusty 14.04 (LTS)

For Ubuntu Trusty, Wily, and Xenial, install the `linux-image-extra-*` kernel
packages, which allows you use the `aufs` storage driver.

To install the `linux-image-extra-*` packages:

1.  Open a terminal on your Ubuntu host.

2.  Update your package manager.

    ```bash
    $ sudo apt-get update
    ```

3.  Install the recommended packages.

    ```bash
    $ sudo apt-get install linux-image-extra-$(uname -r) linux-image-extra-virtual
    ```

4.  Go ahead and [install Docker](ubuntulinux.md#install-the-latest-version).

#### Ubuntu Precise 12.04 (LTS)

For Ubuntu Precise, Docker requires the 3.13 kernel version. If your kernel
version is older than 3.13, you must upgrade it. Refer to this table to see
which packages are required for your environment:

| Package                           | Description |
| --------------------------------- | ----------- |
| `linux-image-generic-lts-trusty`  | Generic Linux kernel image. This kernel has AUFS built in. This is required to run Docker. |
| `linux-headers-generic-lts-trusty`| Allows packages such as ZFS and VirtualBox guest additions which depend on them. If you didn't install the headers for your existing kernel, then you can skip these headers for the"trusty" kernel. If you're unsure, you should include this package for safety. |
| `xserver-xorg-lts-trusty`         | Optional in non-graphical environments without Unity/Xorg. **Required** when running Docker on machine with a graphical environment. |
| `ligbl1-mesa-glx-lts-trusty`      | To learn more about the reasons for these packages, read the installation instructions for backported kernels, specifically the [LTS Enablement Stack](https://wiki.ubuntu.com/Kernel/LTSEnablementStack). Refer to note 5 under each version. |


To upgrade your kernel and install the additional packages, do the following:

1.  Open a terminal on your Ubuntu host.

2.  Update your package manager.

    ```bash
    $ sudo apt-get update
    ```

3.  Install both the required and optional packages.

    ```bash
    $ sudo apt-get install linux-image-generic-lts-trusty
    ```

    Repeat this step for other packages you need to install.

4.  Reboot your host to use the updated kernel.

    ```bash
    $ sudo reboot
    ```

5.  After your system reboots, go ahead and
    [install Docker](ubuntulinux.md#install-the-latest-version).

## Install the latest version

Make sure you have satisfied all the
[prerequisites](ubuntulinux.md#prerequisites), then follow these steps.

>**Note**: For production systems, it is recommended that you
[install a specific version](ubuntulinux.md#install-a-specific-version) so that
you do not accidentally update Docker. You should plan upgrades for production
systems carefully.

1.  Log into your Ubuntu installation as a user with `sudo` privileges.

2.  Update your `APT` package index.

    ```bash
    $ sudo apt-get update
    ```
3.  Install Docker.

    ```bash
    $ sudo apt-get install docker-engine
    ```

4.  Start the `docker` daemon.

    ```bash
    $ sudo service docker start
    ```

5.  Verify that `docker` is installed correctly by running the `hello-world`
    image.

    ```bash
    $ sudo docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

## Install a specific version

To install a specific version of `docker-engine`:

1.  List all available versions using `apt-cache madison`:

    ```bash
    $ apt-cache madison docker-engine

    docker-engine | 1.12.3-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.12.2-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.12.1-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.12.0-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.11.2-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.11.1-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    docker-engine | 1.11.0-0~xenial | https://apt.dockerproject.org/repo ubuntu-xenial/main amd64 Packages
    ```

2.  The second field is the version string. To install exactly `1.12.0-0~xenial`,
    append it after the package name in the `apt-get install` command, separated
    from the package name by an equals sign (`=`).
    ```bash
    $ sudo apt-get install docker-engine=1.12.0-0~xenial
    ```

    If you already have a newer version installed, you will be prompted to
    downgrade Docker. Otherwise, the specific version will be installed.

3.  Follow steps 4 and 5 of
    [Install the latest version](ubuntulinux.md#install-the-latest-version).

## Install a pre-release version

If you want to test Docker on Ubuntu, on a non-production system, follow these
steps. To install a stable released version of Docker afterward, you will need
to revert to the previous configuration.

1.  Edit `/etc/apt/sources.list.d/docker.list`.

    ```bash
    $ sudo nano /etc/apt/sources.list.d/docker.list
    ```

    Change `main` to `testing` at the end of the top line. Save and close the
    file.

2.  Update the package list.

    ```bash
    $ sudo apt-get update
    ```

3.  List the available testing versions.

    ```bash
    $ sudo apt-cache madison docker-engine
    ```

4.  Install a specific version following the same procedure as
    [Install a specific version](ubuntulinux.md#install-a-specific-version).

## Optional configurations

This section contains optional procedures for configuring Ubuntu to work better
with Docker.

* [Manage Docker as a non-root user](ubuntulinux.md#manage-docker-as-a-non-root-user)
* [Adjust memory and swap accounting](ubuntulinux.md#adjust-memory-and-swap-accounting)
* [Enable UFW forwarding](ubuntulinux.md#enable-ufw-forwarding)
* [Configure a DNS server for use by Docker](ubuntulinux.md#configure-a-dns-server-for-use-by-docker)
* [Configure Docker to start on boot](ubuntulinux.md#configure-docker-to-start-on-boot)

### Manage Docker as a non-root user

The `docker` daemon binds to a Unix socket instead of a TCP port. By default
that Unix socket is owned by the user `root` and other users can only access it
using `sudo`. The `docker` daemon always runs as the `root` user.

If you don't want to use `sudo` when you use the `docker` command, create a Unix
group called `docker` and add users to it. When the `docker` daemon starts, it
makes the ownership of the Unix socket read/writable by the `docker` group.

>**Warning**: The `docker` group is equivalent to the `root` user. For details
on how this impacts security in your system, see [*Docker Daemon Attack
Surface*](../../security/security.md#docker-daemon-attack-surface).

To create the `docker` group and add your user:

1.  Log into Ubuntu as a user with `sudo` privileges.

2.  Create the `docker` group.
    
    ```bash
    $ sudo groupadd docker
    ```

3.  Add your user to the `docker` group.

    ```bash
    $ sudo usermod -aG docker $USER
    ```

4.  Log out and log back in so that your group membership is re-evaluated.

5.  Verify that you can `docker` commands without `sudo`.

    ```bash
    $ docker run hello-world
    ```

	  If this fails, you will see an error:

    ```none
		Cannot connect to the Docker daemon. Is 'docker daemon' running on this host?
    ```

	  Check whether the `DOCKER_HOST` environment variable is set for your shell.

    ```bash
    $ env | grep DOCKER_HOST
    ```

	  If it is set, the above command will return a result. If so, unset it.

    ```bash
    $ unset DOCKER_HOST
    ```

    You may need to edit your environment in files such as `~/.bashrc` or
    `~/.profile` to prevent the `DOCKER_HOST` variable from being set
    erroneously.

### Enable memory and swap accounting

You may see messages similar to the following when working with an image:

```none
WARNING: Your kernel does not support cgroup swap limit. WARNING: Your
kernel does not support swap limit capabilities. Limitation discarded.
```

If you don't care about these capabilities, you can ignore the warning. You can
enable these capabilities in your kernel by following these instructions. Memory
and swap accounting incur an overhead of about 1% of the total available
memory and a 10% overall performance degradation, even if Docker is not running.

1.  Log into Ubuntu as a user with `sudo` privileges.

2.  Edit the `/etc/default/grub` file.

3.  Add or edit the `GRUB_CMDLINE_LINUX` line to add the following two key-value
    pairs:

    ```none
    GRUB_CMDLINE_LINUX="cgroup_enable=memory swapaccount=1"
    ```

    Save and close the file.

4.  Update GRUB.

    ```bash
    $ sudo update-grub
    ```

     If your GRUB configuration file has incorrect syntax, an error will occur.
     In this case, steps 3 and 4.

6.  Reboot your system. Memory and swap accounting are enabled and the warning
    does not occur.


### Enable UFW forwarding

If you use [UFW (Uncomplicated Firewall)](https://help.ubuntu.com/community/UFW)
on the same host as you run Docker, you'll need to do additional configuration.
Docker uses a bridge to manage container networking. By default, UFW drops all
forwarding traffic. You must set UFW's forwarding policy appropriately.

In addition, UFW blocks all incoming traffic by default. If you want to access
the Docker Remote API from another host and you have enabled remote access, you
need to configure UFW to allow incoming connections on the Docker port, which
defaults to `2376` if TLS encrypted transport is enabled or `2375` otherwise. By
default, Docker runs **without** TLS enabled. If you do not use TLS, you are
strongly discouraged from allowing access to the Docker Remote API from remote
hosts, to prevent remote privilege-escalation attacks.

To configure UFW and allow incoming connections on the Docker port:

1.  Log into Ubuntu as a user with `sudo` privileges.

2.  Verify that UFW is enabled.

    ```bash
    $ sudo ufw status
    ```

    If `ufw` is not enabled, the remaining steps will not be helpful.

3.  Edit the UFW configuration file, which is usually `/etc/default/ufw` or
`/etc/sysconfig/ufw`. Set the `DEFAULT_FORWARD_POLICY` policy to `ACCEPT`.

    ```none
    DEFAULT_FORWARD_POLICY="ACCEPT"
    ```

    Save and close the file.

4.  If you need to enable access to the Docker Remote API from external hosts
    and understand the security implications (see the section before this
    procedure), then configure UFW to allow incoming connections on the Docker port,
    which is 2375 if you do not use TLS, and 2376 if you do.

    ```bash
    $ sudo ufw allow 2376/tcp
    ```

5.  Reload UFW.
    ```bash
    $ sudo ufw reload
    ```

### Configure a DNS server for use by Docker

Ubuntu systems which use `networkmanager` use a `dnsmasq` instance that runs on
a loopback address such as `127.0.0.1` or `127.0.1.1` and adds this entry to
`/etc/resolv.conf`. The `dnsmasq` service provides a local DNS cache to speed up
DNS look-ups and also provides DHCP services. This configuration will not work
within a Docker container which has its own network namespace. This is because
the Docker container resolves loopback addresses such as `127.0.0.1` to itself,
and it is very unlikely to be running a DNS server on its own loopback address.

If Docker detects that no DNS server referenced in `/etc/resolv.conf` is a fully
functional DNS server, the following warning occurs and Docker uses the public
DNS servers provided by Google at `8.8.8.8` and `8.8.4.4` for DNS resolution.

```none
WARNING: Local (127.0.0.1) DNS resolver found in resolv.conf and containers
can't use it. Using default external servers : [8.8.8.8 8.8.4.4]
```

If you don't use `dnsmasq` or NetworkManager or have never seen this warning,
you can skip the rest of this section. To see if you use `dnsmasq`, use the
following command:

```bash
$ ps aux |grep dnsmasq
```

If this warning occurs and cannot use the public nameservers, such as when you
run a DNS server which resolves hostnames on your internal network, you have
two choices:

- You can specify a DNS server for Docker to use.
- You can disable `dnsmasq` in NetworkManager. If you do this, NetworkManager
  will add your true DNS nameserver to `/etc/resolv.conf`, but you will lose the
  possible benefits of `dnsmasq`.

**You only need to use one of these methods.**

#### Specify DNS servers for Docker

The instructions below work whether your Ubuntu installation uses `upstart` or
`systemd`.

The default location of the configuration file is `/etc/docker/daemon.json`. You
can change the location of the configuration file using the `--config-file`
daemon flag. The documentation below assumes the configuration file is located
at `/etc/docker/daemon.json`.

1.  Log into Ubuntu as a user with `sudo` privileges.

2.  Create or edit the Docker daemon configuration file, which defaults to
    `/etc/docker/daemon.json` file, which controls the Docker daemon
    configuration.

    ```bash
    sudo nano /etc/docker/daemon.json
    ```

2.  Add a `dns` key with one or more IP addresses as values. If the file has
    existing contents, you only need to add or edit the `dns` line.
    
    ```json
    {
    	"dns": ["8.8.8.8", "8.8.4.4"]
    }
    ```

    If your internal DNS server cannot resolve public IP addresses, include at
    least one DNS server which can, so that you can connect to Docker Hub and so
    that your containers can resolve internet domain names.

    Save and close the file.

3.  Restart the Docker daemon.

    ```bash
    $ sudo service docker restart
    ```

4.  Verify that Docker can resolve external IP addresses by trying to pull an
    image:

    ```bash
    $ docker pull hello-world
    ```

5.  If necessary, verify that Docker containers can resolve an internal hostname
    by pinging it.

    ```bash
    $ docker run --rm -it alpine ping -c4 my_internal_host

    PING google.com (192.168.1.2): 56 data bytes
    64 bytes from 192.168.1.2: seq=0 ttl=41 time=7.597 ms
    64 bytes from 192.168.1.2: seq=1 ttl=41 time=7.635 ms
    64 bytes from 192.168.1.2: seq=2 ttl=41 time=7.660 ms
    64 bytes from 192.168.1.2: seq=3 ttl=41 time=7.677 ms
    ```

#### Disable `dnsmasq` in NetworkManager

If you prefer not to change the Docker daemon's configuration to use a specific
IP address, follow these instructions to disable `dnsmasq` in NetworkManager.

1.  Edit the `/etc/NetworkManager/NetworkManager.conf` file.

2.  Comment out the `dns=dnsmasq` line by adding a `#` character to the beginning
    of the line.

    ```none
    # dns=dnsmasq
    ```

    Save and close the file.

4.  Restart both NetworkManager and Docker. As an alternative, you can reboot
    your system.

    ```bash
    $ sudo restart network-manager
    $ sudo restart docker
    ```

### Configure Docker to start on boot

Ubuntu uses `systemd` as its boot and service manager `15.04` onwards and `upstart`
for versions `14.10` and below.

#### `systemd`

```bash
$ sudo systemctl enable docker
```

#### `upstart`

For `14.10` and below, Docker is automatically configured to start on boot using
`upstart`.

## Upgrade Docker

To install the latest version of Docker with `apt-get`. The following example
fetches information about available versions of all system packages, then
updates Docker if a new version is available.

```bash
$ sudo apt-get update
$ sudo apt-get install docker-engine
```

## Uninstallation

To uninstall the Docker package:

```bash
$ sudo apt-get purge docker-engine
```

To uninstall the Docker package and dependencies that are no longer needed:

```bash
$ sudo apt-get autoremove --purge docker-engine
```

Images, containers, volumes, or customized configuration files on your host are
not automatically removed. To delete all images, containers, and volumes run the
following command:

```bash
$ rm -rf /var/lib/docker
```

You must delete any edited configuration files manually.
