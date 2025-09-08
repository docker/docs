---
description: Run the Docker daemon as a non-root user (Rootless mode)
keywords: security, namespaces, rootless
title: Rootless mode
weight: 10
---

Rootless mode lets you run the Docker daemon and containers as a non-root
user to mitigate potential vulnerabilities in the daemon and
the container runtime.

Rootless mode does not require root privileges even during the installation of
the Docker daemon, as long as the [prerequisites](#prerequisites) are met.

## How it works

Rootless mode executes the Docker daemon and containers inside a user namespace.
This is similar to [`userns-remap` mode](../userns-remap.md), except that
with `userns-remap` mode, the daemon itself is running with root privileges,
whereas in rootless mode, both the daemon and the container are running without
root privileges.

Rootless mode does not use binaries with `SETUID` bits or file capabilities,
except `newuidmap` and `newgidmap`, which are needed to allow multiple
UIDs/GIDs to be used in the user namespace.


## Prerequisites

-  You must install `newuidmap` and `newgidmap` on the host. These commands
  are provided by the `uidmap` package on most distributions.

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

The `dockerd-rootless-setuptool.sh install` script (see following) automatically shows help
when the prerequisites are not satisfied.

## Install

> [!NOTE]
>
> If the system-wide Docker daemon is already running, consider disabling it:
>```console
>$ sudo systemctl disable --now docker.service docker.socket
>$ sudo rm /var/run/docker.sock
>```
> Should you choose not to shut down the `docker` service and socket, you will need to use the `--force`
> parameter in the next section. There are no known issues, but until you shutdown and disable you're
> still running rootful Docker. 

{{< tabs >}}
{{< tab name="With packages (RPM/DEB)" >}}

If you installed Docker 20.10 or later with [RPM/DEB packages](/engine/install), you should have `dockerd-rootless-setuptool.sh` in `/usr/bin`.

Run `dockerd-rootless-setuptool.sh install` as a non-root user to set up the daemon:

```console
$ dockerd-rootless-setuptool.sh install
[INFO] Creating /home/testuser/.config/systemd/user/docker.service
...
[INFO] Installed docker.service successfully.
[INFO] To control docker.service, run: `systemctl --user (start|stop|restart) docker.service`
[INFO] To run docker.service on system startup, run: `sudo loginctl enable-linger testuser`

[INFO] Creating CLI context "rootless"
Successfully created context "rootless"
[INFO] Using CLI context "rootless"
Current context is now "rootless"

[INFO] Make sure the following environment variable(s) are set (or add them to ~/.bashrc):
export PATH=/usr/bin:$PATH

[INFO] Some applications may require the following environment variable too:
export DOCKER_HOST=unix:///run/user/1000/docker.sock
```

If `dockerd-rootless-setuptool.sh` is not present, you may need to install the `docker-ce-rootless-extras` package manually, e.g.,

```console
$ sudo apt-get install -y docker-ce-rootless-extras
```

{{< /tab >}}
{{< tab name="Without packages" >}}

If you do not have permission to run package managers like `apt-get` and `dnf`,
consider using the installation script available at [https://get.docker.com/rootless](https://get.docker.com/rootless).
Since static packages are not available for `s390x`, hence it is not supported for `s390x`.

```console
$ curl -fsSL https://get.docker.com/rootless | sh
...
[INFO] Creating /home/testuser/.config/systemd/user/docker.service
...
[INFO] Installed docker.service successfully.
[INFO] To control docker.service, run: `systemctl --user (start|stop|restart) docker.service`
[INFO] To run docker.service on system startup, run: `sudo loginctl enable-linger testuser`

[INFO] Creating CLI context "rootless"
Successfully created context "rootless"
[INFO] Using CLI context "rootless"
Current context is now "rootless"

[INFO] Make sure the following environment variable(s) are set (or add them to ~/.bashrc):
export PATH=/home/testuser/bin:$PATH

[INFO] Some applications may require the following environment variable too:
export DOCKER_HOST=unix:///run/user/1000/docker.sock
```

The binaries will be installed at `~/bin`.

{{< /tab >}}
{{< /tabs >}}

Run `docker info` to confirm that the `docker` client is connecting to the Rootless daemon:
```console
$ docker info
Client: Docker Engine - Community
 Version:    28.3.3
 Context:    rootless
...
Server:
...
 Security Options:
  seccomp
   Profile: builtin
  rootless
  cgroupns
...
```

See [Troubleshooting](./troubleshoot.md) if you faced an error.