---
title: Enhanced Container Isolation
linkTitle: Enhanced Container Isolation
description: Enhanced Container Isolation provides additional security for Docker Desktop by preventing malicious containers from compromising the host
keywords: enhanced container isolation, container security, sysbox runtime, linux user namespaces, hardened desktop
aliases:
 - /desktop/hardened-desktop/enhanced-container-isolation/
 - /security/for-admins/hardened-desktop/enhanced-container-isolation/
 - /security/hardened-desktop/enhanced-container-isolation/how-eci-works
 - /security/hardened-desktop/enhanced-container-isolation/features-benefits
weight: 10
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Enhanced Container Isolation (ECI) prevents malicious containers from compromising Docker Desktop or the host system. It applies advanced security techniques automatically while maintaining full developer productivity and workflow compatibility.

ECI strengthens container isolation and locks in security configurations created by administrators, such as [Registry Access Management policies](/manuals/enterprise/security/hardened-desktop/registry-access-management.md) and [Settings Management](../settings-management/_index.md) controls.

> [!NOTE]
>
> ECI works alongside other Docker security features like reduced Linux capabilities, seccomp, and AppArmor.

## Who should use Enhanced Container Isolation?

Enhanced Container Isolation is designed for:

- Organizations that want to prevent container-based attacks and reduce security vulnerabilities in developer environments
- Security teams that need stronger container isolation without impacting developer workflows
- Enterprises that require additional protection when running untrusted or third-party container images

## How Enhanced Container Isolation works

