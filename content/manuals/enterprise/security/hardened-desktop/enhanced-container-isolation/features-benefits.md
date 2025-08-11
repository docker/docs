---
title: Features and benefits of ECI
linkTitle: Features and benefits
description: Key security features and benefits of Enhanced Container Isolation for Docker Desktop
keywords: set up, enhanced container isolation, rootless, security, features, Docker Desktop
aliases:
 - /desktop/hardened-desktop/enhanced-container-isolation/features-benefits/
 - /security/for-admins/hardened-desktop/enhanced-container-isolation/features-benefits/
weight: 20
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Enhanced Container Isolation provides multiple layers of security protection that work together to prevent container breaches while maintaining full developer productivity. This page explains the key security features and their benefits.

## Linux user namespace on all containers

Enhanced Container Isolation automatically applies [Linux user namespaces](https://man7.org/linux/man-pages/man7/user_namespaces.7.html)
to every container. The root user inside containers maps to unprivileged users in the Docker Desktop Linux VM, preventing privilege escalation attacks.

For example:

```console
$ docker run -it --rm --name=first alpine
/ # cat /proc/self/uid_map
         0     100000      65536
```

The output `0 100000 65536` shows that:

- Container root user (0) maps to unprivileged user 100000 in the VM
- The mapping covers a continuous range of 64K user IDs
- The same principle applies to group IDs

Each container gets an exclusive user ID range managed automatically by Sysbox:

```console
$ docker run -it --rm --name=second alpine
/ # cat /proc/self/uid_map
         0     165536      65536
```

In contrast, without Enhanced Container Isolation, the container's root user is
in fact root on the host (aka "true root") and this applies to all containers:

```console
$ docker run -it --rm alpine
/ # cat /proc/self/uid_map
         0       0     4294967295
```

This shows the container root user maps directly to VM root user (0), providing much less isolation.

## Security benefits

Enhanced Container Isolation provides the following security benefits:

- Container processes never run with true root privileges in the Linux VM
- Processes don't run with any valid user ID in the Linux VM
- Linux capabilities are constrained to resources within the container only
- Significantly improved container-to-host and cross-container isolation

## Secured privileged containers

Privileged containers (`docker run --privileged`) normally pose significant security risks because they provide unrestricted access to the Linux kernel. Without Enhanced Container Isolation, privileged containers can:

- Run as true root with all capabilities
- Bypass seccomp and AppArmor restrictions
- Access all hardware devices
- Modify global kernel settings

Organizations securing developer environments face challenges with privileged containers because they can gain control of the Docker Desktop VM and alter security settings like registry access management and network proxies.

Enhanced Container Isolation ensures privileged containers can only access resources within their container boundary. The Linux user namespace and Sysbox security techniques prevent privileged containers from breaching the Docker Desktop VM.

> [!NOTE]
>
> Enhanced Container Isolation doesn't prevent users from running privileged containers, but makes them secure by containing their access to container resources only. Privileged workloads that modify global kernel settings (loading kernel modules, changing Berkeley Packet Filter settings) receive "permission denied" errors.

For example, Enhanced Container Isolation ensures privileged containers can't
access Docker Desktop network settings in the Linux VM configured via BPF:

```console
$ docker run --privileged djs55/bpftool map show
Error: can't get next map: Operation not permitted
```

Without Enhanced Container Isolation, privileged containers can access and modify these settings:

```console
$ docker run --privileged djs55/bpftool map show
17: ringbuf  name blocked_packets  flags 0x0
        key 0B  value 0B  max_entries 16777216  memlock 0B
18: hash  name allowed_map  flags 0x0
        key 4B  value 4B  max_entries 10000  memlock 81920B
20: lpm_trie  name allowed_trie  flags 0x1
        key 8B  value 8B  max_entries 1024  memlock 16384B
```

Advanced container workloads like Docker-in-Docker and Kubernetes-in-Docker still work with Enhanced Container Isolation but run much more securely than before.

## Namespace isolation enforcement

Enhanced Container Isolation prevents containers from sharing Linux namespaces with the Docker Desktop VM, which would break isolation boundaries.

For example, sharing the PID namespace fails:

```console
$ docker run -it --rm --pid=host alpine
docker: Error response from daemon: failed to create shim task: OCI runtime create failed: error in the container spec: invalid or unsupported container spec: sysbox containers can't share namespaces [pid] with the host (because they use the linux user-namespace for isolation): unknown.
```

Similarly sharing the network namespace fails:

```console
$ docker run -it --rm --network=host alpine
docker: Error response from daemon: failed to create shim task: OCI runtime create failed: error in the container spec: invalid or unsupported container spec: sysbox containers can't share a network namespace with the host (because they use the linux user-namespace for isolation): unknown.
```

In addition, the `--userns=host` flag, used to disable the user namespace on the
container, is ignored:

```console
$ docker run -it --rm --userns=host alpine
/ # cat /proc/self/uid_map
         0     100000      65536
```

Finally, Docker build `--network=host` and Docker buildx entitlements
(`network.host`, `security.insecure`) are not allowed. Builds that require these
won't work properly.

## Protected bind mounts

Enhanced Container Isolation maintains support for standard file sharing while preventing access to sensitive VM directories.

### Allowed: Host directory bind mounts

Users can continue bind mounting host directories configured in **Settings** > **Resources** > **File sharing**:

```console
$ docker run -it --rm -v $HOME:/mnt alpine
/ # ls /mnt
# Lists home directory contents successfully
```

### Blocked: VM configuration bind mounts

Containers can't mount sensitive Docker Desktop VM directories that contain security configurations:

```
$ docker run -it --rm -v /etc/docker/daemon.json:/mnt/daemon.json alpine
docker: Error response from daemon: failed to create shim task: OCI runtime create failed: error in the container spec: can't mount /etc/docker/daemon.json because it's configured as a restricted host mount: unknown
```

This prevents containers from reading or modifying:

- Docker Engine configurations
- Registry access management settings
- Proxy configurations
- Other security-related VM files

> [!NOTE]
>
> By default, Enhanced Container Isolation blocks bind mounting the Docker Engine socket (`/var/run/docker.sock`) as this would grant containers control over Docker Engine. Administrators can create exceptions for trusted container images. See [Docker socket mount permissions](config.md#docker-socket-mount-permissions).

## Advanced system call filtering

Enhanced Container Isolation intercepts and inspects sensitive system calls between containers and the Linux kernel, preventing containers from using legitimate capabilities maliciously.

Even containers with CAP_SYS_ADMIN capability can't use mount operations to breach container boundaries:

```console
$ docker run -it --rm --cap-add SYS_ADMIN -v $HOME:/mnt:ro alpine
/ # mount -o remount,rw /mnt /mnt
mount: permission denied (are you root?)
```

The read-only bind mount can't be changed to read-write from within the container, even with the appropriate capability.

### Allowed: Internal container mounts

Containers can still create and modify mounts within their own filesystem:

```console
/ # mkdir /root/tmpfs
/ # mount -t tmpfs tmpfs /root/tmpfs
/ # mount -o remount,ro /root/tmpfs /root/tmpfs
/ # findmnt | grep tmpfs
├─/root/tmpfs    tmpfs      tmpfs    ro,relatime,uid=100000,gid=100000

/ # mount -o remount,rw /root/tmpfs /root/tmpfs
/ # findmnt | grep tmpfs
├─/root/tmpfs    tmpfs      tmpfs    rw,relatime,uid=100000,gid=100000
```

### Performance optimization

Enhanced Container Isolation performs system call filtering efficiently by:

- Intercepting only control-path system calls (rarely used in most workloads)
- Avoiding interception of data-path system calls (frequent operations)
- Maintaining container performance in the majority of use cases

This approach ensures that containers with all Linux capabilities still can't breach container boundaries while maintaining good performance.

## Automatic filesystem user ID mapping

Enhanced Container Isolation solves the challenge of sharing files between containers that use different user ID ranges through automatic filesystem mapping.

The challenge is, each container gets exclusive user ID mappings (container 1: 100000-165535, container 2: 165536-231071). Without mapping, files created by one container would be inaccessible to others due to different user IDs.

To solve this, Sysbox uses filesystem user ID remapping through:

- ID-mapped mounts: Linux kernel feature (added in 2021) that maps filesystem accesses
- Alternative shiftsfs: Fallback module for older kernel versions

### How it works

Filesystem accesses from containers' real user IDs (e.g., 100000-165535) are mapped to the standard range (0-65535) in the Docker Desktop VM. This allows:

- Volume sharing across containers with different user ID ranges
- Consistent file ownership regardless of container user ID mappings
- Transparent file access without user intervention

Even though filesystem mapping allows containers to access VM files with user ID 0, the bind mount restrictions prevent containers from mounting sensitive VM directories.

## Procfs and sysfs emulation

Enhanced Container Isolation emulates portions of the `/proc` and `/sys` filesystems within containers to hide sensitive host information and provide per-container views of kernel resources.

The `/proc/uptime` file shows container uptime instead of VM uptime:

```console
$ docker run -it --rm alpine
/ # cat /proc/uptime
5.86 5.86
```

Without Enhanced Container Isolation, this would show Docker Desktop VM uptime, potentially
leaking system information.

Several `/proc/sys resources` that aren't namespaced by the Linux kernel are emulated per-container. Each container sees a separate view of these resources, and Sysbox coordinates values when programming kernel settings.

## Performance and compatibility

Enhanced Container Isolation is designed to provide strong security without impacting developer productivity:

- No performance impact: System call filtering targets only control-path calls, leaving data-path operations unaffected
- Full workflow compatibility: Existing development processes, tools, and container images work unchanged
- Advanced workload support: Docker-in-Docker, Kubernetes-in-Docker, and other complex scenarios work securely
- Automatic management: User ID mappings, filesystem access, and security policies are handled automatically
