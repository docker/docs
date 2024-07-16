---
title: Troubleshooting the Docker daemon
description: Learn how to troubleshoot errors and misconfigurations in the Docker daemon
keywords: |
  docker, daemon, configuration, troubleshooting, error, fail to start,
  networking, dns resolver, ip forwarding, dnsmasq, firewall,
  Cannot connect to the Docker daemon. Is 'docker daemon' running on this host?
aliases:
  - /engine/install/troubleshoot/
  - /storage/troubleshooting_volume_errors/
tags: [Troubleshooting]
---

This page describes how to troubleshoot and debug the daemon if you run into
issues.

You can turn on debugging on the daemon to learn about the runtime activity of
the daemon and to aid in troubleshooting. If the daemon is unresponsive, you can
also [force a full stack trace](logs.md#force-a-stack-trace-to-be-logged) of all
threads to be added to the daemon log by sending the `SIGUSR` signal to the
Docker daemon.

## Daemon

### Unable to connect to the Docker daemon

```text
Cannot connect to the Docker daemon. Is 'docker daemon' running on this host?
```

This error may indicate:

- The Docker daemon isn't running on your system. Start the daemon and try
  running the command again.
- Your Docker client is attempting to connect to a Docker daemon on a different
  host, and that host is unreachable.

### Check whether Docker is running

The operating-system independent way to check whether Docker is running is to
ask Docker, using the `docker info` command.

You can also use operating system utilities, such as
`sudo systemctl is-active docker` or `sudo status docker` or
`sudo service docker status`, or checking the service status using Windows
utilities.

Finally, you can check in the process list for the `dockerd` process, using
commands like `ps` or `top`.

#### Check which host your client is connecting to

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

### Troubleshoot conflicts between the `daemon.json` and startup scripts

If you use a `daemon.json` file and also pass options to the `dockerd` command
manually or using start-up scripts, and these options conflict, Docker fails to
start with an error such as:

```text
unable to configure the Docker daemon with file /etc/docker/daemon.json:
the following directives are specified both as a flag and in the configuration
file: hosts: (from flag: [unix:///var/run/docker.sock], from file: [tcp://127.0.0.1:2376])
```

If you see an error similar to this one and you are starting the daemon manually
with flags, you may need to adjust your flags or the `daemon.json` to remove the
conflict.

> **Note**
>
> If you see this specific error message about `hosts`, continue to the
> [next section](#configure-the-daemon-host-with-systemd)
> for a workaround.

If you are starting Docker using your operating system's init scripts, you may
need to override the defaults in these scripts in ways that are specific to the
operating system.

#### Configure the daemon host with systemd

One notable example of a configuration conflict that's difficult to
troubleshoot is when you want to specify a different daemon address from the
default. Docker listens on a socket by default. On Debian and Ubuntu systems
using `systemd`, this means that a host flag `-H` is always used when starting
`dockerd`. If you specify a `hosts` entry in the `daemon.json`, this causes a
configuration conflict and results in the Docker daemon failing to start.

To work around this problem, create a new file
`/etc/systemd/system/docker.service.d/docker.conf` with the following contents,
to remove the `-H` argument that's used when starting the daemon by default.

```systemd
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd
```

There are other times when you might need to configure `systemd` with Docker,
such as [configuring a HTTP or HTTPS proxy](./proxy.md).

> **Note**
>
> If you override this option without specifying a `hosts` entry in the
> `daemon.json` or a `-H` flag when starting Docker manually, Docker fails to
> start.

Run `sudo systemctl daemon-reload` before attempting to start Docker. If Docker
starts successfully, it's now listening on the IP address specified in the
`hosts` key of the `daemon.json` instead of a socket.

> **Important**
>
> Setting `hosts` in the `daemon.json` isn't supported on Docker
> Desktop for Windows or Docker Desktop for Mac.
{ .important }

### Out of memory issues

If your containers attempt to use more memory than the system has available, you
may experience an Out of Memory (OOM) exception, and a container, or the Docker
daemon, might be stopped by the kernel OOM killer. To prevent this from
happening, ensure that your application runs on hosts with adequate memory and
see
[Understand the risks of running out of memory](../containers/resource_constraints.md#understand-the-risks-of-running-out-of-memory).

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

### Kernel cgroup swap limit capabilities

On Ubuntu or Debian hosts, you may see messages similar to the following when
working with an image.

```text
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

   ```text
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

## Networking

### IP forwarding problems

If you manually configure your network using `systemd-network` with systemd
version 219 or later, Docker containers may not be able to access your network.
Beginning with systemd version 220, the forwarding setting for a given network
(`net.ipv4.conf.<interface>.forwarding`) defaults to off. This setting prevents
IP forwarding. It also conflicts with Docker's behavior of enabling the
`net.ipv4.conf.all.forwarding` setting within containers.

To work around this on RHEL, CentOS, or Fedora, edit the `<interface>.network`
file in `/usr/lib/systemd/network/` on your Docker host, for example,
`/usr/lib/systemd/network/80-container-host0.network`.

Add the following block within the `[Network]` section.

```systemd
[Network]
...
IPForward=kernel
# OR
IPForward=true
```

This configuration allows IP forwarding from the container as expected.

### DNS resolver issues

```console
DNS resolver found in resolv.conf and containers can't use it
```

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

```text
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

  Turning off `dnsmasq` adds the IP addresses of actual DNS nameservers to
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

### Turn off `dnsmasq`

{{< tabs >}}
{{< tab name="Ubuntu" >}}

If you prefer not to change the Docker daemon's configuration to use a specific
IP address, follow these instructions to turn off `dnsmasq` in NetworkManager.

1. Edit the `/etc/NetworkManager/NetworkManager.conf` file.

2. Comment out the `dns=dnsmasq` line by adding a `#` character to the beginning
   of the line.

   ```text
   # dns=dnsmasq
   ```

   Save and close the file.

3. Restart both NetworkManager and Docker. As an alternative, you can reboot
   your system.

   ```console
   $ sudo systemctl restart network-manager
   $ sudo systemctl restart docker
   ```

{{< /tab >}}
{{< tab name="RHEL, CentOS, or Fedora" >}}

To turn off `dnsmasq` on RHEL, CentOS, or Fedora:

1. Turn off the `dnsmasq` service:

   ```console
   $ sudo systemctl stop dnsmasq
   $ sudo systemctl disable dnsmasq
   ```

2. Configure the DNS servers manually using the
   [Red Hat documentation](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/configuring_and_managing_networking/configuring-the-order-of-dns-servers_configuring-and-managing-networking).

{{< /tab >}}
{{< /tabs >}}

## Volumes

### Unable to remove filesystem

```text
Error: Unable to remove filesystem
```

Some container-based utilities, such
as [Google cAdvisor](https://github.com/google/cadvisor), mount Docker system
directories, such as `/var/lib/docker/`, into a container. For instance, the
documentation for `cadvisor` instructs you to run the `cadvisor` container as
follows:

```console
$ sudo docker run \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:rw \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --publish=8080:8080 \
  --detach=true \
  --name=cadvisor \
  google/cadvisor:latest
```

When you bind-mount `/var/lib/docker/`, this effectively mounts all resources of
all other running containers as filesystems within the container which mounts
`/var/lib/docker/`. When you attempt to remove any of these containers, the
removal attempt may fail with an error like the following:

```none
Error: Unable to remove filesystem for
74bef250361c7817bee19349c93139621b272bc8f654ae112dd4eb9652af9515:
remove /var/lib/docker/containers/74bef250361c7817bee19349c93139621b272bc8f654ae112dd4eb9652af9515/shm:
Device or resource busy
```

The problem occurs if the container which bind-mounts `/var/lib/docker/`
uses `statfs` or `fstatfs` on filesystem handles within `/var/lib/docker/`
and does not close them.

Typically, we would advise against bind-mounting `/var/lib/docker` in this way.
However, `cAdvisor` requires this bind-mount for core functionality.

If you are unsure which process is causing the path mentioned in the error to
be busy and preventing it from being removed, you can use the `lsof` command
to find its process. For instance, for the error above:

```console
$ sudo lsof /var/lib/docker/containers/74bef250361c7817bee19349c93139621b272bc8f654ae112dd4eb9652af9515/shm
```

To work around this problem, stop the container which bind-mounts
`/var/lib/docker` and try again to remove the other container.
