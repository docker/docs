---
description: Optional post-installation steps for Linux
keywords: Docker, Docker documentation, requirements, apt, installation, ubuntu, install, uninstall, upgrade, update
title: Post-installation steps for Linux
---


This section contains optional procedures for configuring Linux hosts to work
better with Docker.

## Manage Docker as a non-root user

The `docker` daemon binds to a Unix socket instead of a TCP port. By default
that Unix socket is owned by the user `root` and other users can only access it
using `sudo`. The `docker` daemon always runs as the `root` user.

If you don't want to use `sudo` when you use the `docker` command, create a Unix
group called `docker` and add users to it. When the `docker` daemon starts, it
makes the ownership of the Unix socket read/writable by the `docker` group.

> **Warning**: The `docker` group grants privileges equivalent to the `root`
> user. For details on how this impacts security in your system, see
> [*Docker Daemon Attack Surface*](/engine/security/security.md#docker-daemon-attack-surface).

To create the `docker` group and add your user:

1.  Create the `docker` group.

    ```bash
    $ sudo groupadd docker
    ```

2.  Add your user to the `docker` group.

    ```bash
    $ sudo usermod -aG docker $USER
    ```

3.  Log out and log back in so that your group membership is re-evaluated.

4.  Verify that you can `docker` commands without `sudo`.

    ```bash
    $ docker run hello-world
    ```

    This command downloads a test image and runs it in a container. When the
    container runs, it prints an informational message and exits.

## Configure Docker to start on boot

Most current Linux distributions (RHEL, CentOS, Fedora, Ubuntu 16.04 and higher)
use [`systemd`](#systemd) to manage which services start when the system boots.
Ubuntu 14.10 and below use [`upstart`](#upstart). Oracle Linux 6 uses
`chkconfig`.

### `systemd`

```bash
$ sudo systemctl enable docker
```

To disable this behavior, use `disable` instead.

```bash
$ sudo systemctl disable docker
```

If you need to add an HTTP Proxy, set a different directory or partition for the
Docker runtime files, or make other customizations, see
[customize your systemd Docker daemon options](/engine/admin/systemd.md).

### `upstart`

Docker is automatically configured to start on boot using
`upstart`. To disable this behavior, use the following command:

```bash
$ echo manual | sudo tee /etc/init/docker.override
```

### `chkconfig`

```bash
$ sudo chkconfig docker on
```

## Use a different storage engine

For information about the different storage engines, see
[Storage drivers](/engine/userguide/storagedriver/imagesandcontainers.md).
The default storage engine and the list of supported storage engines depend on
your host's Linux distribution and available kernel drivers.

## Troubleshooting

### `Cannot connect to the Docker daemon`

If you see an error such as the following, your Docker client may be configured
to connect to a Docker daemon on a different host, and that host may not be
reachable.

```none
Cannot connect to the Docker daemon. Is 'docker daemon' running on this host?
```

To see which host your client is configured to connect to, check the value of
the `DOCKER_HOST` variable in your environment.

```bash
$ env | grep DOCKER_HOST
```

If this command returns a value, the Docker client is set to connect to a
Docker daemon running on that host. If it is unset, the Docker client is set to
connect to the Docker daemon running on the local host. If it is set in error,
use the following command to unset it:

```bash
$ unset DOCKER_HOST
```

You may need to edit your environment in files such as `~/.bashrc` or
`~/.profile` to prevent the `DOCKER_HOST` variable from being set
erroneously.

If `DOCKER_HOST` is set as intended, verify that the Docker daemon is running
on the remote host and that a firewall or network outage is not preventing you
from connecting.

### IP forwarding problems

If you manually configure your network using `systemd-network` with `systemd`
version 219 or higher, Docker containers may be unable to access your network.
Beginning with `systemd` version 220, the forwarding setting for a given network
(`net.ipv4.conf.<interface>.forwarding`) defaults to *off*. This setting
prevents IP forwarding. It also conflicts with Docker's behavior of enabling
the `net.ipv4.conf.all.forwarding` setting within containers.

To work around this on RHEL, CentOS, or Fedora, edit the `<interface>.network`
file in `/usr/lib/systemd/network/` on your Docker host
(ex: `/usr/lib/systemd/network/80-container-host0.network`) and add the
following block within the `[Network]` section.

```
[Network]
...
IPForward=kernel
# OR
IPForward=true
...
```

This configuration allows IP forwarding from the container as expected.


### `DNS resolver found in resolv.conf and containers can't use it`

Linux systems which use a GUI often have a network manager running, which uses a
`dnsmasq` instance running on a loopback address such as `127.0.0.1` or
`127.0.1.1` to cache DNS requests, and adds this entry to
`/etc/resolv.conf`. The `dnsmasq` service speeds up
DNS look-ups and also provides DHCP services. This configuration will not work
within a Docker container which has its own network namespace, because
the Docker container resolves loopback addresses such as `127.0.0.1` to
**itself**, and it is very unlikely to be running a DNS server on its own
loopback address.

If Docker detects that no DNS server referenced in `/etc/resolv.conf` is a fully
functional DNS server, the following warning occurs and Docker uses the public
DNS servers provided by Google at `8.8.8.8` and `8.8.4.4` for DNS resolution.

```none
WARNING: Local (127.0.0.1) DNS resolver found in resolv.conf and containers
can't use it. Using default external servers : [8.8.8.8 8.8.4.4]
```

If you see this warning, first check to see if you use `dnsmasq`:

```bash
$ ps aux |grep dnsmasq
```

If your container needs to resolve hosts which are internal to your network, the
public nameservers will not be adequate. You have two choices:

- You can specify a DNS server for Docker to use, **or**
- You can disable `dnsmasq` in NetworkManager. If you do this, NetworkManager
  will add your true DNS nameserver to `/etc/resolv.conf`, but you will lose the
  possible benefits of `dnsmasq`.

**You only need to use one of these methods.**

### Specify DNS servers for Docker

The default location of the configuration file is `/etc/docker/daemon.json`. You
can change the location of the configuration file using the `--config-file`
daemon flag. The documentation below assumes the configuration file is located
at `/etc/docker/daemon.json`.

1. .  Create or edit the Docker daemon configuration file, which defaults to
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
    $ docker run --rm -it alpine ping -c4 <my_internal_host>

    PING google.com (192.168.1.2): 56 data bytes
    64 bytes from 192.168.1.2: seq=0 ttl=41 time=7.597 ms
    64 bytes from 192.168.1.2: seq=1 ttl=41 time=7.635 ms
    64 bytes from 192.168.1.2: seq=2 ttl=41 time=7.660 ms
    64 bytes from 192.168.1.2: seq=3 ttl=41 time=7.677 ms
    ```

#### Disable `dnsmasq`

##### Ubuntu

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

##### RHEL, CentOS, or Fedora

To disable `dnsmasq` on RHEL, CentOS, or Fedora:

1.  Disable the `dnsmasq` service:

    ```bash
    $ sudo service dnsmasq stop

    $ sudo systemctl disable dnsmasq
    ```

2.  Configure the DNS servers manually using the
    [Red Hat documentation](https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/6/html/Deployment_Guide/s1-networkscripts-interfaces.html){: target="_blank" class="_"}.

### Allow access to the remote API through a firewall

If you run a firewall on the same host as you run Docker and you want to access
the Docker Remote API from another host and remote access is enabled, you need
to configure your firewall to allow incoming connections on the Docker port,
which defaults to `2376` if TLS encrypted transport is enabled or `2375`
otherwise.

#### Specific instructions for UFW

[UFW (Uncomplicated Firewall)](https://help.ubuntu.com/community/UFW) drops all
forwarding traffic and all incoming traffic by default. If you want to access
the Docker Remote API from another host and you have enabled remote access, you
need to configure UFW to allow incoming connections on the Docker port, which
defaults to `2376` if TLS encrypted transport is enabled or `2375` otherwise. By
default, Docker runs **without** TLS enabled. If you do not use TLS, you are
strongly discouraged from allowing access to the Docker Remote API from remote
hosts, to prevent remote privilege-escalation attacks.

To configure UFW and allow incoming connections on the Docker port:

1.  Verify that UFW is enabled.

    ```bash
    $ sudo ufw status
    ```

    If `ufw` is not enabled, the remaining steps will not be helpful.

2.  Edit the UFW configuration file, which is usually `/etc/default/ufw` or
    `/etc/sysconfig/ufw`. Set the `DEFAULT_FORWARD_POLICY` policy to `ACCEPT`.

    ```none
    DEFAULT_FORWARD_POLICY="ACCEPT"
    ```

    Save and close the file.

3.  If you need to enable access to the Docker Remote API from external hosts
    and understand the security implications (see the section before this
    procedure), then configure UFW to allow incoming connections on the Docker port,
    which is 2375 if you do not use TLS, and 2376 if you do.

    ```bash
    $ sudo ufw allow 2376/tcp
    ```

4.  Reload UFW.

    ```bash
    $ sudo ufw reload
    ```

### `Your kernel does not support cgroup swap limit capabilities`

You may see messages similar to the following when working with an image:

```none
WARNING: Your kernel does not support swap limit capabilities. Limitation discarded.
```

If you don't need these capabilities, you can ignore the warning. You can
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

## Next steps

- Continue with the [User Guide](/engine/userguide/index.md).
