---
description: Optional post-installation steps for Linux
keywords: >
  Docker, Docker documentation, requirements, apt, installation, ubuntu,
  install, uninstall, upgrade, update
title: Post-installation steps for Linux
redirect_from:
  - /engine/installation/linux/docker-ee/linux-postinstall/
  - /engine/installation/linux/linux-postinstall/
  - /install/linux/linux-postinstall/
---

These optional post-installation procedures shows you how to configure your
Linux host machine to work better with Docker.

## Manage Docker as a non-root user

The Docker daemon binds to a Unix socket, not a TCP port. By default it's the
`root` user that owns the Unix socket, and other users can only access it using
`sudo`. The Docker daemon always runs as the `root` user.

If you don't want to preface the `docker` command with `sudo`, create a Unix
group called `docker` and add users to it. When the Docker daemon starts, it
creates a Unix socket accessible by members of the `docker` group.

<!-- prettier-ignore -->
> *Warning**
>
> The `docker` group grants root-level privileges to the user. For
> details on how this impacts security in your system, see
> [Docker Daemon Attack Surface](../security/index.md#docker-daemon-attack-surface).
{: .warning}

> **Note**
>
> To run Docker without root privileges, see
> [Run the Docker daemon as a non-root user (Rootless mode)](../security/rootless.md).

To create the `docker` group and add your user:

1. Create the `docker` group.

   ```console
   $ sudo groupadd docker
   ```

2. Add your user to the `docker` group.

   ```console
   $ sudo usermod -aG docker $USER
   ```

3. Log out and log back in so that your group membership is re-evaluated.

   > If you're running Linux in a virtual machine, it may be necessary to
   > restart the virtual machine for changes to take effect.

   You can also run the following command to activate the changes to groups:

   ```console
   $ newgrp docker
   ```

4. Verify that you can run `docker` commands without `sudo`.

   ```console
   $ docker run hello-world
   ```

   This command downloads a test image and runs it in a container. When the
   container runs, it prints a message and exits.

   If you initially ran Docker CLI commands using `sudo` before adding your user
   to the `docker` group, you may see the following error:

   ```none
   WARNING: Error loading config file: /home/user/.docker/config.json -
   stat /home/user/.docker/config.json: permission denied
   ```

   This error indicates that the permission settings for the `~/.docker/`
   directory are incorrect, due to having used the `sudo` command earlier.

   To fix this problem, either remove the `~/.docker/` directory (it's recreated
   automatically, but any custom settings are lost), or change its ownership and
   permissions using the following commands:

   ```console
   $ sudo chown "$USER":"$USER" /home/"$USER"/.docker -R
   $ sudo chmod g+rwx "$HOME/.docker" -R
   ```

## Configure Docker to start on boot

Many modern Linux distributions use [systemd](../../config/daemon/systemd.md) to
manage which services start when the system boots. On Debian and Ubuntu, the
Docker starts on boot by default. To automatically start Docker and containerd
on boot for other Linux distributions using systemd, run the following commands:

```console
$ sudo systemctl enable docker.service
$ sudo systemctl enable containerd.service
```

To stop this behavior, use `disable` instead.

```console
$ sudo systemctl disable docker.service
$ sudo systemctl disable containerd.service
```

If you need to add an HTTP proxy, set a different directory or partition for the
Docker runtime files, or make other customizations, see
[customize your systemd Docker daemon options](../../config/daemon/systemd.md).

## Use a different storage engine

For information about the different storage engines, see
[Storage drivers](../../storage/storagedriver/index.md). The default storage
engine and the list of supported storage engines depend on your host's Linux
distribution and available kernel drivers.

## Configure default logging driver

Docker provides the [capability](../../config/containers/logging/index.md) to
collect and view log data from all containers running on a host via a series of
logging drivers. The default logging driver, `json-file`, writes log data to
JSON-formatted files on the host filesystem. Over time, these log files expand
in size, leading to potential exhaustion of disk resources.

To avoid issues with overusing disk for log data, consider one of the following
options:

- Configure the `json-file` logging driver to turn on
  [log rotation](../../config/containers/logging/json-file.md)
- Use an
  [alternative logging driver](../../config/containers/logging/configure.md#configure-the-default-logging-driver)
  such as the ["local" logging driver](../../config/containers/logging/local.md)
  that performs log rotation by default
- Use a logging driver that sends logs to a remote logging aggregator.

## Configure where the Docker daemon listens for connections

By default, the Docker daemon listens for connections on a Unix socket to accept
requests from local clients. It's possible to allow Docker to accept requests
from remote hosts by configuring it to listen on an IP address and port as well
as the Unix socket. For more detailed information on this configuration option,
refer to the
[dockerd CLI reference](/engine/reference/commandline/dockerd/#bind-docker-to-another-hostport-or-a-unix-socket).

<!-- prettier-ignore -->
> Secure your connection
>
> Before configuring Docker to accept connections from remote hosts it's
> critically important that you understand the security implications of opening
> Docker to the network. If steps aren't taken to secure the connection, it's
> possible for remote non-root users to gain root access on the host. For more
> information on how to use TLS certificates to secure this connection, check
> [Protect the Docker daemon socket](../security/protect-access.md).
{: .warning}

You can configure Docker to accept remote connections can be done using the
`docker.service` systemd unit file for Linux distributions using systemd. Or you
can use the `daemon.json` file, if your distribution doesn't use systemd.

> systemd vs `daemon.json`
>
> Configuring Docker to listen for connections using both the `systemd` unit
> file and the `daemon.json` file causes a conflict that prevents Docker from
> starting.

### Configuring remote access with `systemd` unit file

1. Use the command `sudo systemctl edit docker.service` to open an override file
   for `docker.service` in a text editor.

2. Add or modify the following lines, substituting your own values.

   ```systemd
   [Service]
   ExecStart=
   ExecStart=/usr/bin/dockerd -H fd:// -H tcp://127.0.0.1:2375
   ```

3. Save the file.

4. Reload the `systemctl` configuration.

   ```console
    $ sudo systemctl daemon-reload
   ```

5. Restart Docker.

   ```console
   $ sudo systemctl restart docker.service
   ```

6. Verify that the change has gone through.

   ```console
   $ sudo netstat -lntp | grep dockerd
   tcp        0      0 127.0.0.1:2375          0.0.0.0:*               LISTEN      3758/dockerd
   ```

### Configuring remote access with `daemon.json`

1. Set the `hosts` array in the `/etc/docker/daemon.json` to connect to the UNIX
   socket and an IP address, as follows:

   ```json
   {
     "hosts": ["unix:///var/run/docker.sock", "tcp://127.0.0.1:2375"]
   }
   ```

2. Restart Docker.

3. Verify that the change has gone through.

   ```console
   $ sudo netstat -lntp | grep dockerd
   tcp        0      0 127.0.0.1:2375          0.0.0.0:*               LISTEN      3758/dockerd
   ```

## Enable IPv6 on the Docker daemon

To enable IPv6 on the Docker daemon, see
[Enable IPv6 support](../../config/daemon/ipv6.md).

## Troubleshooting

### Kernel compatibility

Docker can't run correctly if your kernel is older than version 3.10, or if it's
missing kernel modules. To check kernel compatibility, you can download and run
the
[`check-config.sh`](https://raw.githubusercontent.com/docker/docker/master/contrib/check-config.sh)
script.

```console
$ curl https://raw.githubusercontent.com/docker/docker/master/contrib/check-config.sh > check-config.sh

$ bash ./check-config.sh
```

The script only works on Linux.

### `Cannot connect to the Docker daemon`

```none
Cannot connect to the Docker daemon. Is 'docker daemon' running on this host?
```

This error may indicate:

- The Docker daemon isn't running on your system. Start the daemon and try
  running the command again.
- Your Docker client is attempting to connect to a Docker daemon on a different
  host, and that host is unreachable.

To see which host your client is connecting to, check the value of the
`DOCKER_HOST` variable in your environment.

```console
$ env | grep DOCKER_HOST
```

If this command returns a value, the Docker client is set to connect to a Docker
daemon running on that host. If it's unset, the Docker client is set to connect
to the Docker daemon running on the local host. If it's set in error, use the
following command to unset it:

```console
$ unset DOCKER_HOST
```

You may need to edit your environment in files such as `~/.bashrc` or
`~/.profile` to prevent the `DOCKER_HOST` variable from being set erroneously.

If `DOCKER_HOST` is set as intended, verify that the Docker daemon is running on
the remote host and that a firewall or network outage isn't preventing you from
connecting.

### IP forwarding problems

If you manually configure your network using `systemd-network` with `systemd`
version 219 or later, Docker containers may not be able to access your network.
Beginning with `systemd` version 220, the forwarding setting for a given network
(`net.ipv4.conf.<interface>.forwarding`) defaults to _off_. This setting
prevents IP forwarding. It also conflicts with Docker's behavior of enabling the
`net.ipv4.conf.all.forwarding` setting within containers.

To work around this on RHEL, CentOS, or Fedora, edit the `<interface>.network`
file in `/usr/lib/systemd/network/` on your Docker host (ex:
`/usr/lib/systemd/network/80-container-host0.network`) and add the following
block within the `[Network]` section.

```systemd
[Network]
...
IPForward=kernel
# OR
IPForward=true
```

This configuration allows IP forwarding from the container as expected.

### `DNS resolver found in resolv.conf and containers can't use it`

Linux desktop environments often have a network manager program running, that
uses `dnsmasq` to cache DNS requests by adding them to `/etc/resolv.conf`. The
`dnsmasq` instance runs on a loopback address such as `127.0.0.1` or
`127.0.1.1`. It speeds up DNS look-ups and provides DHCP services. Such a
configuration doesn't work within a Docker container. The Docker container uses
its own network namespace, and resolves loopback addresses such as `127.0.0.1`
to itself, and it's unlikely to be running a DNS server on its own loopback
address.

If Docker detects that no DNS server referenced in `/etc/resolv.conf` is a fully
functional DNS server, the following warning occurs:

```none
WARNING: Local (127.0.0.1) DNS resolver found in resolv.conf and containers
can't use it. Using default external servers : [8.8.8.8 8.8.4.4]
```

If you see this warning, first check to see if you use `dnsmasq`:

```console
$ ps aux | grep dnsmasq
```

If your container needs to resolve hosts which are internal to your network, the
public nameservers aren't adequate. You have two choices:

- Specify DNS servers for Docker to use.
- Turn off `dnsmasq`.

  Turning off `dnsmasq` adds the IP addresses of actual DNS nameserver to
  `/etc/resolv.conf`, and you lose the benefits of `dnsmasq`.

You only need to use one of these methods.

### Specify DNS servers for Docker

The default location of the configuration file is `/etc/docker/daemon.json`. You
can change the location of the configuration file using the `--config-file`
daemon flag. The following instruction assumes that the location of the
configuration file is `/etc/docker/daemon.json`.

1. Create or edit the Docker daemon configuration file, which defaults to
   `/etc/docker/daemon.json` file, which controls the Docker daemon
   configuration.

   ```console
   $ sudo nano /etc/docker/daemon.json
   ```

2. Add a `dns` key with one or more DNS server IP addresses as values.

   ```json
   {
     "dns": ["8.8.8.8", "8.8.4.4"]
   }
   ```

   If the file has existing contents, you only need to add or edit the `dns`
   line. If your internal DNS server can't resolve public IP addresses, include
   at least one DNS server that can. Doing so allows you to connect to Docker
   Hub, and your containers to resolve internet domain names.

   Save and close the file.

3. Restart the Docker daemon.

   ```console
   $ sudo service docker restart
   ```

4. Verify that Docker can resolve external IP addresses by trying to pull an
   image:

   ```console
   $ docker pull hello-world
   ```

5. If necessary, verify that Docker containers can resolve an internal hostname
   by pinging it.

   ```console
   $ docker run --rm -it alpine ping -c4 <my_internal_host>

   PING google.com (192.168.1.2): 56 data bytes
   64 bytes from 192.168.1.2: seq=0 ttl=41 time=7.597 ms
   64 bytes from 192.168.1.2: seq=1 ttl=41 time=7.635 ms
   64 bytes from 192.168.1.2: seq=2 ttl=41 time=7.660 ms
   64 bytes from 192.168.1.2: seq=3 ttl=41 time=7.677 ms
   ```

#### Turn off `dnsmasq`

##### Ubuntu

If you prefer not to change the Docker daemon's configuration to use a specific
IP address, follow these instructions to turn off `dnsmasq` in NetworkManager.

1. Edit the `/etc/NetworkManager/NetworkManager.conf` file.

2. Comment out the `dns=dnsmasq` line by adding a `#` character to the beginning
   of the line.

   ```none
   # dns=dnsmasq
   ```

   Save and close the file.

3. Restart both NetworkManager and Docker. As an alternative, you can reboot
   your system.

   ```console
   $ sudo systemctl restart network-manager
   $ sudo systemctl restart docker
   ```

##### RHEL, CentOS, or Fedora

To turn off `dnsmasq` on RHEL, CentOS, or Fedora:

1. Turn off the `dnsmasq` service:

   ```console
   $ sudo systemctl stop dnsmasq
   $ sudo systemctl disable dnsmasq
   ```

2. Configure the DNS servers manually using the
   [Red Hat documentation](https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/6/html/Deployment_Guide/s1-networkscripts-interfaces.html){:
   target="_blank" rel="noopener" class="_"}.

### Allow access to the remote API through a firewall

If you run a firewall on the same host as you run Docker, and you want to access
the Docker Remote API from another remote host, you must configure your firewall
to allow incoming connections on the Docker port. The default port is `2376` if
you're using TLS encrypted transport, or `2375` otherwise.

Two common firewall daemons are
[Uncomplicated Firewall (UFW)](https://help.ubuntu.com/community/UFW) (often
used for Ubuntu systems) and [firewalld](https://firewalld.org) (often used for
RPM-based systems). Consult the documentation for your OS and firewall, but the
following information might help you get started. These options are permissive
and you may want to use a different configuration that locks your system down
more.

- UFW: Set `DEFAULT_FORWARD_POLICY="ACCEPT"` in your configuration.

- firewalld: Add rules similar to the following to your policy (one for incoming
  requests and one for outgoing requests).

  ```xml
  <direct>
    [ <rule ipv="ipv6" table="filter" chain="FORWARD_direct" priority="0"> -i zt0 -j ACCEPT </rule> ]
    [ <rule ipv="ipv6" table="filter" chain="FORWARD_direct" priority="0"> -o zt0 -j ACCEPT </rule> ]
  </direct>
  ```

  Double check that the interface names and chain names are correct.

### `Your kernel does not support cgroup swap limit capabilities`

On Ubuntu or Debian hosts, You may see messages similar to the following when
working with an image.

```none
WARNING: Your kernel does not support swap limit capabilities. Limitation discarded.
```

If you don't need these capabilities, you can ignore the warning.

You can turn on these capabilities on Ubuntu or Debian by following these
instructions. Memory and swap accounting incur an overhead of about 1% of the
total available memory and a 10% overall performance degradation, even when
Docker isn't running.

1. Log into the Ubuntu or Debian host as a user with `sudo` privileges.

2. Edit the `/etc/default/grub` file. Add or edit the `GRUB_CMDLINE_LINUX` line
   to add the following two key-value pairs:

   ```none
   GRUB_CMDLINE_LINUX="cgroup_enable=memory swapaccount=1"
   ```

   Save and close the file.

3. Update the GRUB boot loader.

   ```console
   $ sudo update-grub
   ```

   An error occurs if your GRUB configuration file has incorrect syntax. In this
   case, repeat steps 2 and 3.

   The changes take effect when you reboot the system.

## Next steps

- Take a look at the [Get started](../../get-started/index.md) training modules
  to learn how to build an image and run it as a containerized application.
- Review the topics in [Develop with Docker](../../develop/index.md) to learn
  how to build new applications using Docker.
