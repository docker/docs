---
title: FAQs and known issues
description: FAQ for Enhanced Container Isolation 
keywords: enhanced container isolation, security, faq, sysbox
toc_max: 2
---

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" data-target="#tab3">FAQs</a></li>
  <li><a data-toggle="tab" data-target="#tab4">Known issues</a></li>
</ul>
<div class="tab-content">
<div id="tab3" class="tab-pane fade in active" markdown="1">

### With Enhanced Container Isolation enabled, can the user still override the `--runtime` flag from the CLI ?

No. With Hardened Desktop enabled, Sysbox is locked as the default (and only) runtime. If a user attempts to override the runtime by launching a container with the standard `runc` runtime, for example `docker run --runtime=runc`, this request is ignored and the container is created through the Sysbox runtime. 

The reason `runc` is disallowed with Enhanced Container Isolation is because it allows users to run as root on the Docker Desktop Linux VM, thereby providing them with implicit control of the VM and the ability to modify the administrative configurations for Docker Desktop, for example.

### With Enhanced Container Isolation enabled, can the user still use the `--privileged` flag from the CLI?

Yes, with Enhanced Container Isolation the container is only privileged within its assigned Linux user-namespace. It is not privileged within the Docker Desktop Linux VM.

For example, the container’s init process will have all Linux capabilities enabled, have read/write access to the kernel’s `/proc` and `/sys`, run without system call or other restrictions normally imposed by Docker on regular containers (for example, seccomp, AppArmor), and see all host devices under the container’s `/dev` directory.

However, because Sysbox launches each container within a dedicated Linux user-namespace and vets sensitive accesses to the kernel, the container can only access resources assigned to it. For example, the container can’t access resources under `/proc` and `/sys` that are not namespaced. Although it can see all host devices under `/dev`, it won’t have permission to access them. Also, while the container can use system calls such as “mount” and “umount”, Sysbox prevents the container from using them to modify the container’s chroot jail.

This makes running a privileged container with Enhanced Container Isolation much safer than a privileged container launched with the standard runc, which offers almost no isolation.

### Why not just restrict usage of the `--privileged` flag with Enhanced Container Isolation?

Privileged containers are typically used to run advanced workloads in containers, for example Docker-in-Docker, to perform kernel operations such as loading modules, or to access hardware devices. We aim to allow running advanced workloads, but deny the ability to perform kernel operations or access hardware devices.

By allowing the `-–privileged` flag but restricting its impact within the container's user-namespace, it’s possible to do this.

<hr>
</div>
<div id="tab4" class="tab-pane fade" markdown="1">

#### Incompatibility with WSL
Enhanced Container Isolation does not currently work when Docker Desktop runs on Windows with WSL/WSL2. This is due to some limitations of the WSL/WSL2 Linux Kernel. As a result, to use Enhanced Container Isolation on Windows, you must configure Docker Desktop to use Hyper-V. This can be enforced using Admin Controls. For more information, see [Settings Management](../settings-management/index.md).

#### Kubernetes pods and extension containers are not yet protected
When Enhanced Container Isolation is enabled, Kubernetes pods and extension containers are not yet protected. A malicious or privileged pod or extension container can compromise the Docker Desktop Linux VM and bypass security controls. 

#### Use in production
Users may experience some differences between running a container in Docker Desktop with Enhanced Container Isolation enabled, and running that same container in production. This is because in production the container may run on another runtime, typically the OCI runc.

<hr>
</div>
</div>