Docker implements ECI using the [Sysbox container runtime](https://github.com/nestybox/sysbox), a
security-enhanced fork of the standard OCI runc runtime. When ECI is turned on, containers created through `docker run` or `docker create` automatically use Sysbox instead of runc without requiring any changes to developer workflows.

> [!NOTE]
>
> When ECI is turned on, the Docker CLI `--runtime` flag is ignored.
Docker's default runtime remains `runc`, but all user containers
implicitly launch with Sysbox.

Even containers using the `--privileged` flag run securely with Enhanced Container Isolation, preventing them from breaching the Docker Desktop virtual machine or other containers.

## Key security features

### Linux user namespace isolation

With Enhanced Container Isolation, all containers leverage Linux user namespaces for stronger isolation. Container root users map to unprivileged users in the Docker Desktop VM:

```console
$ docker run -it --rm --name=first alpine
/ # cat /proc/self/uid_map
         0     100000      65536
```

This output shows that container root (0) maps to unprivileged user 100000 in the VM, with a range of 64K user IDs. Each container gets exclusive mappings:

```console
$ docker run -it --rm --name=second alpine
/ # cat /proc/self/uid_map
         0     165536      65536
```

Without Enhanced Container Isolation, containers run as true root:

```console
$ docker run -it --rm alpine
/ # cat /proc/self/uid_map
         0       0     4294967295
```

By using Linux user namespaces, ECI ensures container processes never run with valid user IDs in the Linux VM, constraining their capabilities to resources within the container only.

### Secured privileged containers

Privileged containers (`docker run --privileged`) normally pose significant security risks because they provide unrestricted access to the Linux kernel. Without ECI, privileged containers can:

- Run as true root with all capabilities
- Bypass seccomp and AppArmor restrictions
- Access all hardware devices
- Modify global kernel settings

Organizations securing developer environments face challenges with privileged containers because they can gain control of the Docker Desktop VM and alter security settings like registry access management and network proxies.

Enhanced Container Isolation transforms privileged containers by ensuring they can only access resources within their container boundary. For example, privileged containers can't access Docker Desktop's network configuration:

```console
$ docker run --privileged djs55/bpftool map show
Error: can't get next map: Operation not permitted
```

Without ECI, privileged containers can easily access and modify these settings:

```console
$ docker run --privileged djs55/bpftool map show
17: ringbuf  name blocked_packets  flags 0x0
        key 0B  value 0B  max_entries 16777216  memlock 0B
18: hash  name allowed_map  flags 0x0
        key 4B  value 4B  max_entries 10000  memlock 81920B
```

Advanced container workloads like Docker-in-Docker and Kubernetes-in-Docker still work with ECI but run much more securely.

> [!NOTE]
>
> ECI doesn't prevent users from running privileged containers, but makes them secure by containing their access. Privileged workloads that modify global kernel settings (loading kernel modules, changing Berkeley Packet Filter settings) receive "permission denied" errors.

### Namespace isolation enforcement

Enhanced Container Isolation prevents containers from sharing Linux namespaces with the Docker Desktop VM, maintaining isolation boundaries:

**PID namespace sharing blocked:**

```console
$ docker run -it --rm --pid=host alpine
docker: Error response from daemon: failed to create shim task: OCI runtime create failed: error in the container spec: invalid or unsupported container spec: sysbox containers can't share namespaces [pid] with the host (because they use the linux user-namespace for isolation): unknown.
```

**Network namespace sharing blocked:**

```console
$ docker run -it --rm --network=host alpine
docker: Error response from daemon: failed to create shim task: OCI runtime create failed: error in the container spec: invalid or unsupported container spec: sysbox containers can't share a network namespace with the host (because they use the linux user-namespace for isolation): unknown.
```

**User namespace override ignored:**

```console
$ docker run -it --rm --userns=host alpine
/ # cat /proc/self/uid_map
         0     100000      65536
```

Docker build operations using `--network-host` and Docker buildx entitlements (`network.host`,
`security.insecure`) are also blocked.

### Protected bind mounts

Enhanced Container Isolation maintains support for standard file sharing while preventing access to sensitive VM directories:

**Host directory mounts continue to work:**

```console
$ docker run -it --rm -v $HOME:/mnt alpine
/ # ls /mnt
# Successfully lists home directory contents
```

**VM configuration mounts are blocked:**

```console
$ docker run -it --rm -v /etc/docker/daemon.json:/mnt/daemon.json alpine
docker: Error response from daemon: failed to create shim task: OCI runtime create failed: error in the container spec: can't mount /etc/docker/daemon.json because it's configured as a restricted host mount: unknown
```

This prevents containers from reading or modifying Docker Engine configurations, registry access management settings, proxy configurations, and other security-related VM files.

> [!NOTE]
>
> By default, ECI blocks bind mounting the Docker Engine socket (/var/run/docker.sock) as this would grant containers control over Docker Engine. Administrators can create exceptions for trusted container images.

### Advanced system call protection

Enhanced Container Isolation intercepts sensitive system calls to prevent containers from using legitimate capabilities maliciously:

```console
$ docker run -it --rm --cap-add SYS_ADMIN -v $HOME:/mnt:ro alpine
/ # mount -o remount,rw /mnt /mnt
mount: permission denied (are you root?)
```

Even with `CAP_SYS_ADMIN` capability, containers can't change read-only bind mounts to read-write, ensuring they can't breach container boundaries.

Containers can still create internal mounts within their filesystem:

```console
/ # mkdir /root/tmpfs
/ # mount -t tmpfs tmpfs /root/tmpfs
/ # mount -o remount,ro /root/tmpfs /root/tmpfs
/ # findmnt | grep tmpfs
├─/root/tmpfs    tmpfs      tmpfs    ro,relatime,uid=100000,gid=100000
```

ECI performs system call filtering efficiently by intercepting only control-path system calls (rarely used) while leaving data-path system calls unaffected, maintaining container performance.

### Automatic filesystem user ID mapping

Enhanced Container Isolation solves file sharing challenges between containers with different user ID ranges through automatic filesystem mapping.

Each container gets exclusive user ID mappings, but Sysbox uses filesystem user ID remapping via Linux kernel ID-mapped mounts (added in 2021) or alternative shiftsfs module. This maps filesystem accesses from containers' real user IDs to standard ranges, enabling:

- Volume sharing across containers with different user ID ranges
- Consistent file ownership regardless of container user ID mappings
- Transparent file access without user intervention

### Information hiding through filesystem emulation

ECI emulates portions of `/proc` and `/sys` filesystems within containers to hide sensitive host information and provide per-container views of kernel resources:

```console
$ docker run -it --rm alpine
/ # cat /proc/uptime
5.86 5.86
```

This shows container uptime instead of Docker Desktop VM uptime, preventing system information from leaking into containers.

Several `/proc/sys` resources that aren't namespaced by the Linux kernel are emulated per-container, with Sysbox coordinating values when programming kernel settings. This enables container workloads that normally require privileged access to run securely.

## Performance and compatibility

Enhanced Container Isolation maintains excellent performance and full compatibility:

- No performance impact: System call filtering targets only control-path calls, leaving data-path operations unaffected
- Full workflow compatibility: Existing development processes, tools, and container images work unchanged
- Advanced workload support: Docker-in-Docker, Kubernetes-in-Docker, and other complex scenarios work securely
- Automatic management: User ID mappings, filesystem access, and security policies are handled automatically
- Standard image support: No special container images or modifications required

> [!IMPORTANT]
>
> ECI protection varies by Docker Desktop version and doesn't yet protect extension containers. Docker builds and Kubernetes in Docker Desktop have varying protection levels depending on the version. For details, see [Enhanced Container Isolation limitations](/manuals/enterprise/security/hardened-desktop/enhanced-container-isolation/limitations.md).
