---
description: Run the Docker daemon as a non-root user (Rootless mode)
keywords: security, namespaces, rootless
title: Run the Docker daemon as a non-root user (Rootless mode)
---

Rootless mode allows running the Docker daemon and containers as a non-root
user, for the sake of mitigating potential vulnerabilities in the daemon and
the container runtime.

Rootless mode does not require root privileges even for installation of the
Docker daemon, as long as [the prerequisites](#prerequisites) are satisfied.

Rootless mode was introduced in Docker Engine 19.03.

> **Note**:
> Rootless mode is an experimental feature and has [limitations](#known-limitations).

## How it works
Rootless mode executes the Docker daemon and containers inside a user namespace.
This is very similar to [`userns-remap` mode](userns-remap.md), except that
with `userns-remap` mode, the daemon itself is running with root privileges, whereas in
rootless mode, both the daemon and the container are running without root privileges.

Rootless mode does not use binaries with SETUID bits or file capabilities,
except `newuidmap` and `newgidmap`, which are needed to allow multiple
UIDs/GIDs to be used in the user namespace.

## Prerequisites

- `newuidmap` and `newgidmap` need to be installed on the host. These commands
  are provided by the `uidmap` package on most distros.

- `/etc/subuid` and `/etc/subgid` should contain at least 65,536 subordinate
  UIDs/GIDs for the user. In the following example, the user `testuser` has
  65,536 subordinate UIDs/GIDs (231072-296607).

```console
$ id -u
1001
$ whoami
testuser
$ grep ^$(whoami): /etc/subuid
testuser:231072:65536
$ grep ^$(whoami): /etc/subgid
testuser::231072:65536
```

### Distribution-specific hint

> Note: Using Ubuntu kernel is recommended.

#### Ubuntu
- No preparation is needed.

- `overlay2` storage driver  is enabled by default
  ([Ubuntu-specific kernel patch](https://kernel.ubuntu.com/git/ubuntu/ubuntu-bionic.git/commit/fs/overlayfs?id=3b7da90f28fe1ed4b79ef2d994c81efbc58f1144)).

- Known to work on Ubuntu 16.04, 18.04, and 20.04.

#### Debian GNU/Linux
- Add `kernel.unprivileged_userns_clone=1` to `/etc/sysctl.conf` (or
  `/etc/sysctl.d`) and run `sudo sysctl --system`.

- To use the `overlay2` storage driver (recommended), run
  `sudo modprobe overlay permit_mounts_in_userns=1`
   ([Debian-specific kernel patch, introduced in Debian 10](https://salsa.debian.org/kernel-team/linux/blob/283390e7feb21b47779b48e0c8eb0cc409d2c815/debian/patches/debian/overlayfs-permit-mounts-in-userns.patch)).
   Put the configuration to `/etc/modprobe.d` for persistence.

- Known to work on Debian 9 and 10. 
  `overlay2` is only supported since Debian 10 and needs `modprobe`
  configuration described above.

#### Arch Linux
- Add `kernel.unprivileged_userns_clone=1` to `/etc/sysctl.conf` (or
  `/etc/sysctl.d`) and run `sudo sysctl --system`

#### openSUSE
- `sudo modprobe ip_tables iptable_mangle iptable_nat iptable_filter` is required.
  This might be required on other distros as well depending on the configuration.

- Known to work on openSUSE 15.

#### Fedora 31 and later
- Fedora 31 uses cgroup v2 by default, which is not yet supported by the containerd runtime.
  Run `sudo grubby --update-kernel=ALL --args="systemd.unified_cgroup_hierarchy=0"`
  to use cgroup v1.
- `sudo dnf install -y iptables` might be needed.

#### CentOS 8
- `sudo dnf install -y iptables` might be needed.

#### CentOS 7
- Add `user.max_user_namespaces=28633` to `/etc/sysctl.conf` (or 
  `/etc/sysctl.d`) and run `sudo sysctl --system`.

- `systemctl --user` does not work by default. 
  Run the daemon directly without systemd:
  `dockerd-rootless.sh --experimental --storage-driver vfs`

- Known to work on CentOS 7.7. Older releases require extra configuration
  steps.

- CentOS 7.6 and older releases require [COPR package `vbatts/shadow-utils-newxidmap`](https://copr.fedorainfracloud.org/coprs/vbatts/shadow-utils-newxidmap/) to be installed.

- CentOS 7.5 and older releases require running
  `sudo grubby --update-kernel=ALL --args="user_namespace.enable=1"` and reboot.

## Known limitations

- Only `vfs` graphdriver is supported. However, on Ubuntu and Debian 10,
  `overlay2` and `overlay` are also supported.
- Following features are not supported:
  - Cgroups (including `docker top`, which depends on the cgroups)
  - AppArmor
  - Checkpoint
  - Overlay network
  - Exposing SCTP ports
- To use `ping` command, see [Routing ping packets](#routing-ping-packets)
- To expose privileged TCP/UDP ports (< 1024), see [Exposing privileged ports](#exposing-privileged-ports)
- `IPAddress` shown in `docker inspect` and is namespaced inside RootlessKit's network namespace.
  This means the IP address is not reachable from the host without `nsenter`-ing into the network namespace.
- Host network (`docker run --net=host`) is namespaced inside RootlessKit as well.

## Install

The installation script is available at [https://get.docker.com/rootless](https://get.docker.com/rootless){: target="_blank" class="_" }.

```console
$ curl -fsSL https://get.docker.com/rootless | sh
```

Make sure to run the script as a non-root user.
To install Rootless Docker as the root user, see [Manual installation](#manual-installation) steps.

The script will show the environment variables that are needed to be set:

```console
$ curl -fsSL https://get.docker.com/rootless | sh
...
# Docker binaries are installed in /home/testuser/bin
# WARN: dockerd is not in your current PATH or pointing to /home/testuser/bin/dockerd
# Make sure the following environment variables are set (or add them to ~/.bashrc):

export PATH=/home/testuser/bin:$PATH
export PATH=$PATH:/sbin
export DOCKER_HOST=unix:///run/user/1001/docker.sock

#
# To control docker service run:
# systemctl --user (start|stop|restart) docker
#
```

### Manual installation
To install the binaries manually without using the installer, extract
`docker-rootless-extras-<version>.tar.gz` along with `docker-<version>.tar.gz`:
from [https://download.docker.com/linux/static/stable/x86_64/](https://download.docker.com/linux/static/stable/x86_64/){: target="_blank" class="_" }

If you already have Docker daemon running as the root, you only need to extract `docker-rootless-extras-<version>.tar.gz`.
The archive can be extracted under an arbitrary directory listed in the `$PATH`. e.g. `/usr/local/bin`, or `$HOME/bin`.

### Nightly channel

To install a nightly version of Rootless Docker, execute the installation script with `CHANNEL="nightly"`:

```console
$ curl -fsSL https://get.docker.com/rootless | CHANNEL="nightly" sh
```

The raw binary archives are available at:
- https://master.dockerproject.org/linux/x86_64/docker-rootless-extras.tgz
- https://master.dockerproject.org/linux/x86_64/docker.tgz

## Usage

### Daemon

Use `systemctl --user` to manage the lifecycle of the daemon:

```console
$ systemctl --user start docker
```

To launch the daemon on system startup, enable the systemd service and lingering:

```console
$ systemctl --user enable docker
$ sudo loginctl enable-linger $(whoami)
```

To run the daemon directly without systemd, you need to run
`dockerd-rootless.sh` instead of `dockerd`:

```console
$ dockerd-rootless.sh --experimental --storage-driver vfs
```

As Rootless mode is experimental, you need to run
`dockerd-rootless.sh` with `--experimental`.
You also need `--storage-driver vfs` unless using Ubuntu or Debian 10 kernel.
You don't need to care about these flags if you manage the daemon using systemd, as
these flags are automatically added to the systemd unit file.

Remarks about directory paths:
- The socket path is set to `$XDG_RUNTIME_DIR/docker.sock` by default.
  `$XDG_RUNTIME_DIR` is typically set to `/run/user/$UID`.
- The data dir is set to `~/.local/share/docker` by default.
- The exec dir is set to `$XDG_RUNTIME_DIR/docker` by default.
- The daemon config dir is set to `~/.config/docker` (not `~/.docker`, which is
  used by the client) by default.

Other remarks:
- The `dockerd-rootless.sh` script executes `dockerd` in its own user, mount,
  and network namespaces. You can enter the namespaces by running 
  `nsenter -U --preserve-credentials -n -m -t $(cat $XDG_RUNTIME_DIR/docker.pid)`.
- `docker info` shows `rootless` in `SecurityOptions`
- `docker info` shows `none` as `Cgroup Driver`

### Client

You need to specify the socket path explicitly.

To specify the socket path via `$DOCKER_HOST`:
```console
$ export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock
$ docker run -d -p 8080:80 nginx
```

To specify the socket path via `docker context`:
```console
$ docker context create rootless --description "for rootless mode" --docker "host=unix://$XDG_RUNTIME_DIR/docker.sock"
rootless
Successfully created context "rootless"
$ docker context use rootless
rootless
Current context is now "rootless"
$ docker run -d -p 8080:80 nginx
```

## Tips

### Rootless Docker in Docker

To run Rootless Docker inside "rootful" Docker, use `docker:<version>-dind-rootless`
image instead of `docker:<version>-dind` image.

```console
$ docker run -d --name dind-rootless --privileged docker:19.03-dind-rootless --experimental
```

`docker:<version>-dind-rootless` image runs as a non-root user (UID 1000).
However, `--privileged` is required for disabling seccomp, AppArmor, and mount
masks.

### Expose Docker API socket via TCP

To expose the Docker API socket via TCP, you need to launch `dockerd-rootless.sh`
with `DOCKERD_ROOTLESS_ROOTLESSKIT_FLAGS="-p 0.0.0.0:2376:2376/tcp"`.

```console
$ DOCKERD_ROOTLESS_ROOTLESSKIT_FLAGS="-p 0.0.0.0:2376:2376/tcp" \
  dockerd-rootless.sh --experimental \
  -H tcp://0.0.0.0:2376 \
  --tlsverify --tlscacert=ca.pem --tlscert=cert.pem --tlskey=key.pem
```

### Expose Docker API socket via SSH

To expose the Docker API socket via SSH, you need to make sure `$DOCKER_HOST`
is set on the remote host.

```console
$ ssh -l <REMOTEUSER> <REMOTEHOST> 'echo $DOCKER_HOST'
unix:///run/user/1001/docker.sock
$ docker -H ssh://<REMOTEUSER>@<REMOTEHOST> run ...
```

### Routing ping packets

On some distributions, `ping` does not work by default.

Add `net.ipv4.ping_group_range = 0   2147483647` to `/etc/sysctl.conf` (or
`/etc/sysctl.d`) and run `sudo sysctl --system` to allow using `ping`.

### Exposing privileged ports

To expose privileged ports (< 1024), set `CAP_NET_BIND_SERVICE` on `rootlesskit` binary.

```console
$ sudo setcap cap_net_bind_service=ep $HOME/bin/rootlesskit
```

Or add `net.ipv4.ip_unprivileged_port_start=0` to `/etc/sysctl.conf` (or
`/etc/sysctl.d`) and run `sudo sysctl --system`.

### Limiting resources

In Docker 19.03, rootless mode ignores cgroup-related `docker run` flags such as
`--cpus`, `--memory`, --pids-limit`.

However, traditional `ulimit` and [`cpulimit`](https://github.com/opsengine/cpulimit)
can be still used, though they work in process-granularity rather than in container-granularity,
and can be arbitrary disabled by the container process.

e.g.
- To limit CPU usage to 0.5 cores (akin to `docker run --cpus 0.5):
  `docker run <IMAGE> cpulimit --limit=50 --include-children <COMMAND>`

- To limit max VSZ to 64MiB (akin to `docker run --memory 64m`):
  `docker run <IMAGE> sh -c "ulimit -v 65536; <COMMAND>"`

- To limit max number of processes to 100 per namespaced UID 2000
  (akin to `docker run --pids-limit=100):
  `docker run --user 2000 --ulimit nproc=100 <IMAGE> <COMMAND>`

### Changing network stack

`dockerd-rootless.sh` uses [slirp4netns](https://github.com/rootless-containers/slirp4netns)
(if installed) or [VPNKit](https://github.com/moby/vpnkit) as the network stack
by default.

These network stacks run in userspace and might have performance overhead.
See [RootlessKit documentation](https://github.com/rootless-containers/rootlesskit/tree/v0.9.5#network-drivers) for further information.

Optionally, you can use `lxc-user-nic` instead for the best performance.
To use `lxc-user-nic`, you need to edit [`/etc/lxc/lxc-usernet`](https://github.com/rootless-containers/rootlesskit/tree/v0.9.5#--netlxc-user-nic-experimental)
and set `$DOCKERD_ROOTLESS_ROOTLESSKIT_NET=lxc-user-nic`.

## Troubleshooting

### Troubles during starting the daemon
#### `[rootlesskit:parent] error: failed to start the child: fork/exec /proc/self/exe: operation not permitted`

This error happens mostly when the value of `/proc/sys/kernel/unprivileged_userns_clone ` is set to 0:

```console
$ cat /proc/sys/kernel/unprivileged_userns_clone
0
```

To fix the issue, add  `kernel.unprivileged_userns_clone=1` to
`/etc/sysctl.conf` (or `/etc/sysctl.d`) and run `sudo sysctl --system`.

#### `[rootlesskit:parent] error: failed to start the child: fork/exec /proc/self/exe: no space left on device`

This error happens mostly when the value of `/proc/sys/user/max_user_namespaces` is too small:

```console
$ cat /proc/sys/user/max_user_namespaces
0
```

To fix the issue, add  `user.max_user_namespaces=28633` to
`/etc/sysctl.conf` (or `/etc/sysctl.d`) and run `sudo sysctl --system`.

#### `[rootlesskit:parent] error: failed to setup UID/GID map: failed to compute uid/gid map: No subuid ranges found for user 1001 ("testuser")`

This error happens when `/etc/subuid` and `/etc/subgid` are not configured.

See [Prerequisites](#prerequisites).

#### `could not get XDG_RUNTIME_DIR`
This error happens when `$XDG_RUNTIME_DIR` is not set.

On a non-systemd host, you need to create a directory and set the path by yourself:
```console
$ export XDG_RUNTIME_DIR=$HOME/.docker/xrd
$ rm -rf $XDG_RUNTIME_DIR
$ mkdir -p $XDG_RUNTIME_DIR
$ dockerd-rootless.sh --experimental
```

> **Note**:
> You have to remove the directory on every logout.

On a systemd host, login to the host via `pam_systemd` (see below).
The value is automatically set to `/run/user/$UID` and cleaned up on every logout.

#### `systemctl --user` fails with `Failed to connect to bus: No such file or directory`

This error happens mostly when you switched from the root user to an non-root user with `sudo`:

```console
# sudo -iu testuser
$ systemctl --user start docker
Failed to connect to bus: No such file or directory
```

Instead of `sudo -iu <USERNAME>`, you need to login via `pam_systemd`, e.g.
- Login via the graphic console
- `ssh <USERNAME>@localhost`
- `machinectl shell <USERNAME>@`

#### The daemon does not start up automatically

You need `sudo loginctl enable-linger $(whoami)` to enable the daemon to start
up automatically. See [Usage](#usage).

#### `rootless mode is supported only when running in experimental mode`

This error happens when the daemon was launched without `--experimental`.
See [Usage](#usage).

### Troubles during `docker pull`
#### `docker: failed to register layer: Error processing tar file(exit status 1): lchown <FILE>: invalid argument`

This error happens when the number of available entries in `/etc/subuid` or `/etc/subgid` is not sufficient.
The number of required entries vary across images, but having 65,536 entries is enough for most images.

See [Prerequisites](#prerequisites).

### Errors during `docker run`

#### `--cpus`, `--memory`, and `--pids-limit` are ignored

Expected behavior in Docker 19.03.
See [Limiting resources](#limiting-resources).

#### `Error response from daemon: cgroups: cgroup mountpoint does not exist: unknown.`

This error happens mostly when the host is running with cgroup v2.
See [Fedora 31 or later](#fedora-31-or-later) to switch the host to use cgroup v1.

### Networking

#### `docker run -p` fails with `cannot expose privileged port ...`

`docker run -p` fails with this error when an privileged port (< 1024) is specified as the host port.

```console
$ docker run -p 80:80 nginx:alpine
docker: Error response from daemon: driver failed programming external connectivity on endpoint focused_swanson (9e2e139a9d8fc92b37c36edfa6214a6e986fa2028c0cc359812f685173fa6df7): Error starting userland proxy: error while calling PortManager.AddPort(): cannot expose privileged port 80, you might need to add "net.ipv4.ip_unprivileged_port_start=0" (currently 1024) to /etc/sysctl.conf, or set CAP_NET_BIND_SERVICE on rootlesskit binary, or choose a larger port number (>= 1024): listen tcp 0.0.0.0:80: bind: permission denied.
```

When this error happened, consider using an unprivileged port instead, e.g. 8080 instead of 80.

```console
$ docker run -p 8080:80 nginx:alpine
```

To allow exposing privileged ports, see [Exposing privileged ports](#exposing-privileged-ports).

#### ping doesn't work

Ping does not work when `/proc/sys/net/ipv4/ping_group_range` is set to `1 0`:

```console
$ cat /proc/sys/net/ipv4/ping_group_range
1       0
```

See [Routing ping packets](#routing-ping-packets).

#### `IPAddress` shown in `docker inspect` is unreachable

Expected behavior, as the daemon is namespaced inside RootlessKit's network namespace.
Use `docker run -p` instead.

#### `--net=host` doesn't listen ports on the host network namespace

Expected behavior, as the daemon is namespaced inside RootlessKit's network namespace.
Use `docker run -p` instead.
