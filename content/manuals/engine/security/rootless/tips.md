---
description: Tips for the Rootless mode
keywords: security, namespaces, rootless
title: Tips
weight: 20
---

## Advanced usage

### Daemon

{{< tabs >}}
{{< tab name="With systemd (Highly recommended)" >}}

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

{{< /tab >}}
{{< tab name="Without systemd" >}}

To run the daemon directly without systemd, you need to run `dockerd-rootless.sh` instead of `dockerd`.

The following environment variables must be set:
- `$HOME`: the home directory
- `$XDG_RUNTIME_DIR`: an ephemeral directory that is only accessible by the expected user, e,g, `~/.docker/run`.
  The directory should be removed on every host shutdown.
  The directory can be on tmpfs, however, should not be under `/tmp`.
  Locating this directory under `/tmp` might be vulnerable to TOCTOU attack.

{{< /tab >}}
{{< /tabs >}}

It's important to note that with directory paths:

- The socket path is set to `$XDG_RUNTIME_DIR/docker.sock` by default.
  `$XDG_RUNTIME_DIR` is typically set to `/run/user/$UID`.
- The data dir is set to `~/.local/share/docker` by default.
  The data dir should not be on NFS.
- The daemon config dir is set to `~/.config/docker` by default.
  This directory is different from `~/.docker` that is used by the client.

### Client

Since Docker Engine v23.0, `dockerd-rootless-setuptool.sh install` automatically configures
the `docker` CLI to use the `rootless` context.

Prior to Docker Engine v23.0, a user had to specify either the socket path or the CLI context explicitly.

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
$ docker run -d --name dind-rootless --privileged docker:25.0-dind-rootless
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

To expose privileged ports (< 1024), set `CAP_NET_BIND_SERVICE` on `rootlesskit` binary and restart the daemon.

```console
$ sudo setcap cap_net_bind_service=ep $(which rootlesskit)
$ systemctl --user restart docker
```

Or add `net.ipv4.ip_unprivileged_port_start=0` to `/etc/sysctl.conf` (or
`/etc/sysctl.d`) and run `sudo sysctl --system`.

### Limiting resources

Limiting resources with cgroup-related `docker run` flags such as `--cpus`, `--memory`, `--pids-limit`
is supported only when running with cgroup v2 and systemd.
See [Changing cgroup version](/manuals/engine/containers/runmetrics.md) to enable cgroup v2.

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

> [!NOTE]
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
