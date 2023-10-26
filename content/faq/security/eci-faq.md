---
title: FAQs
description: Frequently asked questions for Enhanced Container Isolation
keywords: enhanced container isolation, security, faq, sysbox, Docker Desktop
toc_max: 2
aliases:
- /desktop/hardened-desktop/enhanced-container-isolation/faq/
---

### Do I need to change the way I use Docker when Enhanced Container Isolation is switched on?

No, you can continue to use Docker as usual. 

### Do all container workloads work well with Enhanced Container Isolation?

The great majority of container workloads run fine with ECI, but a few do not
(yet). For the few workloads that don't yet work with Enhanced Container
Isolation, Docker is continuing to improve the feature to reduce this to a
minimum.

### Can I run privileged containers with Enhanced Container Isolation?

Yes, you can use the `--privileged` flag in containers but unlike privileged
containers without Enhanced Container Isolation, the container can only use it's elevated privileges to
access resources assigned to the container. It can't access global kernel
resources in the Docker Desktop Linux VM. This allows you to run privileged
containers securely (including Docker-in-Docker). For more information, see [Key features and benefits](features-benefits.md#privileged-containers-are-also-secured).

### Will all privileged container workloads run with Enhanced Container Isolation?

No. Privileged container workloads that wish to access global kernel resources
inside the Docker Desktop Linux VM won't work. For example, you can't use a
privileged container to load a kernel module.

### Why not just restrict usage of the `--privileged` flag?

Privileged containers are typically used to run advanced workloads in
containers, for example Docker-in-Docker or Kubernetes-in-Docker, to
perform kernel operations such as loading modules, or to access hardware
devices.

Enhanced Container Isolation allows the running of advanced workloads, but denies the ability to perform
kernel operations or access hardware devices.

### Does Enhanced Container Isolation restrict bind mounts inside the container?

Yes, it restricts bind mounts of directories located in the Docker Desktop Linux
VM into the container.

It doesn't restrict bind mounts of your host machine files into the container,
as configured via Docker Desktop's **Settings** > **Resources** > **File Sharing**.

### Does Enhanced Container Isolation protect all containers launched with Docker Desktop?

It protects all containers launched by users via `docker create` and `docker run`. It does not yet protect Docker Desktop Kubernetes pods, ExtensioncContainers, and Dev Environments.

### Does Enhanced Container Isolation protect containers launched prior to enabling ECI?

No. Containers created prior to switching on ECI are not protected. Therefore, we
recommend removing all containers prior to switching on ECI. 

### Does Enhanced Container Isolation affect the performance of containers?

Enhanced Container Isolation has very little impact on the performance of
containers. The exception is for containers that perform lots of `mount` and
`umount` system calls, as these are trapped and vetted by the Sysbox container
runtime to ensure they are not being used to breach the container's filesystem.

### With Enhanced Container Isolation, can the user still override the `--runtime` flag from the CLI ?

No. With Enhanced Container Isolation enabled, Sysbox is set as the default (and only) runtime for
containers deployed by Docker Desktop users. If a user attempts to override the
runtime (e.g., `docker run --runtime=runc`), this request is ignored and the
container is created through the Sysbox runtime.

The reason `runc` is disallowed with Enhanced Container Isolation because it
allows users to run as "true root" on the Docker Desktop Linux VM, thereby
providing them with implicit control of the VM and the ability to modify the
administrative configurations for Docker Desktop, for example.

### How is ECI different from Docker Engine's userns-remap mode?

See [How does it work](how-eci-works.md#enhanced-container-isolation-vs-docker-userns-remap-mode).

### How is ECI different from Rootless Docker?

See [How does it work](how-eci-works.md#enhanced-container-isolation-vs-rootless-docker)