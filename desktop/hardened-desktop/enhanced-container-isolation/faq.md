---
title: FAQs and known issues
description: FAQ for Enhanced Container Isolation
keywords: enhanced container isolation, security, faq, sysbox, known issues, Docker Desktop
toc_max: 2
---

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab3">FAQs</a></li>
  <li><a data-toggle="tab" data-target="#tab4">Limitations and Known Issues</a></li>
</ul>
<div class="tab-content">
<div id="tab3" class="tab-pane fade in active" markdown="1">

#### Do I need to change the way I use Docker when Enhanced Container Isolation is enabled?

No, you can continue to use Docker as usual. Enhanced Container Isolation will be mostly transparent to you.

#### Do all container workloads work well with Enhanced Container Isolation?

The great majority of container workloads run fine with ECI, but a few do not
(yet). For the few workloads that don't yet work with Enhanced Container
Isolation, Docker will continue to improve the feature to reduce this to a
minimum.

#### Can I run privileged containers with Enhanced Container Isolation?

Yes, you can use the `--privileged` flag in containers but unlike privileged
containers without Enhanced Container Isolation, the container can only use it's elevated privileges to
access resources assigned to the container. It can't access global kernel
resources in the Docker Desktop Linux VM. This allows you to run privileged
containers securely (including Docker-in-Docker). For more information, see [Key features and benefits](features-benefits.md#privileged-containers-are-also-secured).

#### Will all privileged container workloads run with Enhanced Container Isolation?

No. Privileged container workloads that wish to access global kernel resources
inside the Docker Desktop Linux VM won't work. For example, you can't use a
privileged container to load a kernel module.

#### Why not just restrict usage of the `--privileged` flag?

Privileged containers are typically used to run advanced workloads in
containers, for example Docker-in-Docker or Kubernetes-in-Docker, to
perform kernel operations such as loading modules, or to access hardware
devices.

Enhanced Container Isolation allows running advanced workloads, but denies the ability to perform
kernel operations or access hardware devices.

#### Does Enhanced Container Isolation restrict bind mounts inside the container?

Yes, it restricts bind mounts of directories located in the Docker Desktop Linux
VM into the container.

It does not restrict bind mounts of your host machine files into the container,
as configured via Docker Desktop's **Settings** > **Resources** > **File Sharing**.

#### Does Enhanced Container Isolation protect all containers launched with Docker Desktop?

It protects all containers launched by users via `docker create` and `docker run`. It does not yet protect Docker Desktop Kubernetes pods, Extension
Containers, and Dev Environments.

#### Does Enhanced Container Isolation protect container launched prior to enabling ECI?

No. Containers created prior to enabling ECI are not protected. Therefore, we
recommend removing all containers prior to enabling ECI. In the future Docker
Desktop will likely make this a hard requirement.

#### Does Enhanced Container Isolation affect performance of containers?

Enhanced Container Isolation has very little impact on the performance of
containers. The exception is for containers that perform lots of `mount` and
`umount` system calls, as these are trapped and vetted by the Sysbox container
runtime to ensure they are not being used to breach the container's filesystem.

#### With Enhanced Container Isolation, can the user still override the `--runtime` flag from the CLI ?

No. With Enhanced Container Isolation enabled, Sysbox is set as the default (and only) runtime for
containers deployed by Docker Desktop users. If a user attempts to override the
runtime (e.g., `docker run --runtime=runc`), this request is ignored and the
container is created through the Sysbox runtime.

The reason `runc` is disallowed with Enhanced Container Isolation because it
allows users to run as "true root" on the Docker Desktop Linux VM, thereby
providing them with implicit control of the VM and the ability to modify the
administrative configurations for Docker Desktop, for example.

#### How is ECI different from Docker Engine's userns-remap mode?

See [How does it work](how-eci-works.md#enhanced-container-isolation-vs-docker-userns-remap-mode).

#### How is ECI different from Rootless Docker?

See [How does it work](how-eci-works.md#enhanced-container-isolation-vs-rootless-docker)

<hr>
</div>
<div id="tab4" class="tab-pane fade" markdown="1">

#### ECI support for WSL

Prior to Docker Desktop 4.20, Enhanced Container Isolation (ECI) on
Windows hosts was only supported when Docker Desktop was configured to use
Hyper-V to create the Docker Desktop Linux VM. ECI was not supported when Docker
Desktop was configured to use Windows Subsystem for Linux (aka WSL).

Starting with Docker Desktop 4.20, ECI is supported when Docker Desktop is
configured to use either Hyper-V or WSL version 2.

>**Note**
>
> Docker Desktop requires WSL 2 version 1.1.3.0 or later. To get the current
> version of WSL on your host, type `wsl --version`. If the command fails or if
> it returns a version number prior to 1.1.3.0, update WSL to the latest version
> by typing `wsl --update` in a Windows command or PowerShell terminal.

Note however that ECI on WSL is not as secure as on Hyper-V because:

* While ECI on WSL still hardens containers so that malicious workloads can't
  easily breach Docker Desktop's Linux VM, ECI on WSL can't prevent Docker
  Desktop users from breaching the Docker Desktop Linux VM. Such users can
  trivially access that VM (as root) with the `wsl -d docker-desktop` command,
  and use that access to modify Docker Engine settings inside the VM. This gives
  Docker Desktop users control of the Docker Desktop VM and allows them to
  bypass Docker Desktop configs set by admins via the
  [settings-management](../settings-management/index.md) feature. In contrast,
  ECI on Hyper-V does not allow Docker Desktop users to breach the Docker
  Desktop Linux VM.

* With WSL 2, all WSL 2 distros on the same Windows host share the same instance
  of the Linux kernel. As a result, Docker Desktop can't ensure the integrity of
  the kernel in the Docker Desktop Linux VM since another WSL 2 distro could
  modify shared kernel settings. In contrast, when using Hyper-V, the Docker
  Desktop Linux VM has a dedicated kernel that is solely under the control of
  Docker Desktop.

The table below summarizes this.

| Security Feature                                   | ECI on WSL   | ECI on Hyper-V   | Comment               |
| -------------------------------------------------- | ------------ | ---------------- | --------------------- |
| Strongly secure containers                         | Yes          | Yes              | Makes it harder for malicious container workloads to breach the Docker Desktop Linux VM and host. |
| Docker Desktop Linux VM protected from user access | No           | Yes              | On WSL, users can access Docker Engine directly or bypass Docker Desktop security settings. |
| Docker Desktop Linux VM has a dedicated kernel     | No           | Yes              | On WSL, Docker Desktop can't guarantee the integrity of kernel level configs. |

In general, using ECI with Hyper-V is more secure than with WSL 2. But WSL 2
offers advantages for performance and resource utilization on the host machine,
and it's an excellent way for users to run their favorite Linux distro on
Windows hosts and access Docker from within (see Docker Desktop's WSL distro
integration feature, enabled via the Dashboard's **Settings** > **Resources** > **WSL Integration**).

#### Docker build and buildx has some restrictions
With ECI enabled, Docker build `--network=host` and Docker buildx entitlements
(`network.host`, `security.insecure`) are not allowed. Builds that require
these will not work properly.

#### Kubernetes pods are not yet protected
Kubernetes pods are not yet protected by ECI. A malicious or privileged pod can
compromise the Docker Desktop Linux VM and bypass security controls. We expect
to improve on this in future versions of Docker Desktop.

#### Extension Containers are not yet protected
Extension containers are also not yet protected by ECI. Ensure you extension
containers come from trusted entities to avoid issues. We expect to improve on
this in future versions of Docker Desktop.

#### Docker Desktop dev environments are not yet protected
Containers launched by the Docker Desktop Dev Environments feature are not yet
protected either. We expect to improve on this in future versions of Docker
Desktop.

#### Use in production
In general users should not experience differences between running a container
in Docker Desktop with ECI enabled, which uses the Sysbox runtime, and running
that same container in production, through the standard OCI `runc` runtime.

However in some cases, typically when running advanced or privileged workloads in
containers, users may experience some differences. In particular, the container
may run with ECI but not with `runc`, or vice-versa.

<hr>
</div>
</div>
