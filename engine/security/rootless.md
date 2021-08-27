---
description: Run the Docker daemon as a non-root user (Rootless mode)
keywords: security, namespaces, rootless
title: Run the Docker daemon as a non-root user (Rootless mode)
---

Rootless mode allows running the Docker daemon and containers as a non-root
user to mitigate potential vulnerabilities in the daemon and
the container runtime.

Rootless mode does not require root privileges even during the installation of
the Docker daemon, as long as the [prerequisites](#prerequisites) are met.

Rootless mode was introduced in Docker Engine v19.03 as an experimental feature.
Rootless mode graduated from experimental in Docker Engine v20.10.

## How it works

Rootless mode executes the Docker daemon and containers inside a user namespace.
This is very similar to [`userns-remap` mode](userns-remap.md), except that
with `userns-remap` mode, the daemon itself is running with root privileges,
whereas in rootless mode, both the daemon and the container are running without
root privileges.

Rootless mode does not use binaries with `SETUID` bits or file capabilities,
except `newuidmap` and `newgidmap`, which are needed to allow multiple
UIDs/GIDs to be used in the user namespace.

## Prerequisites

-  You must install `newuidmap` and `newgidmap` on the host. These commands
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
testuser:231072:65536
```

### Distribution-specific hint

> Note: We recommend that you use the Ubuntu kernel.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#hint-ubuntu">Ubuntu</a></li>
  <li><a data-toggle="tab" data-target="#hint-debian">Debian GNU/Linux</a></li>
  <li><a data-toggle="tab" data-target="#hint-arch">Arch Linux</a></li>
  <li><a data-toggle="tab" data-target="#hint-opensuse-sles">openSUSE and SLES</a></li>
  <li><a data-toggle="tab" data-target="#hint-centos8-rhel8-fedora">CentOS 8, RHEL 8 and Fedora</a></li>
  <li><a data-toggle="tab" data-target="#hint-centos7-rhel7">CentOS 7 and RHEL 7</a></li>
</ul>
<div class="tab-content">

<div id="hint-ubuntu" class="tab-pane fade in active" markdown="1">
- Install `dbus-user-session` package if not installed. Run `sudo apt-get install -y dbus-user-session` and relogin.

- `overlay2` storage driver  is enabled by default
  ([Ubuntu-specific kernel patch](https://kernel.ubuntu.com/git/ubuntu/ubuntu-bionic.git/commit/fs/overlayfs?id=3b7da90f28fe1ed4b79ef2d994c81efbc58f1144)).

- Known to work on Ubuntu 18.04, 20.04, and 21.04.
</div>
<div id="hint-debian" class="tab-pane fade in" markdown="1">
- Install `dbus-user-session` package if not installed. Run `sudo apt-get install -y dbus-user-session` and relogin.

- For Debian 10, add `kernel.unprivileged_userns_clone=1` to `/etc/sysctl.conf` (or
  `/etc/sysctl.d`) and run `sudo sysctl --system`. This step is not required on Debian 11.

- Installing `fuse-overlayfs` is recommended. Run `sudo apt-get install -y fuse-overlayfs`.
  Using `overlay2` storage driver with Debian-specific modprobe option `sudo modprobe overlay permit_mounts_in_userns=1` is also possible,
  however, highly discouraged due to [instability](https://github.com/moby/moby/issues/42302).

- Rootless docker requires version of `slirp4netns` greater than `v0.4.0` (when `vpnkit` is not installed).
  Check you have this with 
  
  ```console
  $ slirp4netns --version
  ```
  If you do not have this download and install with `sudo apt-get install -y slirp4netns` or download the latest [release](https://github.com/rootless-containers/slirp4netns/releases).

</div>
<div id="hint-arch" class="tab-pane fade in" markdown="1">
- Installing `fuse-overlayfs` is recommended. Run `sudo pacman -S fuse-overlayfs`.

- Add `kernel.unprivileged_userns_clone=1` to `/etc/sysctl.conf` (or
  `/etc/sysctl.d`) and run `sudo sysctl --system`
</div>
<div id="hint-opensuse-sles" class="tab-pane fade in" markdown="1">
- Installing `fuse-overlayfs` is recommended. Run `sudo zypper install -y fuse-overlayfs`.

- `sudo modprobe ip_tables iptable_mangle iptable_nat iptable_filter` is required.
  This might be required on other distros as well depending on the configuration.

- Known to work on openSUSE 15 and SLES 15.
</div>
<div id="hint-centos8-rhel8-fedora" class="tab-pane fade in" markdown="1">
- Installing `fuse-overlayfs` is recommended. Run `sudo dnf install -y fuse-overlayfs`.

- You might need `sudo dnf install -y iptables`.

- Known to work on CentOS 8, RHEL 8, and Fedora 34.
</div>
<div id="hint-centos7-rhel7" class="tab-pane fade in" markdown="1">
- Add `user.max_user_namespaces=28633` to `/etc/sysctl.conf` (or 
  `/etc/sysctl.d`) and run `sudo sysctl --system`.

- `systemctl --user` does not work by default. 
  Run `dockerd-rootless.sh` directly without systemd.
</div>
</div> <!-- tab-content -->

## Known limitations

- Only the following storage drivers are supported:
  - `overlay2` (only if running with kernel 5.11 or later, or Ubuntu-flavored kernel)
  - `fuse-overlayfs` (only if running with kernel 4.18 or later, and `fuse-overlayfs` is installed)
  - `btrfs` (only if running with kernel 4.18 or later, or `~/.local/share/docker` is mounted with `user_subvol_rm_allowed` mount option)
  - `vfs`
- Cgroup is supported only when running with cgroup v2 and systemd. See [Limiting resources](#limiting-resources).
- Following features are not supported:
  - AppArmor
  - Checkpoint
  - Overlay network
  - Exposing SCTP ports
- To use the `ping` command, see [Routing ping packets](#routing-ping-packets).
- To expose privileged TCP/UDP ports (< 1024), see [Exposing privileged ports](#exposing-privileged-ports).
- `IPAddress` shown in `docker inspect` and is namespaced inside RootlessKit's network namespace.
  This means the IP address is not reachable from the host without `nsenter`-ing into the network namespace.
- Host network (`docker run --net=host`) is also namespaced inside RootlessKit.
- NFS mounts as the docker "data-root" is not supported. This limitation is not specific to rootless mode.

## Install
> **Note**
>
> If the system-wide Docker daemon is already running, consider disabling it:
> `$ sudo systemctl disable --now docker.service docker.socket`

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#install-with-packages">With packages (RPM/DEB)</a></li>
  <li><a data-toggle="tab" data-target="#install-without-packages">Without packages</a></li>
</ul>
<div class="tab-content">

<div id="install-with-packages" class="tab-pane fade in active" markdown="1">
If you installed Docker 20.10 or later with [RPM/DEB packages](/engine/install), you should have `dockerd-rootless-setuptool.sh` in `/usr/bin`.

Run `dockerd-rootless-setuptool.sh install` as a non-root user to set up the daemon:

```console
$ dockerd-rootless-setuptool.sh install
[INFO] Creating /home/testuser/.config/systemd/user/docker.service
...
[INFO] Installed docker.service successfully.
[INFO] To control docker.service, run: `systemctl --user (start|stop|restart) docker.service`
[INFO] To run docker.service on system startup, run: `sudo loginctl enable-linger testuser`

[INFO] Make sure the following environment variables are set (or add them to ~/.bashrc):

export PATH=/usr/bin:$PATH
export DOCKER_HOST=unix:///run/user/1000/docker.sock
```

If `dockerd-rootless-setuptool.sh` is not present, you may need to install the `docker-ce-rootless-extras` package manually, e.g.,

```console
$ sudo apt-get install -y docker-ce-rootless-extras
```

</div>
<div id="install-without-packages" class="tab-pane fade in" markdown="1">
If you do not have permission to run package managers like `apt-get` and `dnf`,
consider using the installation script available at [https://get.docker.com/rootless](https://get.docker.com/rootless){: target="_blank" rel="noopener" class="_" }.
Since static packages are not available for `s390x`, hence it is not supported for `s390x`.

```console
$ curl -fsSL https://get.docker.com/rootless | sh
...
[INFO] Creating /home/testuser/.config/systemd/user/docker.service
...
[INFO] Installed docker.service successfully.
[INFO] To control docker.service, run: `systemctl --user (start|stop|restart) docker.service`
[INFO] To run docker.service on system startup, run: `sudo loginctl enable-linger testuser`

[INFO] Make sure the following environment variables are set (or add them to ~/.bashrc):

export PATH=/home/testuser/bin:$PATH
export DOCKER_HOST=unix:///run/user/1000/docker.sock
```

The binaries will be installed at `~/bin`.
</div>
</div> <!-- tab-content -->

See [Troubleshooting](#troubleshooting) if you faced an error.

## Uninstall

To remove the systemd service of the Docker daemon, run `dockerd-rootless-setuptool.sh uninstall`:

```console
$ dockerd-rootless-setuptool.sh uninstall
+ systemctl --user stop docker.service
+ systemctl --user disable docker.service
Removed /home/testuser/.config/systemd/user/default.target.wants/docker.service.
[INFO] Uninstalled docker.service
[INFO] This uninstallation tool does NOT remove Docker binaries and data.
[INFO] To remove data, run: `/usr/bin/rootlesskit rm -rf /home/testuser/.local/share/docker`
```

Unset environment variables PATH and DOCKER_HOST if you have added them to `~/.bashrc`.

To remove the data directory, run `rootlesskit rm -rf ~/.local/share/docker`.

To remove the binaries, remove `docker-ce-rootless-extras` package if you installed Docker with package managers.
If you installed Docker with https://get.docker.com/rootless ([Install without packages](#install)),
remove the binary files under `~/bin`:
```console
$ cd ~/bin
$ rm -f containerd containerd-shim containerd-shim-runc-v2 ctr docker docker-init docker-proxy dockerd dockerd-rootless-setuptool.sh dockerd-rootless.sh rootlesskit rootlesskit-docker-proxy runc vpnkit
```

## Usage

### Daemon
<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#usage-with-systemd">With systemd (Highly recommended)</a></li>
  <li><a data-toggle="tab" data-target="#usage-without-systemd">Without systemd</a></li>
</ul>
<div class="tab-content">

<div id="usage-with-systemd" class="tab-pane fade in active" markdown="1">
The systemd unit file is installed as  `~/.config/systemd/user/docker.service`.

Use `systemctl --user` to manage the lifecycle of the daemon:

```console
$ systemctl --user start docker
```

To launch the daemon on system startup, enable the systemd service and lingering:

```console
$ systemctl --user enable docker
$ sudo loginctl enable-linger $(whoami)
```

Starting Rootless Docker as a systemd-wide service (`/etc/systemd/system/docker.service`)
is not supported, even with the `User=` directive.

</div>
<div id="usage-without-systemd" class="tab-pane fade in" markdown="1">
To run the daemon directly without systemd, you need to run `dockerd-rootless.sh` instead of `dockerd`.

The following environment variables must be set:
- `$HOME`: the home directory
- `$XDG_RUNTIME_DIR`: an ephemeral directory that is only accessible by the expected user, e,g, `~/.docker/run`.
  The directory should be removed on every host shutdown.
  The directory can be on tmpfs, however, should not be under `/tmp`.
  Locating this directory under `/tmp` might be vulnerable to TOCTOU attack.

</div>
</div> <!-- tab-content -->

Remarks about directory paths:

- The socket path is set to `$XDG_RUNTIME_DIR/docker.sock` by default.
  `$XDG_RUNTIME_DIR` is typically set to `/run/user/$UID`.
- The data dir is set to `~/.local/share/docker` by default.
  The data dir should not be on NFS.
- The daemon config dir is set to `~/.config/docker` by default.
  This directory is different from `~/.docker` that is used by the client.

### Client

You need to specify either the socket path or the CLI context explicitly.

To specify the socket path using `$DOCKER_HOST`:

```console
$ export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock
$ docker run -d -p 8080:80 nginx
```

To specify the CLI context using `docker context`:

```console
$ docker context use rootless
rootless
Current context is now "rootless"
$ docker run -d -p 8080:80 nginx
```

## Best practices

### Rootless Docker in Docker

To run Rootless Docker inside "rootful" Docker, use the `docker:<version>-dind-rootless`
image instead of `docker:<version>-dind`.

```console
$ docker run -d --name dind-rootless --privileged docker:20.10-dind-rootless
```

The `docker:<version>-dind-rootless` image runs as a non-root user (UID 1000).
However, `--privileged` is required for disabling seccomp, AppArmor, and mount
masks.

### Expose Docker API socket through TCP

To expose the Docker API socket through TCP, you need to launch `dockerd-rootless.sh`
with `DOCKERD_ROOTLESS_ROOTLESSKIT_FLAGS="-p 0.0.0.0:2376:2376/tcp"`.

```console
$ DOCKERD_ROOTLESS_ROOTLESSKIT_FLAGS="-p 0.0.0.0:2376:2376/tcp" \
  dockerd-rootless.sh \
  -H tcp://0.0.0.0:2376 \
  --tlsverify --tlscacert=ca.pem --tlscert=cert.pem --tlskey=key.pem
```

### Expose Docker API socket through SSH

To expose the Docker API socket through SSH, you need to make sure `$DOCKER_HOST`
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
Limiting resources with cgroup-related `docker run` flags such as `--cpus`, `--memory`, `--pids-limit`
is supported only when running with cgroup v2 and systemd.
See [Changing cgroup version](../../config/containers/runmetrics.md) to enable cgroup v2.

If `docker info` shows `none` as `Cgroup Driver`, the conditions are not satisfied.
When these conditions are not satisfied, rootless mode ignores the cgroup-related `docker run` flags.
See [Limiting resources without cgroup](#limiting-resources-without-cgroup) for workarounds.

If `docker info` shows `systemd` as `Cgroup Driver`, the conditions are satisfied.
However, typically, only `memory` and `pids` controllers are delegated to non-root users by default.

```console
$ cat /sys/fs/cgroup/user.slice/user-$(id -u).slice/user@$(id -u).service/cgroup.controllers
memory pids
```

To allow delegation of all controllers, you need to change the systemd configuration as follows:

```console
# mkdir -p /etc/systemd/system/user@.service.d
# cat > /etc/systemd/system/user@.service.d/delegate.conf << EOF
[Service]
Delegate=cpu cpuset io memory pids
EOF
# systemctl daemon-reload
```

> **Note**
>
> Delegating `cpuset` requires systemd 244 or later.

#### Limiting resources without cgroup
Even when cgroup is not available, you can still use the traditional `ulimit` and [`cpulimit`](https://github.com/opsengine/cpulimit),
though they work in process-granularity rather than in container-granularity,
and can be arbitrarily disabled by the container process.

For example:

- To limit CPU usage to 0.5 cores (similar to `docker run --cpus 0.5`):
  `docker run <IMAGE> cpulimit --limit=50 --include-children <COMMAND>`
- To limit max VSZ to 64MiB (similar to `docker run --memory 64m`):
  `docker run <IMAGE> sh -c "ulimit -v 65536; <COMMAND>"`

- To limit max number of processes to 100 per namespaced UID 2000
  (similar to `docker run --pids-limit=100`):
  `docker run --user 2000 --ulimit nproc=100 <IMAGE> <COMMAND>`

## Troubleshooting

### Errors when starting the Docker daemon

**[rootlesskit:parent] error: failed to start the child: fork/exec /proc/self/exe: operation not permitted**

This error occurs mostly when the value of `/proc/sys/kernel/unprivileged_userns_clone ` is set to 0:

```console
$ cat /proc/sys/kernel/unprivileged_userns_clone
0
```

To fix this issue, add  `kernel.unprivileged_userns_clone=1` to
`/etc/sysctl.conf` (or `/etc/sysctl.d`) and run `sudo sysctl --system`.

**[rootlesskit:parent] error: failed to start the child: fork/exec /proc/self/exe: no space left on device**

This error occurs mostly when the value of `/proc/sys/user/max_user_namespaces` is too small:

```console
$ cat /proc/sys/user/max_user_namespaces
0
```

To fix this issue, add  `user.max_user_namespaces=28633` to
`/etc/sysctl.conf` (or `/etc/sysctl.d`) and run `sudo sysctl --system`.

**[rootlesskit:parent] error: failed to setup UID/GID map: failed to compute uid/gid map: No subuid ranges found for user 1001 ("testuser")**

This error occurs when `/etc/subuid` and `/etc/subgid` are not configured. See [Prerequisites](#prerequisites).

**could not get XDG_RUNTIME_DIR**

This error occurs when `$XDG_RUNTIME_DIR` is not set.

On a non-systemd host, you need to create a directory and then set the path:

```console
$ export XDG_RUNTIME_DIR=$HOME/.docker/xrd
$ rm -rf $XDG_RUNTIME_DIR
$ mkdir -p $XDG_RUNTIME_DIR
$ dockerd-rootless.sh
```

> **Note**:
> You must remove the directory every time you log out.

On a systemd host, log into the host using `pam_systemd` (see below).
The value is automatically set to `/run/user/$UID` and cleaned up on every logout.

**`systemctl --user` fails with "Failed to connect to bus: No such file or directory"**

This error occurs mostly when you switch from the root user to an non-root user with `sudo`:

```console
# sudo -iu testuser
$ systemctl --user start docker
Failed to connect to bus: No such file or directory
```

Instead of `sudo -iu <USERNAME>`, you need to log in using `pam_systemd`. For example:

- Log in through the graphic console
- `ssh <USERNAME>@localhost`
- `machinectl shell <USERNAME>@`

**The daemon does not start up automatically**

You need `sudo loginctl enable-linger $(whoami)` to enable the daemon to start
up automatically. See [Usage](#usage).

**iptables failed: iptables -t nat -N DOCKER: Fatal: can't open lock file /run/xtables.lock: Permission denied**

This error may happen with an older version of Docker when SELinux is enabled on the host.

The issue has been fixed in Docker 20.10.8.
A known workaround for older version of Docker is to run the following commands to disable SELinux for `iptables`:
```console
$ sudo dnf install -y policycoreutils-python-utils && sudo semanage permissive -a iptables_t
```

### `docker pull` errors

**docker: failed to register layer: Error processing tar file(exit status 1): lchown &lt;FILE&gt;: invalid argument**

This error occurs when the number of available entries in `/etc/subuid` or
`/etc/subgid` is not sufficient. The number of entries required vary across
images. However, 65,536 entries are sufficient for most images. See
[Prerequisites](#prerequisites).

**docker: failed to register layer: ApplyLayer exit status 1 stdout:  stderr: lchown &lt;FILE&gt;: operation not permitted**

This error occurs mostly when `~/.local/share/docker` is located on NFS.

A workaround is to specify non-NFS `data-root` directory in `~/.config/docker/daemon.json` as follows:
```json
{"data-root":"/somewhere-out-of-nfs"}
```

### `docker run` errors

**docker: Error response from daemon: OCI runtime create failed: ...: read unix @-&gt;/run/systemd/private: read: connection reset by peer: unknown.**

This error occurs on cgroup v2 hosts mostly when the dbus daemon is not running for the user.

```console
$ systemctl --user is-active dbus
inactive

$ docker run hello-world
docker: Error response from daemon: OCI runtime create failed: container_linux.go:380: starting container process caused: process_linux.go:385: applying cgroup configuration for process caused: error while starting unit "docker
-931c15729b5a968ce803784d04c7421f791d87e5ca1891f34387bb9f694c488e.scope" with properties [{Name:Description Value:"libcontainer container 931c15729b5a968ce803784d04c7421f791d87e5ca1891f34387bb9f694c488e"} {Name:Slice Value:"use
r.slice"} {Name:PIDs Value:@au [4529]} {Name:Delegate Value:true} {Name:MemoryAccounting Value:true} {Name:CPUAccounting Value:true} {Name:IOAccounting Value:true} {Name:TasksAccounting Value:true} {Name:DefaultDependencies Val
ue:false}]: read unix @->/run/systemd/private: read: connection reset by peer: unknown.
```

To fix the issue, run `sudo apt-get install -y dbus-user-session` or `sudo dnf install -y dbus-daemon`, and then relogin.

If the error still occurs, try running `systemctl --user enable --now dbus` (without sudo).

**`--cpus`, `--memory`, and `--pids-limit` are ignored**

This is an expected behavior on cgroup v1 mode.
To use these flags, the host needs to be configured for enabling cgroup v2.
For more information, see [Limiting resources](#limiting-resources).

### Networking errors

**`docker run -p` fails with `cannot expose privileged port`**

`docker run -p` fails with this error when a privileged port (< 1024) is specified as the host port.

```console
$ docker run -p 80:80 nginx:alpine
docker: Error response from daemon: driver failed programming external connectivity on endpoint focused_swanson (9e2e139a9d8fc92b37c36edfa6214a6e986fa2028c0cc359812f685173fa6df7): Error starting userland proxy: error while calling PortManager.AddPort(): cannot expose privileged port 80, you might need to add "net.ipv4.ip_unprivileged_port_start=0" (currently 1024) to /etc/sysctl.conf, or set CAP_NET_BIND_SERVICE on rootlesskit binary, or choose a larger port number (>= 1024): listen tcp 0.0.0.0:80: bind: permission denied.
```

When you experience this error, consider using an unprivileged port instead. For example, 8080 instead of 80.

```console
$ docker run -p 8080:80 nginx:alpine
```

To allow exposing privileged ports, see [Exposing privileged ports](#exposing-privileged-ports).

**ping doesn't work**

Ping does not work when `/proc/sys/net/ipv4/ping_group_range` is set to `1 0`:

```console
$ cat /proc/sys/net/ipv4/ping_group_range
1       0
```

For details, see [Routing ping packets](#routing-ping-packets).

**`IPAddress` shown in `docker inspect` is unreachable**

This is an expected behavior, as the daemon is namespaced inside RootlessKit's
network namespace. Use `docker run -p` instead.

**`--net=host` doesn't listen ports on the host network namespace**

This is an expected behavior, as the daemon is namespaced inside RootlessKit's
network namespace. Use `docker run -p` instead.

**Network is slow**

Docker with rootless mode uses [slirp4netns](https://github.com/rootless-containers/slirp4netns) as the default network stack if slirp4netns v0.4.0 or later is installed.
If slirp4netns is not installed, Docker falls back to [VPNKit](https://github.com/moby/vpnkit).

Installing slirp4netns may improve the network throughput.
See [RootlessKit documentation](https://github.com/rootless-containers/rootlesskit/tree/v0.13.0#network-drivers) for the benchmark result.

Also, changing MTU value may improve the throughput.
The MTU value can be specified by creating `~/.config/systemd/user/docker.service.d/override.conf` with the following content:

```systemd
[Service]
Environment="DOCKERD_ROOTLESS_ROOTLESSKIT_MTU=<INTEGER>"
```

And then restart the daemon:
```console
$ systemctl --user daemon-reload
$ systemctl --user restart docker
```

**`docker run -p` does not propagate source IP addresses**

This is because Docker with rootless mode uses RootlessKit's builtin port driver by default.

The source IP addresses can be propagated by creating `~/.config/systemd/user/docker.service.d/override.conf` with the following content:

```systemd
[Service]
Environment="DOCKERD_ROOTLESS_ROOTLESSKIT_PORT_DRIVER=slirp4netns"
```

And then restart the daemon:
```console
$ systemctl --user daemon-reload
$ systemctl --user restart docker
```

Note that this configuration decreases throughput.
See [RootlessKit documentation](https://github.com/rootless-containers/rootlesskit/tree/v0.13.0#port-drivers) for the benchmark result.

### Tips for debugging
**Entering into `dockerd` namespaces**

The `dockerd-rootless.sh` script executes `dockerd` in its own user, mount, and network namespaces.

For debugging, you can enter the namespaces by running
`nsenter -U --preserve-credentials -n -m -t $(cat $XDG_RUNTIME_DIR/docker.pid)`.
