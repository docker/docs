---
title: How ECI works
description: Technical details of how Enhanced Container Isolation provides additional security for Docker Desktop
keywords: set up, enhanced container isolation, rootless, security
aliases:
 - /desktop/hardened-desktop/enhanced-container-isolation/how-eci-works/
 - /security/for-admins/hardened-desktop/enhanced-container-isolation/how-eci-works/
weight: 10
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Enhanced Container Isolation uses the [Sysbox
container runtime](https://github.com/nestybox/sysbox) to provide stronger container security. Sysbox is a specialized fork of the standard OCI runc runtime that enhances container isolation without impacting developer workflows.

When [Enhanced Container Isolation is turned on](index.md#how-do-i-enable-enhanced-container-isolation), containers
created through `docker run` or `docker create` automatically
use Sysbox instead of the standard runc runtime. Users can continue working with containers normally without any changes to their workflows.

Even containers using the `--privileged` flag run securely with Enhanced Container Isolation, preventing them from breaching the Docker Desktop virtual machine or other containers.

> [!NOTE]
>
> When Enhanced Container Isolation is turned on, the Docker CLI `--runtime` flag is ignored. Docker's default runtime remains `runc`, but all user containers implicitly launch with Sysbox.

## How Enhanced Container Isolation secures containers

Enhanced Container Isolation applies multiple security techniques simultaneously:

- Linux user namespaces: Root user in containers maps to unprivileged users in the Docker Desktop VM
- Sensitive directory restrictions: Containers can't mount sensitive VM directories
- System call filtering: Sensitive system calls between containers and the Linux kernel are inspected and restricted
- Filesystem user/group mapping: User and group IDs are safely mapped between container namespaces and the Linux VM
- Filesystem emulation: Portions of `/proc` and `/sys` filesystems are emulated inside containers for security
- Exclusive namespace mappings: Each container gets unique user-namespace mappings automatically
- Namespace isolation enforcement: Containers can't share namespaces with the Docker Desktop VM (`--network=host`, `--pid=host` are blocked)
- Configuration protection: Containers can't modify Docker Desktop VM configuration files
- Docker socket restrictions: Containers can't access the Docker Engine socket by default
- VM console restrictions: Direct console access to the Docker Desktop VM is blocked
- Privileged container containment: `--privileged` containers only have privileges within their own namespace
- Workflow compatibility: No changes required to existing development processes, tools, or container images
- Docker-in-Docker security: DinD and Kubernetes-in-Docker work but run unprivileged within the VM

These techniques use recent Linux kernel advances and complement Docker's existing security mechanisms including Linux namespaces, cgroups, restricted capabilities, seccomp, and AppArmor. Together, they create strong isolation between containers and the Linux kernel inside the Docker Desktop VM.

> [!IMPORTANT]
>
> ECI protection varies by Docker Desktop version. Later versions include more comprehensive protection. ECI doesn't yet protect extension containers.

## Enhanced Container Isolation versus user namespace remapping

The Docker Engine includes [userns-remap mode](/engine/security/userns-remap/)
that turns on user namespaces in all containers. However, this feature has several
[limitations](/engine/security/userns-remap/) and isn't supported in Docker Desktop.

Both userns-remap mode and Enhanced Container Isolation improve container isolation using Linux user namespaces, but Enhanced Container Isolation provides significant advantages:

- Automatic exclusive mappings: Each container gets unique user-namespace mappings without manual configuration
- Advanced isolation features: Includes system call filtering, filesystem emulation, and other security enhancements beyond basic user namespace remapping
- Docker Desktop optimization: Designed specifically for Docker Desktop's virtual machine architecture
- Enterprise security focus: Built for organizations with stringent security requirements

### Enhanced Container Isolation versus Rootless Docker

[Rootless Docker](/engine/security/rootless/) allows Docker Engine and containers to run without root privileges on Linux hosts. This lets non-root users install and run Docker natively on Linux systems.

Rootless Docker isn't supported in Docker Desktop because Docker Desktop already provides isolation through virtualization. Docker Desktop runs Docker Engine in a Linux VM, which already allows non-root host users to run Docker safely.

Key differences between Enhanced Container Isolation and Rootless Docker:

- Scope: Enhanced Container Isolation secures containers within the VM, while Rootless Docker secures the entire Docker Engine on the host
- Architecture: Enhanced Container Isolation works within Docker Desktop's VM architecture, while Rootless Docker modifies the host installation
- Limitations: Enhanced Container Isolation avoids the known limitations of Rootless Docker
- Security boundary: Enhanced Container Isolation creates stronger isolation between containers and the Docker Engine, while Rootless Docker focuses on isolating the Docker Engine from the host

Enhanced Container Isolation ensures that containers can't breach the Docker Desktop Linux VM or modify security settings within it, providing an additional security layer specifically designed for Docker Desktop's architecture.
